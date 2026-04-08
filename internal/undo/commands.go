package undo

import (
	"fmt"
	"time"

	"github.com/coliva/tsk/internal/model"
)

// CreateTaskCommand handles task creation with undo
type CreateTaskCommand struct {
	TaskID string
	Title  string
	Status model.Status
}

func NewCreateTaskCommand(taskID, title string, status model.Status) *CreateTaskCommand {
	return &CreateTaskCommand{
		TaskID: taskID,
		Title:  title,
		Status: status,
	}
}

func (c *CreateTaskCommand) Execute(board *model.Board) error {
	task := model.NewTask(c.TaskID, c.Title, c.Status)
	board.AddTask(task)
	return nil
}

func (c *CreateTaskCommand) Undo(board *model.Board) error {
	board.RemoveTask(c.TaskID)
	return nil
}

func (c *CreateTaskCommand) Description() string {
	return fmt.Sprintf("Create task: %s", c.Title)
}

// DeleteTaskCommand handles task deletion with undo
type DeleteTaskCommand struct {
	DeletedTask *model.Task
}

func NewDeleteTaskCommand(task *model.Task) *DeleteTaskCommand {
	// Make a copy of the task for restoration
	taskCopy := *task
	return &DeleteTaskCommand{
		DeletedTask: &taskCopy,
	}
}

func (c *DeleteTaskCommand) Execute(board *model.Board) error {
	board.RemoveTask(c.DeletedTask.ID)
	return nil
}

func (c *DeleteTaskCommand) Undo(board *model.Board) error {
	board.AddTask(c.DeletedTask)
	return nil
}

func (c *DeleteTaskCommand) Description() string {
	return fmt.Sprintf("Delete task: %s", c.DeletedTask.Title)
}

// MoveTaskCommand handles moving tasks between panes with undo
type MoveTaskCommand struct {
	TaskID    string
	OldStatus model.Status
	NewStatus model.Status
}

func NewMoveTaskCommand(taskID string, oldStatus, newStatus model.Status) *MoveTaskCommand {
	return &MoveTaskCommand{
		TaskID:    taskID,
		OldStatus: oldStatus,
		NewStatus: newStatus,
	}
}

func (c *MoveTaskCommand) Execute(board *model.Board) error {
	board.MoveTask(c.TaskID, c.NewStatus)
	return nil
}

func (c *MoveTaskCommand) Undo(board *model.Board) error {
	board.MoveTask(c.TaskID, c.OldStatus)
	return nil
}

func (c *MoveTaskCommand) Description() string {
	return fmt.Sprintf("Move task to %s", c.NewStatus)
}

// EditTaskCommand handles task editing with undo
type EditTaskCommand struct {
	TaskID         string
	OldTitle       string
	OldDescription string
	OldDueDate     *time.Time
	OldLabels      []string
	NewTitle       string
	NewDescription string
	NewDueDate     *time.Time
	NewLabels      []string
}

func NewEditTaskCommand(task *model.Task, newTitle, newDesc string, newDue *time.Time, newLabels []string) *EditTaskCommand {
	var oldDue *time.Time
	if task.DueDate != nil {
		d := *task.DueDate
		oldDue = &d
	}
	var newDueCopy *time.Time
	if newDue != nil {
		d := *newDue
		newDueCopy = &d
	}
	// Copy old labels
	oldLabels := make([]string, len(task.Labels))
	copy(oldLabels, task.Labels)
	// Copy new labels
	newLabelsCopy := make([]string, len(newLabels))
	copy(newLabelsCopy, newLabels)

	return &EditTaskCommand{
		TaskID:         task.ID,
		OldTitle:       task.Title,
		OldDescription: task.Description,
		OldDueDate:     oldDue,
		OldLabels:      oldLabels,
		NewTitle:       newTitle,
		NewDescription: newDesc,
		NewDueDate:     newDueCopy,
		NewLabels:      newLabelsCopy,
	}
}

func (c *EditTaskCommand) Execute(board *model.Board) error {
	task := board.FindTask(c.TaskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", c.TaskID)
	}
	task.Title = c.NewTitle
	task.Description = c.NewDescription
	task.DueDate = c.NewDueDate
	task.Labels = make([]string, len(c.NewLabels))
	copy(task.Labels, c.NewLabels)
	task.UpdatedAt = time.Now()
	return nil
}

func (c *EditTaskCommand) Undo(board *model.Board) error {
	task := board.FindTask(c.TaskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", c.TaskID)
	}
	task.Title = c.OldTitle
	task.Description = c.OldDescription
	task.DueDate = c.OldDueDate
	task.Labels = make([]string, len(c.OldLabels))
	copy(task.Labels, c.OldLabels)
	task.UpdatedAt = time.Now()
	return nil
}

func (c *EditTaskCommand) Description() string {
	return fmt.Sprintf("Edit task: %s", c.OldTitle)
}

// SetPriorityCommand handles priority changes with undo
type SetPriorityCommand struct {
	TaskID      string
	TaskTitle   string
	OldPriority model.Priority
	NewPriority model.Priority
}

