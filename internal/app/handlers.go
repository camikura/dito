package app

import (
	"time"

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
	case "ctrl+q":
		// Quit confirmation: first press shows message, second press quits
		if m.QuitConfirmation {
			return m, tea.Quit
		}
		m.QuitConfirmation = true
		return m, tea.Tick(ui.QuitConfirmationTimeout, func(_ time.Time) tea.Msg {
			return clearQuitConfirmationMsg{}
		})

	case "ctrl+c":
		// In data pane, Ctrl+C copies selected row
		if m.CurrentPane == FocusPaneData {
			return handleDataCopy(m)
		}
		return m, nil

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
	// Calculate tables pane height using the same logic as view.go
	visibleLines := calculateTablesHeight(m)

	// Handle M-< and M-> (Alt+Shift+, and Alt+Shift+.)
	// On Mac, these produce special characters: ¯ (175) and ˘ (728)
	switch msg.String() {
	case "alt+<", "¯":
		// Jump to first table
		m.CursorTable = 0
		m.TablesScrollOffset = 0
		return m, nil

	case "alt+>", "˘":
		// Jump to last table
		if len(m.Tables) > 0 {
			m.CursorTable = len(m.Tables) - 1
			// Adjust scroll offset to show cursor
			if m.CursorTable >= visibleLines {
				m.TablesScrollOffset = m.CursorTable - visibleLines + 1
			}
		}
		return m, nil
	}

	switch msg.Type {
	case tea.KeyUp, tea.KeyCtrlP:
		if m.CursorTable > 0 {
			m.CursorTable--

			// Adjust scroll offset to keep cursor visible
			if m.CursorTable < m.TablesScrollOffset {
				m.TablesScrollOffset = m.CursorTable
			}
		}
		return m, nil

	case tea.KeyDown, tea.KeyCtrlN:
		if m.CursorTable < len(m.Tables)-1 {
			m.CursorTable++

			// Adjust scroll offset to keep cursor visible
			if m.CursorTable >= m.TablesScrollOffset+visibleLines {
				m.TablesScrollOffset = m.CursorTable - visibleLines + 1
			}
		}
		return m, nil

	case tea.KeyEnter:
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

			// Fetch ancestor table schemas if not already loaded (for inherited columns display)
			var ancestorCmds []tea.Cmd
			ancestors := ui.GetAncestorTableNames(tableName)
			for _, ancestor := range ancestors {
				if _, exists := m.TableDetails[ancestor]; !exists {
					ancestorCmds = append(ancestorCmds, db.FetchTableDetails(m.NosqlClient, ancestor))
				}
			}

			// Check if schema is already loaded
			if details, exists := m.TableDetails[tableName]; exists && details != nil && details.Schema != nil {
				// Schema available - fetch data with ORDER BY
				ddl := details.Schema.DDL
				primaryKeys := ui.ParsePrimaryKeysFromDDL(ddl)
				m.CurrentSQL = buildDefaultSQL(tableName, ddl)
				m.SQLCursorPos = ui.RuneLen(m.CurrentSQL)
				dataCmd := db.FetchTableData(m.NosqlClient, tableName, ui.DefaultFetchSize, primaryKeys)
				if len(ancestorCmds) > 0 {
					ancestorCmds = append(ancestorCmds, dataCmd)
					return m, tea.Batch(ancestorCmds...)
				}
				return m, dataCmd
			}

			// Schema not loaded - fetch schema first, data will be fetched when schema arrives
			m.CurrentSQL = "SELECT * FROM " + tableName
			m.SQLCursorPos = ui.RuneLen(m.CurrentSQL)
			m.LoadingData = true
			ancestorCmds = append(ancestorCmds, db.FetchTableDetails(m.NosqlClient, tableName))
			return m, tea.Batch(ancestorCmds...)
		}
		return m, nil
	}

	return m, nil
}

