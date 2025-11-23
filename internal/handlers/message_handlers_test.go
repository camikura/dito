package handlers

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/oracle/nosql-go-sdk/nosqldb"

	"github.com/camikura/dito/internal/app"
	"github.com/camikura/dito/internal/db"
)

func TestHandleWindowSize(t *testing.T) {
	tests := []struct {
		name             string
		initialModel     app.Model
		msg              tea.WindowSizeMsg
		expectedWidth    int
		expectedHeight   int
		expectedViewport int
	}{
		{
			name: "normal window size",
			initialModel: app.Model{
				Width:  0,
				Height: 0,
			},
			msg: tea.WindowSizeMsg{
				Width:  100,
				Height: 30,
			},
			expectedWidth:    100,
			expectedHeight:   30,
			expectedViewport: 19, // (30 - 8) - 3 = 19
		},
		{
			name: "small window size",
			initialModel: app.Model{
				Width:  0,
				Height: 0,
			},
			msg: tea.WindowSizeMsg{
				Width:  50,
				Height: 10,
			},
			expectedWidth:    50,
			expectedHeight:   10,
			expectedViewport: 1, // minimum viewport size is 1
		},
		{
			name: "large window size",
			initialModel: app.Model{
				Width:  0,
				Height: 0,
			},
			msg: tea.WindowSizeMsg{
				Width:  200,
				Height: 60,
			},
			expectedWidth:    200,
			expectedHeight:   60,
			expectedViewport: 49, // (60 - 8) - 3 = 49
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HandleWindowSize(tt.initialModel, tt.msg)

			if result.Width != tt.expectedWidth {
				t.Errorf("HandleWindowSize() Width = %d, want %d", result.Width, tt.expectedWidth)
			}
			if result.Height != tt.expectedHeight {
				t.Errorf("HandleWindowSize() Height = %d, want %d", result.Height, tt.expectedHeight)
			}
			if result.ViewportSize != tt.expectedViewport {
				t.Errorf("HandleWindowSize() ViewportSize = %d, want %d", result.ViewportSize, tt.expectedViewport)
			}
		})
	}
}

func TestHandleKeyPress(t *testing.T) {
	tests := []struct {
		name          string
		initialModel  app.Model
		msg           tea.KeyMsg
		expectedCmd   bool // whether a command should be returned
	}{
		{
			name: "selection screen",
			initialModel: app.Model{
				Screen: app.ScreenSelection,
			},
			msg: tea.KeyMsg{
				Type: tea.KeyRunes,
			},
			expectedCmd: false,
		},
		{
			name: "on-premise config screen",
			initialModel: app.Model{
				Screen: app.ScreenOnPremiseConfig,
			},
			msg: tea.KeyMsg{
				Type: tea.KeyRunes,
			},
			expectedCmd: false,
		},
		{
			name: "cloud config screen",
			initialModel: app.Model{
				Screen: app.ScreenCloudConfig,
			},
			msg: tea.KeyMsg{
				Type: tea.KeyRunes,
			},
			expectedCmd: false,
		},
		{
			name: "table list screen",
			initialModel: app.Model{
				Screen: app.ScreenTableList,
			},
			msg: tea.KeyMsg{
				Type: tea.KeyRunes,
			},
			expectedCmd: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultModel, resultCmd := HandleKeyPress(tt.initialModel, tt.msg)

			// Check that model is returned
			if resultModel.Screen != tt.initialModel.Screen {
				t.Errorf("HandleKeyPress() changed screen unexpectedly")
			}

			// Check command return
			if tt.expectedCmd && resultCmd == nil {
				t.Errorf("HandleKeyPress() should return a command but got nil")
			}
		})
	}
}

