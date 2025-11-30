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
			return m, db.FetchTableData(m.NosqlClient, tableName, ui.DefaultFetchSize, primaryKeys)
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
			existingData.Offset = msg.Offset
		}
	} else {
		// Store new data
		m.TableData[msg.TableName] = &db.TableDataResult{
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

	m.LoadingData = false
	return m, nil
}
