## Why

Command-line users and developers need a fast, keyboard-driven task management tool that works entirely in the terminal. Existing tools like Jira and Trello require context-switching to a web browser, breaking developer flow and adding friction to task tracking workflows.

## What Changes

- Add a terminal-based user interface (TUI) for managing tasks, boards, and workflows
- Support Kanban-style board views with panes (To Do, In Progress, Done)
- Enable task creation, editing, deletion, and status transitions via keyboard shortcuts
- Provide filtering and search capabilities for tasks
- Support data persistence to local storage
- Include color-coded task priorities and labels
- Support vim motions

## Capabilities

### New Capabilities
- `tui-core`: Core terminal UI framework with panels, navigation, and keyboard handling
- `board-management`: Create, view, and manage Kanban boards with fixed panes
- `task-operations`: CRUD operations for tasks including title, description, priority, labels, and due dates
- `task-views`: Different view modes - board view, list view, and task detail view
- `persistence`: Local file-based storage for boards and tasks (JSON format)
- `keyboard-navigation`: Vim-style and arrow-key navigation with customizable shortcuts

### Modified Capabilities
<!-- No existing capabilities to modify - this is a new project -->

## Impact

- **New codebase**: Standalone terminal application
- **Dependencies**: Go + Bubbletea/Bubbles/Lipgloss (Charm ecosystem)
- **Storage**: Local JSON files in user's home directory (~/.tsk/)
- **Platform**: Cross-platform support (macOS, Linux, Windows) via single binary
