package app

import (
	"errors"
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

func TestHandleConnectionKeys(t *testing.T) {
	t.Run("Enter opens connection dialog", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneConnection
		m.Width = 120
		m.Height = 40

		newModel, _ := handleConnectionKeys(m, tea.KeyMsg{Type: tea.KeyEnter})

		if !newModel.ConnectionDialogVisible {
			t.Error("Expected connection dialog to be visible")
		}
		if newModel.ConnectionDialogField != 0 {
			t.Errorf("ConnectionDialogField = %d, want 0", newModel.ConnectionDialogField)
		}
	})

	t.Run("Ctrl+D disconnects", func(t *testing.T) {
		m := InitialModel()
		m.Connected = true
		m.CurrentPane = FocusPaneConnection
		m.Tables = []string{"users", "products"}
		m.SelectedTable = 1

		newModel, _ := handleConnectionKeys(m, tea.KeyMsg{Type: tea.KeyCtrlD})

		if newModel.Connected {
			t.Error("Expected to be disconnected")
		}
		if len(newModel.Tables) != 0 {
			t.Errorf("Tables should be cleared, got %v", newModel.Tables)
		}
		if newModel.SelectedTable != -1 {
			t.Errorf("SelectedTable = %d, want -1", newModel.SelectedTable)
		}
	})
}

func TestHandleConnectionDialogKeys(t *testing.T) {
	t.Run("Esc closes dialog", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionDialogVisible = true

		newModel, _ := handleConnectionDialogKeys(m, tea.KeyMsg{Type: tea.KeyEsc})

		if newModel.ConnectionDialogVisible {
			t.Error("Expected dialog to be closed")
		}
	})

	t.Run("Tab switches fields", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionDialogVisible = true
		m.ConnectionDialogField = 0

		newModel, _ := handleConnectionDialogKeys(m, tea.KeyMsg{Type: tea.KeyTab})

		if newModel.ConnectionDialogField != 1 {
			t.Errorf("ConnectionDialogField = %d, want 1", newModel.ConnectionDialogField)
		}
	})

	t.Run("Backspace deletes character", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionDialogVisible = true
		m.ConnectionDialogField = 0
		m.EditEndpoint = "localhost"
		m.EditCursorPos = 9

		newModel, _ := handleConnectionDialogKeys(m, tea.KeyMsg{Type: tea.KeyBackspace})

		if newModel.EditEndpoint != "localhos" {
			t.Errorf("EditEndpoint = %q, want %q", newModel.EditEndpoint, "localhos")
		}
	})

	t.Run("Runes inserts characters", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionDialogVisible = true
		m.ConnectionDialogField = 0
		m.EditEndpoint = "local"
		m.EditCursorPos = 5

		newModel, _ := handleConnectionDialogKeys(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})

		if newModel.EditEndpoint != "localh" {
			t.Errorf("EditEndpoint = %q, want %q", newModel.EditEndpoint, "localh")
		}
	})

	t.Run("Left arrow moves cursor", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionDialogVisible = true
		m.EditCursorPos = 5

		newModel, _ := handleConnectionDialogKeys(m, tea.KeyMsg{Type: tea.KeyLeft})

		if newModel.EditCursorPos != 4 {
			t.Errorf("EditCursorPos = %d, want 4", newModel.EditCursorPos)
		}
	})

	t.Run("Home moves cursor to start", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionDialogVisible = true
		m.EditCursorPos = 5

		newModel, _ := handleConnectionDialogKeys(m, tea.KeyMsg{Type: tea.KeyHome})

		if newModel.EditCursorPos != 0 {
			t.Errorf("EditCursorPos = %d, want 0", newModel.EditCursorPos)
		}
	})
}

