## Why

The current full-screen modals for simple actions like creating a task or confirming deletion feel heavy and disorienting. Users lose context of their board when these modals appear. The recently implemented search popup demonstrates that centered overlays with the board visible behind provide a better user experience for quick interactions.

## What Changes

- Task creation will use a compact centered popup overlay instead of full-screen modal
- Delete task confirmation will use a compact centered popup overlay instead of full-screen modal
- Both popups will show the board visible behind them (using lipgloss v2 Layer compositing)
- The interactions (Enter to confirm, Escape to cancel) remain the same

## Capabilities

### New Capabilities

None - this modifies existing capabilities.

### Modified Capabilities

- `task-operations`: Change task creation modal to a compact popup overlay
- `tui-core`: Update dialog requirement to distinguish between full-screen dialogs (for complex editing) and popup overlays (for simple inputs/confirmations)

## Impact

- `internal/ui/app.go` - Update task creation mode rendering to use popup overlay
- `internal/ui/modal.go` - Update confirmation modal to use popup overlay pattern
- Reuses the `overlayDialog()` function already implemented for search popup
