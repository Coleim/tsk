## ADDED Requirements

### Requirement: Single-pane view displays one status at a time
The application SHALL display tasks from only one status pane at a time (To Do, In Progress, or Done).

#### Scenario: Single pane layout
- **WHEN** user is on the main screen
- **THEN** only tasks from the current pane are displayed in the task list

#### Scenario: Task card display
- **WHEN** tasks exist in the current pane
- **THEN** each task shows title, priority indicator, and labels (if any)

#### Scenario: Empty pane display
- **WHEN** no tasks exist in current pane
- **THEN** a "No tasks" placeholder is shown

#### Scenario: Task list scrolling
- **WHEN** there are more tasks than can fit in the visible area
- **THEN** the list scrolls to keep the selected task visible

#### Scenario: Scroll follows selection
- **WHEN** user navigates with j/k beyond visible area
- **THEN** the view scrolls to keep the selected task in view

### Requirement: User can switch between panes
The application SHALL allow users to switch between To Do, In Progress, and Done panes.

#### Scenario: Switch to next pane
- **WHEN** user presses 'l'
- **THEN** view switches to the next pane (To Do → In Progress → Done → To Do)

#### Scenario: Switch to previous pane
- **WHEN** user presses 'h'
- **THEN** view switches to the previous pane (Done → In Progress → To Do → Done)

### Requirement: Preview panel shows selected task
The application SHALL display a preview panel showing details of the currently selected task.

#### Scenario: Split layout
- **WHEN** user is on the main board view
- **THEN** the screen is split 30/70 with task list on the left and preview on the right
- **AND** both panels have the same height

#### Scenario: Responsive sizing
- **WHEN** terminal is resized
- **THEN** both panels maintain 30%/70% width proportionally

#### Scenario: Preview updates on navigation
- **WHEN** user navigates to a different task with j/k
- **THEN** preview panel updates to show the newly selected task

#### Scenario: Preview content
- **WHEN** a task is selected
- **THEN** preview shows title, status, priority, due date, colored label badges, and truncated description

#### Scenario: Empty preview
- **WHEN** no task is selected or pane is empty
- **THEN** preview panel shows placeholder or pane statistics

### Requirement: Task detail view for full editing
The application SHALL provide a full-screen detail view for editing a task.
The edit view uses the full terminal width and height for maximum editing space.

#### Scenario: Open task detail
- **WHEN** user presses Enter on a task
- **THEN** a full-screen detail view opens for editing all task properties
- **AND** the view fills the entire terminal, not a centered modal

#### Scenario: Detail view shows colored labels
- **WHEN** task has labels in the detail view
- **THEN** labels appear as colored badges with their assigned colors

#### Scenario: Close detail view
- **WHEN** user presses Escape in detail view
- **THEN** view returns to the single-pane board view

### Requirement: User can search tasks
The application SHALL allow users to search tasks using a compact popup overlay with starts-with matching on titles.

#### Scenario: Open search
- **WHEN** user presses '/'
- **THEN** a centered popup overlay appears with search input
- **AND** the board remains visible behind the popup

#### Scenario: Search matching
- **WHEN** user types search text
- **THEN** tasks whose title starts with the search text are shown first
- **AND** tasks matching in description or labels (contains) are also included

#### Scenario: Search results display
- **WHEN** search has results
- **THEN** maximum 10 results are visible in the popup
- **AND** results can be scrolled with j/k if more exist

#### Scenario: Close search
- **WHEN** user presses Escape in search
- **THEN** search popup closes and full board view is restored

### Requirement: Help overlay displays shortcuts
The application SHALL provide a help overlay showing all keyboard shortcuts.

#### Scenario: Show help
- **WHEN** user presses '?'
- **THEN** a help overlay appears listing all keyboard shortcuts

#### Scenario: Close help
- **WHEN** user presses Escape in help overlay
- **THEN** help overlay closes and returns to previous view

### Requirement: Task list respects current sort order
The application SHALL display tasks in the currently selected sort order.

#### Scenario: Sorted task display
- **WHEN** a sort mode is active
- **THEN** the task list displays tasks in that sorted order
- **AND** sorting applies within each status pane independently

#### Scenario: Sort combined with filter
- **WHEN** both sort and filter are active
- **THEN** filtering is applied first, then sorting is applied to filtered results
