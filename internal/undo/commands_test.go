package undo

import (
	"testing"
	"time"

	"github.com/coliva/tsk/internal/model"
)

// ============ Manager Tests (14.12) ============

func TestNewManager(t *testing.T) {
	m := NewManager(10)
	if m == nil {
		t.Fatal("Manager should not be nil")
	}
	if m.CanUndo() {
		t.Error("New manager should not be able to undo")
	}
	if m.CanRedo() {
		t.Error("New manager should not be able to redo")
	}
}

func TestManagerDefaultMaxSize(t *testing.T) {
	// Test with 0 should default to 20
	m := NewManager(0)
	board := model.NewBoard("test", "Test Board")

	// Add more than 20 commands
	for i := 0; i < 25; i++ {
		cmd := NewCreateTaskCommand("task-"+string(rune('0'+i)), "Task", model.StatusToDo)
		_ = m.Execute(board, cmd)
	}

	// Should have capped at 20
	count := 0
	for m.CanUndo() {
		_, _ = m.Undo(board)
		count++
	}
	if count > 20 {
		t.Errorf("Expected max 20 undo operations, got %d", count)
	}
}

func TestExecuteAndUndo(t *testing.T) {
	m := NewManager(10)
	board := model.NewBoard("test", "Test Board")

	cmd := NewCreateTaskCommand("task-1", "Test Task", model.StatusToDo)
	if err := m.Execute(board, cmd); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Task should exist
	if board.FindTask("task-1") == nil {
		t.Error("Task should exist after execute")
	}
	if !m.CanUndo() {
		t.Error("Should be able to undo after execute")
	}

	// Undo
	undone, err := m.Undo(board)
	if err != nil {
		t.Fatalf("Undo failed: %v", err)
	}
	if undone == nil {
		t.Error("Undone command should not be nil")
	}

	// Task should be removed
	if board.FindTask("task-1") != nil {
		t.Error("Task should not exist after undo")
	}
	if m.CanUndo() {
		t.Error("Should not be able to undo again")
	}
}

func TestRedo(t *testing.T) {
	m := NewManager(10)
	board := model.NewBoard("test", "Test Board")

	cmd := NewCreateTaskCommand("task-1", "Test Task", model.StatusToDo)
	_ = m.Execute(board, cmd)
	_, _ = m.Undo(board)

	if !m.CanRedo() {
		t.Error("Should be able to redo after undo")
	}

	// Redo
	redone, err := m.Redo(board)
	if err != nil {
		t.Fatalf("Redo failed: %v", err)
	}
	if redone == nil {
		t.Error("Redone command should not be nil")
	}

	// Task should exist again
	if board.FindTask("task-1") == nil {
		t.Error("Task should exist after redo")
	}
	if m.CanRedo() {
		t.Error("Should not be able to redo again")
	}
}

func TestNewActionClearsRedoStack(t *testing.T) {
	m := NewManager(10)
	board := model.NewBoard("test", "Test Board")

	// Create and undo a task
	cmd1 := NewCreateTaskCommand("task-1", "Task 1", model.StatusToDo)
	_ = m.Execute(board, cmd1)
	_, _ = m.Undo(board)

	// Now can redo
	if !m.CanRedo() {
		t.Error("Should be able to redo")
	}

	// Execute new command
	cmd2 := NewCreateTaskCommand("task-2", "Task 2", model.StatusToDo)
	_ = m.Execute(board, cmd2)

	// Redo should be cleared
	if m.CanRedo() {
		t.Error("Redo stack should be cleared after new action")
	}
}

// ============ CreateTaskCommand Tests (14.12) ============

func TestCreateTaskCommand(t *testing.T) {
	board := model.NewBoard("test", "Test Board")
	cmd := NewCreateTaskCommand("task-1", "Test Task", model.StatusInProgress)

	// Execute
	if err := cmd.Execute(board); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	task := board.FindTask("task-1")
	if task == nil {
		t.Fatal("Task should exist")
	}
	if task.Title != "Test Task" {
		t.Errorf("Expected title 'Test Task', got '%s'", task.Title)
	}
	if task.Status != model.StatusInProgress {
		t.Errorf("Expected StatusInProgress, got %v", task.Status)
	}

	// Undo
	if err := cmd.Undo(board); err != nil {
		t.Fatalf("Undo failed: %v", err)
	}

	if board.FindTask("task-1") != nil {
		t.Error("Task should not exist after undo")
	}

	// Description
	desc := cmd.Description()
	if desc != "Create task: Test Task" {
		t.Errorf("Unexpected description: %s", desc)
	}
}

