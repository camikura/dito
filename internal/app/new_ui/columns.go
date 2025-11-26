package new_ui

import (
	"github.com/camikura/dito/internal/views"
)

// getColumnsInSchemaOrder returns column names in schema definition order
func getColumnsInSchemaOrder(m Model, tableName string, rows []map[string]interface{}) []string {
	var columns []string
	columnSet := make(map[string]bool)

	// Try to get columns from schema DDL first
	if details, exists := m.TableDetails[tableName]; exists && details != nil && details.Schema != nil {
		if details.Schema.DDL != "" {
			primaryKeys := views.ParsePrimaryKeysFromDDL(details.Schema.DDL)
			cols := views.ParseColumnsFromDDL(details.Schema.DDL, primaryKeys)

			// Extract column names in schema order
			for _, col := range cols {
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
		sortColumns(extraColumns)
		columns = append(columns, extraColumns...)
	}

	// Fallback: if no columns from schema, extract from first row
	if len(columns) == 0 && len(rows) > 0 {
		for col := range rows[0] {
			columns = append(columns, col)
		}
		sortColumns(columns)
	}

	return columns
}

// sortColumns sorts column names (simple alphabetical for now)
func sortColumns(columns []string) {
	// Simple bubble sort
	n := len(columns)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if columns[j] > columns[j+1] {
				columns[j], columns[j+1] = columns[j+1], columns[j]
			}
		}
	}
}
