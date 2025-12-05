package app

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/camikura/dito/internal/ui"
)

// moveCursorUpInText moves cursor up one line in multi-line text
func moveCursorUpInText(text string, cursorPos int) int {
	lines := strings.Split(text, "\n")
	if len(lines) <= 1 {
		return cursorPos
	}

	// Find current line and column
	currentPos := 0
	currentLine := 0
	currentCol := 0
	for i, line := range lines {
		lineLen := len([]rune(line)) + 1 // +1 for newline
		if currentPos+lineLen > cursorPos {
			currentLine = i
			currentCol = cursorPos - currentPos
			break
		}
		currentPos += lineLen
	}

	if currentLine == 0 {
		return 0 // Already on first line, go to start
	}

	// Move to previous line, same column or end of line
	prevLine := lines[currentLine-1]
	prevLineLen := len([]rune(prevLine))
	newCol := currentCol
	if newCol > prevLineLen {
		newCol = prevLineLen
	}

	// Calculate new position
	newPos := 0
	for i := 0; i < currentLine-1; i++ {
		newPos += len([]rune(lines[i])) + 1
	}
	newPos += newCol

	return newPos
}

// moveCursorToLineStart moves cursor to the beginning of the current line (Emacs Ctrl+A)
func moveCursorToLineStart(text string, cursorPos int) int {
	lines := strings.Split(text, "\n")
	if len(lines) <= 1 {
		return 0 // Single line, go to start
	}

	// Find current line start
	currentPos := 0
	for _, line := range lines {
		lineLen := len([]rune(line)) + 1 // +1 for newline
		if currentPos+lineLen > cursorPos {
			return currentPos // Return start of current line
		}
		currentPos += lineLen
	}

	return currentPos // Fallback
}

// moveCursorToLineEnd moves cursor to the end of the current line (Emacs Ctrl+E)
func moveCursorToLineEnd(text string, cursorPos int) int {
	lines := strings.Split(text, "\n")
	if len(lines) <= 1 {
		return len([]rune(text)) // Single line, go to end
	}

	// Find current line end
	currentPos := 0
	for _, line := range lines {
		lineLen := len([]rune(line)) + 1 // +1 for newline
		if currentPos+lineLen > cursorPos {
			// End of line is before the newline character
			return currentPos + len([]rune(line))
		}
		currentPos += lineLen
	}

	return len([]rune(text)) // Fallback
}

// moveCursorDownInText moves cursor down one line in multi-line text
func moveCursorDownInText(text string, cursorPos int) int {
	lines := strings.Split(text, "\n")
	if len(lines) <= 1 {
		return len([]rune(text)) // Go to end
	}

	// Find current line and column
	currentPos := 0
	currentLine := 0
	currentCol := 0
	for i, line := range lines {
		lineLen := len([]rune(line)) + 1 // +1 for newline
		if currentPos+lineLen > cursorPos {
			currentLine = i
			currentCol = cursorPos - currentPos
			break
		}
		currentPos += lineLen
	}

	if currentLine >= len(lines)-1 {
		return len([]rune(text)) // Already on last line, go to end
	}

	// Move to next line, same column or end of line
	nextLine := lines[currentLine+1]
	nextLineLen := len([]rune(nextLine))
	newCol := currentCol
	if newCol > nextLineLen {
		newCol = nextLineLen
	}

	// Calculate new position
	newPos := 0
	for i := 0; i <= currentLine; i++ {
		newPos += len([]rune(lines[i])) + 1
	}
	newPos += newCol

	return newPos
}

