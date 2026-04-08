package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/coliva/tsk/internal/model"
	"github.com/coliva/tsk/internal/storage"
	"github.com/coliva/tsk/internal/styles"
	"github.com/coliva/tsk/internal/undo"
	"github.com/google/uuid"
)

// Messages
type (
	tickMsg        time.Time
	statusClearMsg struct{}
)

// App is the main application model
type App struct {
	state   *model.AppState
	storage *storage.Storage

	// Text input for various modes
	textInput textinput.Model

	// Undo/redo manager
	undoManager *undo.Manager

	// Task detail view
	taskDetail *TaskDetail

	// Task edit view
	taskEdit *TaskEdit

	// Label editor
	labelEditor *LabelEditor

	// Due date editor
	dueDateEditor *DueDateEditor

	// Board selector
	boardSelector *BoardSelector

	// Search component
	search *Search

	// Filter component
	filter *Filter

	// Auto-save ticker
	lastAutoSave time.Time

	// Error state
	err error
}

// NewApp creates a new application
func NewApp(store *storage.Storage) *App {
	ti := textinput.New()
	ti.Placeholder = "Enter text..."
	ti.CharLimit = 256

	return &App{
		state:        model.NewAppState(),
		storage:      store,
		textInput:    ti,
		undoManager:  undo.NewManager(20),
		lastAutoSave: time.Now(),
	}
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	return tea.Batch(
		a.loadInitialBoard(),
		a.tickCmd(),
	)
}

func (a *App) loadInitialBoard() tea.Cmd {
	return func() tea.Msg {
		hasBoards, err := a.storage.HasBoards()
		if err != nil {
			return err
		}

		if !hasBoards {
			// First run - show welcome screen
			a.state.Mode = model.ModeWelcome
			a.textInput.SetValue("")
			a.textInput.Placeholder = "My Tasks"
			a.textInput.Focus()
			return nil
		}

		// Load most recent board
		board, err := a.storage.MostRecentBoard()
		if err != nil {
			return err
		}
		a.state.SetBoard(board)
		return nil
	}
}

func (a *App) tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update handles messages
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case error:
		a.err = msg
		a.state.SetStatusMessage(fmt.Sprintf("Error: %v", msg))
		return a, nil

	case tea.KeyMsg:
		cmd := a.handleKeyMsg(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case tea.WindowSizeMsg:
		a.state.Width = msg.Width
		a.state.Height = msg.Height

	case SearchDebounceMsg:
		// Perform the search after debounce
		if a.search != nil {
			a.search.PerformSearch(msg.Query)
		}

	case tickMsg:
		// Clear old status messages (3 seconds)
		a.state.ClearStatusMessageIfOld(3 * time.Second)

		// Auto-save if dirty (5 seconds)
		if a.state.Dirty && time.Since(a.lastAutoSave) > 5*time.Second {
			if err := a.save(); err == nil {
				a.lastAutoSave = time.Now()
			}
		}

		cmds = append(cmds, a.tickCmd())

	case statusClearMsg:
		a.state.StatusMessage = ""
	}

	return a, tea.Batch(cmds...)
}

