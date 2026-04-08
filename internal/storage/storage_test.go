package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/coliva/tsk/internal/model"
)

// Helper function to create a temp directory for tests
func createTestDir(t *testing.T) (string, func()) {
	t.Helper()
	dir, err := os.MkdirTemp("", "tsk-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return dir, func() { os.RemoveAll(dir) }
}

// ============ Storage Creation Tests (14.6) ============

func TestNewStorageWithPath(t *testing.T) {
	dir, cleanup := createTestDir(t)
	defer cleanup()

	storage, err := NewStorageWithPath(dir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Check directories were created
	boardsDir := filepath.Join(dir, "data", "boards")
	if _, err := os.Stat(boardsDir); os.IsNotExist(err) {
		t.Error("Boards directory should be created")
	}

	backupsDir := filepath.Join(dir, "backups")
	if _, err := os.Stat(backupsDir); os.IsNotExist(err) {
		t.Error("Backups directory should be created")
	}

	_ = storage // Use the storage
}

// ============ Load/Save Tests (14.6-14.7) ============

func TestSaveAndLoadBoard(t *testing.T) {
	dir, cleanup := createTestDir(t)
	defer cleanup()

	storage, err := NewStorageWithPath(dir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Create a board
	board := model.NewBoard("test-board", "Test Board")
	task := model.NewTask("task-1", "Test Task", model.StatusToDo)
	task.Description = "A test description"
	task.Priority = model.PriorityHigh
	task.Labels = []string{"bug", "urgent"}
	board.AddTask(task)

	// Save the board
	if err := storage.SaveBoard(board); err != nil {
		t.Fatalf("Failed to save board: %v", err)
	}

	// Load the board back
	loaded, err := storage.LoadBoard("test-board")
	if err != nil {
		t.Fatalf("Failed to load board: %v", err)
	}

	// Verify loaded data
	if loaded.ID != board.ID {
		t.Errorf("ID mismatch: expected %s, got %s", board.ID, loaded.ID)
	}
	if loaded.Name != board.Name {
		t.Errorf("Name mismatch: expected %s, got %s", board.Name, loaded.Name)
	}
	if len(loaded.Tasks) != 1 {
		t.Fatalf("Task count mismatch: expected 1, got %d", len(loaded.Tasks))
	}

	loadedTask := loaded.Tasks[0]
	if loadedTask.ID != task.ID {
		t.Errorf("Task ID mismatch")
	}
	if loadedTask.Title != task.Title {
		t.Errorf("Task Title mismatch")
	}
	if loadedTask.Priority != task.Priority {
		t.Errorf("Task Priority mismatch")
	}
	if len(loadedTask.Labels) != len(task.Labels) {
		t.Errorf("Task Labels mismatch")
	}
}

func TestLoadNonExistentBoard(t *testing.T) {
	dir, cleanup := createTestDir(t)
	defer cleanup()

	storage, _ := NewStorageWithPath(dir)

	_, err := storage.LoadBoard("non-existent")
	if err == nil {
		t.Error("Loading non-existent board should return error")
	}
}

// ============ ListBoards Tests (14.8) ============

func TestListBoards(t *testing.T) {
	dir, cleanup := createTestDir(t)
	defer cleanup()

	storage, _ := NewStorageWithPath(dir)

	// Empty list initially
	boards, err := storage.ListBoards()
	if err != nil {
		t.Fatalf("ListBoards failed: %v", err)
	}
	if len(boards) != 0 {
		t.Errorf("Expected 0 boards, got %d", len(boards))
	}

	// Create some boards
	board1 := model.NewBoard("board-1", "First Board")
	board1.AddTask(model.NewTask("task-1", "Task 1", model.StatusToDo))
	storage.SaveBoard(board1)

	board2 := model.NewBoard("board-2", "Second Board")
	board2.AddTask(model.NewTask("task-2", "Task 2", model.StatusToDo))
	board2.AddTask(model.NewTask("task-3", "Task 3", model.StatusInProgress))
	storage.SaveBoard(board2)

	boards, err = storage.ListBoards()
	if err != nil {
		t.Fatalf("ListBoards failed: %v", err)
	}
	if len(boards) != 2 {
		t.Errorf("Expected 2 boards, got %d", len(boards))
	}

	// Verify board info
	found := make(map[string]BoardInfo)
	for _, b := range boards {
		found[b.ID] = b
	}

	if info, ok := found["board-1"]; !ok {
		t.Error("board-1 not found")
	} else if info.TaskCount != 1 {
		t.Errorf("board-1 should have 1 task, got %d", info.TaskCount)
	}

	if info, ok := found["board-2"]; !ok {
		t.Error("board-2 not found")
	} else if info.TaskCount != 2 {
		t.Errorf("board-2 should have 2 tasks, got %d", info.TaskCount)
	}
}

func TestMostRecentBoard(t *testing.T) {
	dir, cleanup := createTestDir(t)
	defer cleanup()

	storage, _ := NewStorageWithPath(dir)

	// No boards returns nil
	board, err := storage.MostRecentBoard()
	if err != nil {
		t.Fatalf("MostRecentBoard failed: %v", err)
	}
	if board != nil {
		t.Error("Expected nil when no boards exist")
	}

	// Create boards
	board1 := model.NewBoard("board-1", "First Board")
	storage.SaveBoard(board1)

	board2 := model.NewBoard("board-2", "Second Board")
	storage.SaveBoard(board2) // This is saved later, so it's more recent

	recent, err := storage.MostRecentBoard()
	if err != nil {
		t.Fatalf("MostRecentBoard failed: %v", err)
	}
	if recent.ID != "board-2" {
		t.Errorf("Expected most recent to be board-2, got %s", recent.ID)
	}
}

// ============ Error Handling Tests (14.9) ============

func TestLoadCorruptBoard(t *testing.T) {
	dir, cleanup := createTestDir(t)
	defer cleanup()

	storage, _ := NewStorageWithPath(dir)

	// Write corrupt JSON
	corruptPath := filepath.Join(dir, "data", "boards", "board-corrupt.json")
	os.WriteFile(corruptPath, []byte("{ invalid json"), 0644)

	_, err := storage.LoadBoard("corrupt")
	if err == nil {
		t.Error("Loading corrupt board should return error")
	}
}

func TestAtomicWrite(t *testing.T) {
	dir, cleanup := createTestDir(t)
	defer cleanup()

	storage, _ := NewStorageWithPath(dir)

	board := model.NewBoard("test", "Test Board")
	storage.SaveBoard(board)

	// Verify no temp files remain after save
	entries, _ := os.ReadDir(filepath.Join(dir, "data", "boards"))
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".tmp" {
			t.Error("Temp file should not remain after save")
		}
	}
}

// ============ Backup Tests (14.10) ============

func TestBackupBoard(t *testing.T) {
	dir, cleanup := createTestDir(t)
	defer cleanup()

	storage, _ := NewStorageWithPath(dir)

	board := model.NewBoard("test", "Test Board")
	board.AddTask(model.NewTask("task-1", "Task 1", model.StatusToDo))
	storage.SaveBoard(board)

	// Create backup
	err := storage.BackupBoard("test")
	if err != nil {
		t.Fatalf("BackupBoard failed: %v", err)
	}

	// Check backup file exists
	entries, err := os.ReadDir(filepath.Join(dir, "backups"))
	if err != nil {
		t.Fatalf("Failed to read backups dir: %v", err)
	}

	found := false
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".json" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Backup file should exist")
	}
}

func TestDeleteBoardCreatesBackup(t *testing.T) {
	dir, cleanup := createTestDir(t)
	defer cleanup()

	storage, _ := NewStorageWithPath(dir)

	board := model.NewBoard("test", "Test Board")
	storage.SaveBoard(board)

	// Delete board (should backup first)
	err := storage.DeleteBoard("test")
	if err != nil {
		t.Fatalf("DeleteBoard failed: %v", err)
	}

	// Board should be gone
	_, err = storage.LoadBoard("test")
	if err == nil {
		t.Error("Board should not exist after delete")
	}

	// Backup should exist
	entries, _ := os.ReadDir(filepath.Join(dir, "backups"))
	if len(entries) == 0 {
		t.Error("Backup should be created before delete")
	}
}

// ============ Archive Tests (14.11) ============

func TestArchiveTask(t *testing.T) {
	dir, cleanup := createTestDir(t)
	defer cleanup()

	storage, _ := NewStorageWithPath(dir)

	board := model.NewBoard("test", "Test Board")
	task := model.NewTask("task-1", "Task 1", model.StatusDone)
	board.AddTask(task)
	storage.SaveBoard(board)

	// Archive the task
	err := storage.ArchiveTask(task, board)
	if err != nil {
		t.Fatalf("ArchiveTask failed: %v", err)
	}

	// Load archive and verify
	archive, err := storage.LoadArchive(board.ID)
	if err != nil {
		t.Fatalf("LoadArchive failed: %v", err)
	}

	if len(archive.Tasks) != 1 {
		t.Errorf("Expected 1 archived task, got %d", len(archive.Tasks))
	}

	if archive.Tasks[0].Task.ID != task.ID {
		t.Error("Archived task ID mismatch")
	}
	if archive.Tasks[0].BoardID != board.ID {
		t.Error("BoardID mismatch in archived task")
	}
	if archive.Tasks[0].ArchivedAt.IsZero() {
		t.Error("ArchivedAt should be set")
	}
}

func TestArchiveMultipleTasks(t *testing.T) {
	dir, cleanup := createTestDir(t)
	defer cleanup()

	storage, _ := NewStorageWithPath(dir)

	board := model.NewBoard("test", "Test Board")
	task1 := model.NewTask("task-1", "Task 1", model.StatusDone)
	task2 := model.NewTask("task-2", "Task 2", model.StatusDone)
	board.AddTask(task1)
	board.AddTask(task2)
	storage.SaveBoard(board)

	// Archive both tasks
	err := storage.ArchiveTasks([]*model.Task{task1, task2}, board)
	if err != nil {
		t.Fatalf("ArchiveTasks failed: %v", err)
	}

	archive, _ := storage.LoadArchive(board.ID)
	if len(archive.Tasks) != 2 {
		t.Errorf("Expected 2 archived tasks, got %d", len(archive.Tasks))
	}
}

func TestLoadEmptyArchive(t *testing.T) {
	dir, cleanup := createTestDir(t)
	defer cleanup()

	storage, _ := NewStorageWithPath(dir)

	archive, err := storage.LoadArchive("non-existent")
	if err != nil {
		t.Fatalf("LoadArchive should not fail for non-existent: %v", err)
	}
	if len(archive.Tasks) != 0 {
		t.Error("Empty archive should have no tasks")
	}
}

// ============ Export/Import Tests ============

func TestExportBoard(t *testing.T) {
	dir, cleanup := createTestDir(t)
	defer cleanup()

	storage, _ := NewStorageWithPath(dir)

	board := model.NewBoard("test", "Test Board")
	board.AddTask(model.NewTask("task-1", "Task 1", model.StatusToDo))
	board.AddTask(model.NewTask("task-2", "Task 2", model.StatusInProgress))

	exportPath := filepath.Join(dir, "export.json")
	err := storage.ExportBoard(board, exportPath)
	if err != nil {
		t.Fatalf("ExportBoard failed: %v", err)
	}

	// Verify export file
	data, err := os.ReadFile(exportPath)
	if err != nil {
		t.Fatalf("Failed to read export file: %v", err)
	}

	var exported model.Board
	if err := json.Unmarshal(data, &exported); err != nil {
		t.Fatalf("Failed to parse export file: %v", err)
	}

	if exported.ID != board.ID {
		t.Error("Exported board ID mismatch")
	}
	if len(exported.Tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(exported.Tasks))
	}
}

