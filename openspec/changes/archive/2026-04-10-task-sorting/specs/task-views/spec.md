## ADDED Requirements

### Requirement: Task list respects current sort order
The application SHALL display tasks in the currently selected sort order.

#### Scenario: Sorted task display
- **WHEN** a sort mode is active
- **THEN** the task list displays tasks in that sorted order
- **AND** sorting applies within each status pane independently

#### Scenario: Sort combined with filter
- **WHEN** both sort and filter are active
- **THEN** filtering is applied first, then sorting is applied to filtered results
