package model

import (
	"time"
)

// BoardStatistics holds computed statistics for a board
type BoardStatistics struct {
	// Status counts
	TodoCount       int
	InProgressCount int
	DoneCount       int
	TotalCount      int

	// Priority counts
	HighPriorityCount   int
	MediumPriorityCount int
	LowPriorityCount    int
	NoPriorityCount     int

	// Due date metrics
	OverdueCount   int
	DueTodayCount  int
	DueThisWeek    int
	NoDueDateCount int

	// Completion
	CompletionPercent float64

	// Board name for display
	BoardName string
}

// ComputeStatistics calculates all statistics for a board
func ComputeStatistics(board *Board) BoardStatistics {
	if board == nil {
		return BoardStatistics{}
	}

	stats := BoardStatistics{
		BoardName: board.Name,
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfToday := today.Add(24 * time.Hour)
	endOfWeek := today.Add(7 * 24 * time.Hour)

	for _, task := range board.Tasks {
		stats.TotalCount++

		// Status counts
		switch task.Status {
		case StatusToDo:
			stats.TodoCount++
		case StatusInProgress:
			stats.InProgressCount++
		case StatusDone:
			stats.DoneCount++
		}

		// Priority counts
		switch task.Priority {
		case PriorityHigh:
			stats.HighPriorityCount++
		case PriorityMedium:
			stats.MediumPriorityCount++
		case PriorityLow:
			stats.LowPriorityCount++
		case PriorityNone:
			stats.NoPriorityCount++
		}

		// Due date metrics (only for non-done tasks)
		if task.DueDate != nil && task.Status != StatusDone {
			dueDate := *task.DueDate
			if dueDate.Before(today) {
				stats.OverdueCount++
			} else if dueDate.Before(endOfToday) {
				stats.DueTodayCount++
			} else if dueDate.Before(endOfWeek) {
				stats.DueThisWeek++
			}
		} else if task.DueDate == nil {
			stats.NoDueDateCount++
		}
	}

	// Completion percentage
	if stats.TotalCount > 0 {
		stats.CompletionPercent = float64(stats.DoneCount) / float64(stats.TotalCount) * 100
	}

	return stats
}
