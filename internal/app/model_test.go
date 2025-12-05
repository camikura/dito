package app

import (
	"testing"
)

func TestInitialModel(t *testing.T) {
	model := InitialModel()

	// Check initial pane
	if model.CurrentPane != FocusPaneConnection {
		t.Errorf("InitialModel() CurrentPane = %v, want %v", model.CurrentPane, FocusPaneConnection)
	}

	// Check not connected
	if model.Connection.Connected {
		t.Error("InitialModel() Connection.Connected should be false")
	}

	// Check tables is empty slice
	if model.Tables.Tables == nil {
		t.Error("InitialModel() Tables.Tables should not be nil")
	}
	if len(model.Tables.Tables) != 0 {
		t.Errorf("InitialModel() Tables.Tables length = %d, want 0", len(model.Tables.Tables))
	}

	// Check selection indices
	if model.Tables.SelectedTable != -1 {
		t.Errorf("InitialModel() Tables.SelectedTable = %d, want -1", model.Tables.SelectedTable)
	}
	if model.Tables.CursorTable != 0 {
		t.Errorf("InitialModel() Tables.CursorTable = %d, want 0", model.Tables.CursorTable)
	}
	if model.SQL.PreviousSelectedTable != -1 {
		t.Errorf("InitialModel() SQL.PreviousSelectedTable = %d, want -1", model.SQL.PreviousSelectedTable)
	}

	// Check maps are initialized
	if model.Schema.TableDetails == nil {
		t.Error("InitialModel() Schema.TableDetails map should be initialized")
	}
	if model.Data.TableData == nil {
		t.Error("InitialModel() Data.TableData map should be initialized")
	}

	// Check SQL state
	if model.SQL.CurrentSQL != "" {
		t.Errorf("InitialModel() SQL.CurrentSQL = %q, want empty", model.SQL.CurrentSQL)
	}
	if model.SQL.CustomSQL {
		t.Error("InitialModel() SQL.CustomSQL should be false")
	}
}

func TestNextPane(t *testing.T) {
	tests := []struct {
		name     string
		current  FocusPane
		expected FocusPane
	}{
		{"Connection to Tables", FocusPaneConnection, FocusPaneTables},
		{"Tables to Schema", FocusPaneTables, FocusPaneSchema},
		{"Schema to SQL", FocusPaneSchema, FocusPaneSQL},
		{"SQL to Data", FocusPaneSQL, FocusPaneData},
		{"Data wraps to Connection", FocusPaneData, FocusPaneConnection},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Model{CurrentPane: tt.current}
			m = m.NextPane()
			if m.CurrentPane != tt.expected {
				t.Errorf("NextPane() from %v = %v, want %v", tt.current, m.CurrentPane, tt.expected)
			}
		})
	}
}

func TestPrevPane(t *testing.T) {
	tests := []struct {
		name     string
		current  FocusPane
		expected FocusPane
	}{
		{"Connection wraps to Data", FocusPaneConnection, FocusPaneData},
		{"Tables to Connection", FocusPaneTables, FocusPaneConnection},
		{"Schema to Tables", FocusPaneSchema, FocusPaneTables},
		{"SQL to Schema", FocusPaneSQL, FocusPaneSchema},
		{"Data to SQL", FocusPaneData, FocusPaneSQL},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Model{CurrentPane: tt.current}
			m = m.PrevPane()
			if m.CurrentPane != tt.expected {
				t.Errorf("PrevPane() from %v = %v, want %v", tt.current, m.CurrentPane, tt.expected)
			}
		})
	}
}

func TestFindTableName(t *testing.T) {
	m := Model{
		Tables: TablesState{Tables: []string{"Users", "Products", "Orders", "users.addresses"}},
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"exact match", "Users", "Users"},
		{"case insensitive lowercase", "users", "Users"},
		{"case insensitive uppercase", "PRODUCTS", "Products"},
		{"case insensitive mixed", "oRdErS", "Orders"},
		{"child table exact", "users.addresses", "users.addresses"},
		{"child table case insensitive", "Users.Addresses", "users.addresses"},
		{"not found", "customers", ""},
		{"empty input", "", ""},
		{"partial match not found", "User", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := m.FindTableName(tt.input)
			if result != tt.expected {
				t.Errorf("FindTableName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFindTableName_EmptyTables(t *testing.T) {
	m := Model{Tables: TablesState{Tables: []string{}}}

	result := m.FindTableName("users")
	if result != "" {
		t.Errorf("FindTableName on empty tables = %q, want empty", result)
	}
}

func TestFindTableIndex(t *testing.T) {
	m := Model{
		Tables: TablesState{Tables: []string{"Users", "Products", "Orders", "users.addresses"}},
	}

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"exact match first", "Users", 0},
		{"exact match second", "Products", 1},
		{"exact match third", "Orders", 2},
		{"case insensitive lowercase", "users", 0},
		{"case insensitive uppercase", "PRODUCTS", 1},
		{"child table", "users.addresses", 3},
		{"child table case insensitive", "Users.Addresses", 3},
		{"not found", "customers", -1},
		{"empty input", "", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := m.FindTableIndex(tt.input)
			if result != tt.expected {
				t.Errorf("FindTableIndex(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFindTableIndex_EmptyTables(t *testing.T) {
	m := Model{Tables: TablesState{Tables: []string{}}}

	result := m.FindTableIndex("users")
	if result != -1 {
		t.Errorf("FindTableIndex on empty tables = %d, want -1", result)
	}
}

func TestPaneCycle(t *testing.T) {
	// Test that cycling through all panes returns to start
	m := Model{CurrentPane: FocusPaneConnection}

	// Forward cycle
	for i := 0; i < 5; i++ {
		m = m.NextPane()
	}
	if m.CurrentPane != FocusPaneConnection {
		t.Errorf("After 5 NextPane(), CurrentPane = %v, want %v", m.CurrentPane, FocusPaneConnection)
	}

	// Backward cycle
	for i := 0; i < 5; i++ {
		m = m.PrevPane()
	}
	if m.CurrentPane != FocusPaneConnection {
		t.Errorf("After 5 PrevPane(), CurrentPane = %v, want %v", m.CurrentPane, FocusPaneConnection)
	}
}
