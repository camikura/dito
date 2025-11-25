package new_ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Color definitions
const (
	ColorPrimary  = "#00D9FF" // Cyan for active borders
	ColorInactive = "#AAAAAA" // Light gray for inactive borders
	ColorGreen    = "#00FF00" // Green for connection status
	ColorLabel    = "#00D9FF" // Cyan for section labels
	ColorHelp     = "#888888" // Gray for help text
)

// RenderView renders the new UI
func RenderView(m Model) string {
	if m.Width == 0 {
		return "Loading..."
	}

	// Layout configuration
	leftPaneWidth := 30
	rightPaneWidth := m.Width - leftPaneWidth - 3 // -3 for space between panes

	// Render each pane
	connectionPane := renderConnectionPane(m, leftPaneWidth)
	tablesPane := renderTablesPane(m, leftPaneWidth)
	schemaPane := renderSchemaPane(m, leftPaneWidth)
	sqlPane := renderSQLPane(m, leftPaneWidth)
	dataPane := renderDataPane(m, rightPaneWidth, m.Height)

	// Join left panes vertically
	leftPanes := lipgloss.JoinVertical(
		lipgloss.Left,
		connectionPane,
		tablesPane,
		schemaPane,
		sqlPane,
	)

	// Join left and right panes horizontally
	panes := lipgloss.JoinHorizontal(lipgloss.Top, leftPanes, dataPane)

	// Footer
	footerContent := "Navigate: ↑/↓ | Switch Pane: tab | Detail: <enter> | SQL: e"
	footerPadding := m.Width - len(footerContent) - len("Dito") - 1
	if footerPadding < 0 {
		footerPadding = 0
	}
	footerContent += strings.Repeat(" ", footerPadding) + "Dito"

	// Assemble final output
	var result strings.Builder
	result.WriteString(panes + "\n")
	result.WriteString(footerContent)

	return result.String()
}

func renderConnectionPane(m Model, width int) string {
	var titleText string
	if m.Connected {
		titleText = " Connection ✓ "
	} else {
		titleText = " Connection "
	}

	borderColor := ColorInactive
	if m.CurrentPane == FocusPaneConnection {
		borderColor = ColorPrimary
	}

	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(borderColor))
	title := borderStyle.Render("╭─" + titleText + strings.Repeat("─", width-len(titleText)-3) + "╮")

	content := "(not configured)"
	if m.Endpoint != "" {
		content = m.Endpoint
	}

	// Pad content to width
	contentPadded := " " + content + strings.Repeat(" ", width-len(content)-3)

	leftBorder := borderStyle.Render("│")
	rightBorder := borderStyle.Render("│")
	bottomBorder := borderStyle.Render("╰" + strings.Repeat("─", width-2) + "╯")

	var result strings.Builder
	result.WriteString(title + "\n")
	result.WriteString(leftBorder + contentPadded + rightBorder + "\n")
	result.WriteString(bottomBorder)

	return result.String()
}

func renderTablesPane(m Model, width int) string {
	titleText := " Tables"
	if len(m.Tables) > 0 {
		titleText += " (" + string(rune(len(m.Tables)+48)) + ")"
	}
	titleText += " "

	borderColor := ColorInactive
	if m.CurrentPane == FocusPaneTables {
		borderColor = ColorPrimary
	}

	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(borderColor))
	title := borderStyle.Render("╭─" + titleText + strings.Repeat("─", width-len(titleText)-3) + "╮")

	// Content lines
	contentLines := []string{"No tables"}
	if len(m.Tables) > 0 {
		contentLines = []string{
			"  users",
			"    addresses",
			"    phones",
			"  products",
			"  orders",
		}
	}

	leftBorder := borderStyle.Render("│")
	rightBorder := borderStyle.Render("│")
	bottomBorder := borderStyle.Render("╰" + strings.Repeat("─", width-2) + "╯")

	var result strings.Builder
	result.WriteString(title + "\n")

	// Render content lines (fixed height: 8 lines)
	for i := 0; i < 8; i++ {
		var line string
		if i < len(contentLines) {
			line = " " + contentLines[i] + strings.Repeat(" ", width-len(contentLines[i])-3)
		} else {
			line = strings.Repeat(" ", width-2)
		}
		result.WriteString(leftBorder + line + rightBorder + "\n")
	}

	result.WriteString(bottomBorder)

	return result.String()
}

