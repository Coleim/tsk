# Features

## Undo/Redo: Command Stack Pattern

**Decision**: Implement undo/redo with a command stack (max 20 actions).

**Design**:
```go
type Command interface {
    Execute() error
    Undo() error
    Description() string  // "delete task", "move task"
}

type UndoManager struct {
    undoStack []Command  // Actions that can be undone
    redoStack []Command  // Actions that can be redone
    maxSize   int        // 20
}
```

**Undoable actions**:
- Create task → Delete the created task
- Delete task → Restore the deleted task
- Move task → Move back to previous pane
- Edit title/description → Restore previous text
- Set priority → Restore previous priority
- Add/remove label → Reverse the label change

**Not undoable**:
- Delete board (too destructive — backup exists)
- Quit

**Behavior**:
- New action clears redo stack (no branching history)
- Status bar shows: `u:undo "delete task" • Ctrl+r:redo`

**Rationale**:
- Command pattern enables clean undo/redo
- 20 actions × ~1KB = ~20KB max memory
- Vim-familiar keybindings (`u`, `Ctrl+r`)

---

## Search and Filter

**Decision**: Implement instant search with optional filters, searching across all panes.

### Search Behavior

1. Press `/` to activate search mode
2. Type query — results update in real-time as you type
3. Search matches against: title, description, and labels (case-insensitive)
4. Results shown in a flat list across all panes, grouped by status
5. `Enter` on a result navigates to that task in its pane
6. `Esc` clears search and returns to normal view

### Search UI Mock

```
┌──────────────────────────────────────────────────────────────────────────────┐
│  tsk  │  Work Tasks                                         Press ? help     │
├──────────────────────────────────────────────────────────────────────────────┤
│  Search: bug fix█                                          (5 matches)       │
├──────────────────────────────────────────────────────────────────────────────┤
│  TO DO                                                                       │
│    ● Fix login bug                                                           │
│    ○ Bug in dashboard                                                        │
│  IN PROGRESS                                                                 │
│   ▶● Fix payment bug◀                                                        │
│  DONE                                                                        │
│    ○ Bug fix for signup                                                      │
│    ○ Fix API bug                                                             │
├──────────────────────────────────────────────────────────────────────────────┤
│  5 matches                                                                   │
│  j/k:navigate  Enter:go to task  Esc:clear search                            │
└──────────────────────────────────────────────────────────────────────────────┘
```

### Filter Capabilities

- **By priority**: Show only High/Med/Low/None priority tasks
- **By label**: Show tasks with specific label(s)
- **By pane**: Limit search to current pane only (toggle)

### Filter Activation

- `f` key opens filter menu during search
- Or press `f` from board view to apply filters without search query

### Filter Menu Mock

```
┌────────────────────────────────────┐
│  FILTERS                           │
├────────────────────────────────────┤
│  Priority:                         │
│    [ ] High  [x] Medium  [ ] Low   │
│                                    │
│  Labels:                           │
│    [x] bug   [ ] feature  [ ] docs │
│                                    │
│  Scope:                            │
│    (•) All panes  ( ) Current pane │
├────────────────────────────────────┤
│  Enter:apply  F:clear  Esc:cancel  │
└────────────────────────────────────┘
```

### Performance Considerations (100+ tasks)

- **In-memory filtering**: All tasks loaded in memory; linear scan is O(n) and fast for <1000 tasks
- **Substring search**: Simple `strings.Contains()` — no indexing needed at this scale
- **Debounced input**: 50ms debounce on keystrokes to avoid re-filtering on every character
- **Virtualized rendering**: Only render visible rows (handled by Bubbletea's viewport)
- **Early termination**: Stop searching after finding 100 matches (show "100+ matches" indicator)

### Search Algorithm

```go
func (b *Board) Search(query string, filters Filters) []Task {
    query = strings.ToLower(query)
    var results []Task
    for _, task := range b.Tasks {
        if filters.MatchesPriority(task) && filters.MatchesLabels(task) {
            if query == "" || matchesQuery(task, query) {
                results = append(results, task)
                if len(results) >= 100 { break }
            }
        }
    }
    return results
}

func matchesQuery(t Task, q string) bool {
    return strings.Contains(strings.ToLower(t.Title), q) ||
           strings.Contains(strings.ToLower(t.Description), q) ||
           containsLabel(t.Labels, q)
}
```

### Rationale

- Instant search matches modern UX expectations (VS Code, Sublime, etc.)
- Cross-pane search answers "where is that task?" quickly
- Filters provide power without complexity — hidden until needed
- Performance optimizations handle realistic workloads (few hundred tasks per board)
- No external search engine needed at single-user scale
