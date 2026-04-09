## Context

The `renderMainView()` in app.go calculates layout as:
- Header: 1 line (fixed)
- Tabs: 1 line (fixed)
- Status bar: 2 lines (fixed)
- Content area: remaining height, split 30/70 between task list and preview

Current code:
```go
leftWidth := a.state.Width * 30 / 100
rightWidth := a.state.Width - leftWidth
```

This percentage-based split causes task list to become unreadably narrow at small widths (e.g., 80-char terminal → 24-char task list).

## Goals / Non-Goals

**Goals:**
- Header and status bar always visible at their fixed positions
- Task list has minimum readable width (30 chars minimum)
- Preview takes remaining width flexibly
- Hide preview completely when window too narrow (single-panel mode)
- Graceful handling when window too small

**Non-Goals:**
- Collapsible panels or resizable splits
- Saving/restoring layout preferences
- Different layouts per mode

## Decisions

### 1. Width Thresholds and Modes

| Width | Mode | Task List | Preview |
|-------|------|-----------|---------|
| ≥80 | Full | 30% (min 30) | Remaining |
| 50-79 | Compact | 30 chars fixed | Remaining |
| <50 | Single Panel | Full width | **Hidden** |

### 2. Minimum Width Constants
Define minimum widths:
- Task list minimum: 30 characters
- Preview minimum: 20 characters  
- Single-panel threshold: 50 characters

**Rationale**: 30 chars fits typical task titles. Below 50 total, hide preview entirely.

### 3. Width Calculation Logic
```go
const (
    MinTaskListWidth = 30
    MinPreviewWidth  = 20
    SinglePanelThreshold = 50
)

if width >= SinglePanelThreshold {
    taskListWidth = max(MinTaskListWidth, width * 30 / 100)
    previewWidth = width - taskListWidth
} else {
    // Single panel mode - task list only
    taskListWidth = width
    previewWidth = 0  // Preview hidden
}
```

**Rationale**: Ensures task list never shrinks below readable width. Preview hidden entirely below threshold for clean UX.

### 4. Single-Panel Mode Behavior
When width < 50:
- Show only task list at full width
- Preview panel not rendered at all
- User can press Enter to view task details in detail mode

**Rationale**: Clean single-panel experience for vertical/narrow terminals. Core navigation preserved.

## Risks / Trade-offs

**[Preview disappears at small sizes]** → Acceptable. Task list is primary UI. Enter key opens detail view.

**[Magic numbers]** → Define as constants at top of file for easy tuning.

**[Abrupt transition]** → Preview disappears at exactly 50 chars. Could add intermediate "minimal preview" but keeping it simple.
