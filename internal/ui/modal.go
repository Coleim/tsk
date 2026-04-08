package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
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

// View renders the modal
func (m *Modal) View(screenWidth, screenHeight int) string {
	var content strings.Builder

	// Title
	content.WriteString(styles.ModalTitleStyle.Render(m.Title))
	content.WriteString("\n\n")

	switch m.Type {
	case ModalConfirm:
		content.WriteString(m.Message)
		content.WriteString("\n\n")

		// Confirm/Cancel buttons
		yes := styles.PreviewLabelStyle.Render("[Y] " + m.ConfirmText)
		no := styles.HelpHintStyle.Render("[N] " + m.CancelText)
		content.WriteString(yes + "    " + no)

	case ModalSelect:
		for i, opt := range m.Options {
			if i == m.SelectedIndex {
				content.WriteString(styles.TaskSelectedStyle.Render(" ▶ " + opt))
			} else {
				content.WriteString(styles.TaskNormalStyle.Render("   " + opt))
			}
			content.WriteString("\n")
		}
		content.WriteString("\n")
		content.WriteString(styles.HelpHintStyle.Render("j/k to navigate, Enter to select, Esc to cancel"))
	}

	// Create styled modal box
	modalBox := styles.ModalStyle.Width(m.Width).Render(content.String())

	// Calculate padding to center the modal
	modalHeight := strings.Count(modalBox, "\n") + 1
	paddingY := (screenHeight - modalHeight) / 2
	if paddingY < 0 {
		paddingY = 0
	}

	modalWidth := lipgloss.Width(modalBox)
	paddingX := (screenWidth - modalWidth) / 2
	if paddingX < 0 {
		paddingX = 0
	}

	// Create dimmed overlay background
	overlay := strings.Repeat("\n", paddingY)

	// Indent modal horizontally
	lines := strings.Split(modalBox, "\n")
	indent := strings.Repeat(" ", paddingX)
	for i, line := range lines {
		lines[i] = indent + line
	}
	modalBox = strings.Join(lines, "\n")

	return overlay + modalBox
}