// ============ DeleteTaskCommand Tests (14.13) ============

func TestDeleteTaskCommand(t *testing.T) {
	board := model.NewBoard("test", "Test Board")
	task := model.NewTask("task-1", "Test Task", model.StatusToDo)
	task.Description = "A description"
	task.Priority = model.PriorityHigh
	task.Labels = []string{"bug"}
	board.AddTask(task)

	cmd := NewDeleteTaskCommand(task)

	// Execute
	if err := cmd.Execute(board); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if board.FindTask("task-1") != nil {
		t.Error("Task should not exist after delete")
	}

	// Undo
	if err := cmd.Undo(board); err != nil {
		t.Fatalf("Undo failed: %v", err)
	}

	restored := board.FindTask("task-1")
	if restored == nil {
		t.Fatal("Task should be restored after undo")
	}
	if restored.Description != "A description" {
		t.Error("Task description not restored")
	}
	if restored.Priority != model.PriorityHigh {
		t.Error("Task priority not restored")
	}
	if len(restored.Labels) != 1 || restored.Labels[0] != "bug" {
		t.Error("Task labels not restored")
	}
}

// ============ MoveTaskCommand Tests (14.14) ============

func TestMoveTaskCommand(t *testing.T) {
	board := model.NewBoard("test", "Test Board")
	task := model.NewTask("task-1", "Test Task", model.StatusToDo)
	board.AddTask(task)

	cmd := NewMoveTaskCommand("task-1", model.StatusToDo, model.StatusDone)

	// Execute
	if err := cmd.Execute(board); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if task.Status != model.StatusDone {
		t.Errorf("Expected StatusDone, got %v", task.Status)
	}

	// Undo
	if err := cmd.Undo(board); err != nil {
		t.Fatalf("Undo failed: %v", err)
	}

	if task.Status != model.StatusToDo {
		t.Errorf("Expected StatusToDo after undo, got %v", task.Status)
	}

	// Description
	desc := cmd.Description()
	if desc != "Move task to Done" {
		t.Errorf("Unexpected description: %s", desc)
	}
}

// ============ EditTaskCommand Tests (14.15) ============

func TestEditTaskCommand(t *testing.T) {
	board := model.NewBoard("test", "Test Board")
	task := model.NewTask("task-1", "Original Title", model.StatusToDo)
	task.Description = "Original Description"
	dueDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	task.DueDate = &dueDate
	task.Labels = []string{"original"}
	board.AddTask(task)

	newDue := time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)
	cmd := NewEditTaskCommand(task, "New Title", "New Description", &newDue, []string{"new", "edited"})

	// Execute
	if err := cmd.Execute(board); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if task.Title != "New Title" {
		t.Errorf("Expected 'New Title', got '%s'", task.Title)
	}
	if task.Description != "New Description" {
		t.Errorf("Expected 'New Description', got '%s'", task.Description)
	}
	if len(task.Labels) != 2 {
		t.Errorf("Expected 2 labels, got %d", len(task.Labels))
	}
	if task.DueDate == nil || !task.DueDate.Equal(newDue) {
		t.Error("Due date mismatch")
	}

	// Undo
	if err := cmd.Undo(board); err != nil {
		t.Fatalf("Undo failed: %v", err)
	}

	if task.Title != "Original Title" {
		t.Errorf("Expected 'Original Title' after undo, got '%s'", task.Title)
	}
	if task.Description != "Original Description" {
		t.Errorf("Expected 'Original Description' after undo, got '%s'", task.Description)
	}
	if len(task.Labels) != 1 || task.Labels[0] != "original" {
		t.Error("Labels not restored after undo")
	}
	if task.DueDate == nil || !task.DueDate.Equal(dueDate) {
		t.Error("Due date not restored after undo")
	}
}

