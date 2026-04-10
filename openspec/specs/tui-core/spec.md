## MODIFIED Requirements

### Requirement: Application supports full-screen dialogs
The application SHALL support full-screen dialogs for complex editing operations like task editing, labels management, and filters.
Full-screen dialogs SHALL use a consistent three-zone layout (header, content, footer) and visual grouping for related content.

#### Scenario: Full-screen dialog display
- **WHEN** user triggers a complex editing action (edit task, labels, filter, board switch, due date, help)
- **THEN** a full-screen dialog appears filling the terminal with focus
- **AND** the dialog uses the full terminal width and height
- **AND** the dialog displays a header zone with the dialog title in accent color

#### Scenario: Dialog content zone
- **WHEN** a full-screen dialog has multiple sections
- **THEN** related content is grouped in section cards with rounded borders
- **AND** section cards have consistent internal padding

#### Scenario: Dialog footer zone
- **WHEN** a full-screen dialog is displayed
- **THEN** a keyboard hint bar appears at the bottom
- **AND** a thin separator line appears above the hint bar

#### Scenario: Dialog dismissal
- **WHEN** user presses Escape in a dialog
- **THEN** the dialog closes and focus returns to the previous view

### Requirement: Application supports popup overlays
The application SHALL support compact popup overlays for simple inputs and confirmations, with the board visible behind.
Popup overlays SHALL use consistent widths based on content type and have uniform styling.

#### Scenario: Popup overlay display
- **WHEN** user triggers a simple action (new task, delete confirmation, search)
- **THEN** a centered popup overlay appears
- **AND** the board remains visible behind the popup
- **AND** the popup has a double border in accent color

#### Scenario: Popup width standards
- **WHEN** displaying a popup overlay
- **THEN** simple input popups use 50 character width
- **AND** list-based popups use 50-60 character width
- **AND** date picker popups use 60 character width

#### Scenario: Popup dismissal
- **WHEN** user presses Escape in a popup
- **THEN** the popup closes and the board view is restored

#### Scenario: Delete confirmation popup
- **WHEN** user presses 'd' to delete a task
- **THEN** a confirmation popup appears asking to confirm deletion

## ADDED Requirements

### Requirement: Filter dialog uses visual sections
The filter dialog SHALL display priority and label filter options in visually distinct section cards.

#### Scenario: Filter sections display
- **WHEN** the filter dialog is opened
- **THEN** a "Priority" section card appears with checkbox items for High, Medium, Low, None
- **AND** a "Labels" section card appears below with checkbox items for available labels
- **AND** each section has a header with accent-colored underline

#### Scenario: Filter active summary
- **WHEN** one or more filters are selected
- **THEN** an "Active filters" summary appears at the top in success color
- **AND** the summary lists all selected filter criteria

### Requirement: Edit task dialog uses form sections
The edit task dialog SHALL organize form fields with visual grouping and clear focus indication.

#### Scenario: Edit form layout
- **WHEN** the edit task dialog is opened
- **THEN** form fields (Title, Description, Labels) appear with labeled sections
- **AND** the currently focused field shows a "▶" indicator and accent color label

#### Scenario: Edit form section cards
- **WHEN** multiple fields are displayed
- **THEN** each field appears in its own visual row with consistent spacing
- **AND** a thin divider line separates the form area from keyboard hints

### Requirement: Board selector uses board cards
The board selector dialog SHALL display boards as visually distinct cards with metadata.

#### Scenario: Board card display
- **WHEN** the board selector is opened
- **THEN** each board appears as a card with rounded border
- **AND** the card shows the board name prominently
- **AND** the card shows task count as secondary information

#### Scenario: Current board indicator
- **WHEN** the currently loaded board is displayed in the list
- **THEN** a "•" bullet indicator appears before the board name
- **AND** the card uses a subtle highlight to distinguish it from other boards
