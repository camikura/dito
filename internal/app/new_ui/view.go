package new_ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/camikura/dito/internal/db"
	"github.com/camikura/dito/internal/ui"
)

// Color definitions - using ui package constants
const (
	ColorPrimary   = ui.ColorPrimaryHex
	ColorInactive  = ui.ColorInactiveHex
	ColorGreen     = ui.ColorGreenHex
	ColorLabel     = ui.ColorLabelHex
	ColorSecondary = ui.ColorSecondaryHex
	ColorTertiary  = ui.ColorTertiaryHex
	ColorPK        = ui.ColorPKHex
	ColorIndex     = ui.ColorIndexHex
	ColorHelp      = ui.ColorHelpHex
)

// RenderView renders the new UI
func RenderView(m Model) string {
	if m.Width == 0 {
		return "Loading..."
	}

	// Layout configuration
	leftPaneContentWidth := ui.LeftPaneContentWidth
	leftPaneActualWidth := leftPaneContentWidth + ui.LeftPaneBorderWidth
	rightPaneActualWidth := m.Width - leftPaneActualWidth + 1 // +1 to use full width

	// Render connection pane first to get its actual height
	connectionPane := renderConnectionPane(m, leftPaneContentWidth)
	connectionPaneHeight := strings.Count(connectionPane, "\n") + 1 // Count actual lines

	// Calculate pane heights based on actual connection pane height
	// This ensures heights are always correct even if connection pane height changes
	availableHeight := m.Height - 1 - connectionPaneHeight - 6
	partHeight := availableHeight / ui.PaneHeightTotalParts
	remainder := availableHeight % ui.PaneHeightTotalParts

	tablesHeight := partHeight * ui.PaneHeightTablesParts
	schemaHeight := partHeight * ui.PaneHeightSchemaParts
	sqlHeight := partHeight * ui.PaneHeightSQLParts

	// Distribute remainder
	ui.DistributeSpace(remainder, &tablesHeight, &schemaHeight, &sqlHeight)

	// Ensure minimum heights
	if tablesHeight < 3 {
		tablesHeight = 3
	}
	if schemaHeight < 3 {
		schemaHeight = 3
	}
	if sqlHeight < 2 {
		sqlHeight = 2
	}

	// After applying minimum heights, redistribute unused space
	usedHeight := tablesHeight + schemaHeight + sqlHeight
	if usedHeight < availableHeight {
		ui.DistributeSpace(availableHeight-usedHeight, &tablesHeight, &schemaHeight, &sqlHeight)
	}

	// Render remaining panes with calculated heights
	tablesPane := renderTablesPaneWithHeight(m, leftPaneContentWidth, tablesHeight)
	schemaPane := renderSchemaPaneWithHeight(m, leftPaneContentWidth, schemaHeight)
	sqlPane := renderSQLPaneWithHeight(m, leftPaneContentWidth, sqlHeight)
	dataPane := renderDataPane(m, rightPaneActualWidth, m.Height)

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
		if m.Connected {
			footerContent = "Switch Pane: tab | Disconnect: ctrl+d"
		} else {
			footerContent = "Switch Pane: tab | Connect: <enter>"
		}
	case FocusPaneTables:
		footerContent = "Navigate: ↑/↓ | Switch Pane: tab | Select: <enter>"
	case FocusPaneSQL:
		footerContent = "Switch Pane: tab | Edit: <enter>"
	case FocusPaneData:
		if m.CustomSQL {
			footerContent = "Navigate: ↑/↓ | Switch Pane: tab | Detail: <enter> | Reset: esc"
		} else {
			footerContent = "Navigate: ↑/↓ | Switch Pane: tab | Detail: <enter>"
		}
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

	baseView := result.String()

	// Overlay connection dialog if visible
	if m.ConnectionDialogVisible {
		return renderConnectionDialog(m)
	}

	// Overlay SQL editor dialog if visible
	if m.SQLEditorVisible {
		return renderSQLEditorDialog(m)
	}

	// Overlay record detail dialog if visible
	if m.RecordDetailVisible {
		return renderRecordDetailDialog(m)
	}

	return baseView
}

