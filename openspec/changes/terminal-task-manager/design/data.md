# Data Storage & Persistence

## JSON files in ~/.tsk/

**Decision**: Store data as JSON files in user's home directory with action-based saves.

**Storage structure**:
```
~/.tsk/
├── data/
│   └── boards/
│       ├── board-abc123.json
│       └── board-def456.json
└── backups/
    └── board-abc123-2026-04-08T14:30:00.json
```

## Data Model

```go
type Task struct {
    ID          string    `json:"id"`           // UUID
    Title       string    `json:"title"`
    Description string    `json:"description,omitempty"`
    Status      Status    `json:"status"`       // todo, in_progress, done
    Priority    Priority  `json:"priority"`     // 0=none, 1=low, 2=med, 3=high
    Labels      []string  `json:"labels,omitempty"` // Free-form tags for categorization
    DueDate     *Time     `json:"due_date,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    Position    int       `json:"position"`     // Order within pane
}

type Board struct {
    ID          string              `json:"id"`
    Name        string              `json:"name"`
    Tasks       []Task              `json:"tasks"`
    BoardLabels map[string]*Label   `json:"board_labels,omitempty"` // Board-level label definitions
    CreatedAt   time.Time           `json:"created_at"`
    UpdatedAt   time.Time           `json:"updated_at"`
}

type Label struct {
    Name  string     `json:"name"`
    Color LabelColor `json:"color"` // red, orange, yellow, green, blue, purple, pink, cyan
}

type LabelColor string // One of: red, orange, yellow, green, blue, purple, pink, cyan
```

## Labels

Labels are stored at the board level for reusability across tasks. Each label has a name and a color, assigned automatically from a predefined palette when created. This ensures consistent coloring when the same label appears on multiple tasks.

**Available colors**: red, orange, yellow, green, blue, purple, pink, cyan

When a new label is added to a task, it is automatically created at the board level with the next available color (cycling through the palette).

**Use cases** - Useful for cross-cutting concerns:

| Category | Example labels |
|----------|----------------|
| Type | `bug`, `feature`, `docs`, `chore` |
| Area | `frontend`, `backend`, `api`, `database` |
| Context | `blocked`, `needs-review`, `urgent` |
| Milestone | `v1.0`, `refactor`, `tech-debt` |

Labels appear on task cards and in the preview panel. Use `/` search or filter menu to find tasks by label.

## Write Strategy: Action-based + Auto-save

Saves occur on meaningful actions (not keystrokes):

| Trigger | Saves |
|---------|-------|
| Exit task detail (`Esc`) | Yes |
| Move task (`>` / `<`) | Yes |
| Create task (`n` → `Enter`) | Yes |
| Delete task (`d` → confirm) | Yes |
| Set priority (`1`/`2`/`3`/`0`) | Yes |
| Label change (`L`) | Yes |
| Quit (`q`, `:wq`) | Yes |
| Switch board (`b`) | Yes |
| Auto-save timer | Every 5 seconds |

Does NOT save on:
- Keystrokes while typing
- Navigation (`j`/`k`/`h`/`l`)
- Opening views (`Enter`, `?`)

**Dirty flag**: Only write to disk if changes exist. Skip write if board unchanged.

## Memory Model

- **Load once**: The current board's JSON is loaded into memory on startup (or board switch)
- **Work in-memory**: All task operations (create, move, edit, delete) modify the in-memory `Board` struct
- **No reload on actions**: Changing task status, priority, etc. does NOT trigger a disk read
- **Write on triggers**: Only save triggers (listed above) write the in-memory state to disk

## Scalability

| Task count | JSON size | Load time | Memory | Notes |
|------------|-----------|-----------|--------|-------|
| 100        | ~15 KB    | <10ms     | ~50 KB | Typical user |
| 500        | ~75 KB    | ~20ms     | ~200 KB| Heavy user |
| 1000       | ~150 KB   | ~50ms     | ~400 KB| Still fast |
| 5000+      | ~750 KB   | ~200ms    | ~2 MB  | Consider archiving done tasks |

For massive boards (5000+ tasks), recommended mitigations:
- Archive completed tasks (`A` to archive all Done tasks, `a` to archive selected)
- Split into multiple boards by project
- Future: lazy-load Done pane if >1000 tasks

## Archive Feature

**Purpose**: Move completed tasks out of the active board to reduce clutter and improve performance.

**Keys** (only available in DONE pane):
- `a`: Archive selected task
- `A`: Archive all Done tasks

**Behavior**:
1. Archived tasks are moved to `~/.tsk/data/archive/<board-id>.json`
2. Archive file is append-only (tasks accumulate over time)
3. Status bar shows "Archived 1 task" or "Archived 15 tasks"
4. Undo available for single-task archive (`a`), not for bulk (`A`)

**Archive file structure**:
```json
{
  "board_id": "abc123",
  "board_name": "Work Tasks",
  "tasks": [
    { "...task...", "archived_at": "2026-04-08T15:00:00Z" }
  ]
}
```

**Key visibility**: `a` and `A` keys only appear in status bar when viewing DONE pane.

## Rationale

- Simple, portable, human-readable format
- No database dependencies to install
- Easy to backup, version control, or manually edit
- Action-based saves reduce unnecessary I/O
- Auto-save provides safety net during long edits
- Dirty flag prevents redundant writes

**Alternatives considered**:
- SQLite: More powerful queries but adds dependency overhead
- Instant writes: More I/O, potential performance impact
- Write on exit: High risk of data loss on crash

---

# Board Lifecycle

## Startup Behavior

- If `~/.tsk/data/boards/` is empty → prompt user for board name, then create and open it
- If exactly one board exists → open it directly
- If multiple boards exist → open the most recently modified board (by `updated_at`)

## First-run Welcome Screen

```
┌────────────────────────────────────┐
│  WELCOME TO TSK                    │
├────────────────────────────────────┤
│  Create your first board:          │
│                                    │
│  Name: Work Tasks█                 │
│                                    │
├────────────────────────────────────┤
│  Enter:create  Ctrl+c:quit         │
└────────────────────────────────────┘
```

## Board Creation (`B` key)

1. Prompt appears: "New board name:"
2. User types name and presses `Enter`
3. Board created with empty panes (To Do, In Progress, Done)
4. Automatically switch to the new board
5. Save triggered immediately

## Board Switching (`b` key)

1. Board selector overlay opens with list of all boards
2. Navigate with `j/k` or `↓/↑`
3. `Enter` to select, `Esc` to cancel
4. Current board saved before switching
5. Selected board loaded

**Board selector mock**:
```
┌────────────────────────────────────┐
│  SELECT BOARD                      │
├────────────────────────────────────┤
│  ▶ Work Tasks          (12 tasks)  │
│    Personal            (5 tasks)   │
│    Side Project        (8 tasks)   │
│    Default             (0 tasks)   │
├────────────────────────────────────┤
│  ↑/↓:navigate  Enter:select  Esc   │
└────────────────────────────────────┘
```

## Rationale

- First-run prompt lets users personalize from the start
- Most-recent board heuristic matches user's likely working context
- Save-before-switch prevents data loss
- Board creation auto-switches so user can start working immediately
