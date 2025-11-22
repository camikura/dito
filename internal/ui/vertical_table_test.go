package ui

import (
	"strings"
	"testing"
)

func TestVerticalTable_Render(t *testing.T) {
	tests := []struct {
		name     string
		table    VerticalTable
		contains []string
	}{
		{
			name: "simple vertical table",
			table: VerticalTable{
				Data: map[string]interface{}{
					"id":   1,
					"name": "Alice",
					"age":  30,
				},
				Keys: []string{"id", "name", "age"},
			},
			contains: []string{"id", "name", "age", "1", "Alice", "30"},
		},
		{
			name: "with different key lengths",
			table: VerticalTable{
				Data: map[string]interface{}{
					"id":          1,
					"first_name":  "Bob",
					"description": "Test user",
				},
				Keys: []string{"id", "first_name", "description"},
			},
			contains: []string{"id", "first_name", "description", "1", "Bob", "Test user"},
		},
		{
			name: "empty data",
			table: VerticalTable{
				Data: map[string]interface{}{},
				Keys: []string{},
			},
			contains: []string{"No data"},
		},
		{
			name: "nil data",
			table: VerticalTable{
				Data: nil,
				Keys: []string{"id", "name"},
			},
			contains: []string{"No data"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.table.Render()

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("Render() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestVerticalTable_NoTrailingNewline(t *testing.T) {
	table := VerticalTable{
		Data: map[string]interface{}{
			"id":   1,
			"name": "Test",
		},
		Keys: []string{"id", "name"},
	}

	result := table.Render()

	if strings.HasSuffix(result, "\n") {
		t.Errorf("Result should not have trailing newline")
	}
}

func TestVerticalTable_KeyAlignment(t *testing.T) {
	table := VerticalTable{
		Data: map[string]interface{}{
			"id":          1,
			"name":        "Alice",
			"description": "Test",
		},
		Keys: []string{"id", "name", "description"},
	}

	result := table.Render()
	lines := strings.Split(result, "\n")

	// All lines should exist (3 keys)
	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}

	// Each line should contain both key and value
	for _, line := range lines {
		if line == "" {
			t.Errorf("Line should not be empty")
		}
	}
}

func TestVerticalTable_LongValue(t *testing.T) {
	table := VerticalTable{
		Data: map[string]interface{}{
			"id":          1,
			"description": "This is a very long description that contains a lot of text",
		},
		Keys: []string{"id", "description"},
	}

	result := table.Render()

	// Should contain the full long value
	if !strings.Contains(result, "This is a very long description") {
		t.Errorf("Result should contain the full long value")
	}
}

func TestVerticalTable_KeyOrder(t *testing.T) {
	table := VerticalTable{
		Data: map[string]interface{}{
			"c": "third",
			"a": "first",
			"b": "second",
		},
		Keys: []string{"a", "b", "c"},
	}

	result := table.Render()
	lines := strings.Split(result, "\n")

	// Keys should appear in the specified order
	if !strings.Contains(lines[0], "a") {
		t.Errorf("First line should contain key 'a'")
	}
	if !strings.Contains(lines[1], "b") {
		t.Errorf("Second line should contain key 'b'")
	}
	if !strings.Contains(lines[2], "c") {
		t.Errorf("Third line should contain key 'c'")
	}
}
