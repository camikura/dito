package handlers

import (
	"testing"

	"github.com/camikura/dito/internal/app"
	"github.com/camikura/dito/internal/db"
)

func TestHandleTableList(t *testing.T) {
	tests := []struct {
		name               string
		initialModel       app.Model
		key                string
		expectedScreen     app.Screen
		expectedRightPane  app.RightPaneMode
		expectedTableIndex int
		expectedDataRow    int
		expectQuitCmd      bool
	}{
		{
			name: "quit with ctrl+c",
			initialModel: app.Model{
				Screen:        app.ScreenTableList,
				Tables:        []string{"users", "products"},
				SelectedTable: 0,
				RightPaneMode: app.RightPaneModeSchema,
			},
			key:                "ctrl+c",
			expectedScreen:     app.ScreenTableList,
			expectedRightPane:  app.RightPaneModeSchema,
			expectedTableIndex: 0,
			expectedDataRow:    0,
			expectQuitCmd:      true,
		},
		{
			name: "navigate down in schema mode",
			initialModel: app.Model{
				Screen:        app.ScreenTableList,
				Tables:        []string{"users", "products"},
				SelectedTable: 0,
				RightPaneMode: app.RightPaneModeSchema,
			},
			key:                "down",
			expectedScreen:     app.ScreenTableList,
			expectedRightPane:  app.RightPaneModeSchema,
			expectedTableIndex: 1,
			expectedDataRow:    0,
			expectQuitCmd:      false,
		},
		{
			name: "navigate up in schema mode",
			initialModel: app.Model{
				Screen:        app.ScreenTableList,
				Tables:        []string{"users", "products"},
				SelectedTable: 1,
				RightPaneMode: app.RightPaneModeSchema,
			},
			key:                "up",
			expectedScreen:     app.ScreenTableList,
			expectedRightPane:  app.RightPaneModeSchema,
			expectedTableIndex: 0,
			expectedDataRow:    0,
			expectQuitCmd:      false,
		},
		{
			name: "don't navigate up beyond first table",
			initialModel: app.Model{
				Screen:        app.ScreenTableList,
				Tables:        []string{"users", "products"},
				SelectedTable: 0,
				RightPaneMode: app.RightPaneModeSchema,
			},
			key:                "up",
			expectedScreen:     app.ScreenTableList,
			expectedRightPane:  app.RightPaneModeSchema,
			expectedTableIndex: 0,
			expectedDataRow:    0,
			expectQuitCmd:      false,
		},
		{
			name: "don't navigate down beyond last table",
			initialModel: app.Model{
				Screen:        app.ScreenTableList,
				Tables:        []string{"users", "products"},
				SelectedTable: 1,
				RightPaneMode: app.RightPaneModeSchema,
			},
			key:                "down",
			expectedScreen:     app.ScreenTableList,
			expectedRightPane:  app.RightPaneModeSchema,
			expectedTableIndex: 1,
			expectedDataRow:    0,
			expectQuitCmd:      false,
		},
		{
			name: "enter switches from schema to list mode",
			initialModel: app.Model{
				Screen:        app.ScreenTableList,
				Tables:        []string{"users"},
				SelectedTable: 0,
				RightPaneMode: app.RightPaneModeSchema,
				TableDetails:  make(map[string]*db.TableDetailsResult),
				TableData:     make(map[string]*db.TableDataResult),
			},
			key:                "enter",
			expectedScreen:     app.ScreenTableList,
			expectedRightPane:  app.RightPaneModeList,
			expectedTableIndex: 0,
			expectedDataRow:    0,
			expectQuitCmd:      false,
		},
		{
			name: "enter switches from list to detail mode",
			initialModel: app.Model{
				Screen:        app.ScreenTableList,
				Tables:        []string{"users"},
				SelectedTable: 0,
				RightPaneMode: app.RightPaneModeList,
				TableData: map[string]*db.TableDataResult{
					"users": {
						TableName: "users",
						Rows: []map[string]interface{}{
							{"id": 1, "name": "Alice"},
						},
					},
				},
			},
			key:                "enter",
			expectedScreen:     app.ScreenTableList,
			expectedRightPane:  app.RightPaneModeDetail,
			expectedTableIndex: 0,
			expectedDataRow:    0,
			expectQuitCmd:      false,
		},
		{
			name: "esc switches from detail to list mode",
			initialModel: app.Model{
				Screen:        app.ScreenTableList,
				Tables:        []string{"users"},
				SelectedTable: 0,
				RightPaneMode: app.RightPaneModeDetail,
			},
			key:                "esc",
			expectedScreen:     app.ScreenTableList,
			expectedRightPane:  app.RightPaneModeList,
			expectedTableIndex: 0,
			expectedDataRow:    0,
			expectQuitCmd:      false,
		},
		{
			name: "esc switches from schema to on-premise config and resets connection status",
			initialModel: app.Model{
				Screen:        app.ScreenTableList,
				Tables:        []string{"users"},
				SelectedTable: 0,
				RightPaneMode: app.RightPaneModeSchema,
				OnPremiseConfig: app.OnPremiseConfig{
					Status:        app.StatusConnected,
					ServerVersion: "Oracle NoSQL Database 23.1",
					ErrorMsg:      "",
				},
			},
			key:                "esc",
			expectedScreen:     app.ScreenOnPremiseConfig,
			expectedRightPane:  app.RightPaneModeSchema,
			expectedTableIndex: 0,
			expectedDataRow:    0,
			expectQuitCmd:      false,
		},
		{
			name: "navigate down data row in list mode",
			initialModel: app.Model{
				Screen:          app.ScreenTableList,
				Tables:          []string{"users"},
				SelectedTable:   0,
				RightPaneMode:   app.RightPaneModeList,
				SelectedDataRow: 0,
				ViewportSize:    10,
				TableData: map[string]*db.TableDataResult{
					"users": {
						TableName: "users",
						Rows: []map[string]interface{}{
							{"id": 1, "name": "Alice"},
							{"id": 2, "name": "Bob"},
						},
					},
				},
			},
			key:                "down",
			expectedScreen:     app.ScreenTableList,
			expectedRightPane:  app.RightPaneModeList,
			expectedTableIndex: 0,
			expectedDataRow:    1,
			expectQuitCmd:      false,
		},
		{
			name: "navigate up data row in list mode",
			initialModel: app.Model{
				Screen:          app.ScreenTableList,
				Tables:          []string{"users"},
				SelectedTable:   0,
				RightPaneMode:   app.RightPaneModeList,
				SelectedDataRow: 1,
				ViewportSize:    10,
				TableData: map[string]*db.TableDataResult{
					"users": {
						TableName: "users",
						Rows: []map[string]interface{}{
							{"id": 1, "name": "Alice"},
							{"id": 2, "name": "Bob"},
						},
					},
				},
			},
			key:                "up",
			expectedScreen:     app.ScreenTableList,
			expectedRightPane:  app.RightPaneModeList,
			expectedTableIndex: 0,
			expectedDataRow:    0,
			expectQuitCmd:      false,
		},
		{
			name: "don't navigate up beyond first data row",
			initialModel: app.Model{
				Screen:          app.ScreenTableList,
				Tables:          []string{"users"},
				SelectedTable:   0,
				RightPaneMode:   app.RightPaneModeList,
				SelectedDataRow: 0,
				TableData: map[string]*db.TableDataResult{
					"users": {
						TableName: "users",
						Rows: []map[string]interface{}{
							{"id": 1, "name": "Alice"},
						},
					},
				},
			},
			key:                "up",
			expectedScreen:     app.ScreenTableList,
			expectedRightPane:  app.RightPaneModeList,
			expectedTableIndex: 0,
			expectedDataRow:    0,
			expectQuitCmd:      false,
		},
		{
			name: "scroll left in list mode",
			initialModel: app.Model{
				Screen:           app.ScreenTableList,
				Tables:           []string{"users"},
				SelectedTable:    0,
				RightPaneMode:    app.RightPaneModeList,
				HorizontalOffset: 5,
			},
			key:                "left",
			expectedScreen:     app.ScreenTableList,
			expectedRightPane:  app.RightPaneModeList,
			expectedTableIndex: 0,
			expectedDataRow:    0,
			expectQuitCmd:      false,
		},
		{
			name: "don't scroll left beyond 0",
			initialModel: app.Model{
				Screen:           app.ScreenTableList,
				Tables:           []string{"users"},
				SelectedTable:    0,
				RightPaneMode:    app.RightPaneModeList,
				HorizontalOffset: 0,
			},
			key:                "left",
			expectedScreen:     app.ScreenTableList,
			expectedRightPane:  app.RightPaneModeList,
			expectedTableIndex: 0,
			expectedDataRow:    0,
			expectQuitCmd:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := createKeyMsg(tt.key)
			resultModel, resultCmd := HandleTableList(tt.initialModel, msg)

			// Check screen
			if resultModel.Screen != tt.expectedScreen {
				t.Errorf("HandleTableList() Screen = %v, want %v", resultModel.Screen, tt.expectedScreen)
			}

			// Check right pane mode
			if resultModel.RightPaneMode != tt.expectedRightPane {
				t.Errorf("HandleTableList() RightPaneMode = %v, want %v", resultModel.RightPaneMode, tt.expectedRightPane)
			}

			// Check selected table index (in schema mode)
			if tt.initialModel.RightPaneMode == app.RightPaneModeSchema {
				if resultModel.SelectedTable != tt.expectedTableIndex {
					t.Errorf("HandleTableList() SelectedTable = %v, want %v", resultModel.SelectedTable, tt.expectedTableIndex)
				}
			}

			// Check selected data row (in list/detail mode)
			if tt.initialModel.RightPaneMode == app.RightPaneModeList || tt.initialModel.RightPaneMode == app.RightPaneModeDetail {
				if resultModel.SelectedDataRow != tt.expectedDataRow {
					t.Errorf("HandleTableList() SelectedDataRow = %v, want %v", resultModel.SelectedDataRow, tt.expectedDataRow)
				}
			}

			// Check quit command
			if tt.expectQuitCmd {
				if resultCmd == nil {
					t.Error("HandleTableList() should return tea.Quit command")
				}
			}

			// Special checks
			if tt.key == "enter" && tt.initialModel.RightPaneMode == app.RightPaneModeSchema {
				if resultModel.SelectedDataRow != 0 {
					t.Error("HandleTableList() should reset SelectedDataRow when switching to list mode")
				}
				if resultModel.ViewportOffset != 0 {
					t.Error("HandleTableList() should reset ViewportOffset when switching to list mode")
				}
				if resultModel.HorizontalOffset != 0 {
					t.Error("HandleTableList() should reset HorizontalOffset when switching to list mode")
				}
			}

			if (tt.key == "left" || tt.key == "h") && tt.initialModel.RightPaneMode == app.RightPaneModeList && tt.initialModel.HorizontalOffset > 0 {
				expectedOffset := tt.initialModel.HorizontalOffset - 1
				if resultModel.HorizontalOffset != expectedOffset {
					t.Errorf("HandleTableList() HorizontalOffset = %v, want %v", resultModel.HorizontalOffset, expectedOffset)
				}
			}

			// Check connection status reset when returning to config screen
			if (tt.key == "esc" || tt.key == "u") && tt.initialModel.RightPaneMode == app.RightPaneModeSchema && tt.expectedScreen == app.ScreenOnPremiseConfig {
				if resultModel.OnPremiseConfig.Status != app.StatusDisconnected {
					t.Errorf("HandleTableList() should reset OnPremiseConfig.Status to StatusDisconnected, got %v", resultModel.OnPremiseConfig.Status)
				}
				if resultModel.OnPremiseConfig.ServerVersion != "" {
					t.Errorf("HandleTableList() should clear OnPremiseConfig.ServerVersion, got %v", resultModel.OnPremiseConfig.ServerVersion)
				}
				if resultModel.OnPremiseConfig.ErrorMsg != "" {
					t.Errorf("HandleTableList() should clear OnPremiseConfig.ErrorMsg, got %v", resultModel.OnPremiseConfig.ErrorMsg)
				}
			}
		})
	}
}
