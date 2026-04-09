## MODIFIED Requirements

### Requirement: Application supports color themes
The application SHALL use the current theme's semantic colors for different task priorities and states.

#### Scenario: Priority colors displayed
- **WHEN** tasks have different priorities
- **THEN** high priority shows in theme's error color, medium in warning color, low in muted color, none in text secondary color

#### Scenario: Status colors displayed
- **WHEN** viewing pane tabs
- **THEN** To Do tab uses theme's primary color, In Progress uses warning color, Done uses success color

### Requirement: Application displays two-panel layout
The application SHALL display a two-panel layout with enhanced visual styling.

#### Scenario: Task list panel
- **WHEN** application is running
- **THEN** left panel shows tasks from the current pane with surface background color
- **AND** panel has rounded double-line border with title

#### Scenario: Preview panel
- **WHEN** a task is selected
- **THEN** right panel shows task details with surface background color
- **AND** panel has rounded double-line border with "Preview" title

#### Scenario: Panel proportions
- **WHEN** displaying the layout
- **THEN** task list takes approximately 30% width, preview takes 70%
- **AND** both panels have the same height
- **AND** panels have consistent padding (1 unit vertical, 2 units horizontal)

### Requirement: Application displays pane tabs
The application SHALL display tabs with pill-style visual treatment.

#### Scenario: Pane tabs display
- **WHEN** application is running
- **THEN** tabs show "TO DO (3)  IN PROGRESS (2)  DONE (4)" with counts

#### Scenario: Active tab styling
- **WHEN** a pane is active
- **THEN** its tab has filled background in the status color and bold white text

#### Scenario: Inactive tab styling
- **WHEN** a pane is inactive
- **THEN** its tab has no background fill and muted text color

### Requirement: Task cards have visual polish
The application SHALL render task cards with enhanced visual distinction.

#### Scenario: Selected task card
- **WHEN** a task is selected
- **THEN** it displays with elevated background, bold text, and left border accent in priority color

#### Scenario: Unselected task card
- **WHEN** a task is not selected
- **THEN** it displays with surface background and normal weight text

#### Scenario: Task card spacing
- **WHEN** displaying task cards
- **THEN** each card has 1 unit vertical margin between cards

### Requirement: Application supports full-screen dialogs
The application SHALL render full-screen dialogs with elevated visual treatment.

#### Scenario: Full-screen dialog display
- **WHEN** user triggers a dialog action
- **THEN** dialog appears with elevated background color
- **AND** dialog has thick rounded border in accent color
- **AND** dialog title is bold in accent color

#### Scenario: Dialog dismissal
- **WHEN** user presses Escape in a dialog
- **THEN** the dialog closes and focus returns to the previous view

### Requirement: Application displays two-line status bar
The application SHALL display a polished two-line status bar with visual separation.

#### Scenario: Status bar visual separation
- **WHEN** application is running
- **THEN** status bar has a top border line separating it from content

#### Scenario: Status bar line 1 shows context
- **WHEN** no recent action
- **THEN** line 1 shows task count with muted styling

#### Scenario: Status bar line 2 shows shortcuts
- **WHEN** application is running
- **THEN** line 2 shows shortcuts in secondary text color

### Requirement: Empty states have visual treatment
The application SHALL display polished empty state messages.

#### Scenario: Empty pane display
- **WHEN** a pane has no tasks
- **THEN** centered muted text shows "No tasks" with hint to press 'n' to create

#### Scenario: Empty search results
- **WHEN** search returns no results
- **THEN** centered muted text shows "No matching tasks"
