package views

import (
	"fmt"

	"github.com/camikura/dito/internal/ui"
	"github.com/oracle/nosql-go-sdk/nosqldb"
)

// RecordViewModel represents the data needed to render the record view.
type RecordViewModel struct {
	Rows        []map[string]interface{} // All data rows
	TableSchema *nosqldb.TableResult     // Schema information (for column order)
	SelectedRow int                      // Currently selected row index
	LoadingData bool                     // Whether data is being loaded
	Error       error                    // Error if any
}

// RenderRecordView renders the record view (vertical table).
// This is a pure rendering function that takes model data and returns a string.
func RenderRecordView(m RecordViewModel) string {
	if m.LoadingData {
		return "Loading data..."
	}

	if m.Rows == nil {
		return "No data available"
	}

	if m.Error != nil {
		return fmt.Sprintf("Error: %v", m.Error)
	}

	if len(m.Rows) == 0 {
		return "No data found"
	}

	if m.SelectedRow < 0 || m.SelectedRow >= len(m.Rows) {
		return "Invalid row selection"
	}

	// Get selected row
	selectedRow := m.Rows[m.SelectedRow]

	// Get column names from actual query results
	// First, collect all columns that exist in the data
	dataColumns := make(map[string]bool)
	for key := range selectedRow {
		dataColumns[key] = true
	}

	// Use schema order if available, but only include columns that exist in the data
	var columnNames []string
	if m.TableSchema != nil && m.TableSchema.DDL != "" {
		primaryKeys := ParsePrimaryKeysFromDDL(m.TableSchema.DDL)
		columns := ParseColumnsFromDDL(m.TableSchema.DDL, primaryKeys)
		for _, col := range columns {
			// Only include columns that are actually in the query results
			if dataColumns[col.Name] {
				columnNames = append(columnNames, col.Name)
			}
		}
	} else if len(selectedRow) > 0 {
		// If DDL is not available, get from data (order is undefined)
		for key := range selectedRow {
			columnNames = append(columnNames, key)
		}
	}

	// Render using ui.VerticalTable
	verticalTable := ui.VerticalTable{
		Data: selectedRow,
		Keys: columnNames,
	}
	return verticalTable.Render()
}
