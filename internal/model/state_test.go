package model

import (
	"testing"
	"time"
)

func TestSwitchToTaskInPane(t *testing.T) {
	// Create a board with tasks in different statuses
	board := NewBoard("test-board", "Test Board")

	task1 := NewTask("task-1", "Task 1", StatusToDo)
	board.AddTask(task1)

	task2 := NewTask("task-2", "Task 2", StatusInProgress)
	board.AddTask(task2)

	task3 := NewTask("task-3", "Task 3", StatusInProgress)
	board.AddTask(task3)

	task4 := NewTask("task-4", "Task 4", StatusDone)
	board.AddTask(task4)

	state := &AppState{
		Board:         board,
		CurrentPane:   StatusToDo,
		SelectedIndex: 0,
	}

	// Test: switch to InProgress pane and select task2
	state.SwitchToTaskInPane(StatusInProgress, task2.ID)

	if state.CurrentPane != StatusInProgress {
		t.Errorf("Expected CurrentPane to be %v, got %v", StatusInProgress, state.CurrentPane)
	}
	if state.SelectedIndex != 0 {
		t.Errorf("Expected SelectedIndex to be 0 (task2 is first), got %d", state.SelectedIndex)
	}

	// Test: switch to InProgress pane and select task3 (second task)
	state.SwitchToTaskInPane(StatusInProgress, task3.ID)

	if state.CurrentPane != StatusInProgress {
		t.Errorf("Expected CurrentPane to be %v, got %v", StatusInProgress, state.CurrentPane)
	}
	if state.SelectedIndex != 1 {
		t.Errorf("Expected SelectedIndex to be 1 (task3 is second), got %d", state.SelectedIndex)
	}

	// Test: switch to Done pane and select task4
	state.SwitchToTaskInPane(StatusDone, task4.ID)

	if state.CurrentPane != StatusDone {
		t.Errorf("Expected CurrentPane to be %v, got %v", StatusDone, state.CurrentPane)
	}
	if state.SelectedIndex != 0 {
		t.Errorf("Expected SelectedIndex to be 0 (task4 is first in Done), got %d", state.SelectedIndex)
	}

	// Test: switch with non-existent task ID defaults to index 0
	state.SwitchToTaskInPane(StatusToDo, "non-existent-id")

	if state.CurrentPane != StatusToDo {
		t.Errorf("Expected CurrentPane to be %v, got %v", StatusToDo, state.CurrentPane)
	}
	if state.SelectedIndex != 0 {
		t.Errorf("Expected SelectedIndex to be 0 for non-existent task, got %d", state.SelectedIndex)
	}
}

func TestClampSelection(t *testing.T) {
	board := NewBoard("test-board", "Test Board")
	task1 := NewTask("task-1", "Task 1", StatusToDo)
	board.AddTask(task1)

	state := &AppState{
		Board:         board,
		CurrentPane:   StatusToDo,
		SelectedIndex: 10, // Out of bounds
	}

	state.ClampSelection()

	if state.SelectedIndex != 0 {
		t.Errorf("Expected SelectedIndex to be clamped to 0, got %d", state.SelectedIndex)
	}

	// Test with empty pane
	state.CurrentPane = StatusInProgress
	state.SelectedIndex = 5
	state.ClampSelection()

	if state.SelectedIndex != 0 {
		t.Errorf("Expected SelectedIndex to be 0 for empty pane, got %d", state.SelectedIndex)
	}
}

// ============ Mode Transition Tests (14.15) ============

