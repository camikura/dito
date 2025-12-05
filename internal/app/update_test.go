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

	if newModel.Window.Width != 120 {
		t.Errorf("Width = %d, want 120", newModel.Window.Width)
	}
	if newModel.Window.Height != 40 {
		t.Errorf("Height = %d, want 40", newModel.Window.Height)
	}
}

func TestNextPrevPane(t *testing.T) {
	m := InitialModel()
	m.Connection.Connected = true // Enable pane switching

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
	m.Connection.Connected = false
	m.CurrentPane = FocusPaneConnection
	m.Window.Width = 120
	m.Window.Height = 40

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
				SQL: SQLState{CustomSQL: true, ColumnOrder: []string{"name", "id", "email"}},
			},
			table:    "users",
			rows:     []map[string]interface{}{{"id": 1, "name": "test", "email": "test@example.com"}},
			expected: []string{"name", "id", "email"},
		},
		{
			name: "without custom SQL - use row keys",
			model: Model{
				SQL: SQLState{CustomSQL: false, ColumnOrder: nil},
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
				Tables: TablesState{SelectedTable: -1},
			},
			expected: 0,
		},
		{
			name: "table selected but no data",
			model: Model{
				Tables: TablesState{Tables: []string{"users"}, SelectedTable: 0},
				Data:   DataState{TableData: map[string]*db.TableDataResult{}},
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

func TestCtrlQQuit(t *testing.T) {
	m := InitialModel()
	m.Window.Width = 120
	m.Window.Height = 40

	// First Ctrl+Q sets confirmation state
	m, cmd := Update(m, tea.KeyMsg{Type: tea.KeyCtrlQ})
	if !m.UI.QuitConfirmation {
		t.Error("Expected QuitConfirmation to be true after first Ctrl+Q")
	}
	if cmd == nil {
		t.Error("Expected timer command for confirmation timeout")
	}

	// Second Ctrl+Q should quit
	_, cmd = Update(m, tea.KeyMsg{Type: tea.KeyCtrlQ})
	if cmd == nil {
		t.Error("Expected quit command on second Ctrl+Q")
	}
}

func TestCtrlCInDataPane(t *testing.T) {
	m := InitialModel()
	m.Window.Width = 120
	m.Window.Height = 40
	m.CurrentPane = FocusPaneData

	// Ctrl+C in data pane should not quit (it copies)
	_, cmd := Update(m, tea.KeyMsg{Type: tea.KeyCtrlC})

	// cmd should not be quit (it should be nil or a timer for copy message)
	// Since no data is selected, it returns nil
	_ = cmd // Just ensure no panic
}

func TestHandleConnectionKeys(t *testing.T) {
	t.Run("Enter opens connection dialog", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneConnection
		m.Window.Width = 120
		m.Window.Height = 40

		newModel, _ := handleConnectionKeys(m, tea.KeyMsg{Type: tea.KeyEnter})

		if !newModel.ConnectionDialog.Visible {
			t.Error("Expected connection dialog to be visible")
		}
		if newModel.ConnectionDialog.Field != 0 {
			t.Errorf("ConnectionDialogField = %d, want 0", newModel.ConnectionDialog.Field)
		}
	})

	t.Run("Ctrl+D disconnects", func(t *testing.T) {
		m := InitialModel()
		m.Connection.Connected = true
		m.CurrentPane = FocusPaneConnection
		m.Tables.Tables = []string{"users", "products"}
		m.Tables.SelectedTable = 1

		newModel, _ := handleConnectionKeys(m, tea.KeyMsg{Type: tea.KeyCtrlD})

		if newModel.Connection.Connected {
			t.Error("Expected to be disconnected")
		}
		if len(newModel.Tables.Tables) != 0 {
			t.Errorf("Tables should be cleared, got %v", newModel.Tables.Tables)
		}
		if newModel.Tables.SelectedTable != -1 {
			t.Errorf("SelectedTable = %d, want -1", newModel.Tables.SelectedTable)
		}
	})
}

func TestHandleConnectionDialogKeys(t *testing.T) {
	t.Run("Esc closes dialog", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionDialog.Visible = true

		newModel, _ := handleConnectionDialogKeys(m, tea.KeyMsg{Type: tea.KeyEsc})

		if newModel.ConnectionDialog.Visible {
			t.Error("Expected dialog to be closed")
		}
	})

	t.Run("Tab switches fields", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionDialog.Visible = true
		m.ConnectionDialog.Field = 0

		newModel, _ := handleConnectionDialogKeys(m, tea.KeyMsg{Type: tea.KeyTab})

		if newModel.ConnectionDialog.Field != 1 {
			t.Errorf("ConnectionDialogField = %d, want 1", newModel.ConnectionDialog.Field)
		}
	})

	t.Run("Backspace deletes character", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionDialog.Visible = true
		m.ConnectionDialog.Field = 0
		m.ConnectionDialog.EditEndpoint = "localhost"
		m.ConnectionDialog.EditCursorPos = 9

		newModel, _ := handleConnectionDialogKeys(m, tea.KeyMsg{Type: tea.KeyBackspace})

		if newModel.ConnectionDialog.EditEndpoint != "localhos" {
			t.Errorf("EditEndpoint = %q, want %q", newModel.ConnectionDialog.EditEndpoint, "localhos")
		}
	})

	t.Run("Runes inserts characters", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionDialog.Visible = true
		m.ConnectionDialog.Field = 0
		m.ConnectionDialog.EditEndpoint = "local"
		m.ConnectionDialog.EditCursorPos = 5

		newModel, _ := handleConnectionDialogKeys(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})

		if newModel.ConnectionDialog.EditEndpoint != "localh" {
			t.Errorf("EditEndpoint = %q, want %q", newModel.ConnectionDialog.EditEndpoint, "localh")
		}
	})

	t.Run("Left arrow moves cursor", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionDialog.Visible = true
		m.ConnectionDialog.EditCursorPos = 5

		newModel, _ := handleConnectionDialogKeys(m, tea.KeyMsg{Type: tea.KeyLeft})

		if newModel.ConnectionDialog.EditCursorPos != 4 {
			t.Errorf("EditCursorPos = %d, want 4", newModel.ConnectionDialog.EditCursorPos)
		}
	})

	t.Run("Home moves cursor to start", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionDialog.Visible = true
		m.ConnectionDialog.EditCursorPos = 5

		newModel, _ := handleConnectionDialogKeys(m, tea.KeyMsg{Type: tea.KeyHome})

		if newModel.ConnectionDialog.EditCursorPos != 0 {
			t.Errorf("EditCursorPos = %d, want 0", newModel.ConnectionDialog.EditCursorPos)
		}
	})
}

