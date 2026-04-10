## ADDED Requirements

### Requirement: Shortcut keys are displayed with accent color
The toolbar SHALL render keyboard shortcut keys with the application accent color (`#cba6f7` in dark theme, `#8839ef` in light theme) while displaying action descriptions in the muted text color (`#6c7086` in dark theme, `#9ca0b0` in light theme).

#### Scenario: Shortcut displayed in toolbar
- **WHEN** the status bar renders a shortcut like "j/k:nav"
- **THEN** the key portion "j/k" is rendered with the accent color (`#cba6f7` dark / `#8839ef` light)
- **AND** the action portion "nav" is rendered with the muted text color (`#6c7086` dark / `#9ca0b0` light)
- **AND** the colon separator is rendered with the muted text color

### Requirement: All toolbar shortcuts use consistent key styling
The toolbar SHALL apply accent color styling to all shortcut keys consistently across different panes and modes.

#### Scenario: Done pane shortcuts
- **WHEN** viewing the Done pane
- **THEN** all shortcut keys (j/k, h/l, d, a, A, Enter, b) are rendered with accent color

#### Scenario: Other pane shortcuts
- **WHEN** viewing ToDo or In Progress panes
- **THEN** all shortcut keys (j/k, h/l, n, d, >/<, Enter, 1-3, b) are rendered with accent color