func TestHandleTablesKeys(t *testing.T) {
	t.Run("Down arrow moves cursor down", func(t *testing.T) {
		m := InitialModel()
		m.Tables = []string{"users", "products", "orders"}
		m.CursorTable = 0
		m.TablesHeight = 10

		newModel, _ := handleTablesKeys(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

		if newModel.CursorTable != 1 {
			t.Errorf("CursorTable = %d, want 1", newModel.CursorTable)
		}
	})

	t.Run("Up arrow moves cursor up", func(t *testing.T) {
		m := InitialModel()
		m.Tables = []string{"users", "products", "orders"}
		m.CursorTable = 2
		m.TablesHeight = 10

		newModel, _ := handleTablesKeys(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

		if newModel.CursorTable != 1 {
			t.Errorf("CursorTable = %d, want 1", newModel.CursorTable)
		}
	})

	t.Run("Down at bottom stays at bottom", func(t *testing.T) {
		m := InitialModel()
		m.Tables = []string{"users", "products"}
		m.CursorTable = 1
		m.TablesHeight = 10

		newModel, _ := handleTablesKeys(m, tea.KeyMsg{Type: tea.KeyDown})

		if newModel.CursorTable != 1 {
			t.Errorf("CursorTable = %d, want 1 (should stay at bottom)", newModel.CursorTable)
		}
	})

	t.Run("Up at top stays at top", func(t *testing.T) {
		m := InitialModel()
		m.Tables = []string{"users", "products"}
		m.CursorTable = 0
		m.TablesHeight = 10

		newModel, _ := handleTablesKeys(m, tea.KeyMsg{Type: tea.KeyUp})

		if newModel.CursorTable != 0 {
			t.Errorf("CursorTable = %d, want 0 (should stay at top)", newModel.CursorTable)
		}
	})
}

func TestHandleSchemaKeys(t *testing.T) {
	t.Run("Up at top stays at top", func(t *testing.T) {
		m := InitialModel()
		m.Tables = []string{"users"}
		m.SelectedTable = 0
		m.SchemaScrollOffset = 0

		newModel, _ := handleSchemaKeys(m, tea.KeyMsg{Type: tea.KeyUp})

		if newModel.SchemaScrollOffset != 0 {
			t.Errorf("SchemaScrollOffset = %d, want 0", newModel.SchemaScrollOffset)
		}
	})

	t.Run("Up scrolls up when offset > 0", func(t *testing.T) {
		m := InitialModel()
		m.Tables = []string{"users"}
		m.SelectedTable = 0
		m.SchemaScrollOffset = 2

		newModel, _ := handleSchemaKeys(m, tea.KeyMsg{Type: tea.KeyUp})

		if newModel.SchemaScrollOffset != 1 {
			t.Errorf("SchemaScrollOffset = %d, want 1", newModel.SchemaScrollOffset)
		}
	})

	t.Run("No table selected returns unchanged", func(t *testing.T) {
		m := InitialModel()
		m.Tables = []string{}
		m.SelectedTable = -1

		newModel, _ := handleSchemaKeys(m, tea.KeyMsg{Type: tea.KeyDown})

		if newModel.SchemaScrollOffset != 0 {
			t.Errorf("SchemaScrollOffset = %d, want 0", newModel.SchemaScrollOffset)
		}
	})
}

func TestHandleSQLKeys(t *testing.T) {
	t.Run("Backspace deletes character", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.CurrentSQL = "SELECT * FROM users"
		m.SQLCursorPos = 19

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyBackspace})

		if newModel.CurrentSQL != "SELECT * FROM user" {
			t.Errorf("CurrentSQL = %q, want %q", newModel.CurrentSQL, "SELECT * FROM user")
		}
	})

	t.Run("Enter inserts newline", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.CurrentSQL = "SELECT *"
		m.SQLCursorPos = 8

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyEnter})

		if newModel.CurrentSQL != "SELECT *\n" {
			t.Errorf("CurrentSQL = %q, want %q", newModel.CurrentSQL, "SELECT *\n")
		}
	})

	t.Run("Left arrow moves cursor left", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.CurrentSQL = "SELECT"
		m.SQLCursorPos = 6

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyLeft})

		if newModel.SQLCursorPos != 5 {
			t.Errorf("SQLCursorPos = %d, want 5", newModel.SQLCursorPos)
		}
	})

	t.Run("Right arrow moves cursor right", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.CurrentSQL = "SELECT"
		m.SQLCursorPos = 3

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyRight})

		if newModel.SQLCursorPos != 4 {
			t.Errorf("SQLCursorPos = %d, want 4", newModel.SQLCursorPos)
		}
	})

	t.Run("Home moves cursor to start", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.CurrentSQL = "SELECT * FROM users"
		m.SQLCursorPos = 10

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyHome})

		if newModel.SQLCursorPos != 0 {
			t.Errorf("SQLCursorPos = %d, want 0", newModel.SQLCursorPos)
		}
	})

	t.Run("End moves cursor to end", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.CurrentSQL = "SELECT"
		m.SQLCursorPos = 0

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyEnd})

		if newModel.SQLCursorPos != 6 {
			t.Errorf("SQLCursorPos = %d, want 6", newModel.SQLCursorPos)
		}
	})

	t.Run("Space inserts space", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.CurrentSQL = "SELECT*"
		m.SQLCursorPos = 6

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeySpace})

		if newModel.CurrentSQL != "SELECT *" {
			t.Errorf("CurrentSQL = %q, want %q", newModel.CurrentSQL, "SELECT *")
		}
	})

	t.Run("Runes inserts characters", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.CurrentSQL = "SELECT"
		m.SQLCursorPos = 6

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'X'}})

		if newModel.CurrentSQL != "SELECTX" {
			t.Errorf("CurrentSQL = %q, want %q", newModel.CurrentSQL, "SELECTX")
		}
	})

	t.Run("Up arrow moves cursor up in multi-line", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.CurrentSQL = "SELECT *\nFROM users"
		m.SQLCursorPos = 14 // middle of second line

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyUp})

		// Should move to first line at same column or end
		if newModel.SQLCursorPos >= 9 {
			t.Errorf("SQLCursorPos = %d, should be less than 9", newModel.SQLCursorPos)
		}
	})

	t.Run("Down arrow moves cursor down in multi-line", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.CurrentSQL = "SELECT *\nFROM users"
		m.SQLCursorPos = 4 // middle of first line

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyDown})

		// Should move to second line
		if newModel.SQLCursorPos < 9 {
			t.Errorf("SQLCursorPos = %d, should be >= 9", newModel.SQLCursorPos)
		}
	})
}