func (a *App) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	// Handle board mode
	if a.state.Mode == model.ModeBoard {
		return a.handleBoardMode(msg)
	}

	// Handle due date mode
	if a.state.Mode == model.ModeDueDate {
		return a.handleDueDateMode(msg)
	}

	// Handle labels mode
	if a.state.Mode == model.ModeLabels {
		return a.handleLabelsMode(msg)
	}

	// Handle edit mode
	if a.state.Mode == model.ModeEdit {
		return a.handleEditMode(msg)
	}

	// Handle detail mode
	if a.state.Mode == model.ModeDetail {
		return a.handleDetailMode(msg)
	}

	// Handle modal mode
	if a.state.Mode == model.ModeModal {
		return a.handleModalMode(msg)
	}

	// Handle search mode
	if a.state.Mode == model.ModeSearch {
		return a.handleSearchMode(msg)
	}

	// Handle filter mode
	if a.state.Mode == model.ModeFilter {
		return a.handleFilterMode(msg)
	}

	// Handle mode-specific input
	if a.state.Mode.IsTextInput() {
		return a.handleTextInputMode(msg)
	}

	// Normal mode commands
	switch msg.String() {
	case "q", "ctrl+c":
		// Save before quit
		a.save()
		return tea.Quit

	case "?":
		a.state.ShowHelp = !a.state.ShowHelp
		return nil

	case "j", "down":
		a.state.SelectNext()
	case "k", "up":
		a.state.SelectPrev()
	case "l", "right":
		a.state.NextPane()
	case "h", "left":
		a.state.PrevPane()

	case "n":
		// New task
		a.state.Mode = model.ModeInsert
		a.textInput.SetValue("")
		a.textInput.Placeholder = "Task title..."
		a.textInput.Focus()
		return textinput.Blink

	case "enter":
		// Open task detail view
		task := a.state.SelectedTask()
		if task != nil {
			a.taskDetail = NewTaskDetail(task)
			a.state.Mode = model.ModeDetail
		}

	case "e":
		// Edit task
		task := a.state.SelectedTask()
		if task != nil {
			a.taskEdit = NewTaskEdit(task)
			a.state.Mode = model.ModeEdit
			return a.taskEdit.Focus()
		}

	case "d":
		// Delete task with confirmation
		task := a.state.SelectedTask()
		if task != nil {
			a.showDeleteConfirmation(task, false)
		}

	case ">":
		// Move task to next pane with undo, following the task
		task := a.state.SelectedTask()
		if task != nil && task.Status != model.StatusDone {
			oldStatus := task.Status
			newStatus := task.Status.Next()
			taskID := task.ID
			cmd := undo.NewMoveTaskCommand(task.ID, oldStatus, newStatus)
			if err := a.undoManager.Execute(a.state.Board, cmd); err == nil {
				a.state.SwitchToTaskInPane(newStatus, taskID)
				a.state.MarkDirty()
				a.state.SetStatusMessage(fmt.Sprintf("Moved to %s", newStatus))
			}
		}

	case "<":
		// Move task to previous pane with undo, following the task
		task := a.state.SelectedTask()
		if task != nil && task.Status != model.StatusToDo {
			oldStatus := task.Status
			newStatus := task.Status.Prev()
			taskID := task.ID
			cmd := undo.NewMoveTaskCommand(task.ID, oldStatus, newStatus)
			if err := a.undoManager.Execute(a.state.Board, cmd); err == nil {
				a.state.SwitchToTaskInPane(newStatus, taskID)
				a.state.MarkDirty()
				a.state.SetStatusMessage(fmt.Sprintf("Moved to %s", newStatus))
			}
		}

	case "1":
		a.setPriorityWithUndo(model.PriorityHigh)
	case "2":
		a.setPriorityWithUndo(model.PriorityMedium)
	case "3":
		a.setPriorityWithUndo(model.PriorityLow)
	case "0":
		a.setPriorityWithUndo(model.PriorityNone)

	case "L":
		// Manage labels
		task := a.state.SelectedTask()
		if task != nil {
			a.labelEditor = NewLabelEditor(task)
			a.state.Mode = model.ModeLabels
			return a.labelEditor.Focus()
		}

	case "t":
		// Set due date
		task := a.state.SelectedTask()
		if task != nil {
			a.dueDateEditor = NewDueDateEditor(task)
			a.state.Mode = model.ModeDueDate
			return a.dueDateEditor.Focus()
		}

	case "b":
		// Switch board
		a.save() // Save current board first
		a.boardSelector = NewBoardSelector(a.storage, a.state.Board)
		a.state.Mode = model.ModeBoard
		return nil

	case "B":
		// Create new board
		a.save() // Save current board first
		a.boardSelector = NewBoardSelector(a.storage, a.state.Board)
		a.state.Mode = model.ModeBoard
		return a.boardSelector.SetMode(BoardModeCreate)

	case "s":
		// Sort by priority
		if a.state.Board != nil {
			a.state.Board.SortByPriority(a.state.CurrentPane)
			a.state.MarkDirty()
			a.state.SetStatusMessage("Sorted by priority")
		}

	case "u":
		// Undo
		if cmd, err := a.undoManager.Undo(a.state.Board); err == nil && cmd != nil {
			a.state.ClampSelection()
			a.state.MarkDirty()
			a.state.SetStatusMessage(fmt.Sprintf("Undo: %s", cmd.Description()))
		} else if cmd == nil {
			a.state.SetStatusMessage("Nothing to undo")
		}

	case "ctrl+r":
		// Redo
		if cmd, err := a.undoManager.Redo(a.state.Board); err == nil && cmd != nil {
			a.state.ClampSelection()
			a.state.MarkDirty()
			a.state.SetStatusMessage(fmt.Sprintf("Redo: %s", cmd.Description()))
		} else if cmd == nil {
			a.state.SetStatusMessage("Nothing to redo")
		}

	case "/":
		// Enter search mode
		a.search = NewSearch(a.state.Board)
		a.state.Mode = model.ModeSearch
		return a.search.Focus()

	case "f":
		// Enter filter mode
		a.filter = NewFilter(a.state.Board, a.state.FilterPriority, a.state.FilterLabels)
		a.state.Mode = model.ModeFilter
		return nil

	case "F":
		// Clear filters
		a.state.FilterPriority = nil
		a.state.FilterLabels = nil
		a.state.SetStatusMessage("Filters cleared")

	case "E":
		// Export current board
		if a.state.Board != nil {
			exportPath := a.storage.DefaultExportPath(a.state.Board)
			if err := a.storage.ExportBoard(a.state.Board, exportPath); err != nil {
				a.state.SetStatusMessage(fmt.Sprintf("Export failed: %v", err))
			} else {
				a.state.SetStatusMessage(fmt.Sprintf("Exported to %s", exportPath))
			}
		}

	case "I":
		// Import board (show file picker placeholder)
		a.state.SetStatusMessage("Import: Use 'tsk import <file>' from command line")

	case "a":
		// Archive selected task (DONE pane only)
		if a.state.CurrentPane == model.StatusDone {
			task := a.state.SelectedTask()
			if task != nil {
				a.archiveTask(task)
			}
		}

	case "A":
		// Archive all Done tasks
		if a.state.CurrentPane == model.StatusDone {
			a.archiveAllDoneTasks()
		}

	case "esc":
		a.state.ShowHelp = false
		a.state.ClearSearch()
	}

	return nil
}

