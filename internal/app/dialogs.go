package app

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/camikura/dito/internal/ui"
)

// renderRecordDetailDialog renders a dialog showing the details of the selected record
func renderRecordDetailDialog(m Model) string {
	if !m.RecordDetailVisible {
		return ""
	}

	// Get the selected row
	if m.SelectedTable < 0 || m.SelectedTable >= len(m.Tables) {
		return ""
	}

	tableName := m.Tables[m.SelectedTable]
	data, exists := m.TableData[tableName]
	if !exists || data == nil || len(data.Rows) == 0 {
		return ""
	}

	if m.SelectedDataRow < 0 || m.SelectedDataRow >= len(data.Rows) {
		return ""
	}

	row := data.Rows[m.SelectedDataRow]

	// Get columns in schema order
	columns := getColumnsInSchemaOrder(m, tableName, data.Rows)

	// Calculate dialog dimensions (80% of screen)
	dialogWidth := m.Width * 4 / 5
	dialogHeight := m.Height * 4 / 5

	// Create record detail component
	rd := ui.NewRecordDetail(ui.RecordDetailConfig{
		Row:          row,
		Columns:      columns,
		Width:        dialogWidth,
		Height:       dialogHeight,
		ScrollOffset: m.RecordDetailScroll,
		BorderColor:  ui.ColorPrimaryHex,
	})

	return rd.RenderCentered(m.Width, m.Height)
}

// renderConnectionDialog renders the connection setup dialog
func renderConnectionDialog(m Model) string {
	dialogWidth := ui.ConnectionDialogWidth

	// Styles
	borderStyle := ui.StyleBorderActive
	labelStyle := ui.StyleTitleActive

	var dialog strings.Builder

	// Title
	titleText := " Connection Setup "
	title := ui.StyleTitleBold.Render(titleText)
	titleLen := len([]rune(titleText))

	// Title line: ╭─ + title + ─...─ + ╮
	// dialogWidth = 1(╭) + 1(─) + titleLen + dashesLen + 1(╮)
	dashesLen := dialogWidth - 3 - titleLen
	if dashesLen < 0 {
		dashesLen = 0
	}
	dialog.WriteString(borderStyle.Render("╭─"))
	dialog.WriteString(title)
	dialog.WriteString(borderStyle.Render(strings.Repeat("─", dashesLen)))
	dialog.WriteString(borderStyle.Render("╮"))
	dialog.WriteString("\n")

	// Content width (between left border+space and space+right border)
	contentWidth := dialogWidth - 4

	// Helper function to render a field line with fixed width
	renderFieldLine := func(label string, value string, fieldIndex int, cursorPos int) string {
		var line strings.Builder
		line.WriteString(borderStyle.Render("│"))
		line.WriteString(" ")

		// Calculate label part: "Label: "
		labelPart := label + ": "
		labelLen := len(labelPart)

		// Value area width
		valueAreaWidth := contentWidth - labelLen
		if valueAreaWidth < 1 {
			valueAreaWidth = 1
		}

		// Convert value to runes for proper multi-byte character handling
		valueRunes := []rune(value)
		valueDisplayWidth := lipgloss.Width(value)

		// Build value display
		var valueDisplay string
		if m.ConnectionDialogField == fieldIndex {
			// Focused text field: cursor only (no background)
			if cursorPos < len(valueRunes) {
				beforeCursor := string(valueRunes[:cursorPos])
				cursorChar := string(valueRunes[cursorPos])
				cursorCharWidth := lipgloss.Width(cursorChar)
				afterCursor := string(valueRunes[cursorPos+1:])
				padding := valueAreaWidth - valueDisplayWidth
				if padding < 0 {
					padding = 0
				}
				// Use different cursor style for narrow (width=1) vs wide (width>=2) characters
				var cursorBlock string
				if cursorCharWidth >= 2 {
					cursorBlock = ui.CursorWide.Render(cursorChar)
				} else {
					cursorBlock = ui.CursorNarrow.Render(cursorChar)
				}
				valueDisplay = beforeCursor + cursorBlock + afterCursor + strings.Repeat(" ", padding)
			} else {
				// Cursor at end
				padding := valueAreaWidth - valueDisplayWidth - 1
				if padding < 0 {
					padding = 0
				}
				valueDisplay = value + ui.CursorNarrow.Render(" ") + strings.Repeat(" ", padding)
			}
		} else {
			// Not focused: plain value with padding
			padding := valueAreaWidth - valueDisplayWidth
			if padding < 0 {
				padding = 0
			}
			valueDisplay = value + strings.Repeat(" ", padding)
		}

		line.WriteString(labelStyle.Render(labelPart))
		line.WriteString(valueDisplay)
		line.WriteString(" ")
		line.WriteString(borderStyle.Render("│"))
		return line.String()
	}

	// Empty line
	dialog.WriteString(borderStyle.Render("│"))
	dialog.WriteString(strings.Repeat(" ", dialogWidth-2))
	dialog.WriteString(borderStyle.Render("│"))
	dialog.WriteString("\n")

	// Endpoint field
	dialog.WriteString(renderFieldLine("Endpoint", m.EditEndpoint, 0, m.EditCursorPos))
	dialog.WriteString("\n")

	// Port field
	dialog.WriteString(renderFieldLine("Port", m.EditPort, 1, m.EditCursorPos))
	dialog.WriteString("\n")

	// Empty line
	dialog.WriteString(borderStyle.Render("│"))
	dialog.WriteString(strings.Repeat(" ", dialogWidth-2))
	dialog.WriteString(borderStyle.Render("│"))
	dialog.WriteString("\n")

	// Help text
	helpText := "Connect: <enter> | Close: esc"
	helpDisplayWidth := lipgloss.Width(helpText)
	helpPadding := contentWidth - helpDisplayWidth
	if helpPadding < 0 {
		helpPadding = 0
	}
	dialog.WriteString(borderStyle.Render("│"))
	dialog.WriteString(" ")
	dialog.WriteString(ui.StyleHelpText.Render(helpText))
	dialog.WriteString(strings.Repeat(" ", helpPadding))
	dialog.WriteString(" ")
	dialog.WriteString(borderStyle.Render("│"))
	dialog.WriteString("\n")

	// Bottom border
	dialog.WriteString(borderStyle.Render("╰"))
	dialog.WriteString(borderStyle.Render(strings.Repeat("─", dialogWidth-2)))
	dialog.WriteString(borderStyle.Render("╯"))

	// Center the dialog
	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		dialog.String(),
	)
}
