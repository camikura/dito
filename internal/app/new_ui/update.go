package new_ui

import (
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/camikura/dito/internal/db"
	"github.com/camikura/dito/internal/ui"
)

// Update handles messages and updates the model
func Update(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		// Calculate pane heights dynamically
		// Use actual connection pane height from model, or default to 5 if not yet set
		connectionPaneHeight := m.ConnectionPaneHeight
		if connectionPaneHeight == 0 {
			connectionPaneHeight = 5 // Default for cloud connection
		}
		// Available height for Tables, Schema, and SQL content (2:2:1 ratio)
		// Total: m.Height = leftPanes + footer
		// leftPanes = Connection + Tables(+2) + Schema(+2) + SQL(+2)
		// So: availableHeight = m.Height - 1(footer) - connectionPaneHeight - 6(borders)
		availableHeight := m.Height - 1 - connectionPaneHeight - 6

		// Split available height in 2:2:1 ratio (Tables:Schema:SQL)
		partHeight := availableHeight / ui.PaneHeightTotalParts
		remainder := availableHeight % ui.PaneHeightTotalParts

		m.TablesHeight = partHeight * ui.PaneHeightTablesParts
		m.SchemaHeight = partHeight * ui.PaneHeightSchemaParts
		m.SQLHeight = partHeight * ui.PaneHeightSQLParts

		// Distribute remainder
		ui.DistributeSpace(remainder, &m.TablesHeight, &m.SchemaHeight, &m.SQLHeight)

		// Ensure minimum heights
		if m.TablesHeight < 3 {
			m.TablesHeight = 3
		}
		if m.SchemaHeight < 3 {
			m.SchemaHeight = 3
		}
		if m.SQLHeight < 2 {
			m.SQLHeight = 2
		}

		// After applying minimum heights, redistribute unused space
		usedHeight := m.TablesHeight + m.SchemaHeight + m.SQLHeight
		if usedHeight < availableHeight {
			ui.DistributeSpace(availableHeight-usedHeight, &m.TablesHeight, &m.SchemaHeight, &m.SQLHeight)
		}

		return m, nil

	case tea.KeyMsg:
		return handleKeyPress(m, msg)

	case db.ConnectionResult:
		return handleConnectionResult(m, msg)

	case db.TableListResult:
		return handleTableListResult(m, msg)

	case db.TableDetailsResult:
		return handleTableDetailsResult(m, msg)

	case db.TableDataResult:
		return handleTableDataResult(m, msg)
	}

	return m, nil
}

func handleKeyPress(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Connection dialog takes precedence
	if m.ConnectionDialogVisible {
		return handleConnectionDialogKeys(m, msg)
	}

	// Record detail dialog takes precedence
	if m.RecordDetailVisible {
		return handleRecordDetailKeys(m, msg)
	}

	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "tab":
		// Only allow pane switching when connected
		if m.Connected {
			m = m.NextPane()
		}
		return m, nil

	case "shift+tab":
		// Only allow pane switching when connected
		if m.Connected {
			m = m.PrevPane()
		}
		return m, nil
	}

	// Pane-specific keys
	switch m.CurrentPane {
	case FocusPaneConnection:
		return handleConnectionKeys(m, msg)
	case FocusPaneTables:
		return handleTablesKeys(m, msg)
	case FocusPaneSchema:
		return handleSchemaKeys(m, msg)
	case FocusPaneSQL:
		return handleSQLKeys(m, msg)
	case FocusPaneData:
		return handleDataKeys(m, msg)
	}

	return m, nil
}

func handleConnectionKeys(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Open connection setup dialog
		m.ConnectionDialogVisible = true
		m.ConnectionDialogField = 0
		// Initialize with current values or defaults
		if m.EditEndpoint == "" {
			m.EditEndpoint = "localhost"
		}
		if m.EditPort == "" {
			m.EditPort = "8080"
		}
		m.EditCursorPos = ui.RuneLen(m.EditEndpoint)
		return m, nil

	case "ctrl+d":
		// Disconnect
		if m.Connected {
			m.Connected = false
			m.NosqlClient = nil
			m.Tables = []string{}
			m.SelectedTable = -1
			m.CursorTable = 0
			m.TableDetails = make(map[string]*db.TableDetailsResult)
			m.TableData = make(map[string]*db.TableDataResult)
			m.ConnectionMsg = ""
		}
		return m, nil
	}

	return m, nil
}

