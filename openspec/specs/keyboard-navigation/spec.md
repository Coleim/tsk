## ADDED Requirements

### Requirement: Vim-style navigation with j/k
The application SHALL support vim-style navigation keys for tasks.

#### Scenario: Navigate down with j
- **WHEN** user presses 'j'
- **THEN** selection moves to the next task down in the current pane

#### Scenario: Navigate up with k
- **WHEN** user presses 'k'
- **THEN** selection moves to the previous task up in the current pane

#### Scenario: Go to bottom with G
- **WHEN** user presses 'G' (shift+g)
- **THEN** selection jumps to the last task in the current pane

#### Scenario: Go to top with g
- **WHEN** user presses 'g' twice (gg)
- **THEN** selection jumps to the first task in the current pane

### Requirement: Pane switching with h/l
The application SHALL support h/l to switch between panes.

#### Scenario: Switch to previous pane with h
- **WHEN** user presses 'h'
- **THEN** view switches to the previous pane (Done → In Progress → To Do)

#### Scenario: Switch to next pane with l
- **WHEN** user presses 'l'
- **THEN** view switches to the next pane (To Do → In Progress → Done)

### Requirement: Arrow key navigation
The application SHALL support arrow keys as alternatives to vim keys.

#### Scenario: Navigate with arrows
- **WHEN** user presses up/down arrow keys
- **THEN** navigation behaves the same as j/k

#### Scenario: Switch panes with arrows
- **WHEN** user presses left/right arrow keys
- **THEN** pane switching behaves the same as h/l

### Requirement: Enter key opens task detail
The application SHALL use Enter to open a full-screen task detail view for viewing and editing.

#### Scenario: Enter opens task
- **WHEN** user presses Enter on a task
- **THEN** a full-screen task detail view opens for viewing and editing

### Requirement: New task with n key
The application SHALL create new tasks with the 'n' key.

#### Scenario: New task
- **WHEN** user presses 'n'
- **THEN** new task input appears in current pane

### Requirement: Delete task with d key
The application SHALL delete tasks with the 'd' key after confirmation.

#### Scenario: Delete task
- **WHEN** user presses 'd' on a task
- **THEN** delete confirmation appears

### Requirement: Move task with > and < keys
The application SHALL move tasks between panes with > and <.

#### Scenario: Move task right
- **WHEN** user presses '>' on a task
- **THEN** task moves to the next pane to the right

#### Scenario: Move task left
- **WHEN** user presses '<' on a task
- **THEN** task moves to the previous pane to the left

### Requirement: Escape key cancels and closes
The application SHALL use Escape to cancel or close overlays.

#### Scenario: Escape closes modals
- **WHEN** user presses Escape in a modal or detail view
- **THEN** the view closes and returns to board

### Requirement: Number keys for priority
The application SHALL use number keys 0-3 for quick priority assignment.

#### Scenario: Set high priority
- **WHEN** user presses '1' on a selected task
- **THEN** task priority is set to High

#### Scenario: Set medium priority
- **WHEN** user presses '2' on a selected task
- **THEN** task priority is set to Medium

#### Scenario: Set low priority
- **WHEN** user presses '3' on a selected task
- **THEN** task priority is set to Low

#### Scenario: Clear priority
- **WHEN** user presses '0' on a selected task
- **THEN** task priority is set to None

### Requirement: Label management with L key
The application SHALL manage labels with the 'L' key.

#### Scenario: Open label manager
- **WHEN** user presses 'L' on a selected task
- **THEN** a full-screen label management dialog opens to add or remove labels
- **AND** existing labels are displayed with their assigned colors

### Requirement: Search with slash key
The application SHALL activate search with the '/' key.

#### Scenario: Activate search
- **WHEN** user presses '/'
- **THEN** search input is focused and results update in real-time

#### Scenario: Search matches
- **WHEN** user types a query
- **THEN** tasks matching title, description, or labels are shown (case-insensitive)

#### Scenario: Navigate search results
- **WHEN** search results are displayed
- **THEN** user can navigate with j/k and press Enter to go to the task