func TestImportBoard(t *testing.T) {
	dir, cleanup := createTestDir(t)
	defer cleanup()

	storage, _ := NewStorageWithPath(dir)

	// Create an export file manually
	exportBoard := model.NewBoard("original-id", "Imported Board")
	exportBoard.AddTask(model.NewTask("task-orig", "Task 1", model.StatusToDo))

	data, _ := json.MarshalIndent(exportBoard, "", "  ")
	importPath := filepath.Join(dir, "import.json")
	os.WriteFile(importPath, data, 0644)

	// Import it
	imported, err := storage.ImportBoard(importPath)
	if err != nil {
		t.Fatalf("ImportBoard failed: %v", err)
	}

	// Verify new IDs were assigned
	if imported.ID == "original-id" {
		t.Error("Imported board should have new ID")
	}
	if imported.Name != "Imported Board (imported)" {
		t.Errorf("Expected name with (imported) suffix, got %s", imported.Name)
	}
	if len(imported.Tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(imported.Tasks))
	}
	if imported.Tasks[0].ID == "task-orig" {
		t.Error("Imported task should have new ID")
	}

	// Verify board was saved
	loaded, err := storage.LoadBoard(imported.ID)
	if err != nil {
		t.Fatalf("Imported board not saved: %v", err)
	}
	if loaded.Name != imported.Name {
		t.Error("Saved board name mismatch")
	}
}