func (a *App) handleTextInputMode(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc":
		// Can't escape welcome screen without creating a board
		if a.state.Mode == model.ModeWelcome {
			return nil
		}
		a.state.Mode = model.ModeNormal
		a.textInput.Blur()
		return nil

	case "enter":
		value := strings.TrimSpace(a.textInput.Value())

		if a.state.Mode == model.ModeWelcome {
			// Create first board
			if value == "" {
				value = "My Tasks"
			}
			board := model.NewBoard(uuid.New().String(), value)
			if err := a.storage.SaveBoard(board); err != nil {
				return func() tea.Msg { return err }
			}
			a.state.SetBoard(board)
			a.state.Mode = model.ModeNormal
			a.textInput.Blur()
			a.state.SetStatusMessage("Welcome to tsk! Press 'n' to create your first task.")
			return nil
		}

		if a.state.Mode == model.ModeInsert && value != "" {
			// Create new task with undo
			taskID := uuid.New().String()
			cmd := undo.NewCreateTaskCommand(taskID, value, a.state.CurrentPane)
			if err := a.undoManager.Execute(a.state.Board, cmd); err == nil {
				a.state.MarkDirty()
				a.state.SetStatusMessage(fmt.Sprintf("Created: %s", value))
			}
		}
		a.state.Mode = model.ModeNormal
		a.textInput.Blur()
		return nil
	}

	// Update text input
	var cmd tea.Cmd
	a.textInput, cmd = a.textInput.Update(msg)

	// Update search results live
	if a.state.Mode == model.ModeSearch {
		a.state.UpdateSearch(a.textInput.Value())
	}

	return cmd
}

func (a *App) setPriority(p model.Priority) {
	task := a.state.SelectedTask()
	if task != nil {
		task.SetPriority(p)
		a.state.MarkDirty()
		a.state.SetStatusMessage(fmt.Sprintf("Priority: %s", p))
	}
}

func (a *App) setPriorityWithUndo(p model.Priority) {
	task := a.state.SelectedTask()
	if task != nil {
		cmd := undo.NewSetPriorityCommand(task, p)
		if err := a.undoManager.Execute(a.state.Board, cmd); err == nil {
			a.state.MarkDirty()
			a.state.SetStatusMessage(fmt.Sprintf("Priority: %s", p))
		}
	}
}

func (a *App) handleModalMode(msg tea.KeyMsg) tea.Cmd {
	modal, ok := a.state.ActiveModal.(*Modal)
	if !ok {
		a.closeModal()
		return nil
	}

	switch msg.String() {
	case "esc", "n", "N":
		modal.Cancel()
		a.closeModal()

	case "enter":
		if modal.Type == ModalSelect {
			modal.Select()
		} else {
			modal.Confirm()
		}
		a.closeModal()

	case "y", "Y":
		if modal.Type == ModalConfirm {
			modal.Confirm()
			a.closeModal()
		}

	case "j", "down":
		modal.SelectNext()

	case "k", "up":
		modal.SelectPrev()
	}

	return nil
}

func (a *App) showModal(modal *Modal) {
	a.state.ActiveModal = modal
	a.state.Mode = model.ModeModal
}

func (a *App) closeModal() {
	a.state.ActiveModal = nil
	a.state.Mode = model.ModeNormal
}

func (a *App) showDeleteConfirmation(task *model.Task, fromDetail bool) {
	title := truncateString(task.Title, 30)
	modal := NewConfirmModal("Delete Task", fmt.Sprintf("Delete task '%s'?", title))
	modal.OnConfirm = func() {
		cmd := undo.NewDeleteTaskCommand(task)
		if err := a.undoManager.Execute(a.state.Board, cmd); err == nil {
			a.state.ClampSelection()
			a.state.MarkDirty()
			a.state.SetStatusMessage(fmt.Sprintf("Deleted: %s (u to undo)", task.Title))
		}
	}
	a.showModal(modal)
}

func (a *App) handleDetailMode(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc", "q":
		// Return to board view
		a.taskDetail = nil
		a.state.Mode = model.ModeNormal

	case "enter", "e":
		// Edit task from detail view
		if a.taskDetail != nil && a.taskDetail.Task != nil {
			a.taskEdit = NewTaskEdit(a.taskDetail.Task)
			a.taskDetail = nil
			a.state.Mode = model.ModeEdit
			return a.taskEdit.Focus()
		}

	case "d":
		// Delete task from detail view with confirmation
		if a.taskDetail != nil && a.taskDetail.Task != nil {
			task := a.taskDetail.Task
			a.taskDetail = nil
			a.showDeleteConfirmation(task, true)
		}

	case "L":
		// Manage labels from detail view
		if a.taskDetail != nil && a.taskDetail.Task != nil {
			a.labelEditor = NewLabelEditor(a.taskDetail.Task)
			a.taskDetail = nil
			a.state.Mode = model.ModeLabels
			return a.labelEditor.Focus()
		}

	case "1":
		a.setDetailPriority(model.PriorityHigh)
	case "2":
		a.setDetailPriority(model.PriorityMedium)
	case "3":
		a.setDetailPriority(model.PriorityLow)
	case "0":
		a.setDetailPriority(model.PriorityNone)
	}

	return nil
}

