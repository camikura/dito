package handlers

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/camikura/dito/internal/app"
)

func TestHandleOnPremiseConfig(t *testing.T) {
	tests := []struct {
		name           string
		initialModel   app.Model
		key            string
		expectedScreen app.Screen
		expectedFocus  int
		expectQuitCmd  bool
		expectDbCmd    bool
	}{
		{
			name: "quit with ctrl+c",
			initialModel: app.Model{
				Screen: app.ScreenOnPremiseConfig,
				OnPremiseConfig: app.OnPremiseConfig{
					Endpoint:  "localhost",
					Port:      "8080",
					Focus:     0,
				},
			},
			key:            "ctrl+c",
			expectedScreen: app.ScreenOnPremiseConfig,
			expectedFocus:  0,
			expectQuitCmd:  true,
			expectDbCmd:    false,
		},
		{
			name: "back to selection with esc",
			initialModel: app.Model{
				Screen: app.ScreenOnPremiseConfig,
				OnPremiseConfig: app.OnPremiseConfig{
					Endpoint:  "localhost",
					Port:      "8080",
					Focus:     0,
				},
			},
			key:            "esc",
			expectedScreen: app.ScreenSelection,
			expectedFocus:  0,
			expectQuitCmd:  false,
			expectDbCmd:    false,
		},
		{
			name: "down to next field",
			initialModel: app.Model{
				Screen: app.ScreenOnPremiseConfig,
				OnPremiseConfig: app.OnPremiseConfig{
					Endpoint:  "localhost",
					Port:      "8080",
					Focus:     0,
				},
			},
			key:            "down",
			expectedScreen: app.ScreenOnPremiseConfig,
			expectedFocus:  1,
			expectQuitCmd:  false,
			expectDbCmd:    false,
		},
		{
			name: "up to previous field",
			initialModel: app.Model{
				Screen: app.ScreenOnPremiseConfig,
				OnPremiseConfig: app.OnPremiseConfig{
					Endpoint:  "localhost",
					Port:      "8080",
					Focus:     1,
				},
			},
			key:            "up",
			expectedScreen: app.ScreenOnPremiseConfig,
			expectedFocus:  0,
			expectQuitCmd:  false,
			expectDbCmd:    false,
		},
		{
			name: "space toggles secure checkbox",
			initialModel: app.Model{
				Screen: app.ScreenOnPremiseConfig,
				OnPremiseConfig: app.OnPremiseConfig{
					Endpoint:  "localhost",
					Port:      "8080",
					Secure:    false,
					Focus:     2,
				},
			},
			key:            " ",
			expectedScreen: app.ScreenOnPremiseConfig,
			expectedFocus:  2,
			expectQuitCmd:  false,
			expectDbCmd:    false,
		},
		{
			name: "enter on test connection button",
			initialModel: app.Model{
				Screen: app.ScreenOnPremiseConfig,
				OnPremiseConfig: app.OnPremiseConfig{
					Endpoint:  "localhost",
					Port:      "8080",
					Focus:     3,
				},
			},
			key:            "enter",
			expectedScreen: app.ScreenOnPremiseConfig,
			expectedFocus:  3,
			expectQuitCmd:  false,
			expectDbCmd:    true, // db.Connect command should be returned
		},
		{
			name: "enter on connect button",
			initialModel: app.Model{
				Screen: app.ScreenOnPremiseConfig,
				OnPremiseConfig: app.OnPremiseConfig{
					Endpoint:  "localhost",
					Port:      "8080",
					Focus:     4,
				},
			},
			key:            "enter",
			expectedScreen: app.ScreenOnPremiseConfig,
			expectedFocus:  4,
			expectQuitCmd:  false,
			expectDbCmd:    true, // db.Connect command should be returned
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := createKeyMsg(tt.key)
			resultModel, resultCmd := HandleOnPremiseConfig(tt.initialModel, msg)

			// Check screen
			if resultModel.Screen != tt.expectedScreen {
				t.Errorf("HandleOnPremiseConfig() Screen = %v, want %v", resultModel.Screen, tt.expectedScreen)
			}

			// Check focus
			if resultModel.OnPremiseConfig.Focus != tt.expectedFocus {
				t.Errorf("HandleOnPremiseConfig() Focus = %v, want %v", resultModel.OnPremiseConfig.Focus, tt.expectedFocus)
			}

			// Check quit command
			if tt.expectQuitCmd {
				if resultCmd == nil {
					t.Error("HandleOnPremiseConfig() should return tea.Quit command")
				}
			}

			// Check db command
			if tt.expectDbCmd {
				if resultCmd == nil {
					t.Error("HandleOnPremiseConfig() should return db.Connect command")
				}
			} else if !tt.expectQuitCmd && tt.key != "esc" {
				// Non-quit, non-db commands should return nil
				if resultCmd != nil && tt.key != "enter" {
					t.Error("HandleOnPremiseConfig() should not return a command for this key")
				}
			}

			// Special checks
			if tt.key == "backspace" && tt.initialModel.OnPremiseConfig.Focus == 0 {
				expectedEndpoint := "localhos" // "localhost" with last char deleted
				if resultModel.OnPremiseConfig.Endpoint != expectedEndpoint {
					t.Errorf("HandleOnPremiseConfig() Endpoint = %v, want %v", resultModel.OnPremiseConfig.Endpoint, expectedEndpoint)
				}
			}

			if tt.key == " " && tt.initialModel.OnPremiseConfig.Focus == 2 {
				if resultModel.OnPremiseConfig.Secure == tt.initialModel.OnPremiseConfig.Secure {
					t.Error("HandleOnPremiseConfig() should toggle Secure checkbox")
				}
			}

			if (tt.key == "enter") && (tt.initialModel.OnPremiseConfig.Focus == 3 || tt.initialModel.OnPremiseConfig.Focus == 4) {
				if resultModel.OnPremiseConfig.Status != app.StatusConnecting {
					t.Errorf("HandleOnPremiseConfig() Status should be StatusConnecting")
				}
			}
		})
	}
}

