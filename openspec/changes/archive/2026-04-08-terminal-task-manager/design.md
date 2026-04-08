# Design Overview

## Context

This is a new terminal-based task management application. No existing codebase to consider. The target users are developers and command-line power users who want to manage tasks without leaving their terminal. The application should feel native to terminal workflows with fast keyboard navigation.

## Goals

- Provide a responsive, keyboard-first terminal UI for task management
- Support Kanban-style workflow with fixed panes (To Do, In Progress, Done)
- Enable offline-first operation with local data persistence
- Deliver vim-style navigation for efficient keyboard users
- Keep the application lightweight and fast to start

## Non-Goals

- Cloud sync or multi-device support (out of scope for initial version)
- Team collaboration features (single-user focus initially)
- Integration with external services like Jira, Trello, GitHub (future enhancement)
- Mobile or web interfaces

## Design Documents

| Document | Topics |
|----------|--------|
| [architecture.md](design/architecture.md) | Technology stack (Go + Bubbletea), Elm architecture, project structure |
| [data.md](design/data.md) | JSON storage, data model, write strategy, memory model, scalability, board lifecycle |
| [ui.md](design/ui.md) | Layout, status bar, help panel, keyboard navigation |
| [features.md](design/features.md) | Undo/redo, search and filter |

## Features Summary

### Task Management
- Create, edit, delete tasks
- Move tasks between panes (`>` / `<`)
- Set priority (High/Med/Low/None) with `1/2/3/0`
- Add labels for categorization (`L`)
- Set due dates (`t`)

### Board Management
- Multiple boards with custom names
- Create board (`B`), switch board (`b`)
- Rename (`R`) and delete (`D`) boards
- Auto-select most recent board on startup
- First-run welcome screen prompts for board name

### Navigation
- Vim-style keys (`h/j/k/l`) + arrow key fallbacks
- Single-pane view with preview panel
- Pane tabs showing To Do / In Progress / Done with counts

### Search & Filter
- Instant search (`/`) across title, description, labels
- Filter by priority and labels (`f`)
- Cross-pane search with results grouped by status
- 50ms debounce, max 100 results

### Undo/Redo
- Command stack pattern (max 20 actions)
- `u` to undo, `Ctrl+r` to redo
- Undoable: create, delete, move, edit, priority, labels

### Archive (DONE pane only)
- Archive selected task (`a`)
- Archive all Done tasks (`A`)
- Stored in `~/.tsk/data/archive/<board-id>.json`

### Persistence
- JSON files in `~/.tsk/data/boards/`
- Action-based saves (not on every keystroke)
- 5-second auto-save with dirty flag
- Automatic backups before destructive operations

### UI Polish
- Two-line status bar (feedback + shortcuts)
- Context-sensitive shortcuts (different in DONE pane)
- Help overlay (`?`)
- Data export (`E`) and import (`I`)

## UX Refinements

### Task Movement Follows Focus

When a task is moved to another pane using `>` or `<`, the view automatically switches to the target pane and selects the moved task. This provides immediate feedback and allows the user to continue working with the task they just moved.

**Behavior**:
- Moving a task with `>` switches to the next pane and highlights the moved task
- Moving a task with `<` switches to the previous pane and highlights the moved task
- The moved task maintains its position in the new pane's list

### Status-Specific Tab Colors

Each pane tab uses a distinct color to improve visual identification:
- **To Do**: Light blue (#75)
- **In Progress**: Yellow/amber (#220)
- **Done**: Green (#48)

Active tabs display in their status color; inactive tabs are gray.

### Labels in Task Edit Modal

The task edit modal (`Enter` on a task) includes a labels field alongside title and description. Labels can be entered manually as comma-separated values, or selected from a popup.

**Field order**:
1. Title (required)
2. Description (optional)
3. Labels (optional)

Navigate between fields with `Tab` / `Shift+Tab`.

**Label popup behavior**:
- When focus is on the Labels field, pressing `Tab` opens a small popup listing all board labels
- Navigate through labels with `Tab`
- Select highlighted label with `Enter` (appends to current labels)
- Close popup with `Escape`

**Label storage**:
- Labels are stored at the board level in `BoardLabels` map
- Each label has a name and auto-assigned color
- Labels persist with board data and can be versioned

## Risks / Trade-offs

**[Risk]** Terminal compatibility across different emulators → Test on common terminals (iTerm2, Terminal.app, Windows Terminal, xterm). Use well-supported ANSI escape codes only.

**[Risk]** Large task lists may impact rendering performance → Implement virtualized list rendering for boards with 100+ tasks.

**[Risk]** Data loss from file corruption → Implement atomic writes with temp files and automatic backups before modifications.

**[Trade-off]** JSON storage limits query complexity → Acceptable for single-user scale; can migrate to SQLite if needed later.

**[Trade-off]** No cloud sync limits cross-device usage → Documented as non-goal; users can manually sync files via git or cloud storage.

