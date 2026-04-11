## ADDED Requirements

### Requirement: Reusable color palette package
The project SHALL provide a `pkg/palette` Go package that can be imported by any TUI project for consistent branding.

#### Scenario: Import palette in external project
- **WHEN** another Go TUI project imports `github.com/coliva/tsk/pkg/palette`
- **THEN** it has access to all color constants as strings
- **AND** it can use DarkTheme and LightTheme structs

#### Scenario: No external dependencies
- **WHEN** pkg/palette is imported
- **THEN** it does not require lipgloss or other TUI libraries
- **AND** colors are plain string constants (hex format)

### Requirement: Canonical color palette based on Catppuccin Mocha
The project SHALL use a unified color palette based on Catppuccin Mocha across TUI and website.

#### Scenario: Background colors
- **WHEN** rendering backgrounds
- **THEN** primary background uses `#1e1e2e`
- **AND** surface/secondary background uses `#313244`
- **AND** elevated/card background uses `#45475a`

#### Scenario: Text colors
- **WHEN** rendering text
- **THEN** primary text uses `#cdd6f4`
- **AND** secondary text uses `#a6adc8`
- **AND** muted text uses `#6c7086`

#### Scenario: Accent color
- **WHEN** rendering interactive or highlighted elements
- **THEN** accent color is `#74c7ec` (sapphire blue)
- **AND** accent hover state is `#8fd4f0`

#### Scenario: Semantic colors
- **WHEN** rendering status indicators
- **THEN** success color is `#a6e3a1`
- **AND** warning color is `#f9e2af`
- **AND** error/danger color is `#f38ba8`

#### Scenario: Border colors
- **WHEN** rendering borders
- **THEN** primary border uses `#585b70`
- **AND** active/focus border uses `#89b4fa`
