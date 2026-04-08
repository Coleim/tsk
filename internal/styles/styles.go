package styles

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/coliva/tsk/internal/model"
)

// Colors - Enhanced palette for better visual distinction
var (
	// Primary colors
	ColorPrimary    = lipgloss.Color("39")  // Bright blue
	ColorSecondary  = lipgloss.Color("245") // Medium gray
	ColorAccent     = lipgloss.Color("213") // Bright magenta/pink
	ColorSuccess    = lipgloss.Color("48")  // Bright green
	ColorWarning    = lipgloss.Color("220") // Bright yellow
	ColorError      = lipgloss.Color("196") // Bright red
	ColorMuted      = lipgloss.Color("243") // Muted gray
	ColorBackground = lipgloss.Color("236") // Dark background
	ColorBorder     = lipgloss.Color("62")  // Purple/blue border

	// Status-specific colors
	ColorToDo       = lipgloss.Color("75")  // Light blue
	ColorInProgress = lipgloss.Color("220") // Yellow/amber
	ColorDone       = lipgloss.Color("48")  // Green

	// Priority colors
	ColorPriorityHigh   = lipgloss.Color("196") // Red
	ColorPriorityMedium = lipgloss.Color("214") // Orange
	ColorPriorityLow    = lipgloss.Color("244") // Gray
	ColorPriorityNone   = lipgloss.Color("245") // Gray
)

// Base styles
var (
	// App title
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent)

	// Board name
	BoardNameStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("255"))

	// Help hint in header
	HelpHintStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	// Selected tab
	TabActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("255")).
			Background(ColorBackground).
			Padding(0, 1)

	// Inactive tab
	TabInactiveStyle = lipgloss.NewStyle().
				Foreground(ColorSecondary).
				Padding(0, 1)

	// Task count in tab
	TabCountStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	// Selected task
	TaskSelectedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("255"))

	// Unselected task
	TaskNormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	// Task list panel
	TaskListStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(0, 1)

	// Preview panel
	PreviewStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(1, 2)

	// Preview title
	PreviewTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("255"))

	// Preview field label
	PreviewLabelStyle = lipgloss.NewStyle().
				Foreground(ColorMuted)

	// Preview field value
	PreviewValueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))

	// Status bar line 1 (context/feedback)
	StatusLine1Style = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))

	// Status bar line 2 (shortcuts)
	StatusLine2Style = lipgloss.NewStyle().
				Foreground(ColorMuted)

	// Mode indicator
	ModeInsertStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSuccess)

	ModeSearchStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary)

	// Modal overlay
	ModalStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorAccent).
			Padding(1, 2)

	// Modal title
	ModalTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent)

	// Label tag
	LabelStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Background(ColorBackground).
			Padding(0, 1)

	// Error message
	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError)

	// Success message
	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	// Warning message
	WarningStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorWarning)
)

// PriorityStyle returns the style for a priority level
func PriorityStyle(p model.Priority) lipgloss.Style {
	var color lipgloss.Color
	switch p {
	case model.PriorityHigh:
		color = ColorPriorityHigh
	case model.PriorityMedium:
		color = ColorPriorityMedium
	case model.PriorityLow:
		color = ColorPriorityLow
	default:
		color = ColorPriorityNone
	}
	return lipgloss.NewStyle().Foreground(color)
}

// PriorityIndicator returns a styled priority symbol
func PriorityIndicator(p model.Priority) string {
	return PriorityStyle(p).Render(p.Symbol())
}

// StatusStyle returns the style for a status
func StatusStyle(s model.Status) lipgloss.Style {
	var color lipgloss.Color
	switch s {
	case model.StatusToDo:
		color = ColorToDo
	case model.StatusInProgress:
		color = ColorInProgress
	case model.StatusDone:
		color = ColorDone
	default:
		color = ColorSecondary
	}
	return lipgloss.NewStyle().Foreground(color)
}

// StatusColor returns the color for a status
func StatusColor(s model.Status) lipgloss.Color {
	switch s {
	case model.StatusToDo:
		return ColorToDo
	case model.StatusInProgress:
		return ColorInProgress
	case model.StatusDone:
		return ColorDone
	default:
		return ColorSecondary
	}
}

// TabStyleForStatus returns an active tab style with status-specific color
func TabStyleForStatus(s model.Status, active bool) lipgloss.Style {
	color := StatusColor(s)
	if active {
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(color).
			Background(ColorBackground).
			Padding(0, 1)
	}
	return lipgloss.NewStyle().
		Foreground(ColorSecondary).
		Padding(0, 1)
}