func handleConnectionDialogKeys(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Helper to check if current field is a text field
	isTextField := m.ConnectionDialogField == 0 || m.ConnectionDialogField == 1

	// Helper to move to field and set cursor position
	moveToField := func(field int) {
		m.ConnectionDialogField = field
		if field == 0 {
			m.EditCursorPos = ui.RuneLen(m.EditEndpoint)
		} else if field == 1 {
			m.EditCursorPos = ui.RuneLen(m.EditPort)
		}
	}

	switch msg.Type {
	case tea.KeyEsc:
		// Close dialog
		m.ConnectionDialogVisible = false
		return m, nil

	case tea.KeyEnter:
		// Connect from any field
		m.ConnectionDialogVisible = false
		m.Endpoint = m.EditEndpoint + ":" + m.EditPort
		return m, db.Connect(m.EditEndpoint, m.EditPort, false)

	case tea.KeyTab, tea.KeyDown:
		moveToField((m.ConnectionDialogField + 1) % 2)
		return m, nil

	case tea.KeyShiftTab, tea.KeyUp:
		moveToField((m.ConnectionDialogField + 1) % 2)
		return m, nil

	case tea.KeyBackspace:
		if isTextField {
			if m.ConnectionDialogField == 0 {
				m.EditEndpoint, m.EditCursorPos = ui.Backspace(m.EditEndpoint, m.EditCursorPos)
			} else {
				m.EditPort, m.EditCursorPos = ui.Backspace(m.EditPort, m.EditCursorPos)
			}
		}
		return m, nil

	case tea.KeyDelete:
		if isTextField {
			if m.ConnectionDialogField == 0 {
				m.EditEndpoint = ui.DeleteAt(m.EditEndpoint, m.EditCursorPos)
			} else {
				m.EditPort = ui.DeleteAt(m.EditPort, m.EditCursorPos)
			}
		}
		return m, nil

	case tea.KeyLeft:
		if isTextField && m.EditCursorPos > 0 {
			m.EditCursorPos--
		}
		return m, nil

	case tea.KeyRight:
		if isTextField {
			maxPos := ui.RuneLen(m.EditEndpoint)
			if m.ConnectionDialogField == 1 {
				maxPos = ui.RuneLen(m.EditPort)
			}
			if m.EditCursorPos < maxPos {
				m.EditCursorPos++
			}
		}
		return m, nil

	case tea.KeyHome:
		if isTextField {
			m.EditCursorPos = 0
		}
		return m, nil

	case tea.KeyEnd:
		if isTextField {
			if m.ConnectionDialogField == 0 {
				m.EditCursorPos = ui.RuneLen(m.EditEndpoint)
			} else {
				m.EditCursorPos = ui.RuneLen(m.EditPort)
			}
		}
		return m, nil

	case tea.KeyRunes:
		if isTextField {
			char := string(msg.Runes)
			if m.ConnectionDialogField == 0 {
				m.EditEndpoint, m.EditCursorPos = ui.InsertWithCursor(m.EditEndpoint, m.EditCursorPos, char)
			} else {
				m.EditPort, m.EditCursorPos = ui.InsertWithCursor(m.EditPort, m.EditCursorPos, char)
			}
		}
		return m, nil
	}

	return m, nil
}

