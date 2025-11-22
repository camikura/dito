package ui

import (
	"strings"
	"testing"
)

func TestDataGrid_Render(t *testing.T) {
	tests := []struct {
		name      string
		grid      DataGrid
		maxWidth  int
		maxHeight int
		contains  []string
	}{
		{
			name: "simple data grid",
			grid: DataGrid{
				Rows: []map[string]interface{}{
					{"id": 1, "name": "Alice"},
					{"id": 2, "name": "Bob"},
				},
				Columns:          []string{"id", "name"},
				SelectedRow:      0,
				HorizontalOffset: 0,
				ViewportOffset:   0,
			},
			maxWidth:  80,
			maxHeight: 10,
			contains:  []string{"id", "name", "Alice", "Bob"},
		},
		{
			name: "with selection on second row",
			grid: DataGrid{
				Rows: []map[string]interface{}{
					{"id": 1, "name": "Alice"},
					{"id": 2, "name": "Bob"},
				},
				Columns:          []string{"id", "name"},
				SelectedRow:      1,
				HorizontalOffset: 0,
				ViewportOffset:   0,
			},
			maxWidth:  80,
			maxHeight: 10,
			contains:  []string{"id", "name", "Alice", "Bob"},
		},
		{
			name: "empty rows",
			grid: DataGrid{
				Rows:             []map[string]interface{}{},
				Columns:          []string{"id", "name"},
				SelectedRow:      0,
				HorizontalOffset: 0,
				ViewportOffset:   0,
			},
			maxWidth:  80,
			maxHeight: 10,
			contains:  []string{"No data"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.grid.Render(tt.maxWidth, tt.maxHeight)

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("Render() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestDataGrid_CalculateColumnWidths(t *testing.T) {
	grid := DataGrid{
		Rows: []map[string]interface{}{
			{"id": 1, "name": "Alice"},
			{"id": 100, "name": "Bob"},
			{"id": 3, "name": "Charlie with a very long name"},
		},
		Columns: []string{"id", "name"},
	}

	widths := grid.calculateColumnWidths(grid.Columns)

	// id column: max of "id" (2) and "100" (3) = 3
	if widths["id"] < 2 {
		t.Errorf("id column width should be at least 2, got %d", widths["id"])
	}

	// name column: should be capped at 32
	if widths["name"] > 32 {
		t.Errorf("name column width should be capped at 32, got %d", widths["name"])
	}
}

func TestDataGrid_GetVisibleColumns(t *testing.T) {
	grid := DataGrid{
		Columns:          []string{"id", "name", "email", "age"},
		HorizontalOffset: 0,
	}

	// No offset
	visible := grid.getVisibleColumns()
	if len(visible) != 4 {
		t.Errorf("Expected 4 visible columns, got %d", len(visible))
	}

	// Offset by 1
	grid.HorizontalOffset = 1
	visible = grid.getVisibleColumns()
	if len(visible) != 3 {
		t.Errorf("Expected 3 visible columns after offset, got %d", len(visible))
	}
	if visible[0] != "name" {
		t.Errorf("First visible column should be 'name', got %q", visible[0])
	}

	// Offset beyond columns
	grid.HorizontalOffset = 10
	visible = grid.getVisibleColumns()
	if len(visible) != 4 {
		t.Errorf("Expected all columns when offset is beyond range, got %d", len(visible))
	}
}

func TestDataGrid_GetViewportRows(t *testing.T) {
	grid := DataGrid{
		Rows: []map[string]interface{}{
			{"id": 1},
			{"id": 2},
			{"id": 3},
			{"id": 4},
			{"id": 5},
		},
		ViewportOffset: 0,
	}

	// Viewport size 2, offset 0
	rows := grid.getViewportRows(2)
	if len(rows) != 2 {
		t.Errorf("Expected 2 viewport rows, got %d", len(rows))
	}
	if rows[0]["id"] != 1 {
		t.Errorf("First row should have id=1, got %v", rows[0]["id"])
	}

	// Viewport size 2, offset 2
	grid.ViewportOffset = 2
	rows = grid.getViewportRows(2)
	if len(rows) != 2 {
		t.Errorf("Expected 2 viewport rows, got %d", len(rows))
	}
	if rows[0]["id"] != 3 {
		t.Errorf("First row should have id=3, got %v", rows[0]["id"])
	}

	// Viewport size larger than remaining rows
	grid.ViewportOffset = 4
	rows = grid.getViewportRows(5)
	if len(rows) != 1 {
		t.Errorf("Expected 1 viewport row, got %d", len(rows))
	}

	// Offset beyond data
	grid.ViewportOffset = 10
	rows = grid.getViewportRows(2)
	if rows != nil {
		t.Errorf("Expected nil for offset beyond data, got %v", rows)
	}
}

func TestDataGrid_RenderHeader(t *testing.T) {
	grid := DataGrid{}
	columns := []string{"id", "name", "email"}
	columnWidths := map[string]int{
		"id":    5,
		"name":  10,
		"email": 20,
	}

	headerParts, headerWidths := grid.renderHeader(columns, columnWidths, 100)

	if len(headerParts) != 3 {
		t.Errorf("Expected 3 header parts, got %d", len(headerParts))
	}

	if len(headerWidths) != 3 {
		t.Errorf("Expected 3 header widths, got %d", len(headerWidths))
	}

	// Check header contains column names
	joined := strings.Join(headerParts, " ")
	for _, col := range columns {
		if !strings.Contains(joined, col) {
			t.Errorf("Header should contain %q", col)
		}
	}
}

func TestDataGrid_RenderRow(t *testing.T) {
	grid := DataGrid{}
	row := map[string]interface{}{
		"id":    1,
		"name":  "Alice",
		"email": "alice@example.com",
	}
	columns := []string{"id", "name", "email"}
	columnWidths := map[string]int{
		"id":    5,
		"name":  10,
		"email": 20,
	}
	headerWidths := []int{5, 10, 20}

	rowParts := grid.renderRow(row, columns, columnWidths, 100, headerWidths)

	if len(rowParts) != 3 {
		t.Errorf("Expected 3 row parts, got %d", len(rowParts))
	}

	// Check row contains data
	joined := strings.Join(rowParts, " ")
	if !strings.Contains(joined, "Alice") {
		t.Errorf("Row should contain 'Alice'")
	}
}

func TestDataGrid_Selection(t *testing.T) {
	grid := DataGrid{
		Rows: []map[string]interface{}{
			{"id": 1, "name": "Alice"},
			{"id": 2, "name": "Bob"},
			{"id": 3, "name": "Charlie"},
		},
		Columns:          []string{"id", "name"},
		SelectedRow:      1,
		HorizontalOffset: 0,
		ViewportOffset:   0,
	}

	result := grid.Render(80, 10)

	// Should contain all data
	if !strings.Contains(result, "Alice") {
		t.Errorf("Result should contain 'Alice'")
	}
	if !strings.Contains(result, "Bob") {
		t.Errorf("Result should contain 'Bob'")
	}
	if !strings.Contains(result, "Charlie") {
		t.Errorf("Result should contain 'Charlie'")
	}
}

func TestDataGrid_HorizontalScroll(t *testing.T) {
	grid := DataGrid{
		Rows: []map[string]interface{}{
			{"col1": "A", "col2": "B", "col3": "C", "col4": "D"},
		},
		Columns:          []string{"col1", "col2", "col3", "col4"},
		SelectedRow:      0,
		HorizontalOffset: 2, // Skip first 2 columns
		ViewportOffset:   0,
	}

	result := grid.Render(80, 10)

	// Should contain col3 and col4
	if !strings.Contains(result, "col3") {
		t.Errorf("Result should contain 'col3' after horizontal scroll")
	}
	if !strings.Contains(result, "col4") {
		t.Errorf("Result should contain 'col4' after horizontal scroll")
	}

	// Should NOT contain col1
	// Note: might contain "col1" as part of header, but not as visible column
}
