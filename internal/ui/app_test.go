package ui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/coliva/tsk/internal/model"
	"github.com/coliva/tsk/internal/storage"
	"github.com/coliva/tsk/internal/undo"
)

// Helper to create a test storage
func createTestStorage(t *testing.T) (*storage.Storage, func()) {
	t.Helper()
	dir, err := os.MkdirTemp("", "tsk-ui-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	store, err := storage.NewStorageWithPath(dir)
	if err != nil {
		_ = os.RemoveAll(dir)
		t.Fatalf("Failed to create storage: %v", err)
	}

	return store, func() { _ = os.RemoveAll(dir) }
}

// Helper to create a test app with a board
func createTestApp(t *testing.T) (*App, func()) {
	t.Helper()
	store, cleanup := createTestStorage(t)

	// Create a test board
	board := model.NewBoard("test-board", "Test Board")
	board.AddTask(model.NewTask("task-1", "Task 1", model.StatusToDo))
	board.AddTask(model.NewTask("task-2", "Task 2", model.StatusToDo))
	board.AddTask(model.NewTask("task-3", "Task 3", model.StatusInProgress))
	_ = store.SaveBoard(board)

	app := NewApp(store)
	app.state.SetBoard(board)
	app.state.Width = 120
	app.state.Height = 40

	return app, cleanup
}

// ============ Integration Tests (14.16-14.21) ============

// TestKeyboardNavigation tests j/k moves selection, h/l switches panes (14.17)
func TestKeyboardNavigation(t *testing.T) {
	app, cleanup := createTestApp(t)
	defer cleanup()

	// Initial state
	if app.state.CurrentPane != model.StatusToDo {
		t.Error("Expected initial pane to be ToDo")
	}
	if app.state.SelectedIndex != 0 {
		t.Error("Expected initial selection to be 0")
	}

	// Simulate 'j' key - move down
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if app.state.SelectedIndex != 1 {
		t.Errorf("Expected SelectedIndex 1 after 'j', got %d", app.state.SelectedIndex)
	}

	// Simulate 'k' key - move up
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if app.state.SelectedIndex != 0 {
		t.Errorf("Expected SelectedIndex 0 after 'k', got %d", app.state.SelectedIndex)
	}

	// Simulate 'l' key - move to next pane (InProgress)
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	if app.state.CurrentPane != model.StatusInProgress {
		t.Errorf("Expected pane InProgress after 'l', got %v", app.state.CurrentPane)
	}

	// Simulate 'h' key - move to previous pane (ToDo)
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	if app.state.CurrentPane != model.StatusToDo {
		t.Errorf("Expected pane ToDo after 'h', got %v", app.state.CurrentPane)
	}

	// Test arrow keys
	app.Update(tea.KeyMsg{Type: tea.KeyDown})
	if app.state.SelectedIndex != 1 {
		t.Errorf("Expected SelectedIndex 1 after Down arrow, got %d", app.state.SelectedIndex)
	}

	app.Update(tea.KeyMsg{Type: tea.KeyUp})
	if app.state.SelectedIndex != 0 {
		t.Errorf("Expected SelectedIndex 0 after Up arrow, got %d", app.state.SelectedIndex)
	}

	app.Update(tea.KeyMsg{Type: tea.KeyRight})
	if app.state.CurrentPane != model.StatusInProgress {
		t.Errorf("Expected pane InProgress after Right arrow, got %v", app.state.CurrentPane)
	}

	app.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if app.state.CurrentPane != model.StatusToDo {
		t.Errorf("Expected pane ToDo after Left arrow, got %v", app.state.CurrentPane)
	}
}

// TestUndoRedoCycle tests delete task → undo → task restored (14.18)
func TestUndoRedoCycle(t *testing.T) {
	app, cleanup := createTestApp(t)
	defer cleanup()

	initialTaskCount := len(app.state.Board.Tasks)
	taskID := app.state.CurrentTasks()[0].ID

	// Delete task
	task := app.state.Board.FindTask(taskID)
	if task == nil {
		t.Fatal("Task should exist before delete")
	}

	// Execute delete command via undo manager
	cmd := undo.NewDeleteTaskCommand(task)
	err := app.undoManager.Execute(app.state.Board, cmd)
	if err != nil {
		t.Fatalf("Failed to execute delete: %v", err)
	}

	// Task should be deleted
	if len(app.state.Board.Tasks) != initialTaskCount-1 {
		t.Errorf("Expected %d tasks after delete, got %d", initialTaskCount-1, len(app.state.Board.Tasks))
	}
	if app.state.Board.FindTask(taskID) != nil {
		t.Error("Task should not exist after delete")
	}

	// Undo
	undone, err := app.undoManager.Undo(app.state.Board)
	if err != nil {
		t.Fatalf("Undo failed: %v", err)
	}
	if undone == nil {
		t.Fatal("Expected undone command")
	}

	// Task should be restored
	if len(app.state.Board.Tasks) != initialTaskCount {
		t.Errorf("Expected %d tasks after undo, got %d", initialTaskCount, len(app.state.Board.Tasks))
	}
	if app.state.Board.FindTask(taskID) == nil {
		t.Error("Task should be restored after undo")
	}

	// Redo
	redone, err := app.undoManager.Redo(app.state.Board)
	if err != nil {
		t.Fatalf("Redo failed: %v", err)
	}
	if redone == nil {
		t.Fatal("Expected redone command")
	}

	// Task should be deleted again
	if len(app.state.Board.Tasks) != initialTaskCount-1 {
		t.Errorf("Expected %d tasks after redo, got %d", initialTaskCount-1, len(app.state.Board.Tasks))
	}
}

// TestSearchWorkflow tests search query → results appear (14.19)
func TestSearchWorkflow(t *testing.T) {
	app, cleanup := createTestApp(t)
	defer cleanup()

	// Start search mode with '/' key
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})

	if app.state.Mode != model.ModeSearch {
		t.Errorf("Expected ModeSearch after '/', got %v", app.state.Mode)
	}

	// Type search query via state update
	app.state.UpdateSearch("Task 1")

	// Should find task-1
	results := app.state.SearchResults
	if len(results) == 0 {
		t.Error("Expected search results for 'Task 1'")
	} else {
		found := false
		for _, r := range results {
			if r.ID == "task-1" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected to find task-1 in search results")
		}
	}

	// Exit search with Escape
	app.Update(tea.KeyMsg{Type: tea.KeyEscape})
	if app.state.Mode != model.ModeNormal {
		t.Errorf("Expected ModeNormal after Escape, got %v", app.state.Mode)
	}
}

