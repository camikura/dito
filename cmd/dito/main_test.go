package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/camikura/dito/internal/app"
	"github.com/camikura/dito/internal/db"
)

func TestModelInit(t *testing.T) {
	m := model{Model: app.InitialModel()}

	cmd := m.Init()

	// Init should return nil
	if cmd != nil {
		t.Error("Init() should return nil")
	}
}

func TestModelUpdate(t *testing.T) {
	tests := []struct {
		name          string
		initialModel  model
		msg           tea.Msg
		expectCmd     bool
		expectedScreen app.Screen
	}{
		{
			name: "window size message",
			initialModel: model{
				Model: app.Model{
					Screen: app.ScreenSelection,
					Width:  0,
					Height: 0,
				},
			},
			msg: tea.WindowSizeMsg{
				Width:  100,
				Height: 30,
			},
			expectCmd:      false,
			expectedScreen: app.ScreenSelection,
		},
		{
			name: "key message - selection screen",
			initialModel: model{
				Model: app.Model{
					Screen:  app.ScreenSelection,
					Choices: []string{"Cloud", "On-Premise"},
					Cursor:  0,
				},
			},
			msg: tea.KeyMsg{
				Type: tea.KeyDown,
			},
			expectCmd:      false,
			expectedScreen: app.ScreenSelection,
		},
		{
			name: "connection result - success",
			initialModel: model{
				Model: app.Model{
					Screen: app.ScreenOnPremiseConfig,
					OnPremiseConfig: app.OnPremiseConfig{
						Status: app.StatusConnecting,
					},
				},
			},
			msg: connectionResultMsg{
				Err:     nil,
				Version: "Oracle NoSQL Database 23.1",
				IsTest:  true,
			},
			expectCmd:      false,
			expectedScreen: app.ScreenOnPremiseConfig,
		},
		{
			name: "table list result",
			initialModel: model{
				Model: app.Model{
					Screen: app.ScreenOnPremiseConfig,
					OnPremiseConfig: app.OnPremiseConfig{
						Status: app.StatusConnected,
					},
					TableDetails: make(map[string]*db.TableDetailsResult),
				},
			},
			msg: tableListResultMsg{
				Tables: []string{"users", "products"},
				Err:    nil,
			},
			expectCmd:      true, // Should fetch details for first table
			expectedScreen: app.ScreenTableList,
		},
		{
			name: "table details result",
			initialModel: model{
				Model: app.Model{
					Screen:         app.ScreenTableList,
					Tables:         []string{"users"},
					TableDetails:   make(map[string]*db.TableDetailsResult),
					LoadingDetails: true,
					RightPaneMode:  app.RightPaneModeSchema,
				},
			},
			msg: tableDetailsResultMsg{
				TableName: "users",
				Schema:    nil,
				Err:       nil,
			},
			expectCmd:      false,
			expectedScreen: app.ScreenTableList,
		},
		{
			name: "table data result",
			initialModel: model{
				Model: app.Model{
					Screen:       app.ScreenTableList,
					Tables:       []string{"users"},
					TableData:    make(map[string]*db.TableDataResult),
					LoadingData:  true,
					RightPaneMode: app.RightPaneModeList,
				},
			},
			msg: tableDataResultMsg{
				TableName: "users",
				Rows: []map[string]interface{}{
					{"id": 1, "name": "Alice"},
				},
				Err: nil,
			},
			expectCmd:      false,
			expectedScreen: app.ScreenTableList,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultModel, resultCmd := tt.initialModel.Update(tt.msg)

			// Check that model is returned
			if resultModel == nil {
				t.Error("Update() should return a model")
			}

			// Check command
			if tt.expectCmd && resultCmd == nil {
				t.Error("Update() should return a command")
			}
			if !tt.expectCmd && resultCmd != nil && tt.name != "key message - selection screen" {
				t.Error("Update() should not return a command")
			}

			// Check screen (type assertion to access internal model)
			if m, ok := resultModel.(model); ok {
				if m.Screen != tt.expectedScreen {
					t.Errorf("Update() Screen = %v, want %v", m.Screen, tt.expectedScreen)
				}

				// Additional checks based on message type
				switch msg := tt.msg.(type) {
				case tea.WindowSizeMsg:
					if m.Width != msg.Width {
						t.Errorf("Update() Width = %d, want %d", m.Width, msg.Width)
					}
					if m.Height != msg.Height {
						t.Errorf("Update() Height = %d, want %d", m.Height, msg.Height)
					}
				case connectionResultMsg:
					// Status should always be reset to Disconnected
					if m.OnPremiseConfig.Status != app.StatusDisconnected {
						t.Error("Update() should reset status to disconnected")
					}
					// For test connections or errors, dialog should be visible
					if msg.IsTest || msg.Err != nil {
						if !m.DialogVisible {
							t.Error("Update() should show dialog for test connection or error")
						}
					}
				case tableListResultMsg:
					if msg.Err == nil && len(msg.Tables) > 0 {
						if len(m.Tables) != len(msg.Tables) {
							t.Errorf("Update() Tables length = %d, want %d", len(m.Tables), len(msg.Tables))
						}
					}
				case tableDetailsResultMsg:
					if msg.Err == nil {
						if m.LoadingDetails {
							t.Error("Update() should set LoadingDetails to false")
						}
					}
				case tableDataResultMsg:
					if msg.Err == nil {
						if m.LoadingData {
							t.Error("Update() should set LoadingData to false")
						}
					}
				}
			}
		})
	}
}

