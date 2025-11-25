package new_ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/camikura/dito/internal/db"
	"github.com/camikura/dito/internal/views"
)

// Color definitions
const (
	ColorPrimary   = "#00D9FF" // Cyan for active borders
	ColorInactive  = "#AAAAAA" // Light gray for inactive borders
	ColorGreen     = "#00FF00" // Green for connection status
	ColorLabel     = "#00D9FF" // Cyan for section labels
	ColorSecondary = "#C47D7D" // Muted reddish for schema section labels (Columns:, Indexes:)
	ColorTertiary  = "#7AA2F7" // Soft blue for data types
	ColorPK        = "#7FBA7A" // Muted green for primary key marker
	ColorIndex     = "#E5C07B" // Warm yellow/beige for index field names
	ColorHelp      = "#888888" // Gray for help text
)

// RenderView renders the new UI
func RenderView(m Model) string {
	if m.Width == 0 {
		return "Loading..."
	}

	// Layout configuration
	leftPaneContentWidth := 50 // Content width for left panes
	leftPaneActualWidth := leftPaneContentWidth + 2 // +2 for borders
	rightPaneActualWidth := m.Width - leftPaneActualWidth + 1 // +1 to use full width

	// Render connection pane first to get its actual height
	connectionPane := renderConnectionPane(m, leftPaneContentWidth)
	connectionPaneHeight := strings.Count(connectionPane, "\n") + 1 // Count actual lines

	// Calculate pane heights based on actual connection pane height
	// This ensures heights are always correct even if connection pane height changes
	availableHeight := m.Height - 1 - connectionPaneHeight - 6
	totalParts := 5
	partHeight := availableHeight / totalParts
	remainder := availableHeight % totalParts

	tablesHeight := partHeight * 2
	schemaHeight := partHeight * 2
	sqlHeight := partHeight

	// Distribute remainder (may be up to 4)
	for remainder > 0 {
		if remainder >= 1 {
			tablesHeight++
			remainder--
		}
		if remainder >= 1 {
			schemaHeight++
			remainder--
		}
		if remainder >= 1 {
			sqlHeight++
			remainder--
		}
	}

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

	// After applying minimum heights, check if we have unused space
	usedHeight := tablesHeight + schemaHeight + sqlHeight
	if usedHeight < availableHeight {
		// Distribute unused space in 2:2:1 ratio again
		extraSpace := availableHeight - usedHeight
		for extraSpace > 0 {
			if extraSpace >= 1 {
				tablesHeight++
				extraSpace--
			}
			if extraSpace >= 1 {
				schemaHeight++
				extraSpace--
			}
			if extraSpace >= 1 {
				sqlHeight++
				extraSpace--
			}
		}
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
				primaryKeys := views.ParsePrimaryKeysFromDDL(details.Schema.DDL)
				columns := views.ParseColumnsFromDDL(details.Schema.DDL, primaryKeys)

				// Fixed width for data type column (right-aligned)
				// Longest type is TIMESTAMP(9) = 12 chars
				const typeColumnWidth = 12

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

	titleText := " SQL "
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

	titleText := " Data "
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
	bottomBorder := borderStyle.Render("╰" + strings.Repeat("─", width-2) + "╯")

	var result strings.Builder
	result.WriteString(title + "\n")

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
			// Render grid view
			result.WriteString(renderGridView(m, data, width, contentLines, leftBorder, rightBorder))
		}
	}

	result.WriteString(bottomBorder)

	return result.String()
}

