package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/coliva/tsk/internal/model"
)

// Storage handles persistence of boards to disk
type Storage struct {
	basePath   string
	boardsPath string
	backupPath string
}

// NewStorage creates a new storage instance
func NewStorage() (*Storage, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	basePath := filepath.Join(home, ".tsk")
	s := &Storage{
		basePath:   basePath,
		boardsPath: filepath.Join(basePath, "data", "boards"),
		backupPath: filepath.Join(basePath, "backups"),
	}

	if err := s.ensureDirectories(); err != nil {
		return nil, err
	}

	return s, nil
}

// NewStorageWithPath creates a storage instance with a custom base path (for testing)
func NewStorageWithPath(basePath string) (*Storage, error) {
	s := &Storage{
		basePath:   basePath,
		boardsPath: filepath.Join(basePath, "data", "boards"),
		backupPath: filepath.Join(basePath, "backups"),
	}

	if err := s.ensureDirectories(); err != nil {
		return nil, err
	}

	return s, nil
}

// ensureDirectories creates the required directory structure
func (s *Storage) ensureDirectories() error {
	dirs := []string{s.boardsPath, s.backupPath}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

// boardPath returns the file path for a board
func (s *Storage) boardPath(id string) string {
	return filepath.Join(s.boardsPath, fmt.Sprintf("board-%s.json", id))
}

// LoadBoard loads a board from disk
func (s *Storage) LoadBoard(id string) (*model.Board, error) {
	path := s.boardPath(id)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("board not found: %s", id)
		}
		return nil, fmt.Errorf("failed to read board: %w", err)
	}

	var board model.Board
	if err := json.Unmarshal(data, &board); err != nil {
		return nil, fmt.Errorf("failed to parse board: %w", err)
	}

	return &board, nil
}

// SaveBoard saves a board to disk using atomic write
func (s *Storage) SaveBoard(board *model.Board) error {
	board.UpdatedAt = time.Now()

	data, err := json.MarshalIndent(board, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize board: %w", err)
	}

	path := s.boardPath(board.ID)
	tempPath := path + ".tmp"

	// Write to temp file first
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tempPath, path); err != nil {
		_ = os.Remove(tempPath) // Clean up temp file
		return fmt.Errorf("failed to save board: %w", err)
	}

	return nil
}

// DeleteBoard deletes a board from disk (with backup)
func (s *Storage) DeleteBoard(id string) error {
	// Create backup first
	if err := s.BackupBoard(id); err != nil {
		return fmt.Errorf("failed to backup before delete: %w", err)
	}

	path := s.boardPath(id)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete board: %w", err)
	}

	return nil
}

// BackupBoard creates a timestamped backup of a board
func (s *Storage) BackupBoard(id string) error {
	board, err := s.LoadBoard(id)
	if err != nil {
		return err
	}

	timestamp := time.Now().Format("2006-01-02T15-04-05")
	backupName := fmt.Sprintf("board-%s-%s.json", id, timestamp)
	backupPath := filepath.Join(s.backupPath, backupName)

	data, err := json.MarshalIndent(board, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize backup: %w", err)
	}

	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup: %w", err)
	}

	return nil
}

// BoardInfo contains metadata about a board without loading all tasks
type BoardInfo struct {
	ID        string
	Name      string
	TaskCount int
	UpdatedAt time.Time
}

// ListBoards returns info about all available boards
func (s *Storage) ListBoards() ([]BoardInfo, error) {
	entries, err := os.ReadDir(s.boardsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []BoardInfo{}, nil
		}
		return nil, fmt.Errorf("failed to list boards: %w", err)
	}

	var boards []BoardInfo
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasPrefix(entry.Name(), "board-") || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		// Extract ID from filename
		id := strings.TrimPrefix(entry.Name(), "board-")
		id = strings.TrimSuffix(id, ".json")

		board, err := s.LoadBoard(id)
		if err != nil {
			continue // Skip invalid boards
		}

		boards = append(boards, BoardInfo{
			ID:        board.ID,
			Name:      board.Name,
			TaskCount: board.TotalTaskCount(),
			UpdatedAt: board.UpdatedAt,
		})
	}

	// Sort by most recently updated
	sort.Slice(boards, func(i, j int) bool {
		return boards[i].UpdatedAt.After(boards[j].UpdatedAt)
	})

	return boards, nil
}

// MostRecentBoard returns the most recently modified board, or nil if none exist
func (s *Storage) MostRecentBoard() (*model.Board, error) {
	boards, err := s.ListBoards()
	if err != nil {
		return nil, err
	}

	if len(boards) == 0 {
		return nil, nil
	}

	return s.LoadBoard(boards[0].ID)
}

// HasBoards returns true if at least one board exists
func (s *Storage) HasBoards() (bool, error) {
	boards, err := s.ListBoards()
	if err != nil {
		return false, err
	}
	return len(boards) > 0, nil
}

// BoardCount returns the number of boards
func (s *Storage) BoardCount() (int, error) {
	boards, err := s.ListBoards()
	if err != nil {
		return 0, err
	}
	return len(boards), nil
}

// archivePath returns the archive file path for a board
func (s *Storage) archivePath(boardID string) string {
	archiveDir := filepath.Join(s.basePath, "data", "archive")
	_ = os.MkdirAll(archiveDir, 0755) // Ignore error - will fail later if dir unavailable
	return filepath.Join(archiveDir, fmt.Sprintf("%s.json", boardID))
}

