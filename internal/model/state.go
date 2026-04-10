package model

import (
	"time"
)

// AppState holds the complete application state
type AppState struct {
	// Current board being viewed
	Board *Board

	// Current pane (status) being viewed
	CurrentPane Status

	// Index of selected task in current pane
	SelectedIndex int

	// Current input mode
	Mode Mode

	// Dirty flag - true if unsaved changes exist
	Dirty bool

	// Last save time
	LastSaved time.Time

	// Status bar message
	StatusMessage string

	// Status message timestamp (for auto-clear)
	StatusMessageTime time.Time

	// Terminal dimensions
	Width  int
	Height int

	// Help overlay visible
	ShowHelp bool

	// Modal state (stored as interface{} to avoid circular dependency)
	ActiveModal interface{}

	// Search query (when in search mode)
	SearchQuery string

	// Search results
	SearchResults []*Task

	// Filters
	FilterPriorities []Priority
	FilterLabels     []string

	// Scroll offset for task list
	ScrollOffset int
}

// NewAppState creates a new application state
func NewAppState() *AppState {
	return &AppState{
		CurrentPane:   StatusToDo,
		SelectedIndex: 0,
		Mode:          ModeNormal,
		Dirty:         false,
		FilterLabels:  []string{},
	}
}

// SetBoard sets the current board
func (s *AppState) SetBoard(board *Board) {
	s.Board = board
	s.CurrentPane = StatusToDo
	s.SelectedIndex = 0
	s.ScrollOffset = 0
	s.ClearSearch()
}

// CurrentTasks returns tasks in the current pane, applying any active filters
func (s *AppState) CurrentTasks() []*Task {
	if s.Board == nil {
		return nil
	}
	tasks := s.Board.TasksByStatus(s.CurrentPane)

	// Apply filters if any are active
	if len(s.FilterPriorities) == 0 && len(s.FilterLabels) == 0 {
		return tasks
	}

	var filtered []*Task
	for _, task := range tasks {
		if s.matchesFilters(task) {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

// matchesFilters checks if a task matches the current filters
func (s *AppState) matchesFilters(task *Task) bool {
	// Check priority filter (task must match ANY selected priority)
	if len(s.FilterPriorities) > 0 {
		matchesPriority := false
		for _, p := range s.FilterPriorities {
			if task.Priority == p {
				matchesPriority = true
				break
			}
		}
		if !matchesPriority {
			return false
		}
	}

	// Check label filters (task must have ALL selected labels)
	for _, filterLabel := range s.FilterLabels {
		found := false
		for _, taskLabel := range task.Labels {
			if taskLabel == filterLabel {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// HasActiveFilters returns true if any filters are active
func (s *AppState) HasActiveFilters() bool {
	return len(s.FilterPriorities) > 0 || len(s.FilterLabels) > 0
}

// SelectedTask returns the currently selected task
func (s *AppState) SelectedTask() *Task {
	tasks := s.CurrentTasks()
	if s.SelectedIndex >= 0 && s.SelectedIndex < len(tasks) {
		return tasks[s.SelectedIndex]
	}
	return nil
}

// SelectNext moves selection down
func (s *AppState) SelectNext() {
	tasks := s.CurrentTasks()
	if s.SelectedIndex < len(tasks)-1 {
		s.SelectedIndex++
	}
}

// SelectPrev moves selection up
func (s *AppState) SelectPrev() {
	if s.SelectedIndex > 0 {
		s.SelectedIndex--
	}
}

// SelectFirst moves selection to the first task (gg)
func (s *AppState) SelectFirst() {
	s.SelectedIndex = 0
}

// SelectLast moves selection to the last task (G)
func (s *AppState) SelectLast() {
	tasks := s.CurrentTasks()
	if len(tasks) > 0 {
		s.SelectedIndex = len(tasks) - 1
	}
}

// NextPane switches to the next pane
func (s *AppState) NextPane() {
	s.CurrentPane = s.CurrentPane.Next()
	s.SelectedIndex = 0
	s.ScrollOffset = 0
}

// PrevPane switches to the previous pane
func (s *AppState) PrevPane() {
	s.CurrentPane = s.CurrentPane.Prev()
	s.SelectedIndex = 0
	s.ScrollOffset = 0
}

// SetStatusMessage sets a status bar message
func (s *AppState) SetStatusMessage(msg string) {
	s.StatusMessage = msg
	s.StatusMessageTime = time.Now()
}

// ClearStatusMessageIfOld clears the status message if older than duration
func (s *AppState) ClearStatusMessageIfOld(duration time.Duration) {
	if s.StatusMessage != "" && time.Since(s.StatusMessageTime) > duration {
		s.StatusMessage = ""
	}
}

// MarkDirty marks the state as having unsaved changes
func (s *AppState) MarkDirty() {
	s.Dirty = true
}

// MarkClean marks the state as saved
func (s *AppState) MarkClean() {
	s.Dirty = false
	s.LastSaved = time.Now()
}

// ClearSearch clears search state
func (s *AppState) ClearSearch() {
	s.SearchQuery = ""
	s.SearchResults = nil
}

// UpdateSearch updates search results based on query
func (s *AppState) UpdateSearch(query string) {
	s.SearchQuery = query
	if s.Board != nil {
		s.SearchResults = s.Board.Search(query)
	}
}

// IsInDonePane returns true if currently viewing Done pane
func (s *AppState) IsInDonePane() bool {
	return s.CurrentPane == StatusDone
}

// ClampSelection ensures selected index is within bounds
func (s *AppState) ClampSelection() {
	tasks := s.CurrentTasks()
	if len(tasks) == 0 {
		s.SelectedIndex = 0
	} else if s.SelectedIndex >= len(tasks) {
		s.SelectedIndex = len(tasks) - 1
	}
}

// SwitchToTaskInPane switches to a pane and selects a task by ID
func (s *AppState) SwitchToTaskInPane(pane Status, taskID string) {
	s.CurrentPane = pane
	s.SelectedIndex = 0
	s.ScrollOffset = 0

	if s.Board == nil {
		return
	}

	// Find the task's position in the new pane
	tasks := s.Board.TasksByStatus(pane)
	for i, task := range tasks {
		if task.ID == taskID {
			s.SelectedIndex = i
			break
		}
	}
}
