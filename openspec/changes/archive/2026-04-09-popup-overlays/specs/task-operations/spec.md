## MODIFIED Requirements

### Requirement: User can create tasks
The application SHALL allow users to create new tasks with a title using a compact popup overlay.

#### Scenario: Create task with title
- **WHEN** user presses 'n' in current pane
- **THEN** a centered popup overlay appears with focus on title input
- **AND** the board remains visible behind the popup

#### Scenario: Submit new task
- **WHEN** user enters a title and presses Enter
- **THEN** a new task is created in the current pane
- **AND** the popup closes

#### Scenario: Cancel task creation
- **WHEN** user presses Escape during task creation
- **THEN** the popup closes without creating a task
- **AND** the board view is restored
