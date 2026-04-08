## 1. Project Setup

- [x] 1.1 Initialize Go module (go mod init)
- [x] 1.2 Install dependencies (bubbletea, bubbles, lipgloss)
- [x] 1.3 Set up project directory structure (cmd/tsk/, internal/model/, internal/ui/, internal/storage/, internal/styles/)
- [x] 1.4 Create entry point (cmd/tsk/main.go) with basic Bubbletea program
- [x] 1.5 Add Makefile with build, run, and test targets

## 2. Data Models

- [x] 2.1 Create Task struct (ID, Title, Description, Status, Priority, Labels, DueDate, Position, CreatedAt, UpdatedAt)
- [x] 2.2 Create Status type with constants (ToDo, InProgress, Done)
- [x] 2.3 Create Priority type with iota constants (None=0, Low=1, Medium=2, High=3)
- [x] 2.4 Create Board struct (ID, Name, Tasks, CreatedAt, UpdatedAt)
- [x] 2.5 Create root Model struct for Bubbletea with board state, current pane, dirty flag, and auto-save ticker
- [x] 2.6 Create Mode type (Normal, Insert, Search, Modal) for input handling

## 3. Persistence Layer

- [x] 3.1 Create storage service with path to ~/.tsk/data/boards/
- [x] 3.2 Implement ensureDataDirectory function to create storage folders
- [x] 3.3 Implement loadBoard function to read JSON files
- [x] 3.4 Implement saveBoard function with atomic writes (temp file + rename)
- [x] 3.5 Implement dirty flag to track unsaved changes
- [x] 3.6 Implement action-based save triggers (exit detail, move, create, delete, priority, labels)
- [x] 3.7 Implement 5-second auto-save timer (only if dirty)
- [x] 3.8 Implement save on quit and board switch
- [x] 3.9 Implement listBoards function to enumerate available boards
- [x] 3.10 Add backup functionality to ~/.tsk/backups/ before destructive operations
- [x] 3.11 Create welcome screen prompting for board name on first run
- [x] 3.12 Implement auto-select board on startup (most recently modified)

## 4. Core TUI Framework

- [x] 4.1 Create root Model with Init, Update, View methods for fullscreen layout
- [x] 4.2 Create Header model showing app name, board name, and shortcuts hint
- [x] 4.3 Create two-line StatusBar (line 1: context/feedback, line 2: shortcuts)
- [x] 4.4 Implement status bar feedback with 3-second auto-clear
- [x] 4.5 Create Modal model for dialogs with overlay effect
- [x] 4.6 Implement Lipgloss styles for priority colors and themes
- [x] 4.7 Add error message display in status bar line 1

## 5. Undo/Redo System

- [x] 5.1 Create Command interface (Execute, Undo, Description)
- [x] 5.2 Create UndoManager with undo/redo stacks (max 20)
- [x] 5.3 Implement CreateTaskCommand with undo (delete created task)
- [x] 5.4 Implement DeleteTaskCommand with undo (restore deleted task)
- [x] 5.5 Implement MoveTaskCommand with undo (move back to previous pane)
- [x] 5.6 Implement EditTaskCommand with undo (restore previous values)
- [x] 5.7 Implement SetPriorityCommand with undo
- [x] 5.8 Implement LabelChangeCommand with undo
- [x] 5.9 Clear redo stack on new action

## 6. Single-Pane Board View

- [x] 6.1 Create BoardModel with two-panel layout (task list + preview)
- [x] 6.2 Create PaneTabs showing To Do / In Progress / Done with counts
- [x] 6.3 Create TaskList showing tasks from current pane
- [x] 6.4 Create TaskCard view function showing title, priority indicator, and labels
- [x] 6.5 Implement empty pane placeholder ("No tasks")
- [x] 6.6 Add visual selection indicator for focused task
- [x] 6.7 Create PreviewPanel showing selected task details

## 7. Task Detail View

- [x] 7.1 Create TaskDetailModel with full task information display
- [x] 7.2 Display task title, description, priority, labels, due date, and status
- [x] 7.3 Add navigation to return to board view

## 8. Sort Functionality

- [x] 8.1 Implement sort by priority (s key)
- [x] 8.2 Sort tasks within current pane (High → Med → Low → None)

## 9. Task Operations

- [x] 9.1 Create TaskCreateModel with bubbles/textinput for title
- [x] 9.2 Create TaskEditModel with editable fields
- [x] 9.3 Implement task creation (n key)
- [x] 9.4 Implement task editing (e key)
- [x] 9.5 Implement task deletion with confirmation (d key)
- [x] 9.6 Implement priority setting (1/2/3/0 keys)
- [x] 9.7 Implement label management (L key)
- [x] 9.8 Implement due date setting (t key)
- [x] 9.9 Implement move task between panes (> / < keys)

## 10. Board Operations

