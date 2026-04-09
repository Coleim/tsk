## MODIFIED Requirements

### Requirement: Application supports full-screen dialogs
The application SHALL support full-screen dialogs for complex editing operations like task editing, labels management, and filters.
Simple interactions use popup overlays instead.

#### Scenario: Full-screen dialog display
- **WHEN** user triggers a complex editing action (edit task, labels, filter, board switch, due date, help)
- **THEN** a full-screen dialog appears filling the terminal with focus
- **AND** the dialog uses the full terminal width and height

#### Scenario: Dialog dismissal
- **WHEN** user presses Escape in a dialog
- **THEN** the dialog closes and focus returns to the previous view

### Requirement: Application supports popup overlays
The application SHALL support compact popup overlays for simple inputs and confirmations, with the board visible behind.

#### Scenario: Popup overlay display
- **WHEN** user triggers a simple action (new task, delete confirmation, search)
- **THEN** a centered popup overlay appears
- **AND** the board remains visible behind the popup

#### Scenario: Popup dismissal
- **WHEN** user presses Escape in a popup
- **THEN** the popup closes and the board view is restored

#### Scenario: Delete confirmation popup
- **WHEN** user presses 'd' to delete a task
- **THEN** a confirmation popup appears asking to confirm deletion
- **AND** the board and selected task remain visible behind the popup
- **AND** user can press 'y' to confirm or 'n'/Escape to cancel