func TestModeTransitions(t *testing.T) {
	state := NewAppState()

	// Initial mode should be Normal
	if state.Mode != ModeNormal {
		t.Errorf("Expected initial mode to be ModeNormal, got %v", state.Mode)
	}

	// Transition to Insert mode
	state.Mode = ModeInsert
	if state.Mode != ModeInsert {
		t.Errorf("Expected ModeInsert, got %v", state.Mode)
	}
	if state.Mode.String() != "INSERT" {
		t.Errorf("Expected mode string 'INSERT', got %s", state.Mode.String())
	}

	// Transition back to Normal
	state.Mode = ModeNormal
	if state.Mode != ModeNormal {
		t.Errorf("Expected ModeNormal, got %v", state.Mode)
	}
	if state.Mode.String() != "NORMAL" {
		t.Errorf("Expected mode string 'NORMAL', got %s", state.Mode.String())
	}

	// Transition to Search mode
	state.Mode = ModeSearch
	if state.Mode != ModeSearch {
		t.Errorf("Expected ModeSearch, got %v", state.Mode)
	}
	if state.Mode.String() != "SEARCH" {
		t.Errorf("Expected mode string 'SEARCH', got %s", state.Mode.String())
	}

	// Transition to Modal mode
	state.Mode = ModeModal
	if state.Mode != ModeModal {
		t.Errorf("Expected ModeModal, got %v", state.Mode)
	}
	if state.Mode.String() != "MODAL" {
		t.Errorf("Expected mode string 'MODAL', got %s", state.Mode.String())
	}
}

func TestModeAllStrings(t *testing.T) {
	tests := []struct {
		mode     Mode
		expected string
	}{
		{ModeNormal, "NORMAL"},
		{ModeInsert, "INSERT"},
		{ModeSearch, "SEARCH"},
		{ModeModal, "MODAL"},
		{ModeWelcome, "WELCOME"},
		{ModeDetail, "DETAIL"},
		{ModeEdit, "EDIT"},
		{ModeLabels, "LABELS"},
		{ModeDueDate, "DUE DATE"},
		{ModeBoard, "BOARD"},
		{ModeFilter, "FILTER"},
	}

	for _, tt := range tests {
		if got := tt.mode.String(); got != tt.expected {
			t.Errorf("Mode(%d).String() = %s, want %s", tt.mode, got, tt.expected)
		}
	}
}

func TestNewAppState(t *testing.T) {
	state := NewAppState()

	if state.CurrentPane != StatusToDo {
		t.Errorf("Expected CurrentPane StatusToDo, got %v", state.CurrentPane)
	}
	if state.SelectedIndex != 0 {
		t.Errorf("Expected SelectedIndex 0, got %d", state.SelectedIndex)
	}
	if state.Mode != ModeNormal {
		t.Errorf("Expected Mode ModeNormal, got %v", state.Mode)
	}
	if state.Dirty {
		t.Error("Expected Dirty to be false")
	}
	if state.FilterLabels == nil {
		t.Error("Expected FilterLabels to be initialized")
	}
}

func TestSetBoard(t *testing.T) {
	state := NewAppState()
	board := NewBoard("test-board", "Test Board")
	board.AddTask(NewTask("task-1", "Task 1", StatusInProgress))

	// Set a different pane first
	state.CurrentPane = StatusDone
	state.SelectedIndex = 5
	state.SearchQuery = "test"

	// Set new board
	state.SetBoard(board)

	if state.Board != board {
		t.Error("Board not set correctly")
	}
	if state.CurrentPane != StatusToDo {
		t.Errorf("Expected CurrentPane reset to StatusToDo, got %v", state.CurrentPane)
	}
	if state.SelectedIndex != 0 {
		t.Errorf("Expected SelectedIndex reset to 0, got %d", state.SelectedIndex)
	}
	if state.SearchQuery != "" {
		t.Errorf("Expected SearchQuery to be cleared, got %s", state.SearchQuery)
	}
}

func TestCurrentTasks(t *testing.T) {
	state := NewAppState()
	board := NewBoard("test-board", "Test Board")

	task1 := NewTask("task-1", "Task 1", StatusToDo)
	task2 := NewTask("task-2", "Task 2", StatusToDo)
	task3 := NewTask("task-3", "Task 3", StatusInProgress)

	board.AddTask(task1)
	board.AddTask(task2)
	board.AddTask(task3)

	state.SetBoard(board)

	// Test ToDo pane
	tasks := state.CurrentTasks()
	if len(tasks) != 2 {
		t.Errorf("Expected 2 ToDo tasks, got %d", len(tasks))
	}

	// Test InProgress pane
	state.CurrentPane = StatusInProgress
	tasks = state.CurrentTasks()
	if len(tasks) != 1 {
		t.Errorf("Expected 1 InProgress task, got %d", len(tasks))
	}

	// Test Empty pane
	state.CurrentPane = StatusDone
	tasks = state.CurrentTasks()
	if len(tasks) != 0 {
		t.Errorf("Expected 0 Done tasks, got %d", len(tasks))
	}
}

