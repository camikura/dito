package ui

import (
	"sort"
	"strings"
)

// ColumnInfo represents a column with its name and type.
type ColumnInfo struct {
	Name         string
	Type         string
	IsPrimaryKey bool
	IsInherited  bool // True if this column is inherited from a parent table
}

// GetParentTableName returns the parent table name for a child table.
// For "users.addresses" returns "users", for "users.addresses.phones" returns "users.addresses".
// Returns empty string if the table has no parent.
func GetParentTableName(tableName string) string {
	lastDot := strings.LastIndex(tableName, ".")
	if lastDot == -1 {
		return ""
	}
	return tableName[:lastDot]
}

// GetAncestorTableNames returns all ancestor table names from root to immediate parent.
// For "users.addresses.phones" returns ["users", "users.addresses"].
func GetAncestorTableNames(tableName string) []string {
	var ancestors []string
	current := tableName
	for {
		parent := GetParentTableName(current)
		if parent == "" {
			break
		}
		// Prepend to get root-to-leaf order
		ancestors = append([]string{parent}, ancestors...)
		current = parent
	}
	return ancestors
}

// ParsePrimaryKeysFromDDL extracts primary key column names from DDL string.
func ParsePrimaryKeysFromDDL(ddl string) []string {
	var primaryKeys []string

	// Find PRIMARY KEY(col1, col2, ...) part
	upperDDL := strings.ToUpper(ddl)
	pkIndex := strings.Index(upperDDL, "PRIMARY KEY")
	if pkIndex == -1 {
		return primaryKeys
	}

	// Get content inside parentheses after PRIMARY KEY
	pkPart := ddl[pkIndex:]
	start := strings.Index(pkPart, "(")
	end := strings.LastIndex(pkPart, ")") // Get last parenthesis
	if start == -1 || end == -1 || start >= end {
		return primaryKeys
	}

	// Extract column name list
	keysPart := pkPart[start+1 : end]

	// Handle SHARD() syntax
	// Support format like PRIMARY KEY(SHARD(id), name)
	keysPart = strings.ReplaceAll(keysPart, "SHARD(", "")
	keysPart = strings.ReplaceAll(keysPart, "shard(", "")
	keysPart = strings.ReplaceAll(keysPart, ")", "")

	keys := strings.Split(keysPart, ",")
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key != "" {
			primaryKeys = append(primaryKeys, key)
		}
	}

	return primaryKeys
}

// ParseColumnsFromDDL extracts column information from DDL string.
func ParseColumnsFromDDL(ddl string, primaryKeys []string) []ColumnInfo {
	var columns []ColumnInfo

	// Create PRIMARY KEY map for fast lookup
	pkMap := make(map[string]bool)
	for _, pk := range primaryKeys {
		pkMap[pk] = true
	}

	// Extract column definition part from CREATE TABLE ... ( ... )
	start := strings.Index(ddl, "(")
	end := strings.LastIndex(ddl, ")")
	if start == -1 || end == -1 || start >= end {
		return columns
	}

	columnDefs := ddl[start+1 : end]

	// Exclude PRIMARY KEY definition
	lines := strings.Split(columnDefs, ",")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip PRIMARY KEY line
		if strings.HasPrefix(strings.ToUpper(line), "PRIMARY KEY") {
			continue
		}

		// Extract column name and type
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			name := parts[0]
			typ := parts[1]

			// Check if it's a PRIMARY KEY
			isPK := pkMap[name]

			columns = append(columns, ColumnInfo{
				Name:         name,
				Type:         typ,
				IsPrimaryKey: isPK,
			})
		}
	}

	return columns
}

// GetColumnsInSchemaOrder returns column names in schema definition order.
// It first tries to get columns from DDL, then adds any extra columns from actual data.
func GetColumnsInSchemaOrder(ddl string, rows []map[string]interface{}) []string {
	return GetColumnsInSchemaOrderWithAncestors(ddl, nil, rows)
}

// GetColumnsInSchemaOrderWithAncestors returns column names in schema definition order,
// including inherited primary key columns from ancestor tables.
// ancestorDDLs should be in order from root to immediate parent.
func GetColumnsInSchemaOrderWithAncestors(ddl string, ancestorDDLs []string, rows []map[string]interface{}) []string {
	var columns []string
	columnSet := make(map[string]bool)

	// First, add primary key columns from ancestors (root to parent order)
	for _, ancestorDDL := range ancestorDDLs {
		if ancestorDDL != "" {
			ancestorPKs := ParsePrimaryKeysFromDDL(ancestorDDL)
			ancestorCols := ParseColumnsFromDDL(ancestorDDL, ancestorPKs)
			// Only add primary key columns from ancestors
			for _, col := range ancestorCols {
				if col.IsPrimaryKey && !columnSet[col.Name] {
					columns = append(columns, col.Name)
					columnSet[col.Name] = true
				}
			}
		}
	}

	// Then add this table's own columns from DDL
	if ddl != "" {
		primaryKeys := ParsePrimaryKeysFromDDL(ddl)
		cols := ParseColumnsFromDDL(ddl, primaryKeys)

		// Extract column names in schema order
		for _, col := range cols {
			if !columnSet[col.Name] {
				columns = append(columns, col.Name)
				columnSet[col.Name] = true
			}
		}
	}

	// Add any columns from actual data that are not in schema
	// (e.g., JSON columns, metadata, etc.)
	if len(rows) > 0 {
		var extraColumns []string
		for col := range rows[0] {
			if !columnSet[col] {
				extraColumns = append(extraColumns, col)
			}
		}
		// Sort extra columns alphabetically and append
		sort.Strings(extraColumns)
		columns = append(columns, extraColumns...)
	}

	// Fallback: if no columns from schema, extract from first row
	if len(columns) == 0 && len(rows) > 0 {
		for col := range rows[0] {
			columns = append(columns, col)
		}
		sort.Strings(columns)
	}

	return columns
}

// GetColumnsFromData extracts column names from actual data rows.
// Used for custom SQL queries where we only want columns that exist in the result.
func GetColumnsFromData(rows []map[string]interface{}) []string {
	if len(rows) == 0 {
		return []string{}
	}

	var columns []string
	for col := range rows[0] {
		columns = append(columns, col)
	}
	sort.Strings(columns)
	return columns
}

// GetColumnTypes extracts column types from DDL.
// Returns a map of column name to type (without Primary Key suffix).
func GetColumnTypes(ddl string) map[string]string {
	types := make(map[string]string)

	if ddl == "" {
		return types
	}

	primaryKeys := ParsePrimaryKeysFromDDL(ddl)
	cols := ParseColumnsFromDDL(ddl, primaryKeys)

	for _, col := range cols {
		types[col.Name] = col.Type
	}

	return types
}