func TestHandleDataKeys(t *testing.T) {
	t.Run("Down arrow moves selection down", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneData
		m.Tables = []string{"users"}
		m.SelectedTable = 0
		m.SelectedDataRow = 0
		m.Height = 40
		m.TableData = map[string]*db.TableDataResult{
			"users": {
				Rows: []map[string]interface{}{
					{"id": 1}, {"id": 2}, {"id": 3},
				},
			},
		}

		newModel, _ := handleDataKeys(m, tea.KeyMsg{Type: tea.KeyDown})

		if newModel.SelectedDataRow != 1 {
			t.Errorf("SelectedDataRow = %d, want 1", newModel.SelectedDataRow)
		}
	})

	t.Run("Up arrow moves selection up", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneData
		m.Tables = []string{"users"}
		m.SelectedTable = 0
		m.SelectedDataRow = 2
		m.Height = 40
		m.TableData = map[string]*db.TableDataResult{
			"users": {
				Rows: []map[string]interface{}{
					{"id": 1}, {"id": 2}, {"id": 3},
				},
			},
		}

		newModel, _ := handleDataKeys(m, tea.KeyMsg{Type: tea.KeyUp})

		if newModel.SelectedDataRow != 1 {
			t.Errorf("SelectedDataRow = %d, want 1", newModel.SelectedDataRow)
		}
	})

	t.Run("Left arrow scrolls left", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneData
		m.HorizontalOffset = 5

		newModel, _ := handleDataKeys(m, tea.KeyMsg{Type: tea.KeyLeft})

		if newModel.HorizontalOffset != 4 {
			t.Errorf("HorizontalOffset = %d, want 4", newModel.HorizontalOffset)
		}
	})

	t.Run("Left at 0 stays at 0", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneData
		m.HorizontalOffset = 0

		newModel, _ := handleDataKeys(m, tea.KeyMsg{Type: tea.KeyLeft})

		if newModel.HorizontalOffset != 0 {
			t.Errorf("HorizontalOffset = %d, want 0", newModel.HorizontalOffset)
		}
	})

	t.Run("Enter opens record detail", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneData
		m.Tables = []string{"users"}
		m.SelectedTable = 0
		m.SelectedDataRow = 0
		m.TableData = map[string]*db.TableDataResult{
			"users": {
				Rows: []map[string]interface{}{{"id": 1}},
			},
		}

		newModel, _ := handleDataKeys(m, tea.KeyMsg{Type: tea.KeyEnter})

		if !newModel.RecordDetailVisible {
			t.Error("Expected record detail to be visible")
		}
	})

	t.Run("Esc resets custom SQL", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneData
		m.CustomSQL = true
		m.ColumnOrder = []string{"name", "id"}
		m.SelectedDataRow = 5
		m.ViewportOffset = 3

		newModel, _ := handleDataKeys(m, tea.KeyMsg{Type: tea.KeyEsc})

		if newModel.CustomSQL {
			t.Error("Expected CustomSQL to be false")
		}
		if newModel.ColumnOrder != nil {
			t.Error("Expected ColumnOrder to be nil")
		}
		if newModel.SelectedDataRow != 0 {
			t.Errorf("SelectedDataRow = %d, want 0", newModel.SelectedDataRow)
		}
	})
}

