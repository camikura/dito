package app

import (
	"strings"
	"testing"

	"github.com/camikura/dito/internal/db"
	"github.com/oracle/nosql-go-sdk/nosqldb"
)

func TestRenderTablesPane(t *testing.T) {
	t.Run("no tables shows message", func(t *testing.T) {
		m := InitialModel()
		m.Tables.Tables = []string{}

		result := renderTablesPaneWithHeight(m, 30, 10)

		if !strings.Contains(result, "No tables") {
			t.Error("Expected 'No tables' in output")
		}
	})

	t.Run("shows table names", func(t *testing.T) {
		m := InitialModel()
		m.Tables.Tables = []string{"users", "products"}

		result := renderTablesPaneWithHeight(m, 30, 10)

		if !strings.Contains(result, "users") {
			t.Error("Expected 'users' in output")
		}
		if !strings.Contains(result, "products") {
			t.Error("Expected 'products' in output")
		}
	})

	t.Run("shows selection marker for selected table", func(t *testing.T) {
		m := InitialModel()
		m.Tables.Tables = []string{"users", "products"}
		m.Tables.SelectedTable = 0

		result := renderTablesPaneWithHeight(m, 30, 10)

		if !strings.Contains(result, "*") {
			t.Error("Expected selection marker '*' in output")
		}
	})

	t.Run("shows table count in title", func(t *testing.T) {
		m := InitialModel()
		m.Tables.Tables = []string{"users", "products", "orders"}

		result := renderTablesPaneWithHeight(m, 30, 10)

		if !strings.Contains(result, "Tables") {
			t.Error("Expected 'Tables' in title")
		}
	})

	t.Run("indents child tables", func(t *testing.T) {
		m := InitialModel()
		m.Tables.Tables = []string{"orders", "orders.items"}

		result := renderTablesPaneWithHeight(m, 30, 10)

		// Child table should show only "items" part
		if !strings.Contains(result, "items") {
			t.Error("Expected child table name 'items' in output")
		}
	})
}

func TestRenderSchemaPane(t *testing.T) {
	t.Run("no table selected shows message", func(t *testing.T) {
		m := InitialModel()
		m.Tables.SelectedTable = -1

		result := renderSchemaPaneWithHeight(m, 30, 10)

		if !strings.Contains(result, "Select a table") {
			t.Error("Expected 'Select a table' in output")
		}
	})

	t.Run("shows loading when no details", func(t *testing.T) {
		m := InitialModel()
		m.Tables.Tables = []string{"users"}
		m.Tables.SelectedTable = 0
		m.Schema.TableDetails = make(map[string]*db.TableDetailsResult)

		result := renderSchemaPaneWithHeight(m, 30, 10)

		if !strings.Contains(result, "Loading") {
			t.Error("Expected 'Loading' in output")
		}
	})

	t.Run("shows schema error message", func(t *testing.T) {
		m := InitialModel()
		m.Tables.Tables = []string{"users"}
		m.Tables.SelectedTable = 0
		m.Schema.ErrorMsg = "Failed to load schema"

		result := renderSchemaPaneWithHeight(m, 30, 10)

		if !strings.Contains(result, "Failed to load schema") {
			t.Error("Expected error message in output")
		}
	})

	t.Run("shows columns and indexes sections", func(t *testing.T) {
		m := InitialModel()
		m.Tables.Tables = []string{"users"}
		m.Tables.SelectedTable = 0
		m.Schema.TableDetails = map[string]*db.TableDetailsResult{
			"users": {
				TableName: "users",
				Schema: &nosqldb.TableResult{
					DDL: "CREATE TABLE users (id INTEGER, name STRING, PRIMARY KEY(id))",
				},
				Indexes: []nosqldb.IndexInfo{},
			},
		}

		result := renderSchemaPaneWithHeight(m, 40, 15)

		if !strings.Contains(result, "Columns") {
			t.Error("Expected 'Columns' section in output")
		}
		if !strings.Contains(result, "Indexes") {
			t.Error("Expected 'Indexes' section in output")
		}
	})

	t.Run("shows table name in title", func(t *testing.T) {
		m := InitialModel()
		m.Tables.Tables = []string{"users"}
		m.Tables.SelectedTable = 0
		m.Schema.TableDetails = map[string]*db.TableDetailsResult{
			"users": {
				TableName: "users",
				Schema: &nosqldb.TableResult{
					DDL: "CREATE TABLE users (id INTEGER, PRIMARY KEY(id))",
				},
			},
		}

		result := renderSchemaPaneWithHeight(m, 40, 10)

		if !strings.Contains(result, "Schema (users)") {
			t.Error("Expected table name in schema title")
		}
	})
}
