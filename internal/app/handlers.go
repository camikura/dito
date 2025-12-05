package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/camikura/dito/internal/db"
	"github.com/camikura/dito/internal/ui"
)

func handleKeyPress(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Connection dialog takes precedence
	if m.ConnectionDialog.Visible {
		return handleConnectionDialogKeys(m, msg)
	}

	// Record detail dialog takes precedence
	if m.RecordDetail.Visible {
		return handleRecordDetailKeys(m, msg)
	}

	switch msg.String() {
	case "ctrl+q":
		// Quit confirmation: first press shows message, second press quits
		if m.UI.QuitConfirmation {
			return m, tea.Quit
		}
		m.UI.QuitConfirmation = true
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
		if m.Connection.Connected {
			m = m.NextPane()
		}
		return m, nil

	case "shift+tab":
		// Only allow pane switching when connected
		if m.Connection.Connected {
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
		m.ConnectionDialog.Visible = true
		m.ConnectionDialog.Field = 0
		// Initialize with current values or defaults
		if m.ConnectionDialog.EditEndpoint == "" {
			m.ConnectionDialog.EditEndpoint = "localhost"
		}
		if m.ConnectionDialog.EditPort == "" {
			m.ConnectionDialog.EditPort = "8080"
		}
		m.ConnectionDialog.EditCursorPos = ui.RuneLen(m.ConnectionDialog.EditEndpoint)
		return m, nil

	case "ctrl+d":
		// Disconnect
		if m.Connection.Connected {
			m.Connection.Connected = false
			m.Connection.NosqlClient = nil
			m.Tables.Tables = []string{}
			m.Tables.SelectedTable = -1
			m.Tables.CursorTable = 0
			m.SQL.CurrentSQL = ""
			m.SQL.CursorPos = 0
			// Clear all cached data
			m.Schema.TableDetails = make(map[string]*db.TableDetailsResult)
			m.Data.TableData = make(map[string]*db.TableDataResult)
		}
		return m, nil
	}

	return m, nil
}

func handleConnectionDialogKeys(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Helper to move to field and set cursor position
	moveToField := func(field int) {
		m.ConnectionDialog.Field = field
		if field == 0 {
			m.ConnectionDialog.EditCursorPos = ui.RuneLen(m.ConnectionDialog.EditEndpoint)
		} else if field == 1 {
			m.ConnectionDialog.EditCursorPos = ui.RuneLen(m.ConnectionDialog.EditPort)
		}
	}

	switch msg.Type {
	case tea.KeyEsc:
		// Close dialog
		m.ConnectionDialog.Visible = false
		return m, nil

	case tea.KeyEnter:
		// Connect from any field
		m.ConnectionDialog.Visible = false
		m.Connection.Endpoint = m.ConnectionDialog.EditEndpoint + ":" + m.ConnectionDialog.EditPort
		return m, db.Connect(m.ConnectionDialog.EditEndpoint, m.ConnectionDialog.EditPort, false)

	case tea.KeyTab, tea.KeyDown:
		moveToField((m.ConnectionDialog.Field + 1) % 2)
		return m, nil

	case tea.KeyShiftTab, tea.KeyUp:
		moveToField((m.ConnectionDialog.Field + 1) % 2)
		return m, nil

	case tea.KeyBackspace:
		if m.ConnectionDialog.Field == 0 {
			m.ConnectionDialog.EditEndpoint, m.ConnectionDialog.EditCursorPos = ui.Backspace(m.ConnectionDialog.EditEndpoint, m.ConnectionDialog.EditCursorPos)
		} else {
			m.ConnectionDialog.EditPort, m.ConnectionDialog.EditCursorPos = ui.Backspace(m.ConnectionDialog.EditPort, m.ConnectionDialog.EditCursorPos)
		}
		return m, nil

	case tea.KeyDelete:
		if m.ConnectionDialog.Field == 0 {
			m.ConnectionDialog.EditEndpoint = ui.DeleteAt(m.ConnectionDialog.EditEndpoint, m.ConnectionDialog.EditCursorPos)
		} else {
			m.ConnectionDialog.EditPort = ui.DeleteAt(m.ConnectionDialog.EditPort, m.ConnectionDialog.EditCursorPos)
		}
		return m, nil

	case tea.KeyLeft:
		if m.ConnectionDialog.EditCursorPos > 0 {
			m.ConnectionDialog.EditCursorPos--
		}
		return m, nil

	case tea.KeyRight:
		maxPos := ui.RuneLen(m.ConnectionDialog.EditEndpoint)
		if m.ConnectionDialog.Field == 1 {
			maxPos = ui.RuneLen(m.ConnectionDialog.EditPort)
		}
		if m.ConnectionDialog.EditCursorPos < maxPos {
			m.ConnectionDialog.EditCursorPos++
		}
		return m, nil

	case tea.KeyHome:
		m.ConnectionDialog.EditCursorPos = 0
		return m, nil

	case tea.KeyEnd:
		if m.ConnectionDialog.Field == 0 {
			m.ConnectionDialog.EditCursorPos = ui.RuneLen(m.ConnectionDialog.EditEndpoint)
		} else {
			m.ConnectionDialog.EditCursorPos = ui.RuneLen(m.ConnectionDialog.EditPort)
		}
		return m, nil

	case tea.KeyRunes:
		char := string(msg.Runes)
		if m.ConnectionDialog.Field == 0 {
			m.ConnectionDialog.EditEndpoint, m.ConnectionDialog.EditCursorPos = ui.InsertWithCursor(m.ConnectionDialog.EditEndpoint, m.ConnectionDialog.EditCursorPos, char)
		} else {
			m.ConnectionDialog.EditPort, m.ConnectionDialog.EditCursorPos = ui.InsertWithCursor(m.ConnectionDialog.EditPort, m.ConnectionDialog.EditCursorPos, char)
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
		m.Tables.CursorTable = 0
		m.Tables.ScrollOffset = 0
		return m, nil

	case "alt+>", "˘":
		// Jump to last table
		if len(m.Tables.Tables) > 0 {
			m.Tables.CursorTable = len(m.Tables.Tables) - 1
			// Adjust scroll offset to show cursor
			if m.Tables.CursorTable >= visibleLines {
				m.Tables.ScrollOffset = m.Tables.CursorTable - visibleLines + 1
			}
		}
		return m, nil
	}

	switch msg.Type {
	case tea.KeyUp, tea.KeyCtrlP:
		if m.Tables.CursorTable > 0 {
			m.Tables.CursorTable--

			// Adjust scroll offset to keep cursor visible
			if m.Tables.CursorTable < m.Tables.ScrollOffset {
				m.Tables.ScrollOffset = m.Tables.CursorTable
			}
		}
		return m, nil

	case tea.KeyDown, tea.KeyCtrlN:
		if m.Tables.CursorTable < len(m.Tables.Tables)-1 {
			m.Tables.CursorTable++

			// Adjust scroll offset to keep cursor visible
			if m.Tables.CursorTable >= m.Tables.ScrollOffset+visibleLines {
				m.Tables.ScrollOffset = m.Tables.CursorTable - visibleLines + 1
			}
		}
		return m, nil

	case tea.KeyEnter:
		// Select table and load data (only on Enter)
		if m.Tables.CursorTable < len(m.Tables.Tables) {
			m.Tables.SelectedTable = m.Tables.CursorTable
			tableName := m.Tables.Tables[m.Tables.SelectedTable]

			// Reset state
			m.SQL.CustomSQL = false
			m.Data.SelectedDataRow = 0
			m.Data.ViewportOffset = 0
			m.Data.HorizontalOffset = 0
			m.Schema.ScrollOffset = 0

			// Move focus to Data pane for immediate interaction
			m.CurrentPane = FocusPaneData

			// Fetch ancestor table schemas if not already loaded (for inherited columns display)
			var ancestorCmds []tea.Cmd
			ancestors := ui.GetAncestorTableNames(tableName)
			for _, ancestor := range ancestors {
				if _, exists := m.Schema.TableDetails[ancestor]; !exists {
					ancestorCmds = append(ancestorCmds, db.FetchTableDetails(m.Connection.NosqlClient, ancestor))
				}
			}

			// Check if schema is already loaded
			if details, exists := m.Schema.TableDetails[tableName]; exists && details != nil && details.Schema != nil {
				// Schema available - fetch data with ORDER BY
				ddl := details.Schema.DDL
				primaryKeys := ui.ParsePrimaryKeysFromDDL(ddl)
				m.SQL.CurrentSQL = buildDefaultSQL(tableName, ddl)
				m.SQL.CursorPos = ui.RuneLen(m.SQL.CurrentSQL)
				dataCmd := db.FetchTableData(m.Connection.NosqlClient, tableName, ui.DefaultFetchSize, primaryKeys)
				if len(ancestorCmds) > 0 {
					ancestorCmds = append(ancestorCmds, dataCmd)
					return m, tea.Batch(ancestorCmds...)
				}
				return m, dataCmd
			}

			// Schema not loaded - fetch schema first, data will be fetched when schema arrives
			m.SQL.CurrentSQL = "SELECT * FROM " + tableName
			m.SQL.CursorPos = ui.RuneLen(m.SQL.CurrentSQL)
			m.Data.LoadingData = true
			ancestorCmds = append(ancestorCmds, db.FetchTableDetails(m.Connection.NosqlClient, tableName))
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
		m.Schema.ScrollOffset = 0
		return m, nil

	case "alt+>", "˘":
		// Scroll to bottom
		m.Schema.ScrollOffset = maxScroll
		return m, nil
	}

	switch msg.Type {
	case tea.KeyUp, tea.KeyCtrlP:
		if m.Schema.ScrollOffset > 0 {
			m.Schema.ScrollOffset--
		}
		return m, nil

	case tea.KeyDown, tea.KeyCtrlN:
		if m.Schema.ScrollOffset < maxScroll {
			m.Schema.ScrollOffset++
		}
		return m, nil
	}

	return m, nil
}

func handleSQLKeys(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlR:
		// Execute SQL
		if !m.Connection.Connected || m.SQL.CurrentSQL == "" {
			return m, nil
		}

		// Parse table name from SQL
		tableName := ui.ExtractTableNameFromSQL(m.SQL.CurrentSQL)
		// Use case-insensitive table name matching
		actualTableName := m.FindTableName(tableName)
		if actualTableName != "" {
			tableName = actualTableName
		}

		// Ctrl+R always executes as custom SQL
		tableIndex := m.FindTableIndex(tableName)
		if tableIndex >= 0 {
			m.SQL.CustomSQL = true
			// Parse column order from SQL
			m.SQL.ColumnOrder = db.ParseSelectColumns(m.SQL.CurrentSQL)
			// Save current SelectedTable for later restoration
			if m.SQL.PreviousSelectedTable == -1 {
				m.SQL.PreviousSelectedTable = m.Tables.SelectedTable
			}
			// Update SelectedTable to match the table in SQL
			m.Tables.SelectedTable = tableIndex
		}

		// Fall back to selected table if no table name in SQL
		if tableName == "" {
			tableName = m.SelectedTableName()
		}

		if tableName != "" {
			var cmds []tea.Cmd

			// Fetch schema (always try, even for unknown tables to get error)
			if _, exists := m.Schema.TableDetails[tableName]; !exists {
				cmds = append(cmds, db.FetchTableDetails(m.Connection.NosqlClient, tableName))
			}

			// Execute custom SQL
			cmds = append(cmds, db.ExecuteCustomSQL(m.Connection.NosqlClient, tableName, m.SQL.CurrentSQL, ui.DefaultFetchSize))

			// Reset data row selection to top
			m.Data.SelectedDataRow = 0
			m.Data.ViewportOffset = 0

			// Move focus to Data pane
			m.CurrentPane = FocusPaneData

			return m, tea.Batch(cmds...)
		}
		return m, nil

	case tea.KeyEnter:
		// Insert newline
		m.SQL.CurrentSQL, m.SQL.CursorPos = ui.InsertWithCursor(m.SQL.CurrentSQL, m.SQL.CursorPos, "\n")
		m.SQL.ScrollOffset = updateSQLScrollOffset(m)
		return m, nil

	case tea.KeyBackspace:
		m.SQL.CurrentSQL, m.SQL.CursorPos = ui.Backspace(m.SQL.CurrentSQL, m.SQL.CursorPos)
		m.SQL.ScrollOffset = updateSQLScrollOffset(m)
		return m, nil

	case tea.KeyDelete:
		m.SQL.CurrentSQL = ui.DeleteAt(m.SQL.CurrentSQL, m.SQL.CursorPos)
		m.SQL.ScrollOffset = updateSQLScrollOffset(m)
		return m, nil

	case tea.KeyLeft, tea.KeyCtrlB:
		// Move cursor left (Ctrl+B: Emacs style)
		if m.SQL.CursorPos > 0 {
			m.SQL.CursorPos--
			m.SQL.ScrollOffset = updateSQLScrollOffset(m)
		}
		return m, nil

	case tea.KeyRight, tea.KeyCtrlF:
		// Move cursor right (Ctrl+F: Emacs style)
		if m.SQL.CursorPos < ui.RuneLen(m.SQL.CurrentSQL) {
			m.SQL.CursorPos++
			m.SQL.ScrollOffset = updateSQLScrollOffset(m)
		}
		return m, nil

	case tea.KeyUp, tea.KeyCtrlP:
		// Move cursor up one line (Ctrl+P: Emacs style)
		m.SQL.CursorPos = moveCursorUpInText(m.SQL.CurrentSQL, m.SQL.CursorPos)
		m.SQL.ScrollOffset = updateSQLScrollOffset(m)
		return m, nil

	case tea.KeyDown, tea.KeyCtrlN:
		// Move cursor down one line (Ctrl+N: Emacs style)
		m.SQL.CursorPos = moveCursorDownInText(m.SQL.CurrentSQL, m.SQL.CursorPos)
		m.SQL.ScrollOffset = updateSQLScrollOffset(m)
		return m, nil

	case tea.KeyHome:
		m.SQL.CursorPos = 0
		m.SQL.ScrollOffset = updateSQLScrollOffset(m)
		return m, nil

	case tea.KeyEnd:
		m.SQL.CursorPos = ui.RuneLen(m.SQL.CurrentSQL)
		m.SQL.ScrollOffset = updateSQLScrollOffset(m)
		return m, nil

	case tea.KeyCtrlA:
		// Emacs: move to beginning of current line
		m.SQL.CursorPos = moveCursorToLineStart(m.SQL.CurrentSQL, m.SQL.CursorPos)
		m.SQL.ScrollOffset = updateSQLScrollOffset(m)
		return m, nil

	case tea.KeyCtrlE:
		// Emacs: move to end of current line
		m.SQL.CursorPos = moveCursorToLineEnd(m.SQL.CurrentSQL, m.SQL.CursorPos)
		m.SQL.ScrollOffset = updateSQLScrollOffset(m)
		return m, nil

	case tea.KeySpace:
		m.SQL.CurrentSQL, m.SQL.CursorPos = ui.InsertWithCursor(m.SQL.CurrentSQL, m.SQL.CursorPos, " ")
		m.SQL.ScrollOffset = updateSQLScrollOffset(m)
		return m, nil

	case tea.KeyRunes:
		for _, r := range msg.Runes {
			m.SQL.CurrentSQL, m.SQL.CursorPos = ui.InsertWithCursor(m.SQL.CurrentSQL, m.SQL.CursorPos, string(r))
		}
		m.SQL.ScrollOffset = updateSQLScrollOffset(m)
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
	if m.ConnectionDialog.Visible || m.RecordDetail.Visible {
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
	connectionHeight := m.Window.ConnectionPaneHeight
	if connectionHeight == 0 {
		connectionHeight = ui.DefaultConnectionPaneHeight
	}

	// Calculate pane boundaries (including borders)
	connectionEnd := connectionHeight
	tablesEnd := connectionEnd + m.Window.TablesHeight + ui.PaneBorderHeight
	schemaEnd := tablesEnd + m.Window.SchemaHeight + ui.PaneBorderHeight
	sqlEnd := schemaEnd + m.Window.SQLHeight + ui.PaneBorderHeight

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