func (a *App) setDetailPriority(p model.Priority) {
	if a.taskDetail != nil && a.taskDetail.Task != nil {
		task := a.taskDetail.Task
		cmd := undo.NewSetPriorityCommand(task, p)
		if err := a.undoManager.Execute(a.state.Board, cmd); err == nil {
			a.state.MarkDirty()
			a.state.SetStatusMessage(fmt.Sprintf("Priority: %s", p))
		}
	}
}

func (a *App) handleEditMode(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc":
		// Cancel editing
		a.taskEdit = nil
		a.state.Mode = model.ModeNormal
		return nil

	case "tab":
		// Switch between fields
		return a.taskEdit.NextField()

	case "shift+tab":
		return a.taskEdit.PrevField()

	case "enter":
		// Save changes
		if a.taskEdit != nil && a.taskEdit.Task != nil {
			newTitle := a.taskEdit.GetTitle()
			if newTitle == "" {
				a.state.SetStatusMessage("Title cannot be empty")
				return nil
			}

			newDesc := a.taskEdit.GetDescription()
			newLabels := a.taskEdit.GetLabels()
			task := a.taskEdit.Task

			cmd := undo.NewEditTaskCommand(task, newTitle, newDesc, task.DueDate, newLabels)
			if err := a.undoManager.Execute(a.state.Board, cmd); err == nil {
				a.state.MarkDirty()
				a.state.SetStatusMessage(fmt.Sprintf("Updated: %s", newTitle))
			}
		}
		a.taskEdit = nil
		a.state.Mode = model.ModeNormal
		return nil
	}

	// Pass other keys to the edit component
	return a.taskEdit.Update(msg)
}

func (a *App) handleLabelsMode(msg tea.KeyMsg) tea.Cmd {
	if a.labelEditor == nil {
		a.state.Mode = model.ModeNormal
		return nil
	}

	// Check if the label editor wants to handle input
	if a.labelEditor.editing {
		// Pass to label editor for text input
		switch msg.String() {
		case "esc", "enter":
			done, cmd := a.labelEditor.HandleKey(msg.String())
			if done {
				a.saveLabels()
			}
			return cmd
		default:
			return a.labelEditor.Update(msg)
		}
	}

	// Handle navigation and actions
	done, cmd := a.labelEditor.HandleKey(msg.String())
	if done {
		a.saveLabels()
	}
	return cmd
}

func (a *App) saveLabels() {
	if a.labelEditor != nil && a.labelEditor.Task != nil {
		newLabels := a.labelEditor.GetLabels()
		cmd := undo.NewLabelChangeCommand(a.labelEditor.Task, newLabels)
		if err := a.undoManager.Execute(a.state.Board, cmd); err == nil {
			a.state.MarkDirty()
			a.state.SetStatusMessage("Labels updated")
		}
	}
	a.labelEditor = nil
	a.state.Mode = model.ModeNormal
}

func (a *App) handleDueDateMode(msg tea.KeyMsg) tea.Cmd {
	if a.dueDateEditor == nil {
		a.state.Mode = model.ModeNormal
		return nil
	}

	switch msg.String() {
	case "esc":
		a.dueDateEditor = nil
		a.state.Mode = model.ModeNormal
		return nil

	case "enter":
		if a.dueDateEditor.HasError() {
			return nil
		}
		a.saveDueDate()
		return nil

	case "ctrl+d":
		// Clear due date
		a.dueDateEditor.Clear()
		return nil
	}

	// Check for quick date keywords
	value := a.dueDateEditor.CurrentValue()
	if a.dueDateEditor.HandleQuickDate(value) {
		// If it was a quick date, we're done
		return nil
	}

	// Pass to editor for text input
	return a.dueDateEditor.Update(msg)
}

func (a *App) saveDueDate() {
	if a.dueDateEditor != nil && a.dueDateEditor.Task != nil {
		task := a.dueDateEditor.Task
		newDue := a.dueDateEditor.GetDueDate()

		// Use EditTaskCommand to save the change with undo support
		cmd := undo.NewEditTaskCommand(task, task.Title, task.Description, newDue, task.Labels)
		if err := a.undoManager.Execute(a.state.Board, cmd); err == nil {
			a.state.MarkDirty()
			if newDue != nil {
				a.state.SetStatusMessage(fmt.Sprintf("Due: %s", newDue.Format("Jan 2, 2006")))
			} else {
				a.state.SetStatusMessage("Due date cleared")
			}
		}
	}
	a.dueDateEditor = nil
	a.state.Mode = model.ModeNormal
}

func (a *App) handleSearchMode(msg tea.KeyMsg) tea.Cmd {
	if a.search == nil {
		a.state.Mode = model.ModeNormal
		return nil
	}

	key := msg.String()

	switch key {
	case "esc", "enter", "j", "k", "down", "up", "ctrl+d", "ctrl+u":
		task, done := a.search.HandleKey(key)
		if done {
			a.state.Mode = model.ModeNormal
			if task != nil {
				// Navigate to the selected task
				a.state.CurrentPane = task.Status
				// Find the task's index in the pane
				tasks := a.state.Board.TasksByStatus(task.Status)
				for i, t := range tasks {
					if t.ID == task.ID {
						a.state.SelectedIndex = i
						break
					}
				}
				a.state.SetStatusMessage(fmt.Sprintf("Found: %s", task.Title))
			}
			a.search = nil
			return nil
		}
		return nil
	}

	// Pass to search for text input
	return a.search.Update(msg)
}

