## 1. Model Layer

- [x] 1.1 Create internal/model/statistics.go with BoardStatistics struct
- [x] 1.2 Add ComputeStatistics(board *Board) function returning BoardStatistics
- [x] 1.3 Implement status counts (ToDo, InProgress, Done)
- [x] 1.4 Implement priority counts (High, Medium, Low, None)
- [x] 1.5 Implement due date metrics (overdue, today, this week, no due date)
- [x] 1.6 Add completion percentage calculation

## 2. State Integration

- [x] 2.1 Add ModeStats to internal/model/mode.go
- [x] 2.2 Add ShowStats bool to AppState struct
- [x] 2.3 Add String() case for ModeStats returning "STATS"

## 3. Statistics View Component

- [x] 3.1 Create internal/ui/stats.go with Stats struct
- [x] 3.2 Implement renderBarGraph() helper for horizontal bars using Unicode blocks
- [x] 3.3 Implement renderStatusSection() showing To Do/In Progress/Done bars
- [x] 3.4 Implement renderPrioritySection() with color-coded priority breakdown
- [x] 3.5 Implement renderDueDateSection() showing overdue/today/week/none counts
- [x] 3.6 Implement renderSummarySection() with completion percentage
- [x] 3.7 Implement View() combining all sections with ModalStyle layout
- [x] 3.8 Handle empty board state gracefully

## 4. App Integration

- [x] 4.1 Add 'S' (shift+s) keybinding in handleNormalMode to toggle ShowStats
- [x] 4.2 Add ShowStats render condition in View() (similar to ShowHelp pattern)
- [x] 4.3 Handle Escape key in stats mode to close overlay
- [x] 4.4 Handle 'S' key in stats mode to toggle off

## 5. Styling

- [x] 5.1 Add StatsHeaderStyle() to styles.go if needed
- [x] 5.2 Add BarGraphStyle() variants for each status/priority color
- [x] 5.3 Ensure overdue count uses WarningStyle or ErrorStyle

## 6. Testing

- [x] 6.1 Add unit tests for ComputeStatistics with various board states
- [x] 6.2 Test empty board statistics
- [x] 6.3 Test overdue detection (tasks with past due dates)
- [x] 6.4 Test completion percentage calculation
- [x] 6.5 Manual testing of statistics overlay rendering