func renderConnectionPane(m Model, width int) string {
	borderColor := ColorInactive
	titleColor := ColorInactive
	if m.CurrentPane == FocusPaneConnection {
		borderColor = ColorPrimary
		titleColor = ColorPrimary
	}

	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(borderColor))
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(titleColor))

	var titleText string
	var titleDisplayWidth int
	if m.Connected {
		// Apply green color to checkmark
		checkmark := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorGreen)).Render("✓")
		titleText = titleStyle.Render(" Connection ") + checkmark + " "
		// " Connection ✓ " = 1 + 10 + 1 + 1 + 1 = 14 display chars
		titleDisplayWidth = 14
	} else {
		titleText = titleStyle.Render(" Connection ") + " "
		// " Connection " = 1 + 10 + 1 = 12 display chars
		titleDisplayWidth = 12
	}

	title := borderStyle.Render("╭─") + titleText + borderStyle.Render(strings.Repeat("─", width-titleDisplayWidth-3)+"╮")

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

	// Pad content to width (no left/right padding)
	paddingLen := width - len(content) - 2
	if paddingLen < 0 {
		paddingLen = 0
	}
	contentPadded := content + strings.Repeat(" ", paddingLen)

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
	return renderTablesPaneWithHeight(m, width, 12)
}

func renderTablesPaneWithHeight(m Model, width int, height int) string {
	borderColor := ColorInactive
	titleColor := ColorInactive
	if m.CurrentPane == FocusPaneTables {
		borderColor = ColorPrimary
		titleColor = ColorPrimary
	}

	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(borderColor))
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(titleColor))

	titleText := " Tables"
	if len(m.Tables) > 0 {
		titleText += " (" + string(rune(len(m.Tables)+48)) + ")"
	}
	titleText += " "

	styledTitle := titleStyle.Render(titleText)
	title := borderStyle.Render("╭─") + styledTitle + borderStyle.Render(strings.Repeat("─", width-len(titleText)-3) + "╮")

	// Prepare content lines with tree structure
	type tableLineInfo struct {
		text       string
		isSelected bool
	}
	var contentLines []tableLineInfo
	if len(m.Tables) == 0 {
		contentLines = []tableLineInfo{{text: "No tables", isSelected: false}}
	} else {
		// Render each table with tree structure
		for i, tableName := range m.Tables {
			// Determine indentation based on '.' separator
			indent := ""
			displayName := tableName
			if dotIndex := strings.LastIndex(tableName, "."); dotIndex != -1 {
				// Child table - indent and show only the child part
				indent = " "
				displayName = tableName[dotIndex+1:]
			}

			// Add selection marker (* for selected table, always visible)
			var prefix string
			isSelected := i == m.SelectedTable
			if isSelected {
				prefix = "* "
			} else {
				prefix = "  "
			}

			contentLines = append(contentLines, tableLineInfo{
				text:       prefix + indent + displayName,
				isSelected: isSelected,
			})
		}
	}

	leftBorder := borderStyle.Render("│")
	rightBorder := borderStyle.Render("│")
	bottomBorder := borderStyle.Render("╰" + strings.Repeat("─", width-2) + "╯")

	// Styles for text color
	selectedTextStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")) // White for selected
	unselectedTextStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")) // Gray for unselected

	var result strings.Builder
	result.WriteString(title + "\n")

	// Render content lines (fill allocated height with content or empty lines)
	for i := 0; i < height; i++ {
		contentIndex := i + m.TablesScrollOffset
		if contentIndex < len(contentLines) {
			lineInfo := contentLines[contentIndex]
			// Apply color based on selection
			var styledText string
			if lineInfo.isSelected {
				styledText = selectedTextStyle.Render(lineInfo.text)
			} else {
				styledText = unselectedTextStyle.Render(lineInfo.text)
			}
			// Calculate padding (based on original text length, not styled)
			paddingLen := width - len(lineInfo.text) - 2
			if paddingLen < 0 {
				paddingLen = 0
			}
			line := styledText + strings.Repeat(" ", paddingLen)
			result.WriteString(leftBorder + line + rightBorder + "\n")
		} else {
			// Empty line for remaining allocated height
			emptyLine := strings.Repeat(" ", width-2)
			result.WriteString(leftBorder + emptyLine + rightBorder + "\n")
		}
	}

	result.WriteString(bottomBorder)

	return result.String()
}

