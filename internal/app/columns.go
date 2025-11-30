package app

import (
	"github.com/camikura/dito/internal/ui"
)

// getColumnsInSchemaOrder returns column names in schema definition order.
// For custom SQL with explicit column list, it uses the parsed column order.
// For SELECT * or normal queries, it uses schema definition order.
// For child tables, ancestor primary key columns are placed first.
func getColumnsInSchemaOrder(m Model, tableName string, rows []map[string]interface{}) []string {
	// For custom SQL with explicit column order from SELECT clause (not SELECT *)
	if m.CustomSQL && len(m.ColumnOrder) > 0 {
		return m.ColumnOrder
	}

	// Use schema order for SELECT * and normal queries
	var ddl string
	if details, exists := m.TableDetails[tableName]; exists && details != nil && details.Schema != nil {
		ddl = details.Schema.DDL
	}

	// Get ancestor DDLs for child tables (root to parent order)
	var ancestorDDLs []string
	ancestors := ui.GetAncestorTableNames(tableName)
	for _, ancestor := range ancestors {
		if details, exists := m.TableDetails[ancestor]; exists && details != nil && details.Schema != nil {
			ancestorDDLs = append(ancestorDDLs, details.Schema.DDL)
		} else {
			ancestorDDLs = append(ancestorDDLs, "")
		}
	}

	return ui.GetColumnsInSchemaOrderWithAncestors(ddl, ancestorDDLs, rows)
}