func NewSetPriorityCommand(task *model.Task, newPriority model.Priority) *SetPriorityCommand {
	return &SetPriorityCommand{
		TaskID:      task.ID,
		TaskTitle:   task.Title,
		OldPriority: task.Priority,
		NewPriority: newPriority,
	}
}

func (c *SetPriorityCommand) Execute(board *model.Board) error {
	task := board.FindTask(c.TaskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", c.TaskID)
	}
	task.SetPriority(c.NewPriority)
	return nil
}

func (c *SetPriorityCommand) Undo(board *model.Board) error {
	task := board.FindTask(c.TaskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", c.TaskID)
	}
	task.SetPriority(c.OldPriority)
	return nil
}

func (c *SetPriorityCommand) Description() string {
	return fmt.Sprintf("Set priority to %s: %s", c.NewPriority, c.TaskTitle)
}

// LabelChangeCommand handles label changes with undo
type LabelChangeCommand struct {
	TaskID    string
	TaskTitle string
	OldLabels []string
	NewLabels []string
}

func NewLabelChangeCommand(task *model.Task, newLabels []string) *LabelChangeCommand {
	oldLabels := make([]string, len(task.Labels))
	copy(oldLabels, task.Labels)
	newLabelsCopy := make([]string, len(newLabels))
	copy(newLabelsCopy, newLabels)
	return &LabelChangeCommand{
		TaskID:    task.ID,
		TaskTitle: task.Title,
		OldLabels: oldLabels,
		NewLabels: newLabelsCopy,
	}
}

func (c *LabelChangeCommand) Execute(board *model.Board) error {
	task := board.FindTask(c.TaskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", c.TaskID)
	}
	task.Labels = make([]string, len(c.NewLabels))
	copy(task.Labels, c.NewLabels)
	task.UpdatedAt = time.Now()
	return nil
}

func (c *LabelChangeCommand) Undo(board *model.Board) error {
	task := board.FindTask(c.TaskID)
	if task == nil {
		return fmt.Errorf("task not found: %s", c.TaskID)
	}
	task.Labels = make([]string, len(c.OldLabels))
	copy(task.Labels, c.OldLabels)
	task.UpdatedAt = time.Now()
	return nil
}

func (c *LabelChangeCommand) Description() string {
	return fmt.Sprintf("Change labels: %s", c.TaskTitle)
}

// ArchiveTaskCommand handles archiving a single task with undo
type ArchiveTaskCommand struct {
	ArchivedTask *model.Task
	BoardID      string
	OnUndo       func(taskID, boardID string) error // Callback to remove from archive file
}

func NewArchiveTaskCommand(task *model.Task, boardID string, onUndo func(taskID, boardID string) error) *ArchiveTaskCommand {
	// Make a copy of the task for restoration
	taskCopy := *task
	return &ArchiveTaskCommand{
		ArchivedTask: &taskCopy,
		BoardID:      boardID,
		OnUndo:       onUndo,
	}
}

func (c *ArchiveTaskCommand) Execute(board *model.Board) error {
	board.RemoveTask(c.ArchivedTask.ID)
	return nil
}

func (c *ArchiveTaskCommand) Undo(board *model.Board) error {
	// Remove from archive file first
	if c.OnUndo != nil {
		if err := c.OnUndo(c.ArchivedTask.ID, c.BoardID); err != nil {
			return err
		}
	}
	// Restore to board
	board.AddTask(c.ArchivedTask)
	return nil
}

func (c *ArchiveTaskCommand) Description() string {
	return fmt.Sprintf("Archive task: %s", c.ArchivedTask.Title)
}

// ArchiveTasksCommand handles archiving multiple tasks with undo
type ArchiveTasksCommand struct {
	ArchivedTasks []*model.Task
	BoardID       string
	OnUndo        func(taskIDs []string, boardID string) error // Callback to remove from archive file
}

func NewArchiveTasksCommand(tasks []*model.Task, boardID string, onUndo func(taskIDs []string, boardID string) error) *ArchiveTasksCommand {
	// Make copies of all tasks for restoration
	copies := make([]*model.Task, len(tasks))
	for i, task := range tasks {
		taskCopy := *task
		copies[i] = &taskCopy
	}
	return &ArchiveTasksCommand{
		ArchivedTasks: copies,
		BoardID:       boardID,
		OnUndo:        onUndo,
	}
}

func (c *ArchiveTasksCommand) Execute(board *model.Board) error {
	for _, task := range c.ArchivedTasks {
		board.RemoveTask(task.ID)
	}
	return nil
}

func (c *ArchiveTasksCommand) Undo(board *model.Board) error {
	// Remove from archive file first
	if c.OnUndo != nil {
		taskIDs := make([]string, len(c.ArchivedTasks))
		for i, task := range c.ArchivedTasks {
			taskIDs[i] = task.ID
		}
		if err := c.OnUndo(taskIDs, c.BoardID); err != nil {
			return err
		}
	}
	// Restore to board
	for _, task := range c.ArchivedTasks {
		board.AddTask(task)
	}
	return nil
}

func (c *ArchiveTasksCommand) Description() string {
	return fmt.Sprintf("Archive %d tasks", len(c.ArchivedTasks))
}