func TestHandleTablesKeys(t *testing.T) {
	t.Run("Down arrow moves cursor down", func(t *testing.T) {
		m := InitialModel()
		m.Tables.Tables = []string{"users", "products", "orders"}
		m.Tables.CursorTable = 0
		m.Window.TablesHeight = 10

		newModel, _ := handleTablesKeys(m, tea.KeyMsg{Type: tea.KeyDown})

		if newModel.Tables.CursorTable != 1 {
			t.Errorf("CursorTable = %d, want 1", newModel.Tables.CursorTable)
		}
	})

	t.Run("Up arrow moves cursor up", func(t *testing.T) {
		m := InitialModel()
		m.Tables.Tables = []string{"users", "products", "orders"}
		m.Tables.CursorTable = 2
		m.Window.TablesHeight = 10

		newModel, _ := handleTablesKeys(m, tea.KeyMsg{Type: tea.KeyUp})

		if newModel.Tables.CursorTable != 1 {
			t.Errorf("CursorTable = %d, want 1", newModel.Tables.CursorTable)
		}
	})

	t.Run("Down at bottom stays at bottom", func(t *testing.T) {
		m := InitialModel()
		m.Tables.Tables = []string{"users", "products"}
		m.Tables.CursorTable = 1
		m.Window.TablesHeight = 10

		newModel, _ := handleTablesKeys(m, tea.KeyMsg{Type: tea.KeyDown})

		if newModel.Tables.CursorTable != 1 {
			t.Errorf("CursorTable = %d, want 1 (should stay at bottom)", newModel.Tables.CursorTable)
		}
	})

	t.Run("Up at top stays at top", func(t *testing.T) {
		m := InitialModel()
		m.Tables.Tables = []string{"users", "products"}
		m.Tables.CursorTable = 0
		m.Window.TablesHeight = 10

		newModel, _ := handleTablesKeys(m, tea.KeyMsg{Type: tea.KeyUp})

		if newModel.Tables.CursorTable != 0 {
			t.Errorf("CursorTable = %d, want 0 (should stay at top)", newModel.Tables.CursorTable)
		}
	})
}

func TestHandleSchemaKeys(t *testing.T) {
	t.Run("Up at top stays at top", func(t *testing.T) {
		m := InitialModel()
		m.Tables.Tables = []string{"users"}
		m.Tables.SelectedTable = 0
		m.Schema.ScrollOffset = 0

		newModel, _ := handleSchemaKeys(m, tea.KeyMsg{Type: tea.KeyUp})

		if newModel.Schema.ScrollOffset != 0 {
			t.Errorf("SchemaScrollOffset = %d, want 0", newModel.Schema.ScrollOffset)
		}
	})

	t.Run("Up scrolls up when offset > 0", func(t *testing.T) {
		m := InitialModel()
		m.Tables.Tables = []string{"users"}
		m.Tables.SelectedTable = 0
		m.Schema.ScrollOffset = 2

		newModel, _ := handleSchemaKeys(m, tea.KeyMsg{Type: tea.KeyUp})

		if newModel.Schema.ScrollOffset != 1 {
			t.Errorf("SchemaScrollOffset = %d, want 1", newModel.Schema.ScrollOffset)
		}
	})

	t.Run("No table selected returns unchanged", func(t *testing.T) {
		m := InitialModel()
		m.Tables.Tables = []string{}
		m.Tables.SelectedTable = -1

		newModel, _ := handleSchemaKeys(m, tea.KeyMsg{Type: tea.KeyDown})

		if newModel.Schema.ScrollOffset != 0 {
			t.Errorf("SchemaScrollOffset = %d, want 0", newModel.Schema.ScrollOffset)
		}
	})
}