func TestCurrentTasksWithFilters(t *testing.T) {
	state := NewAppState()
	board := NewBoard("test-board", "Test Board")

	task1 := NewTask("task-1", "High Priority", StatusToDo)
	task1.Priority = PriorityHigh

	task2 := NewTask("task-2", "Low Priority", StatusToDo)
	task2.Priority = PriorityLow

	task3 := NewTask("task-3", "No Priority", StatusToDo)
	task3.Priority = PriorityNone

	board.AddTask(task1)
	board.AddTask(task2)
	board.AddTask(task3)

	state.SetBoard(board)

	// Without filter - all tasks
	tasks := state.CurrentTasks()
	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks without filter, got %d", len(tasks))
	}

	// With priority filter - only high
	state.FilterPriorities = []Priority{PriorityHigh}
	tasks = state.CurrentTasks()
	if len(tasks) != 1 {
		t.Errorf("Expected 1 high priority task, got %d", len(tasks))
	}
	if tasks[0].ID != "task-1" {
		t.Error("Expected task-1 to match high priority filter")
	}

	// Clear filter
	state.FilterPriorities = nil
	tasks = state.CurrentTasks()
	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks after clearing filter, got %d", len(tasks))
	}
}

func TestHasActiveFilters(t *testing.T) {
	state := NewAppState()

	// No filters
	if state.HasActiveFilters() {
		t.Error("Expected no active filters initially")
	}

	// With priority filter
	state.FilterPriorities = []Priority{PriorityHigh}
	if !state.HasActiveFilters() {
		t.Error("Expected active filters with priority set")
	}

	// With label filter
	state.FilterPriorities = nil
	state.FilterLabels = []string{"bug"}
	if !state.HasActiveFilters() {
		t.Error("Expected active filters with labels set")
	}

	// Clear all
	state.FilterLabels = nil
	if state.HasActiveFilters() {
		t.Error("Expected no active filters after clearing")
	}
}

// ============ Sort Mode Tests ============

func TestSortModeString(t *testing.T) {
	tests := []struct {
		mode     SortMode
		expected string
	}{
		{SortCreatedDesc, "Newest First"},
		{SortCreatedAsc, "Oldest First"},
		{SortDueDateAsc, "Due Date (Earliest)"},
		{SortDueDateDesc, "Due Date (Latest)"},
		{SortTitleAsc, "Title A-Z"},
		{SortTitleDesc, "Title Z-A"},
		{SortPriorityDesc, "Priority (High First)"},
		{SortPriorityAsc, "Priority (Low First)"},
	}

	for _, tt := range tests {
		if got := tt.mode.String(); got != tt.expected {
			t.Errorf("SortMode(%d).String() = %q, want %q", tt.mode, got, tt.expected)
		}
	}
}

func TestAllSortModes(t *testing.T) {
	modes := AllSortModes()
	if len(modes) != 8 {
		t.Errorf("Expected 8 sort modes, got %d", len(modes))
	}
}

func TestSortByCreatedDate(t *testing.T) {
	board := NewBoard("test-board", "Test")

	// Create tasks with different creation times
	task1 := NewTask("older", "Older Task", StatusToDo)
	task2 := NewTask("newer", "Newer Task", StatusToDo)

	// Set creation times (newer is more recent)
	task1.CreatedAt = task2.CreatedAt.Add(-time.Hour)

	board.AddTask(task1)
	board.AddTask(task2)

	state := &AppState{
		Board:       board,
		CurrentPane: StatusToDo,
		SortMode:    SortCreatedDesc,
	}

	// Test SortCreatedDesc (newest first)
	tasks := state.CurrentTasks()
	if len(tasks) != 2 {
		t.Fatalf("Expected 2 tasks, got %d", len(tasks))
	}
	if tasks[0].ID != "newer" {
		t.Errorf("SortCreatedDesc: Expected newer task first, got %s", tasks[0].ID)
	}

	// Test SortCreatedAsc (oldest first)
	state.SortMode = SortCreatedAsc
	tasks = state.CurrentTasks()
	if tasks[0].ID != "older" {
		t.Errorf("SortCreatedAsc: Expected older task first, got %s", tasks[0].ID)
	}
}

