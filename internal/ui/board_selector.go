package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/coliva/tsk/internal/model"
	"github.com/coliva/tsk/internal/storage"
	"github.com/coliva/tsk/internal/styles"
)

// BoardSelectorMode defines the current board selector mode
type BoardSelectorMode int

const (
	BoardModeSelect BoardSelectorMode = iota
	BoardModeCreate
	BoardModeRename
	BoardModeDelete
)

// BoardSelector handles board selection and management
type BoardSelector struct {
	storage       *storage.Storage
	boards        []storage.BoardInfo
	selectedIdx   int
	mode          BoardSelectorMode
	input         textinput.Model
	currentBoard  *model.Board
	confirmDelete bool
}

// NewBoardSelector creates a new board selector
func NewBoardSelector(store *storage.Storage, currentBoard *model.Board) *BoardSelector {
	boards, _ := store.ListBoards()

	input := textinput.New()
	input.Placeholder = "Board name..."
	input.CharLimit = 64

	// Find current board index
	selectedIdx := 0
	if currentBoard != nil {
		for i, b := range boards {
			if b.ID == currentBoard.ID {
				selectedIdx = i
				break
			}
		}
	}

	return &BoardSelector{
		storage:      store,
		boards:       boards,
		selectedIdx:  selectedIdx,
		mode:         BoardModeSelect,
		input:        input,
		currentBoard: currentBoard,
	}
}

// SetMode changes the selector mode
func (bs *BoardSelector) SetMode(mode BoardSelectorMode) tea.Cmd {
	bs.mode = mode
	bs.confirmDelete = false

	switch mode {
	case BoardModeCreate:
		bs.input.SetValue("")
		bs.input.Placeholder = "New board name..."
		return bs.input.Focus()
	case BoardModeRename:
		if bs.selectedIdx < len(bs.boards) {
			bs.input.SetValue(bs.boards[bs.selectedIdx].Name)
		}
		bs.input.Placeholder = "Board name..."
		return bs.input.Focus()
	}
	return nil
}

// Update handles text input
func (bs *BoardSelector) Update(msg tea.Msg) tea.Cmd {
	if bs.mode == BoardModeCreate || bs.mode == BoardModeRename {
		var cmd tea.Cmd
		bs.input, cmd = bs.input.Update(msg)
		return cmd
	}
	return nil
}

// SelectNext moves selection down
func (bs *BoardSelector) SelectNext() {
	if bs.selectedIdx < len(bs.boards)-1 {
		bs.selectedIdx++
	}
}

// SelectPrev moves selection up
func (bs *BoardSelector) SelectPrev() {
	if bs.selectedIdx > 0 {
		bs.selectedIdx--
	}
}

// SelectedBoard returns the currently selected board info
func (bs *BoardSelector) SelectedBoard() *storage.BoardInfo {
	if bs.selectedIdx >= 0 && bs.selectedIdx < len(bs.boards) {
		return &bs.boards[bs.selectedIdx]
	}
	return nil
}

// GetInputValue returns the current input value
func (bs *BoardSelector) GetInputValue() string {
	return strings.TrimSpace(bs.input.Value())
}

// Mode returns the current mode
func (bs *BoardSelector) Mode() BoardSelectorMode {
	return bs.mode
}

// ConfirmDelete sets delete confirmation state
func (bs *BoardSelector) ConfirmDelete() {
	bs.confirmDelete = true
}

// IsConfirmingDelete returns true if awaiting delete confirmation
func (bs *BoardSelector) IsConfirmingDelete() bool {
	return bs.confirmDelete
}

// Refresh reloads the board list
func (bs *BoardSelector) Refresh() {
	bs.boards, _ = bs.storage.ListBoards()
	if bs.selectedIdx >= len(bs.boards) {
		bs.selectedIdx = len(bs.boards) - 1
	}
	if bs.selectedIdx < 0 {
		bs.selectedIdx = 0
	}
}

