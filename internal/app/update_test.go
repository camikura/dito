package app

import (
	"reflect"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/camikura/dito/internal/db"
)

func TestSortTablesForTree(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "single table",
			input:    []string{"users"},
			expected: []string{"users"},
		},
		{
			name:     "no child tables",
			input:    []string{"products", "users", "orders"},
			expected: []string{"orders", "products", "users"},
		},
		{
			name:     "parent before child",
			input:    []string{"users.phones", "users"},
			expected: []string{"users", "users.phones"},
		},
		{
			name:     "multiple children",
			input:    []string{"users.phones", "users", "users.addresses"},
			expected: []string{"users", "users.addresses", "users.phones"},
		},
		{
			name:     "mixed parents and children",
			input:    []string{"users.phones", "products", "users", "orders.items", "orders"},
			expected: []string{"orders", "orders.items", "products", "users", "users.phones"},
		},
		{
			name:     "already sorted",
			input:    []string{"orders", "users", "users.addresses"},
			expected: []string{"orders", "users", "users.addresses"},
		},
		{
			name:     "complex hierarchy",
			input:    []string{"a.b", "c", "a", "b.c", "b"},
			expected: []string{"a", "a.b", "b", "b.c", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy to avoid modifying the original
			input := make([]string, len(tt.input))
			copy(input, tt.input)

			result := sortTablesForTree(input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("sortTablesForTree(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSortTablesForTree_DoesNotModifyInput(t *testing.T) {
	input := []string{"users.phones", "users", "products"}
	original := make([]string, len(input))
	copy(original, input)

	_ = sortTablesForTree(input)

	// The original input should not be modified in order
	// Note: sortTablesForTree copies the input, so this should pass
	if !reflect.DeepEqual(input, original) {
		t.Errorf("sortTablesForTree modified original slice: got %v, want %v", input, original)
	}
}

func TestBuildDefaultSQL(t *testing.T) {
	tests := []struct {
		name      string
		tableName string
		ddl       string
		expected  string
	}{
		{
			name:      "no DDL",
			tableName: "users",
			ddl:       "",
			expected:  "SELECT * FROM users",
		},
		{
			name:      "single primary key",
			tableName: "users",
			ddl:       "CREATE TABLE users (id INTEGER, name STRING, PRIMARY KEY(id))",
			expected:  "SELECT * FROM users ORDER BY id",
		},
		{
			name:      "composite primary key",
			tableName: "orders",
			ddl:       "CREATE TABLE orders (user_id INTEGER, order_id INTEGER, amount DOUBLE, PRIMARY KEY(user_id, order_id))",
			expected:  "SELECT * FROM orders ORDER BY user_id, order_id",
		},
		{
			name:      "primary key with SHARD",
			tableName: "items",
			ddl:       "CREATE TABLE items (id INTEGER, name STRING, PRIMARY KEY(SHARD(id), name))",
			expected:  "SELECT * FROM items ORDER BY id, name",
		},
		{
			name:      "child table",
			tableName: "users.addresses",
			ddl:       "CREATE TABLE users.addresses (id INTEGER, street STRING, PRIMARY KEY(id))",
			expected:  "SELECT * FROM users.addresses ORDER BY id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildDefaultSQL(tt.tableName, tt.ddl)
			if result != tt.expected {
				t.Errorf("buildDefaultSQL(%q, %q) = %q, want %q", tt.tableName, tt.ddl, result, tt.expected)
			}
		})
	}
}

func TestMoveCursorUpInText(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		cursorPos int
		expected  int
	}{
		{
			name:      "single line - no change needed",
			text:      "hello",
			cursorPos: 3,
			expected:  3,
		},
		{
			name:      "two lines - move from second to first",
			text:      "hello\nworld",
			cursorPos: 8, // 'r' in world
			expected:  2, // same column in first line
		},
		{
			name:      "already on first line - go to start",
			text:      "hello\nworld",
			cursorPos: 3,
			expected:  0,
		},
		{
			name:      "first line shorter - go to end of first",
			text:      "hi\nhello world",
			cursorPos: 14, // column 11 of second line
			expected:  2,  // end of first line "hi"
		},
		{
			name:      "three lines - move from third to second",
			text:      "aaa\nbbb\nccc",
			cursorPos: 10, // 'c' at position 2
			expected:  6,  // 'b' at position 2
		},
		{
			name:      "empty text",
			text:      "",
			cursorPos: 0,
			expected:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := moveCursorUpInText(tt.text, tt.cursorPos)
			if result != tt.expected {
				t.Errorf("moveCursorUpInText(%q, %d) = %d, want %d", tt.text, tt.cursorPos, result, tt.expected)
			}
		})
	}
}

