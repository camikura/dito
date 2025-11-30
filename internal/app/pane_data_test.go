package app

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"

	"github.com/camikura/dito/internal/db"
)

func TestRenderDataPane(t *testing.T) {
	t.Run("no table selected shows message", func(t *testing.T) {
		m := InitialModel()
		m.SelectedTable = -1

		result := renderDataPane(m, 60, 20)

		if !strings.Contains(result, "Select a table") {
			t.Error("Expected 'Select a table' in output")
		}
	})

	t.Run("shows loading message", func(t *testing.T) {
		m := InitialModel()
		m.Tables = []string{"users"}
		m.SelectedTable = 0
		m.LoadingData = true
		m.TableData = make(map[string]*db.TableDataResult)

		result := renderDataPane(m, 60, 20)

		if !strings.Contains(result, "Loading") {
			t.Error("Expected 'Loading' in output")
		}
	})

	t.Run("shows no rows message", func(t *testing.T) {
		m := InitialModel()
		m.Tables = []string{"users"}
		m.SelectedTable = 0
		m.TableData = map[string]*db.TableDataResult{
			"users": {Rows: []map[string]interface{}{}},
		}

		result := renderDataPane(m, 60, 20)

		if !strings.Contains(result, "No rows") {
			t.Error("Expected 'No rows' in output")
		}
	})

	t.Run("shows error message", func(t *testing.T) {
		m := InitialModel()
		m.Tables = []string{"users"}
		m.SelectedTable = 0
		m.DataErrorMsg = "Query failed"

		result := renderDataPane(m, 60, 20)

		if !strings.Contains(result, "Query failed") {
			t.Error("Expected error message in output")
		}
	})

	t.Run("shows table name in title", func(t *testing.T) {
		m := InitialModel()
		m.Tables = []string{"users"}
		m.SelectedTable = 0
		m.TableData = map[string]*db.TableDataResult{
			"users": {Rows: []map[string]interface{}{}},
		}

		result := renderDataPane(m, 60, 20)

		if !strings.Contains(result, "Data (users)") {
			t.Error("Expected table name in data title")
		}
	})

	t.Run("renders data rows", func(t *testing.T) {
		m := InitialModel()
		m.Tables = []string{"users"}
		m.SelectedTable = 0
		m.Width = 80
		m.Height = 30
		m.TableData = map[string]*db.TableDataResult{
			"users": {
				Rows: []map[string]interface{}{
					{"id": 1, "name": "Alice"},
					{"id": 2, "name": "Bob"},
				},
			},
		}
		m.TableDetails = map[string]*db.TableDetailsResult{
			"users": {TableName: "users"},
		}

		result := renderDataPane(m, 80, 30)

		// Should contain column headers or data
		if !strings.Contains(result, "Data") {
			t.Error("Expected Data title in output")
		}
	})

	t.Run("custom SQL shows extracted table name", func(t *testing.T) {
		m := InitialModel()
		m.CustomSQL = true
		m.CurrentSQL = "SELECT * FROM products"
		m.Tables = []string{"users", "products"}
		m.TableData = map[string]*db.TableDataResult{
			"products": {Rows: []map[string]interface{}{}},
		}

		result := renderDataPane(m, 60, 20)

		if !strings.Contains(result, "products") {
			t.Error("Expected extracted table name in output")
		}
	})
}

func TestRenderBottomBorderWithScrollbar(t *testing.T) {
	borderStyle := lipgloss.NewStyle()

	t.Run("renders border with minimal width", func(t *testing.T) {
		result := renderBottomBorderWithScrollbar(borderStyle, 10, 100, 50, 0)

		if !strings.Contains(result, "╰") || !strings.Contains(result, "╯") {
			t.Error("Expected border corners in output")
		}
	})

	t.Run("handles zero content width", func(t *testing.T) {
		result := renderBottomBorderWithScrollbar(borderStyle, 2, 0, 0, 0)

		// Should not panic and return something
		if result == "" {
			t.Error("Expected non-empty result")
		}
	})
}

func TestGetColumnTypes(t *testing.T) {
	t.Run("returns empty map when no schema", func(t *testing.T) {
		m := InitialModel()
		m.TableDetails = make(map[string]*db.TableDetailsResult)

		result := getColumnTypes(m, "users", []string{"id", "name"})

		if result == nil {
			t.Error("Expected non-nil map")
		}
	})
}
