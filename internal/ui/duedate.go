package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/coliva/tsk/internal/model"
	"github.com/coliva/tsk/internal/styles"
)

// Quick date options
var quickDateOptions = []string{"today", "tomorrow", "next week", "next month", "clear"}

// DueDateEditor handles due date setting for a task
type DueDateEditor struct {
	Task       *model.Task
	input      textinput.Model
	err        string
	hasDate    bool
	parsedDate time.Time
	quickIdx   int // -1 means no quick option selected
}

// NewDueDateEditor creates a new due date editor
func NewDueDateEditor(task *model.Task) *DueDateEditor {
	input := textinput.New()
	input.Placeholder = "YYYY-MM-DD or Tab for quick options"
	input.CharLimit = 20

	hasDate := false
	if task.DueDate != nil {
		input.SetValue(task.DueDate.Format("2006-01-02"))
		hasDate = true
	}

	return &DueDateEditor{
		Task:     task,
		input:    input,
		hasDate:  hasDate,
		quickIdx: -1,
	}
}

// Focus activates the editor
func (de *DueDateEditor) Focus() tea.Cmd {
	return de.input.Focus()
}

// Update handles text input
func (de *DueDateEditor) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	de.input, cmd = de.input.Update(msg)
	de.validateDate()
	return cmd
}

func (de *DueDateEditor) validateDate() {
	de.err = ""
	value := strings.TrimSpace(de.input.Value())
	if value == "" {
		de.hasDate = false
		return
	}

	// Try various date formats
	formats := []string{
		"2006-01-02",
		"2006/01/02",
		"01/02/2006",
		"Jan 2, 2006",
		"January 2, 2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, value); err == nil {
			de.parsedDate = t
			de.hasDate = true
			return
		}
	}

	de.err = "Invalid date format"
	de.hasDate = false
}

// GetDueDate returns the parsed due date or nil
func (de *DueDateEditor) GetDueDate() *time.Time {
	if de.hasDate {
		return &de.parsedDate
	}
	return nil
}

// HasError returns true if there's a validation error
func (de *DueDateEditor) HasError() bool {
	return de.err != ""
}

// View renders the due date editor as popup content
func (de *DueDateEditor) View() string {
	var lines []string

	lines = append(lines, styles.ModalTitleStyle().Render("Set Due Date"))
	lines = append(lines, "")

	// Current due date
	if de.Task.DueDate != nil {
		current := de.Task.DueDate.Format("Monday, January 2, 2006")
		lines = append(lines, styles.FormFieldLabelStyle().Render("Current: ")+styles.PreviewValueStyle().Render(current))
		lines = append(lines, "")
	}

	// Input with proper styling
	lines = append(lines, styles.ActiveIndicator()+styles.FormFieldActiveLabelStyle().Render("New date:"))
	lines = append(lines, "    "+de.input.View())

	// Error or preview
	if de.err != "" {
		lines = append(lines, styles.ErrorStyle().Render(de.err))
	} else if de.hasDate {
		preview := de.parsedDate.Format("Monday, January 2, 2006")
		lines = append(lines, styles.SuccessStyle().Render("✓ "+preview))
	}

	lines = append(lines, "")

	// Quick options with selection
	lines = append(lines, styles.HelpHintStyle().Render("Quick set (Tab to cycle):"))
	var opts []string
	for i, opt := range quickDateOptions {
		if i == de.quickIdx {
			opts = append(opts, styles.FormFieldActiveLabelStyle().Render("["+opt+"]"))
		} else {
			opts = append(opts, styles.HelpHintStyle().Render(opt))
		}
	}
	lines = append(lines, "  "+strings.Join(opts, "  "))
	lines = append(lines, "")

	lines = append(lines, styles.KeyboardHintBarStyle().Render("Enter:save  Esc:cancel  Tab:quick select"))

	content := strings.Join(lines, "\n")

	// Popup style - fixed width 60 for date display
	return styles.ModalStyle().Width(60).Render(content)
}

// HandleQuickDate handles quick date shortcuts
func (de *DueDateEditor) HandleQuickDate(input string) bool {
	now := time.Now()
	var newDate time.Time

	switch strings.ToLower(input) {
	case "today":
		newDate = now
	case "tomorrow":
		newDate = now.AddDate(0, 0, 1)
	case "next week":
		newDate = now.AddDate(0, 0, 7)
	case "next month":
		newDate = now.AddDate(0, 1, 0)
	case "clear":
		de.input.SetValue("")
		de.hasDate = false
		return true
	default:
		return false
	}

	de.input.SetValue(newDate.Format("2006-01-02"))
	de.parsedDate = newDate
	de.hasDate = true
	de.err = ""
	return true
}

// SetQuickDate sets a quick date
func (de *DueDateEditor) SetQuickDate(days int) {
	newDate := time.Now().AddDate(0, 0, days)
	de.input.SetValue(newDate.Format("2006-01-02"))
	de.parsedDate = newDate
	de.hasDate = true
	de.err = ""
}

// Clear clears the due date
func (de *DueDateEditor) Clear() {
	de.input.SetValue("")
	de.hasDate = false
	de.err = ""
}

// CurrentValue returns the current input value
func (de *DueDateEditor) CurrentValue() string {
	return de.input.Value()
}

// NextQuickOption cycles to the next quick date option
func (de *DueDateEditor) NextQuickOption() {
	de.quickIdx = (de.quickIdx + 1) % len(quickDateOptions)
	// Apply the quick option immediately for preview
	de.HandleQuickDate(quickDateOptions[de.quickIdx])
}

// HasQuickSelection returns true if a quick option is selected
func (de *DueDateEditor) HasQuickSelection() bool {
	return de.quickIdx >= 0
}

// SetValue sets the input value and validates
func (de *DueDateEditor) SetValue(value string) {
	de.input.SetValue(value)
	de.validateDate()
}

// String for debugging
func (de *DueDateEditor) String() string {
	if de.hasDate {
		return fmt.Sprintf("Due: %s", de.parsedDate.Format("2006-01-02"))
	}
	return "No due date"
}
