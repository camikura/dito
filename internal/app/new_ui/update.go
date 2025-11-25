package new_ui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/camikura/dito/internal/db"
)

// Update handles messages and updates the model
func Update(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
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
	case FocusPaneData:
		return handleDataKeys(m, msg)
	}

	return m, nil
}

func handleTablesKeys(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.CursorTable > 0 {
			m.CursorTable--
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

			// Get primary keys from schema if available
			var primaryKeys []string
			if details, exists := m.TableDetails[tableName]; exists && details != nil && details.Schema != nil {
				// We'll parse from DDL in a moment - for now just pass empty slice
				primaryKeys = []string{}
			}

			// Load table data
			return m, db.FetchTableData(m.NosqlClient, tableName, 1000, primaryKeys)
		}
		return m, nil
	}

	return m, nil
}

func handleDataKeys(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.SelectedDataRow > 0 {
			m.SelectedDataRow--
		}
		return m, nil

	case "down", "j":
		// TODO: Check row count
		m.SelectedDataRow++
		return m, nil
	}

	return m, nil
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

	m.Tables = msg.Tables
	if len(m.Tables) > 0 {
		m.CursorTable = 0
		// Fetch details for first table
		return m, db.FetchTableDetails(m.NosqlClient, m.Tables[0])
	}

	return m, nil
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
		return m, nil
	}

	m.TableData[msg.TableName] = &msg
	m.LoadingData = false
	return m, nil
}
