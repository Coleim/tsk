## MODIFIED Requirements

### Requirement: Application provides keyboard hint bar
The application SHALL provide a consistent footer bar for keyboard shortcuts across all dialogs and popups.

#### Scenario: Hint bar rendering
- **WHEN** a dialog or popup is displayed
- **THEN** a keyboard hint bar appears at the bottom
- **AND** key bindings are formatted as "key:action" with spaces between items
- **AND** the key portion is rendered with the accent color
- **AND** the colon and action portions are rendered with muted text color

#### Scenario: Hint bar separator
- **WHEN** a hint bar is rendered in a full-screen dialog
- **THEN** a thin horizontal line appears above the hint bar for visual separation
