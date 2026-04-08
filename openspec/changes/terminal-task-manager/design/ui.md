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

**Status bar positioning**: The status bar MUST always be pinned to the bottom of the terminal, regardless of content height. Use calculated padding to push the status bar down when content is short.

## Responsive Layout

The UI adapts to terminal size changes:

**Minimum dimensions**: 60 columns × 10 rows (show warning below this)

**Responsive behavior**:
| Terminal width | Adaptation |
|----------------|------------|
| < 80 cols | Hide preview panel, full-width task list |
| 80-120 cols | 30/70 split (task list / preview) |
| > 120 cols | 30/70 split (task list / preview) |

| Terminal height | Adaptation |
|-----------------|------------|
| < 15 rows | Hide tabs, compact mode |
| 15-30 rows | Standard layout |
| > 30 rows | Show more tasks, no change to layout |

**On resize** (`tea.WindowSizeMsg`):
1. Recalculate panel widths
2. Re-render immediately
3. Clamp selection if visible rows decrease
4. Status bar stays at bottom

**Implementation**:
```go
func (a *App) View() string {
    // Calculate fixed heights
    headerHeight := 1
    tabsHeight := 1
    statusHeight := 2
    contentHeight := a.state.Height - headerHeight - tabsHeight - statusHeight - 2

    // Render content with padding to push status bar down
    contentLines := countLines(content)
    padding := strings.Repeat("\n", max(0, contentHeight - contentLines))

    return lipgloss.JoinVertical(lipgloss.Left,
        header,
        tabs,
        content + padding,
        statusBar,  // Always at bottom
    )
}
```

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
- Task list and preview panel have equal height for visual balance

## Loading State

**Decision**: Show an animated spinner while loading board data.

```
┌─────────────────────────┐
│                         │
│   ⣾ Loading tasks...    │
│                         │
└─────────────────────────┘
```

- Uses `bubbles/spinner` component with Dot style
- Centered both horizontally and vertically
- Displayed during initial board load
- Transitions to main UI once data is ready

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

---

# Color Scheme

## Status-Specific Colors

Each pane/status uses a distinct color for visual identification:

| Status | Color | ANSI Code | Usage |
|--------|-------|-----------|-------|
| To Do | Light blue | #75 | Tab text, status indicators |
| In Progress | Yellow/amber | #220 | Tab text, status indicators |
| Done | Green | #48 | Tab text, status indicators |

**Tab rendering**:
- Active tab: Bold, status color, bracketed `[TO DO (3)]`
- Inactive tab: Gray (#245), no brackets `IN PROGRESS (2)`

**Priority colors** (unchanged):
- High: Red (#196)
- Medium: Orange (#214)
- Low: Green (#48)
- None: Gray (#245)

**UI accent colors**:
- Primary: Bright blue (#39)
- Accent: Magenta/pink (#213) - used for titles, modal borders
- Success: Green (#48)
- Warning: Yellow (#220)
- Error: Red (#196)

---

# Task Edit Modal

## Fields

The task edit modal includes three editable fields:

| Field | Placeholder | Char Limit | Required |
|-------|-------------|------------|----------|
| Title | "Task title..." | 256 | Yes |
| Description | "Description (optional)..." | 1024 | No |
| Labels | "Labels (comma-separated)..." | 512 | No |

**Labels format**: Comma-separated values. Whitespace is trimmed from each label.

Example: `bug, urgent, frontend` → `["bug", "urgent", "frontend"]`

**Navigation**:
- `Tab`: Move to next field
- `Shift+Tab`: Move to previous field
- `Enter`: Save changes
- `Esc`: Cancel without saving

**Visual indicator**: Current field has `▶` prefix on label.

---

# Task Movement Behavior

## Follow Mode

When moving a task between panes with `>` or `<`, the view automatically follows the task:

**Behavior**:
1. Task status is updated to new pane
2. View switches to the target pane
3. Selection moves to the moved task
4. Status message confirms: "Moved to In Progress"

**Implementation**: `SwitchToTaskInPane(status, taskID)` method updates `CurrentPane` and finds the task's new index.

**Rationale**: Provides immediate feedback and allows continued work on the same task without manual navigation.
