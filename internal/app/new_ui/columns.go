package new_ui

import (
	"github.com/camikura/dito/internal/ui"
)

// getColumnsInSchemaOrder returns column names in schema definition order
func getColumnsInSchemaOrder(m Model, tableName string, rows []map[string]interface{}) []string {
	// Get DDL from table details
	var ddl string
	if details, exists := m.TableDetails[tableName]; exists && details != nil && details.Schema != nil {
		ddl = details.Schema.DDL
	}

	return ui.GetColumnsInSchemaOrder(ddl, rows)
}