func renderSchemaPane(m Model, width int) string {
	titleText := " Schema "

	// Schema pane is never focused (display-only)
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorInactive))
	title := borderStyle.Render("╭─" + titleText + strings.Repeat("─", width-len(titleText)-3) + "╮")

	content := "Select a table"

	leftBorder := borderStyle.Render("│")
	rightBorder := borderStyle.Render("│")
	bottomBorder := borderStyle.Render("╰" + strings.Repeat("─", width-2) + "╯")

	var result strings.Builder
	result.WriteString(title + "\n")

	// Render content lines (fixed height: 8 lines)
	for i := 0; i < 8; i++ {
		var line string
		if i == 0 {
			line = " " + content + strings.Repeat(" ", width-len(content)-3)
		} else {
			line = strings.Repeat(" ", width-2)
		}
		result.WriteString(leftBorder + line + rightBorder + "\n")
	}

	result.WriteString(bottomBorder)

	return result.String()
}

func renderSQLPane(m Model, width int) string {
	titleText := " SQL "

	borderColor := ColorInactive
	if m.CurrentPane == FocusPaneSQL {
		borderColor = ColorPrimary
	}

	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(borderColor))
	title := borderStyle.Render("╭─" + titleText + strings.Repeat("─", width-len(titleText)-3) + "╮")

	content := ""
	if m.CurrentSQL != "" {
		content = m.CurrentSQL
	}

	leftBorder := borderStyle.Render("│")
	rightBorder := borderStyle.Render("│")
	bottomBorder := borderStyle.Render("╰" + strings.Repeat("─", width-2) + "╯")

	contentPadded := " " + content + strings.Repeat(" ", width-len(content)-3)

	var result strings.Builder
	result.WriteString(title + "\n")
	result.WriteString(leftBorder + contentPadded + rightBorder + "\n")
	result.WriteString(bottomBorder)

	return result.String()
}

func renderDataPane(m Model, width int, totalHeight int) string {
	titleText := " Data "

	borderColor := ColorInactive
	if m.CurrentPane == FocusPaneData {
		borderColor = ColorPrimary
	}

	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(borderColor))
	title := borderStyle.Render("╭─" + titleText + strings.Repeat("─", width-len(titleText)-3) + "╮")

	content := "No data"

	// Calculate content line count based on total height
	// Left panes: Connection(3) + Tables(10) + Schema(10) + SQL(3) = 26 lines
	// Footer: 1 line
	// Total used: 27 lines
	// Data pane total should be: totalHeight - 1 (footer)
	// Data pane structure: title(1) + content lines(?) + bottom(1)
	// So: content lines = totalHeight - 1 - 1 - 1 = totalHeight - 3
	contentLines := totalHeight - 3
	if contentLines < 24 {
		contentLines = 24 // minimum to match left panes
	}

	leftBorder := borderStyle.Render("│")
	rightBorder := borderStyle.Render("│")
	bottomBorder := borderStyle.Render("╰" + strings.Repeat("─", width-2) + "╯")

	var result strings.Builder
	result.WriteString(title + "\n")

	// Render content lines
	for i := 0; i < contentLines; i++ {
		var line string
		if i == 0 {
			line = " " + content + strings.Repeat(" ", width-len(content)-3)
		} else {
			line = strings.Repeat(" ", width-2)
		}
		result.WriteString(leftBorder + line + rightBorder + "\n")
	}

	result.WriteString(bottomBorder)

	return result.String()
}