func (a *App) handleFilterMode(msg tea.KeyMsg) tea.Cmd {
	if a.filter == nil {
		a.state.Mode = model.ModeNormal
		return nil
	}

	key := msg.String()
	done, apply := a.filter.HandleKey(key)

	if done {
		if apply {
			// Apply filters to state
			a.state.FilterPriority = a.filter.GetSelectedPriority()
			a.state.FilterLabels = a.filter.GetSelectedLabels()
			if a.filter.HasFilters() {
				a.state.SetStatusMessage("Filters applied (F to clear)")
			} else {
				a.state.SetStatusMessage("Filters cleared")
			}
		}
		a.filter = nil
		a.state.Mode = model.ModeNormal
	}

	return nil
}

func (a *App) handleBoardMode(msg tea.KeyMsg) tea.Cmd {
	if a.boardSelector == nil {
		a.state.Mode = model.ModeNormal
		return nil
	}

	mode := a.boardSelector.Mode()

	// Handle text input modes
	if mode == BoardModeCreate || mode == BoardModeRename {
		switch msg.String() {
		case "esc":
			a.boardSelector.SetMode(BoardModeSelect)
			return nil

		case "enter":
			name := a.boardSelector.GetInputValue()
			if name == "" {
				a.state.SetStatusMessage("Board name cannot be empty")
				return nil
			}

			if mode == BoardModeCreate {
				return a.createBoard(name)
			} else {
				return a.renameBoard(name)
			}

		default:
			return a.boardSelector.Update(msg)
		}
	}

	// Handle delete confirmation
	if mode == BoardModeDelete {
		switch msg.String() {
		case "esc":
			a.boardSelector.SetMode(BoardModeSelect)
			return nil

		case "y", "Y":
			return a.deleteSelectedBoard()

		case "d":
			if !a.boardSelector.IsConfirmingDelete() {
				a.boardSelector.ConfirmDelete()
			}
			return nil
		}
		return nil
	}

	// Handle board selection mode
	switch msg.String() {
	case "esc", "q":
		a.boardSelector = nil
		a.state.Mode = model.ModeNormal
		return nil

	case "j", "down":
		a.boardSelector.SelectNext()

	case "k", "up":
		a.boardSelector.SelectPrev()

	case "enter":
		return a.switchToSelectedBoard()

	case "B":
		return a.boardSelector.SetMode(BoardModeCreate)

	case "R":
		if a.boardSelector.SelectedBoard() != nil {
			return a.boardSelector.SetMode(BoardModeRename)
		}

	case "D":
		if a.boardSelector.SelectedBoard() != nil {
			a.boardSelector.SetMode(BoardModeDelete)
		}
	}

	return nil
}

func (a *App) createBoard(name string) tea.Cmd {
	board := model.NewBoard(uuid.New().String(), name)
	if err := a.storage.SaveBoard(board); err != nil {
		a.state.SetStatusMessage(fmt.Sprintf("Error: %v", err))
		return nil
	}

	// Switch to the new board
	a.state.SetBoard(board)
	a.undoManager.Clear()
	a.boardSelector = nil
	a.state.Mode = model.ModeNormal
	a.state.SetStatusMessage(fmt.Sprintf("Created board: %s", name))
	return nil
}

func (a *App) renameBoard(name string) tea.Cmd {
	selected := a.boardSelector.SelectedBoard()
	if selected == nil {
		return nil
	}

	// Load the board to rename
	board, err := a.storage.LoadBoard(selected.ID)
	if err != nil {
		a.state.SetStatusMessage(fmt.Sprintf("Error: %v", err))
		return nil
	}

	board.Name = name
	if err := a.storage.SaveBoard(board); err != nil {
		a.state.SetStatusMessage(fmt.Sprintf("Error: %v", err))
		return nil
	}

	// If it's the current board, update the name
	if a.state.Board != nil && a.state.Board.ID == selected.ID {
		a.state.Board.Name = name
	}

	a.boardSelector.Refresh()
	a.boardSelector.SetMode(BoardModeSelect)
	a.state.SetStatusMessage(fmt.Sprintf("Renamed to: %s", name))
	return nil
}

func (a *App) deleteSelectedBoard() tea.Cmd {
	selected := a.boardSelector.SelectedBoard()
	if selected == nil {
		return nil
	}

	// Cannot delete the only board
	if len(a.boardSelector.boards) <= 1 {
		a.state.SetStatusMessage("Cannot delete the only board")
		a.boardSelector.SetMode(BoardModeSelect)
		return nil
	}

	// Delete the board
	if err := a.storage.DeleteBoard(selected.ID); err != nil {
		a.state.SetStatusMessage(fmt.Sprintf("Error: %v", err))
		return nil
	}

	deletedName := selected.Name
	wasCurrentBoard := a.state.Board != nil && a.state.Board.ID == selected.ID

	a.boardSelector.Refresh()

	// If we deleted the current board, switch to another
	if wasCurrentBoard {
		newBoard := a.boardSelector.SelectedBoard()
		if newBoard != nil {
			board, err := a.storage.LoadBoard(newBoard.ID)
			if err == nil {
				a.state.SetBoard(board)
				a.undoManager.Clear()
			}
		}
	}

	a.boardSelector.SetMode(BoardModeSelect)
	a.state.SetStatusMessage(fmt.Sprintf("Deleted: %s", deletedName))
	return nil
}

