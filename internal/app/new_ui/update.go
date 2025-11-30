package new_ui

import (
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/camikura/dito/internal/db"
	"github.com/camikura/dito/internal/ui"
	"github.com/camikura/dito/internal/views"
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
	// SQL Editor dialog takes precedence
	if m.SQLEditorVisible {
		return handleSQLEditorKeys(m, msg)
	}

	// Record detail dialog takes precedence
	if m.RecordDetailVisible {
		return handleRecordDetailKeys(m, msg)
	}

	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "tab":
		m = m.NextPane()
		return m, nil

	case "shift+tab":
		m = m.PrevPane()
		return m, nil
	}

	// Pane-specific keys
	switch m.CurrentPane {
	case FocusPaneTables:
		return handleTablesKeys(m, msg)
	case FocusPaneSchema:
		return handleSchemaKeys(m, msg)
	case FocusPaneData:
		return handleDataKeys(m, msg)
	}

	return m, nil
}

func handleTablesKeys(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	visibleLines := m.TablesHeight // Tables pane visible height (dynamic)

	switch msg.String() {
	case "up", "k":
		if m.CursorTable > 0 {
			m.CursorTable--
			m.SelectedTable = m.CursorTable // Move selection with cursor
			m.SchemaScrollOffset = 0 // Reset scroll when changing tables

			// Adjust scroll offset to keep cursor visible
			if m.CursorTable < m.TablesScrollOffset {
				m.TablesScrollOffset = m.CursorTable
			}

			// Auto-update schema for table under cursor
			if m.CursorTable < len(m.Tables) {
				tableName := m.Tables[m.CursorTable]
				return m, db.FetchTableDetails(m.NosqlClient, tableName)
			}
		}
		return m, nil

	case "down", "j":
		if m.CursorTable < len(m.Tables)-1 {
			m.CursorTable++
			m.SelectedTable = m.CursorTable // Move selection with cursor
			m.SchemaScrollOffset = 0 // Reset scroll when changing tables

			// Adjust scroll offset to keep cursor visible
			if m.CursorTable >= m.TablesScrollOffset+visibleLines {
				m.TablesScrollOffset = m.CursorTable - visibleLines + 1
			}

			// Auto-update schema for table under cursor
			if m.CursorTable < len(m.Tables) {
				tableName := m.Tables[m.CursorTable]
				return m, db.FetchTableDetails(m.NosqlClient, tableName)
			}
		}
		return m, nil

	case "enter":
		// Select table and load data
		if m.CursorTable < len(m.Tables) {
			m.SelectedTable = m.CursorTable
			tableName := m.Tables[m.SelectedTable]

			// Generate SQL query
			m.CurrentSQL = "SELECT * FROM " + tableName
			m.CustomSQL = false

			// Reset data scrolling state
			m.SelectedDataRow = 0
			m.ViewportOffset = 0
			m.HorizontalOffset = 0

			// Move focus to Data pane for immediate interaction
			m.CurrentPane = FocusPaneData

			// Get primary keys from schema if available
			var primaryKeys []string
			if details, exists := m.TableDetails[tableName]; exists && details != nil && details.Schema != nil && details.Schema.DDL != "" {
				primaryKeys = views.ParsePrimaryKeysFromDDL(details.Schema.DDL)
			}

			// Load table data (100 rows with ORDER BY PK)
			return m, db.FetchTableData(m.NosqlClient, tableName, 100, primaryKeys)
		}
		return m, nil
	}

	return m, nil
}

