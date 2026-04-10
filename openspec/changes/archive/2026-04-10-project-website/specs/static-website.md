## ADDED Requirements

### Requirement: Landing page with hero section
The website SHALL have a landing page with project name, tagline, and call-to-action buttons.

#### Scenario: User visits homepage
- **WHEN** user navigates to coleim.github.io/tsk
- **THEN** they see the project name "tsk", tagline, and Install/GitHub buttons

### Requirement: Features section
The website SHALL display key features with icons or illustrations.

#### Scenario: User views features
- **WHEN** user scrolls to features section
- **THEN** they see feature cards for: Kanban workflow, Vim navigation, Multiple boards, Priorities, Labels, Due dates

### Requirement: Installation section
The website SHALL provide installation commands for all supported methods.

#### Scenario: User views installation options
- **WHEN** user views installation section
- **THEN** they see copy-able commands for: Homebrew, curl script, go install, and source build

### Requirement: Screenshots section
The website SHALL display screenshots of the application.

#### Scenario: User views screenshots
- **WHEN** user scrolls to screenshots section
- **THEN** they see at least 2 screenshots showing: main view with tasks, and task editing/popups

### Requirement: Responsive design
The website SHALL be usable on mobile devices.

#### Scenario: Mobile viewport
- **WHEN** user views site on mobile (< 768px width)
- **THEN** layout adjusts to single column, navigation remains accessible
