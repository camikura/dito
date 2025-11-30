package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/camikura/dito/internal/db"
	"github.com/camikura/dito/internal/ui"
)

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
	case "ctrl+c":
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
			m.CurrentSQL = ""
			m.SQLCursorPos = 0
			// Clear all cached data
			m.TableDetails = make(map[string]*db.TableDetailsResult)
			m.TableData = make(map[string]*db.TableDataResult)
		}
		return m, nil
	}

	return m, nil
}

func handleConnectionDialogKeys(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
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
		if m.ConnectionDialogField == 0 {
			m.EditEndpoint, m.EditCursorPos = ui.Backspace(m.EditEndpoint, m.EditCursorPos)
		} else {
			m.EditPort, m.EditCursorPos = ui.Backspace(m.EditPort, m.EditCursorPos)
		}
		return m, nil

	case tea.KeyDelete:
		if m.ConnectionDialogField == 0 {
			m.EditEndpoint = ui.DeleteAt(m.EditEndpoint, m.EditCursorPos)
		} else {
			m.EditPort = ui.DeleteAt(m.EditPort, m.EditCursorPos)
		}
		return m, nil

	case tea.KeyLeft:
		if m.EditCursorPos > 0 {
			m.EditCursorPos--
		}
		return m, nil

	case tea.KeyRight:
		maxPos := ui.RuneLen(m.EditEndpoint)
		if m.ConnectionDialogField == 1 {
			maxPos = ui.RuneLen(m.EditPort)
		}
		if m.EditCursorPos < maxPos {
			m.EditCursorPos++
		}
		return m, nil

	case tea.KeyHome:
		m.EditCursorPos = 0
		return m, nil

	case tea.KeyEnd:
		if m.ConnectionDialogField == 0 {
			m.EditCursorPos = ui.RuneLen(m.EditEndpoint)
		} else {
			m.EditCursorPos = ui.RuneLen(m.EditPort)
		}
		return m, nil

	case tea.KeyRunes:
		char := string(msg.Runes)
		if m.ConnectionDialogField == 0 {
			m.EditEndpoint, m.EditCursorPos = ui.InsertWithCursor(m.EditEndpoint, m.EditCursorPos, char)
		} else {
			m.EditPort, m.EditCursorPos = ui.InsertWithCursor(m.EditPort, m.EditCursorPos, char)
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
				return m, db.FetchTableData(m.NosqlClient, tableName, ui.DefaultFetchSize, primaryKeys)
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
	if m.SelectedTable >= 0 && m.SelectedTable < len(m.Tables) {
		schemaTableName = m.Tables[m.SelectedTable]
	} else if m.CursorTable < len(m.Tables) {
		schemaTableName = m.Tables[m.CursorTable]
	}

	if schemaTableName == "" {
		return m, nil
	}

	// Calculate max scroll (dynamic based on content)
	maxScroll := 0
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

	switch msg.String() {
	case "up", "k":
		if m.SchemaScrollOffset > 0 {
			m.SchemaScrollOffset--
		}
		return m, nil

	case "down", "j":
		if m.SchemaScrollOffset < maxScroll {
			m.SchemaScrollOffset++
		}
		return m, nil
	}

	return m, nil
}

func handleSQLKeys(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlR:
		// Execute SQL
		if !m.Connected || m.CurrentSQL == "" {
			return m, nil
		}

		// Parse table name from SQL
		tableName := ui.ExtractTableNameFromSQL(m.CurrentSQL)
		// Use case-insensitive table name matching
		actualTableName := m.FindTableName(tableName)
		if actualTableName != "" {
			tableName = actualTableName
		}

		// Check if this is a custom SQL (not the default SELECT * FROM table)
		tableIndex := m.FindTableIndex(tableName)
		if tableIndex >= 0 {
			// Get DDL to generate proper default SQL with ORDER BY
			var ddl string
			if details, exists := m.TableDetails[tableName]; exists && details != nil && details.Schema != nil {
				ddl = details.Schema.DDL
			}
			defaultSQL := buildDefaultSQL(tableName, ddl)
			if m.CurrentSQL != defaultSQL {
				// This is custom SQL
				m.CustomSQL = true
				// Parse column order from SQL
				m.ColumnOrder = db.ParseSelectColumns(m.CurrentSQL)
				// Save current SelectedTable for later restoration
				if !m.CustomSQL || m.PreviousSelectedTable == -1 {
					m.PreviousSelectedTable = m.SelectedTable
				}
				// Update SelectedTable to match the table in SQL
				m.SelectedTable = tableIndex
			} else {
				// This is standard SQL
				m.CustomSQL = false
				m.ColumnOrder = nil
			}
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
			if m.SelectedTable >= 0 && m.SelectedTable < len(m.Tables) {
				tableName := m.Tables[m.SelectedTable]
				if data, exists := m.TableData[tableName]; exists && data != nil {
					remainingRows := totalRows - m.SelectedDataRow - 1
					if remainingRows <= ui.FetchMoreThreshold && data.HasMore && !m.LoadingData {
						// Custom SQL uses OFFSET pagination
						if data.IsCustomSQL && data.CurrentSQL != "" {
							m.LoadingData = true
							return m, db.FetchMoreCustomSQL(m.NosqlClient, tableName, data.CurrentSQL, ui.DefaultFetchSize, data.Offset)
						}
						// Standard queries use PRIMARY KEY cursor pagination
						if data.LastPKValues != nil {
							m.LoadingData = true
							var primaryKeys []string
							if details, exists := m.TableDetails[tableName]; exists && details != nil && details.Schema != nil && details.Schema.DDL != "" {
								primaryKeys = ui.ParsePrimaryKeysFromDDL(details.Schema.DDL)
							}
							return m, db.FetchMoreTableData(m.NosqlClient, tableName, ui.DefaultFetchSize, primaryKeys, data.LastPKValues)
						}
					}
				}
			}
		}
		return m, nil

	case "left", "h":
		if m.HorizontalOffset > 0 {
			m.HorizontalOffset--
		}
		return m, nil

	case "right", "l":
		if m.HorizontalOffset < maxHorizontalOffset {
			m.HorizontalOffset++
		}
		return m, nil

	case "enter":
		// Show record detail dialog
		if totalRows > 0 && m.SelectedDataRow < totalRows {
			m.RecordDetailVisible = true
			m.RecordDetailScroll = 0
		}
		return m, nil

	case "esc":
		// Reset to default SQL (only if custom SQL is active)
		if m.CustomSQL {
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

func handleRecordDetailKeys(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Calculate max scroll for record detail
	maxScroll := calculateRecordDetailMaxScroll(m)

	switch msg.String() {
	case "esc":
		m.RecordDetailVisible = false
		m.RecordDetailScroll = 0
		return m, nil

	case "up", "k":
		if m.RecordDetailScroll > 0 {
			m.RecordDetailScroll--
		}
		return m, nil

	case "down", "j":
		if m.RecordDetailScroll < maxScroll {
			m.RecordDetailScroll++
		}
		return m, nil

	case "home":
		m.RecordDetailScroll = 0
		return m, nil

	case "end":
		m.RecordDetailScroll = maxScroll
		return m, nil

	case "pgup":
		// Scroll up by page
		m.RecordDetailScroll -= ui.PageScrollAmount
		if m.RecordDetailScroll < 0 {
			m.RecordDetailScroll = 0
		}
		return m, nil

	case "pgdown":
		// Scroll down by page
		m.RecordDetailScroll += ui.PageScrollAmount
		if m.RecordDetailScroll > maxScroll {
			m.RecordDetailScroll = maxScroll
		}
		return m, nil
	}

	return m, nil
}