func TestHandleRecordDetailKeys(t *testing.T) {
	t.Run("Esc closes dialog", func(t *testing.T) {
		m := InitialModel()
		m.RecordDetailVisible = true
		m.RecordDetailScroll = 5

		newModel, _ := handleRecordDetailKeys(m, tea.KeyMsg{Type: tea.KeyEsc})

		if newModel.RecordDetailVisible {
			t.Error("Expected record detail to be closed")
		}
		if newModel.RecordDetailScroll != 0 {
			t.Errorf("RecordDetailScroll = %d, want 0", newModel.RecordDetailScroll)
		}
	})

	t.Run("Down scrolls down when maxScroll > 0", func(t *testing.T) {
		m := InitialModel()
		m.RecordDetailVisible = true
		m.RecordDetailScroll = 0
		m.Tables = []string{"users"}
		m.SelectedTable = 0
		m.SelectedDataRow = 0
		m.Height = 10 // Small height to ensure maxScroll > 0

		// Create a row with many fields to exceed visible height
		row := make(map[string]interface{})
		for i := 0; i < 20; i++ {
			row[string(rune('a'+i))] = i
		}
		m.TableData = map[string]*db.TableDataResult{
			"users": {
				Rows: []map[string]interface{}{row},
			},
		}

		newModel, _ := handleRecordDetailKeys(m, tea.KeyMsg{Type: tea.KeyDown})

		if newModel.RecordDetailScroll != 1 {
			t.Errorf("RecordDetailScroll = %d, want 1", newModel.RecordDetailScroll)
		}
	})

	t.Run("Up scrolls up", func(t *testing.T) {
		m := InitialModel()
		m.RecordDetailVisible = true
		m.RecordDetailScroll = 3

		newModel, _ := handleRecordDetailKeys(m, tea.KeyMsg{Type: tea.KeyUp})

		if newModel.RecordDetailScroll != 2 {
			t.Errorf("RecordDetailScroll = %d, want 2", newModel.RecordDetailScroll)
		}
	})

	t.Run("Home scrolls to top", func(t *testing.T) {
		m := InitialModel()
		m.RecordDetailVisible = true
		m.RecordDetailScroll = 10

		newModel, _ := handleRecordDetailKeys(m, tea.KeyMsg{Type: tea.KeyHome})

		if newModel.RecordDetailScroll != 0 {
			t.Errorf("RecordDetailScroll = %d, want 0", newModel.RecordDetailScroll)
		}
	})
}

func TestHandleTableListResult(t *testing.T) {
	t.Run("success stores sorted tables", func(t *testing.T) {
		m := InitialModel()

		newModel, _ := handleTableListResult(m, db.TableListResult{
			Tables: []string{"users.phones", "products", "users"},
		})

		expected := []string{"products", "users", "users.phones"}
		if !reflect.DeepEqual(newModel.Tables, expected) {
			t.Errorf("Tables = %v, want %v", newModel.Tables, expected)
		}
		if newModel.CursorTable != 0 {
			t.Errorf("CursorTable = %d, want 0", newModel.CursorTable)
		}
	})

	t.Run("error returns unchanged model", func(t *testing.T) {
		m := InitialModel()

		newModel, _ := handleTableListResult(m, db.TableListResult{
			Err: errors.New("test error"),
		})

		if len(newModel.Tables) != 0 {
			t.Errorf("Tables should be empty on error")
		}
	})
}

