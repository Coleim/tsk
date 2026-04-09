package styles

import (
	"image/color"
	"os"

	"charm.land/lipgloss/v2"
	"github.com/coliva/tsk/internal/model"
)

// Theme defines the color palette for the application
type Theme struct {
	// Backgrounds
	Background color.Color
	Surface    color.Color
	Elevated   color.Color

	// Text
	TextPrimary   color.Color
	TextSecondary color.Color
	TextMuted     color.Color

	// Semantic
	Accent  color.Color
	Success color.Color
	Warning color.Color
	Error   color.Color

	// Border
	Border      color.Color
	BorderLight color.Color

	// Status
	StatusToDo       color.Color
	StatusInProgress color.Color
	StatusDone       color.Color

	// Priority
	PriorityHigh   color.Color
	PriorityMedium color.Color
	PriorityLow    color.Color
	PriorityNone   color.Color

	// Labels
	LabelRed    color.Color
	LabelOrange color.Color
	LabelYellow color.Color
	LabelGreen  color.Color
	LabelBlue   color.Color
	LabelPurple color.Color
	LabelPink   color.Color
	LabelCyan   color.Color
}

// DarkTheme - Catppuccin-inspired dark theme
var DarkTheme = Theme{
	// Backgrounds
	Background: lipgloss.Color("#1e1e2e"),
	Surface:    lipgloss.Color("#313244"),
	Elevated:   lipgloss.Color("#45475a"),

	// Text
	TextPrimary:   lipgloss.Color("#cdd6f4"),
	TextSecondary: lipgloss.Color("#a6adc8"),
	TextMuted:     lipgloss.Color("#6c7086"),

	// Semantic
	Accent:  lipgloss.Color("#cba6f7"),
	Success: lipgloss.Color("#a6e3a1"),
	Warning: lipgloss.Color("#f9e2af"),
	Error:   lipgloss.Color("#f38ba8"),

	// Border
	Border:      lipgloss.Color("#89b4fa"),
	BorderLight: lipgloss.Color("#585b70"),

	// Status
	StatusToDo:       lipgloss.Color("#89b4fa"),
	StatusInProgress: lipgloss.Color("#f9e2af"),
	StatusDone:       lipgloss.Color("#a6e3a1"),

	// Priority
	PriorityHigh:   lipgloss.Color("#f38ba8"),
	PriorityMedium: lipgloss.Color("#fab387"),
	PriorityLow:    lipgloss.Color("#6c7086"),
	PriorityNone:   lipgloss.Color("#585b70"),

	// Labels
	LabelRed:    lipgloss.Color("#f38ba8"),
	LabelOrange: lipgloss.Color("#fab387"),
	LabelYellow: lipgloss.Color("#f9e2af"),
	LabelGreen:  lipgloss.Color("#a6e3a1"),
	LabelBlue:   lipgloss.Color("#89b4fa"),
	LabelPurple: lipgloss.Color("#cba6f7"),
	LabelPink:   lipgloss.Color("#f5c2e7"),
	LabelCyan:   lipgloss.Color("#94e2d5"),
}

// LightTheme - Light mode colors
var LightTheme = Theme{
	// Backgrounds
	Background: lipgloss.Color("#eff1f5"),
	Surface:    lipgloss.Color("#e6e9ef"),
	Elevated:   lipgloss.Color("#dce0e8"),

	// Text
	TextPrimary:   lipgloss.Color("#4c4f69"),
	TextSecondary: lipgloss.Color("#6c6f85"),
	TextMuted:     lipgloss.Color("#9ca0b0"),

	// Semantic
	Accent:  lipgloss.Color("#8839ef"),
	Success: lipgloss.Color("#40a02b"),
	Warning: lipgloss.Color("#df8e1d"),
	Error:   lipgloss.Color("#d20f39"),

	// Border
	Border:      lipgloss.Color("#1e66f5"),
	BorderLight: lipgloss.Color("#bcc0cc"),

	// Status
	StatusToDo:       lipgloss.Color("#1e66f5"),
	StatusInProgress: lipgloss.Color("#df8e1d"),
	StatusDone:       lipgloss.Color("#40a02b"),

	// Priority
	PriorityHigh:   lipgloss.Color("#d20f39"),
	PriorityMedium: lipgloss.Color("#fe640b"),
	PriorityLow:    lipgloss.Color("#9ca0b0"),
	PriorityNone:   lipgloss.Color("#bcc0cc"),

	// Labels
	LabelRed:    lipgloss.Color("#d20f39"),
	LabelOrange: lipgloss.Color("#fe640b"),
	LabelYellow: lipgloss.Color("#df8e1d"),
	LabelGreen:  lipgloss.Color("#40a02b"),
	LabelBlue:   lipgloss.Color("#1e66f5"),
	LabelPurple: lipgloss.Color("#8839ef"),
	LabelPink:   lipgloss.Color("#ea76cb"),
	LabelCyan:   lipgloss.Color("#179299"),
}

