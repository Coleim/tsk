package ui

import (
	"fmt"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lipglossv1 "github.com/charmbracelet/lipgloss"
	"github.com/coliva/tsk/internal/model"
	"github.com/coliva/tsk/internal/storage"
	"github.com/coliva/tsk/internal/styles"
	"github.com/coliva/tsk/internal/undo"
	"github.com/google/uuid"
)

// Layout constants for responsive design
const (
	MinTaskListWidth     = 39 // Minimum width for task list to be readable
	SinglePanelThreshold = 59 // Below this width, hide preview panel
	UIMargin             = 1  // Small margin around the entire UI
)

// Messages
type (
	tickMsg        time.Time
	statusClearMsg struct{}
	boardLoadedMsg struct {
		board *model.Board
		err   error
	}
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
	search *SimpleSearch

	// Filter component
	filter *Filter

	// Auto-save ticker
	lastAutoSave time.Time

	// Loading state
	loading bool
	spinner spinner.Model

	// Error state
	err error

	// Pending key for multi-key commands (like gg)
	pendingKey string
}

// NewApp creates a new application
func NewApp(store *storage.Storage) *App {
	ti := textinput.New()
	ti.Placeholder = "Enter text..."
	ti.CharLimit = 256

	// Create spinner with dots style (using v1 lipgloss for bubbles compatibility)
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipglossv1.NewStyle().Foreground(lipglossv1.Color("213"))

	return &App{
		state:        model.NewAppState(),
		storage:      store,
		textInput:    ti,
		undoManager:  undo.NewManager(20),
		lastAutoSave: time.Now(),
		loading:      true,
		spinner:      s,
	}
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	return tea.Batch(
		a.spinner.Tick,
		a.loadInitialBoard(),
		a.tickCmd(),
	)
}

