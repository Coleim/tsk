package ui

import (
	"fmt"
	"strings"

	"github.com/coliva/tsk/internal/model"
	"github.com/coliva/tsk/internal/styles"
)

// FilterMode represents the current filter selection area
type FilterMode int

const (
	FilterModePriority FilterMode = iota
	FilterModeLabels
)

// Filter handles task filtering functionality
type Filter struct {
	board              *model.Board
	mode               FilterMode
	priorityIdx        int
	labelIdx           int
	selectedPriorities map[model.Priority]bool
	selectedLabels     map[string]bool
	availableLabels    []string
}

// NewFilter creates a new filter component
func NewFilter(board *model.Board, currentPriority *model.Priority, currentLabels []string) *Filter {
	priorities := make(map[model.Priority]bool)
	labels := make(map[string]bool)

	// Initialize from current filters
	if currentPriority != nil {
		priorities[*currentPriority] = true
	}
	for _, label := range currentLabels {
		labels[label] = true
	}

	// Get available labels from board
	availableLabels := board.AllLabels()

	return &Filter{
		board:              board,
		mode:               FilterModePriority,
		priorityIdx:        0,
		labelIdx:           0,
		selectedPriorities: priorities,
		selectedLabels:     labels,
		availableLabels:    availableLabels,
	}
}

// HandleKey handles key input for navigation and toggling
func (f *Filter) HandleKey(key string) (done bool, apply bool) {
	switch key {
	case "esc":
		return true, false

	case "enter":
		return true, true

	case "j", "down":
		if f.mode == FilterModePriority {
			if f.priorityIdx < 3 { // 4 priorities (High, Med, Low, None)
				f.priorityIdx++
			} else {
				// Move to labels
				f.mode = FilterModeLabels
				f.labelIdx = 0
			}
		} else {
			if f.labelIdx < len(f.availableLabels)-1 {
				f.labelIdx++
			}
		}

	case "k", "up":
		if f.mode == FilterModeLabels {
			if f.labelIdx > 0 {
				f.labelIdx--
			} else {
				// Move to priorities
				f.mode = FilterModePriority
				f.priorityIdx = 3
			}
		} else {
			if f.priorityIdx > 0 {
				f.priorityIdx--
			}
		}

	case "tab":
		// Toggle between priority and labels sections
		if f.mode == FilterModePriority {
			f.mode = FilterModeLabels
			f.labelIdx = 0
		} else {
			f.mode = FilterModePriority
			f.priorityIdx = 0
		}

	case " ", "x":
		// Toggle selection
		if f.mode == FilterModePriority {
			priority := f.indexToPriority(f.priorityIdx)
			f.selectedPriorities[priority] = !f.selectedPriorities[priority]
		} else if len(f.availableLabels) > 0 {
			label := f.availableLabels[f.labelIdx]
			f.selectedLabels[label] = !f.selectedLabels[label]
		}

	case "c":
		// Clear all filters
		f.selectedPriorities = make(map[model.Priority]bool)
		f.selectedLabels = make(map[string]bool)
	}

	return false, false
}

func (f *Filter) indexToPriority(idx int) model.Priority {
	switch idx {
	case 0:
		return model.PriorityHigh
	case 1:
		return model.PriorityMedium
	case 2:
		return model.PriorityLow
	default:
		return model.PriorityNone
	}
}

// GetSelectedPriority returns the selected priority filter (nil if multiple or none)
func (f *Filter) GetSelectedPriority() *model.Priority {
	var selected *model.Priority
	count := 0
	for p, isSelected := range f.selectedPriorities {
		if isSelected {
			count++
			pCopy := p
			selected = &pCopy
		}
	}
	// Only return if exactly one priority is selected
	if count == 1 {
		return selected
	}
	return nil
}

// GetSelectedLabels returns the selected label filters
func (f *Filter) GetSelectedLabels() []string {
	var labels []string
	for label, isSelected := range f.selectedLabels {
		if isSelected {
			labels = append(labels, label)
		}
	}
	return labels
}

// HasFilters returns true if any filters are active
func (f *Filter) HasFilters() bool {
	for _, isSelected := range f.selectedPriorities {
		if isSelected {
			return true
		}
	}
	for _, isSelected := range f.selectedLabels {
		if isSelected {
			return true
		}
	}
	return false
}

// View renders the filter interface
func (f *Filter) View(width, height int) string {
	var lines []string

	// Header
	lines = append(lines, styles.ModalTitleStyle.Render("Filter Tasks"))
	lines = append(lines, "")

	// Priority section
	lines = append(lines, styles.PreviewLabelStyle.Render("Priority:"))
	priorities := []struct {
		p    model.Priority
		name string
	}{
		{model.PriorityHigh, "High"},
		{model.PriorityMedium, "Medium"},
		{model.PriorityLow, "Low"},
		{model.PriorityNone, "None"},
	}

	for i, p := range priorities {
		checkbox := "[ ]"
		if f.selectedPriorities[p.p] {
			checkbox = "[x]"
		}

		prefix := "  "
		style := styles.TaskNormalStyle
		if f.mode == FilterModePriority && f.priorityIdx == i {
			prefix = "▶ "
			style = styles.TaskSelectedStyle
		}

		prioStyle := styles.PriorityStyle(p.p)
		line := prefix + checkbox + " " + prioStyle.Render(p.name)
		if f.mode == FilterModePriority && f.priorityIdx == i {
			line = style.Render(prefix + checkbox + " " + p.name)
		}
		lines = append(lines, line)
	}

	lines = append(lines, "")

	// Labels section
	lines = append(lines, styles.PreviewLabelStyle.Render("Labels:"))
	if len(f.availableLabels) == 0 {
		lines = append(lines, styles.HelpHintStyle.Render("  No labels in board"))
	} else {
		for i, label := range f.availableLabels {
			checkbox := "[ ]"
			if f.selectedLabels[label] {
				checkbox = "[x]"
			}

			prefix := "  "
			style := styles.TaskNormalStyle
			if f.mode == FilterModeLabels && f.labelIdx == i {
				prefix = "▶ "
				style = styles.TaskSelectedStyle
			}

			line := prefix + checkbox + " " + style.Render(label)
			lines = append(lines, line)
		}
	}

	lines = append(lines, "")

	// Active filter summary
	if f.HasFilters() {
		var filterParts []string
		for p, isSelected := range f.selectedPriorities {
			if isSelected {
				filterParts = append(filterParts, p.String())
			}
		}
		for label, isSelected := range f.selectedLabels {
			if isSelected {
				filterParts = append(filterParts, label)
			}
		}
		lines = append(lines, styles.SuccessStyle.Render(fmt.Sprintf("Active: %s", strings.Join(filterParts, ", "))))
		lines = append(lines, "")
	}

	// Help text
	lines = append(lines, styles.HelpHintStyle.Render("j/k: navigate  Space/x: toggle  c: clear  Enter: apply  Esc: cancel"))

	content := strings.Join(lines, "\n")

	// Create modal box
	modalWidth := width - 20
	if modalWidth < 50 {
		modalWidth = 50
	}
	if modalWidth > 70 {
		modalWidth = 70
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