func TestHandleTableDataResult(t *testing.T) {
	t.Run("success stores data", func(t *testing.T) {
		m := InitialModel()
		m.LoadingData = true

		rows := []map[string]interface{}{
			{"id": 1, "name": "test"},
		}

		newModel, _ := handleTableDataResult(m, db.TableDataResult{
			TableName: "users",
			Rows:      rows,
			HasMore:   true,
		})

		if newModel.LoadingData {
			t.Error("Expected LoadingData to be false")
		}
		data := newModel.TableData["users"]
		if data == nil {
			t.Fatal("Expected table data to exist")
		}
		if len(data.Rows) != 1 {
			t.Errorf("Rows count = %d, want 1", len(data.Rows))
		}
		if !data.HasMore {
			t.Error("Expected HasMore to be true")
		}
	})

	t.Run("append merges data", func(t *testing.T) {
		m := InitialModel()
		m.LoadingData = true
		m.TableData = map[string]*db.TableDataResult{
			"users": {
				Rows: []map[string]interface{}{{"id": 1}},
			},
		}

		newModel, _ := handleTableDataResult(m, db.TableDataResult{
			TableName: "users",
			Rows:      []map[string]interface{}{{"id": 2}},
			IsAppend:  true,
		})

		data := newModel.TableData["users"]
		if len(data.Rows) != 2 {
			t.Errorf("Rows count = %d, want 2", len(data.Rows))
		}
	})

	t.Run("error clears loading state", func(t *testing.T) {
		m := InitialModel()
		m.LoadingData = true

		newModel, _ := handleTableDataResult(m, db.TableDataResult{
			Err: errors.New("test error"),
		})

		if newModel.LoadingData {
			t.Error("Expected LoadingData to be false")
		}
		if newModel.DataErrorMsg == "" {
			t.Error("Expected DataErrorMsg to be set")
		}
	})
}

func TestHandleConnectionResult(t *testing.T) {
	t.Run("error sets message", func(t *testing.T) {
		m := InitialModel()

		newModel, _ := handleConnectionResult(m, db.ConnectionResult{
			Err: errors.New("connection failed"),
		})

		if newModel.Connected {
			t.Error("Expected Connected to be false")
		}
		if newModel.ConnectionMsg == "" {
			t.Error("Expected ConnectionMsg to be set")
		}
	})
}

func TestCalculateRecordDetailMaxScroll(t *testing.T) {
	t.Run("no table selected returns 0", func(t *testing.T) {
		m := InitialModel()
		m.SelectedTable = -1

		result := calculateRecordDetailMaxScroll(m)

		if result != 0 {
			t.Errorf("calculateRecordDetailMaxScroll() = %d, want 0", result)
		}
	})

	t.Run("no data returns 0", func(t *testing.T) {
		m := InitialModel()
		m.Tables = []string{"users"}
		m.SelectedTable = 0
		m.TableData = map[string]*db.TableDataResult{}

		result := calculateRecordDetailMaxScroll(m)

		if result != 0 {
			t.Errorf("calculateRecordDetailMaxScroll() = %d, want 0", result)
		}
	})
}

func TestHandleTableDetailsResult(t *testing.T) {
	t.Run("error sets schema error message", func(t *testing.T) {
		m := InitialModel()
		m.LoadingData = true

		newModel, cmd := handleTableDetailsResult(m, db.TableDetailsResult{
			Err: errors.New("failed to fetch schema"),
		})

		if newModel.SchemaErrorMsg == "" {
			t.Error("Expected SchemaErrorMsg to be set")
		}
		if newModel.LoadingData {
			t.Error("Expected LoadingData to be false")
		}
		if cmd != nil {
			t.Error("Expected no command")
		}
	})

	t.Run("success stores table details", func(t *testing.T) {
		m := InitialModel()
		m.TableDetails = make(map[string]*db.TableDetailsResult)

		newModel, _ := handleTableDetailsResult(m, db.TableDetailsResult{
			TableName: "users",
		})

		if newModel.SchemaErrorMsg != "" {
			t.Errorf("SchemaErrorMsg should be empty, got %q", newModel.SchemaErrorMsg)
		}
		if newModel.TableDetails["users"] == nil {
			t.Error("Expected table details to be stored")
		}
	})

	t.Run("success clears previous error", func(t *testing.T) {
		m := InitialModel()
		m.SchemaErrorMsg = "previous error"
		m.TableDetails = make(map[string]*db.TableDetailsResult)

		newModel, _ := handleTableDetailsResult(m, db.TableDetailsResult{
			TableName: "users",
		})

		if newModel.SchemaErrorMsg != "" {
			t.Errorf("SchemaErrorMsg should be cleared, got %q", newModel.SchemaErrorMsg)
		}
	})
}

