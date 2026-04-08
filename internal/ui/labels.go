package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/coliva/tsk/internal/model"
	"github.com/coliva/tsk/internal/styles"
)

// LabelEditor handles label management for a task
type LabelEditor struct {
	Task        *model.Task
	Board       *model.Board
	labels      []string
	input       textinput.Model
	selectedIdx int
	editing     bool
}

// NewLabelEditor creates a new label editor
func NewLabelEditor(task *model.Task, board *model.Board) *LabelEditor {
	labels := make([]string, len(task.Labels))
	copy(labels, task.Labels)

	input := textinput.New()
	input.Placeholder = "Add label..."
	input.CharLimit = 64

	return &LabelEditor{
		Task:        task,
		Board:       board,
		labels:      labels,
		input:       input,
		selectedIdx: -1, // -1 means input field is selected
		editing:     false,
	}
}

// Focus activates the label editor
func (le *LabelEditor) Focus() tea.Cmd {
	return le.input.Focus()
}

// Update handles input
func (le *LabelEditor) Update(msg tea.Msg) tea.Cmd {
	if le.editing {
		var cmd tea.Cmd
		le.input, cmd = le.input.Update(msg)
		return cmd
	}
	return nil
}

// HandleKey handles key input for navigation and actions
func (le *LabelEditor) HandleKey(key string) (done bool, cmd tea.Cmd) {
	switch key {
	case "esc":
		if le.editing {
			le.editing = false
			le.input.SetValue("")
			le.input.Blur()
			return false, nil
		}
		return true, nil // Exit label editor

	case "enter":
		if le.editing {
			// Add the label
			value := strings.TrimSpace(le.input.Value())
			if value != "" && !le.hasLabel(value) {
				le.labels = append(le.labels, value)
			}
			le.input.SetValue("")
			le.editing = false
			le.input.Blur()
			return false, nil
		}
		// Start editing if in the input area
		if le.selectedIdx == -1 {
			le.editing = true
			return false, le.input.Focus()
		}
		return false, nil

	case "d", "backspace":
		// Delete selected label
		if !le.editing && le.selectedIdx >= 0 && le.selectedIdx < len(le.labels) {
			le.labels = append(le.labels[:le.selectedIdx], le.labels[le.selectedIdx+1:]...)
			if le.selectedIdx >= len(le.labels) {
				le.selectedIdx = len(le.labels) - 1
			}
		}
		return false, nil

	case "j", "down":
		if !le.editing {
			if le.selectedIdx < len(le.labels)-1 {
				le.selectedIdx++
			} else {
				le.selectedIdx = -1 // Go back to input
			}
		}
		return false, nil

	case "k", "up":
		if !le.editing {
			if le.selectedIdx == -1 && len(le.labels) > 0 {
				le.selectedIdx = len(le.labels) - 1
			} else if le.selectedIdx > 0 {
				le.selectedIdx--
			}
		}
		return false, nil

	case "a":
		// Quick add mode
		if !le.editing && le.selectedIdx == -1 {
			le.editing = true
			return false, le.input.Focus()
		}
		return false, nil
	}

	return false, nil
}

// GetLabels returns the current labels
func (le *LabelEditor) GetLabels() []string {
	return le.labels
}

func (le *LabelEditor) hasLabel(label string) bool {
	for _, l := range le.labels {
		if strings.EqualFold(l, label) {
			return true
		}
	}
	return false
}

// View renders the label editor (full screen)
func (le *LabelEditor) View(width, height int) string {
	var lines []string

	lines = append(lines, styles.ModalTitleStyle.Render("Manage Labels"))
	lines = append(lines, "")

	// Current labels
	if len(le.labels) == 0 {
		lines = append(lines, styles.HelpHintStyle.Render("(No labels)"))
	} else {
		for i, labelName := range le.labels {
			var line string
			var labelDisplay string
			if le.Board != nil {
				label := le.Board.GetLabel(labelName)
				labelDisplay = styles.LabelBadge(label.Name, label.Color)
			} else {
				labelDisplay = labelName
			}
			if i == le.selectedIdx {
				line = styles.TaskSelectedStyle.Render(" ▶ ") + labelDisplay
			} else {
				line = "   " + labelDisplay
			}
			lines = append(lines, line)
		}
	}

	lines = append(lines, "")

	// Add label input
	inputLabel := "Add:"
	if le.selectedIdx == -1 && !le.editing {
		inputLabel = "▶ Add:"
	}
	lines = append(lines, styles.PreviewLabelStyle.Render(inputLabel))
	if le.editing {
		lines = append(lines, le.input.View())
	} else {
		lines = append(lines, styles.HelpHintStyle.Render("(press Enter or 'a' to add)"))
	}

	lines = append(lines, "")
	lines = append(lines, styles.HelpHintStyle.Render("j/k: navigate  Enter: save  d: delete label  Esc: done"))

	content := strings.Join(lines, "\n")

	// Full-screen layout
	editWidth := width - 4
	if editWidth < 50 {
		editWidth = 50
	}

	box := styles.ModalStyle.Width(editWidth).Height(height - 4).Render(content)

	return box
}
