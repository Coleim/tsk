## 1. Model Layer

- [x] 1.1 Define SortMode enum in internal/model/state.go with values: SortCreatedDesc, SortCreatedAsc, SortDueDateAsc, SortDueDateDesc, SortTitleAsc, SortTitleDesc, SortPriorityDesc, SortPriorityAsc
- [x] 1.2 Add SortMode field to AppState struct
- [x] 1.3 Add String() method to SortMode for display names
- [x] 1.4 Add SortMode field to Board struct for persistence

## 2. Sorting Logic

- [x] 2.1 Create sortTasks() function in internal/model/state.go
- [x] 2.2 Implement sorting by creation date (both directions)
- [x] 2.3 Implement sorting by due date with nil handling (nil last for asc, nil first for desc)
- [x] 2.4 Implement sorting by title (case-insensitive)
- [x] 2.5 Implement sorting by priority
- [x] 2.6 Integrate sortTasks() into CurrentTasks() after filtering

## 3. Sort Selector UI

- [x] 3.1 Create SortSelector struct in internal/ui/sort.go
- [x] 3.2 Add list of sort options with current selection highlight
- [x] 3.3 Implement HandleKey for navigation (j/k), selection (Enter), cancel (Esc)
- [x] 3.4 Implement View() with popup styling matching filter dialog

## 4. App Integration

- [x] 4.1 Add ModeSort to internal/model/mode.go
- [x] 4.2 Replace existing 's' keybinding in app.go to open sort selector (currently calls Board.SortByPriority)
- [x] 4.3 Add handleSortMode() in app.go
- [x] 4.4 Apply selected sort mode to state when Enter pressed
- [x] 4.5 Render sort selector popup in View()
- [x] 4.6 Remove or deprecate Board.SortByPriority() method (no longer used)

## 5. Sort Indicator Display

- [x] 5.1 Add SortIndicatorStyle() to styles.go
- [x] 5.2 Display current sort mode in header when non-default
- [x] 5.3 Show sort hint in help overlay

## 6. Testing

- [x] 6.1 Add unit tests for sortTasks() with each sort mode
- [x] 6.2 Test nil due date handling in sorting
- [x] 6.3 Test sort combined with filter
- [x] 6.4 Manual testing of sort selector UI