func (a *App) switchToSelectedBoard() tea.Cmd {
	selected := a.boardSelector.SelectedBoard()
	if selected == nil {
		return nil
	}

	// Don't switch if already on this board
	if a.state.Board != nil && a.state.Board.ID == selected.ID {
		a.boardSelector = nil
		a.state.Mode = model.ModeNormal
		return nil
	}

	// Save current board
	a.save()

	// Load the new board
	board, err := a.storage.LoadBoard(selected.ID)
	if err != nil {
		a.state.SetStatusMessage(fmt.Sprintf("Error: %v", err))
		return nil
	}

	a.state.SetBoard(board)
	a.undoManager.Clear()
	a.boardSelector = nil
	a.state.Mode = model.ModeNormal
	a.state.SetStatusMessage(fmt.Sprintf("Switched to: %s", board.Name))
	return nil
}

func (a *App) archiveTask(task *model.Task) {
	if a.state.Board == nil || task == nil {
		return
	}

	// Archive the task
	if err := a.storage.ArchiveTask(task, a.state.Board); err != nil {
		a.state.SetStatusMessage(fmt.Sprintf("Error: %v", err))
		return
	}

	// Remove from board
	a.state.Board.RemoveTask(task.ID)
	a.state.ClampSelection()
	a.state.MarkDirty()
	a.state.SetStatusMessage(fmt.Sprintf("Archived: %s", task.Title))
}

func (a *App) archiveAllDoneTasks() {
	if a.state.Board == nil {
		return
	}

	doneTasks := a.state.Board.TasksByStatus(model.StatusDone)
	if len(doneTasks) == 0 {
		a.state.SetStatusMessage("No done tasks to archive")
		return
	}

	// Archive all done tasks
	if err := a.storage.ArchiveTasks(doneTasks, a.state.Board); err != nil {
		a.state.SetStatusMessage(fmt.Sprintf("Error: %v", err))
		return
	}

	// Remove all from board
	for _, task := range doneTasks {
		a.state.Board.RemoveTask(task.ID)
	}
	a.state.ClampSelection()
	a.state.MarkDirty()
	a.state.SetStatusMessage(fmt.Sprintf("Archived %d tasks", len(doneTasks)))
}

func (a *App) save() error {
	if a.state.Board == nil {
		return nil
	}
	if err := a.storage.SaveBoard(a.state.Board); err != nil {
		return err
	}
	a.state.MarkClean()
	return nil
}

// View renders the application
func (a *App) View() string {
	if a.state.Width == 0 {
		return "Loading..."
	}

	// Welcome screen for first run
	if a.state.Mode == model.ModeWelcome {
		return a.renderWelcomeScreen()
	}

	// Board selector
	if a.state.Mode == model.ModeBoard && a.boardSelector != nil {
		return a.boardSelector.View(a.state.Width, a.state.Height)
	}

	// Due date mode
	if a.state.Mode == model.ModeDueDate && a.dueDateEditor != nil {
		return a.dueDateEditor.View(a.state.Width, a.state.Height)
	}

	// Labels mode
	if a.state.Mode == model.ModeLabels && a.labelEditor != nil {
		return a.labelEditor.View(a.state.Width, a.state.Height)
	}

	// Filter mode
	if a.state.Mode == model.ModeFilter && a.filter != nil {
		return a.filter.View(a.state.Width, a.state.Height)
	}

	// Task edit view
	if a.state.Mode == model.ModeEdit && a.taskEdit != nil {
		return a.taskEdit.View(a.state.Width, a.state.Height)
	}

	// Task detail view
	if a.state.Mode == model.ModeDetail && a.taskDetail != nil {
		return a.taskDetail.View(a.state.Width, a.state.Height)
	}

	// Modal overlay (for confirmations, selections)
	if a.state.Mode == model.ModeModal && a.state.ActiveModal != nil {
		if modal, ok := a.state.ActiveModal.(*Modal); ok {
			return modal.View(a.state.Width, a.state.Height)
		}
	}

	// Help overlay (full screen)
	if a.state.ShowHelp {
		return a.renderHelpOverlay()
	}

	// Search mode
	if a.state.Mode == model.ModeSearch && a.search != nil {
		return a.search.View(a.state.Width, a.state.Height)
	}

	// Text input overlay (when in insert/search mode)
	if a.state.Mode.IsTextInput() {
		return a.renderWithTextInput("")
	}

	// Calculate layout heights
	headerHeight := 1
	tabsHeight := 1
	statusHeight := 2
	contentHeight := a.state.Height - headerHeight - tabsHeight - statusHeight - 2 // -2 for newlines
	if contentHeight < 1 {
		contentHeight = 1
	}

	// Two-panel layout: task list (30%) + preview (70%)
	leftWidth := a.state.Width * 30 / 100
	rightWidth := a.state.Width - leftWidth

	taskList := a.renderTaskList(leftWidth, contentHeight)
	preview := a.renderPreview(rightWidth, contentHeight)

	// Join panels horizontally
	content := lipgloss.JoinHorizontal(lipgloss.Top, taskList, preview)

	// Pad content to fill available height
	contentLines := strings.Count(content, "\n") + 1
	padding := ""
	if contentLines < contentHeight {
		padding = strings.Repeat("\n", contentHeight-contentLines)
	}

	// Build the full view
	return lipgloss.JoinVertical(
		lipgloss.Left,
		a.renderHeader(),
		a.renderTabs(),
		content+padding,
		a.renderStatusBar(),
	)
}

