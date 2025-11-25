package new_ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Color definitions
const (
	ColorPrimary   = "#00D9FF" // Cyan for active borders
	ColorInactive  = "#666666" // Gray for inactive borders
	ColorGreen     = "#00FF00" // Green for connection status
	ColorLabel     = "#00D9FF" // Cyan for section labels
	ColorHelp      = "#888888" // Gray for help text
)

// RenderView renders the new UI
func RenderView(m Model) string {
	if m.Width == 0 {
		return "Loading..."
	}

	// Layout configuration
	leftPaneWidth := 30
	rightPaneWidth := m.Width - leftPaneWidth - 3 // -3 for borders

	// Render each pane
	connectionPane := renderConnectionPane(m, leftPaneWidth)
	tablesPane := renderTablesPane(m, leftPaneWidth)
	schemaPane := renderSchemaPane(m, leftPaneWidth)
	sqlPane := renderSQLPane(m, leftPaneWidth)
	dataPane := renderDataPane(m, rightPaneWidth)

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

	// Header (right-aligned connection info)
	headerContent := ""
	if m.Connected {
		headerContent = strings.Repeat(" ", m.Width-len(m.Endpoint)-14) + "Connected to " + m.Endpoint
	}
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorHelp)).
		Padding(0, 1).
		Width(m.Width - 2)
	header := headerStyle.Render(headerContent)

	// Footer
	footerContent := "Navigate: ↑/↓ | Switch Pane: tab | Detail: <enter> | SQL: e"
	footerPadding := m.Width - len(footerContent) - len("Dito") - 4
	if footerPadding < 0 {
		footerPadding = 0
	}
	footerContent += strings.Repeat(" ", footerPadding) + "Dito"
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorHelp)).
		Padding(0, 1).
		Width(m.Width - 2)
	footer := footerStyle.Render(footerContent)

	// Separator
	separator := strings.Repeat("─", m.Width-2)

	// Assemble content
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		separator,
		panes,
		separator,
		footer,
	)

	// Add borders
	borderStyleColor := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorPrimary))
	title := " Dito "
	topBorder := borderStyleColor.Render("╭──" + title + strings.Repeat("─", m.Width-10) + "╮")
	leftBorder := borderStyleColor.Render("│")
	rightBorder := borderStyleColor.Render("│")

	var result strings.Builder
	result.WriteString(topBorder + "\n")

	for _, line := range strings.Split(content, "\n") {
		result.WriteString(leftBorder + line + rightBorder + "\n")
	}

	bottomBorder := borderStyleColor.Render("╰" + strings.Repeat("─", m.Width-2) + "╯")
	result.WriteString(bottomBorder)

	return result.String()
}

func renderConnectionPane(m Model, width int) string {
	var titleText string
	if m.Connected {
		titleText = " Connection ✓ "
	} else {
		titleText = " Connection "
	}

	title := "╭─" + titleText + strings.Repeat("─", width-len(titleText)-3) + "╮"

	content := "(not configured)"
	if m.Endpoint != "" {
		content = m.Endpoint
	}

	borderColor := ColorInactive
	if m.CurrentPane == FocusPaneConnection {
		borderColor = ColorPrimary
	}

	contentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1).
		Width(width - 2)

	borderedContent := contentStyle.Render(content)

	// Add side borders
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(borderColor))
	leftBorder := borderStyle.Render("│")
	rightBorder := borderStyle.Render("│")
	bottomBorder := borderStyle.Render("╰" + strings.Repeat("─", width-2) + "╯")

	var result strings.Builder
	result.WriteString(title + "\n")
	result.WriteString(leftBorder + borderedContent + rightBorder + "\n")
	result.WriteString(bottomBorder)

	return result.String()
}

func renderTablesPane(m Model, width int) string {
	title := "╭─ Tables"
	if len(m.Tables) > 0 {
		title += lipgloss.NewStyle().Render(" (" + string(rune(len(m.Tables)+48)) + ")")
	}
	title += " " + strings.Repeat("─", width-len(title)-1) + "╮"

	content := "No tables"
	if len(m.Tables) > 0 {
		// Placeholder for table list
		content = "  users\n    addresses\n    phones\n  products\n  orders"
	}

	borderColor := ColorInactive
	if m.CurrentPane == FocusPaneTables {
		borderColor = ColorPrimary
	}

	paneStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Width(width).
		Height(8).
		Padding(0, 1)

	return title + "\n" + paneStyle.Render(content)
}

func renderSchemaPane(m Model, width int) string {
	title := "╭─ Schema " + strings.Repeat("─", width-len("╭─ Schema ")-1) + "╮"

	content := "Select a table"

	// Schema pane is never focused (display-only)
	paneStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(ColorInactive)).
		Width(width).
		Height(8).
		Padding(0, 1)

	return title + "\n" + paneStyle.Render(content)
}

func renderSQLPane(m Model, width int) string {
	title := "╭─ SQL " + strings.Repeat("─", width-len("╭─ SQL ")-1) + "╮"

	content := ""
	if m.CurrentSQL != "" {
		content = m.CurrentSQL
	}

	borderColor := ColorInactive
	if m.CurrentPane == FocusPaneSQL {
		borderColor = ColorPrimary
	}

	paneStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Width(width).
		Height(3).
		Padding(0, 1)

	return title + "\n" + paneStyle.Render(content)
}

func renderDataPane(m Model, width int) string {
	title := "╭─ Data " + strings.Repeat("─", width-len("╭─ Data ")-1) + "╮"

	content := "No data"

	borderColor := ColorInactive
	if m.CurrentPane == FocusPaneData {
		borderColor = ColorPrimary
	}

	paneHeight := m.Height - 6

	paneStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Width(width).
		Height(paneHeight).
		Padding(0, 1)

	return title + "\n" + paneStyle.Render(content)
}