func handleTablesKeys(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	visibleLines := m.TablesHeight // Tables pane visible height (dynamic)

	switch msg.String() {
	case "up", "k":
		if m.CursorTable > 0 {
			m.CursorTable--

			// Adjust scroll offset to keep cursor visible
			if m.CursorTable < m.TablesScrollOffset {
				m.TablesScrollOffset = m.CursorTable
			}
		}
		return m, nil

	case "down", "j":
		if m.CursorTable < len(m.Tables)-1 {
			m.CursorTable++

			// Adjust scroll offset to keep cursor visible
			if m.CursorTable >= m.TablesScrollOffset+visibleLines {
				m.TablesScrollOffset = m.CursorTable - visibleLines + 1
			}
		}
		return m, nil

	case "enter":
		// Select table and load data (only on Enter)
		if m.CursorTable < len(m.Tables) {
			m.SelectedTable = m.CursorTable
			tableName := m.Tables[m.SelectedTable]

			// Reset state
			m.CustomSQL = false
			m.SelectedDataRow = 0
			m.ViewportOffset = 0
			m.HorizontalOffset = 0
			m.SchemaScrollOffset = 0

			// Move focus to Data pane for immediate interaction
			m.CurrentPane = FocusPaneData

			// Check if schema is already loaded
			if details, exists := m.TableDetails[tableName]; exists && details != nil && details.Schema != nil {
				// Schema available - fetch data with ORDER BY
				ddl := details.Schema.DDL
				primaryKeys := ui.ParsePrimaryKeysFromDDL(ddl)
				m.CurrentSQL = buildDefaultSQL(tableName, ddl)
				m.SQLCursorPos = ui.RuneLen(m.CurrentSQL)
				return m, db.FetchTableData(m.NosqlClient, tableName, 100, primaryKeys)
			}

			// Schema not loaded - fetch schema first, data will be fetched when schema arrives
			m.CurrentSQL = "SELECT * FROM " + tableName
			m.SQLCursorPos = ui.RuneLen(m.CurrentSQL)
			m.LoadingData = true
			return m, db.FetchTableDetails(m.NosqlClient, tableName)
		}
		return m, nil
	}

	return m, nil
}

func handleSchemaKeys(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Determine which table schema is displayed (same logic as view)
	var schemaTableName string
	if m.CustomSQL && m.CurrentSQL != "" {
		extractedName := ui.ExtractTableNameFromSQL(m.CurrentSQL)
		schemaTableName = m.FindTableName(extractedName)
		if schemaTableName == "" && extractedName != "" {
			schemaTableName = extractedName
		}
	} else if m.SelectedTable >= 0 && m.SelectedTable < len(m.Tables) {
		schemaTableName = m.Tables[m.SelectedTable]
	}

	// Calculate max scroll offset based on content
	var maxScroll int
	if schemaTableName != "" {
		if details, exists := m.TableDetails[schemaTableName]; exists && details != nil {
			// Count content lines
			lineCount := 1 // "Columns:"
			if details.Schema.DDL != "" {
				primaryKeys := ui.ParsePrimaryKeysFromDDL(details.Schema.DDL)
				columns := ui.ParseColumnsFromDDL(details.Schema.DDL, primaryKeys)
				lineCount += len(columns)
			}
			lineCount += 2 // Empty line + "Indexes:"
			lineCount += len(details.Indexes)
			if len(details.Indexes) == 0 {
				lineCount++ // "(none)" line
			}

			// Max scroll = total lines - visible lines (dynamic)
			maxScroll = lineCount - m.SchemaHeight
			if maxScroll < 0 {
				maxScroll = 0
			}
		}
	}

	switch msg.String() {
	case "down", "j": // Scroll down
		if m.SchemaScrollOffset < maxScroll {
			m.SchemaScrollOffset++
		}
		return m, nil

	case "up", "k": // Scroll up
		if m.SchemaScrollOffset > 0 {
			m.SchemaScrollOffset--
		}
		return m, nil
	}

	return m, nil
}

