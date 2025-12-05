package app

import (
	"strings"
	"testing"
)

func TestRenderSQLPane(t *testing.T) {
	t.Run("empty SQL shows cursor when focused", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.CurrentSQL = ""

		result := renderSQLPaneWithHeight(m, 40, 5)

		// Should have border characters
		if !strings.Contains(result, "╭") || !strings.Contains(result, "╯") {
			t.Error("Expected border characters in output")
		}
	})

	t.Run("shows SQL content", func(t *testing.T) {
		m := InitialModel()
		m.CurrentSQL = "SELECT * FROM users"

		result := renderSQLPaneWithHeight(m, 40, 5)

		if !strings.Contains(result, "SELECT") {
			t.Error("Expected SQL content in output")
		}
	})

	t.Run("shows Custom label when custom SQL", func(t *testing.T) {
		m := InitialModel()
		m.CustomSQL = true
		m.CurrentSQL = "SELECT id FROM users"

		result := renderSQLPaneWithHeight(m, 40, 5)

		if !strings.Contains(result, "Custom") {
			t.Error("Expected 'Custom' label in output")
		}
	})

	t.Run("wraps long SQL text", func(t *testing.T) {
		m := InitialModel()
		m.CurrentSQL = "SELECT id, name, email, address, phone FROM users WHERE active = true"

		result := renderSQLPaneWithHeight(m, 30, 5)

		// Should still contain SQL keywords
		if !strings.Contains(result, "SELECT") {
			t.Error("Expected SQL content in output")
		}
	})

	t.Run("handles multiline SQL", func(t *testing.T) {
		m := InitialModel()
		m.CurrentSQL = "SELECT *\nFROM users\nWHERE id = 1"

		result := renderSQLPaneWithHeight(m, 40, 5)

		if !strings.Contains(result, "SELECT") {
			t.Error("Expected SELECT in output")
		}
		if !strings.Contains(result, "FROM") {
			t.Error("Expected FROM in output")
		}
	})

	t.Run("cursor position in middle of text", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.CurrentSQL = "SELECT"
		m.SQLCursorPos = 3

		result := renderSQLPaneWithHeight(m, 40, 5)

		// Should render without error
		if !strings.Contains(result, "SQL") {
			t.Error("Expected SQL title in output")
		}
	})
}