func (a *App) loadInitialBoard() tea.Cmd {
	return func() tea.Msg {
		hasBoards, err := a.storage.HasBoards()
		if err != nil {
			return boardLoadedMsg{err: err}
		}

		if !hasBoards {
			// First run - show welcome screen
			return boardLoadedMsg{board: nil}
		}

		// Load most recent board
		board, err := a.storage.MostRecentBoard()
		if err != nil {
			return boardLoadedMsg{err: err}
		}
		return boardLoadedMsg{board: board}
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
		a.loading = false
		a.state.SetStatusMessage(fmt.Sprintf("Error: %v", msg))
		return a, nil

	case boardLoadedMsg:
		a.loading = false
		if msg.err != nil {
			a.err = msg.err
			a.state.SetStatusMessage(fmt.Sprintf("Error: %v", msg.err))
			return a, nil
		}
		if msg.board == nil {
			// First run - show welcome screen
			a.state.Mode = model.ModeWelcome
			a.textInput.SetValue("")
			a.textInput.Placeholder = "My Tasks"
			a.textInput.Focus()
		} else {
			// Ensure all labels have proper definitions
			msg.board.EnsureLabelsExist()
			a.state.SetBoard(msg.board)
		}
		return a, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		a.spinner, cmd = a.spinner.Update(msg)
		return a, cmd

	case tea.KeyMsg:
		cmd := a.handleKeyMsg(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case tea.WindowSizeMsg:
		a.state.Width = msg.Width
		a.state.Height = msg.Height

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

	// Clear pending key for any key except "g" (which may start a "gg" sequence)
	key := msg.String()
	if key != "g" {
		a.pendingKey = ""
	}

	// Normal mode commands
	switch key {
	case "q", "ctrl+c":
		// Save before quit
		_ = a.save()
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

	case "G":
		// Go to bottom (vim style)
		a.state.SelectLast()

	case "g":
		// First g - wait for second g
		if a.pendingKey == "g" {
			// gg - go to top (vim style)
			a.state.SelectFirst()
			a.pendingKey = ""
		} else {
			a.pendingKey = "g"
		}

	case "n":
		// New task
		a.state.Mode = model.ModeInsert
		a.textInput.SetValue("")
		a.textInput.Placeholder = "Task title..."
		a.textInput.Width = 44 // Fit within 50-char popup (accounting for padding)
		a.textInput.Focus()
		return textinput.Blink

	case "enter":
		// Open task detail view
		task := a.state.SelectedTask()
		if task != nil {
			a.taskDetail = NewTaskDetail(task, a.state.Board)
			a.state.Mode = model.ModeDetail
		}

	case "e":
		// Edit task
		task := a.state.SelectedTask()
		if task != nil {
			a.taskEdit = NewTaskEdit(task, a.state.Board)
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
			a.labelEditor = NewLabelEditor(task, a.state.Board)
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
		_ = a.save() // Save current board first
		a.boardSelector = NewBoardSelector(a.storage, a.state.Board)
		a.state.Mode = model.ModeBoard
		return nil

	case "B":
		// Create new board
		_ = a.save() // Save current board first
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
		a.search = NewSimpleSearch(a.state.Board)
		a.state.Mode = model.ModeSearch
		return a.search.Init()

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
			a.taskEdit = NewTaskEdit(a.taskDetail.Task, a.state.Board)
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
			a.labelEditor = NewLabelEditor(a.taskDetail.Task, a.state.Board)
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
	// Handle label popup if open
	if a.taskEdit.IsPopupOpen() {
		switch msg.String() {
		case "esc":
			// Close popup without selection
			a.taskEdit.CloseLabelPopup()
			return nil
		case "tab":
			// Move to next label in popup
			a.taskEdit.NextPopupLabel()
			return nil
		case "shift+tab":
			// Move to previous label in popup
			a.taskEdit.PrevPopupLabel()
			return nil
		case "enter":
			// Select current label and close popup
			a.taskEdit.SelectPopupLabel()
			return nil
		}
		// Ignore other keys while popup is open
		return nil
	}

	switch msg.String() {
	case "esc":
		// Cancel editing
		a.taskEdit = nil
		a.state.Mode = model.ModeNormal
		return nil

	case "tab":
		// If on labels field, open popup; otherwise switch fields
		if a.taskEdit.IsLabelsField() {
			a.taskEdit.OpenLabelPopup()
			return nil
		}
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

			// Ensure all labels exist in board (gives them colors)
			for _, labelName := range newLabels {
				a.state.Board.GetLabel(labelName)
			}

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

	case "tab":
		// Cycle through quick date options
		a.dueDateEditor.NextQuickOption()
		return nil

	case "enter":
		// If we have a quick selection or valid date, save it
		if a.dueDateEditor.HasQuickSelection() || !a.dueDateEditor.HasError() {
			a.saveDueDate()
		}
		return nil

	case "ctrl+d":
		// Clear due date
		a.dueDateEditor.Clear()
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

	done, task, cmd := a.search.Update(msg)
	if done {
		a.state.Mode = model.ModeNormal
		if task != nil {
			// Navigate to the selected task
			a.state.CurrentPane = task.Status
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
	}
	return cmd
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
	_ = a.save()

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

	// Archive the task to storage
	if err := a.storage.ArchiveTask(task, a.state.Board); err != nil {
		a.state.SetStatusMessage(fmt.Sprintf("Error: %v", err))
		return
	}

	// Use undo command to remove from board (with callback to unarchive)
	cmd := undo.NewArchiveTaskCommand(task, a.state.Board.ID, func(taskID, boardID string) error {
		return a.storage.UnarchiveTask(taskID, boardID)
	})
	if err := a.undoManager.Execute(a.state.Board, cmd); err != nil {
		a.state.SetStatusMessage(fmt.Sprintf("Error: %v", err))
		return
	}

	a.state.ClampSelection()
	a.state.MarkDirty()
	a.state.SetStatusMessage(fmt.Sprintf("Archived: %s (u to undo)", task.Title))
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

	// Archive all done tasks to storage
	if err := a.storage.ArchiveTasks(doneTasks, a.state.Board); err != nil {
		a.state.SetStatusMessage(fmt.Sprintf("Error: %v", err))
		return
	}

	// Use undo command to remove all from board (with callback to unarchive)
	cmd := undo.NewArchiveTasksCommand(doneTasks, a.state.Board.ID, func(taskIDs []string, boardID string) error {
		return a.storage.UnarchiveTasks(taskIDs, boardID)
	})
	if err := a.undoManager.Execute(a.state.Board, cmd); err != nil {
		a.state.SetStatusMessage(fmt.Sprintf("Error: %v", err))
		return
	}

	a.state.ClampSelection()
	a.state.MarkDirty()
	a.state.SetStatusMessage(fmt.Sprintf("Archived %d tasks (u to undo)", len(doneTasks)))
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
		return ""
	}

	// Loading state with spinner
	if a.loading {
		return a.renderLoadingScreen()
	}

	// Welcome screen for first run
	if a.state.Mode == model.ModeWelcome {
		return a.renderWelcomeScreen()
	}

	// Board selector
	if a.state.Mode == model.ModeBoard && a.boardSelector != nil {
		return a.boardSelector.View(a.state.Width, a.state.Height)
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

	// Modal overlay (for confirmations, selections) - popup on main view
	if a.state.Mode == model.ModeModal && a.state.ActiveModal != nil {
		if modal, ok := a.state.ActiveModal.(*Modal); ok {
			mainView := a.renderMainView()
			popup := modal.View(a.state.Width, a.state.Height)
			return overlayDialog(mainView, popup, a.state.Width, a.state.Height)
		}
	}

	// Help overlay (full screen)
	if a.state.ShowHelp {
		return a.renderHelpOverlay()
	}

	// Search mode: overlay popup on main view
	if a.state.Mode == model.ModeSearch && a.search != nil {
		mainView := a.renderMainView()
		popup := a.search.View()
		return overlayDialog(mainView, popup, a.state.Width, a.state.Height)
	}

	// Insert mode: overlay new task popup on main view
	if a.state.Mode == model.ModeInsert {
		mainView := a.renderMainView()
		popup := a.renderNewTaskPopup()
		return overlayDialog(mainView, popup, a.state.Width, a.state.Height)
	}

	// Due date mode: overlay popup on main view
	if a.state.Mode == model.ModeDueDate && a.dueDateEditor != nil {
		mainView := a.renderMainView()
		popup := a.dueDateEditor.View()
		return overlayDialog(mainView, popup, a.state.Width, a.state.Height)
	}

	// Labels mode: overlay popup on main view
	if a.state.Mode == model.ModeLabels && a.labelEditor != nil {
		mainView := a.renderMainView()
		popup := a.labelEditor.View()
		return overlayDialog(mainView, popup, a.state.Width, a.state.Height)
	}

	// Welcome mode (still full-screen for onboarding)
	if a.state.Mode == model.ModeWelcome {
		return a.renderWithTextInput("")
	}

	return a.renderMainView()
}

func (a *App) renderMainView() string {
	// Calculate layout dimensions (sides margin only)
	availableWidth := a.state.Width - (UIMargin * 2)
	availableHeight := a.state.Height

	headerHeight := 1
	tabsHeight := 1
	statusHeight := 2
	contentHeight := availableHeight - headerHeight - tabsHeight - statusHeight
	if contentHeight < 1 {
		contentHeight = 1
	}

	// Responsive width calculation
	var content string
	width := availableWidth

	if width < SinglePanelThreshold {
		// Single-panel mode: task list only, preview hidden
		taskList := a.renderTaskList(width, contentHeight)
		content = taskList
	} else {
		// Two-panel mode: task list + preview
		// Task list gets 30% but minimum MinTaskListWidth
		leftWidth := width * 30 / 100
		if leftWidth < MinTaskListWidth {
			leftWidth = MinTaskListWidth
		}
		rightWidth := width - leftWidth

		taskList := a.renderTaskList(leftWidth, contentHeight)
		preview := a.renderPreview(rightWidth, contentHeight)

		// Join panels horizontally
		content = lipgloss.JoinHorizontal(lipgloss.Top, taskList, preview)
	}

	// Pad content to fill available height, but cap at contentHeight
	contentLines := strings.Count(content, "\n") + 1
	padding := ""
	if contentLines < contentHeight {
		padding = strings.Repeat("\n", contentHeight-contentLines)
	}

	// Cap the final content to prevent overflow
	finalContent := content + padding
	finalLines := strings.Split(finalContent, "\n")
	if len(finalLines) > contentHeight {
		finalLines = finalLines[:contentHeight]
		finalContent = strings.Join(finalLines, "\n")
	}

	// Build the full view
	view := lipgloss.JoinVertical(
		lipgloss.Left,
		a.renderHeader(),
		a.renderTabs(),
		finalContent,
		a.renderStatusBar(),
	)

	// Add side margins only
	return lipgloss.NewStyle().MarginLeft(UIMargin).MarginRight(UIMargin).Render(view)
}

func (a *App) renderHeader() string {
	title := styles.TitleStyle().Render("tsk")
	boardName := ""
	if a.state.Board != nil {
		boardName = styles.BoardNameStyle().Render(" │ " + a.state.Board.Name)
	}

	helpHint := styles.HelpHintStyle().Render("Press ? for help")

	// Calculate spacing (account for UI margin)
	availableWidth := a.state.Width - (UIMargin * 2)
	leftPart := title + boardName
	spacing := availableWidth - lipgloss.Width(leftPart) - lipgloss.Width(helpHint) - 2
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
		return styles.TaskListStyle().Width(width - 2).Height(height).Render("No board loaded")
	}

	tasks := a.state.CurrentTasks()
	if len(tasks) == 0 {
		return styles.TaskListStyle().Width(width - 2).Height(height).Render(styles.EmptyStateStyle().Render("No tasks • Press 'n' to create"))
	}

	// Each bordered task takes 3 lines (top border + content + bottom border)
	// TaskListStyle adds 2 lines for border (no vertical padding)
	// Reserve 2 lines for scroll indicators (1 top, 1 bottom)
	// So: availableLines = height - 2 (task list border) - 2 (indicators) = height - 4
	// visibleTasks = availableLines / 3
	visibleHeight := (height - 4) / 3
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
	// Width calculations for task items:
	// - TaskListStyle has border(2) + padding(4) = 6 chars horizontally
	// - TaskSelectedStyle/TaskNormalStyle have border(2) + padding(2) = 4 chars
	// - Prefix "▶ ● " or "  ● " = 5 visual chars
	// So: outer width -> task width -> inner content -> title space
	taskWidth := width - 8      // Leave room for TaskListStyle border/padding
	innerWidth := taskWidth - 4 // Task border + padding
	prefixWidth := 5            // "▶ ● " or "  ● "
	maxTitleWidth := innerWidth - prefixWidth
	if maxTitleWidth < 10 {
		maxTitleWidth = 10
	}

	// Top scroll indicator
	if startIdx > 0 {
		lines = append(lines, styles.HelpHintStyle().Render(fmt.Sprintf("  ↑ %d more above", startIdx)))
	} else {
		lines = append(lines, styles.HelpHintStyle().Render("  ─ top ─"))
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

		// Style based on selection - both have rounded borders
		var line string
		if i == a.state.SelectedIndex {
			// Selected: arrow + priority + title, with accent border
			content := fmt.Sprintf("▶ %s %s", priority, title)
			line = styles.TaskSelectedStyle().Width(taskWidth).Render(content)
		} else {
			// Unselected: priority + title, with subtle border
			content := fmt.Sprintf("  %s %s", priority, title)
			line = styles.TaskNormalStyle().Width(taskWidth).Render(content)
		}

		lines = append(lines, line)
	}

	// Calculate how many lines we've used:
	// - 1 for top indicator
	// - 3 per task (bordered task = top border + content + bottom border)
	// - 1 for bottom indicator
	// TaskListStyle adds border (1 top + 1 bottom = 2 total)
	usedLines := 1 + (endIdx-startIdx)*3 + 1
	availableInner := height - 2 // Subtract TaskListStyle border(2)
	spacerLines := availableInner - usedLines
	if spacerLines < 0 {
		spacerLines = 0
	}

	// Add spacer to push bottom indicator down
	for i := 0; i < spacerLines; i++ {
		lines = append(lines, "")
	}

	// Bottom scroll indicator (always reserve 1 line)
	if endIdx < len(tasks) {
		lines = append(lines, styles.HelpHintStyle().Render(fmt.Sprintf("  ↓ %d more below", len(tasks)-endIdx)))
	} else {
		lines = append(lines, styles.HelpHintStyle().Render("  ─ end ─"))
	}

	// Join vertically without extra spacing between bordered tasks
	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return styles.TaskListStyle().Width(width - 2).Height(height).Render(content)
}

func (a *App) renderPreview(width, height int) string {
	task := a.state.SelectedTask()
	if task == nil {
		return styles.PreviewStyle().Width(width - 2).Height(height).Render(styles.HelpHintStyle().Render("Select a task to preview"))
	}

	var lines []string

	lines = append(lines, styles.PreviewTitleStyle().Render(task.Title))
	lines = append(lines, "")
	lines = append(lines, styles.PreviewLabelStyle().Render("Status: ")+styles.StatusStyle(task.Status).Render(task.Status.String()))
	lines = append(lines, styles.PreviewLabelStyle().Render("Priority: ")+styles.PriorityStyle(task.Priority).Render(task.Priority.String()))

	if task.DueDate != nil {
		lines = append(lines, styles.PreviewLabelStyle().Render("Due: ")+styles.PreviewValueStyle().Render(task.DueDate.Format("Jan 2, 2006")))
	}

	if len(task.Labels) > 0 && a.state.Board != nil {
		var labelBadges []string
		for _, labelName := range task.Labels {
			label := a.state.Board.GetLabel(labelName)
			labelBadges = append(labelBadges, styles.LabelBadge(label.Name, label.Color))
		}
		lines = append(lines, styles.PreviewLabelStyle().Render("Labels: ")+strings.Join(labelBadges, " "))
	}

	if task.Description != "" {
		lines = append(lines, "")
		lines = append(lines, styles.PreviewValueStyle().Render(task.Description))
	}

	content := strings.Join(lines, "\n")
	return styles.PreviewStyle().Width(width - 2).Height(height).Render(content)
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
		line1 = styles.WarningStyle().Render("[FILTERED]") + "  " + line1
	}

	// Add mode indicator
	switch a.state.Mode {
	case model.ModeInsert:
		line1 = styles.ModeInsertStyle().Render("-- INSERT --") + "  " + line1
	case model.ModeSearch:
		line1 = styles.ModeSearchStyle().Render("-- SEARCH --") + "  " + line1
	}

	// Line 2: Shortcuts (context-sensitive)
	var shortcuts string
	if a.state.CurrentPane == model.StatusDone {
		shortcuts = "j/k:nav  h/l:pane  d:del  a:archive  A:archive all  Enter:edit  b:board"
	} else {
		shortcuts = "j/k:nav  h/l:pane  n:new  d:del  >/<:move  Enter:edit  1-3:priority  b:board"
	}

	return styles.StatusLine1Style().Render(line1) + "\n" + styles.StatusLine2Style().Render(shortcuts)
}

func (a *App) renderHelpOverlay() string {
	help := `
                              KEYBOARD SHORTCUTS                              
──────────────────────────────────────────────────────────────────────────────
  NAVIGATION                           TASK ACTIONS                           
  j/k or ↓/↑    Move between tasks     n          Create new task             
  h/l or ←/→    Switch panes           Enter      Edit task                   
  G             Go to bottom           d          Delete task                 
  gg            Go to top                 
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
	// Full-screen help
	editWidth := a.state.Width - 4
	if editWidth < 50 {
		editWidth = 50
	}
	return styles.ModalStyle().Width(editWidth).Height(a.state.Height - 4).Render(help)
}

func (a *App) renderWithTextInput(base string) string {
	// Create a full-screen dialog with text input
	var title string
	if a.state.Mode == model.ModeInsert {
		title = "New Task"
	} else {
		title = "Search"
	}

	modal := styles.ModalTitleStyle().Render(title) + "\n\n" + a.textInput.View() + "\n\n" +
		styles.HelpHintStyle().Render("Enter to confirm, Esc to cancel")

	// Full-screen layout
	editWidth := a.state.Width - 4
	if editWidth < 50 {
		editWidth = 50
	}
	return styles.ModalStyle().Width(editWidth).Height(a.state.Height - 4).Render(modal)
}

// renderNewTaskPopup renders a compact popup for new task creation
func (a *App) renderNewTaskPopup() string {
	var b strings.Builder

	// Title
	b.WriteString(styles.ModalTitleStyle().Render("New Task"))
	b.WriteString("\n\n")

	// Input
	b.WriteString(a.textInput.View())
	b.WriteString("\n\n")

	// Help
	b.WriteString(styles.HelpHintStyle().Render("Enter:create Esc:cancel"))

	// Create popup box (50 chars wide like search popup)
	return styles.ModalStyle().
		Width(50).
		Render(b.String())
}

func (a *App) renderLoadingScreen() string {
	// Center the loading indicator
	loadingText := a.spinner.View() + " Loading tasks..."

	// Create a nicely styled loading box
	loadingBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Padding(2, 4).
		Render(loadingText)

	// Center horizontally and vertically
	boxWidth := lipgloss.Width(loadingBox)
	boxHeight := lipgloss.Height(loadingBox)

	paddingX := (a.state.Width - boxWidth) / 2
	paddingY := (a.state.Height - boxHeight) / 2
	if paddingX < 0 {
		paddingX = 0
	}
	if paddingY < 0 {
		paddingY = 0
	}

	vertPad := strings.Repeat("\n", paddingY)
	horizPad := strings.Repeat(" ", paddingX)

	return vertPad + horizPad + loadingBox
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

	content := styles.ModalTitleStyle().Render(welcome) + "\n\n" +
		styles.PreviewLabelStyle().Render("Welcome! Let's create your first board.") + "\n\n" +
		styles.HelpHintStyle().Render("Board name:") + "\n" +
		a.textInput.View() + "\n\n" +
		styles.HelpHintStyle().Render("Press Enter to continue (or leave blank for 'My Tasks')")

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

	modal := styles.ModalStyle().Width(modalWidth).Render(content)

	// Add vertical and horizontal padding
	vertPad := strings.Repeat("\n", paddingY)
	horizPad := strings.Repeat(" ", paddingX)

	return vertPad + horizPad + modal
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

// overlayDialog overlays a dialog popup on top of the background content
// using lipgloss v2's Layer and Canvas compositing
func overlayDialog(background, dialog string, width, height int) string {
	// Calculate dialog dimensions
	dialogWidth := lipgloss.Width(dialog)
	dialogHeight := lipgloss.Height(dialog)

	// Center the dialog
	x := (width - dialogWidth) / 2
	y := (height - dialogHeight) / 2
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	// Create layers
	bgLayer := lipgloss.NewLayer(background)
	dialogLayer := lipgloss.NewLayer(dialog).X(x).Y(y).Z(1)

	// Create compositor and render
	compositor := lipgloss.NewCompositor(bgLayer, dialogLayer)
	return compositor.Render()
}