func handleSchemaKeys(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Determine which table schema is displayed (same logic as view)
	schemaTableName := m.SelectedTableName()
	if schemaTableName == "" {
		schemaTableName = m.CursorTableName()
	}
	if schemaTableName == "" {
		return m, nil
	}

	// Calculate schema pane height and content line count
	schemaHeight := calculateSchemaHeight(m)
	lineCount := calculateSchemaContentLineCount(m, schemaTableName)

	// Max scroll = total lines - visible lines
	maxScroll := lineCount - schemaHeight
	if maxScroll < 0 {
		maxScroll = 0
	}

	// Handle M-< and M-> (Alt+Shift+, and Alt+Shift+.)
	// On Mac, these produce special characters: ¯ (175) and ˘ (728)
	switch msg.String() {
	case "alt+<", "¯":
		// Scroll to top
		m.SchemaScrollOffset = 0
		return m, nil

	case "alt+>", "˘":
		// Scroll to bottom
		m.SchemaScrollOffset = maxScroll
		return m, nil
	}

	switch msg.Type {
	case tea.KeyUp, tea.KeyCtrlP:
		if m.SchemaScrollOffset > 0 {
			m.SchemaScrollOffset--
		}
		return m, nil

	case tea.KeyDown, tea.KeyCtrlN:
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

		// Ctrl+R always executes as custom SQL
		tableIndex := m.FindTableIndex(tableName)
		if tableIndex >= 0 {
			m.CustomSQL = true
			// Parse column order from SQL
			m.ColumnOrder = db.ParseSelectColumns(m.CurrentSQL)
			// Save current SelectedTable for later restoration
			if m.PreviousSelectedTable == -1 {
				m.PreviousSelectedTable = m.SelectedTable
			}
			// Update SelectedTable to match the table in SQL
			m.SelectedTable = tableIndex
		}

		// Fall back to selected table if no table name in SQL
		if tableName == "" {
			tableName = m.SelectedTableName()
		}

		if tableName != "" {
			var cmds []tea.Cmd

			// Fetch schema (always try, even for unknown tables to get error)
			if _, exists := m.TableDetails[tableName]; !exists {
				cmds = append(cmds, db.FetchTableDetails(m.NosqlClient, tableName))
			}

			// Execute custom SQL
			cmds = append(cmds, db.ExecuteCustomSQL(m.NosqlClient, tableName, m.CurrentSQL, ui.DefaultFetchSize))

			// Reset data row selection to top
			m.SelectedDataRow = 0
			m.ViewportOffset = 0

			// Move focus to Data pane
			m.CurrentPane = FocusPaneData

			return m, tea.Batch(cmds...)
		}
		return m, nil

	case tea.KeyEnter:
		// Insert newline
		m.CurrentSQL, m.SQLCursorPos = ui.InsertWithCursor(m.CurrentSQL, m.SQLCursorPos, "\n")
		m.SQLScrollOffset = updateSQLScrollOffset(m)
		return m, nil

	case tea.KeyBackspace:
		m.CurrentSQL, m.SQLCursorPos = ui.Backspace(m.CurrentSQL, m.SQLCursorPos)
		m.SQLScrollOffset = updateSQLScrollOffset(m)
		return m, nil

	case tea.KeyDelete:
		m.CurrentSQL = ui.DeleteAt(m.CurrentSQL, m.SQLCursorPos)
		m.SQLScrollOffset = updateSQLScrollOffset(m)
		return m, nil

	case tea.KeyLeft, tea.KeyCtrlB:
		// Move cursor left (Ctrl+B: Emacs style)
		if m.SQLCursorPos > 0 {
			m.SQLCursorPos--
			m.SQLScrollOffset = updateSQLScrollOffset(m)
		}
		return m, nil

	case tea.KeyRight, tea.KeyCtrlF:
		// Move cursor right (Ctrl+F: Emacs style)
		if m.SQLCursorPos < ui.RuneLen(m.CurrentSQL) {
			m.SQLCursorPos++
			m.SQLScrollOffset = updateSQLScrollOffset(m)
		}
		return m, nil

	case tea.KeyUp, tea.KeyCtrlP:
		// Move cursor up one line (Ctrl+P: Emacs style)
		m.SQLCursorPos = moveCursorUpInText(m.CurrentSQL, m.SQLCursorPos)
		m.SQLScrollOffset = updateSQLScrollOffset(m)
		return m, nil

	case tea.KeyDown, tea.KeyCtrlN:
		// Move cursor down one line (Ctrl+N: Emacs style)
		m.SQLCursorPos = moveCursorDownInText(m.CurrentSQL, m.SQLCursorPos)
		m.SQLScrollOffset = updateSQLScrollOffset(m)
		return m, nil

	case tea.KeyHome:
		m.SQLCursorPos = 0
		m.SQLScrollOffset = updateSQLScrollOffset(m)
		return m, nil

	case tea.KeyEnd:
		m.SQLCursorPos = ui.RuneLen(m.CurrentSQL)
		m.SQLScrollOffset = updateSQLScrollOffset(m)
		return m, nil

	case tea.KeyCtrlA:
		// Emacs: move to beginning of current line
		m.SQLCursorPos = moveCursorToLineStart(m.CurrentSQL, m.SQLCursorPos)
		m.SQLScrollOffset = updateSQLScrollOffset(m)
		return m, nil

	case tea.KeyCtrlE:
		// Emacs: move to end of current line
		m.SQLCursorPos = moveCursorToLineEnd(m.CurrentSQL, m.SQLCursorPos)
		m.SQLScrollOffset = updateSQLScrollOffset(m)
		return m, nil

	case tea.KeySpace:
		m.CurrentSQL, m.SQLCursorPos = ui.InsertWithCursor(m.CurrentSQL, m.SQLCursorPos, " ")
		m.SQLScrollOffset = updateSQLScrollOffset(m)
		return m, nil

	case tea.KeyRunes:
		for _, r := range msg.Runes {
			m.CurrentSQL, m.SQLCursorPos = ui.InsertWithCursor(m.CurrentSQL, m.SQLCursorPos, string(r))
		}
		m.SQLScrollOffset = updateSQLScrollOffset(m)
		return m, nil
	}

	return m, nil
}