// renderGridView renders the data grid with headers, separator, and rows
func renderGridView(m Model, data *db.TableDataResult, width int, contentLines int, leftBorder, rightBorder string) string {
	var result strings.Builder

	// Extract column names from first row
	var columns []string
	if len(data.Rows) > 0 {
		for col := range data.Rows[0] {
			columns = append(columns, col)
		}
		// Sort columns for consistent display
		// For better UX, we should order by schema, but for now just sort alphabetically
		sortColumns(columns)
	}

	// Get column types from schema
	tableName := m.Tables[m.SelectedTable]
	columnTypes := getColumnTypes(m, tableName, columns)

	// Calculate column widths
	contentWidth := width - 2 // Subtract borders
	columnWidths := calculateColumnWidths(columns, data.Rows, contentWidth)

	// Line 1: Header row
	headerLine := renderHeaderRow(columns, columnWidths, contentWidth)
	result.WriteString(leftBorder + headerLine + rightBorder + "\n")

	// Line 2: Separator (── ───── ────)
	separatorLine := renderSeparator(columnWidths, contentWidth)
	result.WriteString(leftBorder + separatorLine + rightBorder + "\n")

	// Lines 3+: Data rows
	dataLinesAvailable := contentLines - 2 // Subtract header and separator
	for i := 0; i < dataLinesAvailable; i++ {
		rowIndex := i + m.ViewportOffset
		if rowIndex < len(data.Rows) {
			isSelected := rowIndex == m.SelectedDataRow
			rowLine := renderDataRow(columns, data.Rows[rowIndex], columnWidths, columnTypes, contentWidth, isSelected)
			result.WriteString(leftBorder + rowLine + rightBorder + "\n")
		} else {
			// Empty line
			emptyLine := strings.Repeat(" ", contentWidth)
			result.WriteString(leftBorder + emptyLine + rightBorder + "\n")
		}
	}

	return result.String()
}

// getColumnTypes extracts column types from schema information
func getColumnTypes(m Model, tableName string, columns []string) map[string]string {
	types := make(map[string]string)

	if details, exists := m.TableDetails[tableName]; exists && details != nil && details.Schema != nil {
		if details.Schema.DDL != "" {
			// Parse column types from DDL
			primaryKeys := views.ParsePrimaryKeysFromDDL(details.Schema.DDL)
			cols := views.ParseColumnsFromDDL(details.Schema.DDL, primaryKeys)
			for _, col := range cols {
				// Remove " (Primary Key)" suffix if present
				colType := col.Type
				if idx := strings.Index(colType, " (Primary Key)"); idx != -1 {
					colType = colType[:idx]
				}
				types[col.Name] = colType
			}
		}
	}

	return types
}

// sortColumns sorts column names (simple alphabetical for now)
func sortColumns(columns []string) {
	// Simple bubble sort
	n := len(columns)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if columns[j] > columns[j+1] {
				columns[j], columns[j+1] = columns[j+1], columns[j]
			}
		}
	}
}

// calculateColumnWidths calculates optimal width for each column
func calculateColumnWidths(columns []string, rows []map[string]interface{}, totalWidth int) []int {
	if len(columns) == 0 {
		return []int{}
	}

	numCols := len(columns)
	separatorSpace := numCols - 1 // 1 space between each column
	availableForColumns := totalWidth - separatorSpace

	// If we can't even fit all columns with minimum width, give equal minimal space
	if availableForColumns < numCols*3 {
		minWidth := availableForColumns / numCols
		if minWidth < 1 {
			minWidth = 1
		}
		widths := make([]int, numCols)
		for i := range widths {
			widths[i] = minWidth
		}
		// Distribute any remainder
		remainder := availableForColumns - (minWidth * numCols)
		for i := 0; i < remainder; i++ {
			widths[i]++
		}
		return widths
	}

	// Start with header widths
	widths := make([]int, numCols)
	for i, col := range columns {
		widths[i] = len(col)
		if widths[i] < 3 {
			widths[i] = 3 // Minimum width
		}
	}

	// Check data widths (sample first 100 rows for performance)
	sampleSize := len(rows)
	if sampleSize > 100 {
		sampleSize = 100
	}
	for i := 0; i < sampleSize; i++ {
		row := rows[i]
		for j, col := range columns {
			if val, exists := row[col]; exists && val != nil {
				valStr := fmt.Sprintf("%v", val)
				if len(valStr) > widths[j] {
					widths[j] = len(valStr)
				}
			}
		}
	}

	// Calculate total needed
	totalNeeded := 0
	for _, w := range widths {
		totalNeeded += w
	}

	// If total exceeds available width, proportionally reduce
	if totalNeeded > availableForColumns {
		// Proportionally scale down
		scale := float64(availableForColumns) / float64(totalNeeded)
		for i := range widths {
			widths[i] = int(float64(widths[i]) * scale)
			if widths[i] < 3 {
				widths[i] = 3
			}
		}

		// After applying minimum widths, verify total fits
		actualTotal := 0
		for _, w := range widths {
			actualTotal += w
		}

		// If still too wide, reduce from largest columns iteratively
		for actualTotal > availableForColumns {
			// Find largest column
			maxIdx := 0
			maxWidth := widths[0]
			for i, w := range widths {
				if w > maxWidth {
					maxWidth = w
					maxIdx = i
				}
			}
			// Reduce it by 1 if possible
			if widths[maxIdx] > 3 {
				widths[maxIdx]--
				actualTotal--
			} else {
				// All columns at minimum, forcefully reduce
				widths[maxIdx]--
				actualTotal--
				if widths[maxIdx] < 1 {
					widths[maxIdx] = 1
				}
			}
		}
	}

	return widths
}

