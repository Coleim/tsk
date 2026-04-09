## Context

The TUI uses `overlayDialog()` in app.go for popup compositing - it layers dialog content over a dimmed background using Lipgloss's Layer/Compositor. Currently ModeSearch, ModeInsert, and ModeModal use this pattern, but ModeDueDate and ModeLabels render full-screen.

Style definitions are in `internal/styles/styles.go`. The popup label styles use `CurrentTheme.Accent` as background for selected items, which is visually harsh. Unselected items show background color leaking past label badges.

## Goals / Non-Goals

**Goals:**
- Consistent popup presentation for Due Date and Labels editors
- Clean selection styling without jarring brightness or background artifacts
- Maintain existing keyboard navigation and behavior

**Non-Goals:**
- Changing the DueDateEditor or LabelEditor internal logic
- Refactoring the overlayDialog function itself
- Adding new features to these editors

## Decisions

### 1. Popup Dimensions
Due Date and Labels editors will use fixed-width popups (50-60 chars) with auto-height content, similar to ModeInsert popup. The View() methods will remove their width/height parameters and return compact content.

**Rationale**: Consistent with existing popup pattern. Centered overlays feel more modern than full-screen modals for quick editors.

### 2. Labels Selection: Arrow-only
Replace `TaskSelectedStyle()` with plain arrow prefix (`▶ `) for selected labels in LabelEditor.

**Rationale**: TaskSelectedStyle has a full border meant for task cards. For label lists, a simple arrow is cleaner and matches the edit popup pattern.

### 3. Popup Label Selection Brightness
Change `PopupSelectedItemStyle()` to use `CurrentTheme.Surface` or a softer accent instead of `CurrentTheme.Accent` for background.

**Rationale**: Accent colors are meant for small highlights, not full line backgrounds. Using Surface with bold/foreground-accent is less harsh.

### 4. Unselected Label Background Fix
Remove Padding from `PopupItemStyle()` or ensure consistent background handling. The issue is that PopupStyle() sets a background, and padding on items shows that background while the LabelBadge has no background.

**Rationale**: Labels should sit directly on the popup background without extraneous styling that creates visual artifacts.

## Risks / Trade-offs

**[Visual consistency]** → Changes to popup styles affect all places using PopupSelectedItemStyle. Verify search popup and other users still look correct.

**[Width calculation]** → Fixed popup widths may not fit very long label names. Keep width generous (60 chars) or use min-width approach.
