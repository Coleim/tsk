package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/coliva/tsk/internal/storage"
	"github.com/coliva/tsk/internal/styles"
	"github.com/coliva/tsk/internal/ui"
)

// Version information - injected at build time via ldflags
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Initialize theme from TSK_THEME environment variable
	styles.InitTheme()

	// Initialize storage
	store, err := storage.NewStorage()
	if err != nil {
		fmt.Printf("Error initializing storage: %v\n", err)
		os.Exit(1)
	}

	// Handle command-line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "import":
			if len(os.Args) < 3 {
				fmt.Println("Usage: tsk import <file.json>")
				os.Exit(1)
			}
			board, err := store.ImportBoard(os.Args[2])
			if err != nil {
				fmt.Printf("Error importing board: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Successfully imported '%s' with %d tasks\n", board.Name, len(board.Tasks))
			return

		case "export":
			if len(os.Args) < 3 {
				fmt.Println("Usage: tsk export <board-name> [output.json]")
				os.Exit(1)
			}
			// Export handled in TUI for now
			fmt.Println("Use the E key in the application to export the current board")
			return

		case "help", "-h", "--help":
			printHelp()
			return

		case "version", "-v", "--version":
			fmt.Printf("tsk version %s\ncommit: %s\nbuilt: %s\n", version, commit, date)
			return

		case "list", "boards":
			boards, err := store.ListBoards()
			if err != nil {
				fmt.Printf("Error listing boards: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Found %d boards:\n", len(boards))
			for _, b := range boards {
				fmt.Printf("  - %s (%s): %d tasks\n", b.Name, b.ID, b.TaskCount)
			}
			return
		}
	}

	// Create and run the app
	app := ui.NewApp(store)
	p := tea.NewProgram(app, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println(`tsk - Terminal Task Manager

Usage:
  tsk                      Start the TUI application
  tsk import <file.json>   Import a board from JSON file
  tsk help                 Show this help message
  tsk version              Show version information

Keyboard shortcuts (in app):
  j/k         Navigate tasks
  h/l         Switch panes
  n           New task
  d           Delete task
  > / <       Move task between panes
  Enter       Edit task
  1/2/3/0     Set priority
  L           Manage labels
  f           Filter tasks
  /           Search tasks
  E           Export current board
  b           Switch boards
  B           Create new board
  ?           Show help
  q           Quit`)
}
