package app

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/camikura/dito/internal/db"
	"github.com/camikura/dito/internal/ui"
)

func renderConnectionPane(m Model, width int) string {
	borderStyle := ui.StyleBorderInactive
	titleStyle := ui.StyleTitleInactive
	if m.CurrentPane == FocusPaneConnection {
		borderStyle = ui.StyleBorderActive
		titleStyle = ui.StyleTitleActive
	}

	var titleText string
	var titleDisplayWidth int
	if m.Connected {
		checkmark := ui.StyleCheckmark.Render("✓")
		titleText = titleStyle.Render(" Connection ") + checkmark + " "
		titleDisplayWidth = 14
	} else {
		titleText = titleStyle.Render(" Connection ")
		titleDisplayWidth = 12
	}

	// Title line: ╭─ + titleText + ─...─ + ╮
	// Total width = width, so dashes = width - 2(╭╮) - 1(─ after ╭) - titleDisplayWidth
	dashesLen := width - 3 - titleDisplayWidth
	if dashesLen < 0 {
		dashesLen = 0
	}
	title := borderStyle.Render("╭─") + titleText + borderStyle.Render(strings.Repeat("─", dashesLen)+"╮")

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
	borderStyle := ui.StyleBorderInactive
	titleStyle := ui.StyleTitleInactive
	if m.CurrentPane == FocusPaneTables {
		borderStyle = ui.StyleBorderActive
		titleStyle = ui.StyleTitleActive
	}

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
		isSelected bool // * marker (Enter pressed)
		isCursor   bool // cursor position (up/down navigation)
	}

	// Determine if selection marker should be shown
	// Hide * when custom SQL targets a table not in the list
	showSelectionMarker := true
	if m.CustomSQL && m.CurrentSQL != "" {
		extractedName := ui.ExtractTableNameFromSQL(m.CurrentSQL)
		if extractedName != "" && m.FindTableName(extractedName) == "" {
			// Custom SQL targets a table not in the list
			showSelectionMarker = false
		}
	}

	var contentLines []tableLineInfo
	if len(m.Tables) == 0 {
		contentLines = []tableLineInfo{{text: "No tables", isSelected: false, isCursor: false}}
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

			// Add selection marker (* for selected table via Enter)
			var prefix string
			isSelected := showSelectionMarker && i == m.SelectedTable
			if isSelected {
				prefix = "* "
			} else {
				prefix = "  "
			}

			contentLines = append(contentLines, tableLineInfo{
				text:       prefix + indent + displayName,
				isSelected: isSelected,
				isCursor:   i == m.CursorTable,
			})
		}
	}

	leftBorder := borderStyle.Render("│")
	rightBorder := borderStyle.Render("│")
	bottomBorder := borderStyle.Render("╰" + strings.Repeat("─", width-2) + "╯")

	var result strings.Builder
	result.WriteString(title + "\n")

	// Render content lines (fill allocated height with content or empty lines)
	isFocused := m.CurrentPane == FocusPaneTables
	for i := 0; i < height; i++ {
		contentIndex := i + m.TablesScrollOffset
		if contentIndex < len(contentLines) {
			lineInfo := contentLines[contentIndex]
			// Apply color based on state
			var styledText string
			if isFocused && lineInfo.isCursor {
				styledText = ui.StyleTableCursor.Render(lineInfo.text)
			} else if lineInfo.isSelected {
				styledText = ui.StyleTableSelected.Render(lineInfo.text)
			} else {
				styledText = ui.StyleTableNormal.Render(lineInfo.text)
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
	// Determine which table to show schema for
	// Use SelectedTable, or extract from custom SQL if applicable
	var schemaTableName string
	if m.CustomSQL && m.CurrentSQL != "" {
		// Extract table name from custom SQL and find exact match from tables list
		extractedName := ui.ExtractTableNameFromSQL(m.CurrentSQL)
		schemaTableName = m.FindTableName(extractedName)
		// Use extracted name if not found in tables list
		if schemaTableName == "" && extractedName != "" {
			schemaTableName = extractedName
		}
	} else if m.SelectedTable >= 0 && m.SelectedTable < len(m.Tables) {
		schemaTableName = m.Tables[m.SelectedTable]
	}

	// Title includes table name if available
	titleText := " Schema "
	if schemaTableName != "" {
		titleText = " Schema (" + schemaTableName + ") "
	}

	// Schema pane can be focused for scrolling
	borderStyle := ui.StyleBorderInactive
	if m.CurrentPane == FocusPaneSchema {
		borderStyle = ui.StyleBorderActive
	}
	title := borderStyle.Render("╭─" + titleText + strings.Repeat("─", width-len(titleText)-3) + "╮")

	// Prepare content lines
	var contentLines []string
	var schemaError string
	if m.SchemaErrorMsg != "" {
		schemaError = m.SchemaErrorMsg
	}
	if schemaTableName == "" {
		if m.CustomSQL && m.CurrentSQL != "" {
			// Custom SQL with table not found in tables list
			contentLines = []string{"No schema"}
		} else {
			contentLines = []string{"Select a table"}
		}
	} else if schemaError != "" {
		// Show error message
		contentLines = []string{schemaError}
	} else {
		details, exists := m.TableDetails[schemaTableName]
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

	var result strings.Builder
	result.WriteString(title + "\n")

	// Render content lines (fill allocated height with content or empty lines)
	for i := 0; i < height; i++ {
		var line string
		contentIndex := i + m.SchemaScrollOffset
		if contentIndex < len(contentLines) {
			content := contentLines[contentIndex]

			// Apply color to label lines (Columns:, Indexes:)
			if content == "Columns:" || content == "Indexes:" {
				paddingLen := width - len(content) - 2
				if paddingLen < 0 {
					paddingLen = 0
				}
				line = ui.StyleSchemaLabel.Render(content) + strings.Repeat(" ", paddingLen)
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
						fieldsDisplay += ui.StyleSchemaIndex.Render(field)
					}
				} else {
					// Single field
					fieldsDisplay = ui.StyleSchemaIndex.Render(fields)
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
					pkField = ui.StyleSchemaPK.Render(pkMarker) + " "
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
				typeField := ui.StyleSchemaType.Render(colType)

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
			// Other content (like "Select a table", "Loading...", error messages)
			if len(content) > width-2 {
				content = content[:width-5] + "..."
			}
			paddingLen := width - len(content) - 2
			if paddingLen < 0 {
				paddingLen = 0
			}
			// Use red for errors, gray for other messages
			if schemaError != "" && content == schemaError {
				line = ui.StyleErrorLight.Render(content) + strings.Repeat(" ", paddingLen)
			} else {
				line = ui.StyleGrayText.Render(content) + strings.Repeat(" ", paddingLen)
			}
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
	isFocused := m.CurrentPane == FocusPaneSQL
	borderStyle := ui.StyleBorderInactive
	titleStyle := ui.StyleTitleInactive
	if isFocused {
		borderStyle = ui.StyleBorderActive
		titleStyle = ui.StyleTitleActive
	}

	// Add [Custom] label if custom SQL is active
	titleText := " SQL "
	if m.CustomSQL {
		titleText = " SQL [Custom] "
	}
	styledTitle := titleStyle.Render(titleText)
	title := borderStyle.Render("╭─") + styledTitle + borderStyle.Render(strings.Repeat("─", width-len(titleText)-3) + "╮")

	leftBorder := borderStyle.Render("│")
	rightBorder := borderStyle.Render("│")
	bottomBorder := borderStyle.Render("╰" + strings.Repeat("─", width-2) + "╯")

	contentWidth := width - 2 // Width inside borders

	// Wrap SQL text to fit content width and track cursor position
	type wrappedLine struct {
		text      string
		cursorCol int // -1 if cursor is not on this line, >= 0 means cursor position
	}

	var wrappedLines []wrappedLine
	sqlRunes := []rune(m.CurrentSQL)

	if len(sqlRunes) == 0 {
		// Empty SQL
		if isFocused {
			wrappedLines = append(wrappedLines, wrappedLine{text: "", cursorCol: 0})
		} else {
			wrappedLines = append(wrappedLines, wrappedLine{text: "", cursorCol: -1})
		}
	} else {
		// Wrap text
		lineStart := 0
		lineWidth := 0
		for i, r := range sqlRunes {
			charWidth := lipgloss.Width(string(r))

			if r == '\n' {
				// End of logical line
				line := string(sqlRunes[lineStart:i])
				cursorCol := -1
				if isFocused && m.SQLCursorPos >= lineStart && m.SQLCursorPos <= i {
					cursorCol = m.SQLCursorPos - lineStart
				}
				wrappedLines = append(wrappedLines, wrappedLine{text: line, cursorCol: cursorCol})
				lineStart = i + 1
				lineWidth = 0
			} else if lineWidth+charWidth > contentWidth && lineWidth > 0 {
				// Wrap line
				line := string(sqlRunes[lineStart:i])
				cursorCol := -1
				if isFocused && m.SQLCursorPos >= lineStart && m.SQLCursorPos < i {
					cursorCol = m.SQLCursorPos - lineStart
				}
				wrappedLines = append(wrappedLines, wrappedLine{text: line, cursorCol: cursorCol})
				lineStart = i
				lineWidth = charWidth
			} else {
				lineWidth += charWidth
			}
		}

		// Add remaining text
		if lineStart <= len(sqlRunes) {
			line := string(sqlRunes[lineStart:])
			lineDisplayWidth := lipgloss.Width(line)
			cursorCol := -1
			if isFocused && m.SQLCursorPos >= lineStart {
				cursorCol = m.SQLCursorPos - lineStart
				// If cursor is at end and line is full width, move cursor to next line
				if cursorCol == len([]rune(line)) && lineDisplayWidth >= contentWidth {
					wrappedLines = append(wrappedLines, wrappedLine{text: line, cursorCol: -1})
					wrappedLines = append(wrappedLines, wrappedLine{text: "", cursorCol: 0})
				} else {
					wrappedLines = append(wrappedLines, wrappedLine{text: line, cursorCol: cursorCol})
				}
			} else {
				wrappedLines = append(wrappedLines, wrappedLine{text: line, cursorCol: cursorCol})
			}
		}
	}

	var result strings.Builder
	result.WriteString(title + "\n")

	// Render wrapped lines
	for i := 0; i < height; i++ {
		var lineContent string
		var lineDisplayWidth int

		if i < len(wrappedLines) {
			wl := wrappedLines[i]
			lineRunes := []rune(wl.text)

			if wl.cursorCol >= 0 {
				// This line has the cursor
				if wl.cursorCol < len(lineRunes) {
					beforeCursor := string(lineRunes[:wl.cursorCol])
					cursorChar := string(lineRunes[wl.cursorCol])
					afterCursor := string(lineRunes[wl.cursorCol+1:])

					var cursorBlock string
					if lipgloss.Width(cursorChar) > 1 {
						cursorBlock = ui.CursorWide.Render(cursorChar)
					} else {
						cursorBlock = ui.CursorNarrow.Render(cursorChar)
					}
					lineContent = beforeCursor + cursorBlock + afterCursor
					lineDisplayWidth = lipgloss.Width(wl.text)
				} else {
					// Cursor at end of line
					lineContent = wl.text + ui.CursorNarrow.Render(" ")
					lineDisplayWidth = lipgloss.Width(wl.text) + 1
				}
			} else {
				lineContent = wl.text
				lineDisplayWidth = lipgloss.Width(wl.text)
			}
		}

		paddingLen := contentWidth - lineDisplayWidth
		if paddingLen < 0 {
			paddingLen = 0
		}
		result.WriteString(leftBorder + lineContent + strings.Repeat(" ", paddingLen) + rightBorder + "\n")
	}
	result.WriteString(bottomBorder)

	return result.String()
}

func renderDataPane(m Model, width int, totalHeight int) string {
	borderStyle := ui.StyleBorderInactive
	titleStyle := ui.StyleTitleInactive
	if m.CurrentPane == FocusPaneData {
		borderStyle = ui.StyleBorderActive
		titleStyle = ui.StyleTitleActive
	}

	// Build title with table name if available
	// For custom SQL, show extracted table name even if not in tables list
	var dataTableName string
	if m.CustomSQL && m.CurrentSQL != "" {
		// Use extracted name directly (may include child table like orders.addresses)
		dataTableName = ui.ExtractTableNameFromSQL(m.CurrentSQL)
	} else if m.SelectedTable >= 0 && m.SelectedTable < len(m.Tables) {
		dataTableName = m.Tables[m.SelectedTable]
	}

	var titleText string
	if dataTableName != "" {
		titleText = fmt.Sprintf(" Data (%s) ", dataTableName)
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

	// Determine which table's data to display
	// For custom SQL, use the extracted table name; otherwise use SelectedTable
	var dataLookupTableName string
	if m.CustomSQL && m.CurrentSQL != "" {
		dataLookupTableName = ui.ExtractTableNameFromSQL(m.CurrentSQL)
	} else if m.SelectedTable >= 0 && m.SelectedTable < len(m.Tables) {
		dataLookupTableName = m.Tables[m.SelectedTable]
	}

	// Prepare content
	if dataLookupTableName == "" {
		// No table selected
		for i := 0; i < contentLines; i++ {
			line := "Select a table"
			if i > 0 {
				line = ""
			}
			paddingLen := width - len(line) - 2
			if paddingLen < 0 {
				paddingLen = 0
			}
			styledLine := line
			if line != "" {
				styledLine = ui.StyleGrayText.Render(line)
			}
			result.WriteString(leftBorder + styledLine + strings.Repeat(" ", paddingLen) + rightBorder + "\n")
		}
	} else {
		data, exists := m.TableData[dataLookupTableName]
		if m.DataErrorMsg != "" {
			// Show error message
			errMsg := m.DataErrorMsg
			maxLen := width - 4
			if len(errMsg) > maxLen {
				errMsg = errMsg[:maxLen-3] + "..."
			}
			for i := 0; i < contentLines; i++ {
				line := ""
				styledLine := ""
				if i == 0 {
					line = errMsg
					styledLine = ui.StyleErrorLight.Render(line)
				}
				paddingLen := width - len(line) - 2
				if paddingLen < 0 {
					paddingLen = 0
				}
				result.WriteString(leftBorder + styledLine + strings.Repeat(" ", paddingLen) + rightBorder + "\n")
			}
		} else if !exists || data == nil {
			// No data loaded yet
			message := "No data"
			if m.LoadingData {
				message = "Loading..."
			}
			for i := 0; i < contentLines; i++ {
				line := ""
				styledLine := ""
				if i == 0 {
					line = message
					styledLine = ui.StyleGrayText.Render(line)
				}
				paddingLen := width - len(line) - 2
				if paddingLen < 0 {
					paddingLen = 0
				}
				result.WriteString(leftBorder + styledLine + strings.Repeat(" ", paddingLen) + rightBorder + "\n")
			}
		} else if len(data.Rows) == 0 {
			// No rows in result
			for i := 0; i < contentLines; i++ {
				line := ""
				styledLine := ""
				if i == 0 {
					line = "No rows"
					styledLine = ui.StyleGrayText.Render(line)
				}
				paddingLen := width - len(line) - 2
				if paddingLen < 0 {
					paddingLen = 0
				}
				result.WriteString(leftBorder + styledLine + strings.Repeat(" ", paddingLen) + rightBorder + "\n")
			}
		} else {
			// Render grid view and get scroll info
			gridContent, scrollInfo := renderGridViewWithScrollInfo(m, dataLookupTableName, data, width, contentLines, leftBorder, borderStyle)
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
func renderGridViewWithScrollInfo(m Model, tableName string, data *db.TableDataResult, width int, contentLines int, leftBorder string, borderStyle lipgloss.Style) (string, scrollInfo) {
	// Get column names in schema definition order
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
	grid.ShowLoading = m.LoadingData
	grid.HasMore = data.HasMore

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
