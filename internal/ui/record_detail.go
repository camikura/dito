package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RecordDetailConfig holds configuration for the record detail dialog.
type RecordDetailConfig struct {
	Row          map[string]interface{} // Row data to display
	Columns      []string               // Column names in display order
	Width        int                    // Dialog width (0 = auto 80% of screen)
	Height       int                    // Dialog height (0 = auto 80% of screen)
	ScrollOffset int                    // Current scroll position
	Title        string                 // Dialog title (default: "Record Details")
	BorderColor  string                 // Border color (default: ColorPrimaryHex)
	Padding      int                    // Inner padding (default: 1)
}

// RecordDetail represents a dialog for displaying record details.
type RecordDetail struct {
	config       RecordDetailConfig
	content      string // Rendered content from VerticalTable
	lines        []string
	contentWidth int
	innerWidth   int
	vScrollBar   *VerticalScrollBar
}

// NewRecordDetail creates a new RecordDetail component.
func NewRecordDetail(config RecordDetailConfig) *RecordDetail {
	// Set defaults
	if config.Title == "" {
		config.Title = " Record Details "
	}
	if config.BorderColor == "" {
		config.BorderColor = ColorPrimaryHex
	}
	if config.Padding == 0 {
		config.Padding = 1
	}

	// Create vertical table and render content
	vt := VerticalTable{
		Data: config.Row,
		Keys: config.Columns,
	}
	content := vt.Render()
	lines := strings.Split(content, "\n")

	// Calculate content dimensions
	contentWidth := config.Width - 2 // Subtract left/right borders
	if contentWidth < 1 {
		contentWidth = 1
	}

	// Inner width = content width - padding on both sides
	innerWidth := contentWidth - (config.Padding * 2)
	if innerWidth < 1 {
		innerWidth = 1
	}

	// Calculate content height (excluding borders)
	contentHeight := config.Height - 2
	if contentHeight < 1 {
		contentHeight = 1
	}

	// Clamp scroll offset
	maxScroll := len(lines) - contentHeight
	if maxScroll < 0 {
		maxScroll = 0
	}
	if config.ScrollOffset > maxScroll {
		config.ScrollOffset = maxScroll
	}
	if config.ScrollOffset < 0 {
		config.ScrollOffset = 0
	}

	// Create vertical scrollbar
	vScrollBar := NewVerticalScrollBar(len(lines), contentHeight, config.ScrollOffset, contentHeight)

	return &RecordDetail{
		config:       config,
		content:      content,
		lines:        lines,
		contentWidth: contentWidth,
		innerWidth:   innerWidth,
		vScrollBar:   vScrollBar,
	}
}

// TotalLines returns the total number of content lines.
func (r *RecordDetail) TotalLines() int {
	return len(r.lines)
}

// MaxScroll returns the maximum scroll offset.
func (r *RecordDetail) MaxScroll() int {
	contentHeight := r.config.Height - 2
	if contentHeight < 1 {
		contentHeight = 1
	}
	maxScroll := len(r.lines) - contentHeight
	if maxScroll < 0 {
		return 0
	}
	return maxScroll
}

// Render renders the dialog as a string.
func (r *RecordDetail) Render() string {
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(r.config.BorderColor))
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(r.config.BorderColor)).Bold(true)

	var dialog strings.Builder

	// Title
	titleText := r.config.Title
	title := titleStyle.Render(titleText)
	titleLen := len([]rune(titleText))

	// Top border: ╭ + title + ─ ... ─ + ╮
	dashesLen := r.contentWidth - titleLen
	if dashesLen < 0 {
		dashesLen = 0
	}
	dialog.WriteString(borderStyle.Render("╭"))
	dialog.WriteString(title)
	dialog.WriteString(borderStyle.Render(strings.Repeat("─", dashesLen)))
	dialog.WriteString(borderStyle.Render("╮"))
	dialog.WriteString("\n")

	// Content height (excluding borders)
	contentHeight := r.config.Height - 2
	if contentHeight < 1 {
		contentHeight = 1
	}

	// Padding string
	paddingStr := strings.Repeat(" ", r.config.Padding)

	// Content lines
	for i := 0; i < contentHeight; i++ {
		lineIndex := r.config.ScrollOffset + i
		var lineContent string
		if lineIndex < len(r.lines) {
			lineContent = r.lines[lineIndex]
		} else {
			lineContent = ""
		}

		// Calculate visible width (excluding ANSI escape codes)
		visibleWidth := lipgloss.Width(lineContent)

		// Pad to fill inner width
		if visibleWidth < r.innerWidth {
			lineContent = lineContent + strings.Repeat(" ", r.innerWidth-visibleWidth)
		}

		// Get right border character (with scrollbar indicator)
		rightBorderChar := r.vScrollBar.GetCharAt(i)

		dialog.WriteString(borderStyle.Render("│"))
		dialog.WriteString(paddingStr)
		dialog.WriteString(lineContent)
		dialog.WriteString(paddingStr)
		dialog.WriteString(borderStyle.Render(rightBorderChar))
		dialog.WriteString("\n")
	}

	// Bottom border: ╰ + ─ ... ─ + ╯
	dialog.WriteString(borderStyle.Render("╰"))
	dialog.WriteString(borderStyle.Render(strings.Repeat("─", r.contentWidth)))
	dialog.WriteString(borderStyle.Render("╯"))

	return dialog.String()
}

// RenderCentered renders the dialog centered on screen using lipgloss.Place.
func (r *RecordDetail) RenderCentered(screenWidth, screenHeight int) string {
	return lipgloss.Place(
		screenWidth,
		screenHeight,
		lipgloss.Center,
		lipgloss.Center,
		r.Render(),
	)
}
