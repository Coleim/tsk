## Why

Task items in the list appear cramped and lack visual distinction. The selection arrow `▶` sits directly against the priority icon with no spacing, making the UI feel cluttered. Individual tasks blend together without clear visual boundaries.

## What Changes

- Add visual containers (subtle borders or backgrounds) around individual task items for better separation
- Increase spacing between selection indicator and priority icon
- Apply consistent padding within task items for a more polished look

## Capabilities

### New Capabilities
<!-- None - this is a styling enhancement -->

### Modified Capabilities
- `tui-core`: Updating task item visual presentation in the task list

## Impact

- `internal/ui/app.go`: Task rendering in `renderTaskList()`
- `internal/styles/styles.go`: New or modified styles for task items
