package ui

import (
	"fmt"
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/coliva/tsk/internal/model"
	"github.com/coliva/tsk/internal/styles"
)

// StatsView renders the statistics overlay
type StatsView struct {
	stats   model.BoardStatistics
	width   int
	height  int
	compact bool // use compact layout when height is limited
}

// NewStatsView creates a new statistics view
func NewStatsView(board *model.Board, width, height int) *StatsView {
	return &StatsView{
		stats:   model.ComputeStatistics(board),
		width:   width,
		height:  height,
		compact: height < 40, // use compact mode more aggressively
	}
}

// View renders the statistics overlay
func (s *StatsView) View() string {
	var sections []string

	// Header
	header := s.renderHeader()
	sections = append(sections, header)

	// Always use compact layout - it's guaranteed to fit in any terminal
	// Full layout with cards is too tall and causes overflow issues
	sections = append(sections, s.renderCompactView())

	// Footer
	sections = append(sections, "")
	sections = append(sections, styles.DialogSeparator(50))
	sections = append(sections, styles.KeyboardHintBarStyle().Render("Press S, Esc, or q to close"))

	content := strings.Join(sections, "\n")

	// Calculate modal width - make it comfortably sized
	editWidth := s.width - 6
	if editWidth < 60 {
		editWidth = 60
	}
	if editWidth > 85 {
		editWidth = 85
	}

	return styles.ModalStyle().Width(editWidth).Render(content)
}

func (s *StatsView) renderHeader() string {
	title := styles.ModalTitleStyle().Render("📊 Board Statistics")
	boardName := styles.HelpHintStyle().Render(fmt.Sprintf("\"%s\"", s.stats.BoardName))
	return title + "\n" + boardName
}

// renderCompactView renders all stats in a condensed format for small terminals
func (s *StatsView) renderCompactView() string {
	var lines []string

	barWidth := 35

	// Task Distribution - inline
	lines = append(lines, "")
	lines = append(lines, styles.SectionCardTitleStyle().Render("Tasks"))
	maxStatus := max(s.stats.TodoCount, s.stats.InProgressCount, s.stats.DoneCount)
	if maxStatus == 0 {
		maxStatus = 1
	}
	todoBar := s.renderBarGraph(s.stats.TodoCount, maxStatus, barWidth, styles.CurrentTheme.StatusToDo)
	lines = append(lines, fmt.Sprintf("  To Do        %s  %d", todoBar, s.stats.TodoCount))
	inProgressBar := s.renderBarGraph(s.stats.InProgressCount, maxStatus, barWidth, styles.CurrentTheme.StatusInProgress)
	lines = append(lines, fmt.Sprintf("  In Progress  %s  %d", inProgressBar, s.stats.InProgressCount))
	doneBar := s.renderBarGraph(s.stats.DoneCount, maxStatus, barWidth, styles.CurrentTheme.StatusDone)
	lines = append(lines, fmt.Sprintf("  Done         %s  %d", doneBar, s.stats.DoneCount))
	lines = append(lines, fmt.Sprintf("  Total: %d tasks", s.stats.TotalCount))

	// Priority - with spacing
	lines = append(lines, "")
	lines = append(lines, styles.SectionCardTitleStyle().Render("Priority"))
	highLabel := lipgloss.NewStyle().Foreground(styles.CurrentTheme.PriorityHigh).Render("● High")
	medLabel := lipgloss.NewStyle().Foreground(styles.CurrentTheme.PriorityMedium).Render("◐ Medium")
	lowLabel := lipgloss.NewStyle().Foreground(styles.CurrentTheme.PriorityLow).Render("● Low")
	noneLabel := lipgloss.NewStyle().Foreground(styles.CurrentTheme.PriorityNone).Render("  None")
	lines = append(lines, fmt.Sprintf("  %s: %d   %s: %d   %s: %d   %s: %d",
		highLabel, s.stats.HighPriorityCount,
		medLabel, s.stats.MediumPriorityCount,
		lowLabel, s.stats.LowPriorityCount,
		noneLabel, s.stats.NoPriorityCount))

	// Due Dates - with spacing
	lines = append(lines, "")
	lines = append(lines, styles.SectionCardTitleStyle().Render("Due Dates"))
	overdueText := fmt.Sprintf("⚠ Overdue: %d", s.stats.OverdueCount)
	if s.stats.OverdueCount > 0 {
		overdueText = styles.ErrorStyle().Render(overdueText)
	}
	lines = append(lines, fmt.Sprintf("  %s   📅 Today: %d   📆 This Week: %d   ∅ No Date: %d",
		overdueText, s.stats.DueTodayCount, s.stats.DueThisWeek, s.stats.NoDueDateCount))

	// Completion
	lines = append(lines, "")
	lines = append(lines, styles.SectionCardTitleStyle().Render("Completion"))
	compBarWidth := 40
	filled := 0
	if s.stats.TotalCount > 0 {
		filled = int(float64(compBarWidth) * s.stats.CompletionPercent / 100)
	}
	filledStr := strings.Repeat("█", filled)
	emptyStr := strings.Repeat("░", compBarWidth-filled)
	bar := lipgloss.NewStyle().Foreground(styles.CurrentTheme.Success).Render(filledStr) +
		lipgloss.NewStyle().Foreground(styles.CurrentTheme.Elevated).Render(emptyStr)
	lines = append(lines, fmt.Sprintf("  %s  %.1f%%", bar, s.stats.CompletionPercent))
	lines = append(lines, fmt.Sprintf("  %d of %d tasks completed", s.stats.DoneCount, s.stats.TotalCount))

	if s.stats.TotalCount == 0 {
		lines = append(lines, "")
		lines = append(lines, s.renderEmptyState())
	}

	return strings.Join(lines, "\n")
}

func (s *StatsView) renderEmptyState() string {
	return styles.HelpHintStyle().Render("No tasks yet. Press 'n' to create your first task!")
}

// renderBarGraph renders a horizontal bar graph
func (s *StatsView) renderBarGraph(value, maxValue, width int, color color.Color) string {
	if maxValue == 0 {
		maxValue = 1
	}

	filled := width * value / maxValue
	if filled > width {
		filled = width
	}

	filledStr := strings.Repeat("█", filled)
	emptyStr := strings.Repeat("░", width-filled)

	filledStyled := lipgloss.NewStyle().Foreground(color).Render(filledStr)
	emptyStyled := lipgloss.NewStyle().Foreground(styles.CurrentTheme.Elevated).Render(emptyStr)

	return filledStyled + emptyStyled
}

// Helper function for max of multiple ints
func max(nums ...int) int {
	if len(nums) == 0 {
		return 0
	}
	m := nums[0]
	for _, n := range nums[1:] {
		if n > m {
			m = n
		}
	}
	return m
}