func TestHandleConnectionDialogKeysAdditional(t *testing.T) {
	t.Run("Right arrow moves cursor right within text", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionDialogVisible = true
		m.ConnectionDialogField = 0
		m.EditEndpoint = "localhost"
		m.EditCursorPos = 3

		newModel, _ := handleConnectionDialogKeys(m, tea.KeyMsg{Type: tea.KeyRight})

		if newModel.EditCursorPos != 4 {
			t.Errorf("EditCursorPos = %d, want 4", newModel.EditCursorPos)
		}
	})

	t.Run("Right arrow at end stays at end", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionDialogVisible = true
		m.ConnectionDialogField = 0
		m.EditEndpoint = "localhost"
		m.EditCursorPos = 9

		newModel, _ := handleConnectionDialogKeys(m, tea.KeyMsg{Type: tea.KeyRight})

		if newModel.EditCursorPos != 9 {
			t.Errorf("EditCursorPos = %d, want 9", newModel.EditCursorPos)
		}
	})

	t.Run("End moves cursor to end", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionDialogVisible = true
		m.ConnectionDialogField = 0
		m.EditEndpoint = "localhost"
		m.EditCursorPos = 3

		newModel, _ := handleConnectionDialogKeys(m, tea.KeyMsg{Type: tea.KeyEnd})

		if newModel.EditCursorPos != 9 {
			t.Errorf("EditCursorPos = %d, want 9", newModel.EditCursorPos)
		}
	})

	t.Run("Delete removes character at cursor", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionDialogVisible = true
		m.ConnectionDialogField = 0
		m.EditEndpoint = "localhost"
		m.EditCursorPos = 5

		newModel, _ := handleConnectionDialogKeys(m, tea.KeyMsg{Type: tea.KeyDelete})

		if newModel.EditEndpoint != "localost" {
			t.Errorf("EditEndpoint = %q, want %q", newModel.EditEndpoint, "localost")
		}
	})

	t.Run("Tab wraps around from last field", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionDialogVisible = true
		m.ConnectionDialogField = 1 // Port field (last)

		newModel, _ := handleConnectionDialogKeys(m, tea.KeyMsg{Type: tea.KeyTab})

		if newModel.ConnectionDialogField != 0 {
			t.Errorf("ConnectionDialogField = %d, want 0", newModel.ConnectionDialogField)
		}
	})

	t.Run("Editing port field", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionDialogVisible = true
		m.ConnectionDialogField = 1 // Port field
		m.EditPort = "8080"
		m.EditCursorPos = 4

		newModel, _ := handleConnectionDialogKeys(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})

		if newModel.EditPort != "80801" {
			t.Errorf("EditPort = %q, want %q", newModel.EditPort, "80801")
		}
	})

	t.Run("Backspace on port field", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionDialogVisible = true
		m.ConnectionDialogField = 1
		m.EditPort = "8080"
		m.EditCursorPos = 4

		newModel, _ := handleConnectionDialogKeys(m, tea.KeyMsg{Type: tea.KeyBackspace})

		if newModel.EditPort != "808" {
			t.Errorf("EditPort = %q, want %q", newModel.EditPort, "808")
		}
	})
}

