## ADDED Requirements

### Requirement: Application renders in terminal
The application SHALL render a full-screen terminal user interface using the Bubbletea framework.

#### Scenario: Application startup
- **WHEN** user runs the application command
- **THEN** a loading indicator with animated spinner appears while loading tasks
- **AND THEN** a full-screen TUI appears with the single-pane board view and preview panel

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
- **THEN** right panel shows task details (title, status, priority, due date, colored label badges, description)

#### Scenario: Panel proportions
- **WHEN** displaying the layout
- **THEN** task list takes approximately 30% width, preview takes 70%
- **AND** both panels have the same height

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

### Requirement: Application supports full-screen dialogs
The application SHALL support full-screen dialogs for task creation, editing, labels, filters, and confirmations.
All dialogs use the full terminal width and height for maximum usability.

#### Scenario: Full-screen dialog display
- **WHEN** user triggers a dialog action (new task, edit, delete, labels, filter, board switch, due date, help)
- **THEN** a full-screen dialog appears filling the terminal with focus
- **AND** the dialog uses the full terminal width and height

#### Scenario: Dialog dismissal
- **WHEN** user presses Escape in a dialog
- **THEN** the dialog closes and focus returns to the previous view

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

### Requirement: Application shows loading state
The application SHALL display a loading indicator while loading board data.

#### Scenario: Loading indicator on startup
- **WHEN** application is loading board data
- **THEN** an animated spinner with "Loading tasks..." text is displayed
- **AND** the spinner is centered on the screen
