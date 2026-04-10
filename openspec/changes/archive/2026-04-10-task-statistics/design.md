## Context

The task manager currently shows minimal statistics—just a task count per pane in the status bar. Users managing many tasks need a comprehensive view of their workload distribution to identify bottlenecks (too many tasks in "In Progress"), overdue items, and completion velocity.

The project uses Charm's bubbletea/lipgloss stack. The help overlay (`?` key) demonstrates the full-screen overlay pattern we'll follow.

## Goals / Non-Goals

**Goals:**
- Provide a visually appealing statistics overlay accessible via keyboard shortcut
- Show task distribution across status columns with horizontal bar graphs
- Display priority breakdown with color-coded bars
- Show due date statistics (overdue, due today, due this week, no due date)
- Keep rendering fast (<16ms) for smooth TUI experience

**Non-Goals:**
- Persistent analytics/history tracking (would need storage changes)
- Export/share statistics
- Time-series graphs (completion over days/weeks)
- Per-board comparison statistics
- Animated transitions (harmonica) - keep it simple for v1

## Decisions

### Decision 1: Full-screen overlay pattern

**Choice**: Use the same overlay pattern as help screen (`ShowStats bool` in state, render in View()).

**Rationale**: 
- Consistent with existing codebase patterns
- Reuses ModalStyle for consistent look
- Simple toggle with single key

**Alternatives considered**:
- Side panel: Would compete with preview panel, complex layout
- Popup modal: Too small for graphs

### Decision 2: ASCII bar graphs with lipgloss styling

**Choice**: Render horizontal bar graphs using Unicode block characters (█▓▒░) styled with lipgloss colors.

**Rationale**:
- Pure lipgloss implementation, no new dependencies
- Block characters are widely supported in modern terminals
- Colors can match existing theme (Catppuccin)

**Alternatives considered**:
- bubbles/progress: Designed for single progress bars, not comparisons
- Third-party charting: Adds dependency, may not match styling

### Decision 3: Compute statistics on-demand

**Choice**: Calculate all statistics when rendering the stats view, not cached.

**Rationale**:
- Task lists are small (typically <100 tasks per board)
- Keeps data fresh—no cache invalidation needed
- Simple implementation

**Alternatives considered**:
- Cached stats updated on task changes: Adds complexity, premature optimization

### Decision 4: Shift+S keybinding

**Choice**: `S` (shift+s) toggles statistics overlay.

**Rationale**:
- `s` is taken by sort selector
- `S` suggests "Statistics" or "Summary"
- Uppercase keys typically toggle views/overlays

### Decision 5: Statistics struct in model package

**Choice**: Create `BoardStatistics` struct in `internal/model/statistics.go` with computed metrics.

**Rationale**:
- Separates data computation from rendering
- Easy to unit test
- Can be extended later (e.g., velocity tracking)

## Risks / Trade-offs

**[Bar graph readability at small widths]** → Use minimum 40-char graph width, scale proportionally. Show numeric labels always.

**[Unicode block character support]** → Fallback to ASCII (`#`, `-`) for terminals that don't support Unicode. Could add detection later.

**[Color accessibility]** → Use existing theme colors which have been chosen for contrast. Add labels/values, don't rely only on color.

## Open Questions

- Should statistics show current board only or all boards? → Start with current board only.
