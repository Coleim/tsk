package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/coliva/tsk/internal/model"
	"github.com/coliva/tsk/internal/styles"
)

// EditField represents which field is being edited
type EditField int

const (
	EditFieldTitle EditField = iota
	EditFieldDescription
	EditFieldLabels
)

// TaskEdit represents the task editing state
type TaskEdit struct {
	Task             *model.Task
	Board            *model.Board
	titleInput       textinput.Model
	descInput        textinput.Model
	labelInput       textinput.Model
	activeField      EditField
	showLabelPopup   bool // Whether label selection popup is visible
	popupSelectedIdx int  // Currently selected index in popup
}

// NewTaskEdit creates a new task editor
func NewTaskEdit(task *model.Task, board *model.Board) *TaskEdit {
	titleInput := textinput.New()
	titleInput.SetValue(task.Title)
	titleInput.Placeholder = "Task title..."
	titleInput.CharLimit = 256
	titleInput.Focus()

	descInput := textinput.New()
	descInput.SetValue(task.Description)
	descInput.Placeholder = "Description (optional)..."
	descInput.CharLimit = 1024

	labelInput := textinput.New()
	labelInput.SetValue(strings.Join(task.Labels, ", "))
	labelInput.Placeholder = "Labels (Tab to autocomplete)..."
	labelInput.CharLimit = 512

	return &TaskEdit{
		Task:             task,
		Board:            board,
		titleInput:       titleInput,
		descInput:        descInput,
		labelInput:       labelInput,
		activeField:      EditFieldTitle,
		showLabelPopup:   false,
		popupSelectedIdx: 0,
	}
}

// Focus activates the text inputs
func (te *TaskEdit) Focus() tea.Cmd {
	switch te.activeField {
	case EditFieldTitle:
		return te.titleInput.Focus()
	case EditFieldDescription:
		return te.descInput.Focus()
	case EditFieldLabels:
		return te.labelInput.Focus()
	}
	return te.titleInput.Focus()
}

// Update handles text input updates
func (te *TaskEdit) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch te.activeField {
	case EditFieldTitle:
		te.titleInput, cmd = te.titleInput.Update(msg)
	case EditFieldDescription:
		te.descInput, cmd = te.descInput.Update(msg)
	case EditFieldLabels:
		te.labelInput, cmd = te.labelInput.Update(msg)
	}

	return cmd
}

// OpenLabelPopup opens the label selection popup
func (te *TaskEdit) OpenLabelPopup() {
	if te.Board == nil {
		return
	}
	labels := te.Board.AllLabels()
	if len(labels) == 0 {
		return
	}
	te.showLabelPopup = true
	te.popupSelectedIdx = 0
}

// CloseLabelPopup closes the label selection popup without selecting
func (te *TaskEdit) CloseLabelPopup() {
	te.showLabelPopup = false
}

// IsPopupOpen returns true if the label popup is currently visible
func (te *TaskEdit) IsPopupOpen() bool {
	return te.showLabelPopup
}

// NextPopupLabel moves selection to next label in popup
func (te *TaskEdit) NextPopupLabel() {
	if te.Board == nil {
		return
	}
	labels := te.Board.AllLabels()
	if len(labels) == 0 {
		return
	}
	te.popupSelectedIdx = (te.popupSelectedIdx + 1) % len(labels)
}

// PrevPopupLabel moves selection to previous label in popup
func (te *TaskEdit) PrevPopupLabel() {
	if te.Board == nil {
		return
	}
	labels := te.Board.AllLabels()
	if len(labels) == 0 {
		return
	}
	te.popupSelectedIdx--
	if te.popupSelectedIdx < 0 {
		te.popupSelectedIdx = len(labels) - 1
	}
}

// SelectPopupLabel adds the currently selected label to the field and closes popup
func (te *TaskEdit) SelectPopupLabel() {
	if te.Board == nil || !te.showLabelPopup {
		return
	}
	labels := te.Board.AllLabels()
	if len(labels) == 0 || te.popupSelectedIdx >= len(labels) {
		te.showLabelPopup = false
		return
	}

	newLabel := labels[te.popupSelectedIdx]
	current := te.labelInput.Value()

	// Check if label already exists
	existing := te.GetLabels()
	alreadyHas := false
	for _, l := range existing {
		if l == newLabel {
			alreadyHas = true
			break
		}
	}

	if !alreadyHas {
		if current == "" {
			te.labelInput.SetValue(newLabel)
		} else {
			te.labelInput.SetValue(current + ", " + newLabel)
		}
	}

	te.showLabelPopup = false
}

// NextField moves to the next field
func (te *TaskEdit) NextField() tea.Cmd {
	switch te.activeField {
	case EditFieldTitle:
		te.activeField = EditFieldDescription
		te.titleInput.Blur()
		return te.descInput.Focus()
	case EditFieldDescription:
		te.activeField = EditFieldLabels
		te.descInput.Blur()
		return te.labelInput.Focus()
	case EditFieldLabels:
		// Wrap around
		te.activeField = EditFieldTitle
		te.labelInput.Blur()
		return te.titleInput.Focus()
	}
	return nil
}