func handleSQLKeys(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlR:
		// Execute SQL - require connection and SQL
		if !m.Connected || m.CurrentSQL == "" {
			return m, nil
		}

		m.CustomSQL = true

		// Save current selection before custom SQL
		m.PreviousSelectedTable = m.SelectedTable

		// Reset data state
		m.SelectedDataRow = 0
		m.ViewportOffset = 0
		m.HorizontalOffset = 0
		m.SchemaScrollOffset = 0

		// Extract table name from SQL for schema display
		extractedName := ui.ExtractTableNameFromSQL(m.CurrentSQL)
		// Find exact table name from tables list (case-insensitive match)
		tableName := m.FindTableName(extractedName)
		// Update SelectedTable if table is in the list
		tableIndex := m.FindTableIndex(extractedName)
		if tableIndex >= 0 {
			m.SelectedTable = tableIndex
		}
		// Use extracted name if not found in tables list (for child tables etc.)
		if tableName == "" && extractedName != "" {
			tableName = extractedName
		}
		// Fall back to selected table if no table name in SQL
		if tableName == "" && m.SelectedTable >= 0 && m.SelectedTable < len(m.Tables) {
			tableName = m.Tables[m.SelectedTable]
		}

		if tableName != "" {
			var cmds []tea.Cmd

			// Fetch schema (always try, even for unknown tables to get error)
			if _, exists := m.TableDetails[tableName]; !exists {
				cmds = append(cmds, db.FetchTableDetails(m.NosqlClient, tableName))
			}

			// Execute custom SQL
			cmds = append(cmds, db.ExecuteCustomSQL(m.NosqlClient, tableName, m.CurrentSQL, ui.DefaultFetchSize))

			// Move focus to Data pane
			m.CurrentPane = FocusPaneData

			return m, tea.Batch(cmds...)
		}
		return m, nil

	case tea.KeyEnter:
		// Insert newline
		m.CurrentSQL, m.SQLCursorPos = ui.InsertWithCursor(m.CurrentSQL, m.SQLCursorPos, "\n")
		return m, nil

	case tea.KeyBackspace:
		m.CurrentSQL, m.SQLCursorPos = ui.Backspace(m.CurrentSQL, m.SQLCursorPos)
		return m, nil

	case tea.KeyDelete:
		m.CurrentSQL = ui.DeleteAt(m.CurrentSQL, m.SQLCursorPos)
		return m, nil

	case tea.KeyLeft:
		if m.SQLCursorPos > 0 {
			m.SQLCursorPos--
		}
		return m, nil

	case tea.KeyRight:
		if m.SQLCursorPos < ui.RuneLen(m.CurrentSQL) {
			m.SQLCursorPos++
		}
		return m, nil

	case tea.KeyUp:
		// Move cursor up one line
		m.SQLCursorPos = moveCursorUpInText(m.CurrentSQL, m.SQLCursorPos)
		return m, nil

	case tea.KeyDown:
		// Move cursor down one line
		m.SQLCursorPos = moveCursorDownInText(m.CurrentSQL, m.SQLCursorPos)
		return m, nil

	case tea.KeyHome, tea.KeyCtrlA:
		m.SQLCursorPos = 0
		return m, nil

	case tea.KeyEnd, tea.KeyCtrlE:
		m.SQLCursorPos = ui.RuneLen(m.CurrentSQL)
		return m, nil

	case tea.KeySpace:
		m.CurrentSQL, m.SQLCursorPos = ui.InsertWithCursor(m.CurrentSQL, m.SQLCursorPos, " ")
		return m, nil

	case tea.KeyRunes:
		for _, r := range msg.Runes {
			m.CurrentSQL, m.SQLCursorPos = ui.InsertWithCursor(m.CurrentSQL, m.SQLCursorPos, string(r))
		}
		return m, nil
	}

	return m, nil
}

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

