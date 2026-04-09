## Context

Currently search is a full-screen overlay that hides the board. The Search component (`internal/ui/search.go`) renders results grouped by status with contains-matching across title, description, and labels.

## Goals / Non-Goals

**Goals:**
- Compact centered popup that keeps board visible
- "Starts with" matching on task titles for faster lookup
- Max 10 visible results to keep popup compact

**Non-Goals:**
- Changing search across description/labels (keep contains for those)
- Adding fuzzy matching or ranking

## Decisions

### Decision: Popup positioning
Use `lipgloss.Place()` to center the search popup. Width fixed at 60-70 chars.

**Rationale**: Consistent with other popups in the app. Large enough for task titles, small enough to not obscure the board.

### Decision: Starts-with matching for titles
Check `strings.HasPrefix(title, query)` for title matches, keep `strings.Contains` for description/labels.

**Rationale**: User explicitly asked for starts-with. Title is the primary search target; description/labels are secondary.

## Risks / Trade-offs

- **Risk**: Starts-with is more restrictive → users may not find tasks
  - *Mitigation*: Description and labels still use contains-matching
