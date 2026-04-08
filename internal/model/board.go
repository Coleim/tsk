package model

import (
	"sort"
	"strings"
	"time"
)

// Board represents a Kanban board with tasks
type Board struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Tasks     []*Task   `json:"tasks"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewBoard creates a new board with the given name
func NewBoard(id, name string) *Board {
	now := time.Now()
	return &Board{
		ID:        id,
		Name:      name,
		Tasks:     []*Task{},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// AddTask adds a task to the board
func (b *Board) AddTask(task *Task) {
	// Set position to end of its status group
	task.Position = len(b.TasksByStatus(task.Status))
	b.Tasks = append(b.Tasks, task)
	b.UpdatedAt = time.Now()
}

// RemoveTask removes a task by ID and returns it
func (b *Board) RemoveTask(id string) *Task {
	for i, task := range b.Tasks {
		if task.ID == id {
			removed := task
			b.Tasks = append(b.Tasks[:i], b.Tasks[i+1:]...)
			b.UpdatedAt = time.Now()
			return removed
		}
	}
	return nil
}

// FindTask finds a task by ID
func (b *Board) FindTask(id string) *Task {
	for _, task := range b.Tasks {
		if task.ID == id {
			return task
		}
	}
	return nil
}

// TasksByStatus returns all tasks with the given status, sorted by position
func (b *Board) TasksByStatus(status Status) []*Task {
	var result []*Task
	for _, task := range b.Tasks {
		if task.Status == status {
			result = append(result, task)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Position < result[j].Position
	})
	return result
}

// TaskCount returns the number of tasks with the given status
func (b *Board) TaskCount(status Status) int {
	count := 0
	for _, task := range b.Tasks {
		if task.Status == status {
			count++
		}
	}
	return count
}

// TotalTaskCount returns the total number of tasks
func (b *Board) TotalTaskCount() int {
	return len(b.Tasks)
}

// MoveTask moves a task to a new status
func (b *Board) MoveTask(id string, newStatus Status) bool {
	task := b.FindTask(id)
	if task == nil {
		return false
	}
	task.SetStatus(newStatus)
	task.Position = len(b.TasksByStatus(newStatus)) - 1
	b.UpdatedAt = time.Now()
	return true
}

// ReorderTask moves a task to a new position within its status
func (b *Board) ReorderTask(id string, newPosition int) bool {
	task := b.FindTask(id)
	if task == nil {
		return false
	}
	tasks := b.TasksByStatus(task.Status)
	if newPosition < 0 || newPosition >= len(tasks) {
		return false
	}
	// Reorder positions
	oldPos := task.Position
	if newPosition > oldPos {
		for _, t := range tasks {
			if t.Position > oldPos && t.Position <= newPosition {
				t.Position--
			}
		}
	} else {
		for _, t := range tasks {
			if t.Position >= newPosition && t.Position < oldPos {
				t.Position++
			}
		}
	}
	task.Position = newPosition
	b.UpdatedAt = time.Now()
	return true
}

// SortByPriority sorts tasks within a status by priority (high first)
func (b *Board) SortByPriority(status Status) {
	tasks := b.TasksByStatus(status)
	sort.SliceStable(tasks, func(i, j int) bool {
		return tasks[i].Priority > tasks[j].Priority
	})
	for i, task := range tasks {
		task.Position = i
	}
	b.UpdatedAt = time.Now()
}

// Search finds tasks matching the query across title, description, and labels
func (b *Board) Search(query string) []*Task {
	if query == "" {
		return b.Tasks
	}
	query = strings.ToLower(query)
	var results []*Task
	for _, task := range b.Tasks {
		if b.matchesQuery(task, query) {
			results = append(results, task)
			if len(results) >= 100 {
				break
			}
		}
	}
	return results
}

func (b *Board) matchesQuery(task *Task, query string) bool {
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

// AllLabels returns all unique labels used in the board
func (b *Board) AllLabels() []string {
	labelSet := make(map[string]bool)
	for _, task := range b.Tasks {
		for _, label := range task.Labels {
			labelSet[label] = true
		}
	}
	var labels []string
	for label := range labelSet {
		labels = append(labels, label)
	}
	sort.Strings(labels)
	return labels
}
