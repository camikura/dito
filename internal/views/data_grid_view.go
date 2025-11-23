package views

import (
	"fmt"

	"github.com/camikura/dito/internal/ui"
	"github.com/oracle/nosql-go-sdk/nosqldb"
)

// DataGridViewModel represents the data needed to render the data grid view.
type DataGridViewModel struct {
	Rows             []map[string]interface{} // Data rows to display
	TableSchema      *nosqldb.TableResult     // Schema information (for column order)
	SelectedRow      int                      // Currently selected row index
	HorizontalOffset int                      // Horizontal scroll offset
	ViewportOffset   int                      // Vertical scroll offset
	Width            int                      // Available width for rendering
	Height           int                      // Available height for rendering
	LoadingData      bool                     // Whether data is being loaded
	Error            error                    // Error if any
	SQL              string                   // SQL query that was executed (for error display)
}

// RenderDataGridView renders the data grid view.
// This is a pure rendering function that takes model data and returns a string.
func RenderDataGridView(m DataGridViewModel) string {
	if m.LoadingData {
		return "Loading data..."
	}

	if m.Rows == nil {
		return "No data available"
	}

	if m.Error != nil {
		return fmt.Sprintf("Error: %v\n\nSQL:\n%s", m.Error, m.SQL)
	}

	if len(m.Rows) == 0 {
		return fmt.Sprintf("No data found\n\nSQL:\n%s", m.SQL)
	}

	// Get column names from actual query results
	// First, collect all columns that exist in the data
	dataColumns := make(map[string]bool)
	if len(m.Rows) > 0 {
		for key := range m.Rows[0] {
			dataColumns[key] = true
		}
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
	} else if len(m.Rows) > 0 {
		// If DDL is not available, get from data (order is undefined)
		for key := range m.Rows[0] {
			columnNames = append(columnNames, key)
		}
	}

	// Render using ui.DataGrid
	grid := ui.DataGrid{
		Rows:             m.Rows,
		Columns:          columnNames,
		SelectedRow:      m.SelectedRow,
		HorizontalOffset: m.HorizontalOffset,
		ViewportOffset:   m.ViewportOffset,
	}
	return grid.Render(m.Width, m.Height)
}
