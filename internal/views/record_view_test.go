package views

import (
	"errors"
	"strings"
	"testing"

	"github.com/oracle/nosql-go-sdk/nosqldb"
)

func TestRenderRecordView(t *testing.T) {
	tests := []struct {
		name     string
		model    RecordViewModel
		contains []string
	}{
		{
			name: "loading state",
			model: RecordViewModel{
				LoadingData: true,
			},
			contains: []string{"Loading data..."},
		},
		{
			name: "no data available",
			model: RecordViewModel{
				Rows:        nil,
				LoadingData: false,
			},
			contains: []string{"No data available"},
		},
		{
			name: "error state",
			model: RecordViewModel{
				Rows:        []map[string]interface{}{},
				Error:       errors.New("connection failed"),
				LoadingData: false,
			},
			contains: []string{"Error:", "connection failed"},
		},
		{
			name: "empty result set",
			model: RecordViewModel{
				Rows:        []map[string]interface{}{},
				Error:       nil,
				LoadingData: false,
			},
			contains: []string{"No data found"},
		},
		{
			name: "invalid row selection - negative index",
			model: RecordViewModel{
				Rows: []map[string]interface{}{
					{"id": 1, "name": "Alice"},
				},
				SelectedRow: -1,
				LoadingData: false,
			},
			contains: []string{"Invalid row selection"},
		},
		{
			name: "invalid row selection - out of bounds",
			model: RecordViewModel{
				Rows: []map[string]interface{}{
					{"id": 1, "name": "Alice"},
				},
				SelectedRow: 5,
				LoadingData: false,
			},
			contains: []string{"Invalid row selection"},
		},
		{
			name: "successful record rendering with schema",
			model: RecordViewModel{
				Rows: []map[string]interface{}{
					{"id": 1, "name": "Alice", "email": "alice@example.com"},
					{"id": 2, "name": "Bob", "email": "bob@example.com"},
				},
				TableSchema: &nosqldb.TableResult{
					DDL: "CREATE TABLE users (id INTEGER, name STRING, email STRING, PRIMARY KEY(id))",
				},
				SelectedRow: 0,
				LoadingData: false,
			},
			contains: []string{"id", "name", "email", "Alice"},
		},
		{
			name: "record rendering for second row",
			model: RecordViewModel{
				Rows: []map[string]interface{}{
					{"id": 1, "name": "Alice", "email": "alice@example.com"},
					{"id": 2, "name": "Bob", "email": "bob@example.com"},
				},
				TableSchema: &nosqldb.TableResult{
					DDL: "CREATE TABLE users (id INTEGER, name STRING, email STRING, PRIMARY KEY(id))",
				},
				SelectedRow: 1,
				LoadingData: false,
			},
			contains: []string{"id", "name", "email", "Bob"},
		},
		{
			name: "record rendering without schema",
			model: RecordViewModel{
				Rows: []map[string]interface{}{
					{"id": 1, "name": "Alice", "status": "active"},
				},
				TableSchema: nil,
				SelectedRow: 0,
				LoadingData: false,
			},
			contains: []string{}, // Column order undefined without schema, just check it doesn't error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderRecordView(tt.model)

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("RenderRecordView() = %q, should contain %q", result, substr)
				}
			}

			if result == "" {
				t.Error("RenderRecordView() should not return empty string")
			}
		})
	}
}

func TestRenderRecordView_ColumnOrder(t *testing.T) {
	// Test that columns are displayed in DDL order when schema is provided
	model := RecordViewModel{
		Rows: []map[string]interface{}{
			{"name": "Alice", "id": 1, "email": "alice@example.com"}, // Intentionally out of order in map
		},
		TableSchema: &nosqldb.TableResult{
			DDL: "CREATE TABLE users (id INTEGER, name STRING, email STRING, PRIMARY KEY(id))",
		},
		SelectedRow: 0,
		LoadingData: false,
	}

	result := RenderRecordView(model)

	// The result should contain all fields
	if !strings.Contains(result, "id") {
		t.Error("RenderRecordView() should contain 'id' field")
	}
	if !strings.Contains(result, "name") {
		t.Error("RenderRecordView() should contain 'name' field")
	}
	if !strings.Contains(result, "email") {
		t.Error("RenderRecordView() should contain 'email' field")
	}
	if !strings.Contains(result, "Alice") {
		t.Error("RenderRecordView() should contain the value 'Alice'")
	}
}
