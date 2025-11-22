package db

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/oracle/nosql-go-sdk/nosqldb"
	"github.com/oracle/nosql-go-sdk/nosqldb/nosqlerr"
)

// ConnectionResult represents the result of a connection attempt.
type ConnectionResult struct {
	Err      error
	Version  string
	Client   *nosqldb.Client
	Endpoint string
	IsTest   bool // true for test connections (no screen transition)
}

// TableListResult represents the result of fetching table list.
type TableListResult struct {
	Tables []string
	Err    error
}

// TableDetailsResult represents the result of fetching table details.
type TableDetailsResult struct {
	TableName string
	Schema    *nosqldb.TableResult
	Indexes   []nosqldb.IndexInfo
	Err       error
}

// TableDataResult represents the result of fetching table data.
type TableDataResult struct {
	TableName    string
	Rows         []map[string]interface{}
	LastPKValues map[string]interface{} // Last row's PRIMARY KEY values (used as cursor)
	HasMore      bool                   // Whether more data is available
	Err          error
	IsAppend     bool   // Whether to append to existing data
	SQL          string // Debug: executed SQL
	DisplaySQL   string // Display: SQL without LIMIT clause
}

// Connect attempts to connect to NoSQL database.
// Returns a tea.Cmd that produces a ConnectionResult message.
func Connect(endpoint, port string, isTest bool) tea.Cmd {
	return func() tea.Msg {
		// Connection configuration
		endpointURL := fmt.Sprintf("http://%s:%s", endpoint, port)
		cfg := nosqldb.Config{
			Mode:     "onprem",
			Endpoint: endpointURL,
		}

		// Create client
		client, err := nosqldb.NewClient(cfg)
		if err != nil {
			return ConnectionResult{Err: err, IsTest: isTest}
		}

		// Simple test (fetch table list)
		req := &nosqldb.ListTablesRequest{}
		_, err = client.ListTables(req)
		if err != nil {
			client.Close()
			// Get error details
			if nosqlErr, ok := err.(*nosqlerr.Error); ok {
				return ConnectionResult{Err: fmt.Errorf("NoSQL Error: %s", nosqlErr.Error()), IsTest: isTest}
			}
			return ConnectionResult{Err: err, IsTest: isTest}
		}

		// Close client for test connection
		if isTest {
			client.Close()
			return ConnectionResult{
				Version: "Connected",
				Err:     nil,
				IsTest:  true,
			}
		}

		// Connection successful - return client (don't close)
		return ConnectionResult{
			Version:  "Connected",
			Err:      nil,
			Client:   client,
			Endpoint: fmt.Sprintf("%s:%s", endpoint, port),
			IsTest:   false,
		}
	}
}

// FetchTables fetches the list of tables from NoSQL database.
// Returns a tea.Cmd that produces a TableListResult message.
func FetchTables(client *nosqldb.Client) tea.Cmd {
	return func() tea.Msg {
		req := &nosqldb.ListTablesRequest{}
		result, err := client.ListTables(req)
		if err != nil {
			return TableListResult{Err: err}
		}

		// Filter out system tables (SYS$*)
		var userTables []string
		for _, table := range result.Tables {
			if !strings.HasPrefix(table, "SYS$") {
				userTables = append(userTables, table)
			}
		}

		return TableListResult{Tables: userTables, Err: nil}
	}
}

// FetchTableDetails fetches table schema and index information.
// Returns a tea.Cmd that produces a TableDetailsResult message.
func FetchTableDetails(client *nosqldb.Client, tableName string) tea.Cmd {
	return func() tea.Msg {
		// Get table information
		tableReq := &nosqldb.GetTableRequest{
			TableName: tableName,
		}
		tableResult, err := client.GetTable(tableReq)
		if err != nil {
			return TableDetailsResult{TableName: tableName, Err: err}
		}

		// Get index information
		indexReq := &nosqldb.GetIndexesRequest{
			TableName: tableName,
		}
		indexResult, err := client.GetIndexes(indexReq)
		if err != nil {
			// Ignore index fetch errors, return schema information only
			return TableDetailsResult{TableName: tableName, Schema: tableResult, Indexes: nil, Err: nil}
		}

		return TableDetailsResult{TableName: tableName, Schema: tableResult, Indexes: indexResult.Indexes, Err: nil}
	}
}

// FetchTableData fetches table data (initial fetch, sorted by PRIMARY KEY).
// Returns a tea.Cmd that produces a TableDataResult message.
func FetchTableData(client *nosqldb.Client, tableName string, limit int, primaryKeys []string) tea.Cmd {
	return fetchTableDataWithCursor(client, tableName, limit, primaryKeys, nil, false)
}

