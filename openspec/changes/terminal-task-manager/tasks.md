## 1. Project Setup

- [ ] 1.1 Initialize Go module (go mod init)
- [ ] 1.2 Install dependencies (bubbletea, bubbles, lipgloss)
- [ ] 1.3 Set up project directory structure (cmd/tsk/, internal/model/, internal/ui/, internal/storage/, internal/styles/)
- [ ] 1.4 Create entry point (cmd/tsk/main.go) with basic Bubbletea program
- [ ] 1.5 Add Makefile with build, run, and test targets

## 2. Data Models

- [ ] 2.1 Create Task struct (ID, Title, Description, Status, Priority, Labels, DueDate, Position, CreatedAt, UpdatedAt)
- [ ] 2.2 Create Status type with constants (ToDo, InProgress, Done)
- [ ] 2.3 Create Priority type with iota constants (None=0, Low=1, Medium=2, High=3)
- [ ] 2.4 Create Board struct (ID, Name, Tasks, CreatedAt, UpdatedAt)
- [ ] 2.5 Create root Model struct for Bubbletea with board state, current pane, dirty flag, and auto-save ticker
- [ ] 2.6 Create Mode type (Normal, Insert, Search, Modal) for input handling

## 3. Persistence Layer

- [ ] 3.1 Create storage service with path to ~/.tsk/data/boards/
- [ ] 3.2 Implement ensureDataDirectory function to create storage folders
- [ ] 3.3 Implement loadBoard function to read JSON files
- [ ] 3.4 Implement saveBoard function with atomic writes (temp file + rename)
- [ ] 3.5 Implement dirty flag to track unsaved changes
- [ ] 3.6 Implement action-based save triggers (exit detail, move, create, delete, priority, labels)
- [ ] 3.7 Implement 5-second auto-save timer (only if dirty)
- [ ] 3.8 Implement save on quit and board switch
- [ ] 3.9 Implement listBoards function to enumerate available boards
- [ ] 3.10 Add backup functionality to ~/.tsk/backups/ before destructive operations
- [ ] 3.11 Create welcome screen prompting for board name on first run
- [ ] 3.12 Implement auto-select board on startup (most recently modified)

## 4. Core TUI Framework

- [ ] 4.1 Create root Model with Init, Update, View methods for fullscreen layout
- [ ] 4.2 Create Header model showing app name, board name, and shortcuts hint
- [ ] 4.3 Create two-line StatusBar (line 1: context/feedback, line 2: shortcuts)
- [ ] 4.4 Implement status bar feedback with 3-second auto-clear
- [ ] 4.5 Create Modal model for dialogs with overlay effect
- [ ] 4.6 Implement Lipgloss styles for priority colors and themes
- [ ] 4.7 Add error message display in status bar line 1

## 5. Undo/Redo System

- [ ] 5.1 Create Command interface (Execute, Undo, Description)
- [ ] 5.2 Create UndoManager with undo/redo stacks (max 20)
- [ ] 5.3 Implement CreateTaskCommand with undo (delete created task)
- [ ] 5.4 Implement DeleteTaskCommand with undo (restore deleted task)
- [ ] 5.5 Implement MoveTaskCommand with undo (move back to previous pane)
- [ ] 5.6 Implement EditTaskCommand with undo (restore previous values)
- [ ] 5.7 Implement SetPriorityCommand with undo
- [ ] 5.8 Implement LabelChangeCommand with undo
- [ ] 5.9 Clear redo stack on new action

## 6. Single-Pane Board View

- [ ] 6.1 Create BoardModel with two-panel layout (task list + preview)
- [ ] 6.2 Create PaneTabs showing To Do / In Progress / Done with counts
- [ ] 6.3 Create TaskList showing tasks from current pane
- [ ] 6.4 Create TaskCard view function showing title, priority indicator, and labels
- [ ] 6.5 Implement empty pane placeholder ("No tasks")
- [ ] 6.6 Add visual selection indicator for focused task
- [ ] 6.7 Create PreviewPanel showing selected task details

## 7. Task Detail View

- [ ] 7.1 Create TaskDetailModel with full task information display
- [ ] 7.2 Display task title, description, priority, labels, due date, and status
- [ ] 7.3 Add navigation to return to board view

## 8. Sort Functionality

- [ ] 8.1 Implement sort by priority (s key)
- [ ] 8.2 Sort tasks within current pane (High → Med → Low → None)

## 9. Task Operations

- [ ] 9.1 Create TaskCreateModel with bubbles/textinput for title
- [ ] 9.2 Create TaskEditModel with editable fields
- [ ] 9.3 Implement task creation (n key)
- [ ] 9.4 Implement task editing (e key)
- [ ] 9.5 Implement task deletion with confirmation (d key)
- [ ] 9.6 Implement priority setting (1/2/3/0 keys)
- [ ] 9.7 Implement label management (L key)
- [ ] 9.8 Implement due date setting (t key)
- [ ] 9.9 Implement move task between panes (> / < keys)

## 10. Board Operations

- [ ] 10.1 Create BoardSelectorModel showing board name and task count
- [ ] 10.2 Implement board creation (B key) with text input
- [ ] 10.3 Implement auto-switch to new board after creation
- [ ] 10.4 Implement board switching (b key) with selector overlay
- [ ] 10.5 Implement save current board before switching
- [ ] 10.6 Implement board renaming (R key)
- [ ] 10.7 Implement board deletion with confirmation (D key)

