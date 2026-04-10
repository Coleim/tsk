package ui

import (
	"strings"

	"github.com/coliva/tsk/internal/model"
	"github.com/coliva/tsk/internal/styles"
)

// SortSelector handles task sort mode selection
type SortSelector struct {
	selectedIdx int
	currentMode model.SortMode
	modes       []model.SortMode
}

// NewSortSelector creates a new sort selector component
func NewSortSelector(currentMode model.SortMode) *SortSelector {
	modes := model.AllSortModes()
	selectedIdx := 0
	// Set initial selection to current mode
	for i, m := range modes {
		if m == currentMode {
			selectedIdx = i
			break
		}
	}
	return &SortSelector{
		selectedIdx: selectedIdx,
		currentMode: currentMode,
		modes:       modes,
	}
}

// HandleKey handles key input for navigation and selection
func (s *SortSelector) HandleKey(key string) (done bool, apply bool) {
	switch key {
	case "esc":
		return true, false

	case "enter":
		return true, true

	case "j", "down":
		if s.selectedIdx < len(s.modes)-1 {
			s.selectedIdx++
		}

	case "k", "up":
		if s.selectedIdx > 0 {
			s.selectedIdx--
		}
	}

	return false, false
}

// GetSelectedMode returns the currently selected sort mode
func (s *SortSelector) GetSelectedMode() model.SortMode {
	if s.selectedIdx >= 0 && s.selectedIdx < len(s.modes) {
		return s.modes[s.selectedIdx]
	}
	return s.currentMode
}

// View renders the sort selector interface
func (s *SortSelector) View(width, height int) string {
	var lines []string

	// Header
	lines = append(lines, styles.ModalTitleStyle().Render("Sort Tasks"))
	lines = append(lines, "")

	// Current sort display
	lines = append(lines, styles.HelpHintStyle().Render("Current: "+s.currentMode.String()))
	lines = append(lines, "")

	// Sort options in section card
	optionLines := []string{}
	optionLines = append(optionLines, styles.SectionCardTitleStyle().Render("Sort By"))

	for i, mode := range s.modes {
		var prefix string
		var modeDisplay string

		// Check indicator
		var indicator string
		if mode == s.currentMode {
			indicator = styles.CheckboxCheckedStyle().Render("●")
		} else {
			indicator = styles.CheckboxUncheckedStyle().Render("○")
		}

		if i == s.selectedIdx {
			prefix = styles.ActiveIndicator()
			modeDisplay = styles.FormFieldActiveLabelStyle().Render(mode.String())
		} else {
			prefix = styles.InactiveIndicator()
			modeDisplay = styles.PopupItemStyle().Render(mode.String())
		}

		optionLines = append(optionLines, prefix+indicator+" "+modeDisplay)
	}

	lines = append(lines, styles.SectionCardStyle().Render(strings.Join(optionLines, "\n")))
	lines = append(lines, "")

	// Separator line above hints
	lines = append(lines, styles.DialogSeparator(40))

	// Help text
	lines = append(lines, styles.KeyboardHintBarStyle().Render("j/k:navigate  Enter:apply  Esc:cancel"))

	content := strings.Join(lines, "\n")

	// Full-screen layout
	editWidth := width - 4
	if editWidth < 50 {
		editWidth = 50
	}

	box := styles.ModalStyle().Width(editWidth).Height(height - 4).Render(content)

	return box
}
