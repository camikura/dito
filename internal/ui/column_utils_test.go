package ui

import (
	"testing"
)

func TestParsePrimaryKeysFromDDL(t *testing.T) {
	tests := []struct {
		name     string
		ddl      string
		expected []string
	}{
		{
			name:     "simple primary key",
			ddl:      "CREATE TABLE users (id INTEGER, name STRING, PRIMARY KEY(id))",
			expected: []string{"id"},
		},
		{
			name:     "composite primary key",
			ddl:      "CREATE TABLE orders (user_id INTEGER, order_id INTEGER, amount DOUBLE, PRIMARY KEY(user_id, order_id))",
			expected: []string{"user_id", "order_id"},
		},
		{
			name:     "primary key with SHARD",
			ddl:      "CREATE TABLE items (id INTEGER, name STRING, PRIMARY KEY(SHARD(id), name))",
			expected: []string{"id", "name"},
		},
		{
			name:     "no primary key",
			ddl:      "CREATE TABLE simple (id INTEGER, name STRING)",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParsePrimaryKeysFromDDL(tt.ddl)
			if len(result) != len(tt.expected) {
				t.Errorf("ParsePrimaryKeysFromDDL() got %v, want %v", result, tt.expected)
				return
			}
			for i, key := range result {
				if key != tt.expected[i] {
					t.Errorf("ParsePrimaryKeysFromDDL()[%d] = %q, want %q", i, key, tt.expected[i])
				}
			}
		})
	}
}

func TestParseColumnsFromDDL(t *testing.T) {
	tests := []struct {
		name        string
		ddl         string
		primaryKeys []string
		expected    []ColumnInfo
	}{
		{
			name:        "simple table",
			ddl:         "CREATE TABLE users (id INTEGER, name STRING, PRIMARY KEY(id))",
			primaryKeys: []string{"id"},
			expected: []ColumnInfo{
				{Name: "id", Type: "INTEGER", IsPrimaryKey: true},
				{Name: "name", Type: "STRING", IsPrimaryKey: false},
			},
		},
		{
			name:        "table with multiple types",
			ddl:         "CREATE TABLE products (id LONG, name STRING, price DOUBLE, active BOOLEAN, PRIMARY KEY(id))",
			primaryKeys: []string{"id"},
			expected: []ColumnInfo{
				{Name: "id", Type: "LONG", IsPrimaryKey: true},
				{Name: "name", Type: "STRING", IsPrimaryKey: false},
				{Name: "price", Type: "DOUBLE", IsPrimaryKey: false},
				{Name: "active", Type: "BOOLEAN", IsPrimaryKey: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseColumnsFromDDL(tt.ddl, tt.primaryKeys)
			if len(result) != len(tt.expected) {
				t.Errorf("ParseColumnsFromDDL() got %d columns, want %d", len(result), len(tt.expected))
				return
			}
			for i, col := range result {
				if col.Name != tt.expected[i].Name {
					t.Errorf("column[%d].Name = %q, want %q", i, col.Name, tt.expected[i].Name)
				}
				if col.Type != tt.expected[i].Type {
					t.Errorf("column[%d].Type = %q, want %q", i, col.Type, tt.expected[i].Type)
				}
				if col.IsPrimaryKey != tt.expected[i].IsPrimaryKey {
					t.Errorf("column[%d].IsPrimaryKey = %v, want %v", i, col.IsPrimaryKey, tt.expected[i].IsPrimaryKey)
				}
			}
		})
	}
}

func TestGetColumnsInSchemaOrder(t *testing.T) {
	tests := []struct {
		name     string
		ddl      string
		rows     []map[string]interface{}
		expected []string
	}{
		{
			name: "columns from DDL",
			ddl:  "CREATE TABLE users (id INTEGER, name STRING, email STRING, PRIMARY KEY(id))",
			rows: []map[string]interface{}{
				{"id": 1, "name": "Alice", "email": "alice@example.com"},
			},
			expected: []string{"id", "name", "email"},
		},
		{
			name: "extra columns from data",
			ddl:  "CREATE TABLE users (id INTEGER, name STRING, PRIMARY KEY(id))",
			rows: []map[string]interface{}{
				{"id": 1, "name": "Alice", "metadata": "{}"},
			},
			expected: []string{"id", "name", "metadata"},
		},
		{
			name: "no DDL, columns from data",
			ddl:  "",
			rows: []map[string]interface{}{
				{"b": 1, "a": 2, "c": 3},
			},
			expected: []string{"a", "b", "c"}, // sorted alphabetically
		},
		{
			name:     "empty DDL and no rows",
			ddl:      "",
			rows:     []map[string]interface{}{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetColumnsInSchemaOrder(tt.ddl, tt.rows)
			if len(result) != len(tt.expected) {
				t.Errorf("GetColumnsInSchemaOrder() got %v, want %v", result, tt.expected)
				return
			}
			for i, col := range result {
				if col != tt.expected[i] {
					t.Errorf("GetColumnsInSchemaOrder()[%d] = %q, want %q", i, col, tt.expected[i])
				}
			}
		})
	}
}

func TestGetColumnTypes(t *testing.T) {
	tests := []struct {
		name     string
		ddl      string
		expected map[string]string
	}{
		{
			name: "simple table",
			ddl:  "CREATE TABLE users (id INTEGER, name STRING, price DOUBLE, PRIMARY KEY(id))",
			expected: map[string]string{
				"id":    "INTEGER",
				"name":  "STRING",
				"price": "DOUBLE",
			},
		},
		{
			name:     "empty DDL",
			ddl:      "",
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetColumnTypes(tt.ddl)
			if len(result) != len(tt.expected) {
				t.Errorf("GetColumnTypes() got %d types, want %d", len(result), len(tt.expected))
				return
			}
			for col, typ := range tt.expected {
				if result[col] != typ {
					t.Errorf("GetColumnTypes()[%q] = %q, want %q", col, result[col], typ)
				}
			}
		})
	}
}
