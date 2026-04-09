## Why

The TUI has inconsistent popup/modal presentation. Due Date and Labels editors render full-screen while Insert Task and Search use centered popup overlays. Additionally, the label selection styling has visual artifacts that detract from the polished look.

## What Changes

- **Due Date Editor**: Convert from full-screen to centered popup overlay using `overlayDialog()`
- **Labels Editor**: Convert from full-screen to centered popup overlay using `overlayDialog()`
- **Labels Editor selection**: Remove bordered style for selected items, use arrow indicator only
- **Edit Task label popup**: Reduce background brightness on selected items (too bright/harsh)
- **Edit Task label popup**: Fix gray background extending past unselected label text

## Capabilities

### New Capabilities
<!-- None - this is UI polish, not new functionality -->

### Modified Capabilities
<!-- These are implementation styling fixes, not spec-level behavior changes -->

## Impact

- `internal/ui/duedate.go`: Modify View() to return smaller content for popup compositing
- `internal/ui/labels.go`: Modify View() to return smaller content; fix selected item styling
- `internal/ui/app.go`: Wrap DueDateEditor and LabelEditor views with overlayDialog()
- `internal/styles/styles.go`: Adjust PopupSelectedItemStyle and PopupItemStyle colors
- `internal/ui/edit.go`: May need adjustments to label popup rendering
