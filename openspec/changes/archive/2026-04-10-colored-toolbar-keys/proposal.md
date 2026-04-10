## Why

The bottom toolbar displays keyboard shortcuts in a uniform muted color, making it harder to quickly scan and identify keys. Highlighting the key portion with the accent color improves visual hierarchy and readability.

## What Changes

- Parse shortcut strings (e.g., "j/k:nav") to separate key and action portions
- Render keys with the accent color while keeping actions in the muted color
- Apply consistently across all context-sensitive shortcut displays

## Capabilities

### New Capabilities
- `shortcut-key-styling`: Styling logic to render keyboard shortcut keys with accent color

### Modified Capabilities
- `visual-components`: Add styling for shortcut keys in toolbar

## Impact

- `internal/ui/app.go`: Modify `renderStatusBar()` to use styled shortcut rendering
- `internal/styles/styles.go`: Add `ShortcutKeyStyle()` function for key styling
