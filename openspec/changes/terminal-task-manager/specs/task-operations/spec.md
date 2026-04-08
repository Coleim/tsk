## ADDED Requirements

### Requirement: User can create tasks
The application SHALL allow users to create new tasks with a title.

#### Scenario: Create task with title
- **WHEN** user presses 'n' in current pane
- **THEN** a task creation modal appears with focus on title input

#### Scenario: Submit new task
- **WHEN** user enters a title and presses Enter
- **THEN** a new task is created in the current pane

#### Scenario: Cancel task creation
- **WHEN** user presses Escape during task creation
- **THEN** the modal closes without creating a task

### Requirement: User can edit task title
The application SHALL allow users to edit the title of existing tasks.

#### Scenario: Edit task title
- **WHEN** user presses Enter on a selected task to open detail view
- **THEN** user can edit the title inline

#### Scenario: Save edited title
- **WHEN** user modifies the title and presses Enter
- **THEN** the task title is updated

### Requirement: User can edit task description
The application SHALL allow users to add and edit a description for tasks.

#### Scenario: Add description
- **WHEN** user opens task detail and navigates to description field
- **THEN** user can type a multi-line description

#### Scenario: View description in detail
- **WHEN** user opens a task with a description
- **THEN** the full description is visible in the detail view

### Requirement: User can set task priority
The application SHALL allow users to set priority levels (High, Medium, Low, None) on tasks.

#### Scenario: Set priority
- **WHEN** user presses '1', '2', '3', or '0' on a selected task
- **THEN** task priority is set to High, Medium, Low, or None respectively

#### Scenario: Priority visual indicator
- **WHEN** task has a priority set
- **THEN** the task displays with a colored indicator matching priority

### Requirement: User can add labels to tasks
The application SHALL allow users to add text labels to tasks for categorization.
Labels are stored at the board level for reusability across tasks, with each label having a unique color.

#### Scenario: Add label
- **WHEN** user presses 'L' on a task and enters a label name
- **THEN** the label is added to the task
- **AND** if the label is new, it is created at the board level with an auto-assigned color

#### Scenario: Remove label
- **WHEN** user presses 'L' and selects an existing label
- **THEN** the label is removed from the task

#### Scenario: Labels display on task
- **WHEN** task has labels
- **THEN** labels appear as colored badges below the task title
- **AND** each label displays with its assigned color (red, orange, yellow, green, blue, purple, pink, or cyan)

#### Scenario: Label color consistency
- **WHEN** the same label is used on multiple tasks
- **THEN** the label displays with the same color on all tasks

### Requirement: User can set due dates
The application SHALL allow users to set due dates on tasks.

#### Scenario: Set due date
- **WHEN** user presses 't' on a task and enters a date
- **THEN** the due date is set on the task

#### Scenario: Due date display
- **WHEN** task has a due date
- **THEN** date appears on the task card

#### Scenario: Overdue indicator
- **WHEN** task due date is in the past
- **THEN** due date displays in red with "overdue" styling

### Requirement: User can delete tasks
The application SHALL allow users to delete tasks with confirmation.

#### Scenario: Delete task
- **WHEN** user presses 'd' on a selected task
- **THEN** a confirmation appears: "Delete task '<title>'?"

#### Scenario: Confirm deletion
- **WHEN** user confirms deletion
- **THEN** the task is removed from the board

#### Scenario: Cancel deletion
- **WHEN** user cancels deletion
- **THEN** the task remains and focus returns to it

### Requirement: User can move tasks between panes
The application SHALL allow users to move tasks to different panes (status change).

#### Scenario: Move task right
- **WHEN** user presses '>' on a task
- **THEN** task moves to the next pane (To Do → In Progress → Done)

#### Scenario: Move task left
- **WHEN** user presses '<' on a task
- **THEN** task moves to the previous pane (Done → In Progress → To Do)

#### Scenario: Follow moved task
- **WHEN** user moves a task to another pane
- **THEN** the view switches to that pane with the moved task selected

#### Scenario: Moved task visible
- **WHEN** user moves a task to a pane with many tasks
- **THEN** the list scrolls to show the moved task as selected