// TestBoardSwitching tests saves current → loads selected (14.20)
func TestBoardSwitching(t *testing.T) {
	store, cleanup := createTestStorage(t)
	defer cleanup()

	// Create two boards
	board1 := model.NewBoard("board-1", "Board 1")
	board1.AddTask(model.NewTask("task-1", "Task on Board 1", model.StatusToDo))
	_ = store.SaveBoard(board1)

	board2 := model.NewBoard("board-2", "Board 2")
	board2.AddTask(model.NewTask("task-2", "Task on Board 2", model.StatusToDo))
	_ = store.SaveBoard(board2)

	app := NewApp(store)
	app.state.SetBoard(board1)
	app.state.Width = 120
	app.state.Height = 40

	// Verify current board
	if app.state.Board.ID != "board-1" {
		t.Errorf("Expected board-1, got %s", app.state.Board.ID)
	}

	// Make changes to mark dirty
	app.state.MarkDirty()

	// Switch to board 2
	newBoard, err := store.LoadBoard("board-2")
	if err != nil {
		t.Fatalf("Failed to load board-2: %v", err)
	}

	// Before switching, should trigger save if dirty
	if app.state.Dirty {
		_ = store.SaveBoard(board1)
	}

	// Switch board
	app.state.SetBoard(newBoard)

	// Verify switched
	if app.state.Board.ID != "board-2" {
		t.Errorf("Expected board-2 after switch, got %s", app.state.Board.ID)
	}

	// Selection should be reset
	if app.state.SelectedIndex != 0 {
		t.Errorf("Expected SelectedIndex 0 after board switch, got %d", app.state.SelectedIndex)
	}
}

// TestAutoSaveFlag tests dirty flag behavior (14.21)
func TestAutoSaveFlag(t *testing.T) {
	app, cleanup := createTestApp(t)
	defer cleanup()

	// Initially not dirty
	if app.state.Dirty {
		t.Error("Expected Dirty to be false initially")
	}

	// Mark dirty
	app.state.MarkDirty()

	if !app.state.Dirty {
		t.Error("Expected Dirty to be true after MarkDirty")
	}

	// Clear dirty (simulating save)
	app.state.MarkClean()

	if app.state.Dirty {
		t.Error("Expected Dirty to be false after MarkClean")
	}
}

// TestHelpToggle tests help overlay behavior
func TestHelpToggle(t *testing.T) {
	app, cleanup := createTestApp(t)
	defer cleanup()

	// Initially help is hidden
	if app.state.ShowHelp {
		t.Error("Expected ShowHelp to be false initially")
	}

	// Toggle help with '?' key
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})

	if !app.state.ShowHelp {
		t.Error("Expected ShowHelp to be true after '?'")
	}

	// Toggle again
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})

	if app.state.ShowHelp {
		t.Error("Expected ShowHelp to be false after second '?'")
	}
}