## 11. Keyboard Navigation

- [ ] 11.1 Create keyboard handler with mode-aware key dispatch
- [ ] 11.2 Implement Normal mode handlers (navigation, commands)
- [ ] 11.3 Implement Insert mode handlers (text input, Esc/Enter transitions)
- [ ] 11.4 Implement Search mode handlers (text input, result navigation)
- [ ] 11.5 Implement Modal mode handlers (confirm/cancel)
- [ ] 11.6 Add "-- INSERT --" indicator to status bar when in Insert mode
- [ ] 11.7 Implement vim navigation (j/k for tasks, h/l for panes)
- [ ] 11.8 Implement arrow key navigation (↑/↓ for tasks, ←/→ for panes)
- [ ] 11.9 Implement undo (u key) using UndoManager
- [ ] 11.10 Implement redo (Ctrl+r) using UndoManager
- [ ] 11.11 Implement Escape for cancel/close and mode transitions
- [ ] 11.12 Implement q and :wq for quit (Normal mode only)

## 12. Search and Filter

- [ ] 12.1 Create SearchModel with bubbles/textinput
- [ ] 12.2 Implement search activation (/ key)
- [ ] 12.3 Implement real-time search with 50ms debounce
- [ ] 12.4 Search across title, description, and labels (case-insensitive)
- [ ] 12.5 Display results grouped by pane (To Do / In Progress / Done)
- [ ] 12.6 Implement navigation within search results (j/k, Enter to go to task)
- [ ] 12.7 Implement Esc to clear search and return to board
- [ ] 12.8 Cap results at 100 matches with "100+ matches" indicator
- [ ] 12.9 Create FilterModel with priority and label checkboxes
- [ ] 12.10 Implement filter activation (f key)
- [ ] 12.11 Implement filter by priority (High/Med/Low/None)
- [ ] 12.12 Implement filter by label (multi-select)
- [ ] 12.13 Implement clear filters (F key)
- [ ] 12.14 Add active filter indicator to status bar

## 13. Help and Polish

- [ ] 13.1 Create HelpModel with grouped keyboard shortcuts (Navigation, Task Actions, Board, Priority, Search, Other)
- [ ] 13.2 Implement help toggle (? key opens/closes)
- [ ] 13.3 Implement Escape to close help panel
- [ ] 13.4 Implement archive selected task (a key, DONE pane only)
- [ ] 13.5 Implement archive all Done tasks (A key, DONE pane only)
- [ ] 13.6 Create archive storage at ~/.tsk/data/archive/<board-id>.json
- [ ] 13.7 Show archive keys in status bar only when in DONE pane
- [ ] 13.8 Implement data export (E key)
- [ ] 13.9 Implement data import (I key)
- [ ] 13.10 Add terminal resize handling
- [ ] 13.11 Add graceful error handling throughout

## 14. Testing

### Unit Tests (internal/model/)
- [ ] 14.1 Test Task struct: creation, validation, JSON serialization
- [ ] 14.2 Test Board struct: add/remove/move tasks, find by ID
- [ ] 14.3 Test Status transitions (ToDo → InProgress → Done)
- [ ] 14.4 Test Priority ordering and comparisons
- [ ] 14.5 Test Label operations (add, remove, contains)

### Unit Tests (internal/storage/)
- [ ] 14.6 Test loadBoard with valid JSON
- [ ] 14.7 Test loadBoard with missing/corrupt file (error handling)
- [ ] 14.8 Test saveBoard atomic write (temp file + rename)
- [ ] 14.9 Test listBoards directory enumeration
- [ ] 14.10 Test backup creation before destructive operations
- [ ] 14.11 Test archive file append behavior

### Unit Tests (internal/ui/)
- [ ] 14.12 Test UndoManager: push, undo, redo, max size limit
- [ ] 14.13 Test Command implementations (CreateTask, DeleteTask, MoveTask)
- [ ] 14.14 Test Search function with various queries and filters
- [ ] 14.15 Test mode transitions (Normal → Insert → Normal)

### Integration Tests (Bubbletea tea.Test)
- [ ] 14.16 Test full workflow: create board → add task → move to Done → archive
- [ ] 14.17 Test keyboard navigation: j/k moves selection, h/l switches panes
- [ ] 14.18 Test undo/redo cycle: delete task → undo → task restored
- [ ] 14.19 Test search: type query → results appear → Enter navigates to task
- [ ] 14.20 Test board switching: saves current → loads selected
- [ ] 14.21 Test auto-save: dirty flag set → 5s timer → save called

### Manual Testing Checklist
- [ ] 14.22 Test on macOS Terminal.app
- [ ] 14.23 Test on iTerm2
- [ ] 14.24 Test on Windows Terminal
- [ ] 14.25 Test on Linux xterm/gnome-terminal
- [ ] 14.26 Test terminal resize handling
- [ ] 14.27 Test with 500+ tasks (performance)

## 15. Documentation

- [ ] 15.1 Create README.md with installation instructions
- [ ] 15.2 Add usage section with key bindings table
- [ ] 15.3 Add screenshots/GIFs of main workflows
- [ ] 15.4 Document data storage location (~/.tsk/)
- [ ] 15.5 Add LICENSE file (MIT)
- [ ] 15.6 Add CHANGELOG.md