func TestSortByDueDate(t *testing.T) {
	board := NewBoard("test-board", "Test")

	// Create tasks with different due dates
	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)
	nextWeek := now.Add(7 * 24 * time.Hour)

	task1 := NewTask("no-due", "No Due Date", StatusToDo)
	// task1 has no due date (nil)

	task2 := NewTask("next-week", "Next Week", StatusToDo)
	task2.DueDate = &nextWeek

	task3 := NewTask("tomorrow", "Tomorrow", StatusToDo)
	task3.DueDate = &tomorrow

	board.AddTask(task1)
	board.AddTask(task2)
	board.AddTask(task3)

	state := &AppState{
		Board:       board,
		CurrentPane: StatusToDo,
		SortMode:    SortDueDateAsc,
	}

	// Test SortDueDateAsc (earliest first, nil last)
	tasks := state.CurrentTasks()
	if len(tasks) != 3 {
		t.Fatalf("Expected 3 tasks, got %d", len(tasks))
	}
	if tasks[0].ID != "tomorrow" {
		t.Errorf("SortDueDateAsc: Expected tomorrow task first, got %s", tasks[0].ID)
	}
	if tasks[1].ID != "next-week" {
		t.Errorf("SortDueDateAsc: Expected next-week task second, got %s", tasks[1].ID)
	}
	if tasks[2].ID != "no-due" {
		t.Errorf("SortDueDateAsc: Expected no-due task last, got %s", tasks[2].ID)
	}

	// Test SortDueDateDesc (latest first, nil first)
	state.SortMode = SortDueDateDesc
	tasks = state.CurrentTasks()
	if tasks[0].ID != "no-due" {
		t.Errorf("SortDueDateDesc: Expected no-due task first, got %s", tasks[0].ID)
	}
	if tasks[1].ID != "next-week" {
		t.Errorf("SortDueDateDesc: Expected next-week task second, got %s", tasks[1].ID)
	}
	if tasks[2].ID != "tomorrow" {
		t.Errorf("SortDueDateDesc: Expected tomorrow task last, got %s", tasks[2].ID)
	}
}

func TestSortByTitle(t *testing.T) {
	board := NewBoard("test-board", "Test")

	task1 := NewTask("charlie", "Charlie Task", StatusToDo)
	task2 := NewTask("alpha", "alpha task", StatusToDo) // lowercase
	task3 := NewTask("beta", "Beta Task", StatusToDo)

	board.AddTask(task1)
	board.AddTask(task2)
	board.AddTask(task3)

	state := &AppState{
		Board:       board,
		CurrentPane: StatusToDo,
		SortMode:    SortTitleAsc,
	}

	// Test SortTitleAsc (A-Z, case-insensitive)
	tasks := state.CurrentTasks()
	if len(tasks) != 3 {
		t.Fatalf("Expected 3 tasks, got %d", len(tasks))
	}
	if tasks[0].ID != "alpha" {
		t.Errorf("SortTitleAsc: Expected alpha task first, got %s", tasks[0].ID)
	}
	if tasks[1].ID != "beta" {
		t.Errorf("SortTitleAsc: Expected beta task second, got %s", tasks[1].ID)
	}
	if tasks[2].ID != "charlie" {
		t.Errorf("SortTitleAsc: Expected charlie task third, got %s", tasks[2].ID)
	}

	// Test SortTitleDesc (Z-A)
	state.SortMode = SortTitleDesc
	tasks = state.CurrentTasks()
	if tasks[0].ID != "charlie" {
		t.Errorf("SortTitleDesc: Expected charlie task first, got %s", tasks[0].ID)
	}
	if tasks[2].ID != "alpha" {
		t.Errorf("SortTitleDesc: Expected alpha task last, got %s", tasks[2].ID)
	}
}

