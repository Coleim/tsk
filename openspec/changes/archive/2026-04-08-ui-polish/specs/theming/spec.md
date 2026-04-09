## ADDED Requirements

### Requirement: Application supports theme selection
The application SHALL support dark and light color themes selectable via environment variable.

#### Scenario: Default dark theme
- **WHEN** application starts without TSK_THEME set
- **THEN** the dark color theme is applied

#### Scenario: Light theme selection
- **WHEN** application starts with TSK_THEME=light
- **THEN** the light color theme is applied

#### Scenario: Invalid theme value
- **WHEN** application starts with TSK_THEME set to an invalid value
- **THEN** the dark theme is applied as fallback

### Requirement: Application uses semantic colors
The application SHALL use semantic color names that map to theme-appropriate values.

#### Scenario: Semantic color consistency
- **WHEN** a component uses the "accent" semantic color
- **THEN** the color rendered matches the current theme's accent color

#### Scenario: Theme color inheritance
- **WHEN** theme is changed
- **THEN** all components using semantic colors update accordingly

### Requirement: Colors meet accessibility standards
The application SHALL use colors that provide sufficient contrast for readability.

#### Scenario: Text contrast in dark theme
- **WHEN** using dark theme
- **THEN** primary text has minimum 7:1 contrast ratio against background

#### Scenario: Text contrast in light theme
- **WHEN** using light theme
- **THEN** primary text has minimum 7:1 contrast ratio against background

### Requirement: Theme provides elevation levels
The application SHALL use distinct background colors to create visual depth hierarchy.

#### Scenario: Background level display
- **WHEN** rendering the main view
- **THEN** the base background uses the theme's background color

#### Scenario: Surface level display
- **WHEN** rendering panels and cards
- **THEN** they use the theme's surface color (slightly elevated from background)

#### Scenario: Elevated level display
- **WHEN** rendering modals and popups
- **THEN** they use the theme's elevated color (most prominent)
