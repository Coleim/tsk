## MODIFIED Requirements

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