func (a *App) renderHeader() string {
	title := styles.TitleStyle.Render("tsk")
	boardName := ""
	if a.state.Board != nil {
		boardName = styles.BoardNameStyle.Render(" │ " + a.state.Board.Name)
	}

	helpHint := styles.HelpHintStyle.Render("Press ? for help")

	// Calculate spacing
	leftPart := title + boardName
	spacing := a.state.Width - lipgloss.Width(leftPart) - lipgloss.Width(helpHint) - 2
	if spacing < 0 {
		spacing = 0
	}

	return leftPart + strings.Repeat(" ", spacing) + helpHint
}

func (a *App) renderTabs() string {
	var tabs []string
	for _, status := range model.AllStatuses() {
		count := 0
		if a.state.Board != nil {
			count = a.state.Board.TaskCount(status)
		}

		label := fmt.Sprintf("%s (%d)", status.String(), count)

		if status == a.state.CurrentPane {
			tabs = append(tabs, styles.TabStyleForStatus(status, true).Render("["+label+"]"))
		} else {
			tabs = append(tabs, styles.TabStyleForStatus(status, false).Render(label))
		}
	}
	return strings.Join(tabs, "  ")
}

func (a *App) renderTaskList(width, height int) string {
	if a.state.Board == nil {
		return styles.TaskListStyle.Width(width - 2).Render("No board loaded")
	}

	tasks := a.state.CurrentTasks()
	if len(tasks) == 0 {
		return styles.TaskListStyle.Width(width - 2).Render(styles.HelpHintStyle.Render("No tasks - press 'n' to create one"))
	}

	// Reserve 2 lines for scroll indicators (always present to avoid layout shift)
	visibleHeight := height - 4
	if visibleHeight < 1 {
		visibleHeight = 1
	}

	// Adjust scroll offset only when selection goes outside viewport
	if a.state.SelectedIndex < a.state.ScrollOffset {
		// Selection above viewport - scroll up
		a.state.ScrollOffset = a.state.SelectedIndex
	} else if a.state.SelectedIndex >= a.state.ScrollOffset+visibleHeight {
		// Selection below viewport - scroll down
		a.state.ScrollOffset = a.state.SelectedIndex - visibleHeight + 1
	}

	// Clamp scroll offset
	if a.state.ScrollOffset < 0 {
		a.state.ScrollOffset = 0
	}
	maxOffset := len(tasks) - visibleHeight
	if maxOffset < 0 {
		maxOffset = 0
	}
	if a.state.ScrollOffset > maxOffset {
		a.state.ScrollOffset = maxOffset
	}

	startIdx := a.state.ScrollOffset
	endIdx := startIdx + visibleHeight
	if endIdx > len(tasks) {
		endIdx = len(tasks)
	}

	var lines []string
	maxTitleWidth := width - 7 // Account for " ▶● " prefix

	// Top scroll indicator (always present)
	if startIdx > 0 {
		lines = append(lines, styles.HelpHintStyle.Render(fmt.Sprintf("  ↑ %d more above", startIdx)))
	} else {
		lines = append(lines, "") // Empty line to reserve space
	}

	for i := startIdx; i < endIdx; i++ {
		task := tasks[i]

		// Priority indicator
		priority := styles.PriorityIndicator(task.Priority)

		// Task title - truncate if needed
		title := task.Title
		if len(title) > maxTitleWidth {
			title = title[:maxTitleWidth-3] + "..."
		}

		// Pad title to fill width
		padding := ""
		if len(title) < maxTitleWidth {
			padding = strings.Repeat(" ", maxTitleWidth-len(title))
		}

		// Style based on selection
		var line string
		if i == a.state.SelectedIndex {
			line = fmt.Sprintf(" ▶%s %s%s", priority, styles.TaskSelectedStyle.Render(title), padding)
		} else {
			line = fmt.Sprintf("  %s %s%s", priority, styles.TaskNormalStyle.Render(title), padding)
		}

		lines = append(lines, line)
	}

	// Bottom scroll indicator (always present)
	if endIdx < len(tasks) {
		lines = append(lines, styles.HelpHintStyle.Render(fmt.Sprintf("  ↓ %d more below", len(tasks)-endIdx)))
	} else {
		lines = append(lines, "") // Empty line to reserve space
	}

	content := strings.Join(lines, "\n")
	return styles.TaskListStyle.Width(width - 2).Render(content)
}

func (a *App) renderPreview(width, height int) string {
	task := a.state.SelectedTask()
	if task == nil {
		return styles.PreviewStyle.Width(width - 2).Render(styles.HelpHintStyle.Render("Select a task to preview"))
	}

	var lines []string

	lines = append(lines, styles.PreviewTitleStyle.Render(task.Title))
	lines = append(lines, "")
	lines = append(lines, styles.PreviewLabelStyle.Render("Status: ")+styles.StatusStyle(task.Status).Render(task.Status.String()))
	lines = append(lines, styles.PreviewLabelStyle.Render("Priority: ")+styles.PriorityStyle(task.Priority).Render(task.Priority.String()))

	if task.DueDate != nil {
		lines = append(lines, styles.PreviewLabelStyle.Render("Due: ")+styles.PreviewValueStyle.Render(task.DueDate.Format("Jan 2, 2006")))
	}

	if len(task.Labels) > 0 {
		labels := strings.Join(task.Labels, ", ")
		lines = append(lines, styles.PreviewLabelStyle.Render("Labels: ")+styles.PreviewValueStyle.Render(labels))
	}

	if task.Description != "" {
		lines = append(lines, "")
		lines = append(lines, styles.PreviewValueStyle.Render(task.Description))
	}

	content := strings.Join(lines, "\n")
	return styles.PreviewStyle.Width(width - 2).Render(content)
}

