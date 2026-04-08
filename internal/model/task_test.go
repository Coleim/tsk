package model

import (
	"encoding/json"
	"testing"
	"time"
)

// ============ Task Tests (14.1) ============

func TestNewTask(t *testing.T) {
	task := NewTask("test-id", "Test Task", StatusToDo)

	if task.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", task.ID)
	}
	if task.Title != "Test Task" {
		t.Errorf("Expected Title 'Test Task', got '%s'", task.Title)
	}
	if task.Status != StatusToDo {
		t.Errorf("Expected Status StatusToDo, got %v", task.Status)
	}
	if task.Priority != PriorityNone {
		t.Errorf("Expected Priority PriorityNone, got %v", task.Priority)
	}
	if len(task.Labels) != 0 {
		t.Errorf("Expected empty Labels, got %v", task.Labels)
	}
	if task.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}

func TestTaskJSONSerialization(t *testing.T) {
	task := NewTask("test-id", "Test Task", StatusInProgress)
	task.Description = "A test description"
	task.Priority = PriorityHigh
	task.Labels = []string{"urgent", "bug"}
	dueDate := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	task.DueDate = &dueDate

	// Serialize to JSON
	data, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Failed to marshal task: %v", err)
	}

	// Deserialize from JSON
	var restored Task
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("Failed to unmarshal task: %v", err)
	}

	// Verify fields
	if restored.ID != task.ID {
		t.Errorf("ID mismatch: expected %s, got %s", task.ID, restored.ID)
	}
	if restored.Title != task.Title {
		t.Errorf("Title mismatch: expected %s, got %s", task.Title, restored.Title)
	}
	if restored.Description != task.Description {
		t.Errorf("Description mismatch: expected %s, got %s", task.Description, restored.Description)
	}
	if restored.Status != task.Status {
		t.Errorf("Status mismatch: expected %v, got %v", task.Status, restored.Status)
	}
	if restored.Priority != task.Priority {
		t.Errorf("Priority mismatch: expected %v, got %v", task.Priority, restored.Priority)
	}
	if len(restored.Labels) != len(task.Labels) {
		t.Errorf("Labels length mismatch: expected %d, got %d", len(task.Labels), len(restored.Labels))
	}
	if restored.DueDate == nil || !restored.DueDate.Equal(*task.DueDate) {
		t.Errorf("DueDate mismatch")
	}
}

func TestTaskClone(t *testing.T) {
	original := NewTask("test-id", "Test Task", StatusToDo)
	original.Labels = []string{"label1", "label2"}
	dueDate := time.Now()
	original.DueDate = &dueDate

	clone := original.Clone()

	// Verify it's a separate copy
	if clone == original {
		t.Error("Clone should be a different object")
	}

	// Modify clone and verify original is unchanged
	clone.Title = "Modified"
	if original.Title == clone.Title {
		t.Error("Modifying clone title should not affect original")
	}

	clone.Labels[0] = "modified"
	if original.Labels[0] == "modified" {
		t.Error("Modifying clone labels should not affect original")
	}

	*clone.DueDate = time.Now().Add(24 * time.Hour)
	if original.DueDate.Equal(*clone.DueDate) {
		t.Error("Modifying clone due date should not affect original")
	}
}

// ============ Board Tests (14.2) ============

func TestBoardAddTask(t *testing.T) {
	board := NewBoard("board-1", "Test Board")
	task := NewTask("task-1", "Task 1", StatusToDo)

	board.AddTask(task)

	if len(board.Tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(board.Tasks))
	}
	if board.Tasks[0] != task {
		t.Error("Task not added correctly")
	}
}