// View renders the board selector (full screen)
func (bs *BoardSelector) View(width, height int) string {
	var lines []string

	switch bs.mode {
	case BoardModeCreate:
		lines = append(lines, styles.ModalTitleStyle().Render("Create New Board"))
		lines = append(lines, "")

		// Form section with proper styling
		formLines := []string{}
		formLines = append(formLines, styles.ActiveIndicator()+styles.FormFieldActiveLabelStyle().Render("Name:"))
		formLines = append(formLines, "    "+bs.input.View())
		lines = append(lines, styles.SectionCardStyle().Render(strings.Join(formLines, "\n")))

		lines = append(lines, "")
		lines = append(lines, styles.DialogSeparator(40))
		lines = append(lines, styles.KeyboardHintBarStyle().Render("Enter:create  Esc:cancel"))

	case BoardModeRename:
		lines = append(lines, styles.ModalTitleStyle().Render("Rename Board"))
		lines = append(lines, "")

		// Form section with proper styling
		formLines := []string{}
		formLines = append(formLines, styles.ActiveIndicator()+styles.FormFieldActiveLabelStyle().Render("New name:"))
		formLines = append(formLines, "    "+bs.input.View())
		lines = append(lines, styles.SectionCardStyle().Render(strings.Join(formLines, "\n")))

		lines = append(lines, "")
		lines = append(lines, styles.DialogSeparator(40))
		lines = append(lines, styles.KeyboardHintBarStyle().Render("Enter:rename  Esc:cancel"))

	case BoardModeDelete:
		lines = append(lines, styles.ModalTitleStyle().Render("Delete Board"))
		lines = append(lines, "")
		if bs.SelectedBoard() != nil {
			boardName := bs.SelectedBoard().Name

			// Warning section
			warnLines := []string{}
			warnLines = append(warnLines, styles.ErrorStyle().Render("⚠ Delete '"+boardName+"'?"))
			warnLines = append(warnLines, "")
			warnLines = append(warnLines, styles.HelpHintStyle().Render("This cannot be undone!"))
			lines = append(lines, styles.SectionCardStyle().Render(strings.Join(warnLines, "\n")))

			lines = append(lines, "")
			if bs.confirmDelete {
				lines = append(lines, styles.ErrorStyle().Render("Press 'y' to confirm or Esc to cancel"))
			} else {
				lines = append(lines, styles.HelpHintStyle().Render("Press 'd' again to confirm, Esc to cancel"))
			}
		}

	default: // BoardModeSelect
		lines = append(lines, styles.ModalTitleStyle().Render("Switch Board"))
		lines = append(lines, "")

		if len(bs.boards) == 0 {
			lines = append(lines, styles.HelpHintStyle().Render("No boards found"))
		} else {
			// Each board as a card-style entry
			for i, board := range bs.boards {
				var cardContent string
				isCurrent := bs.currentBoard != nil && board.ID == bs.currentBoard.ID

				// Board name with current indicator
				var nameDisplay string
				if isCurrent {
					nameDisplay = "• " + board.Name
				} else {
					nameDisplay = "  " + board.Name
				}

				// Task count as secondary info
				taskInfo := styles.HelpHintStyle().Render(fmt.Sprintf("%d tasks", board.TaskCount))

				if i == bs.selectedIdx {
					// Selected board - use arrow indicator and accent color
					cardContent = styles.ActiveIndicator() + styles.FormFieldActiveLabelStyle().Render(nameDisplay) + "  " + taskInfo
					lines = append(lines, styles.TaskSelectedStyle().Render(cardContent))
				} else if isCurrent {
					// Current board but not selected - subtle highlight
					cardContent = styles.InactiveIndicator() + styles.PreviewValueStyle().Render(nameDisplay) + "  " + taskInfo
					lines = append(lines, cardContent)
				} else {
					// Normal board - no special styling
					cardContent = styles.InactiveIndicator() + styles.PopupItemStyle().Render(nameDisplay) + "  " + taskInfo
					lines = append(lines, cardContent)
				}
			}
		}

		lines = append(lines, "")
		lines = append(lines, styles.DialogSeparator(40))
		lines = append(lines, styles.KeyboardHintBarStyle().Render("j/k:navigate  Enter:switch  B:new  R:rename  D:delete  Esc:cancel"))
	}

	content := strings.Join(lines, "\n")

	// Full-screen layout
	editWidth := width - 4
	if editWidth < 50 {
		editWidth = 50
	}

	box := styles.ModalStyle().Width(editWidth).Height(height - 4).Render(content)

	return box
}
