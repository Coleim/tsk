## Why

The current layout uses fixed percentages (30/70) for task list and preview panels, which can cause readability issues at small window sizes. The header and status bars should remain fixed regardless of window size, with only the content area adapting.

## What Changes

- **Fixed header/status bars**: Header always at top, status bar always at bottom, regardless of window height
- **Minimum task list width**: Task list maintains minimum readable width (e.g., 25-30 chars)
- **Flexible preview**: Preview panel takes remaining width after task list minimum is satisfied
- **Small window handling**: When window is too small for both panels, switch to single-panel mode or show graceful degradation message

## Capabilities

### New Capabilities
<!-- None - this is layout refinement, not new functionality -->

### Modified Capabilities
<!-- These are implementation layout fixes, not spec-level behavior changes -->

## Impact

- `internal/ui/app.go`: Modify `renderMainView()` width calculations, ensure header/status always render
- Possibly add minimum window size handling or single-panel fallback mode