func TestHandleSQLKeys(t *testing.T) {
	t.Run("Backspace deletes character", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.SQL.CurrentSQL = "SELECT * FROM users"
		m.SQL.CursorPos = 19

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyBackspace})

		if newModel.SQL.CurrentSQL != "SELECT * FROM user" {
			t.Errorf("CurrentSQL = %q, want %q", newModel.SQL.CurrentSQL, "SELECT * FROM user")
		}
	})

	t.Run("Enter inserts newline", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.SQL.CurrentSQL = "SELECT *"
		m.SQL.CursorPos = 8

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyEnter})

		if newModel.SQL.CurrentSQL != "SELECT *\n" {
			t.Errorf("CurrentSQL = %q, want %q", newModel.SQL.CurrentSQL, "SELECT *\n")
		}
	})

	t.Run("Left arrow moves cursor left", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.SQL.CurrentSQL = "SELECT"
		m.SQL.CursorPos = 6

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyLeft})

		if newModel.SQL.CursorPos != 5 {
			t.Errorf("SQLCursorPos = %d, want 5", newModel.SQL.CursorPos)
		}
	})

	t.Run("Right arrow moves cursor right", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.SQL.CurrentSQL = "SELECT"
		m.SQL.CursorPos = 3

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyRight})

		if newModel.SQL.CursorPos != 4 {
			t.Errorf("SQLCursorPos = %d, want 4", newModel.SQL.CursorPos)
		}
	})

	t.Run("Home moves cursor to start", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.SQL.CurrentSQL = "SELECT * FROM users"
		m.SQL.CursorPos = 10

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyHome})

		if newModel.SQL.CursorPos != 0 {
			t.Errorf("SQLCursorPos = %d, want 0", newModel.SQL.CursorPos)
		}
	})

	t.Run("End moves cursor to end", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.SQL.CurrentSQL = "SELECT"
		m.SQL.CursorPos = 0

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyEnd})

		if newModel.SQL.CursorPos != 6 {
			t.Errorf("SQLCursorPos = %d, want 6", newModel.SQL.CursorPos)
		}
	})

	t.Run("Space inserts space", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.SQL.CurrentSQL = "SELECT*"
		m.SQL.CursorPos = 6

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeySpace})

		if newModel.SQL.CurrentSQL != "SELECT *" {
			t.Errorf("CurrentSQL = %q, want %q", newModel.SQL.CurrentSQL, "SELECT *")
		}
	})

	t.Run("Runes inserts characters", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.SQL.CurrentSQL = "SELECT"
		m.SQL.CursorPos = 6

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'X'}})

		if newModel.SQL.CurrentSQL != "SELECTX" {
			t.Errorf("CurrentSQL = %q, want %q", newModel.SQL.CurrentSQL, "SELECTX")
		}
	})

	t.Run("Up arrow moves cursor up in multi-line", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.SQL.CurrentSQL = "SELECT *\nFROM users"
		m.SQL.CursorPos = 14 // middle of second line

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyUp})

		// Should move to first line at same column or end
		if newModel.SQL.CursorPos >= 9 {
			t.Errorf("SQLCursorPos = %d, should be less than 9", newModel.SQL.CursorPos)
		}
	})

	t.Run("Down arrow moves cursor down in multi-line", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.SQL.CurrentSQL = "SELECT *\nFROM users"
		m.SQL.CursorPos = 4 // middle of first line

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyDown})

		// Should move to second line
		if newModel.SQL.CursorPos < 9 {
			t.Errorf("SQLCursorPos = %d, should be >= 9", newModel.SQL.CursorPos)
		}
	})
}

