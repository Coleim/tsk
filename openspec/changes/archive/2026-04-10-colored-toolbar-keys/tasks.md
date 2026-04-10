## 1. Style Function

- [x] 1.1 Add `ShortcutKeyStyle()` function to `internal/styles/styles.go` that returns a style using `CurrentTheme.Accent`

## 2. Shortcut Formatting

- [x] 2.1 Add `formatShortcut(key, action string)` helper function in `internal/ui/app.go` that applies key styling
- [x] 2.2 Update Done pane shortcuts in `renderStatusBar()` to use `formatShortcut()`
- [x] 2.3 Update other pane shortcuts in `renderStatusBar()` to use `formatShortcut()`

## 3. Verification

- [x] 3.1 Run application and verify shortcut keys display with accent color in all panes