func renderSchemaPane(m Model, width int) string {
	return renderSchemaPaneWithHeight(m, width, 12)
}

func renderSchemaPaneWithHeight(m Model, width int, height int) string {
	// Title includes table name if available
	titleText := " Schema "
	if len(m.Tables) > 0 && m.CursorTable < len(m.Tables) {
		tableName := m.Tables[m.CursorTable]
		titleText = " Schema (" + tableName + ") "
	}

	// Schema pane can be focused for scrolling
	var borderStyle lipgloss.Style
	if m.CurrentPane == FocusPaneSchema {
		borderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorPrimary))
	} else {
		borderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorInactive))
	}
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
				primaryKeys := ui.ParsePrimaryKeysFromDDL(details.Schema.DDL)
				columns := ui.ParseColumnsFromDDL(details.Schema.DDL, primaryKeys)

				// First pass: find the longest column name
				maxColNameLen := 0
				for _, col := range columns {
					if len(col.Name) > maxColNameLen {
						maxColNameLen = len(col.Name)
					}
				}

				// Format each column: PK|||Name|||Type|||maxLen (use ||| as separator)
				for _, col := range columns {
					pkMarker := " " // Single space when not PK
					if col.IsPrimaryKey {
						pkMarker = "P" // Single "P" for primary key
					}
					contentLines = append(contentLines, fmt.Sprintf("%s|||%s|||%s|||%d", pkMarker, col.Name, col.Type, maxColNameLen))
				}
			}

			// Add indexes section
			contentLines = append(contentLines, "")
			contentLines = append(contentLines, "Indexes:")
			if len(details.Indexes) > 0 {
				for _, index := range details.Indexes {
					fields := strings.Join(index.FieldNames, ", ")
					// Format: IndexName|||Fields (use ||| as separator to apply color in rendering)
					contentLines = append(contentLines, fmt.Sprintf("IDX|||%s|||%s", index.IndexName, fields))
				}
			} else {
				contentLines = append(contentLines, "  (none)")
			}
		}
	}

	leftBorder := borderStyle.Render("│")
	rightBorder := borderStyle.Render("│")
	bottomBorder := borderStyle.Render("╰" + strings.Repeat("─", width-2) + "╯")

	// Styles for rendering
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSecondary))
	typeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorTertiary))
	pkStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorPK))
	indexStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorIndex))

	var result strings.Builder
	result.WriteString(title + "\n")

	// Render content lines (fill allocated height with content or empty lines)
	for i := 0; i < height; i++ {
		var line string
		contentIndex := i + m.SchemaScrollOffset
		if contentIndex < len(contentLines) {
			content := contentLines[contentIndex]

			// Apply yellow color to label lines (Columns:, Indexes:)
			if content == "Columns:" || content == "Indexes:" {
				paddingLen := width - len(content) - 2
				if paddingLen < 0 {
					paddingLen = 0
				}
				line = labelStyle.Render(content) + strings.Repeat(" ", paddingLen)
			} else if strings.HasPrefix(content, "IDX|||") {
			// Index line: IDX|||IndexName|||Fields
			parts := strings.Split(content, "|||")
			if len(parts) >= 3 {
				indexName := parts[1]
				fields := parts[2]

				// Format: "  indexName fields" with field names in index color and commas in white
				var fieldsDisplay string
				if strings.Contains(fields, ", ") {
					// Multiple fields: color each field name separately, keep commas white
					fieldList := strings.Split(fields, ", ")
					for i, field := range fieldList {
						if i > 0 {
							fieldsDisplay += ", " // White comma
						}
						fieldsDisplay += indexStyle.Render(field)
					}
				} else {
					// Single field
					fieldsDisplay = indexStyle.Render(fields)
				}

				displayText := "  " + indexName + " " + fieldsDisplay
				displayLen := 2 + len(indexName) + 1 + len(fields)

				availableWidth := width - 2
				rightPadding := availableWidth - displayLen
				if rightPadding < 0 {
					rightPadding = 0
				}
				line = displayText + strings.Repeat(" ", rightPadding)
			} else {
				// Fallback
				paddingLen := width - len(content) - 2
				if paddingLen < 0 {
					paddingLen = 0
				}
				line = content + strings.Repeat(" ", paddingLen)
			}
		} else if strings.Contains(content, "|||") {
			// Column line with PK, name, type, and maxColNameLen separated by |||
			parts := strings.Split(content, "|||")
			if len(parts) >= 4 {
				pkMarker := parts[0]  // "P" or " "
				colName := parts[1]
				colType := parts[2]
				maxColNameLen, _ := strconv.Atoi(parts[3])

				// Fixed column widths for alignment
				const pkColWidth = 2               // Fixed width for PK marker (1 char + 1 space)
				nameColWidth := maxColNameLen + 1  // Use actual max column name length + 1 space

				// PK marker with fixed width
				var pkField string
				if pkMarker == "P" {
					pkField = pkStyle.Render(pkMarker) + " "
				} else {
					pkField = strings.Repeat(" ", pkColWidth)
				}

				// Pad column name to fixed width
				namePadding := nameColWidth - len(colName)
				if namePadding < 0 {
					namePadding = 0
				}
				nameField := colName + strings.Repeat(" ", namePadding)

				// Type field (no fixed width, left-aligned)
				typeField := typeStyle.Render(colType)

				// Build line with fixed-width columns: PK + Name + Type
				alignedLine := pkField + nameField + typeField

				// Calculate right padding
				displayLen := pkColWidth + nameColWidth + len(colType)
				availableWidth := width - 2 // -2 for borders
				rightPadding := availableWidth - displayLen
				if rightPadding < 0 {
					rightPadding = 0
				}
				line = alignedLine + strings.Repeat(" ", rightPadding)
			} else {
				// Fallback if parsing fails
				paddingLen := width - len(content) - 2
				if paddingLen < 0 {
					paddingLen = 0
				}
				line = content + strings.Repeat(" ", paddingLen)
			}
		} else {
			// Other content (like "Select a table", "Loading...")
			if len(content) > width-2 {
				content = content[:width-5] + "..."
			}
			paddingLen := width - len(content) - 2
			if paddingLen < 0 {
				paddingLen = 0
			}
			line = content + strings.Repeat(" ", paddingLen)
			}
			result.WriteString(leftBorder + line + rightBorder + "\n")
		} else {
			// Empty line for remaining allocated height
			emptyLine := strings.Repeat(" ", width-2)
			result.WriteString(leftBorder + emptyLine + rightBorder + "\n")
		}
	}

	result.WriteString(bottomBorder)

	return result.String()
}

