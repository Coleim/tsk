## ADDED Requirements

### Requirement: Data persists between sessions
The application SHALL save all data to local storage and restore it on startup.

#### Scenario: Save on action completion
- **WHEN** user completes a meaningful action (exit task detail, move task, create task, delete task, set priority, change labels)
- **THEN** changes are saved to disk

#### Scenario: Auto-save timer
- **WHEN** 5 seconds have passed since last save and board has unsaved changes
- **THEN** changes are automatically saved to disk

#### Scenario: Save on quit
- **WHEN** user quits the application
- **THEN** any pending changes are saved to disk

#### Scenario: Save on board switch
- **WHEN** user switches to a different board
- **THEN** current board changes are saved before loading new board

#### Scenario: Load on startup
- **WHEN** application starts
- **THEN** all boards and tasks are loaded from storage

### Requirement: Skip save when no changes
The application SHALL track a dirty flag and skip disk writes when no changes exist.

#### Scenario: No write when unchanged
- **WHEN** a save is triggered but board has no changes since last save
- **THEN** no disk write occurs

#### Scenario: Dirty flag cleared after save
- **WHEN** board is successfully saved to disk
- **THEN** dirty flag is reset to false

### Requirement: Data stored in user home directory
The application SHALL store data in ~/.tsk/ directory.

#### Scenario: Storage location
- **WHEN** data is saved
- **THEN** files are written to ~/.tsk/data/boards/

#### Scenario: Directory creation
- **WHEN** application runs for the first time
- **THEN** ~/.tsk/data/boards/ directory is created if it doesn't exist

### Requirement: Data format is JSON
The application SHALL store data in human-readable JSON format.

#### Scenario: Board file format
- **WHEN** a board is saved
- **THEN** it is stored as ~/.tsk/data/boards/<board-id>.json

#### Scenario: JSON structure
- **WHEN** reading a board file
- **THEN** it contains id, name, tasks array, created_at, and updated_at fields

#### Scenario: Task structure
- **WHEN** reading a task from JSON
- **THEN** it contains id, title, description, status, priority, labels, due_date, position, created_at, updated_at

### Requirement: Data writes are atomic
The application SHALL use atomic writes to prevent data corruption.

#### Scenario: Safe write operation
- **WHEN** data is saved
- **THEN** it is written to a temp file first, then moved to final location

#### Scenario: Write failure handling
- **WHEN** a write operation fails
- **THEN** the original file is preserved and an error is shown

### Requirement: Automatic backups before destructive operations
The application SHALL create backups before delete operations.

#### Scenario: Pre-delete backup
- **WHEN** user deletes a board
- **THEN** a backup is saved to ~/.tsk/backups/ with timestamp

#### Scenario: Backup retention
- **WHEN** backups are created
- **THEN** the last 10 backups are retained, older ones are deleted

### Requirement: Application creates default board on first run
The application SHALL create a default "My Tasks" board on first run.

#### Scenario: First run initialization
- **WHEN** application runs with no existing data
- **THEN** a "My Tasks" board is created with empty To Do, In Progress, Done panes

### Requirement: Data export capability
The application SHALL allow users to export boards as JSON files.

#### Scenario: Export board
- **WHEN** user presses 'E' for export
- **THEN** current board is exported to a user-specified location

### Requirement: Data import capability
The application SHALL allow users to import boards from JSON files.

#### Scenario: Import board
- **WHEN** user presses 'I' for import and selects a file
- **THEN** the board is imported and added to the board list
