# UI Layout & Navigation

## Layout: Single-Pane with Preview

**Decision**: Display one pane at a time (To Do / In Progress / Done) with a task preview panel and two-line status bar.

**Main layout**:
```
┌──────────────────────────────────────────────────────────────────────────────┐
│  tsk  │  Board Name                                         Press ? help     │
├─────────────────────────────────────────────────┬────────────────────────────┤
│  [TO DO] (3)     IN PROGRESS (2)     DONE (4)   │  PREVIEW                   │
├─────────────────────────────────────────────────┼────────────────────────────┤
│                                                 │                            │
│   ● Task 1                                      │  Task title                │
│  ▶● Task 2 ◀  (selected)                        │  Status: TO DO             │
│   ○ Task 3                                      │  Priority: HIGH            │
│                                                 │  Due: Apr 15               │
│                                                 │  Labels: [test]            │
│                                                 │  Description...            │
├─────────────────────────────────────────────────┴────────────────────────────┤
│  3 tasks in TO DO                                           u:undo  Ctrl+r   │
│  j/k:nav  h/l:pane  n:new  d:del  >/<:move  Enter:edit  1-3:priority  b:board │
└──────────────────────────────────────────────────────────────────────────────┘
```

## Two-line Status Bar

- Line 1: Context info (task count) or action feedback + undo/redo hints
- Line 2: Keyboard shortcuts (always visible) — help hint `?` shown in header

## Help Panel (`?` key)

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                              KEYBOARD SHORTCUTS                              │
├──────────────────────────────────────────────────────────────────────────────┤
│  NAVIGATION                           TASK ACTIONS                           │
│  j/k or ↓/↑    Move between tasks     n          Create new task             │
│  h/l or ←/→    Switch panes           Enter      Edit task                   │
│                                       d          Delete task                 │
│  MOVE TASKS                           > / <      Move task right/left        │
│  >             Move to next pane                                             │
│  <             Move to previous pane  PRIORITY                               │
│                                       1          High priority               │
│  BOARD                                2          Medium priority             │
│  b             Switch board           3          Low priority                │
│  B             Create new board       0          Clear priority              │
│  R             Rename board                                                  │
│  D             Delete board           LABELS                                 │
│                                       L          Manage labels               │
│  SEARCH & FILTER                                                             │
│  /             Search tasks           OTHER                                  │
│  s             Sort by priority       u          Undo                        │
│  F             Clear filters          Ctrl+r     Redo                        │
│                                       q / :wq    Quit                        │
│                                       ?          Toggle this help            │
├──────────────────────────────────────────────────────────────────────────────┤
│                            Press ? or Esc to close                           │
└──────────────────────────────────────────────────────────────────────────────┘
```

## Layout Rationale

- Single pane = full width for task list, less visual noise
- Preview panel shows task details without opening modal
- Tabs show all panes with counts, clear mental model
- `h/l` to switch panes feels like navigating vim buffers
- Two-line status bar keeps shortcuts always visible

---

# Keyboard Navigation

## Vim-style with Fallbacks

**Decision**: Minimal vim keybindings (h/j/k/l) with arrow key alternatives.

## Key Mappings

**Navigation**:
- `j/k` or `↓/↑`: Navigate tasks within current pane
- `h/l` or `←/→`: Switch between panes (To Do ↔ In Progress ↔ Done)

**Task Actions**:
- `Enter`: Open task for full editing
- `n`: New task in current pane
- `d`: Delete (with confirmation)
- `>` / `<`: Move task to next/previous pane
- `Esc`: Close/cancel

**Priority & Labels**:
- `1/2/3/0`: Set priority (High/Med/Low/None)
- `L`: Manage labels

**Global**:
- `/`: Search/filter
- `s`: Sort by priority
- `u`: Undo last action
- `Ctrl+r`: Redo
- `q` or `:wq`: Quit app
- `?`: Help overlay

**DONE pane only**:
- `a`: Archive selected task
- `A`: Archive all Done tasks

## Context-sensitive Status Bar

The status bar shows different shortcuts based on current pane:

**In TO DO or IN PROGRESS pane**:
```
j/k:nav  h/l:pane  n:new  d:del  >/<:move  Enter:edit  1-3:priority  b:board
```

**In DONE pane**:
```
j/k:nav  h/l:pane  d:del  a:archive  A:archive all  Enter:edit  b:board
```

## Rationale

- 17 keys to learn, covers all essential use cases
- Vim users get h/j/k/l, everyone else uses arrows
- `Esc` for safe close, `q` only quits from main view
- `u` and `Ctrl+r` match vim undo/redo

## Input Modes

The application uses distinct modes to separate commands from text input:

| Mode | Active when | Keys are... | Exit with |
|------|-------------|-------------|-----------|
| **Normal** | Board view, navigating tasks | Commands (`j`, `k`, `n`, `d`, etc.) | — |
| **Insert** | Creating/editing task title/description | Text input | `Enter` (save) or `Esc` (cancel) |
| **Search** | Search input active (`/`) | Text input | `Enter` (go to result) or `Esc` (cancel) |
| **Modal** | Dialog open (delete confirm, filters) | Limited commands | `Enter` (confirm) or `Esc` (cancel) |

**Visual indicators**:
- **Normal mode**: No cursor in text field, task list is navigable
- **Insert mode**: Blinking cursor in text input, status bar shows "-- INSERT --"
- **Search mode**: Cursor in search box, status bar shows search shortcuts

**Mode transitions**:
```
Normal ──n──> Insert (new task)
Normal ──Enter──> Insert (edit task)
Normal ──/──> Search
Normal ──d──> Modal (delete confirmation)
Insert ──Esc──> Normal (discard)
Insert ──Enter──> Normal (save)
Search ──Esc──> Normal
Modal ──Esc──> Normal
```

**Key behavior by mode**:
| Key | Normal mode | Insert mode | Search mode |
|-----|-------------|-------------|-------------|
| `j` | Move down | Type "j" | Type "j" |
| `k` | Move up | Type "k" | Type "k" |
| `Esc` | (no effect) | Cancel edit | Cancel search |
| `Enter` | Edit task | Save text | Go to result |
| `n` | New task | Type "n" | Type "n" |

**Implementation**: Bubbletea's `Update` function checks current mode before interpreting key events:
```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if m.mode == ModeInsert {
            return m.handleInsertMode(msg)
        }
        return m.handleNormalMode(msg)
    }
    return m, nil
}
```