func renderSQLPane(m Model, width int) string {
	return renderSQLPaneWithHeight(m, width, 6)
}

func renderSQLPaneWithHeight(m Model, width int, height int) string {
	borderColor := ColorInactive
	titleColor := ColorInactive
	if m.CurrentPane == FocusPaneSQL {
		borderColor = ColorPrimary
		titleColor = ColorPrimary
	}

	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(borderColor))
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(titleColor))

	// Add [Custom] label if custom SQL is active
	titleText := " SQL "
	if m.CustomSQL {
		titleText = " SQL [Custom] "
	}
	styledTitle := titleStyle.Render(titleText)
	title := borderStyle.Render("╭─") + styledTitle + borderStyle.Render(strings.Repeat("─", width-len(titleText)-3) + "╮")

	content := ""
	if m.CurrentSQL != "" {
		content = m.CurrentSQL
		// Truncate if too long
		if len(content) > width-2 {
			content = content[:width-5] + "..."
		}
	}

	leftBorder := borderStyle.Render("│")
	rightBorder := borderStyle.Render("│")
	bottomBorder := borderStyle.Render("╰" + strings.Repeat("─", width-2) + "╯")

	// Calculate padding, ensuring it's not negative (no left/right padding)
	paddingLen := width - len(content) - 2
	if paddingLen < 0 {
		paddingLen = 0
	}
	contentPadded := content + strings.Repeat(" ", paddingLen)

	var result strings.Builder
	result.WriteString(title + "\n")

	// Render SQL content with dynamic height (fill with empty lines to use allocated space)
	for i := 0; i < height; i++ {
		if i == 0 {
			result.WriteString(leftBorder + contentPadded + rightBorder + "\n")
		} else {
			emptyLine := strings.Repeat(" ", width-2)
			result.WriteString(leftBorder + emptyLine + rightBorder + "\n")
		}
	}
	result.WriteString(bottomBorder)

	return result.String()
}

