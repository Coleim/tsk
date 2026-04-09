package ui

import (
	"fmt"
	"strings"

	"github.com/coliva/tsk/internal/model"
	"github.com/coliva/tsk/internal/styles"
)

// TaskDetail represents the task detail view state
type TaskDetail struct {
	Task  *model.Task
	Board *model.Board
}

// NewTaskDetail creates a new task detail view
func NewTaskDetail(task *model.Task, board *model.Board) *TaskDetail {
	return &TaskDetail{
		Task:  task,
		Board: board,
	}
}

// View renders the task detail view
func (td *TaskDetail) View(width, height int) string {
	if td.Task == nil {
		return "No task selected"
	}

	var lines []string

	// Header with back navigation hint
	header := styles.HelpHintStyle().Render("Press Esc or Enter to return")
	lines = append(lines, header)
	lines = append(lines, "")

	// Task title
	lines = append(lines, styles.ModalTitleStyle().Render(td.Task.Title))
	lines = append(lines, "")

	// Status
	statusLine := styles.PreviewLabelStyle().Render("Status:    ") +
		styles.StatusStyle(td.Task.Status).Render(td.Task.Status.String())
	lines = append(lines, statusLine)

	// Priority
	priorityLine := styles.PreviewLabelStyle().Render("Priority:  ") +
		styles.PriorityStyle(td.Task.Priority).Render(td.Task.Priority.String()) +
		" " + styles.PriorityIndicator(td.Task.Priority)
	lines = append(lines, priorityLine)

	// Due date
	if td.Task.DueDate != nil {
		dueLine := styles.PreviewLabelStyle().Render("Due:       ") +
			styles.PreviewValueStyle().Render(td.Task.DueDate.Format("Monday, January 2, 2006"))
		lines = append(lines, dueLine)
	}

	// Labels
	if len(td.Task.Labels) > 0 && td.Board != nil {
		var labelBadges []string
		for _, labelName := range td.Task.Labels {
			label := td.Board.GetLabel(labelName)
			labelBadges = append(labelBadges, styles.LabelBadge(label.Name, label.Color))
		}
		labelsLine := styles.PreviewLabelStyle().Render("Labels:    ") + strings.Join(labelBadges, " ")
		lines = append(lines, labelsLine)
	}

	// Created/Updated times
	lines = append(lines, "")
	createdLine := styles.HelpHintStyle().Render(fmt.Sprintf("Created: %s",
		td.Task.CreatedAt.Format("Jan 2, 2006 3:04 PM")))
	lines = append(lines, createdLine)

	updatedLine := styles.HelpHintStyle().Render(fmt.Sprintf("Updated: %s",
		td.Task.UpdatedAt.Format("Jan 2, 2006 3:04 PM")))
	lines = append(lines, updatedLine)

	// Description
	lines = append(lines, "")
	if td.Task.Description != "" {
		lines = append(lines, styles.PreviewLabelStyle().Render("Description:"))
		lines = append(lines, "")

		// Word wrap description
		wrapped := wordWrap(td.Task.Description, width-4)
		for _, line := range wrapped {
			lines = append(lines, styles.PreviewValueStyle().Render(line))
		}
	} else {
		lines = append(lines, styles.HelpHintStyle().Render("(No description)"))
	}

	// Add shortcuts at bottom
	lines = append(lines, "")
	lines = append(lines, styles.StatusLine2Style().Render("───────────────────────────────────────"))
	lines = append(lines, styles.StatusLine2Style().Render("e:edit  d:delete  1-3:priority  L:labels  t:due date"))

	content := strings.Join(lines, "\n")

	// Full-screen layout
	editWidth := width - 4
	if editWidth < 50 {
		editWidth = 50
	}

	box := styles.ModalStyle().Width(editWidth).Height(height - 4).Render(content)

	return box
}

// wordWrap wraps text to the specified width
func wordWrap(text string, width int) []string {
	if width <= 0 {
		width = 60
	}

	var lines []string
	words := strings.Fields(text)

	if len(words) == 0 {
		return lines
	}

	currentLine := words[0]
	for _, word := range words[1:] {
		if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}
