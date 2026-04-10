package ui

import (
	"fmt"
	"os"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/coliva/tsk/internal/model"
	"github.com/coliva/tsk/internal/storage"
)

// createLargeBoard creates a board with the specified number of tasks
func createLargeBoard(taskCount int) *model.Board {
	board := model.NewBoard("perf-test", "Performance Test Board")

	statuses := []model.Status{model.StatusToDo, model.StatusInProgress, model.StatusDone}

	for i := 0; i < taskCount; i++ {
		status := statuses[i%3]
		task := model.NewTask(
			fmt.Sprintf("task-%04d", i),
			fmt.Sprintf("Task %04d - Performance Test", i),
			status,
		)
		task.Priority = model.Priority(i % 4)
		task.Labels = []string{"perf-test", fmt.Sprintf("batch-%d", i/100)}
		task.Description = "Performance test task description for benchmarking purposes"
		board.AddTask(task)
	}

	return board
}

// BenchmarkBoardLoad measures time to load a large board
func BenchmarkBoardLoad(b *testing.B) {
	dir, _ := os.MkdirTemp("", "tsk-bench-*")
	defer func() { _ = os.RemoveAll(dir) }()

	store, _ := storage.NewStorageWithPath(dir)
	board := createLargeBoard(500)
	_ = store.SaveBoard(board)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := store.LoadBoard("perf-test")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBoardSave measures time to save a large board
func BenchmarkBoardSave(b *testing.B) {
	dir, _ := os.MkdirTemp("", "tsk-bench-*")
	defer func() { _ = os.RemoveAll(dir) }()

	store, _ := storage.NewStorageWithPath(dir)
	board := createLargeBoard(500)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := store.SaveBoard(board)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkNavigation measures navigation performance
func BenchmarkNavigation(b *testing.B) {
	board := createLargeBoard(500)
	state := model.NewAppState()
	state.SetBoard(board)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate j/k navigation
		state.SelectNext()
		state.SelectNext()
		state.SelectPrev()
		state.NextPane()
		state.PrevPane()
	}
}

// BenchmarkSearch measures search performance
func BenchmarkSearch(b *testing.B) {
	board := createLargeBoard(500)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = board.Search("Task 250")
	}
}

// BenchmarkSearchAll measures searching all tasks
func BenchmarkSearchAll(b *testing.B) {
	board := createLargeBoard(500)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = board.Search("Performance")
	}
}

// BenchmarkFilter measures filter performance
func BenchmarkFilter(b *testing.B) {
	board := createLargeBoard(500)
	state := model.NewAppState()
	state.SetBoard(board)

	state.FilterPriorities = []model.Priority{model.PriorityHigh}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = state.CurrentTasks()
	}
}

// BenchmarkAppView measures View rendering performance
func BenchmarkAppView(b *testing.B) {
	dir, _ := os.MkdirTemp("", "tsk-bench-*")
	defer func() { _ = os.RemoveAll(dir) }()

	store, _ := storage.NewStorageWithPath(dir)
	board := createLargeBoard(500)
	_ = store.SaveBoard(board)

	app := NewApp(store)
	app.state.SetBoard(board)
	app.state.Width = 120
	app.state.Height = 40

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = app.View()
	}
}

// BenchmarkKeyHandling measures key event handling
func BenchmarkKeyHandling(b *testing.B) {
	dir, _ := os.MkdirTemp("", "tsk-bench-*")
	defer func() { _ = os.RemoveAll(dir) }()

	store, _ := storage.NewStorageWithPath(dir)
	board := createLargeBoard(500)
	_ = store.SaveBoard(board)

	app := NewApp(store)
	app.state.SetBoard(board)
	app.state.Width = 120
	app.state.Height = 40

	keys := []tea.KeyMsg{
		{Type: tea.KeyRunes, Runes: []rune{'j'}},
		{Type: tea.KeyRunes, Runes: []rune{'k'}},
		{Type: tea.KeyRunes, Runes: []rune{'l'}},
		{Type: tea.KeyRunes, Runes: []rune{'h'}},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			app.Update(key)
		}
	}
}

// TestPerformanceThresholds validates performance meets requirements
func TestPerformanceThresholds(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping performance test in short mode")
	}

	dir, _ := os.MkdirTemp("", "tsk-perf-*")
	defer func() { _ = os.RemoveAll(dir) }()

	store, _ := storage.NewStorageWithPath(dir)
	board := createLargeBoard(600)
	_ = store.SaveBoard(board)

	// Test 1: Board load time < 500ms
	t.Run("BoardLoadUnder500ms", func(t *testing.T) {
		start := time.Now()
		_, err := store.LoadBoard("perf-test")
		elapsed := time.Since(start)

		if err != nil {
			t.Fatalf("Failed to load board: %v", err)
		}
		if elapsed > 500*time.Millisecond {
			t.Errorf("Board load took %v, expected < 500ms", elapsed)
		}
		t.Logf("Board load: %v", elapsed)
	})

	// Test 2: Search completes < 100ms
	t.Run("SearchUnder100ms", func(t *testing.T) {
		start := time.Now()
		results := board.Search("Performance")
		elapsed := time.Since(start)

		if elapsed > 100*time.Millisecond {
			t.Errorf("Search took %v, expected < 100ms", elapsed)
		}
		t.Logf("Search: %v (%d results)", elapsed, len(results))
	})

	// Test 3: Navigation < 10ms
	t.Run("NavigationUnder10ms", func(t *testing.T) {
		state := model.NewAppState()
		state.SetBoard(board)

		start := time.Now()
		for i := 0; i < 100; i++ {
			state.SelectNext()
		}
		elapsed := time.Since(start)

		avgPerOp := elapsed / 100
		if avgPerOp > 10*time.Millisecond {
			t.Errorf("Navigation avg %v per op, expected < 10ms", avgPerOp)
		}
		t.Logf("Navigation avg: %v per op", avgPerOp)
	})

	// Test 4: View render < 50ms
	t.Run("ViewRenderUnder50ms", func(t *testing.T) {
		app := NewApp(store)
		app.state.SetBoard(board)
		app.state.Width = 120
		app.state.Height = 40

		start := time.Now()
		_ = app.View()
		elapsed := time.Since(start)

		if elapsed > 50*time.Millisecond {
			t.Errorf("View render took %v, expected < 50ms", elapsed)
		}
		t.Logf("View render: %v", elapsed)
	})

	// Test 5: Filter < 50ms
	t.Run("FilterUnder50ms", func(t *testing.T) {
		state := model.NewAppState()
		state.SetBoard(board)
		state.FilterPriorities = []model.Priority{model.PriorityHigh}

		start := time.Now()
		_ = state.CurrentTasks()
		elapsed := time.Since(start)

		if elapsed > 50*time.Millisecond {
			t.Errorf("Filter took %v, expected < 50ms", elapsed)
		}
		t.Logf("Filter: %v", elapsed)
	})
}
