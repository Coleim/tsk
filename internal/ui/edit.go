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
	titleInput  textinput.Model
	descInput   textinput.Model
	labelInput  textinput.Model
	activeField EditField
}

// NewTaskEdit creates a new task editor
func NewTaskEdit(task *model.Task) *TaskEdit {
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
	labelInput.Placeholder = "Labels (comma-separated)..."
	labelInput.CharLimit = 512

	return &TaskEdit{
		Task:        task,
		titleInput:  titleInput,
		descInput:   descInput,
		labelInput:  labelInput,
		activeField: EditFieldTitle,
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

// View renders the edit view
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
	lines = append(lines, "")

	// Help text
	lines = append(lines, styles.HelpHintStyle.Render("Tab: switch field  Enter: save  Esc: cancel"))

	content := strings.Join(lines, "\n")

	// Create modal box
	modalWidth := width - 20
	if modalWidth < 50 {
		modalWidth = 50
	}
	if modalWidth > 80 {
		modalWidth = 80
	}

	box := styles.ModalStyle.Width(modalWidth).Render(content)

	// Center vertically
	boxHeight := strings.Count(box, "\n") + 1
	paddingY := (height - boxHeight) / 2
	if paddingY < 0 {
		paddingY = 0
	}

	return strings.Repeat("\n", paddingY) + box
}