// TestMoveTaskBetweenPanes tests > and < keys
func TestMoveTaskBetweenPanes(t *testing.T) {
	app, cleanup := createTestApp(t)
	defer cleanup()

	// Get initial task in ToDo
	tasks := app.state.CurrentTasks()
	if len(tasks) < 1 {
		t.Fatal("Need at least 1 task for test")
	}
	taskID := tasks[0].ID

	// Verify task is in ToDo
	task := app.state.Board.FindTask(taskID)
	if task.Status != model.StatusToDo {
		t.Errorf("Expected task status ToDo, got %v", task.Status)
	}

	// Move task forward
	cmd := undo.NewMoveTaskCommand(taskID, model.StatusToDo, model.StatusInProgress)
	err := app.undoManager.Execute(app.state.Board, cmd)
	if err != nil {
		t.Fatalf("Failed to move task: %v", err)
	}

	// Verify task moved
	task = app.state.Board.FindTask(taskID)
	if task.Status != model.StatusInProgress {
		t.Errorf("Expected task status InProgress after move, got %v", task.Status)
	}

	// Move task backward
	cmd = undo.NewMoveTaskCommand(taskID, model.StatusInProgress, model.StatusToDo)
	err = app.undoManager.Execute(app.state.Board, cmd)
	if err != nil {
		t.Fatalf("Failed to move task back: %v", err)
	}

	// Verify task moved back
	task = app.state.Board.FindTask(taskID)
	if task.Status != model.StatusToDo {
		t.Errorf("Expected task status ToDo after move back, got %v", task.Status)
	}
}

// TestPriorityCommands tests priority setting
func TestPriorityCommands(t *testing.T) {
	app, cleanup := createTestApp(t)
	defer cleanup()

	tasks := app.state.CurrentTasks()
	if len(tasks) < 1 {
		t.Fatal("Need at least 1 task")
	}
	task := tasks[0]

	// Set priority to High
	cmd := undo.NewSetPriorityCommand(task, model.PriorityHigh)
	_ = app.undoManager.Execute(app.state.Board, cmd)

	if task.Priority != model.PriorityHigh {
		t.Errorf("Expected PriorityHigh, got %v", task.Priority)
	}

	// Set priority to Medium
	cmd = undo.NewSetPriorityCommand(task, model.PriorityMedium)
	_ = app.undoManager.Execute(app.state.Board, cmd)

	if task.Priority != model.PriorityMedium {
		t.Errorf("Expected PriorityMedium, got %v", task.Priority)
	}

	// Set priority to Low
	cmd = undo.NewSetPriorityCommand(task, model.PriorityLow)
	_ = app.undoManager.Execute(app.state.Board, cmd)

	if task.Priority != model.PriorityLow {
		t.Errorf("Expected PriorityLow, got %v", task.Priority)
	}

	// Remove priority
	cmd = undo.NewSetPriorityCommand(task, model.PriorityNone)
	_ = app.undoManager.Execute(app.state.Board, cmd)

	if task.Priority != model.PriorityNone {
		t.Errorf("Expected PriorityNone, got %v", task.Priority)
	}
}

// TestCreateTaskWorkflow tests new task creation flow (14.16 partial)
func TestCreateTaskWorkflow(t *testing.T) {
	app, cleanup := createTestApp(t)
	defer cleanup()

	initialCount := len(app.state.Board.Tasks)

	// Create a task via undo command
	cmd := undo.NewCreateTaskCommand("new-task-id", "New Test Task", model.StatusToDo)
	err := app.undoManager.Execute(app.state.Board, cmd)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Verify task was created
	if len(app.state.Board.Tasks) != initialCount+1 {
		t.Errorf("Expected %d tasks after create, got %d", initialCount+1, len(app.state.Board.Tasks))
	}

	// Find the new task
	newTask := app.state.Board.FindTask("new-task-id")
	if newTask == nil {
		t.Fatal("New task not found")
	}
	if newTask.Title != "New Test Task" {
		t.Errorf("Expected title 'New Test Task', got '%s'", newTask.Title)
	}
	if newTask.Status != model.StatusToDo {
		t.Errorf("Expected status ToDo, got %v", newTask.Status)
	}
}