// calculateMaxHorizontalOffset calculates the maximum horizontal scroll offset
// so that the rightmost column is visible at the right edge of the pane
func calculateMaxHorizontalOffset(m Model) int {
	if m.SelectedTable < 0 || m.SelectedTable >= len(m.Tables) {
		return 0
	}

	tableName := m.Tables[m.SelectedTable]
	data, exists := m.TableData[tableName]
	if !exists || data == nil || len(data.Rows) == 0 {
		return 0
	}

	// Calculate total content width
	columns := getColumnsInSchemaOrder(m, tableName, data.Rows)
	columnTypes := getColumnTypes(m, tableName, columns)
	grid := ui.NewGrid(columns, columnTypes, data.Rows)
	totalWidth := grid.TotalContentWidth()

	// Calculate viewport width (right pane width - borders)
	leftPaneContentWidth := ui.LeftPaneContentWidth
	rightPaneWidth := m.Width - leftPaneContentWidth - 2 // -2 for borders

	// Max offset is total width minus viewport width
	maxOffset := totalWidth - rightPaneWidth
	if maxOffset < 0 {
		maxOffset = 0
	}

	return maxOffset
}

// calculateRecordDetailMaxScroll calculates the maximum scroll position for record detail dialog
func calculateRecordDetailMaxScroll(m Model) int {
	if m.SelectedTable < 0 || m.SelectedTable >= len(m.Tables) {
		return 0
	}

	tableName := m.Tables[m.SelectedTable]
	data, exists := m.TableData[tableName]
	if !exists || data == nil || m.SelectedDataRow >= len(data.Rows) {
		return 0
	}

	row := data.Rows[m.SelectedDataRow]

	// Get columns in order
	columns := getColumnsInSchemaOrder(m, tableName, data.Rows)

	// Calculate dialog dimensions (must match dialogs.go)
	dialogWidth := m.Width * 4 / 5
	dialogHeight := m.Height * 4 / 5
	contentWidth := dialogWidth - 2 // borders
	innerWidth := contentWidth - 2  // padding (1 on each side)
	contentHeight := dialogHeight - 2

	// Use VerticalTable to get actual rendered line count (with wrapping)
	vt := ui.VerticalTable{
		Data:     row,
		Keys:     columns,
		MaxWidth: innerWidth,
	}
	content := vt.Render()
	lineCount := len(strings.Split(content, "\n"))

	maxScroll := lineCount - contentHeight
	if maxScroll < 0 {
		maxScroll = 0
	}

	return maxScroll
}

// buildDefaultSQL generates the default SELECT statement for a table.
// If primary keys are available from DDL, adds ORDER BY clause.
func buildDefaultSQL(tableName string, ddl string) string {
	sql := "SELECT * FROM " + tableName
	if ddl != "" {
		primaryKeys := ui.ParsePrimaryKeysFromDDL(ddl)
		if len(primaryKeys) > 0 {
			sql += " ORDER BY " + strings.Join(primaryKeys, ", ")
		}
	}
	return sql
}

// calculatePaneHeights calculates pane heights using the same logic as view.go
func calculatePaneHeights(m Model) (tablesHeight, schemaHeight, sqlHeight int) {
	// Render connection pane and count its actual lines (same as view.go)
	connectionPane := renderConnectionPane(m, ui.LeftPaneContentWidth)
	connectionPaneHeight := strings.Count(connectionPane, "\n") + 1

	availableHeight := m.Height - 1 - connectionPaneHeight - 6
	partHeight := availableHeight / ui.PaneHeightTotalParts

	tablesHeight = partHeight * ui.PaneHeightTablesParts
	schemaHeight = partHeight * ui.PaneHeightSchemaParts
	sqlHeight = partHeight * ui.PaneHeightSQLParts

	remainder := availableHeight % ui.PaneHeightTotalParts
	ui.DistributeSpace(remainder, &tablesHeight, &schemaHeight, &sqlHeight)

	if tablesHeight < 3 {
		tablesHeight = 3
	}
	if schemaHeight < 3 {
		schemaHeight = 3
	}
	if sqlHeight < 2 {
		sqlHeight = 2
	}

	usedHeight := tablesHeight + schemaHeight + sqlHeight
	if usedHeight < availableHeight {
		ui.DistributeSpace(availableHeight-usedHeight, &tablesHeight, &schemaHeight, &sqlHeight)
	}

	return tablesHeight, schemaHeight, sqlHeight
}

