## MODIFIED Requirements

### Requirement: User can search tasks
The application SHALL allow users to search tasks using a compact popup overlay with starts-with matching on titles.

#### Scenario: Open search
- **WHEN** user presses '/'
- **THEN** a centered popup overlay appears with search input
- **AND** the board remains visible behind the popup

#### Scenario: Search matching
- **WHEN** user types search text
- **THEN** tasks whose title starts with the search text are shown first
- **AND** tasks matching in description or labels (contains) are also included

#### Scenario: Search results display
- **WHEN** search has results
- **THEN** maximum 10 results are visible in the popup
- **AND** results can be scrolled with j/k if more exist

#### Scenario: Close search
- **WHEN** user presses Escape in search
- **THEN** search popup closes and full board view is restored