func TestHandleCloudConfig(t *testing.T) {
	tests := []struct {
		name           string
		initialModel   app.Model
		key            string
		expectedScreen app.Screen
		expectedFocus  int
		expectQuitCmd  bool
	}{
		{
			name: "quit with ctrl+c",
			initialModel: app.Model{
				Screen: app.ScreenCloudConfig,
				CloudConfig: app.CloudConfig{
					Region:      "us-ashburn-1",
					Compartment: "ocid1.compartment...",
					AuthMethod:  0,
					ConfigFile:  "DEFAULT",
					Focus:       0,
				},
			},
			key:            "ctrl+c",
			expectedScreen: app.ScreenCloudConfig,
			expectedFocus:  0,
			expectQuitCmd:  true,
		},
		{
			name: "back to selection with esc",
			initialModel: app.Model{
				Screen: app.ScreenCloudConfig,
				CloudConfig: app.CloudConfig{
					Region:      "us-ashburn-1",
					Compartment: "ocid1.compartment...",
					AuthMethod:  0,
					ConfigFile:  "DEFAULT",
					Focus:       0,
				},
			},
			key:            "esc",
			expectedScreen: app.ScreenSelection,
			expectedFocus:  0,
			expectQuitCmd:  false,
		},
		{
			name: "down to next field",
			initialModel: app.Model{
				Screen: app.ScreenCloudConfig,
				CloudConfig: app.CloudConfig{
					Region:      "us-ashburn-1",
					Compartment: "ocid1.compartment...",
					AuthMethod:  0,
					ConfigFile:  "DEFAULT",
					Focus:       0,
				},
			},
			key:            "down",
			expectedScreen: app.ScreenCloudConfig,
			expectedFocus:  1,
			expectQuitCmd:  false,
		},
		{
			name: "space selects radio button",
			initialModel: app.Model{
				Screen: app.ScreenCloudConfig,
				CloudConfig: app.CloudConfig{
					Region:      "us-ashburn-1",
					Compartment: "ocid1.compartment...",
					AuthMethod:  0,
					ConfigFile:  "DEFAULT",
					Focus:       2,
				},
			},
			key:            " ",
			expectedScreen: app.ScreenCloudConfig,
			expectedFocus:  2,
			expectQuitCmd:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := createKeyMsg(tt.key)
			resultModel, resultCmd := HandleCloudConfig(tt.initialModel, msg)

			// Check screen
			if resultModel.Screen != tt.expectedScreen {
				t.Errorf("HandleCloudConfig() Screen = %v, want %v", resultModel.Screen, tt.expectedScreen)
			}

			// Check focus
			if resultModel.CloudConfig.Focus != tt.expectedFocus {
				t.Errorf("HandleCloudConfig() Focus = %v, want %v", resultModel.CloudConfig.Focus, tt.expectedFocus)
			}

			// Check quit command
			if tt.expectQuitCmd {
				if resultCmd == nil {
					t.Error("HandleCloudConfig() should return tea.Quit command")
				}
			}

			// Special checks
			if tt.key == " " && tt.initialModel.CloudConfig.Focus >= 2 && tt.initialModel.CloudConfig.Focus <= 4 {
				expectedAuthMethod := tt.initialModel.CloudConfig.Focus - 2
				if resultModel.CloudConfig.AuthMethod != expectedAuthMethod {
					t.Errorf("HandleCloudConfig() AuthMethod = %v, want %v", resultModel.CloudConfig.AuthMethod, expectedAuthMethod)
				}
			}
		})
	}
}

// Helper function to create KeyMsg
func createKeyMsg(key string) tea.KeyMsg {
	msg := tea.KeyMsg{Type: tea.KeyRunes}
	switch key {
	case "ctrl+c":
		msg.Type = tea.KeyCtrlC
	case "esc":
		msg.Type = tea.KeyEsc
	case "tab":
		msg.Type = tea.KeyTab
	case "shift+tab":
		msg.Type = tea.KeyShiftTab
	case "enter":
		msg.Type = tea.KeyEnter
	case " ":
		msg.Type = tea.KeySpace
	case "left":
		msg.Type = tea.KeyLeft
	case "right":
		msg.Type = tea.KeyRight
	case "backspace":
		msg.Type = tea.KeyBackspace
	case "up":
		msg.Type = tea.KeyUp
	case "down":
		msg.Type = tea.KeyDown
	default:
		msg.Type = tea.KeyRunes
		msg.Runes = []rune{rune(key[0])}
	}
	return msg
}
