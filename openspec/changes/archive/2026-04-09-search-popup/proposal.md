## Why

The current search is a full-screen overlay that hides the board context. A compact popup would be less disruptive and more consistent with modern search UIs. Additionally, "starts with" matching provides more intuitive results for quick task lookup.

## What Changes

- Replace full-screen search dialog with a centered popup overlay
- Change search matching from "contains" to "starts with" for task titles
- Keep the board visible behind the search popup for context
- Limit visible results to improve popup compactness

## Capabilities

### New Capabilities
None

### Modified Capabilities
- `task-views`: Search behavior changes from full-screen to popup overlay, and matching changes from contains to starts-with

## Impact

- `internal/ui/search.go`: Popup rendering, starts-with filtering
- `internal/ui/app.go`: Search overlay positioning
