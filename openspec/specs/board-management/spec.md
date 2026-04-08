## ADDED Requirements

### Requirement: Application auto-selects board on startup
The application SHALL automatically open a board without requiring user interaction.

#### Scenario: First run (no boards)
- **WHEN** user starts the app and no boards exist
- **THEN** a welcome screen prompts user for a board name
- **AND** after entering a name, the board is created and opened

#### Scenario: Single board exists
- **WHEN** user starts the app and exactly one board exists
- **THEN** that board is opened directly

#### Scenario: Multiple boards exist
- **WHEN** user starts the app and multiple boards exist
- **THEN** the most recently modified board (by updated_at) is opened

### Requirement: User can view boards
The application SHALL display a board with three fixed panes: To Do, In Progress, Done.

#### Scenario: Board displays panes
- **WHEN** user opens a board
- **THEN** pane tabs show "To Do", "In Progress", "Done" with task counts

#### Scenario: Empty pane display
- **WHEN** current pane has no tasks
- **THEN** pane displays with placeholder text "No tasks"

### Requirement: User can create boards
The application SHALL allow users to create new boards with custom names.

#### Scenario: Create new board
- **WHEN** user presses 'B' and enters a board name
- **THEN** a new board is created with empty panes (To Do, In Progress, Done)
- **AND** the app switches to the new board

#### Scenario: Board name validation
- **WHEN** user tries to create a board with an empty name
- **THEN** the system displays an error and prevents creation

### Requirement: User can switch between boards
The application SHALL allow users to switch between multiple boards.

#### Scenario: Board selection
- **WHEN** user presses 'b' to open board selector
- **THEN** a list of available boards appears with task counts

#### Scenario: Board switch saves current board
- **WHEN** user selects a different board from the list
- **THEN** the current board is saved to disk before switching

#### Scenario: Board switch
- **WHEN** user selects a different board from the list
- **THEN** the view updates to show the selected board

### Requirement: User can rename boards
The application SHALL allow users to rename existing boards.

#### Scenario: Rename board
- **WHEN** user presses 'R'
- **THEN** a prompt appears to enter the new board name

### Requirement: User can delete boards
The application SHALL allow users to delete boards with confirmation.

#### Scenario: Delete board with confirmation
- **WHEN** user presses 'D' to delete board
- **THEN** a confirmation dialog appears asking "Delete board '<name>'?"

#### Scenario: Delete board confirmed
- **WHEN** user confirms board deletion
- **THEN** the board and all its tasks are removed and user switches to another board
