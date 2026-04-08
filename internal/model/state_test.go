package model

import (
	"testing"
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
	highPriority := PriorityHigh
	state.FilterPriority = &highPriority
	tasks = state.CurrentTasks()
	if len(tasks) != 1 {
		t.Errorf("Expected 1 high priority task, got %d", len(tasks))
	}
	if tasks[0].ID != "task-1" {
		t.Error("Expected task-1 to match high priority filter")
	}

	// Clear filter
	state.FilterPriority = nil
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
	highPriority := PriorityHigh
	state.FilterPriority = &highPriority
	if !state.HasActiveFilters() {
		t.Error("Expected active filters with priority set")
	}

	// With label filter
	state.FilterPriority = nil
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