func renderDataPane(m Model, width int, totalHeight int) string {
	borderColor := ColorInactive
	titleColor := ColorInactive
	if m.CurrentPane == FocusPaneData {
		borderColor = ColorPrimary
		titleColor = ColorPrimary
	}

	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(borderColor))
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(titleColor))

	// Build title with table name if available
	var titleText string
	if m.SelectedTable >= 0 && m.SelectedTable < len(m.Tables) {
		tableName := m.Tables[m.SelectedTable]
		titleText = fmt.Sprintf(" Data (%s) ", tableName)
	} else {
		titleText = " Data "
	}
	styledTitle := titleStyle.Render(titleText)
	title := borderStyle.Render("╭─") + styledTitle + borderStyle.Render(strings.Repeat("─", width-len(titleText)-3) + "╮")

	// Data pane should match left panes height
	// Total height = totalHeight (m.Height passed in)
	// Data pane takes: totalHeight - 1 (footer)
	// Data pane structure: title(1) + content lines(?) + bottom(1)
	// So: content lines = totalHeight - 1 - 2 = totalHeight - 3
	contentLines := totalHeight - 3
	if contentLines < 5 {
		contentLines = 5
	}

	leftBorder := borderStyle.Render("│")
	rightBorder := borderStyle.Render("│")

	var result strings.Builder
	result.WriteString(title + "\n")

	// Track scroll info for bottom border
	var totalContentWidth, viewportWidth int

	// Prepare content
	if m.SelectedTable < 0 || m.SelectedTable >= len(m.Tables) {
		// No table selected
		for i := 0; i < contentLines; i++ {
			line := "No data"
			if i > 0 {
				line = ""
			}
			paddingLen := width - len(line) - 2
			if paddingLen < 0 {
				paddingLen = 0
			}
			result.WriteString(leftBorder + line + strings.Repeat(" ", paddingLen) + rightBorder + "\n")
		}
	} else {
		tableName := m.Tables[m.SelectedTable]
		data, exists := m.TableData[tableName]
		if !exists || data == nil {
			// No data loaded yet
			message := "No data"
			if m.LoadingData {
				message = "Loading..."
			}
			for i := 0; i < contentLines; i++ {
				line := ""
				if i == 0 {
					line = message
				}
				paddingLen := width - len(line) - 2
				if paddingLen < 0 {
					paddingLen = 0
				}
				result.WriteString(leftBorder + line + strings.Repeat(" ", paddingLen) + rightBorder + "\n")
			}
		} else if len(data.Rows) == 0 {
			// No rows in result
			for i := 0; i < contentLines; i++ {
				line := "No rows"
				if i > 0 {
					line = ""
				}
				paddingLen := width - len(line) - 2
				if paddingLen < 0 {
					paddingLen = 0
				}
				result.WriteString(leftBorder + line + strings.Repeat(" ", paddingLen) + rightBorder + "\n")
			}
		} else {
			// Render grid view and get scroll info
			gridContent, scrollInfo := renderGridViewWithScrollInfo(m, data, width, contentLines, leftBorder, borderStyle)
			result.WriteString(gridContent)
			totalContentWidth = scrollInfo.totalWidth
			viewportWidth = scrollInfo.viewportWidth
		}
	}

	// Render bottom border with scrollbar
	bottomBorder := renderBottomBorderWithScrollbar(borderStyle, width, totalContentWidth, viewportWidth, m.HorizontalOffset)
	result.WriteString(bottomBorder)

	return result.String()
}