// CurrentTheme holds the active theme
var CurrentTheme = &DarkTheme

// InitTheme initializes the theme based on TSK_THEME environment variable
func InitTheme() {
	theme := os.Getenv("TSK_THEME")
	switch theme {
	case "light":
		CurrentTheme = &LightTheme
	default:
		CurrentTheme = &DarkTheme
	}
}

// Legacy color variables for backward compatibility (will reference current theme)
var (
	ColorPrimary    = lipgloss.Color("39")
	ColorSecondary  = lipgloss.Color("245")
	ColorAccent     = lipgloss.Color("213")
	ColorSuccess    = lipgloss.Color("48")
	ColorWarning    = lipgloss.Color("220")
	ColorError      = lipgloss.Color("196")
	ColorMuted      = lipgloss.Color("243")
	ColorBackground = lipgloss.Color("236")
	ColorBorder     = lipgloss.Color("62")

	ColorToDo       = lipgloss.Color("75")
	ColorInProgress = lipgloss.Color("220")
	ColorDone       = lipgloss.Color("48")

	ColorPriorityHigh   = lipgloss.Color("196")
	ColorPriorityMedium = lipgloss.Color("214")
	ColorPriorityLow    = lipgloss.Color("244")
	ColorPriorityNone   = lipgloss.Color("245")

	ColorLabelRed    = lipgloss.Color("196")
	ColorLabelOrange = lipgloss.Color("214")
	ColorLabelYellow = lipgloss.Color("226")
	ColorLabelGreen  = lipgloss.Color("48")
	ColorLabelBlue   = lipgloss.Color("39")
	ColorLabelPurple = lipgloss.Color("135")
	ColorLabelPink   = lipgloss.Color("213")
	ColorLabelCyan   = lipgloss.Color("51")
)

// Style getter functions using current theme

// TitleStyle returns the app title style
func TitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(CurrentTheme.Accent)
}

// BoardNameStyle returns the board name style
func BoardNameStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(CurrentTheme.TextPrimary)
}

// HelpHintStyle returns the help hint style
func HelpHintStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(CurrentTheme.TextMuted)
}

// TabActiveStyle returns the active tab style
func TabActiveStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(CurrentTheme.TextPrimary).
		Background(CurrentTheme.Surface).
		Padding(0, 1)
}

// TabInactiveStyle returns the inactive tab style
func TabInactiveStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(CurrentTheme.TextMuted).
		Padding(0, 1)
}

// TabCountStyle returns the tab count style
func TabCountStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(CurrentTheme.TextMuted)
}

// TaskSelectedStyle returns the selected task style (with full border)
func TaskSelectedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(CurrentTheme.Accent).
		Background(CurrentTheme.Surface).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#cba6f7")). // Charm mauve/lavender
		Padding(0, 1)
}

// TaskNormalStyle returns the normal task style (with subtle border)
func TaskNormalStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(CurrentTheme.TextSecondary).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#6c7086")). // Charm overlay0
		Padding(0, 1)
}

// TaskListStyle returns the task list panel style
func TaskListStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(CurrentTheme.BorderLight).
		Padding(0, 1)
}

// PreviewStyle returns the preview panel style
func PreviewStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(CurrentTheme.BorderLight).
		Padding(1, 2)
}

// PreviewTitleStyle returns the preview title style
func PreviewTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(CurrentTheme.TextPrimary)
}

// PreviewLabelStyle returns the preview label style
func PreviewLabelStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(CurrentTheme.TextMuted)
}

// PreviewValueStyle returns the preview value style
func PreviewValueStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(CurrentTheme.TextSecondary)
}

// StatusLine1Style returns the status bar line 1 style
func StatusLine1Style() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(CurrentTheme.TextSecondary)
}

// StatusLine2Style returns the status bar line 2 style
func StatusLine2Style() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(CurrentTheme.TextMuted)
}

// StatusBarStyle returns the status bar container style
func StatusBarStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderTop(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(CurrentTheme.BorderLight).
		Background(CurrentTheme.Surface)
}

// ModeInsertStyle returns the insert mode indicator style
func ModeInsertStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(CurrentTheme.Success)
}

// ModeSearchStyle returns the search mode indicator style
func ModeSearchStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(CurrentTheme.StatusToDo)
}

// ModalStyle returns the modal overlay style
func ModalStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(CurrentTheme.Accent).
		Padding(1, 2)
}

// ModalTitleStyle returns the modal title style
func ModalTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(CurrentTheme.Accent)
}

// LabelStyle returns the label tag style
func LabelStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(CurrentTheme.Accent).
		Background(CurrentTheme.Surface).
		Padding(0, 1)
}

// ErrorStyle returns the error message style
func ErrorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(CurrentTheme.Error)
}

// SuccessStyle returns the success message style
func SuccessStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(CurrentTheme.Success)
}