func TestImportBoardInvalidFile(t *testing.T) {
	dir, cleanup := createTestDir(t)
	defer cleanup()

	storage, _ := NewStorageWithPath(dir)

	// Non-existent file
	_, err := storage.ImportBoard("/non/existent/path.json")
	if err == nil {
		t.Error("Importing non-existent file should fail")
	}

	// Invalid JSON
	invalidPath := filepath.Join(dir, "invalid.json")
	os.WriteFile(invalidPath, []byte("not json"), 0644)

	_, err = storage.ImportBoard(invalidPath)
	if err == nil {
		t.Error("Importing invalid JSON should fail")
	}
}

func TestDefaultExportPath(t *testing.T) {
	dir, cleanup := createTestDir(t)
	defer cleanup()

	storage, _ := NewStorageWithPath(dir)

	board := model.NewBoard("test", "My Test Board")
	path := storage.DefaultExportPath(board)

	if path != "tsk-export-My-Test-Board.json" {
		t.Errorf("Expected sanitized filename, got %s", path)
	}
}

func TestHasBoards(t *testing.T) {
	dir, cleanup := createTestDir(t)
	defer cleanup()

	storage, _ := NewStorageWithPath(dir)

	has, _ := storage.HasBoards()
	if has {
		t.Error("Should not have boards initially")
	}

	board := model.NewBoard("test", "Test")
	storage.SaveBoard(board)

	has, _ = storage.HasBoards()
	if !has {
		t.Error("Should have boards after saving")
	}
}

func TestBoardCount(t *testing.T) {
	dir, cleanup := createTestDir(t)
	defer cleanup()

	storage, _ := NewStorageWithPath(dir)

	count, _ := storage.BoardCount()
	if count != 0 {
		t.Errorf("Expected 0, got %d", count)
	}

	storage.SaveBoard(model.NewBoard("b1", "Board 1"))
	storage.SaveBoard(model.NewBoard("b2", "Board 2"))

	count, _ = storage.BoardCount()
	if count != 2 {
		t.Errorf("Expected 2, got %d", count)
	}
}