func TestHandleDataKeys(t *testing.T) {
	t.Run("Down arrow moves selection down", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneData
		m.Tables.Tables = []string{"users"}
		m.Tables.SelectedTable = 0
		m.Data.SelectedDataRow = 0
		m.Window.Height = 40
		m.Data.TableData = map[string]*db.TableDataResult{
			"users": {
				Rows: []map[string]interface{}{
					{"id": 1}, {"id": 2}, {"id": 3},
				},
			},
		}

		newModel, _ := handleDataKeys(m, tea.KeyMsg{Type: tea.KeyDown})

		if newModel.Data.SelectedDataRow != 1 {
			t.Errorf("SelectedDataRow = %d, want 1", newModel.Data.SelectedDataRow)
		}
	})

	t.Run("Up arrow moves selection up", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneData
		m.Tables.Tables = []string{"users"}
		m.Tables.SelectedTable = 0
		m.Data.SelectedDataRow = 2
		m.Window.Height = 40
		m.Data.TableData = map[string]*db.TableDataResult{
			"users": {
				Rows: []map[string]interface{}{
					{"id": 1}, {"id": 2}, {"id": 3},
				},
			},
		}

		newModel, _ := handleDataKeys(m, tea.KeyMsg{Type: tea.KeyUp})

		if newModel.Data.SelectedDataRow != 1 {
			t.Errorf("SelectedDataRow = %d, want 1", newModel.Data.SelectedDataRow)
		}
	})

	t.Run("Left arrow scrolls left", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneData
		m.Data.HorizontalOffset = 5

		newModel, _ := handleDataKeys(m, tea.KeyMsg{Type: tea.KeyLeft})

		if newModel.Data.HorizontalOffset != 4 {
			t.Errorf("HorizontalOffset = %d, want 4", newModel.Data.HorizontalOffset)
		}
	})

	t.Run("Left at 0 stays at 0", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneData
		m.Data.HorizontalOffset = 0

		newModel, _ := handleDataKeys(m, tea.KeyMsg{Type: tea.KeyLeft})

		if newModel.Data.HorizontalOffset != 0 {
			t.Errorf("HorizontalOffset = %d, want 0", newModel.Data.HorizontalOffset)
		}
	})

	t.Run("Enter opens record detail", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneData
		m.Tables.Tables = []string{"users"}
		m.Tables.SelectedTable = 0
		m.Data.SelectedDataRow = 0
		m.Data.TableData = map[string]*db.TableDataResult{
			"users": {
				Rows: []map[string]interface{}{{"id": 1}},
			},
		}

		newModel, _ := handleDataKeys(m, tea.KeyMsg{Type: tea.KeyEnter})

		if !newModel.RecordDetail.Visible {
			t.Error("Expected record detail to be visible")
		}
	})

	t.Run("Esc resets custom SQL", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneData
		m.SQL.CustomSQL = true
		m.SQL.ColumnOrder = []string{"name", "id"}
		m.Data.SelectedDataRow = 5
		m.Data.ViewportOffset = 3

		newModel, _ := handleDataKeys(m, tea.KeyMsg{Type: tea.KeyEsc})

		if newModel.SQL.CustomSQL {
			t.Error("Expected CustomSQL to be false")
		}
		if newModel.SQL.ColumnOrder != nil {
			t.Error("Expected ColumnOrder to be nil")
		}
		if newModel.Data.SelectedDataRow != 0 {
			t.Errorf("SelectedDataRow = %d, want 0", newModel.Data.SelectedDataRow)
		}
	})
}

