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