// clearCopyMessageMsg is sent to clear the copy message
type clearCopyMessageMsg struct{}

// clearQuitConfirmationMsg is sent to clear the quit confirmation state
type clearQuitConfirmationMsg struct{}

// handleMouseClick handles mouse click events for pane selection
func handleMouseClick(m Model, msg tea.MouseMsg) (Model, tea.Cmd) {
	// Only handle left click (button press)
	if msg.Button != tea.MouseButtonLeft || msg.Action != tea.MouseActionPress {
		return m, nil
	}

	// Ignore if dialogs are visible
	if m.ConnectionDialogVisible || m.RecordDetailVisible {
		return m, nil
	}

	x, y := msg.X, msg.Y

	// Calculate pane boundaries
	leftPaneWidth := ui.LeftPaneContentWidth

	// If click is in the right pane (Data pane)
	if x >= leftPaneWidth {
		m.CurrentPane = FocusPaneData
		return m, nil
	}

	// Click is in the left pane area
	// Need to calculate vertical boundaries for each pane
	// Connection pane is at the top, height varies but typically 5 lines
	connectionHeight := m.ConnectionPaneHeight
	if connectionHeight == 0 {
		connectionHeight = ui.DefaultConnectionPaneHeight
	}

	// Calculate pane boundaries (including borders)
	connectionEnd := connectionHeight
	tablesEnd := connectionEnd + m.TablesHeight + ui.PaneBorderHeight
	schemaEnd := tablesEnd + m.SchemaHeight + ui.PaneBorderHeight
	sqlEnd := schemaEnd + m.SQLHeight + ui.PaneBorderHeight

	// Determine which pane was clicked based on Y position
	if y < connectionEnd {
		m.CurrentPane = FocusPaneConnection
	} else if y < tablesEnd {
		m.CurrentPane = FocusPaneTables
	} else if y < schemaEnd {
		m.CurrentPane = FocusPaneSchema
	} else if y < sqlEnd {
		m.CurrentPane = FocusPaneSQL
	}

	return m, nil
}