func TestBoardRemoveTask(t *testing.T) {
	board := NewBoard("board-1", "Test Board")
	task1 := NewTask("task-1", "Task 1", StatusToDo)
	task2 := NewTask("task-2", "Task 2", StatusToDo)

	board.AddTask(task1)
	board.AddTask(task2)

	removed := board.RemoveTask("task-1")

	if removed == nil || removed.ID != "task-1" {
		t.Error("RemoveTask should return the removed task")
	}
	if len(board.Tasks) != 1 {
		t.Errorf("Expected 1 task remaining, got %d", len(board.Tasks))
	}
	if board.Tasks[0].ID != "task-2" {
		t.Error("Wrong task remaining")
	}

	// Try removing non-existent task
	removed = board.RemoveTask("non-existent")
	if removed != nil {
		t.Error("Removing non-existent task should return nil")
	}
}

func TestBoardFindTask(t *testing.T) {
	board := NewBoard("board-1", "Test Board")
	task := NewTask("task-1", "Task 1", StatusToDo)
	board.AddTask(task)

	found := board.FindTask("task-1")
	if found == nil || found.ID != "task-1" {
		t.Error("FindTask should return the task")
	}

	notFound := board.FindTask("non-existent")
	if notFound != nil {
		t.Error("FindTask should return nil for non-existent task")
	}
}

func TestBoardMoveTask(t *testing.T) {
	board := NewBoard("board-1", "Test Board")
	task := NewTask("task-1", "Task 1", StatusToDo)
	board.AddTask(task)

	success := board.MoveTask("task-1", StatusInProgress)
	if !success {
		t.Error("MoveTask should return true")
	}
	if task.Status != StatusInProgress {
		t.Errorf("Task status should be StatusInProgress, got %v", task.Status)
	}

	// Try moving non-existent task
	success = board.MoveTask("non-existent", StatusDone)
	if success {
		t.Error("Moving non-existent task should return false")
	}
}

func TestBoardTasksByStatus(t *testing.T) {
	board := NewBoard("board-1", "Test Board")

	task1 := NewTask("task-1", "Task 1", StatusToDo)
	task2 := NewTask("task-2", "Task 2", StatusInProgress)
	task3 := NewTask("task-3", "Task 3", StatusToDo)

	board.AddTask(task1)
	board.AddTask(task2)
	board.AddTask(task3)

	todoTasks := board.TasksByStatus(StatusToDo)
	if len(todoTasks) != 2 {
		t.Errorf("Expected 2 ToDo tasks, got %d", len(todoTasks))
	}

	inProgressTasks := board.TasksByStatus(StatusInProgress)
	if len(inProgressTasks) != 1 {
		t.Errorf("Expected 1 InProgress task, got %d", len(inProgressTasks))
	}

	doneTasks := board.TasksByStatus(StatusDone)
	if len(doneTasks) != 0 {
		t.Errorf("Expected 0 Done tasks, got %d", len(doneTasks))
	}
}

func TestBoardTaskCount(t *testing.T) {
	board := NewBoard("board-1", "Test Board")

	board.AddTask(NewTask("task-1", "Task 1", StatusToDo))
	board.AddTask(NewTask("task-2", "Task 2", StatusToDo))
	board.AddTask(NewTask("task-3", "Task 3", StatusInProgress))

	if board.TaskCount(StatusToDo) != 2 {
		t.Errorf("Expected 2 ToDo tasks, got %d", board.TaskCount(StatusToDo))
	}
	if board.TaskCount(StatusInProgress) != 1 {
		t.Errorf("Expected 1 InProgress task, got %d", board.TaskCount(StatusInProgress))
	}
	if board.TotalTaskCount() != 3 {
		t.Errorf("Expected 3 total tasks, got %d", board.TotalTaskCount())
	}
}

// ============ Status Tests (14.3) ============

func TestStatusTransitions(t *testing.T) {
	// Test Next transitions
	if StatusToDo.Next() != StatusInProgress {
		t.Error("ToDo.Next() should be InProgress")
	}
	if StatusInProgress.Next() != StatusDone {
		t.Error("InProgress.Next() should be Done")
	}
	if StatusDone.Next() != StatusDone {
		t.Error("Done.Next() should be Done (no change)")
	}

	// Test Prev transitions
	if StatusDone.Prev() != StatusInProgress {
		t.Error("Done.Prev() should be InProgress")
	}
	if StatusInProgress.Prev() != StatusToDo {
		t.Error("InProgress.Prev() should be ToDo")
	}
	if StatusToDo.Prev() != StatusToDo {
		t.Error("ToDo.Prev() should be ToDo (no change)")
	}
}