// ArchivedTask represents an archived task
type ArchivedTask struct {
	*model.Task
	ArchivedAt time.Time `json:"archivedAt"`
	BoardID    string    `json:"boardId"`
	BoardName  string    `json:"boardName"`
}

// Archive represents the archive file structure
type Archive struct {
	BoardID string         `json:"boardId"`
	Tasks   []ArchivedTask `json:"tasks"`
}

// LoadArchive loads the archive for a board
func (s *Storage) LoadArchive(boardID string) (*Archive, error) {
	path := s.archivePath(boardID)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Archive{BoardID: boardID, Tasks: []ArchivedTask{}}, nil
		}
		return nil, fmt.Errorf("failed to read archive: %w", err)
	}

	var archive Archive
	if err := json.Unmarshal(data, &archive); err != nil {
		return nil, fmt.Errorf("failed to parse archive: %w", err)
	}

	return &archive, nil
}

// SaveArchive saves the archive for a board
func (s *Storage) SaveArchive(archive *Archive) error {
	data, err := json.MarshalIndent(archive, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize archive: %w", err)
	}

	path := s.archivePath(archive.BoardID)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write archive: %w", err)
	}

	return nil
}

// ArchiveTask archives a single task
func (s *Storage) ArchiveTask(task *model.Task, board *model.Board) error {
	archive, err := s.LoadArchive(board.ID)
	if err != nil {
		return err
	}

	archivedTask := ArchivedTask{
		Task:       task,
		ArchivedAt: time.Now(),
		BoardID:    board.ID,
		BoardName:  board.Name,
	}

	archive.Tasks = append(archive.Tasks, archivedTask)

	return s.SaveArchive(archive)
}

// ArchiveTasks archives multiple tasks
func (s *Storage) ArchiveTasks(tasks []*model.Task, board *model.Board) error {
	if len(tasks) == 0 {
		return nil
	}

	archive, err := s.LoadArchive(board.ID)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, task := range tasks {
		archivedTask := ArchivedTask{
			Task:       task,
			ArchivedAt: now,
			BoardID:    board.ID,
			BoardName:  board.Name,
		}
		archive.Tasks = append(archive.Tasks, archivedTask)
	}

	return s.SaveArchive(archive)
}

// UnarchiveTask removes a task from the archive by ID
func (s *Storage) UnarchiveTask(taskID string, boardID string) error {
	archive, err := s.LoadArchive(boardID)
	if err != nil {
		return err
	}

	// Filter out the task
	newTasks := make([]ArchivedTask, 0, len(archive.Tasks))
	for _, t := range archive.Tasks {
		if t.ID != taskID {
			newTasks = append(newTasks, t)
		}
	}
	archive.Tasks = newTasks

	return s.SaveArchive(archive)
}

// UnarchiveTasks removes multiple tasks from the archive by IDs
func (s *Storage) UnarchiveTasks(taskIDs []string, boardID string) error {
	archive, err := s.LoadArchive(boardID)
	if err != nil {
		return err
	}

	// Create a set of IDs to remove
	removeSet := make(map[string]bool)
	for _, id := range taskIDs {
		removeSet[id] = true
	}

	// Filter out the tasks
	newTasks := make([]ArchivedTask, 0, len(archive.Tasks))
	for _, t := range archive.Tasks {
		if !removeSet[t.ID] {
			newTasks = append(newTasks, t)
		}
	}
	archive.Tasks = newTasks

	return s.SaveArchive(archive)
}

// ExportBoard exports a board to a JSON file at the specified path
func (s *Storage) ExportBoard(board *model.Board, exportPath string) error {
	data, err := json.MarshalIndent(board, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize board for export: %w", err)
	}

	if err := os.WriteFile(exportPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}

	return nil
}

// DefaultExportPath returns the default export path for a board
func (s *Storage) DefaultExportPath(board *model.Board) string {
	// Use current working directory with sanitized board name
	safeName := strings.ReplaceAll(board.Name, " ", "-")
	safeName = strings.ReplaceAll(safeName, "/", "-")
	return fmt.Sprintf("tsk-export-%s.json", safeName)
}

// ImportBoard imports a board from a JSON file and creates it with a new ID
func (s *Storage) ImportBoard(importPath string) (*model.Board, error) {
	data, err := os.ReadFile(importPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read import file: %w", err)
	}

	var board model.Board
	if err := json.Unmarshal(data, &board); err != nil {
		return nil, fmt.Errorf("failed to parse import file: %w", err)
	}

	// Generate new ID to avoid conflicts
	board.ID = generateBoardID()
	board.Name = board.Name + " (imported)"
	board.CreatedAt = time.Now()
	board.UpdatedAt = time.Now()

	// Regenerate task IDs to avoid conflicts
	for _, task := range board.Tasks {
		task.ID = generateTaskID()
		task.CreatedAt = time.Now()
		task.UpdatedAt = time.Now()
	}

	// Save the imported board
	if err := s.SaveBoard(&board); err != nil {
		return nil, fmt.Errorf("failed to save imported board: %w", err)
	}

	return &board, nil
}

func generateBoardID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func generateTaskID() string {
	return fmt.Sprintf("task-%d", time.Now().UnixNano())
}