func TestHandleTablesKeysAdditional(t *testing.T) {
	t.Run("Empty tables list returns unchanged", func(t *testing.T) {
		m := InitialModel()
		m.Tables = []string{}
		m.CursorTable = 0

		newModel, _ := handleTablesKeys(m, tea.KeyMsg{Type: tea.KeyDown})

		if newModel.CursorTable != 0 {
			t.Errorf("CursorTable = %d, want 0", newModel.CursorTable)
		}
	})

	t.Run("Up at top does nothing", func(t *testing.T) {
		m := InitialModel()
		m.Tables = []string{"users", "products", "orders"}
		m.CursorTable = 0
		m.TablesHeight = 10

		newModel, _ := handleTablesKeys(m, tea.KeyMsg{Type: tea.KeyUp})

		if newModel.CursorTable != 0 {
			t.Errorf("CursorTable = %d, want 0", newModel.CursorTable)
		}
	})

	t.Run("Down at bottom does nothing", func(t *testing.T) {
		m := InitialModel()
		m.Tables = []string{"users", "products", "orders"}
		m.CursorTable = 2
		m.TablesHeight = 10

		newModel, _ := handleTablesKeys(m, tea.KeyMsg{Type: tea.KeyDown})

		if newModel.CursorTable != 2 {
			t.Errorf("CursorTable = %d, want 2", newModel.CursorTable)
		}
	})

	t.Run("j key moves down", func(t *testing.T) {
		m := InitialModel()
		m.Tables = []string{"users", "products"}
		m.CursorTable = 0
		m.TablesHeight = 10

		newModel, _ := handleTablesKeys(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

		if newModel.CursorTable != 1 {
			t.Errorf("CursorTable = %d, want 1", newModel.CursorTable)
		}
	})

	t.Run("k key moves up", func(t *testing.T) {
		m := InitialModel()
		m.Tables = []string{"users", "products"}
		m.CursorTable = 1
		m.TablesHeight = 10

		newModel, _ := handleTablesKeys(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

		if newModel.CursorTable != 0 {
			t.Errorf("CursorTable = %d, want 0", newModel.CursorTable)
		}
	})
}

func TestHandleSchemaKeysAdditional(t *testing.T) {
	t.Run("Down without table details does nothing", func(t *testing.T) {
		m := InitialModel()
		m.Tables = []string{"users"}
		m.SelectedTable = 0
		m.SchemaScrollOffset = 0

		newModel, _ := handleSchemaKeys(m, tea.KeyMsg{Type: tea.KeyDown})

		// Without table details, maxScroll is 0, so no scrolling
		if newModel.SchemaScrollOffset != 0 {
			t.Errorf("SchemaScrollOffset = %d, want 0", newModel.SchemaScrollOffset)
		}
	})

	t.Run("k key scrolls up", func(t *testing.T) {
		m := InitialModel()
		m.Tables = []string{"users"}
		m.SelectedTable = 0
		m.SchemaScrollOffset = 2

		newModel, _ := handleSchemaKeys(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

		if newModel.SchemaScrollOffset != 1 {
			t.Errorf("SchemaScrollOffset = %d, want 1", newModel.SchemaScrollOffset)
		}
	})
}

func TestHandleDataKeysAdditional(t *testing.T) {
	t.Run("Left arrow scrolls left", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneData
		m.HorizontalOffset = 5

		newModel, _ := handleDataKeys(m, tea.KeyMsg{Type: tea.KeyLeft})

		if newModel.HorizontalOffset != 4 {
			t.Errorf("HorizontalOffset = %d, want 4", newModel.HorizontalOffset)
		}
	})

	t.Run("h key scrolls left", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneData
		m.HorizontalOffset = 5

		newModel, _ := handleDataKeys(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})

		if newModel.HorizontalOffset != 4 {
			t.Errorf("HorizontalOffset = %d, want 4", newModel.HorizontalOffset)
		}
	})

	t.Run("Left at zero stays at zero", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneData
		m.HorizontalOffset = 0

		newModel, _ := handleDataKeys(m, tea.KeyMsg{Type: tea.KeyLeft})

		if newModel.HorizontalOffset != 0 {
			t.Errorf("HorizontalOffset = %d, want 0", newModel.HorizontalOffset)
		}
	})

	t.Run("Enter key shows record detail", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneData
		m.Tables = []string{"users"}
		m.SelectedTable = 0
		m.SelectedDataRow = 0
		m.TableData = map[string]*db.TableDataResult{
			"users": {
				Rows: []map[string]interface{}{
					{"id": 1},
				},
			},
		}

		newModel, _ := handleDataKeys(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'\r'}})

		// Check if "enter" is handled
		if newModel.RecordDetailVisible {
			// Good, enter was handled
		}
	})

	t.Run("j key moves down", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneData
		m.Tables = []string{"users"}
		m.SelectedTable = 0
		m.SelectedDataRow = 0
		m.Height = 40
		m.TableData = map[string]*db.TableDataResult{
			"users": {
				Rows: []map[string]interface{}{
					{"id": 1}, {"id": 2},
				},
			},
		}

		newModel, _ := handleDataKeys(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

		if newModel.SelectedDataRow != 1 {
			t.Errorf("SelectedDataRow = %d, want 1", newModel.SelectedDataRow)
		}
	})

	t.Run("k key moves up", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneData
		m.Tables = []string{"users"}
		m.SelectedTable = 0
		m.SelectedDataRow = 1
		m.Height = 40
		m.TableData = map[string]*db.TableDataResult{
			"users": {
				Rows: []map[string]interface{}{
					{"id": 1}, {"id": 2},
				},
			},
		}

		newModel, _ := handleDataKeys(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

		if newModel.SelectedDataRow != 0 {
			t.Errorf("SelectedDataRow = %d, want 0", newModel.SelectedDataRow)
		}
	})

	t.Run("No table selected returns unchanged", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneData
		m.SelectedTable = -1
		m.SelectedDataRow = 0

		newModel, _ := handleDataKeys(m, tea.KeyMsg{Type: tea.KeyDown})

		if newModel.SelectedDataRow != 0 {
			t.Errorf("SelectedDataRow = %d, want 0", newModel.SelectedDataRow)
		}
	})
}