// TestFilterWorkflow tests filter activation and clearing
func TestFilterWorkflow(t *testing.T) {
	store, cleanup := createTestStorage(t)
	defer cleanup()

	// Create board with tasks of different priorities
	board := model.NewBoard("test-board", "Test Board")

	task1 := model.NewTask("task-1", "High Priority Task", model.StatusToDo)
	task1.Priority = model.PriorityHigh

	task2 := model.NewTask("task-2", "Low Priority Task", model.StatusToDo)
	task2.Priority = model.PriorityLow

	task3 := model.NewTask("task-3", "No Priority Task", model.StatusToDo)
	task3.Priority = model.PriorityNone

	board.AddTask(task1)
	board.AddTask(task2)
	board.AddTask(task3)
	_ = store.SaveBoard(board)

	app := NewApp(store)
	app.state.SetBoard(board)

	// All tasks visible
	if len(app.state.CurrentTasks()) != 3 {
		t.Errorf("Expected 3 tasks without filter, got %d", len(app.state.CurrentTasks()))
	}

	// Apply filter
	highPriority := model.PriorityHigh
	app.state.FilterPriority = &highPriority

	if len(app.state.CurrentTasks()) != 1 {
		t.Errorf("Expected 1 task with High priority filter, got %d", len(app.state.CurrentTasks()))
	}

	// Clear filter
	app.state.FilterPriority = nil
	app.state.FilterLabels = nil

	if len(app.state.CurrentTasks()) != 3 {
		t.Errorf("Expected 3 tasks after clearing filter, got %d", len(app.state.CurrentTasks()))
	}
}

// TestAppView tests that View() doesn't panic
func TestAppView(t *testing.T) {
	app, cleanup := createTestApp(t)
	defer cleanup()

	// View should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("View() panicked: %v", r)
		}
	}()

	view := app.View()
	if view == "" {
		t.Error("View() returned empty string")
	}
}

// TestExportPath verifies export path sanitization
func TestExportPath(t *testing.T) {
	store, cleanup := createTestStorage(t)
	defer cleanup()

	board := model.NewBoard("test-board", "My Board")
	path := store.DefaultExportPath(board)

	// Should be sanitized
	if filepath.Base(path) != "tsk-export-My-Board.json" {
		t.Errorf("Expected sanitized filename, got %s", path)
	}
}

// TestDetailModeLabelEditor tests pressing 'L' in detail mode opens label editor
func TestDetailModeLabelEditor(t *testing.T) {
	app, cleanup := createTestApp(t)
	defer cleanup()

	// Open task detail view with Enter
	app.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if app.state.Mode != model.ModeDetail {
		t.Fatalf("Expected ModeDetail after Enter, got %v", app.state.Mode)
	}

	// Get task being viewed
	if app.taskDetail == nil || app.taskDetail.Task == nil {
		t.Fatal("Expected taskDetail to be set")
	}
	taskID := app.taskDetail.Task.ID

	// Press 'L' to open label editor
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'L'}})

	if app.state.Mode != model.ModeLabels {
		t.Errorf("Expected ModeLabels after 'L' in detail mode, got %v", app.state.Mode)
	}

	if app.labelEditor == nil {
		t.Error("Expected labelEditor to be set after 'L'")
	}

	// Label editor should be for the same task
	if app.labelEditor != nil && app.labelEditor.Task.ID != taskID {
		t.Errorf("Expected label editor for task %s, got %s", taskID, app.labelEditor.Task.ID)
	}
}

// TestDeleteConfirmation tests pressing 'd' shows confirmation modal
func TestDeleteConfirmation(t *testing.T) {
	app, cleanup := createTestApp(t)
	defer cleanup()

	initialCount := len(app.state.Board.Tasks)
	taskTitle := app.state.CurrentTasks()[0].Title

	// Press 'd' to delete - should show confirmation
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

	if app.state.Mode != model.ModeModal {
		t.Fatalf("Expected ModeModal after 'd', got %v", app.state.Mode)
	}

	if app.state.ActiveModal == nil {
		t.Fatal("Expected ActiveModal to be set")
	}

	// Task should NOT be deleted yet
	if len(app.state.Board.Tasks) != initialCount {
		t.Error("Task should not be deleted before confirmation")
	}

	// Press 'n' to cancel
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

	if app.state.Mode != model.ModeNormal {
		t.Errorf("Expected ModeNormal after cancel, got %v", app.state.Mode)
	}

	// Task should still exist
	if len(app.state.Board.Tasks) != initialCount {
		t.Error("Task should not be deleted after cancel")
	}

	// Press 'd' again and confirm with 'y'
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

	if app.state.Mode != model.ModeNormal {
		t.Errorf("Expected ModeNormal after confirm, got %v", app.state.Mode)
	}

	// Task should be deleted
	if len(app.state.Board.Tasks) != initialCount-1 {
		t.Errorf("Expected %d tasks after deletion, got %d", initialCount-1, len(app.state.Board.Tasks))
	}

	// Status message should indicate deletion
	if app.state.StatusMessage == "" {
		t.Error("Expected status message after deletion")
	}

	// Verify the deleted task's title was in the status message
	if !strings.Contains(app.state.StatusMessage, taskTitle) {
		t.Errorf("Expected status message to contain '%s', got '%s'", taskTitle, app.state.StatusMessage)
	}
}
