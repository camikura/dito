package db

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/oracle/nosql-go-sdk/nosqldb"
	"github.com/oracle/nosql-go-sdk/nosqldb/nosqlerr"
	"github.com/oracle/nosql-go-sdk/nosqldb/types"
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
	IsAppend     bool     // Whether to append to existing data
	SQL          string   // Debug: executed SQL
	DisplaySQL   string   // Display: SQL without LIMIT clause
	IsCustomSQL  bool     // Whether this is a custom SQL query (not auto-generated)
	ColumnOrder  []string // Column order from SELECT clause (for custom SQL)
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

// parseSelectColumns extracts column names/aliases from SELECT clause in order.
// Handles: SELECT col1, col2 as alias, t.col3 FROM ...
func parseSelectColumns(sql string) []string {
	upperSQL := strings.ToUpper(sql)

	// Find SELECT and FROM positions
	selectIdx := strings.Index(upperSQL, "SELECT")
	if selectIdx == -1 {
		return nil
	}
	fromIdx := strings.Index(upperSQL, "FROM")
	if fromIdx == -1 || fromIdx <= selectIdx+6 {
		return nil
	}

	// Extract the part between SELECT and FROM
	selectPart := strings.TrimSpace(sql[selectIdx+6 : fromIdx])
	if selectPart == "*" || strings.HasPrefix(selectPart, "* ") {
		return nil // SELECT * doesn't specify column order
	}

	// Split by comma, handling potential nested parentheses
	var columns []string
	depth := 0
	current := ""
	for _, ch := range selectPart {
		if ch == '(' {
			depth++
			current += string(ch)
		} else if ch == ')' {
			depth--
			current += string(ch)
		} else if ch == ',' && depth == 0 {
			columns = append(columns, strings.TrimSpace(current))
			current = ""
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		columns = append(columns, strings.TrimSpace(current))
	}

	// Extract column name or alias from each column expression
	var result []string
	for _, col := range columns {
		col = strings.TrimSpace(col)
		if col == "" {
			continue
		}

		// Check for AS alias (case insensitive)
		upperCol := strings.ToUpper(col)
		if asIdx := strings.LastIndex(upperCol, " AS "); asIdx != -1 {
			// Use the alias
			alias := strings.TrimSpace(col[asIdx+4:])
			result = append(result, alias)
		} else {
			// No alias, use the column name (handle table.column format)
			parts := strings.Split(col, ".")
			colName := strings.TrimSpace(parts[len(parts)-1])
			result = append(result, colName)
		}
	}

	return result
}

// ExecuteCustomSQL executes custom SQL query and returns results.
// Returns a tea.Cmd that produces a TableDataResult message.
func ExecuteCustomSQL(client *nosqldb.Client, tableName string, sql string, limit int) tea.Cmd {
	return func() tea.Msg {
		// Parse column order from SELECT clause
		columnOrder := parseSelectColumns(sql)

		// Add LIMIT clause to SQL if not present
		statement := sql
		displayStatement := sql
		if !strings.Contains(strings.ToUpper(sql), "LIMIT") {
			statement = fmt.Sprintf("%s LIMIT %d", sql, limit)
		}

		prepReq := &nosqldb.PrepareRequest{
			Statement: statement,
		}
		prepResult, err := client.Prepare(prepReq)
		if err != nil {
			return TableDataResult{TableName: tableName, Err: err, IsAppend: false, SQL: statement, DisplaySQL: displayStatement, IsCustomSQL: true}
		}

		queryReq := &nosqldb.QueryRequest{
			PreparedStatement: &prepResult.PreparedStatement,
		}

		// Fetch all results (using SDK's internal pagination)
		var rows []map[string]interface{}
		for {
			queryResult, err := client.Query(queryReq)
			if err != nil {
				return TableDataResult{TableName: tableName, Err: err, IsAppend: false, SQL: statement, DisplaySQL: displayStatement, IsCustomSQL: true}
			}

			// Get results
			results, err := queryResult.GetResults()
			if err != nil {
				return TableDataResult{TableName: tableName, Err: err, IsAppend: false, SQL: statement, DisplaySQL: displayStatement, IsCustomSQL: true}
			}

			for _, result := range results {
				row := result.Map()
				// Convert SDK-specific types (e.g., *types.MapValue) to native Go types
				convertedRow := convertRowValues(row)
				rows = append(rows, convertedRow)
			}

			// Exit if no continuation token
			if queryReq.IsDone() {
				break
			}
		}

		// Check if more pages exist
		hasMore := len(rows) == limit

		return TableDataResult{
			TableName:    tableName,
			Rows:         rows,
			LastPKValues: nil, // Custom queries don't support cursor-based pagination
			HasMore:      hasMore,
			Err:          nil,
			IsAppend:     false,
			SQL:          statement,
			DisplaySQL:   displayStatement,
			IsCustomSQL:  true, // This is a custom SQL query
			ColumnOrder:  columnOrder,
		}
	}
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
			return TableDataResult{TableName: tableName, Err: err, IsAppend: isAppend, SQL: statement, DisplaySQL: displayStatement, IsCustomSQL: false}
		}

		queryReq := &nosqldb.QueryRequest{
			PreparedStatement: &prepResult.PreparedStatement,
		}

		// Fetch all results (using SDK's internal pagination)
		var rows []map[string]interface{}
		for {
			queryResult, err := client.Query(queryReq)
			if err != nil {
				return TableDataResult{TableName: tableName, Err: err, IsAppend: isAppend, SQL: statement, DisplaySQL: displayStatement, IsCustomSQL: false}
			}

			// Get results
			results, err := queryResult.GetResults()
			if err != nil {
				return TableDataResult{TableName: tableName, Err: err, IsAppend: isAppend, SQL: statement, DisplaySQL: displayStatement, IsCustomSQL: false}
			}

			for _, result := range results {
				row := result.Map()
				// Convert SDK-specific types (e.g., *types.MapValue) to native Go types
				convertedRow := convertRowValues(row)
				rows = append(rows, convertedRow)
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
			IsCustomSQL:  false, // This is an auto-generated SQL query
		}
	}
}