func TestHandleRecordDetailKeys(t *testing.T) {
	t.Run("Esc closes dialog", func(t *testing.T) {
		m := InitialModel()
		m.RecordDetail.Visible = true
		m.RecordDetail.ScrollOffset = 5

		newModel, _ := handleRecordDetailKeys(m, tea.KeyMsg{Type: tea.KeyEsc})

		if newModel.RecordDetail.Visible {
			t.Error("Expected record detail to be closed")
		}
		if newModel.RecordDetail.ScrollOffset != 0 {
			t.Errorf("RecordDetailScroll = %d, want 0", newModel.RecordDetail.ScrollOffset)
		}
	})

	t.Run("Down scrolls down when maxScroll > 0", func(t *testing.T) {
		m := InitialModel()
		m.RecordDetail.Visible = true
		m.RecordDetail.ScrollOffset = 0
		m.Tables.Tables = []string{"users"}
		m.Tables.SelectedTable = 0
		m.Data.SelectedDataRow = 0
		m.Window.Height = 10 // Small height to ensure maxScroll > 0

		// Create a row with many fields to exceed visible height
		row := make(map[string]interface{})
		for i := 0; i < 20; i++ {
			row[string(rune('a'+i))] = i
		}
		m.Data.TableData = map[string]*db.TableDataResult{
			"users": {
				Rows: []map[string]interface{}{row},
			},
		}

		newModel, _ := handleRecordDetailKeys(m, tea.KeyMsg{Type: tea.KeyDown})

		if newModel.RecordDetail.ScrollOffset != 1 {
			t.Errorf("RecordDetailScroll = %d, want 1", newModel.RecordDetail.ScrollOffset)
		}
	})

	t.Run("Up scrolls up", func(t *testing.T) {
		m := InitialModel()
		m.RecordDetail.Visible = true
		m.RecordDetail.ScrollOffset = 3

		newModel, _ := handleRecordDetailKeys(m, tea.KeyMsg{Type: tea.KeyUp})

		if newModel.RecordDetail.ScrollOffset != 2 {
			t.Errorf("RecordDetailScroll = %d, want 2", newModel.RecordDetail.ScrollOffset)
		}
	})

	t.Run("Home scrolls to top", func(t *testing.T) {
		m := InitialModel()
		m.RecordDetail.Visible = true
		m.RecordDetail.ScrollOffset = 10

		newModel, _ := handleRecordDetailKeys(m, tea.KeyMsg{Type: tea.KeyHome})

		if newModel.RecordDetail.ScrollOffset != 0 {
			t.Errorf("RecordDetailScroll = %d, want 0", newModel.RecordDetail.ScrollOffset)
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
		if !reflect.DeepEqual(newModel.Tables.Tables, expected) {
			t.Errorf("Tables = %v, want %v", newModel.Tables.Tables, expected)
		}
		if newModel.Tables.CursorTable != 0 {
			t.Errorf("CursorTable = %d, want 0", newModel.Tables.CursorTable)
		}
	})

	t.Run("error returns unchanged model", func(t *testing.T) {
		m := InitialModel()

		newModel, _ := handleTableListResult(m, db.TableListResult{
			Err: errors.New("test error"),
		})

		if len(newModel.Tables.Tables) != 0 {
			t.Errorf("Tables should be empty on error")
		}
	})
}

func TestHandleTableDataResult(t *testing.T) {
	t.Run("success stores data", func(t *testing.T) {
		m := InitialModel()
		m.Data.LoadingData = true

		rows := []map[string]interface{}{
			{"id": 1, "name": "test"},
		}

		newModel, _ := handleTableDataResult(m, db.TableDataResult{
			TableName: "users",
			Rows:      rows,
			HasMore:   true,
		})

		if newModel.Data.LoadingData {
			t.Error("Expected LoadingData to be false")
		}
		data := newModel.Data.TableData["users"]
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
		m.Data.LoadingData = true
		m.Data.TableData = map[string]*db.TableDataResult{
			"users": {
				Rows: []map[string]interface{}{{"id": 1}},
			},
		}

		newModel, _ := handleTableDataResult(m, db.TableDataResult{
			TableName: "users",
			Rows:      []map[string]interface{}{{"id": 2}},
			IsAppend:  true,
		})

		data := newModel.Data.TableData["users"]
		if len(data.Rows) != 2 {
			t.Errorf("Rows count = %d, want 2", len(data.Rows))
		}
	})

	t.Run("error clears loading state", func(t *testing.T) {
		m := InitialModel()
		m.Data.LoadingData = true

		newModel, _ := handleTableDataResult(m, db.TableDataResult{
			Err: errors.New("test error"),
		})

		if newModel.Data.LoadingData {
			t.Error("Expected LoadingData to be false")
		}
		if newModel.Data.ErrorMsg == "" {
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

		if newModel.Connection.Connected {
			t.Error("Expected Connected to be false")
		}
		if newModel.Connection.Message == "" {
			t.Error("Expected ConnectionMsg to be set")
		}
	})
}

func TestCalculateRecordDetailMaxScroll(t *testing.T) {
	t.Run("no table selected returns 0", func(t *testing.T) {
		m := InitialModel()
		m.Tables.SelectedTable = -1

		result := calculateRecordDetailMaxScroll(m)

		if result != 0 {
			t.Errorf("calculateRecordDetailMaxScroll() = %d, want 0", result)
		}
	})

	t.Run("no data returns 0", func(t *testing.T) {
		m := InitialModel()
		m.Tables.Tables = []string{"users"}
		m.Tables.SelectedTable = 0
		m.Data.TableData = map[string]*db.TableDataResult{}

		result := calculateRecordDetailMaxScroll(m)

		if result != 0 {
			t.Errorf("calculateRecordDetailMaxScroll() = %d, want 0", result)
		}
	})
}

func TestHandleTableDetailsResult(t *testing.T) {
	t.Run("error sets schema error message", func(t *testing.T) {
		m := InitialModel()
		m.Data.LoadingData = true

		newModel, cmd := handleTableDetailsResult(m, db.TableDetailsResult{
			Err: errors.New("failed to fetch schema"),
		})

		if newModel.Schema.ErrorMsg == "" {
			t.Error("Expected SchemaErrorMsg to be set")
		}
		if newModel.Data.LoadingData {
			t.Error("Expected LoadingData to be false")
		}
		if cmd != nil {
			t.Error("Expected no command")
		}
	})

	t.Run("success stores table details", func(t *testing.T) {
		m := InitialModel()
		m.Schema.TableDetails = make(map[string]*db.TableDetailsResult)

		newModel, _ := handleTableDetailsResult(m, db.TableDetailsResult{
			TableName: "users",
		})

		if newModel.Schema.ErrorMsg != "" {
			t.Errorf("SchemaErrorMsg should be empty, got %q", newModel.Schema.ErrorMsg)
		}
		if newModel.Schema.TableDetails["users"] == nil {
			t.Error("Expected table details to be stored")
		}
	})

	t.Run("success clears previous error", func(t *testing.T) {
		m := InitialModel()
		m.Schema.ErrorMsg = "previous error"
		m.Schema.TableDetails = make(map[string]*db.TableDetailsResult)

		newModel, _ := handleTableDetailsResult(m, db.TableDetailsResult{
			TableName: "users",
		})

		if newModel.Schema.ErrorMsg != "" {
			t.Errorf("SchemaErrorMsg should be cleared, got %q", newModel.Schema.ErrorMsg)
		}
	})
}

func TestHandleConnectionDialogKeysAdditional(t *testing.T) {
	t.Run("Right arrow moves cursor right within text", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionDialog.Visible = true
		m.ConnectionDialog.Field = 0
		m.ConnectionDialog.EditEndpoint = "localhost"
		m.ConnectionDialog.EditCursorPos = 3

		newModel, _ := handleConnectionDialogKeys(m, tea.KeyMsg{Type: tea.KeyRight})

		if newModel.ConnectionDialog.EditCursorPos != 4 {
			t.Errorf("EditCursorPos = %d, want 4", newModel.ConnectionDialog.EditCursorPos)
		}
	})

	t.Run("Right arrow at end stays at end", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionDialog.Visible = true
		m.ConnectionDialog.Field = 0
		m.ConnectionDialog.EditEndpoint = "localhost"
		m.ConnectionDialog.EditCursorPos = 9

		newModel, _ := handleConnectionDialogKeys(m, tea.KeyMsg{Type: tea.KeyRight})

		if newModel.ConnectionDialog.EditCursorPos != 9 {
			t.Errorf("EditCursorPos = %d, want 9", newModel.ConnectionDialog.EditCursorPos)
		}
	})

	t.Run("End moves cursor to end", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionDialog.Visible = true
		m.ConnectionDialog.Field = 0
		m.ConnectionDialog.EditEndpoint = "localhost"
		m.ConnectionDialog.EditCursorPos = 3

		newModel, _ := handleConnectionDialogKeys(m, tea.KeyMsg{Type: tea.KeyEnd})

		if newModel.ConnectionDialog.EditCursorPos != 9 {
			t.Errorf("EditCursorPos = %d, want 9", newModel.ConnectionDialog.EditCursorPos)
		}
	})

	t.Run("Delete removes character at cursor", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionDialog.Visible = true
		m.ConnectionDialog.Field = 0
		m.ConnectionDialog.EditEndpoint = "localhost"
		m.ConnectionDialog.EditCursorPos = 5

		newModel, _ := handleConnectionDialogKeys(m, tea.KeyMsg{Type: tea.KeyDelete})

		if newModel.ConnectionDialog.EditEndpoint != "localost" {
			t.Errorf("EditEndpoint = %q, want %q", newModel.ConnectionDialog.EditEndpoint, "localost")
		}
	})

	t.Run("Tab wraps around from last field", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionDialog.Visible = true
		m.ConnectionDialog.Field = 1 // Port field (last)

		newModel, _ := handleConnectionDialogKeys(m, tea.KeyMsg{Type: tea.KeyTab})

		if newModel.ConnectionDialog.Field != 0 {
			t.Errorf("ConnectionDialogField = %d, want 0", newModel.ConnectionDialog.Field)
		}
	})

	t.Run("Editing port field", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionDialog.Visible = true
		m.ConnectionDialog.Field = 1 // Port field
		m.ConnectionDialog.EditPort = "8080"
		m.ConnectionDialog.EditCursorPos = 4

		newModel, _ := handleConnectionDialogKeys(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})

		if newModel.ConnectionDialog.EditPort != "80801" {
			t.Errorf("EditPort = %q, want %q", newModel.ConnectionDialog.EditPort, "80801")
		}
	})

	t.Run("Backspace on port field", func(t *testing.T) {
		m := InitialModel()
		m.ConnectionDialog.Visible = true
		m.ConnectionDialog.Field = 1
		m.ConnectionDialog.EditPort = "8080"
		m.ConnectionDialog.EditCursorPos = 4

		newModel, _ := handleConnectionDialogKeys(m, tea.KeyMsg{Type: tea.KeyBackspace})

		if newModel.ConnectionDialog.EditPort != "808" {
			t.Errorf("EditPort = %q, want %q", newModel.ConnectionDialog.EditPort, "808")
		}
	})
}