func TestSortByPriority(t *testing.T) {
	board := NewBoard("test-board", "Test")

	task1 := NewTask("low", "Low Priority", StatusToDo)
	task1.Priority = PriorityLow

	task2 := NewTask("high", "High Priority", StatusToDo)
	task2.Priority = PriorityHigh

	task3 := NewTask("medium", "Medium Priority", StatusToDo)
	task3.Priority = PriorityMedium

	task4 := NewTask("none", "No Priority", StatusToDo)
	task4.Priority = PriorityNone

	board.AddTask(task1)
	board.AddTask(task2)
	board.AddTask(task3)
	board.AddTask(task4)

	state := &AppState{
		Board:       board,
		CurrentPane: StatusToDo,
		SortMode:    SortPriorityDesc,
	}

	// Test SortPriorityDesc (high first)
	tasks := state.CurrentTasks()
	if len(tasks) != 4 {
		t.Fatalf("Expected 4 tasks, got %d", len(tasks))
	}
	if tasks[0].ID != "high" {
		t.Errorf("SortPriorityDesc: Expected high task first, got %s", tasks[0].ID)
	}
	if tasks[1].ID != "medium" {
		t.Errorf("SortPriorityDesc: Expected medium task second, got %s", tasks[1].ID)
	}
	if tasks[2].ID != "low" {
		t.Errorf("SortPriorityDesc: Expected low task third, got %s", tasks[2].ID)
	}
	if tasks[3].ID != "none" {
		t.Errorf("SortPriorityDesc: Expected none task last, got %s", tasks[3].ID)
	}

	// Test SortPriorityAsc (low first)
	state.SortMode = SortPriorityAsc
	tasks = state.CurrentTasks()
	if tasks[0].ID != "none" {
		t.Errorf("SortPriorityAsc: Expected none task first, got %s", tasks[0].ID)
	}
	if tasks[3].ID != "high" {
		t.Errorf("SortPriorityAsc: Expected high task last, got %s", tasks[3].ID)
	}
}

func TestSortCombinedWithFilter(t *testing.T) {
	board := NewBoard("test-board", "Test")

	task1 := NewTask("high-1", "High Priority 1", StatusToDo)
	task1.Priority = PriorityHigh

	task2 := NewTask("low-1", "Low Priority 1", StatusToDo)
	task2.Priority = PriorityLow

	task3 := NewTask("high-2", "High Priority 2", StatusToDo)
	task3.Priority = PriorityHigh
	// Make high-2 older than high-1
	task3.CreatedAt = task1.CreatedAt.Add(-time.Hour)

	board.AddTask(task1)
	board.AddTask(task2)
	board.AddTask(task3)

	state := &AppState{
		Board:            board,
		CurrentPane:      StatusToDo,
		SortMode:         SortCreatedDesc, // Newest first
		FilterPriorities: []Priority{PriorityHigh},
	}

	// Should filter to only high priority tasks, then sort by creation date
	tasks := state.CurrentTasks()
	if len(tasks) != 2 {
		t.Fatalf("Expected 2 high priority tasks, got %d", len(tasks))
	}
	// high-1 is newer, should be first
	if tasks[0].ID != "high-1" {
		t.Errorf("Expected high-1 (newer) first after filter+sort, got %s", tasks[0].ID)
	}
	if tasks[1].ID != "high-2" {
		t.Errorf("Expected high-2 (older) second after filter+sort, got %s", tasks[1].ID)
	}
}

func TestSortDoesNotModifyOriginal(t *testing.T) {
	board := NewBoard("test-board", "Test")

	task1 := NewTask("b", "B Task", StatusToDo)
	task2 := NewTask("a", "A Task", StatusToDo)

	board.AddTask(task1)
	board.AddTask(task2)

	state := &AppState{
		Board:       board,
		CurrentPane: StatusToDo,
		SortMode:    SortTitleAsc,
	}

	// Get sorted tasks
	sortedTasks := state.CurrentTasks()

	// Verify sorted order
	if sortedTasks[0].ID != "a" {
		t.Errorf("Expected sorted tasks[0] to be 'a', got %s", sortedTasks[0].ID)
	}

	// Verify original board tasks order is unchanged
	originalTasks := board.TasksByStatus(StatusToDo)
	if originalTasks[0].ID != "b" {
		t.Errorf("Expected original tasks[0] to still be 'b', got %s", originalTasks[0].ID)
	}
}

func TestSetBoardLoadsSortMode(t *testing.T) {
	board := NewBoard("test-board", "Test")
	board.SortMode = SortTitleAsc

	state := NewAppState()
	state.SetBoard(board)

	if state.SortMode != SortTitleAsc {
		t.Errorf("Expected SortMode to be loaded from board, got %v", state.SortMode)
	}
}