func TestStatusString(t *testing.T) {
	if StatusToDo.String() != "To Do" {
		t.Errorf("Expected 'To Do', got '%s'", StatusToDo.String())
	}
	if StatusInProgress.String() != "In Progress" {
		t.Errorf("Expected 'In Progress', got '%s'", StatusInProgress.String())
	}
	if StatusDone.String() != "Done" {
		t.Errorf("Expected 'Done', got '%s'", StatusDone.String())
	}
}

func TestAllStatuses(t *testing.T) {
	statuses := AllStatuses()
	if len(statuses) != 3 {
		t.Errorf("Expected 3 statuses, got %d", len(statuses))
	}
	if statuses[0] != StatusToDo || statuses[1] != StatusInProgress || statuses[2] != StatusDone {
		t.Error("AllStatuses should return [ToDo, InProgress, Done]")
	}
}

// ============ Priority Tests (14.4) ============

func TestPriorityOrdering(t *testing.T) {
	if PriorityHigh <= PriorityMedium {
		t.Error("High priority should be greater than Medium")
	}
	if PriorityMedium <= PriorityLow {
		t.Error("Medium priority should be greater than Low")
	}
	if PriorityLow <= PriorityNone {
		t.Error("Low priority should be greater than None")
	}
}

func TestPriorityString(t *testing.T) {
	tests := []struct {
		priority Priority
		expected string
	}{
		{PriorityHigh, "High"},
		{PriorityMedium, "Medium"},
		{PriorityLow, "Low"},
		{PriorityNone, "None"},
	}

	for _, tt := range tests {
		if got := tt.priority.String(); got != tt.expected {
			t.Errorf("Priority(%d).String() = %s, want %s", tt.priority, got, tt.expected)
		}
	}
}

func TestPrioritySymbol(t *testing.T) {
	tests := []struct {
		priority Priority
		expected string
	}{
		{PriorityHigh, "●"},
		{PriorityMedium, "◐"},
		{PriorityLow, "●"},
		{PriorityNone, " "},
	}

	for _, tt := range tests {
		if got := tt.priority.Symbol(); got != tt.expected {
			t.Errorf("Priority(%d).Symbol() = %s, want %s", tt.priority, got, tt.expected)
		}
	}
}

func TestSetPriority(t *testing.T) {
	task := NewTask("test", "Test", StatusToDo)
	originalUpdate := task.UpdatedAt

	time.Sleep(1 * time.Millisecond) // Ensure time difference
	task.SetPriority(PriorityHigh)

	if task.Priority != PriorityHigh {
		t.Error("Priority should be set to High")
	}
	if !task.UpdatedAt.After(originalUpdate) {
		t.Error("UpdatedAt should be updated")
	}
}

// ============ Label Tests (14.5) ============

func TestLabelOperations(t *testing.T) {
	task := NewTask("test", "Test", StatusToDo)

	// Add label
	task.AddLabel("bug")
	if !task.HasLabel("bug") {
		t.Error("Task should have 'bug' label")
	}
	if len(task.Labels) != 1 {
		t.Errorf("Expected 1 label, got %d", len(task.Labels))
	}

	// Add duplicate label (should be no-op)
	task.AddLabel("bug")
	if len(task.Labels) != 1 {
		t.Error("Adding duplicate label should not increase count")
	}

	// Add another label
	task.AddLabel("urgent")
	if len(task.Labels) != 2 {
		t.Errorf("Expected 2 labels, got %d", len(task.Labels))
	}

	// Check HasLabel
	if task.HasLabel("nonexistent") {
		t.Error("HasLabel should return false for non-existent label")
	}

	// Remove label
	task.RemoveLabel("bug")
	if task.HasLabel("bug") {
		t.Error("Task should not have 'bug' label after removal")
	}
	if len(task.Labels) != 1 {
		t.Errorf("Expected 1 label after removal, got %d", len(task.Labels))
	}

	// Remove non-existent label (should be no-op)
	task.RemoveLabel("nonexistent")
	if len(task.Labels) != 1 {
		t.Error("Removing non-existent label should not change count")
	}
}