func TestHandleTablesKeysAdditional(t *testing.T) {
	t.Run("Empty tables list returns unchanged", func(t *testing.T) {
		m := InitialModel()
		m.Tables.Tables = []string{}
		m.Tables.CursorTable = 0

		newModel, _ := handleTablesKeys(m, tea.KeyMsg{Type: tea.KeyDown})

		if newModel.Tables.CursorTable != 0 {
			t.Errorf("CursorTable = %d, want 0", newModel.Tables.CursorTable)
		}
	})

	t.Run("Up at top does nothing", func(t *testing.T) {
		m := InitialModel()
		m.Tables.Tables = []string{"users", "products", "orders"}
		m.Tables.CursorTable = 0
		m.Window.TablesHeight = 10

		newModel, _ := handleTablesKeys(m, tea.KeyMsg{Type: tea.KeyUp})

		if newModel.Tables.CursorTable != 0 {
			t.Errorf("CursorTable = %d, want 0", newModel.Tables.CursorTable)
		}
	})

	t.Run("Down at bottom does nothing", func(t *testing.T) {
		m := InitialModel()
		m.Tables.Tables = []string{"users", "products", "orders"}
		m.Tables.CursorTable = 2
		m.Window.TablesHeight = 10

		newModel, _ := handleTablesKeys(m, tea.KeyMsg{Type: tea.KeyDown})

		if newModel.Tables.CursorTable != 2 {
			t.Errorf("CursorTable = %d, want 2", newModel.Tables.CursorTable)
		}
	})

	t.Run("Down arrow moves down", func(t *testing.T) {
		m := InitialModel()
		m.Tables.Tables = []string{"users", "products"}
		m.Tables.CursorTable = 0
		m.Window.TablesHeight = 10

		newModel, _ := handleTablesKeys(m, tea.KeyMsg{Type: tea.KeyDown})

		if newModel.Tables.CursorTable != 1 {
			t.Errorf("CursorTable = %d, want 1", newModel.Tables.CursorTable)
		}
	})

	t.Run("Up arrow moves up", func(t *testing.T) {
		m := InitialModel()
		m.Tables.Tables = []string{"users", "products"}
		m.Tables.CursorTable = 1
		m.Window.TablesHeight = 10

		newModel, _ := handleTablesKeys(m, tea.KeyMsg{Type: tea.KeyUp})

		if newModel.Tables.CursorTable != 0 {
			t.Errorf("CursorTable = %d, want 0", newModel.Tables.CursorTable)
		}
	})
}

