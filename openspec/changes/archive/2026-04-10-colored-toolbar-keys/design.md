## Context

The application uses Charm's bubbletea/lipgloss stack. The bottom status bar displays context-sensitive shortcuts as plain text with a single muted color (`TextMuted`). The accent color (`CurrentTheme.Accent` - purple `#cba6f7` in dark theme) is used throughout the app for highlighting selected items and important elements.

Current shortcut format: `"j/k:nav  h/l:pane  n:new  d:del  ..."`

## Goals / Non-Goals

**Goals:**
- Render shortcut keys with the accent color for better visual scanning
- Maintain consistent styling with the rest of the application
- Keep implementation minimal and localized

**Non-Goals:**
- Changing the shortcut format or content
- Adding new shortcut categories or groups
- Modifying the help overlay styling

## Decisions

### 1. Add `ShortcutKeyStyle()` function in styles.go
Use the existing `Accent` color from the theme for consistency. This follows the pattern of other style functions in the codebase.

**Alternative considered**: Inline styling in app.go - rejected because it violates the separation of concerns established by the styles package.

### 2. Create helper function to format shortcuts
Add a `formatShortcut(key, action string)` helper that applies colored key styling.

**Alternative considered**: Regex-based parsing of the full shortcut string - rejected as overly complex and brittle.

### 3. Update shortcut string construction
Replace hardcoded shortcut strings with function calls to `formatShortcut()`.

## Risks / Trade-offs

**[Risk]** More verbose code for shortcut construction → Accepted trade-off for better visual design and maintainability.

**[Risk]** Slight increase in render complexity → Negligible impact; lipgloss style operations are efficient.
