package ui

import (
	"strings"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/coliva/tsk/internal/styles"
)

// ModalType defines the type of modal dialog
type ModalType int

const (
	ModalConfirm ModalType = iota
	ModalInput
	ModalSelect
)

// Modal represents a dialog overlay
type Modal struct {
	Type    ModalType
	Title   string
	Message string

	// For ModalConfirm
	ConfirmText string
	CancelText  string

	// For ModalSelect
	Options       []string
	SelectedIndex int

	// Callbacks
	OnConfirm func()
	OnCancel  func()
	OnSelect  func(index int)

	// Dimensions
	Width  int
	Height int
}

// NewConfirmModal creates a confirmation dialog
func NewConfirmModal(title, message string) *Modal {
	return &Modal{
		Type:        ModalConfirm,
		Title:       title,
		Message:     message,
		ConfirmText: "Yes",
		CancelText:  "No",
		Width:       50,
	}
}

// NewSelectModal creates a selection dialog
func NewSelectModal(title string, options []string) *Modal {
	return &Modal{
		Type:          ModalSelect,
		Title:         title,
		Options:       options,
		SelectedIndex: 0,
		Width:         50,
	}
}

// SelectNext moves selection down
func (m *Modal) SelectNext() {
	if m.Type == ModalSelect && m.SelectedIndex < len(m.Options)-1 {
		m.SelectedIndex++
	}
}

// SelectPrev moves selection up
func (m *Modal) SelectPrev() {
	if m.Type == ModalSelect && m.SelectedIndex > 0 {
		m.SelectedIndex--
	}
}

// Confirm triggers the confirm action
func (m *Modal) Confirm() {
	if m.OnConfirm != nil {
		m.OnConfirm()
	}
}

// Cancel triggers the cancel action
func (m *Modal) Cancel() {
	if m.OnCancel != nil {
		m.OnCancel()
	}
}

// Select triggers the select action for current selection
func (m *Modal) Select() {
	if m.OnSelect != nil && m.Type == ModalSelect {
		m.OnSelect(m.SelectedIndex)
	}
}

// View renders the modal as a compact popup
func (m *Modal) View(screenWidth, screenHeight int) string {
	var lines []string

	// Title
	lines = append(lines, styles.ModalTitleStyle().Render(m.Title))
	lines = append(lines, "")

	switch m.Type {
	case ModalConfirm:
		lines = append(lines, m.Message)
		lines = append(lines, "")

		// Confirm/Cancel buttons
		yes := "[Y] " + m.ConfirmText
		no := "[N] " + m.CancelText
		buttonsLine := styles.PreviewLabelStyle().Render(yes) + "  " + styles.HelpHintStyle().Render(no)
		lines = append(lines, buttonsLine)

	case ModalSelect:
		maxWidth := lipgloss.Width(m.Title)
		for i, opt := range m.Options {
			var line string
			if i == m.SelectedIndex {
				line = styles.TaskSelectedStyle().Render(" ▶ " + opt)
			} else {
				line = styles.TaskNormalStyle().Render("   " + opt)
			}
			lines = append(lines, line)
			if lipgloss.Width(line) > maxWidth {
				maxWidth = lipgloss.Width(line)
			}
		}
		lines = append(lines, "")
		helpLine := styles.HelpHintStyle().Render("j/k:navigate Enter:select Esc:cancel")
		lines = append(lines, helpLine)
		if lipgloss.Width(helpLine) > maxWidth {
			maxWidth = lipgloss.Width(helpLine)
		}

		// Build content and set width for select modals
		content := strings.Join(lines, "\n")
		if maxWidth < 30 {
			maxWidth = 30
		}
		return styles.ModalStyle().Width(maxWidth).Render(content)
	}

	// For confirm modals, don't set width - let it expand naturally
	content := strings.Join(lines, "\n")
	return styles.ModalStyle().Render(content)
}