// scrollInfo holds scroll-related information for the grid
type scrollInfo struct {
	totalWidth     int
	viewportWidth  int
	totalRows      int
	viewportRows   int
	verticalOffset int
}

// renderBottomBorderWithScrollbar renders the bottom border with an integrated scrollbar
func renderBottomBorderWithScrollbar(borderStyle lipgloss.Style, width int, totalContentWidth int, viewportWidth int, offset int) string {
	// Border structure: ╰ + content + ╯
	// Content width = width - 2 (for ╰ and ╯)
	contentWidth := width - 2
	if contentWidth < 1 {
		return borderStyle.Render("╰╯")
	}

	// Create scrollbar
	scrollbar := ui.NewScrollBar(totalContentWidth, viewportWidth, offset, contentWidth)
	scrollbarLine := scrollbar.Render()

	return borderStyle.Render("╰") + borderStyle.Render(scrollbarLine) + borderStyle.Render("╯")
}

// renderGridViewWithScrollInfo renders the data grid and returns scroll information
func renderGridViewWithScrollInfo(m Model, data *db.TableDataResult, width int, contentLines int, leftBorder string, borderStyle lipgloss.Style) (string, scrollInfo) {
	// Get column names in schema definition order
	tableName := m.Tables[m.SelectedTable]
	columns := getColumnsInSchemaOrder(m, tableName, data.Rows)

	// Get column types from schema
	columnTypes := getColumnTypes(m, tableName, columns)

	// Create Grid component
	contentWidth := width - 2 // Subtract borders
	grid := ui.NewGrid(columns, columnTypes, data.Rows)
	grid.Width = contentWidth
	grid.Height = contentLines
	grid.HorizontalOffset = m.HorizontalOffset
	grid.VerticalOffset = m.ViewportOffset
	grid.SelectedRow = m.SelectedDataRow

	// Calculate viewport rows (contentLines minus header and separator)
	viewportRows := contentLines - 2
	if viewportRows < 1 {
		viewportRows = 1
	}

	// Get scroll info before rendering
	info := scrollInfo{
		totalWidth:     grid.TotalContentWidth(),
		viewportWidth:  contentWidth,
		totalRows:      len(data.Rows),
		viewportRows:   viewportRows,
		verticalOffset: m.ViewportOffset,
	}

	// Create vertical scrollbar
	vScrollBar := ui.NewVerticalScrollBar(info.totalRows, info.viewportRows, info.verticalOffset, contentLines)

	// Render grid content
	gridContent := grid.Render()

	// Add borders to each line with vertical scrollbar on right
	var result strings.Builder
	lines := strings.Split(gridContent, "\n")
	for i := 0; i < contentLines; i++ {
		// Get right border character (with scrollbar indicator)
		rightBorderChar := vScrollBar.GetCharAt(i)
		rightBorder := borderStyle.Render(rightBorderChar)

		if i < len(lines) {
			result.WriteString(leftBorder + lines[i] + rightBorder + "\n")
		} else {
			// Empty line to fill height
			emptyLine := strings.Repeat(" ", contentWidth)
			result.WriteString(leftBorder + emptyLine + rightBorder + "\n")
		}
	}

	return result.String(), info
}


// getColumnTypes extracts column types from schema information
func getColumnTypes(m Model, tableName string, columns []string) map[string]string {
	// Get DDL from table details
	var ddl string
	if details, exists := m.TableDetails[tableName]; exists && details != nil && details.Schema != nil {
		ddl = details.Schema.DDL
	}

	return ui.GetColumnTypes(ddl)
}

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
		BorderColor:  ColorPrimary,
	})

	return rd.RenderCentered(m.Width, m.Height)
}