#### Scenario: Clear search
- **WHEN** user presses Esc during search
- **THEN** search is cleared and view returns to normal board

### Requirement: Filter tasks
The application SHALL allow filtering by priority and labels.

#### Scenario: Open filter menu
- **WHEN** user presses 'f'
- **THEN** filter menu appears with priority and label options

#### Scenario: Apply filters
- **WHEN** user selects filters and presses Enter
- **THEN** only matching tasks are displayed

#### Scenario: Clear filters
- **WHEN** user presses 'F'
- **THEN** all filters are cleared

### Requirement: Sort by priority with s key
The application SHALL open a sort mode selector popup when user presses 's' key.

#### Scenario: Open sort selector
- **WHEN** user presses 's'
- **THEN** a sort selector popup appears with all available sort options
- **AND** the currently active sort mode is highlighted

#### Scenario: Select sort mode
- **WHEN** user navigates to a sort option with j/k and presses Enter
- **THEN** tasks in current pane are sorted according to the selected mode
- **AND** the popup closes

#### Scenario: Cancel sort selection
- **WHEN** user presses Escape in the sort selector
- **THEN** the popup closes without changing sort mode

#### Scenario: Sort options available
- **WHEN** sort selector is open
- **THEN** user can choose from: Created (Newest), Created (Oldest), Due Date (Soonest), Due Date (Latest), Title (A-Z), Title (Z-A), Priority (High-Low), Priority (Low-High)

### Requirement: Undo with u key
The application SHALL undo the last action with the 'u' key.

#### Scenario: Undo action
- **WHEN** user presses 'u' and undo stack is not empty
- **THEN** last action is reversed and status bar shows "↩ Undid <action>"

#### Scenario: Nothing to undo
- **WHEN** user presses 'u' and undo stack is empty
- **THEN** status bar shows "Nothing to undo"

### Requirement: Redo with Ctrl+r
The application SHALL redo the last undone action with Ctrl+r.

#### Scenario: Redo action
- **WHEN** user presses Ctrl+r and redo stack is not empty
- **THEN** last undone action is re-applied and status bar shows "↪ Redid <action>"

#### Scenario: Nothing to redo
- **WHEN** user presses Ctrl+r and redo stack is empty
- **THEN** status bar shows "Nothing to redo"

### Requirement: Quit application with q
The application SHALL allow quitting with 'q' or ':wq'.

#### Scenario: Quit from main view
- **WHEN** user presses 'q' or types ':wq' on main board view
- **THEN** application exits

### Requirement: Help overlay with ? key
The application SHALL show help with the '?' key.

#### Scenario: Show help
- **WHEN** user presses '?'
- **THEN** a full-screen help dialog appears showing all keyboard shortcuts grouped by category

#### Scenario: Close help with ?
- **WHEN** help dialog is open and user presses '?'
- **THEN** help dialog closes (toggle behavior)

#### Scenario: Close help with Escape
- **WHEN** help dialog is open and user presses Escape
- **THEN** help dialog closes

### Requirement: Archive tasks with a/A keys
The application SHALL allow archiving completed tasks when in DONE pane.

#### Scenario: Archive selected task
- **WHEN** user is in DONE pane and presses 'a'
- **THEN** the selected task is moved to archive file
- **AND** status bar shows "Archived: <task title> (u to undo)"

#### Scenario: Undo archive
- **WHEN** user presses 'u' after archiving a task
- **THEN** the task is restored to DONE pane
- **AND** status bar shows "Restored: <task title>"

#### Scenario: Archive all done tasks
- **WHEN** user is in DONE pane and presses 'A'
- **THEN** all tasks in DONE pane are moved to archive file
- **AND** status bar shows "Archived N tasks (u to undo)"

#### Scenario: Undo archive all
- **WHEN** user presses 'u' after archiving all done tasks
- **THEN** all archived tasks are restored to DONE pane
- **AND** status bar shows "Restored N tasks"

#### Scenario: Archive keys not available outside DONE
- **WHEN** user is in TO DO or IN PROGRESS pane
- **THEN** 'a' and 'A' keys have no effect
- **AND** archive keys are not shown in status bar
