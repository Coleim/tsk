package model

import (
	"time"
)

// Status represents the state of a task in the Kanban workflow
type Status string

const (
	StatusToDo       Status = "todo"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

// AllStatuses returns all valid statuses in order
func AllStatuses() []Status {
	return []Status{StatusToDo, StatusInProgress, StatusDone}
}

// String returns the display name for a status
func (s Status) String() string {
	switch s {
	case StatusToDo:
		return "To Do"
	case StatusInProgress:
		return "In Progress"
	case StatusDone:
		return "Done"
	default:
		return string(s)
	}
}

// Next returns the next status in the workflow (for > key)
func (s Status) Next() Status {
	switch s {
	case StatusToDo:
		return StatusInProgress
	case StatusInProgress:
		return StatusDone
	default:
		return s
	}
}

// Prev returns the previous status in the workflow (for < key)
func (s Status) Prev() Status {
	switch s {
	case StatusDone:
		return StatusInProgress
	case StatusInProgress:
		return StatusToDo
	default:
		return s
	}
}

// Priority represents task priority level
type Priority int

const (
	PriorityNone   Priority = 0
	PriorityLow    Priority = 1
	PriorityMedium Priority = 2
	PriorityHigh   Priority = 3
)

// String returns the display name for a priority
func (p Priority) String() string {
	switch p {
	case PriorityHigh:
		return "High"
	case PriorityMedium:
		return "Medium"
	case PriorityLow:
		return "Low"
	default:
		return "None"
	}
}

// Symbol returns a visual indicator for the priority
func (p Priority) Symbol() string {
	switch p {
	case PriorityHigh:
		return "●"
	case PriorityMedium:
		return "◐"
	case PriorityLow:
		return "●"
	default:
		return " "
	}
}

// Task represents a single task item
type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Status      Status     `json:"status"`
	Priority    Priority   `json:"priority"`
	Labels      []string   `json:"labels,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Position    int        `json:"position"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// NewTask creates a new task with the given title and status
func NewTask(id, title string, status Status) *Task {
	now := time.Now()
	return &Task{
		ID:        id,
		Title:     title,
		Status:    status,
		Priority:  PriorityNone,
		Labels:    []string{},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// HasLabel checks if the task has a specific label
func (t *Task) HasLabel(label string) bool {
	for _, l := range t.Labels {
		if l == label {
			return true
		}
	}
	return false
}

// AddLabel adds a label to the task if not already present
func (t *Task) AddLabel(label string) {
	if !t.HasLabel(label) {
		t.Labels = append(t.Labels, label)
		t.UpdatedAt = time.Now()
	}
}

// RemoveLabel removes a label from the task
func (t *Task) RemoveLabel(label string) {
	for i, l := range t.Labels {
		if l == label {
			t.Labels = append(t.Labels[:i], t.Labels[i+1:]...)
			t.UpdatedAt = time.Now()
			return
		}
	}
}

// SetPriority sets the task priority
func (t *Task) SetPriority(p Priority) {
	t.Priority = p
	t.UpdatedAt = time.Now()
}

// SetStatus sets the task status
func (t *Task) SetStatus(s Status) {
	t.Status = s
	t.UpdatedAt = time.Now()
}

// Clone creates a deep copy of the task
func (t *Task) Clone() *Task {
	clone := *t
	clone.Labels = make([]string, len(t.Labels))
	copy(clone.Labels, t.Labels)
	if t.DueDate != nil {
		dueCopy := *t.DueDate
		clone.DueDate = &dueCopy
	}
	return &clone
}