func handleDataKeys(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Get total row count for current table
	var totalRows int
	if m.SelectedTable >= 0 && m.SelectedTable < len(m.Tables) {
		tableName := m.Tables[m.SelectedTable]
		if data, exists := m.TableData[tableName]; exists && data != nil {
			totalRows = len(data.Rows)
		}
	}

	// Calculate visible lines for data rows
	// Data pane structure: title(1) + content lines + bottom(1)
	// Content lines = header(1) + separator(1) + data rows
	// contentLines = m.Height - 1(footer) - 2(title+bottom) = m.Height - 3
	contentLines := m.Height - 3
	if contentLines < 5 {
		contentLines = 5
	}
	// Data visible lines = content lines - 2 (header + separator)
	dataVisibleLines := contentLines - 2
	if dataVisibleLines < 1 {
		dataVisibleLines = 1
	}

	// Calculate max horizontal offset
	maxHorizontalOffset := calculateMaxHorizontalOffset(m)

	switch msg.String() {
	case "up", "k":
		if m.SelectedDataRow > 0 {
			m.SelectedDataRow--

			// Calculate middle position of visible area
			middlePosition := dataVisibleLines / 2

			// Calculate maximum viewport offset
			maxViewportOffset := totalRows - dataVisibleLines
			if maxViewportOffset < 0 {
				maxViewportOffset = 0
			}

			// Scrolling logic (symmetric to down):
			// When above middle, keep cursor at middle by adjusting viewport
			// But never exceed maxViewportOffset (when at bottom)
			// When at or below middle, viewport stays at 0
			if m.SelectedDataRow > middlePosition {
				// Still above middle - keep cursor at middle
				m.ViewportOffset = m.SelectedDataRow - middlePosition
				// But don't exceed max offset
				if m.ViewportOffset > maxViewportOffset {
					m.ViewportOffset = maxViewportOffset
				}
			} else {
				// At or below middle - viewport is 0
				m.ViewportOffset = 0
			}
		}
		return m, nil

	case "down", "j":
		if totalRows > 0 && m.SelectedDataRow < totalRows-1 {
			m.SelectedDataRow++

			// Calculate middle position of visible area
			middlePosition := dataVisibleLines / 2

			// Calculate maximum viewport offset (when last row is at bottom of screen)
			maxViewportOffset := totalRows - dataVisibleLines
			if maxViewportOffset < 0 {
				maxViewportOffset = 0
			}

			// Scrolling logic:
			// 1. First: cursor moves to middle (no scroll, VP stays 0)
			// 2. Middle: cursor stays at middle, viewport scrolls
			// 3. End: viewport stops at max, cursor moves to bottom
			if m.SelectedDataRow > middlePosition && m.ViewportOffset < maxViewportOffset {
				// Cursor has passed middle position and we can still scroll
				// Keep cursor at middle by adjusting viewport
				m.ViewportOffset = m.SelectedDataRow - middlePosition
				if m.ViewportOffset > maxViewportOffset {
					m.ViewportOffset = maxViewportOffset
				}
			}

			// Check if we need to fetch more data
			tableName := m.Tables[m.SelectedTable]
			if data, exists := m.TableData[tableName]; exists && data != nil {
				remainingRows := totalRows - m.SelectedDataRow - 1
				if remainingRows <= ui.FetchMoreThreshold && data.HasMore && !m.LoadingData && data.LastPKValues != nil {
					m.LoadingData = true
					// Get PRIMARY KEYs from schema
					var primaryKeys []string
					if details, exists := m.TableDetails[tableName]; exists && details != nil && details.Schema != nil && details.Schema.DDL != "" {
						primaryKeys = ui.ParsePrimaryKeysFromDDL(details.Schema.DDL)
					}
					return m, db.FetchMoreTableData(m.NosqlClient, tableName, 100, primaryKeys, data.LastPKValues)
				}
			}
		}
		return m, nil

	case "left", "h":
		// Scroll left
		if m.HorizontalOffset > 0 {
			m.HorizontalOffset -= 5 // Scroll by 5 characters
			if m.HorizontalOffset < 0 {
				m.HorizontalOffset = 0
			}
		}
		return m, nil

	case "right", "l":
		// Scroll right (limited to max offset)
		if m.HorizontalOffset < maxHorizontalOffset {
			m.HorizontalOffset += 5 // Scroll by 5 characters
			if m.HorizontalOffset > maxHorizontalOffset {
				m.HorizontalOffset = maxHorizontalOffset
			}
		}
		return m, nil

	case "enter":
		// Show record detail dialog
		if totalRows > 0 {
			m.RecordDetailVisible = true
			m.RecordDetailScroll = 0
		}
		return m, nil

	case "esc":
		// Reset to default SQL mode (only when custom SQL is active)
		if m.CustomSQL {
			// Restore previous selection
			if m.PreviousSelectedTable >= 0 && m.PreviousSelectedTable < len(m.Tables) {
				m.SelectedTable = m.PreviousSelectedTable
			}
			m.PreviousSelectedTable = -1

			m.CustomSQL = false
			m.ColumnOrder = nil
			m.SelectedDataRow = 0
			m.ViewportOffset = 0
			m.HorizontalOffset = 0
			m.SchemaErrorMsg = ""
			m.DataErrorMsg = ""

			// Reload data with default SQL if a table is selected
			if m.SelectedTable >= 0 && m.SelectedTable < len(m.Tables) {
				tableName := m.Tables[m.SelectedTable]

				var ddl string
				var primaryKeys []string
				if details, exists := m.TableDetails[tableName]; exists && details != nil && details.Schema != nil {
					ddl = details.Schema.DDL
					primaryKeys = ui.ParsePrimaryKeysFromDDL(ddl)
				}

				m.CurrentSQL = buildDefaultSQL(tableName, ddl)
				m.SQLCursorPos = ui.RuneLen(m.CurrentSQL)
				return m, db.FetchTableData(m.NosqlClient, tableName, ui.DefaultFetchSize, primaryKeys)
			}
			m.CurrentSQL = ""
			m.SQLCursorPos = 0
		}
		return m, nil
	}

	return m, nil
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

	// Get columns in same order as renderGridView
	columns := getColumnsInSchemaOrder(m, tableName, data.Rows)
	columnTypes := getColumnTypes(m, tableName, columns)

	// Calculate data pane width (must match renderDataPane calculation)
	leftPaneActualWidth := ui.LeftPaneContentWidth + ui.LeftPaneBorderWidth
	rightPaneActualWidth := m.Width - leftPaneActualWidth + 1
	contentWidth := rightPaneActualWidth - 2

	// Use the Grid component to calculate max offset (ensures consistency)
	grid := ui.NewGrid(columns, columnTypes, data.Rows)
	grid.Width = contentWidth

	return grid.MaxHorizontalOffset()
}

