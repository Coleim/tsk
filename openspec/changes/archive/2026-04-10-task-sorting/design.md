## Context

The task manager currently displays tasks in the order they were added to each status pane. Users with many tasks need ways to quickly find and prioritize tasks based on different criteria. The filtering system exists but sorting would complement it by organizing the filtered (or unfiltered) results.

Current state:
- Tasks are stored in arrays per status (To Do, In Progress, Done)
- 's' key currently sorts by priority only (high to low) - **modifies storage order**
- Filter system ('f') modifies the task list view with multi-select priority/label filtering - **display only**
- DueDate and CreatedAt fields exist on tasks

Existing keyboard commands reference:
- `f` - Open filter popup (priority + labels)
- `F` - Clear filters
- `s` - Sort by priority (currently calls `Board.SortByPriority()`, modifies storage)
- `/` - Search tasks
- `?` - Help overlay

## Goals / Non-Goals

**Goals:**
- Allow users to sort tasks by creation date, due date, title, and priority
- Provide ascending/descending options for each sort type
- Show clear indication of current sort mode
- Make sorting accessible via keyboard shortcut
- Persist sort preference per board

**Non-Goals:**
- Multi-field sorting (e.g., sort by priority then by date)
- Custom sort orders or drag-and-drop reordering
- Sort persistence across sessions (can be added later)

## Decisions

### Decision 1: Sort at display time, not storage

**Choice**: Apply sorting in `CurrentTasks()` after filtering, don't modify stored task order.

**Rationale**: 
- Maintains data integrity - original order preserved
- Consistent with filter behavior (filters don't modify storage)
- Allows different views without data conflicts

**Alternatives considered**:
- Sort in storage: Would affect all views, harder to implement undo

### Decision 2: Replace 's' single-action sort with sort selector

**Choice**: Press 's' to open sort mode selector popup (replacing current single-action priority sort).

**Rationale**:
- 's' is already mnemonic for "sort" and used for sorting
- Consistent with 'f' for filter pattern (popup selector)
- More powerful than single-action while using same key

**Alternatives considered**:
- New key (e.g., 'S' shift-s): Inconsistent, 's' already means sort
- Keep 's' as priority sort, add 'S' for selector: Fragmented UX

### Decision 3: Sort mode as enum with compound options

**Choice**: Define `SortMode` enum with combined field+direction values (e.g., `SortCreatedDesc`, `SortTitleAsc`).

**Rationale**:
- Simple to store and compare
- No need for separate direction tracking
- Clear intent in code

### Decision 4: Default sort is by creation date (newest first)

**Choice**: New boards default to `SortCreatedDesc`.

**Rationale**:
- Shows recent work first
- Matches user mental model of "what did I add recently"

## Risks / Trade-offs

**[Performance with large lists]** → Sorting happens on every view render. Mitigate by only re-sorting when sort mode changes or tasks change.

**[Sort stability]** → Go's sort.Slice is not stable. Use sort.SliceStable when secondary ordering matters (e.g., preserve insertion order for equal values).

**[Nil due dates]** → Tasks without due dates need clear handling. Decision: sort nil due dates last when sorting by due date ascending, first when descending.