// renderSQLEditorDialog renders the SQL editor dialog
func renderSQLEditorDialog(m Model) string {
	// Calculate dialog size
	dialogWidth := m.Width - 10
	if dialogWidth < 60 {
		dialogWidth = 60
	}
	if dialogWidth > 100 {
		dialogWidth = 100
	}

	dialogHeight := m.Height - 10
	if dialogHeight < 10 {
		dialogHeight = 10
	}
	if dialogHeight > 20 {
		dialogHeight = 20
	}

	// Calculate content dimensions
	contentWidth := dialogWidth - 4  // borders + padding
	contentHeight := dialogHeight - 6 // title, help, borders, padding

	// Border style
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorPrimary))

	// Build dialog
	var dialog strings.Builder

	// Title
	titleText := " SQL Editor "
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorPrimary)).Bold(true)
	title := titleStyle.Render(titleText)
	titleLen := len([]rune(titleText))

	// Top border: ╭ + title + ─...─ + ╮
	dashesLen := dialogWidth - 2 - titleLen
	if dashesLen < 0 {
		dashesLen = 0
	}
	dialog.WriteString(borderStyle.Render("╭"))
	dialog.WriteString(title)
	dialog.WriteString(borderStyle.Render(strings.Repeat("─", dashesLen)))
	dialog.WriteString(borderStyle.Render("╮"))
	dialog.WriteString("\n")

	// Parse SQL into lines and calculate cursor position
	sqlLines := strings.Split(m.EditSQL, "\n")
	cursorLine := 0
	cursorCol := m.SQLCursorPos
	currentPos := 0
	for i, line := range sqlLines {
		lineLen := len(line) + 1 // +1 for newline
		if currentPos+lineLen > m.SQLCursorPos {
			cursorLine = i
			cursorCol = m.SQLCursorPos - currentPos
			break
		}
		currentPos += lineLen
	}

	// Render SQL content with cursor
	for i := 0; i < contentHeight; i++ {
		dialog.WriteString(borderStyle.Render("│"))
		dialog.WriteString(" ")

		var lineContent string
		if i < len(sqlLines) {
			line := sqlLines[i]
			if i == cursorLine {
				// Insert cursor indicator
				if cursorCol <= len(line) {
					lineContent = line[:cursorCol] + "_" + line[cursorCol:]
				} else {
					lineContent = line + "_"
				}
			} else {
				lineContent = line
			}
		}

		// Pad or truncate to fit content width
		visibleLen := len([]rune(lineContent))
		if visibleLen < contentWidth {
			lineContent += strings.Repeat(" ", contentWidth-visibleLen)
		} else if visibleLen > contentWidth {
			lineContent = string([]rune(lineContent)[:contentWidth])
		}

		dialog.WriteString(lineContent)
		dialog.WriteString(" ")
		dialog.WriteString(borderStyle.Render("│"))
		dialog.WriteString("\n")
	}

	// Empty line before help
	dialog.WriteString(borderStyle.Render("│"))
	dialog.WriteString(strings.Repeat(" ", dialogWidth-2))
	dialog.WriteString(borderStyle.Render("│"))
	dialog.WriteString("\n")

	// Help text
	helpText := "Execute: ctrl+r | Cancel: esc"
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorHelp))
	helpPadding := dialogWidth - 4 - len(helpText)
	if helpPadding < 0 {
		helpPadding = 0
	}
	dialog.WriteString(borderStyle.Render("│"))
	dialog.WriteString(" ")
	dialog.WriteString(helpStyle.Render(helpText))
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