// ============ Board Search Tests ============

func TestBoardSearch(t *testing.T) {
	board := NewBoard("board-1", "Test Board")

	task1 := NewTask("task-1", "Fix bug in login", StatusToDo)
	task1.Description = "Users cannot log in"
	task1.Labels = []string{"bug", "urgent"}

	task2 := NewTask("task-2", "Add new feature", StatusInProgress)
	task2.Description = "Implement search functionality"
	task2.Labels = []string{"feature"}

	task3 := NewTask("task-3", "Update documentation", StatusDone)

	board.AddTask(task1)
	board.AddTask(task2)
	board.AddTask(task3)

	// Search by title
	results := board.Search("bug")
	if len(results) != 1 || results[0].ID != "task-1" {
		t.Error("Search by title should find task-1")
	}

	// Search by description
	results = board.Search("search functionality")
	if len(results) != 1 || results[0].ID != "task-2" {
		t.Error("Search by description should find task-2")
	}

	// Search by label
	results = board.Search("urgent")
	if len(results) != 1 || results[0].ID != "task-1" {
		t.Error("Search by label should find task-1")
	}

	// Search case-insensitive
	results = board.Search("FIX BUG")
	if len(results) != 1 || results[0].ID != "task-1" {
		t.Error("Search should be case-insensitive")
	}

	// Search with empty query returns all
	results = board.Search("")
	if len(results) != 3 {
		t.Errorf("Empty search should return all tasks, got %d", len(results))
	}

	// Search with no results
	results = board.Search("nonexistent query")
	if len(results) != 0 {
		t.Error("Search with no matches should return empty")
	}
}

func TestBoardSortByPriority(t *testing.T) {
	board := NewBoard("board-1", "Test Board")

	task1 := NewTask("task-1", "Low priority", StatusToDo)
	task1.Priority = PriorityLow

	task2 := NewTask("task-2", "High priority", StatusToDo)
	task2.Priority = PriorityHigh

	task3 := NewTask("task-3", "Medium priority", StatusToDo)
	task3.Priority = PriorityMedium

	board.AddTask(task1)
	board.AddTask(task2)
	board.AddTask(task3)

	board.SortByPriority(StatusToDo)

	tasks := board.TasksByStatus(StatusToDo)
	if tasks[0].Priority != PriorityHigh {
		t.Error("First task should have High priority")
	}
	if tasks[1].Priority != PriorityMedium {
		t.Error("Second task should have Medium priority")
	}
	if tasks[2].Priority != PriorityLow {
		t.Error("Third task should have Low priority")
	}
}

func TestBoardAllLabels(t *testing.T) {
	board := NewBoard("board-1", "Test Board")

	task1 := NewTask("task-1", "Task 1", StatusToDo)
	task1.Labels = []string{"bug", "urgent"}

	task2 := NewTask("task-2", "Task 2", StatusToDo)
	task2.Labels = []string{"feature", "bug"} // "bug" is duplicate

	board.AddTask(task1)
	board.AddTask(task2)

	labels := board.AllLabels()

	// Should have 3 unique labels: bug, feature, urgent (sorted)
	if len(labels) != 3 {
		t.Errorf("Expected 3 unique labels, got %d: %v", len(labels), labels)
	}

	expected := []string{"bug", "feature", "urgent"}
	for i, label := range expected {
		if labels[i] != label {
			t.Errorf("Expected label %s at position %d, got %s", label, i, labels[i])
		}
	}
}
