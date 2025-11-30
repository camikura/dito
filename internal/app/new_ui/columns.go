package new_ui

import (
	"github.com/camikura/dito/internal/ui"
)

// getColumnsInSchemaOrder returns column names in schema definition order.
// For custom SQL queries, it uses the column order from SELECT clause if available.
func getColumnsInSchemaOrder(m Model, tableName string, rows []map[string]interface{}) []string {
	// For custom SQL with column order from SELECT clause
	if m.CustomSQL && len(m.ColumnOrder) > 0 {
		return m.ColumnOrder
	}

	// For custom SQL without parsed column order, fall back to data columns
	if m.CustomSQL {
		return ui.GetColumnsFromData(rows)
	}

	// Get DDL from table details
	var ddl string
	if details, exists := m.TableDetails[tableName]; exists && details != nil && details.Schema != nil {
		ddl = details.Schema.DDL
	}

	return ui.GetColumnsInSchemaOrder(ddl, rows)
}
