## Why

Users need visibility into their task management patterns to understand productivity trends and identify bottlenecks. Currently, task counts are only shown per-pane in the status bar. A dedicated statistics view with visual graphs would help users see completion rates, priority distribution, and overdue task counts at a glance—making the TUI both more informative and delightful.

## What Changes

- Add a statistics overlay accessible via a keyboard shortcut
- Display task counts by status as horizontal bar graphs
- Show priority distribution as a visual breakdown
- Display overdue and upcoming due date statistics
- Include completion velocity metrics (tasks done today/this week)
- Use Charm's lipgloss for beautiful, colorful rendering with ASCII/Unicode bar characters

## Capabilities

### New Capabilities
- `task-statistics`: Visual statistics overlay showing task metrics, graphs, and completion trends

### Modified Capabilities
- `keyboard-navigation`: Add 'S' (shift+s) keybinding to toggle statistics overlay

## Impact

- **Code**: New `internal/ui/stats.go` component, new `ModeStats` in mode.go
- **Dependencies**: Uses existing lipgloss (already in project), may add bubbles progress component
- **UI**: New full-screen overlay view, similar pattern to help overlay