func (a *App) renderStatusBar() string {
	// Line 1: Context or feedback message
	line1 := ""
	if a.state.StatusMessage != "" {
		line1 = a.state.StatusMessage
	} else if a.state.Board != nil {
		count := len(a.state.CurrentTasks())
		totalCount := a.state.Board.TaskCount(a.state.CurrentPane)
		if a.state.HasActiveFilters() && count != totalCount {
			line1 = fmt.Sprintf("%d/%d tasks in %s (filtered)", count, totalCount, a.state.CurrentPane)
		} else {
			line1 = fmt.Sprintf("%d tasks in %s", count, a.state.CurrentPane)
		}
		if a.state.Dirty {
			line1 += " (unsaved)"
		}
	}

	// Add filter indicator
	if a.state.HasActiveFilters() {
		line1 = styles.WarningStyle.Render("[FILTERED]") + "  " + line1
	}

	// Add mode indicator
	if a.state.Mode == model.ModeInsert {
		line1 = styles.ModeInsertStyle.Render("-- INSERT --") + "  " + line1
	} else if a.state.Mode == model.ModeSearch {
		line1 = styles.ModeSearchStyle.Render("-- SEARCH --") + "  " + line1
	}

	// Line 2: Shortcuts (context-sensitive)
	var shortcuts string
	if a.state.CurrentPane == model.StatusDone {
		shortcuts = "j/k:nav  h/l:pane  d:del  a:archive  A:archive all  Enter:edit  b:board"
	} else {
		shortcuts = "j/k:nav  h/l:pane  n:new  d:del  >/<:move  Enter:edit  1-3:priority  b:board"
	}

	return styles.StatusLine1Style.Render(line1) + "\n" + styles.StatusLine2Style.Render(shortcuts)
}

func (a *App) renderHelpOverlay() string {
	help := `
                              KEYBOARD SHORTCUTS                              
──────────────────────────────────────────────────────────────────────────────
  NAVIGATION                           TASK ACTIONS                           
  j/k or ↓/↑    Move between tasks     n          Create new task             
  h/l or ←/→    Switch panes           Enter      Edit task                   
                                       d          Delete task                 
  MOVE TASKS                           > / <      Move task right/left        
  >             Move to next pane                                             
  <             Move to previous pane  PRIORITY                               
                                       1          High priority               
  BOARD                                2          Medium priority             
  b             Switch board           3          Low priority                
  B             Create new board       0          Clear priority              
  R             Rename board                                                  
  D             Delete board           LABELS                                 
                                       L          Manage labels               
  SEARCH & FILTER                                                             
  /             Search tasks           OTHER                                  
  s             Sort by priority       u          Undo                        
  F             Clear filters          Ctrl+r     Redo                        
                                       q / :wq    Quit                        
                                       ?          Toggle this help            
──────────────────────────────────────────────────────────────────────────────
                            Press ? or Esc to close                           
`
	return styles.ModalStyle.Render(help)
}

func (a *App) renderWithTextInput(base string) string {
	// Create a modal with text input
	var title string
	if a.state.Mode == model.ModeInsert {
		title = "New Task"
	} else {
		title = "Search"
	}

	modal := styles.ModalTitleStyle.Render(title) + "\n\n" + a.textInput.View() + "\n\n" +
		styles.HelpHintStyle.Render("Enter to confirm, Esc to cancel")

	return styles.ModalStyle.Render(modal)
}

func (a *App) renderWelcomeScreen() string {
	welcome := `
████████╗███████╗██╗  ██╗
╚══██╔══╝██╔════╝██║ ██╔╝
   ██║   ███████╗█████╔╝ 
   ██║   ╚════██║██╔═██╗ 
   ██║   ███████║██║  ██╗
   ╚═╝   ╚══════╝╚═╝  ╚═╝
   Terminal Task Manager
`

	content := styles.ModalTitleStyle.Render(welcome) + "\n\n" +
		styles.PreviewLabelStyle.Render("Welcome! Let's create your first board.") + "\n\n" +
		styles.HelpHintStyle.Render("Board name:") + "\n" +
		a.textInput.View() + "\n\n" +
		styles.HelpHintStyle.Render("Press Enter to continue (or leave blank for 'My Tasks')")

	// Center the content
	modalWidth := 60
	modalHeight := 16

	paddingX := (a.state.Width - modalWidth) / 2
	paddingY := (a.state.Height - modalHeight) / 2
	if paddingX < 0 {
		paddingX = 0
	}
	if paddingY < 0 {
		paddingY = 0
	}

	modal := styles.ModalStyle.Width(modalWidth).Render(content)

	// Add vertical padding
	vertPad := strings.Repeat("\n", paddingY)

	return vertPad + modal
}

// truncateString truncates a string to maxLen characters, adding "..." if truncated
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