// convertRowValues converts SDK-specific types in a row to native Go types.
// This is needed for JSON fields which are returned as *types.MapValue by the SDK.
func convertRowValues(row map[string]interface{}) map[string]interface{} {
	converted := make(map[string]interface{})
	for key, val := range row {
		converted[key] = convertValue(val)
	}
	return converted
}

// convertValue recursively converts SDK-specific types to native Go types.
func convertValue(val interface{}) interface{} {
	if val == nil {
		return nil
	}

	// Handle *types.MapValue (JSON object)
	if mapVal, ok := val.(*types.MapValue); ok {
		return convertMapValue(mapVal)
	}

	// Handle types.MapValue (non-pointer)
	if mapVal, ok := val.(types.MapValue); ok {
		return convertMapValue(&mapVal)
	}

	// Use reflection to handle other SDK types
	v := reflect.ValueOf(val)

	// If it's a pointer, dereference it
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	// Handle map types
	if v.Kind() == reflect.Map {
		result := make(map[string]interface{})
		for _, key := range v.MapKeys() {
			keyStr := fmt.Sprintf("%v", key.Interface())
			result[keyStr] = convertValue(v.MapIndex(key).Interface())
		}
		return result
	}

	// Handle slice/array types
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		length := v.Len()
		result := make([]interface{}, length)
		for i := 0; i < length; i++ {
			result[i] = convertValue(v.Index(i).Interface())
		}
		return result
	}

	// If we reach here, it might be an SDK type we don't recognize
	// Try to convert it using JSON marshal/unmarshal as a last resort
	jsonBytes, err := json.Marshal(val)
	if err == nil {
		var result interface{}
		if err := json.Unmarshal(jsonBytes, &result); err == nil {
			return result
		}
	}

	// Return primitive types as-is
	return val
}

// convertMapValue converts *types.MapValue to map[string]interface{}
func convertMapValue(mapVal *types.MapValue) map[string]interface{} {
	if mapVal == nil {
		return nil
	}

	// Try to marshal to JSON and back to get a clean map[string]interface{}
	jsonBytes, err := json.Marshal(mapVal)
	if err != nil {
		// Fallback: return empty map for conversion errors
		return make(map[string]interface{})
	}

	var result map[string]interface{}
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		// Fallback: return empty map for unmarshal errors
		return make(map[string]interface{})
	}

	// Ensure we return an initialized map, not nil
	if result == nil {
		return make(map[string]interface{})
	}

	return result
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