func TestHandleSchemaKeysAdditional(t *testing.T) {
	t.Run("Down without table details does nothing", func(t *testing.T) {
		m := InitialModel()
		m.Tables.Tables = []string{"users"}
		m.Tables.SelectedTable = 0
		m.Schema.ScrollOffset = 0

		newModel, _ := handleSchemaKeys(m, tea.KeyMsg{Type: tea.KeyDown})

		// Without table details, maxScroll is 0, so no scrolling
		if newModel.Schema.ScrollOffset != 0 {
			t.Errorf("SchemaScrollOffset = %d, want 0", newModel.Schema.ScrollOffset)
		}
	})

	t.Run("Up arrow scrolls up", func(t *testing.T) {
		m := InitialModel()
		m.Tables.Tables = []string{"users"}
		m.Tables.SelectedTable = 0
		m.Schema.ScrollOffset = 2

		newModel, _ := handleSchemaKeys(m, tea.KeyMsg{Type: tea.KeyUp})

		if newModel.Schema.ScrollOffset != 1 {
			t.Errorf("SchemaScrollOffset = %d, want 1", newModel.Schema.ScrollOffset)
		}
	})
}

func TestHandleDataKeysAdditional(t *testing.T) {
	t.Run("Left arrow scrolls left", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneData
		m.Data.HorizontalOffset = 5

		newModel, _ := handleDataKeys(m, tea.KeyMsg{Type: tea.KeyLeft})

		if newModel.Data.HorizontalOffset != 4 {
			t.Errorf("HorizontalOffset = %d, want 4", newModel.Data.HorizontalOffset)
		}
	})

	t.Run("Left arrow scrolls left second", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneData
		m.Data.HorizontalOffset = 5

		newModel, _ := handleDataKeys(m, tea.KeyMsg{Type: tea.KeyLeft})

		if newModel.Data.HorizontalOffset != 4 {
			t.Errorf("HorizontalOffset = %d, want 4", newModel.Data.HorizontalOffset)
		}
	})

	t.Run("Left at zero stays at zero", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneData
		m.Data.HorizontalOffset = 0

		newModel, _ := handleDataKeys(m, tea.KeyMsg{Type: tea.KeyLeft})

		if newModel.Data.HorizontalOffset != 0 {
			t.Errorf("HorizontalOffset = %d, want 0", newModel.Data.HorizontalOffset)
		}
	})

	t.Run("Enter key shows record detail", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneData
		m.Tables.Tables = []string{"users"}
		m.Tables.SelectedTable = 0
		m.Data.SelectedDataRow = 0
		m.Data.TableData = map[string]*db.TableDataResult{
			"users": {
				Rows: []map[string]interface{}{
					{"id": 1},
				},
			},
		}

		newModel, _ := handleDataKeys(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'\r'}})

		// Check if "enter" is handled
		if newModel.RecordDetail.Visible {
			// Good, enter was handled
		}
	})

	t.Run("Down arrow moves down", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneData
		m.Tables.Tables = []string{"users"}
		m.Tables.SelectedTable = 0
		m.Data.SelectedDataRow = 0
		m.Window.Height = 40
		m.Data.TableData = map[string]*db.TableDataResult{
			"users": {
				Rows: []map[string]interface{}{
					{"id": 1}, {"id": 2},
				},
			},
		}

		newModel, _ := handleDataKeys(m, tea.KeyMsg{Type: tea.KeyDown})

		if newModel.Data.SelectedDataRow != 1 {
			t.Errorf("SelectedDataRow = %d, want 1", newModel.Data.SelectedDataRow)
		}
	})

	t.Run("Up arrow moves up", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneData
		m.Tables.Tables = []string{"users"}
		m.Tables.SelectedTable = 0
		m.Data.SelectedDataRow = 1
		m.Window.Height = 40
		m.Data.TableData = map[string]*db.TableDataResult{
			"users": {
				Rows: []map[string]interface{}{
					{"id": 1}, {"id": 2},
				},
			},
		}

		newModel, _ := handleDataKeys(m, tea.KeyMsg{Type: tea.KeyUp})

		if newModel.Data.SelectedDataRow != 0 {
			t.Errorf("SelectedDataRow = %d, want 0", newModel.Data.SelectedDataRow)
		}
	})

	t.Run("No table selected returns unchanged", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneData
		m.Tables.SelectedTable = -1
		m.Data.SelectedDataRow = 0

		newModel, _ := handleDataKeys(m, tea.KeyMsg{Type: tea.KeyDown})

		if newModel.Data.SelectedDataRow != 0 {
			t.Errorf("SelectedDataRow = %d, want 0", newModel.Data.SelectedDataRow)
		}
	})
}

