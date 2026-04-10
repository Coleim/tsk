# tsk - Terminal Task Manager

> **Note:** This project was built to test and demonstrate the [OpenSpec](https://github.com/openspec-dev/openspec) framework — an experimental approach to AI-assisted software development using structured specifications. The app itself is fully functional!

A fast, keyboard-driven terminal task manager with Kanban-style workflow. Built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)

## Overview

**tsk** is a terminal-based task manager that brings Kanban-style workflow to your command line. No browser, no electron, no bloat — just a fast TUI that stays out of your way.

![alt text](/docs/assets/screenshot_01.png)
![alt text](/docs/assets/screenshot_02.png)

## Features

- **Kanban Workflow**: Three panes (To Do, In Progress, Done) for visual task tracking
- **Vim-Style Navigation**: Move between tasks with `j`/`k` and panes with `h`/`l`
- **Multiple Boards**: Create and switch between separate task boards
- **Priority Levels**: Set High/Medium/Low/None priority with quick keys
- **Labels**: Categorize tasks with custom labels
- **Due Dates**: Track deadlines with due date support
- **Search & Filter**: Find tasks instantly across all panes
- **Undo/Redo**: Recover from mistakes with full undo support
- **Archive**: Clean up completed tasks while preserving history
- **Export/Import**: Share boards as JSON files
- **Offline-First**: All data stored locally in `~/.tsk/`

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap Coleim/tsk https://github.com/Coleim/tsk.git
brew install tsk
```

### Shell Script

```bash
curl -sSL https://raw.githubusercontent.com/Coleim/tsk/main/install.sh | bash
```

Or install a specific version:

```bash
curl -sSL https://raw.githubusercontent.com/Coleim/tsk/main/install.sh | bash -s -- v1.0.0
```

### Go Install

```bash
go install github.com/coliva/tsk/cmd/tsk@latest
```

### From Source

```bash
git clone https://github.com/Coleim/tsk.git
cd tsk
make build
sudo mv bin/tsk /usr/local/bin/
```

## Quick Start

```bash
# Start tsk
tsk

# First run will prompt for a board name
# Then you're ready to add tasks!
```

## Keyboard Shortcuts

### Navigation

| Key | Action |
|-----|--------|
| `j` / `↓` | Move down in task list |
| `k` / `↑` | Move up in task list |
| `h` / `←` | Switch to previous pane |
| `l` / `→` | Switch to next pane |

### Task Actions

| Key | Action |
|-----|--------|
| `n` | Create new task |
| `Enter` | Edit task |
| `d` | Delete task |
| `>` | Move task to next pane |
| `<` | Move task to previous pane |
| `1` | Set HIGH priority |
| `2` | Set MEDIUM priority |
| `3` | Set LOW priority |
| `0` | Clear priority |
| `L` | Manage labels |
| `t` | Set due date |

### Board Management

| Key | Action |
|-----|--------|
| `b` | Switch boards |
| `B` | Create new board |
| `R` | Rename current board |
| `D` | Delete current board |

### Search & Filter

| Key | Action |
|-----|--------|
| `/` | Search tasks |
| `f` | Open filter panel |
| `F` | Clear all filters |
| `s` | Sort by priority |

### Archive & Export (in Done pane)

| Key | Action |
|-----|--------|
| `a` | Archive selected task |
| `A` | Archive all done tasks |
| `E` | Export current board |

### General

| Key | Action |
|-----|--------|
| `u` | Undo last action |
| `Ctrl+r` | Redo |
| `?` | Show help |
| `q` | Quit |

## Data Storage

All data is stored in `~/.tsk/`:

```
~/.tsk/
├── data/
│   ├── boards/      # Board JSON files
│   └── archive/     # Archived tasks
└── backups/         # Automatic backups
```

### Export/Import

Export your current board:
```bash
# In the app, press 'E' to export
# Creates: tsk-export-BoardName.json
```

Import a board:
```bash
tsk import tsk-export-BoardName.json
```

## Command Line

```bash
tsk                      # Start the TUI
tsk import <file>        # Import a board from JSON
tsk help                 # Show help
tsk version              # Show version
```

## Configuration

tsk works out of the box with sensible defaults. Data is auto-saved:
- After each action (create, edit, delete, move)
- Every 5 seconds if there are unsaved changes
- On quit or board switch

## Themes

tsk supports dark and light themes. Set the `TSK_THEME` environment variable to switch:

```bash
# Use dark theme (default - Catppuccin-inspired)
tsk

# Use light theme
TSK_THEME=light tsk

# Or export for all sessions
export TSK_THEME=light
```

### Dark Theme (Default)
- Background: Deep purple-blue (#1e1e2e)
- Accent: Lavender (#cba6f7)
- Optimized for low-light environments

### Light Theme
- Background: Light gray (#eff1f5)
- Accent: Purple (#8839ef)
- Optimized for well-lit environments

## Development

### Prerequisites

- Go 1.21+
- golangci-lint (`brew install golangci-lint`)

### Building

```bash
make build    # Build binary to bin/tsk
make run      # Build and run
make test     # Run tests
make lint     # Run linter
make check    # Run lint + tests
```

### Pre-commit Hooks

Install Git hooks to run lint and tests before each commit:

```bash
./scripts/install-hooks.sh
```

This ensures code quality checks pass locally before pushing to CI.

### Makefile Targets

| Target | Description |
|--------|-------------|
| `make build` | Build the binary |
| `make test` | Run unit tests |
| `make lint` | Run golangci-lint |
| `make check` | Run lint + tests (pre-commit) |
| `make bench` | Run benchmarks |
| `make perf-test` | Run performance threshold tests |

## License

MIT License - see [LICENSE](LICENSE) for details.

---

<sub>Built with ❤️ using [OpenSpec](https://github.com/openspec-dev/openspec) — an experimental AI-assisted development framework.</sub>