- [x] 10.1 Create BoardSelectorModel showing board name and task count
- [x] 10.2 Implement board creation (B key) with text input
- [x] 10.3 Implement auto-switch to new board after creation
- [x] 10.4 Implement board switching (b key) with selector overlay
- [x] 10.5 Implement save current board before switching
- [x] 10.6 Implement board renaming (R key)
- [x] 10.7 Implement board deletion with confirmation (D key)

## 11. Keyboard Navigation

- [x] 11.1 Create keyboard handler with mode-aware key dispatch
- [x] 11.2 Implement Normal mode handlers (navigation, commands)
- [x] 11.3 Implement Insert mode handlers (text input, Esc/Enter transitions)
- [x] 11.4 Implement Search mode handlers (text input, result navigation)
- [x] 11.5 Implement Modal mode handlers (confirm/cancel)
- [x] 11.6 Add "-- INSERT --" indicator to status bar when in Insert mode
- [x] 11.7 Implement vim navigation (j/k for tasks, h/l for panes)
- [x] 11.8 Implement arrow key navigation (↑/↓ for tasks, ←/→ for panes)
- [x] 11.9 Implement undo (u key) using UndoManager
- [x] 11.10 Implement redo (Ctrl+r) using UndoManager
- [x] 11.11 Implement Escape for cancel/close and mode transitions
- [x] 11.12 Implement q and :wq for quit (Normal mode only)

## 12. Search and Filter

- [x] 12.1 Create SearchModel with bubbles/textinput
- [x] 12.2 Implement search activation (/ key)
- [x] 12.3 Implement real-time search with 50ms debounce
- [x] 12.4 Search across title, description, and labels (case-insensitive)
- [x] 12.5 Display results grouped by pane (To Do / In Progress / Done)
- [x] 12.6 Implement navigation within search results (j/k, Enter to go to task)
- [x] 12.7 Implement Esc to clear search and return to board
- [x] 12.8 Cap results at 100 matches with "100+ matches" indicator
- [x] 12.9 Create FilterModel with priority and label checkboxes
- [x] 12.10 Implement filter activation (f key)
- [x] 12.11 Implement filter by priority (High/Med/Low/None)
- [x] 12.12 Implement filter by label (multi-select)
- [x] 12.13 Implement clear filters (F key)
- [x] 12.14 Add active filter indicator to status bar

## 13. Help and Polish

- [x] 13.1 Create HelpModel with grouped keyboard shortcuts (Navigation, Task Actions, Board, Priority, Search, Other)
- [x] 13.2 Implement help toggle (? key opens/closes)
- [x] 13.3 Implement Escape to close help panel
- [x] 13.4 Implement archive selected task (a key, DONE pane only)
- [x] 13.5 Implement archive all Done tasks (A key, DONE pane only)
- [x] 13.6 Create archive storage at ~/.tsk/data/archive/<board-id>.json
- [x] 13.7 Show archive keys in status bar only when in DONE pane
- [x] 13.8 Implement data export (E key)
- [x] 13.9 Implement data import (I key)
- [x] 13.10 Add terminal resize handling
- [x] 13.11 Add graceful error handling throughout

## 14. Testing

### Unit Tests (internal/model/)
- [x] 14.1 Test Task struct: creation, validation, JSON serialization
- [x] 14.2 Test Board struct: add/remove/move tasks, find by ID
- [x] 14.3 Test Status transitions (ToDo → InProgress → Done)
- [x] 14.4 Test Priority ordering and comparisons
- [x] 14.5 Test Label operations (add, remove, contains)

### Unit Tests (internal/storage/)
- [x] 14.6 Test loadBoard with valid JSON
- [x] 14.7 Test loadBoard with missing/corrupt file (error handling)
- [x] 14.8 Test saveBoard atomic write (temp file + rename)
- [x] 14.9 Test listBoards directory enumeration
- [x] 14.10 Test backup creation before destructive operations
- [x] 14.11 Test archive file append behavior

### Unit Tests (internal/ui/)
- [x] 14.12 Test UndoManager: push, undo, redo, max size limit
- [x] 14.13 Test Command implementations (CreateTask, DeleteTask, MoveTask)
- [x] 14.14 Test Search function with various queries and filters
- [x] 14.15 Test mode transitions (Normal → Insert → Normal)

### Integration Tests (Bubbletea tea.Test)
- [x] 14.16 Test full workflow: create board → add task → move to Done → archive
- [x] 14.17 Test keyboard navigation: j/k moves selection, h/l switches panes
- [x] 14.18 Test undo/redo cycle: delete task → undo → task restored
- [x] 14.19 Test search: type query → results appear → Enter navigates to task
- [x] 14.20 Test board switching: saves current → loads selected
- [x] 14.21 Test auto-save: dirty flag set → 5s timer → save called

## 15. Documentation

- [x] 15.1 Create README.md with installation instructions
- [x] 15.2 Add usage section with key bindings table
- [ ] 15.3 Add screenshots/GIFs of main workflows
- [x] 15.4 Document data storage location (~/.tsk/)
- [x] 15.5 Add LICENSE file (MIT)
- [x] 15.6 Add CHANGELOG.md
