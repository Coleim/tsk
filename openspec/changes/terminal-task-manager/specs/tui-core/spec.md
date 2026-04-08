## ADDED Requirements

### Requirement: Application renders in terminal
The application SHALL render a full-screen terminal user interface using the Bubbletea framework.

#### Scenario: Application startup
- **WHEN** user runs the application command
- **THEN** a full-screen TUI appears with the single-pane board view and preview panel

#### Scenario: Terminal resize handling
- **WHEN** user resizes the terminal window
- **THEN** the UI re-renders to fit the new dimensions

### Requirement: Application displays two-panel layout
The application SHALL display a two-panel layout: task list (left) and preview (right).

#### Scenario: Task list panel
- **WHEN** application is running
- **THEN** left panel shows tasks from the current pane (To Do, In Progress, or Done)

#### Scenario: Preview panel
- **WHEN** a task is selected
- **THEN** right panel shows task details (title, status, priority, due date, labels, description)

#### Scenario: Panel proportions
- **WHEN** displaying the layout
- **THEN** task list takes approximately 60% width, preview takes 40%

### Requirement: Application displays pane tabs
The application SHALL display tabs showing all panes with the current pane highlighted.

#### Scenario: Pane tabs display
- **WHEN** application is running
- **THEN** tabs show "[TO DO] (3)  IN PROGRESS (2)  DONE (4)" with current pane bracketed

#### Scenario: Tab updates on switch
- **WHEN** user switches panes with h/l
- **THEN** the bracket indicator moves to the new current pane

### Requirement: Application displays header
The application SHALL display a header bar showing the application name, current board name, and available keyboard shortcuts hint.

#### Scenario: Header displays application info
- **WHEN** application is running
- **THEN** header shows "tsk" title, current board name, and "Press ? for help"

### Requirement: Application displays two-line status bar
The application SHALL display a two-line status bar at the bottom with context and shortcuts.

#### Scenario: Status bar line 1 shows context
- **WHEN** no recent action
- **THEN** line 1 shows task count in current pane (e.g., "3 tasks in TO DO")

#### Scenario: Status bar line 1 shows feedback
- **WHEN** user completes an action
- **THEN** line 1 shows action feedback and undo/redo hints (e.g., "✓ Deleted task • u:undo")

#### Scenario: Status bar line 2 shows shortcuts
- **WHEN** application is running
- **THEN** line 2 always shows keyboard shortcuts for current context

#### Scenario: Feedback auto-clears
- **WHEN** 3 seconds pass after action feedback
- **THEN** line 1 returns to showing task count

### Requirement: Application supports modal dialogs
The application SHALL support modal dialogs for task creation, editing, and confirmations.

#### Scenario: Modal overlays main view
- **WHEN** user triggers a modal action (new task, edit, delete)
- **THEN** a modal dialog appears centered over the main view with focus

#### Scenario: Modal dismissal
- **WHEN** user presses Escape in a modal
- **THEN** the modal closes and focus returns to the previous view

### Requirement: Application supports color themes
The application SHALL use semantic colors for different task priorities and states.

#### Scenario: Priority colors displayed
- **WHEN** tasks have different priorities
- **THEN** high priority shows in red, medium in yellow, low in green, none in white

### Requirement: Application handles errors gracefully
The application SHALL display user-friendly error messages when operations fail.

#### Scenario: Error notification
- **WHEN** an operation fails (e.g., save error)
- **THEN** an error message appears in status bar line 1 with details