// WarningStyle returns the warning message style
func WarningStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(CurrentTheme.Warning)
}

// EmptyStateStyle returns the empty state style
func EmptyStateStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(CurrentTheme.TextMuted).
		Italic(true)
}

// PanelTitleStyle returns the panel title style
func PanelTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(CurrentTheme.Accent).
		MarginBottom(1)
}

// PriorityStyle returns the style for a priority level
func PriorityStyle(p model.Priority) lipgloss.Style {
	var c color.Color
	switch p {
	case model.PriorityHigh:
		c = CurrentTheme.PriorityHigh
	case model.PriorityMedium:
		c = CurrentTheme.PriorityMedium
	case model.PriorityLow:
		c = CurrentTheme.PriorityLow
	default:
		c = CurrentTheme.PriorityNone
	}
	return lipgloss.NewStyle().Foreground(c)
}

// PriorityIndicator returns a styled priority symbol
func PriorityIndicator(p model.Priority) string {
	return PriorityStyle(p).Render(p.Symbol())
}

// PriorityColor returns the color for a priority level
func PriorityColor(p model.Priority) color.Color {
	switch p {
	case model.PriorityHigh:
		return CurrentTheme.PriorityHigh
	case model.PriorityMedium:
		return CurrentTheme.PriorityMedium
	case model.PriorityLow:
		return CurrentTheme.PriorityLow
	default:
		return CurrentTheme.PriorityNone
	}
}

// StatusStyle returns the style for a status
func StatusStyle(s model.Status) lipgloss.Style {
	var c color.Color
	switch s {
	case model.StatusToDo:
		c = CurrentTheme.StatusToDo
	case model.StatusInProgress:
		c = CurrentTheme.StatusInProgress
	case model.StatusDone:
		c = CurrentTheme.StatusDone
	default:
		c = CurrentTheme.TextSecondary
	}
	return lipgloss.NewStyle().Foreground(c)
}

// StatusColor returns the color for a status
func StatusColor(s model.Status) color.Color {
	switch s {
	case model.StatusToDo:
		return CurrentTheme.StatusToDo
	case model.StatusInProgress:
		return CurrentTheme.StatusInProgress
	case model.StatusDone:
		return CurrentTheme.StatusDone
	default:
		return CurrentTheme.TextSecondary
	}
}

// TabStyleForStatus returns an active tab style with status-specific color
func TabStyleForStatus(s model.Status, active bool) lipgloss.Style {
	color := StatusColor(s)
	if active {
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(color)
	}
	return lipgloss.NewStyle().
		Foreground(CurrentTheme.TextMuted)
}

// LabelColor returns the color for a label color
func LabelColor(c model.LabelColor) color.Color {
	switch c {
	case model.LabelColorRed:
		return CurrentTheme.LabelRed
	case model.LabelColorOrange:
		return CurrentTheme.LabelOrange
	case model.LabelColorYellow:
		return CurrentTheme.LabelYellow
	case model.LabelColorGreen:
		return CurrentTheme.LabelGreen
	case model.LabelColorBlue:
		return CurrentTheme.LabelBlue
	case model.LabelColorPurple:
		return CurrentTheme.LabelPurple
	case model.LabelColorPink:
		return CurrentTheme.LabelPink
	case model.LabelColorCyan:
		return CurrentTheme.LabelCyan
	default:
		return CurrentTheme.TextSecondary
	}
}

// ColoredLabelStyle returns a styled label with its color
func ColoredLabelStyle(c model.LabelColor) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(LabelColor(c)).
		Bold(true)
}

// LabelBadge returns a styled label badge
func LabelBadge(name string, c model.LabelColor) string {
	return ColoredLabelStyle(c).Render("[" + name + "]")
}

// PopupStyle returns the style for popup overlays
func PopupStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(CurrentTheme.Accent).
		Padding(0, 1).
		MarginLeft(2)
}

// PopupTitleStyle returns the style for popup titles
func PopupTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(CurrentTheme.Accent).
		Bold(true).
		MarginBottom(1)
}

// PopupItemStyle returns the style for normal popup items
func PopupItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(CurrentTheme.TextSecondary)
}

// PopupSelectedItemStyle returns the style for selected popup items
func PopupSelectedItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(CurrentTheme.Accent).
		Bold(true)
}

// TaskCardStyle returns style for task cards with priority accent
func TaskCardStyle(selected bool, priority model.Priority) lipgloss.Style {
	priorityColor := PriorityColor(priority)
	if selected {
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(CurrentTheme.TextPrimary).
			Background(CurrentTheme.Elevated).
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(priorityColor).
			Padding(0, 1).
			MarginBottom(1)
	}
	return lipgloss.NewStyle().
		Foreground(CurrentTheme.TextSecondary).
		Background(CurrentTheme.Surface).
		BorderLeft(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(priorityColor).
		Padding(0, 1).
		MarginBottom(1)
}