// renderConnectionDialog renders the connection setup dialog
func renderConnectionDialog(m Model) string {
	dialogWidth := 50

	// Border style
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorPrimary))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorLabel))
	selectedBgStyle := lipgloss.NewStyle().Background(lipgloss.Color(ColorPrimary)).Foreground(lipgloss.Color("#000000"))
	cursorStyle := lipgloss.NewStyle().Background(lipgloss.Color("#ffffff")).Foreground(lipgloss.Color("#000000"))

	var dialog strings.Builder

	// Title
	titleText := " Connection Setup "
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorPrimary)).Bold(true)
	title := titleStyle.Render(titleText)
	titleLen := len([]rune(titleText))

	dashesLen := dialogWidth - 2 - titleLen
	if dashesLen < 0 {
		dashesLen = 0
	}
	dialog.WriteString(borderStyle.Render("╭"))
	dialog.WriteString(title)
	dialog.WriteString(borderStyle.Render(strings.Repeat("─", dashesLen)))
	dialog.WriteString(borderStyle.Render("╮"))
	dialog.WriteString("\n")

	// Content width
	contentWidth := dialogWidth - 4

	// Helper function to render a field
	renderField := func(label string, value string, fieldIndex int, isEditing bool, cursorPos int) string {
		var line strings.Builder
		line.WriteString(borderStyle.Render("│"))
		line.WriteString(" ")

		// Marker for current field
		marker := "  "
		if m.ConnectionDialogField == fieldIndex {
			marker = "* "
		}

		labelText := labelStyle.Render(marker + label + ": ")
		labelLen := len(marker) + len(label) + 2

		valueWidth := contentWidth - labelLen
		if valueWidth < 10 {
			valueWidth = 10
		}

		var valueText string
		if m.ConnectionDialogField == fieldIndex && isEditing {
			// Show cursor in editing mode
			if cursorPos < len(value) {
				valueText = value[:cursorPos] + cursorStyle.Render(string(value[cursorPos])) + value[cursorPos+1:]
			} else {
				valueText = value + cursorStyle.Render(" ")
			}
		} else if m.ConnectionDialogField == fieldIndex {
			// Highlight selected field
			valueText = selectedBgStyle.Render(value + strings.Repeat(" ", valueWidth-len(value)))
		} else {
			valueText = value
		}

		// Pad to width
		displayLen := len(value)
		if m.ConnectionDialogField == fieldIndex && !isEditing {
			displayLen = valueWidth
		}
		padding := contentWidth - labelLen - displayLen
		if padding < 0 {
			padding = 0
		}

		line.WriteString(labelText)
		line.WriteString(valueText)
		if m.ConnectionDialogField != fieldIndex || isEditing {
			line.WriteString(strings.Repeat(" ", padding))
		}
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
	dialog.WriteString(renderField("Endpoint", m.EditEndpoint, 0, m.ConnectionDialogEditing, m.EditCursorPos))
	dialog.WriteString("\n")

	// Port field
	dialog.WriteString(renderField("Port", m.EditPort, 1, m.ConnectionDialogEditing, m.EditCursorPos))
	dialog.WriteString("\n")

	// Empty line
	dialog.WriteString(borderStyle.Render("│"))
	dialog.WriteString(strings.Repeat(" ", dialogWidth-2))
	dialog.WriteString(borderStyle.Render("│"))
	dialog.WriteString("\n")

	// Connect button
	dialog.WriteString(borderStyle.Render("│"))
	dialog.WriteString(" ")
	buttonText := "[ Connect ]"
	buttonPadding := (contentWidth - len(buttonText)) / 2
	if m.ConnectionDialogField == 2 {
		dialog.WriteString(strings.Repeat(" ", buttonPadding))
		dialog.WriteString(selectedBgStyle.Render(buttonText))
		dialog.WriteString(strings.Repeat(" ", contentWidth-buttonPadding-len(buttonText)))
	} else {
		dialog.WriteString(strings.Repeat(" ", buttonPadding))
		dialog.WriteString(buttonText)
		dialog.WriteString(strings.Repeat(" ", contentWidth-buttonPadding-len(buttonText)))
	}
	dialog.WriteString(" ")
	dialog.WriteString(borderStyle.Render("│"))
	dialog.WriteString("\n")

	// Empty line
	dialog.WriteString(borderStyle.Render("│"))
	dialog.WriteString(strings.Repeat(" ", dialogWidth-2))
	dialog.WriteString(borderStyle.Render("│"))
	dialog.WriteString("\n")

	// Help text
	var helpText string
	if m.ConnectionDialogEditing {
		helpText = "Confirm: enter | Cancel: esc"
	} else {
		helpText = "Edit: enter | Navigate: ↑/↓ | Close: esc"
	}
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorHelp))
	helpPadding := dialogWidth - 4 - len(helpText)
	if helpPadding < 0 {
		helpPadding = 0
	}
	dialog.WriteString(borderStyle.Render("│"))
	dialog.WriteString(" ")
	dialog.WriteString(helpStyle.Render(helpText))
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
