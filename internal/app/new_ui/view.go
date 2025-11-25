package new_ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/camikura/dito/internal/views"
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

	// Footer (changes based on focused pane)
	var footerContent string
	switch m.CurrentPane {
	case FocusPaneConnection:
		footerContent = "Switch Pane: tab"
	case FocusPaneTables:
		footerContent = "Navigate: ↑/↓ | Switch Pane: tab | Select: <enter>"
	case FocusPaneSQL:
		footerContent = "Switch Pane: tab | Edit SQL: e"
	case FocusPaneData:
		footerContent = "Navigate: ↑/↓ | Switch Pane: tab | Detail: <enter>"
	}

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
	if m.ConnectionMsg != "" {
		// Show error message if connection failed
		content = m.ConnectionMsg
		if len(content) > width-4 {
			content = content[:width-7] + "..."
		}
	} else if m.Endpoint != "" {
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

	// Prepare content lines
	var contentLines []string
	if len(m.Tables) == 0 {
		contentLines = []string{"No tables"}
	} else {
		// Render each table with proper selection/cursor highlighting
		for i, tableName := range m.Tables {
			var prefix string
			if i == m.SelectedTable {
				prefix = "* "
			} else {
				prefix = "  "
			}
			contentLines = append(contentLines, prefix+tableName)
		}
	}

	leftBorder := borderStyle.Render("│")
	rightBorder := borderStyle.Render("│")
	bottomBorder := borderStyle.Render("╰" + strings.Repeat("─", width-2) + "╯")

	// Styles for selection
	selectedStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(ColorPrimary)).
		Foreground(lipgloss.Color("#000000")).
		Width(width - 2)

	var result strings.Builder
	result.WriteString(title + "\n")

	// Render content lines (fixed height: 8 lines)
	for i := 0; i < 8; i++ {
		var line string
		if i < len(contentLines) {
			content := contentLines[i]
			// Apply full-span background if this is the cursor position
			if i == m.CursorTable && len(m.Tables) > 0 {
				line = selectedStyle.Render(" " + content + strings.Repeat(" ", width-len(content)-3))
				result.WriteString(leftBorder + line + rightBorder + "\n")
			} else {
				line = " " + content + strings.Repeat(" ", width-len(content)-3)
				result.WriteString(leftBorder + line + rightBorder + "\n")
			}
		} else {
			line = strings.Repeat(" ", width-2)
			result.WriteString(leftBorder + line + rightBorder + "\n")
		}
	}

	result.WriteString(bottomBorder)

	return result.String()
}

func renderSchemaPane(m Model, width int) string {
	titleText := " Schema "

	// Schema pane is never focused (display-only)
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorInactive))
	title := borderStyle.Render("╭─" + titleText + strings.Repeat("─", width-len(titleText)-3) + "╮")

	// Prepare content lines
	var contentLines []string
	if len(m.Tables) == 0 || m.CursorTable >= len(m.Tables) {
		contentLines = []string{"Select a table"}
	} else {
		tableName := m.Tables[m.CursorTable]
		details, exists := m.TableDetails[tableName]
		if !exists || details == nil {
			contentLines = []string{"Loading..."}
		} else {
			// Render schema information
			contentLines = append(contentLines, "Columns:")
			if details.Schema.DDL != "" {
				primaryKeys := views.ParsePrimaryKeysFromDDL(details.Schema.DDL)
				columns := views.ParseColumnsFromDDL(details.Schema.DDL, primaryKeys)
				for _, col := range columns {
					colLine := "  " + col.Name + " " + col.Type
					if len(colLine) > width-4 {
						colLine = colLine[:width-7] + "..."
					}
					contentLines = append(contentLines, colLine)
				}
			}
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
			content := contentLines[i]
			if len(content) > width-4 {
				content = content[:width-7] + "..."
			}
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

	// Data pane should fill entire screen height minus footer
	// Data pane structure: title(1) + content lines(?) + bottom(1)
	// Total height: totalHeight - 1 (footer)
	// So: content lines = totalHeight - 1 - 1 - 1 = totalHeight - 3
	contentLines := totalHeight - 3
	if contentLines < 5 {
		contentLines = 5
	}

	// Prepare content
	var displayLines []string
	if m.SelectedTable < 0 || m.SelectedTable >= len(m.Tables) {
		displayLines = []string{"No data"}
	} else {
		tableName := m.Tables[m.SelectedTable]
		data, exists := m.TableData[tableName]
		if !exists || data == nil {
			if m.LoadingData {
				displayLines = []string{"Loading..."}
			} else {
				displayLines = []string{"No data"}
			}
		} else {
			// Render data rows
			if len(data.Rows) == 0 {
				displayLines = []string{"No rows"}
			} else {
				// Simple data rendering - just show row count for now
				displayLines = append(displayLines, "Rows: "+string(rune(len(data.Rows)+48)))
				displayLines = append(displayLines, "")
				// TODO: Implement proper grid rendering in Phase 3
				displayLines = append(displayLines, "(Grid view coming in Phase 3)")
			}
		}
	}

	leftBorder := borderStyle.Render("│")
	rightBorder := borderStyle.Render("│")
	bottomBorder := borderStyle.Render("╰" + strings.Repeat("─", width-2) + "╯")

	var result strings.Builder
	result.WriteString(title + "\n")

	// Render content lines
	for i := 0; i < contentLines; i++ {
		var line string
		if i < len(displayLines) {
			content := displayLines[i]
			if len(content) > width-4 {
				content = content[:width-7] + "..."
			}
			line = " " + content + strings.Repeat(" ", width-len(content)-3)
		} else {
			line = strings.Repeat(" ", width-2)
		}
		result.WriteString(leftBorder + line + rightBorder + "\n")
	}

	result.WriteString(bottomBorder)

	return result.String()
}
