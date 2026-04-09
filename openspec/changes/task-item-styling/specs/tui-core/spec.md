## ADDED Requirements

### Requirement: Task items display with visual distinction
The application SHALL render task items with clear visual separation and proper spacing between elements.

#### Scenario: Selected task highlighting
- **WHEN** a task is selected in the task list
- **THEN** the selected task row displays with a background highlight
- **AND** the selection indicator `▶` appears with spacing before the priority icon

#### Scenario: Task element spacing
- **WHEN** rendering a task item
- **THEN** there SHALL be at least one space between the selection indicator and priority icon
- **AND** there SHALL be at least one space between the priority icon and task title

#### Scenario: Unselected task display
- **WHEN** a task is not selected
- **THEN** it displays without background highlighting
- **AND** maintains consistent spacing with selected tasks
