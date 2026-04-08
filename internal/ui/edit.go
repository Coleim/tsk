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
	Task        *model.Task
	Board       *model.Board
	titleInput  textinput.Model
	descInput   textinput.Model
	labelInput  textinput.Model
	activeField EditField
	labelIdx    int // Index for cycling through labels with Tab
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
		Task:        task,
		Board:       board,
		titleInput:  titleInput,
		descInput:   descInput,
		labelInput:  labelInput,
		activeField: EditFieldTitle,
		labelIdx:    -1,
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

// CycleLabel cycles through existing board labels when Tab is pressed on labels field
func (te *TaskEdit) CycleLabel() {
	if te.Board == nil {
		return
	}
	labels := te.Board.AllLabels()
	if len(labels) == 0 {
		return
	}
	te.labelIdx = (te.labelIdx + 1) % len(labels)

	// Get current labels and append the new one
	current := te.labelInput.Value()
	newLabel := labels[te.labelIdx]

	if current == "" {
		te.labelInput.SetValue(newLabel)
	} else {
		// Check if label already exists
		existing := te.GetLabels()
		for _, l := range existing {
			if l == newLabel {
				return // Already have this label
			}
		}
		te.labelInput.SetValue(current + ", " + newLabel)
	}
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
	if te.activeField == EditFieldLabels {
		lines = append(lines, te.labelInput.View())
		lines = append(lines, styles.HelpHintStyle.Render("  (Tab to cycle through existing labels)"))
	} else {
		lines = append(lines, te.labelInput.View())
	}
	lines = append(lines, "")

	// Help text
	if te.activeField == EditFieldLabels {
		lines = append(lines, styles.HelpHintStyle.Render("Tab: cycle labels  Shift+Tab: prev field  Enter: save  Esc: cancel"))
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
