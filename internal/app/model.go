package app

import (
	"github.com/oracle/nosql-go-sdk/nosqldb"

	"github.com/camikura/dito/internal/db"
)

// Screen represents the current screen type
type Screen int

const (
	ScreenSelection Screen = iota
	ScreenOnPremiseConfig
	ScreenCloudConfig
	ScreenTableList
)

// ConnectionStatus represents the connection state
type ConnectionStatus int

const (
	StatusDisconnected ConnectionStatus = iota
	StatusConnecting
	StatusConnected
	StatusError
)

// RightPaneMode represents the right pane display mode
type RightPaneMode int

const (
	RightPaneModeSchema RightPaneMode = iota
	RightPaneModeList   // Data list view
	RightPaneModeDetail // Record detail view
)

// OnPremiseConfig holds on-premise connection configuration
type OnPremiseConfig struct {
	Endpoint      string
	Port          string
	Secure        bool
	Focus         int // Currently focused field
	Status        ConnectionStatus
	ErrorMsg      string
	ServerVersion string
	CursorPos     int // Text input cursor position
}

// CloudConfig holds cloud connection configuration
type CloudConfig struct {
	Region        string
	Compartment   string
	AuthMethod    int // 0: OCI Config Profile, 1: Instance Principal, 2: Resource Principal
	ConfigFile    string
	Focus         int // Currently focused field
	Status        ConnectionStatus
	ErrorMsg      string
	ServerVersion string
	CursorPos     int // Text input cursor position
}

// Model is the main application model
type Model struct {
	Screen          Screen
	Choices         []string
	Cursor          int
	Selected        map[int]struct{}
	OnPremiseConfig OnPremiseConfig
	CloudConfig     CloudConfig
	Width           int
	Height          int
	// Table list screen
	NosqlClient    *nosqldb.Client
	Tables         []string
	SelectedTable  int
	Endpoint       string // Connection endpoint (for status display)
	TableDetails   map[string]*db.TableDetailsResult
	LoadingDetails bool
	// Data display
	RightPaneMode    RightPaneMode
	TableData        map[string]*db.TableDataResult
	DataOffset       int // Data fetch offset (for infinite scroll)
	FetchSize        int // Number of rows to fetch at once
	LoadingData      bool
	SelectedDataRow  int // Selected row in data view mode (absolute position)
	ViewportOffset   int // Display start position
	ViewportSize     int // Number of rows to display at once
	HorizontalOffset int // Horizontal scroll offset (column-based, 0-indexed)
}

// InitialModel returns the initial application model
func InitialModel() Model {
	return Model{
		Screen:   ScreenSelection,
		Choices:  []string{"Oracle NoSQL Cloud Service", "On-Premise"},
		Selected: make(map[int]struct{}),
		OnPremiseConfig: OnPremiseConfig{
			Endpoint:  "localhost",
			Port:      "8080",
			Secure:    false,
			Focus:     0,
			Status:    StatusDisconnected,
			CursorPos: 9, // End of "localhost"
		},
		CloudConfig: CloudConfig{
			Region:      "us-ashburn-1",
			Compartment: "",
			AuthMethod:  0, // OCI Config Profile
			ConfigFile:  "DEFAULT",
			Focus:       0,
			Status:      StatusDisconnected,
			CursorPos:   12, // End of "us-ashburn-1"
		},
		TableDetails:  make(map[string]*db.TableDetailsResult),
		RightPaneMode: RightPaneModeSchema,
		TableData:     make(map[string]*db.TableDataResult),
		DataOffset:    0,
		FetchSize:     100, // Fetch 100 rows at once (infinite scroll)
		ViewportSize:  10,
	}
}
