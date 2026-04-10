package model

// Mode represents the current input mode of the application
type Mode int

const (
	// ModeNormal is the default mode for navigation and commands
	ModeNormal Mode = iota
	// ModeInsert is for text input (creating/editing tasks)
	ModeInsert
	// ModeSearch is for search input
	ModeSearch
	// ModeModal is for dialog interactions
	ModeModal
	// ModeWelcome is for first-run board name input
	ModeWelcome
	// ModeDetail is for viewing task details
	ModeDetail
	// ModeEdit is for editing a task
	ModeEdit
	// ModeLabels is for managing task labels
	ModeLabels
	// ModeDueDate is for setting task due date
	ModeDueDate
	// ModeBoard is for board management
	ModeBoard
	// ModeFilter is for filtering tasks
	ModeFilter
	// ModeSort is for sorting tasks
	ModeSort
	// ModeStats is for viewing statistics
	ModeStats
)

// String returns the display name for the mode
func (m Mode) String() string {
	switch m {
	case ModeNormal:
		return "NORMAL"
	case ModeInsert:
		return "INSERT"
	case ModeSearch:
		return "SEARCH"
	case ModeModal:
		return "MODAL"
	case ModeWelcome:
		return "WELCOME"
	case ModeDetail:
		return "DETAIL"
	case ModeEdit:
		return "EDIT"
	case ModeLabels:
		return "LABELS"
	case ModeDueDate:
		return "DUE DATE"
	case ModeBoard:
		return "BOARD"
	case ModeFilter:
		return "FILTER"
	case ModeSort:
		return "SORT"
	case ModeStats:
		return "STATS"
	default:
		return ""
	}
}

// IsTextInput returns true if the mode accepts text input
func (m Mode) IsTextInput() bool {
	return m == ModeInsert || m == ModeSearch || m == ModeWelcome
}