// FetchMoreTableData fetches additional table data (using PRIMARY KEY cursor).
// Returns a tea.Cmd that produces a TableDataResult message.
func FetchMoreTableData(client *nosqldb.Client, tableName string, limit int, primaryKeys []string, lastPKValues map[string]interface{}) tea.Cmd {
	return fetchTableDataWithCursor(client, tableName, limit, primaryKeys, lastPKValues, true)
}

// fetchTableDataWithCursor is an internal function to fetch table data with PRIMARY KEY cursor support.
func fetchTableDataWithCursor(client *nosqldb.Client, tableName string, limit int, primaryKeys []string, lastPKValues map[string]interface{}, isAppend bool) tea.Cmd {
	return func() tea.Msg {
		// Explicitly sort by PRIMARY KEY order
		var orderByClause string
		if len(primaryKeys) > 0 {
			orderByClause = " ORDER BY " + strings.Join(primaryKeys, ", ")
		}

		// Build WHERE clause (if PRIMARY KEY cursor exists)
		var whereClause string
		if lastPKValues != nil && len(lastPKValues) > 0 && len(primaryKeys) > 0 {
			// Build condition for composite PRIMARY KEY
			// Example: WHERE pk1 > ? OR (pk1 = ? AND pk2 > ?) OR (pk1 = ? AND pk2 = ? AND pk3 > ?)
			var conditions []string
			for i := 0; i < len(primaryKeys); i++ {
				var cond string
				if i == 0 {
					// First key: pk1 > ?
					val := lastPKValues[primaryKeys[i]]
					cond = fmt.Sprintf("%s > %s", primaryKeys[i], formatValue(val))
				} else {
					// Following keys: (pk1 = ? AND pk2 = ? AND ... AND pkN > ?)
					var parts []string
					for j := 0; j < i; j++ {
						val := lastPKValues[primaryKeys[j]]
						parts = append(parts, fmt.Sprintf("%s = %s", primaryKeys[j], formatValue(val)))
					}
					val := lastPKValues[primaryKeys[i]]
					parts = append(parts, fmt.Sprintf("%s > %s", primaryKeys[i], formatValue(val)))
					cond = "(" + strings.Join(parts, " AND ") + ")"
				}
				conditions = append(conditions, cond)
			}
			whereClause = " WHERE " + strings.Join(conditions, " OR ")
		}

		statement := fmt.Sprintf("SELECT * FROM %s%s%s LIMIT %d", tableName, whereClause, orderByClause, limit)

		// Display SQL (without LIMIT clause)
		displayStatement := fmt.Sprintf("SELECT * FROM %s%s%s", tableName, whereClause, orderByClause)

		prepReq := &nosqldb.PrepareRequest{
			Statement: statement,
		}
		prepResult, err := client.Prepare(prepReq)
		if err != nil {
			return TableDataResult{TableName: tableName, Err: err, IsAppend: isAppend, SQL: statement, DisplaySQL: displayStatement}
		}

		queryReq := &nosqldb.QueryRequest{
			PreparedStatement: &prepResult.PreparedStatement,
		}

		// Fetch all results (using SDK's internal pagination)
		var rows []map[string]interface{}
		for {
			queryResult, err := client.Query(queryReq)
			if err != nil {
				return TableDataResult{TableName: tableName, Err: err, IsAppend: isAppend, SQL: statement, DisplaySQL: displayStatement}
			}

			// Get results
			results, err := queryResult.GetResults()
			if err != nil {
				return TableDataResult{TableName: tableName, Err: err, IsAppend: isAppend, SQL: statement, DisplaySQL: displayStatement}
			}

			for _, result := range results {
				rows = append(rows, result.Map())
			}

			// Exit if no continuation token
			if queryReq.IsDone() {
				break
			}
		}

		// Save last row's PRIMARY KEY values
		var newLastPKValues map[string]interface{}
		if len(rows) > 0 && len(primaryKeys) > 0 {
			lastRow := rows[len(rows)-1]
			newLastPKValues = make(map[string]interface{})
			for _, pk := range primaryKeys {
				if val, exists := lastRow[pk]; exists {
					newLastPKValues[pk] = val
				}
			}
		}

		// Check if more pages exist
		// If fetched rows == limit, more data may be available
		hasMore := len(rows) == limit

		return TableDataResult{
			TableName:    tableName,
			Rows:         rows,
			LastPKValues: newLastPKValues,
			HasMore:      hasMore,
			Err:          nil,
			IsAppend:     isAppend,
			SQL:          statement,
			DisplaySQL:   displayStatement,
		}
	}
}

// formatValue formats a value for SQL statement.
func formatValue(val interface{}) string {
	switch v := val.(type) {
	case string:
		// Escape single quotes by doubling them
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
	case int, int32, int64, float32, float64:
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("'%v'", v)
	}
}
