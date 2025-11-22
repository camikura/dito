package views

import (
	"strings"
	"testing"

	"github.com/oracle/nosql-go-sdk/nosqldb"
)

func TestRenderSchemaView(t *testing.T) {
	tests := []struct {
		name     string
		model    SchemaViewModel
		contains []string
	}{
		{
			name: "parent table with children",
			model: SchemaViewModel{
				TableName: "users",
				AllTables: []string{"users", "users.orders", "users.addresses", "products"},
				TableSchema: &nosqldb.TableResult{
					DDL: "CREATE TABLE users (id INTEGER, name STRING, PRIMARY KEY(id))",
				},
				Indexes:        []nosqldb.IndexInfo{},
				LoadingDetails: false,
			},
			contains: []string{"Table:", "users", "Parent:", "-", "Children:", "orders, addresses", "Columns:", "id", "name", "Indexes:"},
		},
		{
			name: "child table",
			model: SchemaViewModel{
				TableName: "users.orders",
				AllTables: []string{"users", "users.orders", "products"},
				TableSchema: &nosqldb.TableResult{
					DDL: "CREATE TABLE users.orders (orderId INTEGER, amount DOUBLE, PRIMARY KEY(orderId))",
				},
				Indexes:        []nosqldb.IndexInfo{},
				LoadingDetails: false,
			},
			contains: []string{"Table:", "users.orders", "Parent:", "users", "Children:", "-", "Columns:"},
		},
		{
			name: "table with indexes",
			model: SchemaViewModel{
				TableName: "products",
				AllTables: []string{"products"},
				TableSchema: &nosqldb.TableResult{
					DDL: "CREATE TABLE products (id INTEGER, name STRING, PRIMARY KEY(id))",
				},
				Indexes: []nosqldb.IndexInfo{
					{IndexName: "name_idx", FieldNames: []string{"name"}},
					{IndexName: "composite_idx", FieldNames: []string{"name", "id"}},
				},
				LoadingDetails: false,
			},
			contains: []string{"Indexes:", "name_idx", "(name)", "composite_idx", "(name, id)"},
		},
		{
			name: "loading state",
			model: SchemaViewModel{
				TableName:      "users",
				AllTables:      []string{"users"},
				TableSchema:    nil,
				Indexes:        nil,
				LoadingDetails: true,
			},
			contains: []string{"Table:", "users", "Loading..."},
		},
		{
			name: "no schema available",
			model: SchemaViewModel{
				TableName:      "users",
				AllTables:      []string{"users"},
				TableSchema:    nil,
				Indexes:        nil,
				LoadingDetails: false,
			},
			contains: []string{"Table:", "users", "Schema information will be displayed here", "Index information will be displayed here"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderSchemaView(tt.model)

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("RenderSchemaView() = %q, should contain %q", result, substr)
				}
			}

			if result == "" {
				t.Error("RenderSchemaView() should not return empty string")
			}
		})
	}
}

func TestParsePrimaryKeysFromDDL(t *testing.T) {
	tests := []struct {
		name string
		ddl  string
		want []string
	}{
		{
			name: "simple primary key",
			ddl:  "CREATE TABLE users (id INTEGER, name STRING, PRIMARY KEY(id))",
			want: []string{"id"},
		},
		{
			name: "composite primary key",
			ddl:  "CREATE TABLE orders (userId INTEGER, orderId INTEGER, amount DOUBLE, PRIMARY KEY(userId, orderId))",
			want: []string{"userId", "orderId"},
		},
		{
			name: "primary key with SHARD",
			ddl:  "CREATE TABLE users (id INTEGER, name STRING, PRIMARY KEY(SHARD(id)))",
			want: []string{"id"},
		},
		{
			name: "composite primary key with SHARD",
			ddl:  "CREATE TABLE orders (userId INTEGER, orderId INTEGER, PRIMARY KEY(SHARD(userId), orderId))",
			want: []string{"userId", "orderId"},
		},
		{
			name: "no primary key",
			ddl:  "CREATE TABLE temp (data STRING)",
			want: []string{},
		},
		{
			name: "empty DDL",
			ddl:  "",
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParsePrimaryKeysFromDDL(tt.ddl)

			if len(got) != len(tt.want) {
				t.Errorf("ParsePrimaryKeysFromDDL() returned %d keys, want %d", len(got), len(tt.want))
				return
			}

			for i, key := range got {
				if key != tt.want[i] {
					t.Errorf("ParsePrimaryKeysFromDDL()[%d] = %q, want %q", i, key, tt.want[i])
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
		wantCount   int
		wantNames   []string
	}{
		{
			name:        "simple table",
			ddl:         "CREATE TABLE users (id INTEGER, name STRING, email STRING, PRIMARY KEY(id))",
			primaryKeys: []string{"id"},
			wantCount:   3,
			wantNames:   []string{"id", "name", "email"},
		},
		{
			name:        "table with composite primary key",
			ddl:         "CREATE TABLE orders (userId INTEGER, orderId INTEGER, amount DOUBLE, status STRING, PRIMARY KEY(userId, orderId))",
			primaryKeys: []string{"userId", "orderId"},
			wantCount:   4,
			wantNames:   []string{"userId", "orderId", "amount", "status"},
		},
		{
			name:        "empty DDL",
			ddl:         "",
			primaryKeys: []string{},
			wantCount:   0,
			wantNames:   []string{},
		},
		{
			name:        "DDL without columns",
			ddl:         "CREATE TABLE empty (PRIMARY KEY(id))",
			primaryKeys: []string{"id"},
			wantCount:   0,
			wantNames:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseColumnsFromDDL(tt.ddl, tt.primaryKeys)

			if len(got) != tt.wantCount {
				t.Errorf("ParseColumnsFromDDL() returned %d columns, want %d", len(got), tt.wantCount)
				return
			}

			for i, col := range got {
				if i < len(tt.wantNames) && col.Name != tt.wantNames[i] {
					t.Errorf("ParseColumnsFromDDL()[%d].Name = %q, want %q", i, col.Name, tt.wantNames[i])
				}
			}

			// Verify that primary key columns have the (Primary Key) suffix
			for _, col := range got {
				isPK := false
				for _, pk := range tt.primaryKeys {
					if col.Name == pk {
						isPK = true
						break
					}
				}
				if isPK && !strings.Contains(col.Type, "Primary Key") {
					t.Errorf("ParseColumnsFromDDL() column %q is a primary key but type %q doesn't contain 'Primary Key'", col.Name, col.Type)
				}
			}
		})
	}
}
