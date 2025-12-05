package app

import (
	"strings"

	"github.com/oracle/nosql-go-sdk/nosqldb"

	"github.com/camikura/dito/internal/db"
)

// FocusPane represents which pane is currently focused
type FocusPane int

const (
	FocusPaneConnection FocusPane = iota
	FocusPaneTables
	FocusPaneSchema
	FocusPaneSQL
	FocusPaneData
)

// WindowState holds window dimensions and calculated pane heights
type WindowState struct {
	Width  int
	Height int

	// Pane heights (calculated dynamically)
	ConnectionPaneHeight int
	TablesHeight         int
	SchemaHeight         int
	SQLHeight            int
}

// ConnectionState holds connection-related state
type ConnectionState struct {
	Endpoint    string
	Connected   bool
	Message     string // Connection status message
	NosqlClient *nosqldb.Client
}

// TablesState holds tables pane state
type TablesState struct {
	Tables        []string
	SelectedTable int // Index of selected table (marked with *)
	CursorTable   int // Index of table under cursor
	ScrollOffset  int // Scroll offset for tables pane
}

// SchemaState holds schema pane state
type SchemaState struct {
	TableDetails   map[string]*db.TableDetailsResult
	LoadingDetails bool
	ErrorMsg       string // Error message from schema fetch
	ScrollOffset   int    // Scroll offset for schema pane
}

// SQLState holds SQL pane state
type SQLState struct {
	CurrentSQL            string
	CustomSQL             bool
	ColumnOrder           []string // Column order from custom SQL SELECT clause
	PreviousSelectedTable int      // Saved SelectedTable before custom SQL
	ScrollOffset          int      // Scroll offset for SQL pane
	CursorPos             int      // Cursor position for inline editing
}

// DataState holds data pane state
type DataState struct {
	TableData        map[string]*db.TableDataResult
	LoadingData      bool
	ErrorMsg         string // Error message from data fetch
	SelectedDataRow  int
	ViewportOffset   int
	HorizontalOffset int
}

// ConnectionDialogState holds connection setup dialog state
type ConnectionDialogState struct {
	Visible      bool
	Field        int    // 0: Endpoint, 1: Port, 2: Connect button
	EditEndpoint string // Endpoint being edited
	EditPort     string // Port being edited
	EditCursorPos int   // Cursor position in current field
}

// RecordDetailDialogState holds record detail dialog state
type RecordDetailDialogState struct {
	Visible      bool
	ScrollOffset int
}

// UIState holds temporary UI state (messages, confirmations)
type UIState struct {
	CopyMessage      string // Temporary message shown after copy operation
	QuitConfirmation bool   // Whether quit confirmation is pending
}

// Model represents the application state
type Model struct {
	Window           WindowState
	Connection       ConnectionState
	Tables           TablesState
	Schema           SchemaState
	SQL              SQLState
	Data             DataState
	ConnectionDialog ConnectionDialogState
	RecordDetail     RecordDetailDialogState
	UI               UIState

	// Focus management
	CurrentPane FocusPane
}

// InitialModel creates the initial model for new UI
func InitialModel() Model {
	return Model{
		CurrentPane: FocusPaneConnection,
		Tables: TablesState{
			Tables:        []string{},
			SelectedTable: -1,
			CursorTable:   0,
		},
		Schema: SchemaState{
			TableDetails: make(map[string]*db.TableDetailsResult),
		},
		SQL: SQLState{
			PreviousSelectedTable: -1,
		},
		Data: DataState{
			TableData: make(map[string]*db.TableDataResult),
		},
	}
}

// NextPane moves focus to the next focusable pane
func (m Model) NextPane() Model {
	// Focus order: Connection → Tables → Schema → SQL → Data
	switch m.CurrentPane {
	case FocusPaneConnection:
		m.CurrentPane = FocusPaneTables
	case FocusPaneTables:
		m.CurrentPane = FocusPaneSchema
	case FocusPaneSchema:
		m.CurrentPane = FocusPaneSQL
	case FocusPaneSQL:
		m.CurrentPane = FocusPaneData
	case FocusPaneData:
		m.CurrentPane = FocusPaneConnection
	}
	return m
}

// PrevPane moves focus to the previous focusable pane
func (m Model) PrevPane() Model {
	// Focus order: Connection ← Tables ← Schema ← SQL ← Data
	switch m.CurrentPane {
	case FocusPaneConnection:
		m.CurrentPane = FocusPaneData
	case FocusPaneTables:
		m.CurrentPane = FocusPaneConnection
	case FocusPaneSchema:
		m.CurrentPane = FocusPaneTables
	case FocusPaneSQL:
		m.CurrentPane = FocusPaneSchema
	case FocusPaneData:
		m.CurrentPane = FocusPaneSQL
	}
	return m
}

// FindTableName finds the actual table name from the tables list using case-insensitive matching.
// Returns the matched table name from the list, or empty string if not found.
func (m Model) FindTableName(name string) string {
	if name == "" {
		return ""
	}
	nameLower := strings.ToLower(name)
	for _, t := range m.Tables.Tables {
		if strings.ToLower(t) == nameLower {
			return t
		}
	}
	return ""
}

// FindTableIndex finds the index of a table name in the tables list using case-insensitive matching.
// Returns the index, or -1 if not found.
func (m Model) FindTableIndex(name string) int {
	if name == "" {
		return -1
	}
	nameLower := strings.ToLower(name)
	for i, t := range m.Tables.Tables {
		if strings.ToLower(t) == nameLower {
			return i
		}
	}
	return -1
}

// HasValidSelectedTable returns true if SelectedTable points to a valid table.
func (m Model) HasValidSelectedTable() bool {
	return m.Tables.SelectedTable >= 0 && m.Tables.SelectedTable < len(m.Tables.Tables)
}

// HasValidCursorTable returns true if CursorTable points to a valid table.
func (m Model) HasValidCursorTable() bool {
	return m.Tables.CursorTable >= 0 && m.Tables.CursorTable < len(m.Tables.Tables)
}

// SelectedTableName returns the name of the selected table, or empty string if none.
func (m Model) SelectedTableName() string {
	if !m.HasValidSelectedTable() {
		return ""
	}
	return m.Tables.Tables[m.Tables.SelectedTable]
}

// CursorTableName returns the name of the table under cursor, or empty string if none.
func (m Model) CursorTableName() string {
	if !m.HasValidCursorTable() {
		return ""
	}
	return m.Tables.Tables[m.Tables.CursorTable]
}

// GetTableDetails returns the table details for the given table name, or nil if not found.
func (m Model) GetTableDetails(tableName string) *db.TableDetailsResult {
	if tableName == "" {
		return nil
	}
	details, exists := m.Schema.TableDetails[tableName]
	if !exists || details == nil {
		return nil
	}
	return details
}

// GetTableData returns the table data for the given table name, or nil if not found.
func (m Model) GetTableData(tableName string) *db.TableDataResult {
	if tableName == "" {
		return nil
	}
	data, exists := m.Data.TableData[tableName]
	if !exists || data == nil {
		return nil
	}
	return data
}

// GetSelectedTableData returns the data for the selected table, or nil if none.
func (m Model) GetSelectedTableData() *db.TableDataResult {
	return m.GetTableData(m.SelectedTableName())
}

// GetSelectedTableDetails returns the details for the selected table, or nil if none.
func (m Model) GetSelectedTableDetails() *db.TableDetailsResult {
	return m.GetTableDetails(m.SelectedTableName())
}
