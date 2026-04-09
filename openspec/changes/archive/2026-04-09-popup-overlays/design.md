## Context

The application currently uses full-screen modals for task creation and delete confirmation. This creates a jarring experience where users lose sight of the board context. The search feature was recently refactored to use a centered popup overlay with the board visible behind it, using lipgloss v2's Layer/Compositor API. This pattern should be extended to task creation and delete confirmation.

Current state:
- `renderWithTextInput()` renders a full-screen modal for new task creation
- `Modal.View()` renders a full-screen modal for confirmations (delete, etc.)
- `overlayDialog()` function exists and uses lipgloss v2 Layer compositing for search popup

## Goals / Non-Goals

**Goals:**
- Task creation uses a compact centered popup overlay with board visible
- Delete confirmation uses a compact centered popup overlay with board visible
- Consistent visual pattern with the search popup
- Reuse the existing `overlayDialog()` function

**Non-Goals:**
- Changing the task edit view (remains full-screen for complex editing)
- Changing the task detail view (remains full-screen)
- Changing help overlay behavior
- Adding new functionality beyond the UI presentation change

## Decisions

### 1. Reuse overlayDialog for all popup overlays

**Decision:** Use the existing `overlayDialog(background, popup, width, height)` function for both task creation and delete confirmation.

**Rationale:** This function already handles lipgloss v2 Layer compositing correctly for the search popup. Reusing it ensures consistent behavior and reduces code duplication.

**Alternative considered:** Create separate overlay functions for each popup type. Rejected because the overlay logic is identical.

### 2. Popup sizing strategy

**Decision:** Use fixed widths appropriate to content:
- Task creation popup: 50 characters wide (matches search popup)
- Delete confirmation popup: 50 characters wide

**Rationale:** Fixed width provides consistent appearance. 50 chars is wide enough for titles and messages but narrow enough to clearly appear as a popup.

### 3. Render main view as background

**Decision:** In both cases, render `renderMainView()` first, then overlay the popup on top.

**Rationale:** This matches the search popup pattern and ensures users see their current board context.

## Risks / Trade-offs

**[Text input width limitation]** → The text input will be constrained to ~46 chars (50 - padding). Accept this as reasonable for task titles; users wanting longer titles can edit after creation.

**[Confirmation message truncation]** → Long confirmation messages may need wrapping. Mitigate by keeping delete confirmation message concise.
