package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/coliva/tsk/internal/model"
	"github.com/coliva/tsk/internal/styles"
)

// SimpleSearch is a compact search popup
type SimpleSearch struct {
	input       textinput.Model
	results     []*model.Task
	selectedIdx int
	board       *model.Board
}

// NewSimpleSearch creates a new simple search component
func NewSimpleSearch(board *model.Board) *SimpleSearch {
	ti := textinput.New()
	ti.Placeholder = "Type to search..."
	ti.CharLimit = 100
	ti.Width = 40
	ti.Focus()

	return &SimpleSearch{
		input:       ti,
		results:     nil,
		selectedIdx: 0,
		board:       board,
	}
}

// Init returns the initial command
func (s *SimpleSearch) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages
func (s *SimpleSearch) Update(msg tea.Msg) (bool, *model.Task, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return true, nil, nil
		case "enter":
			if len(s.results) > 0 && s.selectedIdx < len(s.results) {
				return true, s.results[s.selectedIdx], nil
			}
			return true, nil, nil
		case "up", "ctrl+p":
			if s.selectedIdx > 0 {
				s.selectedIdx--
			}
			return false, nil, nil
		case "down", "ctrl+n":
			if s.selectedIdx < len(s.results)-1 {
				s.selectedIdx++
			}
			return false, nil, nil
		}
	}

	// Update text input
	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)

	// Search on every keystroke
	s.doSearch()

	return false, nil, cmd
}

func (s *SimpleSearch) doSearch() {
	query := strings.ToLower(strings.TrimSpace(s.input.Value()))
	if query == "" || s.board == nil {
		s.results = nil
		s.selectedIdx = 0
		return
	}

	var results []*model.Task
	for _, status := range model.AllStatuses() {
		for _, task := range s.board.TasksByStatus(status) {
			// Title: starts-with matching
			if strings.HasPrefix(strings.ToLower(task.Title), query) {
				results = append(results, task)
				if len(results) >= 8 {
					break
				}
			}
		}
		if len(results) >= 8 {
			break
		}
	}

	s.results = results
	if s.selectedIdx >= len(results) {
		s.selectedIdx = 0
	}
}

// View renders the search popup
func (s *SimpleSearch) View() string {
	var b strings.Builder

	// Title
	b.WriteString(styles.ModalTitleStyle().Render("Search"))
	b.WriteString("\n\n")

	// Input
	b.WriteString(s.input.View())
	b.WriteString("\n")

	// Results
	if len(s.results) == 0 {
		if s.input.Value() != "" {
			b.WriteString("\n")
			b.WriteString(styles.HelpHintStyle().Render("No matches"))
		}
	} else {
		b.WriteString("\n")
		for i, task := range s.results {
			prefix := "  "
			style := styles.TaskNormalStyle()
			if i == s.selectedIdx {
				prefix = "> "
				style = styles.TaskSelectedStyle()
			}
			title := task.Title
			if len(title) > 35 {
				title = title[:32] + "..."
			}
			fmt.Fprintf(&b, "%s%s\n", prefix, style.Render(title))
		}
	}

	b.WriteString("\n")
	b.WriteString(styles.HelpHintStyle().Render("↑↓:select Enter:go Esc:close"))

	// Create popup box
	return styles.ModalStyle().
		Width(50).
		Render(b.String())
}
