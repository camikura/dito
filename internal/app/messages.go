package app

import (
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/camikura/dito/internal/db"
	"github.com/camikura/dito/internal/ui"
)

func handleConnectionResult(m Model, msg db.ConnectionResult) (Model, tea.Cmd) {
	if msg.Err != nil {
		// Connection failed
		m.Connection.Connected = false
		m.Connection.Message = msg.Err.Error()
		return m, nil
	}

	// Connection successful
	m.Connection.Connected = true
	m.Connection.NosqlClient = msg.Client
	m.Connection.Endpoint = msg.Endpoint
	m.Connection.Message = ""

	// Fetch table list
	return m, db.FetchTables(msg.Client)
}

func handleTableListResult(m Model, msg db.TableListResult) (Model, tea.Cmd) {
	if msg.Err != nil {
		// TODO: Show error
		return m, nil
	}

	// Sort tables for tree display (parents before children)
	m.Tables.Tables = sortTablesForTree(msg.Tables)
	if len(m.Tables.Tables) > 0 {
		m.Tables.CursorTable = 0
		// SelectedTable stays at -1 until user presses Enter
	}

	return m, nil
}

// sortTablesForTree sorts table names so parent tables appear before their children
// e.g., ["users.phones", "users", "products", "users.addresses"] ->
//
//	["products", "users", "users.addresses", "users.phones"]
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
		m.Schema.ErrorMsg = msg.Err.Error()
		m.Data.LoadingData = false
		return m, nil
	}

	// Clear any previous error
	m.Schema.ErrorMsg = ""
	m.Schema.TableDetails[msg.TableName] = &msg

	// If this is the selected table and we're waiting for data, fetch it now
	if m.Data.LoadingData && !m.SQL.CustomSQL && m.Tables.SelectedTable >= 0 && m.Tables.SelectedTable < len(m.Tables.Tables) {
		tableName := m.Tables.Tables[m.Tables.SelectedTable]
		if tableName == msg.TableName && msg.Schema != nil {
			// Update SQL with ORDER BY
			primaryKeys := ui.ParsePrimaryKeysFromDDL(msg.Schema.DDL)
			m.SQL.CurrentSQL = buildDefaultSQL(tableName, msg.Schema.DDL)
			m.SQL.CursorPos = ui.RuneLen(m.SQL.CurrentSQL)
			// Now fetch data with proper ORDER BY
			return m, db.FetchTableData(m.Connection.NosqlClient, tableName, ui.DefaultFetchSize, primaryKeys)
		}
	}

	return m, nil
}

func handleTableDataResult(m Model, msg db.TableDataResult) (Model, tea.Cmd) {
	if msg.Err != nil {
		m.Data.LoadingData = false
		m.Data.ErrorMsg = msg.Err.Error()
		return m, nil
	}

	// Clear any previous error
	m.Data.ErrorMsg = ""

	// If this is an append operation (additional data fetch), merge with existing data
	if msg.IsAppend {
		if existingData, exists := m.Data.TableData[msg.TableName]; exists && existingData != nil {
			// Append new rows to existing rows
			existingData.Rows = append(existingData.Rows, msg.Rows...)
			// Update pagination info
			existingData.LastPKValues = msg.LastPKValues
			existingData.HasMore = msg.HasMore
			existingData.Offset = msg.Offset
			// Viewport offset stays unchanged - cursor remains at center
			// and new data appears below in the previously empty space
		}
	} else {
		// Store new data
		m.Data.TableData[msg.TableName] = &db.TableDataResult{
			TableName:    msg.TableName,
			Rows:         msg.Rows,
			LastPKValues: msg.LastPKValues,
			HasMore:      msg.HasMore,
			IsCustomSQL:  msg.IsCustomSQL,
			ColumnOrder:  msg.ColumnOrder,
			CurrentSQL:   msg.CurrentSQL,
			Offset:       msg.Offset,
		}
	}

	m.Data.LoadingData = false
	return m, nil
}
