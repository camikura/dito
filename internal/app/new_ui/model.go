package new_ui

import (
	"github.com/oracle/nosql-go-sdk/nosqldb"

	"github.com/camikura/dito/internal/db"
)

// FocusPane represents which pane is currently focused
type FocusPane int

const (
	FocusPaneConnection FocusPane = iota
	FocusPaneTables
	FocusPaneSQL
	FocusPaneData
)

// Model represents the new UI model
type Model struct {
	// Window dimensions
	Width  int
	Height int

	// Focus management
	CurrentPane FocusPane

	// Connection info
	Endpoint      string
	Connected     bool
	ConnectionMsg string

	// Tables
	Tables        []string
	SelectedTable int  // Index of selected table (marked with *)
	CursorTable   int  // Index of table under cursor
	NosqlClient   *nosqldb.Client

	// Schema (display only, auto-updated from cursor position)
	TableDetails   map[string]*db.TableDetailsResult
	LoadingDetails bool

	// SQL
	CurrentSQL string
	CustomSQL  bool

	// Data
	TableData       map[string]*db.TableDataResult
	LoadingData     bool
	SelectedDataRow int
	ViewportOffset  int
	HorizontalOffset int
}

// InitialModel creates the initial model for new UI
func InitialModel() Model {
	return Model{
		CurrentPane:    FocusPaneConnection,
		Connected:      false,
		Tables:         []string{},
		SelectedTable:  -1,
		CursorTable:    0,
		TableDetails:   make(map[string]*db.TableDetailsResult),
		TableData:      make(map[string]*db.TableDataResult),
		CurrentSQL:     "",
		CustomSQL:      false,
	}
}

// NextPane moves focus to the next focusable pane
func (m Model) NextPane() Model {
	// Focus order: Connection → Tables → SQL → Data
	switch m.CurrentPane {
	case FocusPaneConnection:
		m.CurrentPane = FocusPaneTables
	case FocusPaneTables:
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
	// Focus order: Connection ← Tables ← SQL ← Data
	switch m.CurrentPane {
	case FocusPaneConnection:
		m.CurrentPane = FocusPaneData
	case FocusPaneTables:
		m.CurrentPane = FocusPaneConnection
	case FocusPaneSQL:
		m.CurrentPane = FocusPaneTables
	case FocusPaneData:
		m.CurrentPane = FocusPaneSQL
	}
	return m
}