func TestHandleConnectionResult(t *testing.T) {
	tests := []struct {
		name           string
		initialModel   app.Model
		msg            db.ConnectionResult
		expectedStatus app.ConnectionStatus
		expectCmd      bool
	}{
		{
			name: "successful connection - test mode",
			initialModel: app.Model{
				OnPremiseConfig: app.OnPremiseConfig{
					Status: app.StatusConnecting,
				},
			},
			msg: db.ConnectionResult{
				Err:     nil,
				Version: "Oracle NoSQL Database 23.1",
				IsTest:  true,
			},
			expectedStatus: app.StatusConnected,
			expectCmd:      false,
		},
		{
			name: "successful connection - not test mode",
			initialModel: app.Model{
				OnPremiseConfig: app.OnPremiseConfig{
					Status: app.StatusConnecting,
				},
			},
			msg: db.ConnectionResult{
				Err:      nil,
				Version:  "Oracle NoSQL Database 23.1",
				IsTest:   false,
				Endpoint: "localhost:8080",
			},
			expectedStatus: app.StatusConnected,
			expectCmd:      true, // should fetch tables
		},
		{
			name: "failed connection",
			initialModel: app.Model{
				OnPremiseConfig: app.OnPremiseConfig{
					Status: app.StatusConnecting,
				},
			},
			msg: db.ConnectionResult{
				Err:    errors.New("connection refused"),
				IsTest: false,
			},
			expectedStatus: app.StatusError,
			expectCmd:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultModel, resultCmd := HandleConnectionResult(tt.initialModel, tt.msg)

			if resultModel.OnPremiseConfig.Status != tt.expectedStatus {
				t.Errorf("HandleConnectionResult() Status = %v, want %v",
					resultModel.OnPremiseConfig.Status, tt.expectedStatus)
			}

			if tt.expectCmd && resultCmd == nil {
				t.Errorf("HandleConnectionResult() should return a command but got nil")
			}
			if !tt.expectCmd && resultCmd != nil {
				t.Errorf("HandleConnectionResult() should not return a command but got one")
			}

			// Check error message is set when there's an error
			if tt.msg.Err != nil && resultModel.OnPremiseConfig.ErrorMsg == "" {
				t.Errorf("HandleConnectionResult() should set ErrorMsg when error occurs")
			}

			// Check version is set when connection is successful
			if tt.msg.Err == nil && tt.msg.Version != "" &&
				resultModel.OnPremiseConfig.ServerVersion != tt.msg.Version {
				t.Errorf("HandleConnectionResult() ServerVersion = %v, want %v",
					resultModel.OnPremiseConfig.ServerVersion, tt.msg.Version)
			}
		})
	}
}

func TestHandleTableListResult(t *testing.T) {
	tests := []struct {
		name           string
		initialModel   app.Model
		msg            db.TableListResult
		expectedScreen app.Screen
		expectCmd      bool
	}{
		{
			name: "successful table list fetch",
			initialModel: app.Model{
				Screen: app.ScreenOnPremiseConfig,
				OnPremiseConfig: app.OnPremiseConfig{
					Status: app.StatusConnected,
				},
				TableDetails: make(map[string]*db.TableDetailsResult),
			},
			msg: db.TableListResult{
				Tables: []string{"users", "products", "orders"},
				Err:    nil,
			},
			expectedScreen: app.ScreenTableList,
			expectCmd:      true, // should fetch details for first table
		},
		{
			name: "successful table list fetch - empty tables",
			initialModel: app.Model{
				Screen: app.ScreenOnPremiseConfig,
				OnPremiseConfig: app.OnPremiseConfig{
					Status: app.StatusConnected,
				},
				TableDetails: make(map[string]*db.TableDetailsResult),
			},
			msg: db.TableListResult{
				Tables: []string{},
				Err:    nil,
			},
			expectedScreen: app.ScreenTableList,
			expectCmd:      false, // no tables to fetch details for
		},
		{
			name: "failed table list fetch",
			initialModel: app.Model{
				Screen: app.ScreenOnPremiseConfig,
				OnPremiseConfig: app.OnPremiseConfig{
					Status: app.StatusConnected,
				},
			},
			msg: db.TableListResult{
				Tables: nil,
				Err:    errors.New("failed to fetch tables"),
			},
			expectedScreen: app.ScreenOnPremiseConfig,
			expectCmd:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultModel, resultCmd := HandleTableListResult(tt.initialModel, tt.msg)

			if resultModel.Screen != tt.expectedScreen {
				t.Errorf("HandleTableListResult() Screen = %v, want %v",
					resultModel.Screen, tt.expectedScreen)
			}

			if tt.expectCmd && resultCmd == nil {
				t.Errorf("HandleTableListResult() should return a command but got nil")
			}
			if !tt.expectCmd && resultCmd != nil {
				t.Errorf("HandleTableListResult() should not return a command but got one")
			}

			// Check error message is set when there's an error
			if tt.msg.Err != nil && resultModel.OnPremiseConfig.ErrorMsg == "" {
				t.Errorf("HandleTableListResult() should set ErrorMsg when error occurs")
			}

			// Check tables are set when successful
			if tt.msg.Err == nil && len(tt.msg.Tables) > 0 {
				if len(resultModel.Tables) != len(tt.msg.Tables) {
					t.Errorf("HandleTableListResult() Tables length = %d, want %d",
						len(resultModel.Tables), len(tt.msg.Tables))
				}
			}
		})
	}
}