// PrevField moves to the previous field
func (te *TaskEdit) PrevField() tea.Cmd {
	switch te.activeField {
	case EditFieldTitle:
		// Wrap around to labels
		te.activeField = EditFieldLabels
		te.titleInput.Blur()
		return te.labelInput.Focus()
	case EditFieldDescription:
		te.activeField = EditFieldTitle
		te.descInput.Blur()
		return te.titleInput.Focus()
	case EditFieldLabels:
		te.activeField = EditFieldDescription
		te.labelInput.Blur()
		return te.descInput.Focus()
	}
	return nil
}

// GetTitle returns the edited title
func (te *TaskEdit) GetTitle() string {
	return strings.TrimSpace(te.titleInput.Value())
}

// GetDescription returns the edited description
func (te *TaskEdit) GetDescription() string {
	return strings.TrimSpace(te.descInput.Value())
}

// GetLabels returns the edited labels as a slice
func (te *TaskEdit) GetLabels() []string {
	input := strings.TrimSpace(te.labelInput.Value())
	if input == "" {
		return nil
	}

	parts := strings.Split(input, ",")
	var labels []string
	for _, p := range parts {
		label := strings.TrimSpace(p)
		if label != "" {
			labels = append(labels, label)
		}
	}
	return labels
}

// View renders the edit view (full screen)
func (te *TaskEdit) View(width, height int) string {
	var lines []string

	// Header
	lines = append(lines, styles.ModalTitleStyle.Render("Edit Task"))
	lines = append(lines, "")

	// Title field
	titleLabel := "Title:"
	if te.activeField == EditFieldTitle {
		titleLabel = "▶ Title:"
	}
	lines = append(lines, styles.PreviewLabelStyle.Render(titleLabel))
	lines = append(lines, te.titleInput.View())
	lines = append(lines, "")

	// Description field
	descLabel := "Description:"
	if te.activeField == EditFieldDescription {
		descLabel = "▶ Description:"
	}
	lines = append(lines, styles.PreviewLabelStyle.Render(descLabel))
	lines = append(lines, te.descInput.View())
	lines = append(lines, "")

	// Labels field
	labelsLabel := "Labels:"
	if te.activeField == EditFieldLabels {
		labelsLabel = "▶ Labels:"
	}
	lines = append(lines, styles.PreviewLabelStyle.Render(labelsLabel))
	lines = append(lines, te.labelInput.View())

	// Show popup or hint based on state
	if te.showLabelPopup && te.Board != nil {
		// Render the label popup
		lines = append(lines, te.renderLabelPopup())
	} else if te.activeField == EditFieldLabels {
		lines = append(lines, styles.HelpHintStyle.Render("  (Tab to open label picker)"))
	}
	lines = append(lines, "")

	// Help text
	if te.showLabelPopup {
		lines = append(lines, styles.HelpHintStyle.Render("Tab: next  Shift+Tab: prev  Enter: select  Esc: close"))
	} else if te.activeField == EditFieldLabels {
		lines = append(lines, styles.HelpHintStyle.Render("Tab: open labels  Shift+Tab: prev field  Enter: save  Esc: cancel"))
	} else {
		lines = append(lines, styles.HelpHintStyle.Render("Tab: next field  Shift+Tab: prev field  Enter: save  Esc: cancel"))
	}

	content := strings.Join(lines, "\n")

	// Full-screen layout
	editWidth := width - 4
	if editWidth < 50 {
		editWidth = 50
	}

	box := styles.ModalStyle.Width(editWidth).Height(height - 4).Render(content)

	return box
}

// IsLabelsField returns true if the labels field is active
func (te *TaskEdit) IsLabelsField() bool {
	return te.activeField == EditFieldLabels
}

// renderLabelPopup renders the label selection popup using Lipgloss
func (te *TaskEdit) renderLabelPopup() string {
	if te.Board == nil {
		return ""
	}

	labels := te.Board.AllLabels()
	if len(labels) == 0 {
		return styles.HelpHintStyle.Render("  (No labels available)")
	}

	var items []string
	items = append(items, styles.PopupTitleStyle.Render("Select Label"))

	for i, labelName := range labels {
		label := te.Board.GetLabel(labelName)
		badge := styles.LabelBadge(label.Name, label.Color)

		if i == te.popupSelectedIdx {
			// Highlight selected item with background
			items = append(items, styles.PopupSelectedItemStyle.Render("▶ "+badge))
		} else {
			items = append(items, styles.PopupItemStyle.Render("  "+badge))
		}
	}

	content := strings.Join(items, "\n")
	return styles.PopupStyle.Render(content)
}
