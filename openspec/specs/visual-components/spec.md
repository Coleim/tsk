## ADDED Requirements

### Requirement: Application provides section card component
The application SHALL provide a reusable section card style that visually groups related content with a rounded border.

#### Scenario: Section card rendering
- **WHEN** a view renders a section card
- **THEN** the card displays with a rounded border in the theme's BorderLight color
- **AND** internal content has consistent padding (1 line vertical, 2 characters horizontal)

#### Scenario: Section card title
- **WHEN** a section card has a title
- **THEN** the title appears above the content with accent color styling
- **AND** a subtle separator line appears between title and content

### Requirement: Application provides form field styles
The application SHALL provide consistent visual styling for form input fields across all dialogs.

#### Scenario: Form field label
- **WHEN** a form field is rendered
- **THEN** the label appears in muted text color above the input

#### Scenario: Active form field
- **WHEN** a form field has focus
- **THEN** a "▶" indicator appears before the label
- **AND** the field label uses primary text color instead of muted

#### Scenario: Form field input display
- **WHEN** a text input is displayed
- **THEN** it renders with the standard Bubbletea textinput styling
- **AND** the placeholder text uses muted color

### Requirement: Application provides keyboard hint bar
The application SHALL provide a consistent footer bar for keyboard shortcuts across all dialogs and popups.

#### Scenario: Hint bar rendering
- **WHEN** a dialog or popup is displayed
- **THEN** a keyboard hint bar appears at the bottom
- **AND** key bindings are formatted as "key:action" with spaces between items
- **AND** the key portion is rendered with the accent color
- **AND** the colon and action portions are rendered with muted text color

#### Scenario: Hint bar separator
- **WHEN** a hint bar is rendered in a full-screen dialog
- **THEN** a thin horizontal line appears above the hint bar for visual separation

### Requirement: Application provides selection indicator styles
The application SHALL provide consistent visual indication for selected/focused items in lists.

#### Scenario: List item selection
- **WHEN** a list item is focused/selected
- **THEN** a "▶" arrow appears before the item text
- **AND** the item uses accent color for text
- **AND** a subtle surface background color highlights the row

#### Scenario: Non-selected list item
- **WHEN** a list item is not selected
- **THEN** the item uses secondary text color
- **AND** two spaces appear before the item (aligning with selected indicator width)

### Requirement: Application provides checkbox styles
The application SHALL provide consistent styling for checkbox lists (filters, multi-select).

#### Scenario: Unchecked checkbox
- **WHEN** a checkbox item is unchecked
- **THEN** it displays as "[ ]" followed by the item label

#### Scenario: Checked checkbox
- **WHEN** a checkbox item is checked
- **THEN** it displays as "[✓]" with the checkmark in success color
- **AND** the checkbox brackets use muted color

### Requirement: Application provides dialog layout helper
The application SHALL provide a standard three-zone layout function for full-screen dialogs.

#### Scenario: Dialog zones
- **WHEN** a full-screen dialog uses the standard layout
- **THEN** it renders three zones: header (title), content (scrollable body), footer (hints)
- **AND** the content zone expands to fill available height
- **AND** the dialog fills the terminal width with appropriate margins