// calculateSchemaHeight calculates the schema pane height using the same logic as view.go
func calculateSchemaHeight(m Model) int {
	_, schemaHeight, _ := calculatePaneHeights(m)
	return schemaHeight
}

// calculateTablesHeight calculates the tables pane height using the same logic as view.go
func calculateTablesHeight(m Model) int {
	tablesHeight, _, _ := calculatePaneHeights(m)
	return tablesHeight
}

// calculateSQLHeight calculates the SQL pane height using the same logic as view.go
func calculateSQLHeight(m Model) int {
	_, _, sqlHeight := calculatePaneHeights(m)
	return sqlHeight
}

// calculateSchemaContentLineCount calculates the number of content lines in the schema pane
// This must match the logic in pane_tables.go renderSchemaPaneWithHeight
func calculateSchemaContentLineCount(m Model, schemaTableName string) int {
	details, exists := m.TableDetails[schemaTableName]
	if !exists || details == nil {
		return 1 // "Loading..." or similar
	}

	lineCount := 1 // "Columns:"

	// Count inherited primary key columns from ancestors
	ancestors := ui.GetAncestorTableNames(schemaTableName)
	for _, ancestorName := range ancestors {
		if ancestorDetails, ancestorExists := m.TableDetails[ancestorName]; ancestorExists && ancestorDetails != nil && ancestorDetails.Schema != nil && ancestorDetails.Schema.DDL != "" {
			ancestorPKs := ui.ParsePrimaryKeysFromDDL(ancestorDetails.Schema.DDL)
			ancestorCols := ui.ParseColumnsFromDDL(ancestorDetails.Schema.DDL, ancestorPKs)
			for _, col := range ancestorCols {
				if col.IsPrimaryKey {
					lineCount++
				}
			}
		}
	}

	// Count this table's own columns
	if details.Schema != nil && details.Schema.DDL != "" {
		primaryKeys := ui.ParsePrimaryKeysFromDDL(details.Schema.DDL)
		columns := ui.ParseColumnsFromDDL(details.Schema.DDL, primaryKeys)
		lineCount += len(columns)
	}

	lineCount += 2 // Empty line + "Indexes:"
	lineCount += len(details.Indexes)
	if len(details.Indexes) == 0 {
		lineCount++ // "(none)" line
	}

	return lineCount
}

// updateSQLScrollOffset updates the SQL scroll offset to keep cursor visible
// Returns the updated scroll offset
func updateSQLScrollOffset(m Model) int {
	sqlHeight := calculateSQLHeight(m)
	contentWidth := ui.LeftPaneContentWidth - 2 // Width inside borders

	// Calculate which wrapped line the cursor is on
	sqlRunes := []rune(m.CurrentSQL)
	cursorLineIndex := 0

	if len(sqlRunes) > 0 {
		lineStart := 0
		lineWidth := 0
		lineCount := 0

		for i, r := range sqlRunes {
			charWidth := lipgloss.Width(string(r))

			if r == '\n' {
				if m.SQLCursorPos >= lineStart && m.SQLCursorPos <= i {
					cursorLineIndex = lineCount
					break
				}
				lineCount++
				lineStart = i + 1
				lineWidth = 0
			} else if lineWidth+charWidth > contentWidth && lineWidth > 0 {
				if m.SQLCursorPos >= lineStart && m.SQLCursorPos < i {
					cursorLineIndex = lineCount
					break
				}
				lineCount++
				lineStart = i
				lineWidth = charWidth
			} else {
				lineWidth += charWidth
			}
		}

		// Check if cursor is in the last segment
		if m.SQLCursorPos >= lineStart {
			cursorLineIndex = lineCount
		}
	}

	// Adjust scroll offset to keep cursor visible
	scrollOffset := m.SQLScrollOffset
	if cursorLineIndex < scrollOffset {
		scrollOffset = cursorLineIndex
	} else if cursorLineIndex >= scrollOffset+sqlHeight {
		scrollOffset = cursorLineIndex - sqlHeight + 1
	}

	return scrollOffset
}