func TestModelView(t *testing.T) {
	tests := []struct {
		name     string
		model    model
		contains []string
	}{
		{
			name: "loading state",
			model: model{
				Model: app.Model{
					Width:  0,
					Height: 0,
				},
			},
			contains: []string{"Loading..."},
		},
		{
			name: "selection screen",
			model: model{
				Model: app.Model{
					Screen:  app.ScreenSelection,
					Choices: []string{"Oracle NoSQL Cloud Service", "On-Premise"},
					Cursor:  0,
					Width:   100,
					Height:  30,
				},
			},
			contains: []string{"Dito", "Select Connection"},
		},
		{
			name: "on-premise config screen",
			model: model{
				Model: app.Model{
					Screen: app.ScreenOnPremiseConfig,
					OnPremiseConfig: app.OnPremiseConfig{
						Endpoint:  "localhost",
						Port:      "8080",
						Secure:    false,
						Focus:     0,
						CursorPos: 9,
					},
					Width:  100,
					Height: 30,
				},
			},
			contains: []string{"Dito", "On-Premise Connection", "localhost", "8080"},
		},
		{
			name: "cloud config screen",
			model: model{
				Model: app.Model{
					Screen: app.ScreenCloudConfig,
					CloudConfig: app.CloudConfig{
						Region:      "us-ashburn-1",
						Compartment: "",
						AuthMethod:  0,
						ConfigFile:  "DEFAULT",
						Focus:       0,
						CursorPos:   12,
					},
					Width:  100,
					Height: 30,
				},
			},
			contains: []string{"Dito", "Cloud Connection", "us-ashburn-1"},
		},
		{
			name: "table list screen",
			model: model{
				Model: app.Model{
					Screen:        app.ScreenTableList,
					Tables:        []string{"users", "products"},
					SelectedTable: 0,
					Endpoint:      "localhost:8080",
					RightPaneMode: app.RightPaneModeSchema,
					TableDetails:  make(map[string]*db.TableDetailsResult),
					TableData:     make(map[string]*db.TableDataResult),
					Width:         100,
					Height:        30,
				},
			},
			contains: []string{"Dito", "Connected to localhost:8080", "Tables"},
		},
		{
			name: "unknown screen",
			model: model{
				Model: app.Model{
					Screen: 999, // invalid screen
					Width:  100,
					Height: 30,
				},
			},
			contains: []string{"Unknown screen"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.model.View()

			// Check that result is not empty (unless loading)
			if result == "" && tt.name != "loading state" {
				t.Error("View() should not return empty string")
			}

			// Check that expected strings are present
			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("View() should contain %q", substr)
				}
			}

			// Check that border characters are present (unless loading)
			if tt.name != "loading state" && tt.name != "unknown screen" {
				if !strings.Contains(result, "╭") || !strings.Contains(result, "╮") ||
					!strings.Contains(result, "╰") || !strings.Contains(result, "╯") {
					t.Error("View() should contain border characters")
				}
			}
		})
	}
}
