# Task Statistics

Statistics overlay for visualizing board metrics and task distribution.

## Requirements

### Requirement: Statistics overlay displays task distribution
The application SHALL display a statistics overlay showing task counts by status as horizontal bar graphs.

#### Scenario: Open statistics overlay
- **WHEN** user presses 'S' (shift+s) on the main board view
- **THEN** a full-screen statistics overlay appears showing board metrics

#### Scenario: Status distribution bars
- **WHEN** statistics overlay is displayed
- **THEN** horizontal bar graphs show task counts for To Do, In Progress, and Done
- **AND** bars are proportional to the total task count
- **AND** each bar displays the numeric count

#### Scenario: Close statistics overlay
- **WHEN** user presses Escape or 'S' in the statistics overlay
- **THEN** the overlay closes and returns to the main board view

### Requirement: Statistics overlay displays priority breakdown
The application SHALL display priority distribution with color-coded visualization.

#### Scenario: Priority distribution display
- **WHEN** statistics overlay is displayed
- **THEN** a priority breakdown shows counts for High, Medium, Low, and None priorities
- **AND** each priority uses its established color from the theme
- **AND** bars or indicators show relative proportions

### Requirement: Statistics overlay displays due date metrics
The application SHALL display due date statistics including overdue and upcoming tasks.

#### Scenario: Due date statistics display
- **WHEN** statistics overlay is displayed
- **THEN** the overlay shows count of overdue tasks (past due date)
- **AND** the overlay shows count of tasks due today
- **AND** the overlay shows count of tasks due this week
- **AND** the overlay shows count of tasks with no due date

#### Scenario: Overdue warning styling
- **WHEN** there are overdue tasks
- **THEN** the overdue count is displayed with warning/error styling

### Requirement: Statistics overlay displays completion summary
The application SHALL display a summary of task completion.

#### Scenario: Completion percentage
- **WHEN** statistics overlay is displayed
- **THEN** an overall completion percentage is shown (Done / Total tasks)
- **AND** the percentage is displayed with a progress indicator

#### Scenario: Empty board statistics
- **WHEN** statistics overlay is displayed for a board with no tasks
- **THEN** the overlay shows appropriate empty state messaging
- **AND** all metrics display zero values gracefully
