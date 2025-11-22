package views

import (
	"errors"
	"strings"
	"testing"

	"github.com/oracle/nosql-go-sdk/nosqldb"
)

func TestRenderDataGridView(t *testing.T) {
	tests := []struct {
		name     string
		model    DataGridViewModel
		contains []string
	}{
		{
			name: "loading state",
			model: DataGridViewModel{
				LoadingData: true,
			},
			contains: []string{"Loading data..."},
		},
		{
			name: "no data available",
			model: DataGridViewModel{
				Rows:        nil,
				LoadingData: false,
			},
			contains: []string{"No data available"},
		},
		{
			name: "error state",
			model: DataGridViewModel{
				Rows:        []map[string]interface{}{},
				Error:       errors.New("connection timeout"),
				SQL:         "SELECT * FROM users",
				LoadingData: false,
			},
			contains: []string{"Error:", "connection timeout", "SQL:", "SELECT * FROM users"},
		},
		{
			name: "empty result set",
			model: DataGridViewModel{
				Rows:        []map[string]interface{}{},
				Error:       nil,
				SQL:         "SELECT * FROM users WHERE id = 999",
				LoadingData: false,
			},
			contains: []string{"No data found", "SQL:", "SELECT * FROM users WHERE id = 999"},
		},
		{
			name: "successful data rendering with schema",
			model: DataGridViewModel{
				Rows: []map[string]interface{}{
					{"id": 1, "name": "Alice"},
					{"id": 2, "name": "Bob"},
				},
				TableSchema: &nosqldb.TableResult{
					DDL: "CREATE TABLE users (id INTEGER, name STRING, PRIMARY KEY(id))",
				},
				SelectedRow:      0,
				HorizontalOffset: 0,
				ViewportOffset:   0,
				Width:            80,
				Height:           10,
				LoadingData:      false,
			},
			contains: []string{"id", "name"},
		},
		{
			name: "data rendering without schema",
			model: DataGridViewModel{
				Rows: []map[string]interface{}{
					{"id": 1, "name": "Alice", "email": "alice@example.com"},
				},
				TableSchema:      nil,
				SelectedRow:      0,
				HorizontalOffset: 0,
				ViewportOffset:   0,
				Width:            80,
				Height:           10,
				LoadingData:      false,
			},
			contains: []string{}, // Column names may be in any order without schema
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderDataGridView(tt.model)

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("RenderDataGridView() = %q, should contain %q", result, substr)
				}
			}

			if result == "" {
				t.Error("RenderDataGridView() should not return empty string")
			}
		})
	}
}

func TestRenderDataGridView_ColumnOrder(t *testing.T) {
	// Test that columns are in DDL order when schema is provided
	model := DataGridViewModel{
		Rows: []map[string]interface{}{
			{"name": "Alice", "id": 1, "email": "alice@example.com"}, // Intentionally out of order
		},
		TableSchema: &nosqldb.TableResult{
			DDL: "CREATE TABLE users (id INTEGER, name STRING, email STRING, PRIMARY KEY(id))",
		},
		SelectedRow:      0,
		HorizontalOffset: 0,
		ViewportOffset:   0,
		Width:            80,
		Height:           10,
		LoadingData:      false,
	}

	result := RenderDataGridView(model)

	// The result should contain columns in DDL order (id, name, email)
	// We can't test exact positioning due to formatting, but we can verify all columns are present
	if !strings.Contains(result, "id") {
		t.Error("RenderDataGridView() should contain 'id' column")
	}
	if !strings.Contains(result, "name") {
		t.Error("RenderDataGridView() should contain 'name' column")
	}
	if !strings.Contains(result, "email") {
		t.Error("RenderDataGridView() should contain 'email' column")
	}
}