func TestHandleTableDetailsResult(t *testing.T) {
	tests := []struct {
		name         string
		initialModel app.Model
		msg          db.TableDetailsResult
		expectCmd    bool
	}{
		{
			name: "successful table details fetch - schema view mode",
			initialModel: app.Model{
				TableDetails:   make(map[string]*db.TableDetailsResult),
				LoadingDetails: true,
				RightPaneMode:  app.RightPaneModeSchema,
			},
			msg: db.TableDetailsResult{
				TableName: "users",
				Schema: &nosqldb.TableResult{
					DDL: "CREATE TABLE users (id INTEGER, name STRING, PRIMARY KEY(id))",
				},
				Err: nil,
			},
			expectCmd: false,
		},
		{
			name: "successful table details fetch - list view mode",
			initialModel: app.Model{
				TableDetails:   make(map[string]*db.TableDetailsResult),
				TableData:      make(map[string]*db.TableDataResult),
				LoadingDetails: true,
				RightPaneMode:  app.RightPaneModeList,
			},
			msg: db.TableDetailsResult{
				TableName: "users",
				Schema: &nosqldb.TableResult{
					DDL: "CREATE TABLE users (id INTEGER, name STRING, PRIMARY KEY(id))",
				},
				Err: nil,
			},
			expectCmd: true, // should fetch data
		},
		{
			name: "failed table details fetch",
			initialModel: app.Model{
				TableDetails:   make(map[string]*db.TableDetailsResult),
				LoadingDetails: true,
			},
			msg: db.TableDetailsResult{
				TableName: "users",
				Err:       errors.New("failed to fetch details"),
			},
			expectCmd: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultModel, resultCmd := HandleTableDetailsResult(tt.initialModel, tt.msg)

			// Check LoadingDetails is set to false
			if resultModel.LoadingDetails != false {
				t.Errorf("HandleTableDetailsResult() LoadingDetails should be false")
			}

			// Check table details are stored when successful
			if tt.msg.Err == nil {
				if _, exists := resultModel.TableDetails[tt.msg.TableName]; !exists {
					t.Errorf("HandleTableDetailsResult() should store table details")
				}
			}

			if tt.expectCmd && resultCmd == nil {
				t.Errorf("HandleTableDetailsResult() should return a command but got nil")
			}
			if !tt.expectCmd && resultCmd != nil {
				t.Errorf("HandleTableDetailsResult() should not return a command but got one")
			}
		})
	}
}

func TestHandleTableDataResult(t *testing.T) {
	tests := []struct {
		name         string
		initialModel app.Model
		msg          db.TableDataResult
		expectedRows int
	}{
		{
			name: "successful table data fetch - new data",
			initialModel: app.Model{
				TableData:   make(map[string]*db.TableDataResult),
				LoadingData: true,
			},
			msg: db.TableDataResult{
				TableName: "users",
				Rows: []map[string]interface{}{
					{"id": 1, "name": "Alice"},
					{"id": 2, "name": "Bob"},
				},
				IsAppend: false,
				Err:      nil,
			},
			expectedRows: 2,
		},
		{
			name: "successful table data fetch - append data",
			initialModel: app.Model{
				TableData: map[string]*db.TableDataResult{
					"users": {
						TableName: "users",
						Rows: []map[string]interface{}{
							{"id": 1, "name": "Alice"},
						},
					},
				},
				LoadingData: true,
			},
			msg: db.TableDataResult{
				TableName: "users",
				Rows: []map[string]interface{}{
					{"id": 2, "name": "Bob"},
				},
				IsAppend: true,
				Err:      nil,
			},
			expectedRows: 2, // 1 existing + 1 new
		},
		{
			name: "failed table data fetch",
			initialModel: app.Model{
				TableData:   make(map[string]*db.TableDataResult),
				LoadingData: true,
			},
			msg: db.TableDataResult{
				TableName: "users",
				Err:       errors.New("failed to fetch data"),
			},
			expectedRows: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultModel, resultCmd := HandleTableDataResult(tt.initialModel, tt.msg)

			// Check LoadingData is set to false
			if resultModel.LoadingData != false {
				t.Errorf("HandleTableDataResult() LoadingData should be false")
			}

			// Check command is not returned (data handlers don't return commands)
			if resultCmd != nil {
				t.Errorf("HandleTableDataResult() should not return a command")
			}

			// Check table data is stored
			if data, exists := resultModel.TableData[tt.msg.TableName]; exists {
				if tt.msg.Err == nil && len(data.Rows) != tt.expectedRows {
					t.Errorf("HandleTableDataResult() Rows length = %d, want %d",
						len(data.Rows), tt.expectedRows)
				}
			} else {
				t.Errorf("HandleTableDataResult() should store table data")
			}
		})
	}
}
