package app

import (
	"strings"
	"testing"
)

func TestRenderConnectionPane(t *testing.T) {
	t.Run("not connected shows not configured", func(t *testing.T) {
		m := InitialModel()
		m.Connected = false
		m.Endpoint = ""

		result := renderConnectionPane(m, 30)

		if !strings.Contains(result, "not configured") {
			t.Error("Expected 'not configured' in output")
		}
	})

	t.Run("connected shows endpoint", func(t *testing.T) {
		m := InitialModel()
		m.Connected = true
		m.Endpoint = "localhost:8080"

		result := renderConnectionPane(m, 30)

		if !strings.Contains(result, "localhost:8080") {
			t.Error("Expected endpoint in output")
		}
		if !strings.Contains(result, "✓") {
			t.Error("Expected checkmark in output when connected")
		}
	})

	t.Run("shows error message when present", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionMsg = "Connection failed"

		result := renderConnectionPane(m, 30)

		if !strings.Contains(result, "Connection failed") {
			t.Error("Expected error message in output")
		}
	})

	t.Run("truncates long error message", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionMsg = "This is a very long error message that should be truncated"

		result := renderConnectionPane(m, 30)

		if !strings.Contains(result, "...") {
			t.Error("Expected truncated message with ellipsis")
		}
	})

	t.Run("active border when focused", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneConnection

		result := renderConnectionPane(m, 30)

		// Should contain border characters
		if !strings.Contains(result, "╭") || !strings.Contains(result, "╯") {
			t.Error("Expected border characters in output")
		}
	})
}