// renderHeaderRow renders the column headers
func renderHeaderRow(columns []string, widths []int, totalWidth int) string {
	var parts []string
	for i, col := range columns {
		w := widths[i]
		if len(col) > w {
			col = col[:w-1] + "…"
		} else {
			col = col + strings.Repeat(" ", w-len(col))
		}
		parts = append(parts, col)
	}
	line := strings.Join(parts, " ")

	// Ensure exact width
	currentLen := len(line)
	if currentLen < totalWidth {
		line += strings.Repeat(" ", totalWidth-currentLen)
	} else if currentLen > totalWidth {
		line = line[:totalWidth]
	}

	return line
}

// renderSeparator renders the separator line (── ───── ────)
func renderSeparator(widths []int, totalWidth int) string {
	var parts []string
	for _, w := range widths {
		parts = append(parts, strings.Repeat("─", w))
	}
	line := strings.Join(parts, " ")

	// Ensure exact width
	currentLen := len(line)
	if currentLen < totalWidth {
		line += strings.Repeat(" ", totalWidth-currentLen)
	} else if currentLen > totalWidth {
		line = line[:totalWidth]
	}

	return line
}

// renderDataRow renders a single data row
func renderDataRow(columns []string, row map[string]interface{}, widths []int, columnTypes map[string]string, totalWidth int, isSelected bool) string {
	var parts []string
	for i, col := range columns {
		w := widths[i]
		val := ""
		if v, exists := row[col]; exists && v != nil {
			val = fmt.Sprintf("%v", v)
		}

		// Check if this column is numeric type
		colType := columnTypes[col]
		isNumeric := isNumericType(colType)

		if len(val) > w {
			val = val[:w-1] + "…"
		} else {
			// Right-align numeric columns, left-align others
			if isNumeric {
				val = strings.Repeat(" ", w-len(val)) + val
			} else {
				val = val + strings.Repeat(" ", w-len(val))
			}
		}
		parts = append(parts, val)
	}
	line := strings.Join(parts, " ")

	// Ensure exact width BEFORE applying styles
	currentLen := len(line)
	if currentLen < totalWidth {
		line += strings.Repeat(" ", totalWidth-currentLen)
	} else if currentLen > totalWidth {
		line = line[:totalWidth]
	}

	// Apply full-span background color for selected row
	if isSelected {
		line = lipgloss.NewStyle().Background(lipgloss.Color(ColorPrimary)).Foreground(lipgloss.Color("#000000")).Render(line)
	}

	return line
}

// isNumericType checks if a column type is numeric
func isNumericType(colType string) bool {
	upperType := strings.ToUpper(colType)
	return strings.Contains(upperType, "INTEGER") ||
		strings.Contains(upperType, "LONG") ||
		strings.Contains(upperType, "DOUBLE") ||
		strings.Contains(upperType, "FLOAT") ||
		strings.Contains(upperType, "NUMBER")
}