func TestEditTaskCommandNotFound(t *testing.T) {
	board := model.NewBoard("test", "Test Board")
	// Task doesn't exist
	cmd := &EditTaskCommand{
		TaskID:         "non-existent",
		OldTitle:       "Old",
		OldDescription: "",
		NewTitle:       "New",
		NewDescription: "",
	}

	err := cmd.Execute(board)
	if err == nil {
		t.Error("Execute should fail for non-existent task")
	}
}

// ============ SetPriorityCommand Tests (14.12) ============

func TestSetPriorityCommand(t *testing.T) {
	board := model.NewBoard("test", "Test Board")
	task := model.NewTask("task-1", "Test Task", model.StatusToDo)
	task.Priority = model.PriorityLow
	board.AddTask(task)

	cmd := NewSetPriorityCommand(task, model.PriorityHigh)

	// Execute
	if err := cmd.Execute(board); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if task.Priority != model.PriorityHigh {
		t.Errorf("Expected PriorityHigh, got %v", task.Priority)
	}

	// Undo
	if err := cmd.Undo(board); err != nil {
		t.Fatalf("Undo failed: %v", err)
	}

	if task.Priority != model.PriorityLow {
		t.Errorf("Expected PriorityLow after undo, got %v", task.Priority)
	}

	// Description
	desc := cmd.Description()
	if desc != "Set priority to High: Test Task" {
		t.Errorf("Unexpected description: %s", desc)
	}
}

// ============ LabelChangeCommand Tests (14.12) ============

func TestLabelChangeCommand(t *testing.T) {
	board := model.NewBoard("test", "Test Board")
	task := model.NewTask("task-1", "Test Task", model.StatusToDo)
	task.Labels = []string{"bug", "urgent"}
	board.AddTask(task)

	cmd := NewLabelChangeCommand(task, []string{"feature", "enhancement"})

	// Execute
	if err := cmd.Execute(board); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if len(task.Labels) != 2 {
		t.Errorf("Expected 2 labels, got %d", len(task.Labels))
	}
	if task.Labels[0] != "feature" || task.Labels[1] != "enhancement" {
		t.Error("Labels not updated correctly")
	}

	// Undo
	if err := cmd.Undo(board); err != nil {
		t.Fatalf("Undo failed: %v", err)
	}

	if len(task.Labels) != 2 {
		t.Errorf("Expected 2 labels after undo, got %d", len(task.Labels))
	}
	if task.Labels[0] != "bug" || task.Labels[1] != "urgent" {
		t.Error("Labels not restored after undo")
	}
}

// ============ Multiple Undo/Redo Operations Test ============

func TestMultipleUndoRedo(t *testing.T) {
	m := NewManager(10)
	board := model.NewBoard("test", "Test Board")

	// Execute multiple commands
	_ = m.Execute(board, NewCreateTaskCommand("task-1", "Task 1", model.StatusToDo))
	_ = m.Execute(board, NewCreateTaskCommand("task-2", "Task 2", model.StatusToDo))
	_ = m.Execute(board, NewCreateTaskCommand("task-3", "Task 3", model.StatusToDo))

	if board.TotalTaskCount() != 3 {
		t.Errorf("Expected 3 tasks, got %d", board.TotalTaskCount())
	}

	// Undo all
	_, _ = m.Undo(board)
	_, _ = m.Undo(board)
	_, _ = m.Undo(board)

	if board.TotalTaskCount() != 0 {
		t.Errorf("Expected 0 tasks after undo all, got %d", board.TotalTaskCount())
	}

	// Redo all
	_, _ = m.Redo(board)
	_, _ = m.Redo(board)
	_, _ = m.Redo(board)

	if board.TotalTaskCount() != 3 {
		t.Errorf("Expected 3 tasks after redo all, got %d", board.TotalTaskCount())
	}
}

func TestUndoOnEmptyStack(t *testing.T) {
	m := NewManager(10)
	board := model.NewBoard("test", "Test Board")

	cmd, err := m.Undo(board)
	if err != nil {
		t.Error("Undo on empty stack should not return error")
	}
	if cmd != nil {
		t.Error("Undo on empty stack should return nil command")
	}
}

func TestRedoOnEmptyStack(t *testing.T) {
	m := NewManager(10)
	board := model.NewBoard("test", "Test Board")

	cmd, err := m.Redo(board)
	if err != nil {
		t.Error("Redo on empty stack should not return error")
	}
	if cmd != nil {
		t.Error("Redo on empty stack should return nil command")
	}
}
