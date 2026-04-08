# Manual Testing Checklist for tsk

Run through this checklist on each target terminal to verify compatibility.

## Pre-requisites

```bash
# Build the binary
make build
# Or run directly
go run ./cmd/tsk
```

## Terminal Compatibility Tests (14.22-14.25)

### For each terminal (macOS Terminal.app, iTerm2, Windows Terminal, Linux xterm/gnome-terminal):

#### Display & Rendering
- [ ] App launches full-screen without artifacts
- [ ] Three-column layout renders correctly
- [ ] Colors display correctly (blue/yellow/green for statuses)
- [ ] Priority indicators (●/◐/○) render correctly
- [ ] Box drawing characters (borders) render correctly
- [ ] Unicode labels display correctly

#### Navigation
- [ ] `j`/`k` moves task selection up/down
- [ ] `h`/`l` switches between panes
- [ ] Arrow keys work identically
- [ ] `g`/`G` jumps to first/last task
- [ ] Tab cycles through panes

#### Task Operations
- [ ] `n` opens new task input
- [ ] `e` opens edit mode
- [ ] `d` prompts for delete confirmation
- [ ] `>` moves task to next status
- [ ] `<` moves task to previous status
- [ ] `1`/`2`/`3`/`0` sets priority

#### Search & Filter
- [ ] `/` opens search mode
- [ ] Search results update as you type
- [ ] `Enter` navigates to selected result
- [ ] `f` opens filter panel
- [ ] `F` clears all filters

#### Board Management
- [ ] `b` opens board selector
- [ ] `B` creates new board
- [ ] Board switching preserves state

#### Other
- [ ] `?` toggles help overlay
- [ ] `u` undoes last action
- [ ] `Ctrl+r` redoes
- [ ] `q` or `:wq` quits

### Terminal Resize Test (14.26)

1. Start tsk at normal size
2. Resize terminal window smaller
   - [ ] UI adapts without crashing
   - [ ] Content remains visible
3. Resize terminal window larger
   - [ ] UI expands to fill space
   - [ ] No artifacts or missing borders
4. Quickly resize multiple times
   - [ ] App remains responsive
   - [ ] No flickering or corruption

## Performance Test (14.27)

```bash
# Generate 500+ tasks
./scripts/perf-test.sh

# Then run tsk and select "Performance Test" board
go run ./cmd/tsk
```

### Performance Checklist
- [ ] App loads board within 2 seconds
- [ ] Navigation (j/k) responds in <100ms
- [ ] Pane switching (h/l) is smooth
- [ ] Search responds within 500ms
- [ ] Scrolling through 200 tasks is smooth
- [ ] No visible lag when selecting tasks
- [ ] Memory usage stays under 100MB

### Cleanup
```bash
rm ~/.tsk/data/boards/board-perf-test.json
```

## Sign-off

| Terminal | Tester | Date | Pass/Fail | Notes |
|----------|--------|------|-----------|-------|
| macOS Terminal.app | | | | |
| iTerm2 | | | | |
| Windows Terminal | | | | |
| Linux gnome-terminal | | | | |
| Linux xterm | | | | |
