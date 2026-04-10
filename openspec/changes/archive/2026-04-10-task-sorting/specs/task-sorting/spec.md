## ADDED Requirements

### Requirement: User can sort tasks by creation date
The application SHALL allow users to sort tasks by creation date, either newest first or oldest first.

#### Scenario: Sort by newest first
- **WHEN** user selects "Created (Newest)" sort
- **THEN** tasks are ordered with most recently created tasks at the top

#### Scenario: Sort by oldest first
- **WHEN** user selects "Created (Oldest)" sort
- **THEN** tasks are ordered with oldest tasks at the top

### Requirement: User can sort tasks by due date
The application SHALL allow users to sort tasks by due date, with options for soonest or latest first.

#### Scenario: Sort by soonest due date
- **WHEN** user selects "Due Date (Soonest)" sort
- **THEN** tasks with the nearest due date appear at the top
- **AND** tasks without a due date appear at the bottom

#### Scenario: Sort by latest due date
- **WHEN** user selects "Due Date (Latest)" sort
- **THEN** tasks with the furthest due date appear at the top
- **AND** tasks without a due date appear at the top

### Requirement: User can sort tasks alphabetically by title
The application SHALL allow users to sort tasks alphabetically by title.

#### Scenario: Sort A to Z
- **WHEN** user selects "Title (A-Z)" sort
- **THEN** tasks are ordered alphabetically from A to Z

#### Scenario: Sort Z to A
- **WHEN** user selects "Title (Z-A)" sort
- **THEN** tasks are ordered alphabetically from Z to A

### Requirement: User can sort tasks by priority
The application SHALL allow users to sort tasks by priority level.

#### Scenario: Sort high to low priority
- **WHEN** user selects "Priority (High-Low)" sort
- **THEN** tasks are ordered with High priority first, then Medium, Low, and None

#### Scenario: Sort low to high priority
- **WHEN** user selects "Priority (Low-High)" sort
- **THEN** tasks are ordered with None priority first, then Low, Medium, and High

### Requirement: Sort mode is selectable via popup
The application SHALL provide a popup selector for choosing sort mode.

#### Scenario: Open sort selector
- **WHEN** user presses 's' on the main board view
- **THEN** a popup appears listing all available sort options
- **AND** the current sort mode is highlighted

#### Scenario: Select sort mode
- **WHEN** user navigates to a sort option and presses Enter
- **THEN** the sort mode is applied to the task list
- **AND** the popup closes

#### Scenario: Cancel sort selection
- **WHEN** user presses Escape in the sort popup
- **THEN** the popup closes without changing the sort mode

### Requirement: Current sort mode is displayed
The application SHALL display the current sort mode in the UI.

#### Scenario: Sort indicator in header
- **WHEN** a non-default sort is active
- **THEN** the current sort mode is displayed in the header area

#### Scenario: Default sort indicator
- **WHEN** the default sort (Created Newest) is active
- **THEN** no special sort indicator is shown (or shows default label)