func handleConnectionResult(m Model, msg db.ConnectionResult) (Model, tea.Cmd) {
	if msg.Err != nil {
		// Connection failed
		m.Connected = false
		m.ConnectionMsg = msg.Err.Error()
		return m, nil
	}

	// Connection successful
	m.Connected = true
	m.NosqlClient = msg.Client
	m.Endpoint = msg.Endpoint
	m.ConnectionMsg = ""

	// Fetch table list
	return m, db.FetchTables(msg.Client)
}

func handleTableListResult(m Model, msg db.TableListResult) (Model, tea.Cmd) {
	if msg.Err != nil {
		// TODO: Show error
		return m, nil
	}

	// Sort tables for tree display (parents before children)
	m.Tables = sortTablesForTree(msg.Tables)
	if len(m.Tables) > 0 {
		m.CursorTable = 0
		// SelectedTable stays at -1 until user presses Enter
	}

	return m, nil
}

// sortTablesForTree sorts table names so parent tables appear before their children
// e.g., ["users.phones", "users", "products", "users.addresses"] ->
//       ["products", "users", "users.addresses", "users.phones"]
func sortTablesForTree(tables []string) []string {
	sorted := make([]string, len(tables))
	copy(sorted, tables)

	sort.Slice(sorted, func(i, j int) bool {
		a, b := sorted[i], sorted[j]

		// Get parent names
		parentA := a
		if dotIndex := strings.LastIndex(a, "."); dotIndex != -1 {
			parentA = a[:dotIndex]
		}
		parentB := b
		if dotIndex := strings.LastIndex(b, "."); dotIndex != -1 {
			parentB = b[:dotIndex]
		}

		// If one is parent of the other, parent comes first
		if a == parentB {
			return true // a is parent of b
		}
		if b == parentA {
			return false // b is parent of a
		}

		// If they have the same parent, sort alphabetically
		if parentA == parentB {
			return a < b
		}

		// Different parents - sort by parent name, then by full name
		if parentA != a && parentB != b {
			// Both are children - compare parents first
			if parentA != parentB {
				return parentA < parentB
			}
		}

		// One is parent, one is not - sort by full name
		return a < b
	})

	return sorted
}

func handleTableDetailsResult(m Model, msg db.TableDetailsResult) (Model, tea.Cmd) {
	if msg.Err != nil {
		m.SchemaErrorMsg = msg.Err.Error()
		m.LoadingData = false
		return m, nil
	}

	// Clear any previous error
	m.SchemaErrorMsg = ""
	m.TableDetails[msg.TableName] = &msg

	// If this is the selected table and we're waiting for data, fetch it now
	if m.LoadingData && !m.CustomSQL && m.SelectedTable >= 0 && m.SelectedTable < len(m.Tables) {
		tableName := m.Tables[m.SelectedTable]
		if tableName == msg.TableName && msg.Schema != nil {
			// Update SQL with ORDER BY
			primaryKeys := ui.ParsePrimaryKeysFromDDL(msg.Schema.DDL)
			m.CurrentSQL = buildDefaultSQL(tableName, msg.Schema.DDL)
			m.SQLCursorPos = ui.RuneLen(m.CurrentSQL)
			// Now fetch data with proper ORDER BY
			return m, db.FetchTableData(m.NosqlClient, tableName, 100, primaryKeys)
		}
	}

	return m, nil
}

