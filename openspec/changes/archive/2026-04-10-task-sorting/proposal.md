## Why

Users need to organize and prioritize their task list beyond the current priority filter. Sorting by creation date helps find recent or old tasks, sorting by due date helps manage deadlines, and alphabetical sorting aids task discovery in large boards.

## What Changes

- Add sort mode selector accessible from main board view
- Support sorting by: creation date (newest/oldest), due date (soonest/latest), title (A-Z/Z-A), priority (high to low/low to high)
- Display current sort indicator in UI
- Persist sort preference per board

## Capabilities

### New Capabilities
- `task-sorting`: Sorting options for task lists by various criteria (date, title, priority)

### Modified Capabilities
- `task-views`: Add sort indicator display and sorted task list rendering
- `keyboard-navigation`: Update 's' key from single priority sort to full sort selector popup

## Impact

- `internal/model/state.go`: Add sort mode field to AppState
- `internal/model/board.go`: Add sort preference field
- `internal/ui/app.go`: Add sort mode handler and keyboard binding
- `internal/ui/sort.go`: New sort mode selector component
- `internal/styles/styles.go`: Add sort indicator styles
