package ui

import (
	"sort"
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
func NewFilter(board *model.Board, currentPriorities []model.Priority, currentLabels []string) *Filter {
	priorities := make(map[model.Priority]bool)
	labels := make(map[string]bool)

	// Initialize from current filters
	for _, p := range currentPriorities {
		priorities[p] = true
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

// GetSelectedPriorities returns all selected priority filters
func (f *Filter) GetSelectedPriorities() []model.Priority {
	var priorities []model.Priority
	for p, isSelected := range f.selectedPriorities {
		if isSelected {
			priorities = append(priorities, p)
		}
	}
	return priorities
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
	lines = append(lines, styles.ModalTitleStyle().Render("Filter Tasks"))
	lines = append(lines, "")

	// Active filter summary at top (prominent position)
	if f.HasFilters() {
		var filterParts []string
		// Iterate in consistent order (High, Medium, Low, None)
		for _, p := range []model.Priority{model.PriorityHigh, model.PriorityMedium, model.PriorityLow, model.PriorityNone} {
			if f.selectedPriorities[p] {
				filterParts = append(filterParts, p.String())
			}
		}
		// Collect and sort labels for consistent display
		var sortedLabels []string
		for label, isSelected := range f.selectedLabels {
			if isSelected {
				sortedLabels = append(sortedLabels, label)
			}
		}
		sort.Strings(sortedLabels)
		filterParts = append(filterParts, sortedLabels...)
		lines = append(lines, styles.SuccessStyle().Render("✓ Active: "+strings.Join(filterParts, ", ")))
		lines = append(lines, "")
	}

	// Priority section with section card styling
	priorityLines := []string{}
	priorityLines = append(priorityLines, styles.SectionCardTitleStyle().Render("Priority"))

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
		var checkbox string
		if f.selectedPriorities[p.p] {
			checkbox = styles.CheckboxCheckedStyle().Render("[✓]")
		} else {
			checkbox = styles.CheckboxUncheckedStyle().Render("[ ]")
		}

		var prefix string
		var nameDisplay string
		if f.mode == FilterModePriority && f.priorityIdx == i {
			prefix = styles.ActiveIndicator()
			nameDisplay = styles.FormFieldActiveLabelStyle().Render(p.name)
		} else {
			prefix = styles.InactiveIndicator()
			nameDisplay = styles.PriorityStyle(p.p).Render(p.name)
		}

		priorityLines = append(priorityLines, prefix+checkbox+" "+nameDisplay)
	}

	lines = append(lines, styles.SectionCardStyle().Render(strings.Join(priorityLines, "\n")))
	lines = append(lines, "")

	// Labels section with section card styling
	labelLines := []string{}
	labelLines = append(labelLines, styles.SectionCardTitleStyle().Render("Labels"))

	if len(f.availableLabels) == 0 {
		labelLines = append(labelLines, styles.HelpHintStyle().Render("No labels in board"))
	} else {
		for i, label := range f.availableLabels {
			var checkbox string
			if f.selectedLabels[label] {
				checkbox = styles.CheckboxCheckedStyle().Render("[✓]")
			} else {
				checkbox = styles.CheckboxUncheckedStyle().Render("[ ]")
			}

			var prefix string
			var labelDisplay string
			if f.mode == FilterModeLabels && f.labelIdx == i {
				prefix = styles.ActiveIndicator()
				labelDisplay = styles.FormFieldActiveLabelStyle().Render(label)
			} else {
				prefix = styles.InactiveIndicator()
				labelDisplay = styles.PopupItemStyle().Render(label)
			}

			labelLines = append(labelLines, prefix+checkbox+" "+labelDisplay)
		}
	}

	lines = append(lines, styles.SectionCardStyle().Render(strings.Join(labelLines, "\n")))
	lines = append(lines, "")

	// Separator line above hints
	lines = append(lines, styles.DialogSeparator(40))

	// Help text with consistent styling
	lines = append(lines, styles.KeyboardHintBarStyle().Render("j/k:navigate  Space/x:toggle  c:clear  Enter:apply  Esc:cancel"))

	content := strings.Join(lines, "\n")

	// Full-screen layout
	editWidth := width - 4
	if editWidth < 50 {
		editWidth = 50
	}

	box := styles.ModalStyle().Width(editWidth).Height(height - 4).Render(content)

	return box
}
