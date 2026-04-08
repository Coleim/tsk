package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/coliva/tsk/internal/model"
	"github.com/coliva/tsk/internal/styles"
)

// SearchDebounceMsg is sent after the debounce delay
type SearchDebounceMsg struct {
	Query string
}

// SearchResult groups tasks with their status for display
type SearchResult struct {
	Task   *model.Task
	Status model.Status
}

// Search handles search functionality
type Search struct {
	input          textinput.Model
	results        []SearchResult
	selectedIdx    int
	board          *model.Board
	lastQuery      string
	hasMoreResults bool
}

// NewSearch creates a new search component
func NewSearch(board *model.Board) *Search {
	input := textinput.New()
	input.Placeholder = "Search tasks..."
	input.CharLimit = 256
	input.Width = 50
	input.Focus()

	return &Search{
		input:       input,
		results:     nil,
		selectedIdx: 0,
		board:       board,
	}
}

// Focus activates the search input
func (s *Search) Focus() tea.Cmd {
	return s.input.Focus()
}

// Update handles input updates
func (s *Search) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)

	// Check if query changed
	query := s.input.Value()
	if query != s.lastQuery {
		s.lastQuery = query
		// Return debounce command
		return tea.Batch(cmd, s.debounce(query))
	}

	return cmd
}

// debounce returns a command that sends a SearchDebounceMsg after 50ms
func (s *Search) debounce(query string) tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
		return SearchDebounceMsg{Query: query}
	})
}

// PerformSearch executes the search and updates results
func (s *Search) PerformSearch(query string) {
	// Only search if this is the current query (debounce check)
	if query != s.lastQuery {
		return
	}

	if query == "" || s.board == nil {
		s.results = nil
		s.selectedIdx = 0
		s.hasMoreResults = false
		return
	}

	queryLower := strings.ToLower(query)
	var results []SearchResult
	count := 0
	maxResults := 100

	// Search all tasks, group by status
	for _, status := range model.AllStatuses() {
		tasks := s.board.TasksByStatus(status)
		for _, task := range tasks {
			if s.matchesQuery(task, queryLower) {
				results = append(results, SearchResult{
					Task:   task,
					Status: status,
				})
				count++
				if count >= maxResults {
					s.hasMoreResults = true
					break
				}
			}
		}
		if count >= maxResults {
			break
		}
	}

	s.results = results
	s.selectedIdx = 0
	if s.selectedIdx >= len(s.results) {
		s.selectedIdx = 0
	}
}

func (s *Search) matchesQuery(task *model.Task, query string) bool {
	if strings.Contains(strings.ToLower(task.Title), query) {
		return true
	}
	if strings.Contains(strings.ToLower(task.Description), query) {
		return true
	}
	for _, label := range task.Labels {
		if strings.Contains(strings.ToLower(label), query) {
			return true
		}
	}
	return false
}

// HandleKey handles navigation keys
func (s *Search) HandleKey(key string) (selectedTask *model.Task, done bool) {
	switch key {
	case "esc":
		return nil, true

	case "enter":
		if len(s.results) > 0 && s.selectedIdx < len(s.results) {
			return s.results[s.selectedIdx].Task, true
		}
		return nil, true

	case "j", "down":
		if s.selectedIdx < len(s.results)-1 {
			s.selectedIdx++
		}

	case "k", "up":
		if s.selectedIdx > 0 {
			s.selectedIdx--
		}

	case "ctrl+d":
		// Page down
		s.selectedIdx += 10
		if s.selectedIdx >= len(s.results) {
			s.selectedIdx = len(s.results) - 1
		}
		if s.selectedIdx < 0 {
			s.selectedIdx = 0
		}

	case "ctrl+u":
		// Page up
		s.selectedIdx -= 10
		if s.selectedIdx < 0 {
			s.selectedIdx = 0
		}
	}

	return nil, false
}

// SelectedTask returns the currently selected task
func (s *Search) SelectedTask() *model.Task {
	if len(s.results) > 0 && s.selectedIdx < len(s.results) {
		return s.results[s.selectedIdx].Task
	}
	return nil
}

// Query returns the current search query
func (s *Search) Query() string {
	return s.input.Value()
}

// View renders the search interface
func (s *Search) View(width, height int) string {
	var lines []string

	// Header
	lines = append(lines, styles.ModalTitleStyle.Render("Search Tasks"))
	lines = append(lines, "")

	// Search input
	lines = append(lines, s.input.View())
	lines = append(lines, "")

	// Results
	if len(s.results) == 0 {
		if s.input.Value() != "" {
			lines = append(lines, styles.HelpHintStyle.Render("No results found"))
		} else {
			lines = append(lines, styles.HelpHintStyle.Render("Type to search across title, description, and labels"))
		}
	} else {
		// Show result count
		countStr := fmt.Sprintf("%d result", len(s.results))
		if len(s.results) != 1 {
			countStr += "s"
		}
		if s.hasMoreResults {
			countStr = "100+ results (showing first 100)"
		}
		lines = append(lines, styles.HelpHintStyle.Render(countStr))
		lines = append(lines, "")

		// Calculate available height for results
		headerLines := 6 // title, empty, input, empty, count, empty
		footerLines := 3 // empty, help, padding
		availableHeight := height - headerLines - footerLines - 4 // 4 for modal borders
		if availableHeight < 5 {
			availableHeight = 5
		}

		// Group results by status
		var currentStatus model.Status = ""
		displayedLines := 0

		// Calculate scroll offset to keep selection visible
		startIdx := 0
		if s.selectedIdx >= availableHeight-2 {
			startIdx = s.selectedIdx - availableHeight + 3
		}

		for i, result := range s.results {
			if displayedLines >= availableHeight {
				lines = append(lines, styles.HelpHintStyle.Render(fmt.Sprintf("... and %d more", len(s.results)-i)))
				break
			}

			// Show status header when it changes
			if result.Status != currentStatus {
				if i >= startIdx {
					statusStyle := styles.StatusStyle(result.Status)
					lines = append(lines, statusStyle.Render(fmt.Sprintf("─── %s ───", result.Status)))
					displayedLines++
				}
				currentStatus = result.Status
			}

			if i < startIdx {
				continue
			}

			// Render task
			prefix := "  "
			taskStyle := styles.TaskNormalStyle
			if i == s.selectedIdx {
				prefix = "▶ "
				taskStyle = styles.TaskSelectedStyle
			}

			// Add priority indicator
			prioIndicator := ""
			if result.Task.Priority != model.PriorityNone {
				prioIndicator = styles.PriorityIndicator(result.Task.Priority) + " "
			}

			taskLine := prefix + prioIndicator + taskStyle.Render(result.Task.Title)
			lines = append(lines, taskLine)
			displayedLines++
		}
	}

	lines = append(lines, "")
	lines = append(lines, styles.HelpHintStyle.Render("j/k: navigate  Enter: go to task  Esc: cancel"))

	content := strings.Join(lines, "\n")

	// Create modal box
	modalWidth := width - 20
	if modalWidth < 60 {
		modalWidth = 60
	}
	if modalWidth > 100 {
		modalWidth = 100
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