func handleTableDataResult(m Model, msg db.TableDataResult) (Model, tea.Cmd) {
	if msg.Err != nil {
		m.LoadingData = false
		m.DataErrorMsg = msg.Err.Error()
		return m, nil
	}

	// Clear any previous error
	m.DataErrorMsg = ""

	// If this is an append operation (additional data fetch), merge with existing data
	if msg.IsAppend {
		if existingData, exists := m.TableData[msg.TableName]; exists && existingData != nil {
			// Append new rows to existing rows
			existingData.Rows = append(existingData.Rows, msg.Rows...)
			// Update pagination info
			existingData.LastPKValues = msg.LastPKValues
			existingData.HasMore = msg.HasMore
		}
	} else {
		// Initial fetch: replace entire data
		m.TableData[msg.TableName] = &msg
	}

	// Store column order from custom SQL
	if msg.IsCustomSQL && len(msg.ColumnOrder) > 0 {
		m.ColumnOrder = msg.ColumnOrder
	} else if !msg.IsCustomSQL {
		m.ColumnOrder = nil
	}

	m.LoadingData = false
	return m, nil
}

// handleRecordDetailKeys handles keys when record detail dialog is visible
func handleRecordDetailKeys(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Calculate max scroll for validation
	maxScroll := calculateRecordDetailMaxScroll(m)

	switch msg.String() {
	case "ctrl+c":
		// Allow quitting even when dialog is open
		return m, tea.Quit

	case "esc", "q":
		// Close dialog
		m.RecordDetailVisible = false
		m.RecordDetailScroll = 0
		return m, nil

	case "down", "j":
		// Scroll down within dialog (only if content overflows)
		if maxScroll > 0 && m.RecordDetailScroll < maxScroll {
			m.RecordDetailScroll++
		}
		return m, nil

	case "up", "k":
		// Scroll up within dialog
		if m.RecordDetailScroll > 0 {
			m.RecordDetailScroll--
		}
		return m, nil

	case "ctrl+d", "pgdown":
		// Page down within the current record
		pageSize := m.Height * 4 / 5 / 2 // Half of dialog height
		if pageSize < 1 {
			pageSize = 1
		}
		m.RecordDetailScroll += pageSize
		if m.RecordDetailScroll > maxScroll {
			m.RecordDetailScroll = maxScroll
		}
		return m, nil

	case "ctrl+u", "pgup":
		// Page up within the current record
		pageSize := m.Height * 4 / 5 / 2 // Half of dialog height
		if pageSize < 1 {
			pageSize = 1
		}
		m.RecordDetailScroll -= pageSize
		if m.RecordDetailScroll < 0 {
			m.RecordDetailScroll = 0
		}
		return m, nil

	case "g", "home":
		// Go to top
		m.RecordDetailScroll = 0
		return m, nil

	case "G", "end":
		// Go to bottom
		m.RecordDetailScroll = maxScroll
		return m, nil
	}

	return m, nil
}

// calculateRecordDetailMaxScroll calculates the maximum scroll offset for record detail dialog
func calculateRecordDetailMaxScroll(m Model) int {
	if m.SelectedTable < 0 || m.SelectedTable >= len(m.Tables) {
		return 0
	}

	tableName := m.Tables[m.SelectedTable]
	data, exists := m.TableData[tableName]
	if !exists || data == nil || len(data.Rows) == 0 {
		return 0
	}

	if m.SelectedDataRow < 0 || m.SelectedDataRow >= len(data.Rows) {
		return 0
	}

	row := data.Rows[m.SelectedDataRow]
	columns := getColumnsInSchemaOrder(m, tableName, data.Rows)

	// Create vertical table to get content
	vt := ui.VerticalTable{
		Data: row,
		Keys: columns,
	}
	content := vt.Render()
	lines := strings.Split(content, "\n")

	// Calculate visible height (dialog is 80% of screen, minus borders)
	// Must match view.go: contentHeight = dialogHeight - 2
	dialogHeight := m.Height * 4 / 5
	contentHeight := dialogHeight - 2 // Subtract top border (1) + bottom border (1)

	maxScroll := len(lines) - contentHeight
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