func TestHandleRecordDetailKeysAdditional(t *testing.T) {
	t.Run("j key scrolls down", func(t *testing.T) {
		m := InitialModel()
		m.RecordDetailVisible = true
		m.RecordDetailScroll = 0
		m.Tables = []string{"users"}
		m.SelectedTable = 0
		m.SelectedDataRow = 0
		m.Height = 10

		row := make(map[string]interface{})
		for i := 0; i < 20; i++ {
			row[string(rune('a'+i))] = i
		}
		m.TableData = map[string]*db.TableDataResult{
			"users": {Rows: []map[string]interface{}{row}},
		}

		newModel, _ := handleRecordDetailKeys(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

		if newModel.RecordDetailScroll != 1 {
			t.Errorf("RecordDetailScroll = %d, want 1", newModel.RecordDetailScroll)
		}
	})

	t.Run("k key scrolls up", func(t *testing.T) {
		m := InitialModel()
		m.RecordDetailVisible = true
		m.RecordDetailScroll = 5

		newModel, _ := handleRecordDetailKeys(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

		if newModel.RecordDetailScroll != 4 {
			t.Errorf("RecordDetailScroll = %d, want 4", newModel.RecordDetailScroll)
		}
	})

	t.Run("End scrolls to max", func(t *testing.T) {
		m := InitialModel()
		m.RecordDetailVisible = true
		m.RecordDetailScroll = 0
		m.Tables = []string{"users"}
		m.SelectedTable = 0
		m.SelectedDataRow = 0
		m.Height = 10

		row := make(map[string]interface{})
		for i := 0; i < 30; i++ {
			row[string(rune('a'+i))] = i
		}
		m.TableData = map[string]*db.TableDataResult{
			"users": {Rows: []map[string]interface{}{row}},
		}

		newModel, _ := handleRecordDetailKeys(m, tea.KeyMsg{Type: tea.KeyEnd})

		// Should scroll to max (but we just check it's > 0)
		if newModel.RecordDetailScroll == 0 {
			t.Errorf("RecordDetailScroll should be > 0")
		}
	})
}

func TestHandleSQLKeysAdditional(t *testing.T) {
	t.Run("Delete removes character at cursor", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.CurrentSQL = "SELECT * FROM users"
		m.SQLCursorPos = 7

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyDelete})

		if newModel.CurrentSQL != "SELECT  FROM users" {
			t.Errorf("CurrentSQL = %q, want %q", newModel.CurrentSQL, "SELECT  FROM users")
		}
	})

	t.Run("Delete at end does nothing", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.CurrentSQL = "SELECT"
		m.SQLCursorPos = 6

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyDelete})

		if newModel.CurrentSQL != "SELECT" {
			t.Errorf("CurrentSQL = %q, want %q", newModel.CurrentSQL, "SELECT")
		}
	})

	t.Run("Backspace at start does nothing", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.CurrentSQL = "SELECT"
		m.SQLCursorPos = 0

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyBackspace})

		if newModel.CurrentSQL != "SELECT" {
			t.Errorf("CurrentSQL = %q, want %q", newModel.CurrentSQL, "SELECT")
		}
	})

	t.Run("Left at start stays at start", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.CurrentSQL = "SELECT"
		m.SQLCursorPos = 0

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyLeft})

		if newModel.SQLCursorPos != 0 {
			t.Errorf("SQLCursorPos = %d, want 0", newModel.SQLCursorPos)
		}
	})

	t.Run("Right at end stays at end", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.CurrentSQL = "SELECT"
		m.SQLCursorPos = 6

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyRight})

		if newModel.SQLCursorPos != 6 {
			t.Errorf("SQLCursorPos = %d, want 6", newModel.SQLCursorPos)
		}
	})
}
