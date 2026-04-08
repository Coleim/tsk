# Architecture

## Technology Stack: Go + Bubbletea

**Decision**: Build with Go and the Bubbletea framework (Elm architecture for terminal UIs).

**Rationale**:
- Bubbletea provides a clean, composable Elm-style architecture
- Go compiles to a single binary with no runtime dependencies
- Excellent performance and low memory footprint
- Strong standard library for file I/O and JSON handling
- Lipgloss companion library for styling
- Bubbles library provides ready-made components (text input, lists, tables)

**Alternatives considered**:
- Rust + Ratatui: Similar benefits but Go has gentler learning curve
- TypeScript + Ink: Requires Node.js runtime, slower startup
- Python + Textual: Good option but interpreted language has performance overhead

## Elm Architecture (Model-Update-View)

**Decision**: Use Bubbletea's Elm architecture with composable models.

**Project structure**:
```
cmd/
└── tsk/
    └── main.go       # Entry point
internal/
├── model/            # Data types (Board, Task, Status)
├── ui/               # Bubbletea components
│   ├── board/        # Board view model
│   ├── task/         # Task views and modals
│   ├── list/         # List view model
│   └── common/       # Shared components (header, status bar)
├── storage/          # Persistence layer
└── styles/           # Lipgloss styles and themes
```

**Rationale**:
- Elm architecture provides predictable state management
- Composable models allow complex UIs from simple pieces
- Idiomatic Go project structure with cmd/ and internal/
- Easy to test each component's Update function

## Testing Strategy

**Approach**: Bottom-up testing with emphasis on unit tests for business logic.

### Test Structure
```
internal/
├── model/
│   ├── task.go
│   └── task_test.go      # Unit tests for Task, Board structs
├── storage/
│   ├── storage.go
│   └── storage_test.go   # Unit tests with temp directories
├── ui/
│   ├── board/
│   │   ├── board.go
│   │   └── board_test.go # Bubbletea tea.Test for Update/View
│   └── undo/
│       ├── undo.go
│       └── undo_test.go  # Command pattern tests
└── integration_test.go   # Full workflow tests
```

### Testing Techniques

**Unit tests** (standard `go test`):
- Test pure functions: `Search()`, `SortByPriority()`, `MatchesQuery()`
- Test struct methods: `Board.AddTask()`, `Board.MoveTask()`
- Use table-driven tests for edge cases

**Bubbletea component tests** (`tea.Test`):
```go
func TestBoardNavigation(t *testing.T) {
    m := NewBoardModel(testBoard)
    // Simulate key press
    m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
    if m.selectedIndex != 1 {
        t.Errorf("expected selection to move down")
    }
}
```

**Storage tests** (with temp directories):
```go
func TestSaveLoad(t *testing.T) {
    dir := t.TempDir()
    store := NewStorage(dir)
    board := &Board{Name: "Test"}
    store.SaveBoard(board)
    loaded, _ := store.LoadBoard(board.ID)
    // Assert equality
}
```

**Integration tests** (full workflows):
- Use `tea.Test` to simulate complete user journeys
- Verify state after multi-step operations
- Test mode transitions and keyboard handling

### What NOT to Test
- Lipgloss styling (visual, not logic)
- Terminal rendering internals (trust Bubbletea)
- Third-party library behavior

### Coverage Target
- **Model layer**: 90%+ (pure logic, easy to test)
- **Storage layer**: 80%+ (I/O, error paths)
- **UI layer**: 60%+ (key handlers, state transitions)
- **Overall**: 75%+
