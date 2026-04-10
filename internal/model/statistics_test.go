package model

import (
	"testing"
	"time"
)

func TestComputeStatistics_EmptyBoard(t *testing.T) {
	board := NewBoard("test-board", "Test Board")
	stats := ComputeStatistics(board)

	if stats.TotalCount != 0 {
		t.Errorf("Expected TotalCount 0, got %d", stats.TotalCount)
	}
	if stats.CompletionPercent != 0 {
		t.Errorf("Expected CompletionPercent 0, got %f", stats.CompletionPercent)
	}
	if stats.BoardName != "Test Board" {
		t.Errorf("Expected BoardName 'Test Board', got '%s'", stats.BoardName)
	}
}

func TestComputeStatistics_NilBoard(t *testing.T) {
	stats := ComputeStatistics(nil)

	if stats.TotalCount != 0 {
		t.Errorf("Expected TotalCount 0 for nil board, got %d", stats.TotalCount)
	}
}

func TestComputeStatistics_StatusCounts(t *testing.T) {
	board := NewBoard("test-board", "Test Board")

	// Add tasks in different statuses
	board.AddTask(NewTask("1", "Task 1", StatusToDo))
	board.AddTask(NewTask("2", "Task 2", StatusToDo))
	board.AddTask(NewTask("3", "Task 3", StatusInProgress))
	board.AddTask(NewTask("4", "Task 4", StatusDone))
	board.AddTask(NewTask("5", "Task 5", StatusDone))
	board.AddTask(NewTask("6", "Task 6", StatusDone))

	stats := ComputeStatistics(board)

	if stats.TodoCount != 2 {
		t.Errorf("Expected TodoCount 2, got %d", stats.TodoCount)
	}
	if stats.InProgressCount != 1 {
		t.Errorf("Expected InProgressCount 1, got %d", stats.InProgressCount)
	}
	if stats.DoneCount != 3 {
		t.Errorf("Expected DoneCount 3, got %d", stats.DoneCount)
	}
	if stats.TotalCount != 6 {
		t.Errorf("Expected TotalCount 6, got %d", stats.TotalCount)
	}
}

func TestComputeStatistics_PriorityCounts(t *testing.T) {
	board := NewBoard("test-board", "Test Board")

	task1 := NewTask("1", "Task 1", StatusToDo)
	task1.Priority = PriorityHigh
	board.AddTask(task1)

	task2 := NewTask("2", "Task 2", StatusToDo)
	task2.Priority = PriorityHigh
	board.AddTask(task2)

	task3 := NewTask("3", "Task 3", StatusToDo)
	task3.Priority = PriorityMedium
	board.AddTask(task3)

	task4 := NewTask("4", "Task 4", StatusToDo)
	task4.Priority = PriorityLow
	board.AddTask(task4)

	task5 := NewTask("5", "Task 5", StatusToDo)
	task5.Priority = PriorityNone
	board.AddTask(task5)

	stats := ComputeStatistics(board)

	if stats.HighPriorityCount != 2 {
		t.Errorf("Expected HighPriorityCount 2, got %d", stats.HighPriorityCount)
	}
	if stats.MediumPriorityCount != 1 {
		t.Errorf("Expected MediumPriorityCount 1, got %d", stats.MediumPriorityCount)
	}
	if stats.LowPriorityCount != 1 {
		t.Errorf("Expected LowPriorityCount 1, got %d", stats.LowPriorityCount)
	}
	if stats.NoPriorityCount != 1 {
		t.Errorf("Expected NoPriorityCount 1, got %d", stats.NoPriorityCount)
	}
}

func TestComputeStatistics_DueDateMetrics(t *testing.T) {
	board := NewBoard("test-board", "Test Board")

	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	today := now
	tomorrow := now.Add(24 * time.Hour)
	nextWeek := now.Add(5 * 24 * time.Hour)

	// Overdue task
	task1 := NewTask("1", "Overdue Task", StatusToDo)
	task1.DueDate = &yesterday
	board.AddTask(task1)

	// Due today
	task2 := NewTask("2", "Due Today", StatusToDo)
	task2.DueDate = &today
	board.AddTask(task2)

	// Due this week
	task3 := NewTask("3", "Due This Week", StatusToDo)
	task3.DueDate = &tomorrow
	board.AddTask(task3)

	task4 := NewTask("4", "Due This Week 2", StatusToDo)
	task4.DueDate = &nextWeek
	board.AddTask(task4)

	// No due date
	task5 := NewTask("5", "No Due Date", StatusToDo)
	board.AddTask(task5)

	// Done task with past due date - should NOT count as overdue
	task6 := NewTask("6", "Done Overdue", StatusDone)
	task6.DueDate = &yesterday
	board.AddTask(task6)

	stats := ComputeStatistics(board)

	if stats.OverdueCount != 1 {
		t.Errorf("Expected OverdueCount 1, got %d", stats.OverdueCount)
	}
	if stats.DueTodayCount != 1 {
		t.Errorf("Expected DueTodayCount 1, got %d", stats.DueTodayCount)
	}
	if stats.DueThisWeek != 2 {
		t.Errorf("Expected DueThisWeek 2, got %d", stats.DueThisWeek)
	}
	if stats.NoDueDateCount != 1 {
		t.Errorf("Expected NoDueDateCount 1, got %d", stats.NoDueDateCount)
	}
}

func TestComputeStatistics_CompletionPercentage(t *testing.T) {
	board := NewBoard("test-board", "Test Board")

	// 3 done out of 6 = 50%
	board.AddTask(NewTask("1", "Task 1", StatusToDo))
	board.AddTask(NewTask("2", "Task 2", StatusToDo))
	board.AddTask(NewTask("3", "Task 3", StatusInProgress))
	board.AddTask(NewTask("4", "Task 4", StatusDone))
	board.AddTask(NewTask("5", "Task 5", StatusDone))
	board.AddTask(NewTask("6", "Task 6", StatusDone))

	stats := ComputeStatistics(board)

	if stats.CompletionPercent != 50 {
		t.Errorf("Expected CompletionPercent 50, got %f", stats.CompletionPercent)
	}
}

func TestComputeStatistics_CompletionPercentage_AllDone(t *testing.T) {
	board := NewBoard("test-board", "Test Board")

	board.AddTask(NewTask("1", "Task 1", StatusDone))
	board.AddTask(NewTask("2", "Task 2", StatusDone))

	stats := ComputeStatistics(board)

	if stats.CompletionPercent != 100 {
		t.Errorf("Expected CompletionPercent 100, got %f", stats.CompletionPercent)
	}
}

func TestComputeStatistics_CompletionPercentage_NoneDone(t *testing.T) {
	board := NewBoard("test-board", "Test Board")

	board.AddTask(NewTask("1", "Task 1", StatusToDo))
	board.AddTask(NewTask("2", "Task 2", StatusInProgress))

	stats := ComputeStatistics(board)

	if stats.CompletionPercent != 0 {
		t.Errorf("Expected CompletionPercent 0, got %f", stats.CompletionPercent)
	}
}