func TestMoveCursorDownInText(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		cursorPos int
		expected  int
	}{
		{
			name:      "single line - go to end",
			text:      "hello",
			cursorPos: 2,
			expected:  5,
		},
		{
			name:      "two lines - move from first to second",
			text:      "hello\nworld",
			cursorPos: 2, // 'l' in hello
			expected:  8, // same column in world
		},
		{
			name:      "already on last line - go to end",
			text:      "hello\nworld",
			cursorPos: 8,
			expected:  11, // end of text
		},
		{
			name:      "first line longer - go to end of second",
			text:      "hello world\nhi",
			cursorPos: 8, // middle of first line
			expected:  14, // end of "hi"
		},
		{
			name:      "three lines - move from first to second",
			text:      "aaa\nbbb\nccc",
			cursorPos: 2, // 'a' at position 2
			expected:  6, // 'b' at position 2
		},
		{
			name:      "empty text",
			text:      "",
			cursorPos: 0,
			expected:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := moveCursorDownInText(tt.text, tt.cursorPos)
			if result != tt.expected {
				t.Errorf("moveCursorDownInText(%q, %d) = %d, want %d", tt.text, tt.cursorPos, result, tt.expected)
			}
		})
	}
}

func TestUpdateWindowSize(t *testing.T) {
	m := InitialModel()

	// Send window size message
	newModel, _ := Update(m, tea.WindowSizeMsg{Width: 120, Height: 40})

	if newModel.Width != 120 {
		t.Errorf("Width = %d, want 120", newModel.Width)
	}
	if newModel.Height != 40 {
		t.Errorf("Height = %d, want 40", newModel.Height)
	}
}

func TestNextPrevPane(t *testing.T) {
	m := InitialModel()
	m.Connected = true // Enable pane switching

	// Test NextPane cycle
	panes := []FocusPane{
		FocusPaneConnection,
		FocusPaneTables,
		FocusPaneSchema,
		FocusPaneSQL,
		FocusPaneData,
		FocusPaneConnection, // wrap around
	}

	for i := 0; i < len(panes)-1; i++ {
		if m.CurrentPane != panes[i] {
			t.Errorf("Step %d: CurrentPane = %v, want %v", i, m.CurrentPane, panes[i])
		}
		m = m.NextPane()
	}

	// Verify wrap around
	if m.CurrentPane != FocusPaneConnection {
		t.Errorf("After full cycle: CurrentPane = %v, want %v", m.CurrentPane, FocusPaneConnection)
	}
}

func TestPaneSwitchingDisabledWhenDisconnected(t *testing.T) {
	m := InitialModel()
	m.Connected = false
	m.CurrentPane = FocusPaneConnection
	m.Width = 120
	m.Height = 40

	// Try to switch panes with Tab
	newModel, _ := Update(m, tea.KeyMsg{Type: tea.KeyTab})

	// Should stay on Connection pane
	if newModel.CurrentPane != FocusPaneConnection {
		t.Errorf("CurrentPane = %v, want %v (should not switch when disconnected)", newModel.CurrentPane, FocusPaneConnection)
	}
}

func TestGetColumnsInSchemaOrder(t *testing.T) {
	tests := []struct {
		name     string
		model    Model
		table    string
		rows     []map[string]interface{}
		expected []string
	}{
		{
			name: "with custom SQL column order",
			model: Model{
				CustomSQL:   true,
				ColumnOrder: []string{"name", "id", "email"},
			},
			table:    "users",
			rows:     []map[string]interface{}{{"id": 1, "name": "test", "email": "test@example.com"}},
			expected: []string{"name", "id", "email"},
		},
		{
			name: "without custom SQL - use row keys",
			model: Model{
				CustomSQL:   false,
				ColumnOrder: nil,
			},
			table:    "users",
			rows:     []map[string]interface{}{{"id": 1, "name": "test"}},
			expected: []string{"id", "name"}, // alphabetical order from rows
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getColumnsInSchemaOrder(tt.model, tt.table, tt.rows)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("getColumnsInSchemaOrder() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCalculateMaxHorizontalOffset(t *testing.T) {
	tests := []struct {
		name     string
		model    Model
		expected int
	}{
		{
			name: "no table selected",
			model: Model{
				SelectedTable: -1,
			},
			expected: 0,
		},
		{
			name: "table selected but no data",
			model: Model{
				Tables:        []string{"users"},
				SelectedTable: 0,
				TableData:     map[string]*db.TableDataResult{},
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateMaxHorizontalOffset(tt.model)
			if result != tt.expected {
				t.Errorf("calculateMaxHorizontalOffset() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestQuitKey(t *testing.T) {
	m := InitialModel()
	m.Width = 120
	m.Height = 40

	// Press 'q' to quit
	_, cmd := Update(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	// cmd should be tea.Quit
	if cmd == nil {
		t.Error("Expected quit command, got nil")
	}
}

func TestCtrlCQuit(t *testing.T) {
	m := InitialModel()
	m.Width = 120
	m.Height = 40

	// Press Ctrl+C to quit
	_, cmd := Update(m, tea.KeyMsg{Type: tea.KeyCtrlC})

	// cmd should be tea.Quit
	if cmd == nil {
		t.Error("Expected quit command, got nil")
	}
}
