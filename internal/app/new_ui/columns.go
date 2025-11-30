package new_ui

import (
	"github.com/camikura/dito/internal/ui"
)

// getColumnsInSchemaOrder returns column names in schema definition order.
// For custom SQL queries, it only returns columns that exist in the actual data.
func getColumnsInSchemaOrder(m Model, tableName string, rows []map[string]interface{}) []string {
	// For custom SQL, only show columns that are actually in the data
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