func handleSchemaKeys(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Calculate max scroll offset based on content
	var maxScroll int
	if len(m.Tables) > 0 && m.CursorTable < len(m.Tables) {
		tableName := m.Tables[m.CursorTable]
		if details, exists := m.TableDetails[tableName]; exists && details != nil {
			// Count content lines
			lineCount := 1 // "Columns:"
			if details.Schema.DDL != "" {
				primaryKeys := views.ParsePrimaryKeysFromDDL(details.Schema.DDL)
				columns := views.ParseColumnsFromDDL(details.Schema.DDL, primaryKeys)
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
						primaryKeys = views.ParsePrimaryKeysFromDDL(details.Schema.DDL)
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

	case "e":
		// Open SQL editor dialog
		m.SQLEditorVisible = true
		m.EditSQL = m.CurrentSQL
		m.SQLCursorPos = len(m.EditSQL)
		return m, nil

	case "esc":
		// Reset to default SQL mode (only when custom SQL is active)
		if m.CustomSQL && m.SelectedTable >= 0 && m.SelectedTable < len(m.Tables) {
			tableName := m.Tables[m.SelectedTable]
			m.CustomSQL = false
			m.ColumnOrder = nil
			m.CurrentSQL = "SELECT * FROM " + tableName
			m.SelectedDataRow = 0
			m.ViewportOffset = 0
			m.HorizontalOffset = 0

			// Reload data with default SQL
			var primaryKeys []string
			if details, exists := m.TableDetails[tableName]; exists && details != nil && details.Schema != nil && details.Schema.DDL != "" {
				primaryKeys = views.ParsePrimaryKeysFromDDL(details.Schema.DDL)
			}
			return m, db.FetchTableData(m.NosqlClient, tableName, ui.DefaultFetchSize, primaryKeys)
		}
		return m, nil
	}

	return m, nil
}

// handleSQLEditorKeys handles key events in the SQL editor dialog
func handleSQLEditorKeys(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		// Close dialog without executing
		m.SQLEditorVisible = false
		return m, nil

	case tea.KeyCtrlR:
		// Execute SQL
		m.SQLEditorVisible = false
		m.CurrentSQL = m.EditSQL
		m.CustomSQL = true

		// Reset data state
		m.SelectedDataRow = 0
		m.ViewportOffset = 0
		m.HorizontalOffset = 0

		// Get current table name
		if m.SelectedTable >= 0 && m.SelectedTable < len(m.Tables) {
			tableName := m.Tables[m.SelectedTable]
			return m, db.ExecuteCustomSQL(m.NosqlClient, tableName, m.EditSQL, ui.DefaultFetchSize)
		}
		return m, nil

	case tea.KeyEnter:
		// Insert newline
		m.EditSQL, m.SQLCursorPos = ui.InsertWithCursor(m.EditSQL, m.SQLCursorPos, "\n")
		return m, nil

	case tea.KeyBackspace:
		m.EditSQL, m.SQLCursorPos = ui.Backspace(m.EditSQL, m.SQLCursorPos)
		return m, nil

	case tea.KeyDelete:
		m.EditSQL = ui.DeleteAt(m.EditSQL, m.SQLCursorPos)
		return m, nil

	case tea.KeyLeft:
		if m.SQLCursorPos > 0 {
			m.SQLCursorPos--
		}
		return m, nil

	case tea.KeyRight:
		if m.SQLCursorPos < len(m.EditSQL) {
			m.SQLCursorPos++
		}
		return m, nil

	case tea.KeyHome, tea.KeyCtrlA:
		m.SQLCursorPos = 0
		return m, nil

	case tea.KeyEnd, tea.KeyCtrlE:
		m.SQLCursorPos = len(m.EditSQL)
		return m, nil

	case tea.KeySpace:
		m.EditSQL, m.SQLCursorPos = ui.InsertWithCursor(m.EditSQL, m.SQLCursorPos, " ")
		return m, nil

	case tea.KeyRunes:
		for _, r := range msg.Runes {
			m.EditSQL, m.SQLCursorPos = ui.InsertWithCursor(m.EditSQL, m.SQLCursorPos, string(r))
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
		m.SelectedTable = 0 // Initialize selection to first table
		// Fetch details for first table
		return m, db.FetchTableDetails(m.NosqlClient, m.Tables[0])
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
		// TODO: Show error
		return m, nil
	}

	m.TableDetails[msg.TableName] = &msg
	return m, nil
}

func handleTableDataResult(m Model, msg db.TableDataResult) (Model, tea.Cmd) {
	if msg.Err != nil {
		// TODO: Show error
		m.LoadingData = false
		return m, nil
	}

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
