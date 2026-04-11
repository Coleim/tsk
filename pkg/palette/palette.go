// Package palette provides a unified color palette for TUI applications.
// Colors are based on Catppuccin Mocha (dark) and Catppuccin Latte (light).
// This package has no dependencies, making it portable across TUI projects.
package palette

// Theme defines the color palette for an application.
// All colors are hex strings (e.g., "#1e1e2e").
type Theme struct {
	// Backgrounds
	Background string
	Surface    string
	Elevated   string

	// Text
	TextPrimary   string
	TextSecondary string
	TextMuted     string

	// Semantic
	Accent      string
	AccentHover string
	Success     string
	Warning     string
	Error       string

	// Border
	Border      string
	BorderLight string

	// Status
	StatusToDo       string
	StatusInProgress string
	StatusDone       string

	// Priority
	PriorityHigh   string
	PriorityMedium string
	PriorityLow    string
	PriorityNone   string

	// Labels
	LabelRed    string
	LabelOrange string
	LabelYellow string
	LabelGreen  string
	LabelBlue   string
	LabelPurple string
	LabelPink   string
	LabelCyan   string
}

// DarkTheme - Catppuccin Mocha inspired dark theme
var DarkTheme = Theme{
	// Backgrounds
	Background: "#1e1e2e",
	Surface:    "#313244",
	Elevated:   "#45475a",

	// Text
	TextPrimary:   "#cdd6f4",
	TextSecondary: "#a6adc8",
	TextMuted:     "#6c7086",

	// Semantic
	Accent:      "#74c7ec",
	AccentHover: "#8fd4f0",
	Success:     "#a6e3a1",
	Warning:     "#f9e2af",
	Error:       "#f38ba8",

	// Border
	Border:      "#89b4fa",
	BorderLight: "#585b70",

	// Status
	StatusToDo:       "#89b4fa",
	StatusInProgress: "#f9e2af",
	StatusDone:       "#a6e3a1",

	// Priority
	PriorityHigh:   "#f38ba8",
	PriorityMedium: "#fab387",
	PriorityLow:    "#6c7086",
	PriorityNone:   "#585b70",

	// Labels
	LabelRed:    "#f38ba8",
	LabelOrange: "#fab387",
	LabelYellow: "#f9e2af",
	LabelGreen:  "#a6e3a1",
	LabelBlue:   "#89b4fa",
	LabelPurple: "#cba6f7",
	LabelPink:   "#f5c2e7",
	LabelCyan:   "#94e2d5",
}

// LightTheme - Catppuccin Latte inspired light theme
var LightTheme = Theme{
	// Backgrounds
	Background: "#eff1f5",
	Surface:    "#e6e9ef",
	Elevated:   "#dce0e8",

	// Text
	TextPrimary:   "#4c4f69",
	TextSecondary: "#6c6f85",
	TextMuted:     "#9ca0b0",

	// Semantic
	Accent:      "#8839ef",
	AccentHover: "#9a4dff",
	Success:     "#40a02b",
	Warning:     "#df8e1d",
	Error:       "#d20f39",

	// Border
	Border:      "#1e66f5",
	BorderLight: "#bcc0cc",

	// Status
	StatusToDo:       "#1e66f5",
	StatusInProgress: "#df8e1d",
	StatusDone:       "#40a02b",

	// Priority
	PriorityHigh:   "#d20f39",
	PriorityMedium: "#fe640b",
	PriorityLow:    "#9ca0b0",
	PriorityNone:   "#bcc0cc",

	// Labels
	LabelRed:    "#d20f39",
	LabelOrange: "#fe640b",
	LabelYellow: "#df8e1d",
	LabelGreen:  "#40a02b",
	LabelBlue:   "#1e66f5",
	LabelPurple: "#8839ef",
	LabelPink:   "#ea76cb",
	LabelCyan:   "#179299",
}