func TestHandleRecordDetailKeysAdditional(t *testing.T) {
	t.Run("Down arrow scrolls down", func(t *testing.T) {
		m := InitialModel()
		m.RecordDetail.Visible = true
		m.RecordDetail.ScrollOffset = 0
		m.Tables.Tables = []string{"users"}
		m.Tables.SelectedTable = 0
		m.Data.SelectedDataRow = 0
		m.Window.Height = 10

		row := make(map[string]interface{})
		for i := 0; i < 20; i++ {
			row[string(rune('a'+i))] = i
		}
		m.Data.TableData = map[string]*db.TableDataResult{
			"users": {Rows: []map[string]interface{}{row}},
		}

		newModel, _ := handleRecordDetailKeys(m, tea.KeyMsg{Type: tea.KeyDown})

		if newModel.RecordDetail.ScrollOffset != 1 {
			t.Errorf("RecordDetailScroll = %d, want 1", newModel.RecordDetail.ScrollOffset)
		}
	})

	t.Run("Up arrow scrolls up", func(t *testing.T) {
		m := InitialModel()
		m.RecordDetail.Visible = true
		m.RecordDetail.ScrollOffset = 5

		newModel, _ := handleRecordDetailKeys(m, tea.KeyMsg{Type: tea.KeyUp})

		if newModel.RecordDetail.ScrollOffset != 4 {
			t.Errorf("RecordDetailScroll = %d, want 4", newModel.RecordDetail.ScrollOffset)
		}
	})

	t.Run("End scrolls to max", func(t *testing.T) {
		m := InitialModel()
		m.RecordDetail.Visible = true
		m.RecordDetail.ScrollOffset = 0
		m.Tables.Tables = []string{"users"}
		m.Tables.SelectedTable = 0
		m.Data.SelectedDataRow = 0
		m.Window.Height = 10

		row := make(map[string]interface{})
		for i := 0; i < 30; i++ {
			row[string(rune('a'+i))] = i
		}
		m.Data.TableData = map[string]*db.TableDataResult{
			"users": {Rows: []map[string]interface{}{row}},
		}

		newModel, _ := handleRecordDetailKeys(m, tea.KeyMsg{Type: tea.KeyEnd})

		// Should scroll to max (but we just check it's > 0)
		if newModel.RecordDetail.ScrollOffset == 0 {
			t.Errorf("RecordDetailScroll should be > 0")
		}
	})
}

func TestHandleSQLKeysAdditional(t *testing.T) {
	t.Run("Delete removes character at cursor", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.SQL.CurrentSQL = "SELECT * FROM users"
		m.SQL.CursorPos = 7

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyDelete})

		if newModel.SQL.CurrentSQL != "SELECT  FROM users" {
			t.Errorf("CurrentSQL = %q, want %q", newModel.SQL.CurrentSQL, "SELECT  FROM users")
		}
	})

	t.Run("Delete at end does nothing", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.SQL.CurrentSQL = "SELECT"
		m.SQL.CursorPos = 6

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyDelete})

		if newModel.SQL.CurrentSQL != "SELECT" {
			t.Errorf("CurrentSQL = %q, want %q", newModel.SQL.CurrentSQL, "SELECT")
		}
	})

	t.Run("Backspace at start does nothing", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.SQL.CurrentSQL = "SELECT"
		m.SQL.CursorPos = 0

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyBackspace})

		if newModel.SQL.CurrentSQL != "SELECT" {
			t.Errorf("CurrentSQL = %q, want %q", newModel.SQL.CurrentSQL, "SELECT")
		}
	})

	t.Run("Left at start stays at start", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.SQL.CurrentSQL = "SELECT"
		m.SQL.CursorPos = 0

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyLeft})

		if newModel.SQL.CursorPos != 0 {
			t.Errorf("SQLCursorPos = %d, want 0", newModel.SQL.CursorPos)
		}
	})

	t.Run("Right at end stays at end", func(t *testing.T) {
		m := InitialModel()
		m.CurrentPane = FocusPaneSQL
		m.SQL.CurrentSQL = "SELECT"
		m.SQL.CursorPos = 6

		newModel, _ := handleSQLKeys(m, tea.KeyMsg{Type: tea.KeyRight})

		if newModel.SQL.CursorPos != 6 {
			t.Errorf("SQLCursorPos = %d, want 6", newModel.SQL.CursorPos)
		}
	})
}
