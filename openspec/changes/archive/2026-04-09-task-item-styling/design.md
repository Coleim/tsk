## Context

The current task list displays tasks as plain text lines with minimal visual separation. Each task line consists of:
- Selection indicator (` ▶` when selected or `  ` when not)
- Priority indicator (colored symbol)
- Task title

The selection arrow sits directly adjacent to the priority icon creating a cramped appearance. Tasks lack individual visual boundaries, blending together in the list.

## Goals / Non-Goals

**Goals:**
- Improve visual separation between task items
- Add proper spacing between selection indicator and priority icon
- Make the selected task more visually prominent
- Maintain terminal compatibility and performance

**Non-Goals:**
- Changing task colors or themes
- Adding animations to task items
- Redesigning the overall layout

## Decisions

### 1. Use background highlighting for selected task
Instead of adding borders around each task (which consumes vertical space), use a subtle background color for the selected task row. This provides visual distinction without reducing visible task count.

**Alternatives considered:**
- Box borders per task: Rejected - takes too much vertical space
- Full task row boxes: Rejected - too visually heavy

### 2. Add spacing between selection indicator and priority
Change from ` ▶P ` to ` ▶  P ` (add one space after arrow before priority). This creates breathing room between the cursor and content.

### 3. Keep tasks as single-line items
Maintain the compact single-line format. Adding more height per task would reduce visible tasks in the viewport.

## Risks / Trade-offs

- Background colors may not render well in all terminals → Use theme-aware colors that degrade gracefully
- Spacing changes reduce available title width → Accept slightly shorter titles (2 chars less)
