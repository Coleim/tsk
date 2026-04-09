package ui

import (
	"strings"

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

// View renders the modal (full screen)
func (m *Modal) View(screenWidth, screenHeight int) string {
	var content strings.Builder

	// Title
	content.WriteString(styles.ModalTitleStyle().Render(m.Title))
	content.WriteString("\n\n")

	switch m.Type {
	case ModalConfirm:
		content.WriteString(m.Message)
		content.WriteString("\n\n")

		// Confirm/Cancel buttons
		yes := styles.PreviewLabelStyle().Render("[Y] " + m.ConfirmText)
		no := styles.HelpHintStyle().Render("[N] " + m.CancelText)
		content.WriteString(yes + "    " + no)

	case ModalSelect:
		for i, opt := range m.Options {
			if i == m.SelectedIndex {
				content.WriteString(styles.TaskSelectedStyle().Render(" ▶ " + opt))
			} else {
				content.WriteString(styles.TaskNormalStyle().Render("   " + opt))
			}
			content.WriteString("\n")
		}
		content.WriteString("\n")
		content.WriteString(styles.HelpHintStyle().Render("j/k to navigate, Enter to select, Esc to cancel"))
	}

	// Full-screen layout
	editWidth := screenWidth - 4
	if editWidth < 50 {
		editWidth = 50
	}

	modalBox := styles.ModalStyle().Width(editWidth).Height(screenHeight - 4).Render(content.String())

	return modalBox
}
