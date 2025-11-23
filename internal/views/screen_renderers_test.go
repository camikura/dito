package views

import (
	"strings"
	"testing"

	"github.com/camikura/dito/internal/app"
)

func TestRenderSelectionScreen(t *testing.T) {
	tests := []struct {
		name     string
		vm       ScreenViewModel
		contains []string
	}{
		{
			name: "selection screen with two choices",
			vm: ScreenViewModel{
				Width:  100,
				Height: 30,
				Model: app.Model{
					Choices: []string{"Oracle NoSQL Cloud Service", "On-Premise"},
					Cursor:  0,
				},
			},
			contains: []string{"Dito", "Select Connection", "Oracle NoSQL Cloud Service", "On-Premise"},
		},
		{
			name: "selection screen with single choice",
			vm: ScreenViewModel{
				Width:  80,
				Height: 25,
				Model: app.Model{
					Choices: []string{"On-Premise"},
					Cursor:  0,
				},
			},
			contains: []string{"Dito", "Select Connection", "On-Premise"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderSelectionScreen(tt.vm)

			// Check that all expected strings are present
			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("RenderSelectionScreen() should contain %q", substr)
				}
			}

			// Check that result is not empty
			if result == "" {
				t.Error("RenderSelectionScreen() should not return empty string")
			}

			// Check that border characters are present
			if !strings.Contains(result, "╭") || !strings.Contains(result, "╮") ||
				!strings.Contains(result, "╰") || !strings.Contains(result, "╯") {
				t.Error("RenderSelectionScreen() should contain border characters")
			}
		})
	}
}

func TestRenderOnPremiseConfigScreen(t *testing.T) {
	tests := []struct {
		name     string
		vm       ScreenViewModel
		contains []string
	}{
		{
			name: "on-premise config screen - disconnected",
			vm: ScreenViewModel{
				Width:  100,
				Height: 30,
				Model: app.Model{
					OnPremiseConfig: app.OnPremiseConfig{
						Endpoint:  "localhost",
						Port:      "8080",
						Secure:    false,
						Status:    app.StatusDisconnected,
						CursorPos: 9,
					},
				},
			},
			contains: []string{"Dito", "On-Premise Connection", "localhost", "8080"},
		},
		{
			name: "on-premise config screen - connecting",
			vm: ScreenViewModel{
				Width:  100,
				Height: 30,
				Model: app.Model{
					OnPremiseConfig: app.OnPremiseConfig{
						Endpoint:  "localhost",
						Port:      "8080",
						Secure:    false,
						Status:    app.StatusConnecting,
						CursorPos: 9,
					},
				},
			},
			contains: []string{"Dito", "On-Premise Connection", "Connecting"},
		},
		{
			name: "on-premise config screen - connected",
			vm: ScreenViewModel{
				Width:  100,
				Height: 30,
				Model: app.Model{
					OnPremiseConfig: app.OnPremiseConfig{
						Endpoint:      "localhost",
						Port:          "8080",
						Secure:        false,
						Status:        app.StatusConnected,
						ServerVersion: "Oracle NoSQL Database 23.1",
						CursorPos:     9,
					},
				},
			},
			contains: []string{"Dito", "On-Premise Connection", "Oracle NoSQL Database 23.1"},
		},
		{
			name: "on-premise config screen - error",
			vm: ScreenViewModel{
				Width:  100,
				Height: 30,
				Model: app.Model{
					OnPremiseConfig: app.OnPremiseConfig{
						Endpoint:  "localhost",
						Port:      "8080",
						Secure:    false,
						Status:    app.StatusError,
						ErrorMsg:  "Connection refused",
						CursorPos: 9,
					},
				},
			},
			contains: []string{"Dito", "On-Premise Connection", "Connection refused"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderOnPremiseConfigScreen(tt.vm)

			// Check that all expected strings are present
			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("RenderOnPremiseConfigScreen() should contain %q", substr)
				}
			}

			// Check that result is not empty
			if result == "" {
				t.Error("RenderOnPremiseConfigScreen() should not return empty string")
			}

			// Check that border characters are present
			if !strings.Contains(result, "╭") || !strings.Contains(result, "╮") ||
				!strings.Contains(result, "╰") || !strings.Contains(result, "╯") {
				t.Error("RenderOnPremiseConfigScreen() should contain border characters")
			}
		})
	}
}

func TestRenderCloudConfigScreen(t *testing.T) {
	tests := []struct {
		name     string
		vm       ScreenViewModel
		contains []string
	}{
		{
			name: "cloud config screen - disconnected",
			vm: ScreenViewModel{
				Width:  100,
				Height: 30,
				Model: app.Model{
					CloudConfig: app.CloudConfig{
						Region:      "us-ashburn-1",
						Compartment: "ocid1.compartment...",
						AuthMethod:  0,
						ConfigFile:  "DEFAULT",
						Status:      app.StatusDisconnected,
						CursorPos:   12,
					},
				},
			},
			contains: []string{"Dito", "Cloud Connection", "us-ashburn-1"},
		},
		{
			name: "cloud config screen - connected",
			vm: ScreenViewModel{
				Width:  100,
				Height: 30,
				Model: app.Model{
					CloudConfig: app.CloudConfig{
						Region:        "us-ashburn-1",
						Compartment:   "ocid1.compartment...",
						AuthMethod:    0,
						ConfigFile:    "DEFAULT",
						Status:        app.StatusConnected,
						ServerVersion: "Oracle NoSQL Cloud Service",
						CursorPos:     12,
					},
				},
			},
			contains: []string{"Dito", "Cloud Connection", "Oracle NoSQL Cloud Service"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderCloudConfigScreen(tt.vm)

			// Check that all expected strings are present
			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("RenderCloudConfigScreen() should contain %q", substr)
				}
			}

			// Check that result is not empty
			if result == "" {
				t.Error("RenderCloudConfigScreen() should not return empty string")
			}

			// Check that border characters are present
			if !strings.Contains(result, "╭") || !strings.Contains(result, "╮") ||
				!strings.Contains(result, "╰") || !strings.Contains(result, "╯") {
				t.Error("RenderCloudConfigScreen() should contain border characters")
			}
		})
	}
}

func TestBuildConnectionStatusMessage(t *testing.T) {
	tests := []struct {
		name          string
		status        app.ConnectionStatus
		serverVersion string
		errorMsg      string
		width         int
		contains      []string
		notContains   []string
	}{
		{
			name:          "status connecting",
			status:        app.StatusConnecting,
			serverVersion: "",
			errorMsg:      "",
			width:         100,
			contains:      []string{"Connecting"},
			notContains:   []string{},
		},
		{
			name:          "status connected with version",
			status:        app.StatusConnected,
			serverVersion: "Oracle NoSQL Database 23.1",
			errorMsg:      "",
			width:         100,
			contains:      []string{"Oracle NoSQL Database 23.1"},
			notContains:   []string{"Connected"},
		},
		{
			name:          "status connected without version",
			status:        app.StatusConnected,
			serverVersion: "",
			errorMsg:      "",
			width:         100,
			contains:      []string{"Connected"},
			notContains:   []string{},
		},
		{
			name:          "status error with message",
			status:        app.StatusError,
			serverVersion: "",
			errorMsg:      "Connection refused",
			width:         100,
			contains:      []string{"Connection refused"},
			notContains:   []string{"Connection failed"},
		},
		{
			name:          "status error without message",
			status:        app.StatusError,
			serverVersion: "",
			errorMsg:      "",
			width:         100,
			contains:      []string{"Connection failed"},
			notContains:   []string{},
		},
		{
			name:          "status error with long message - truncated",
			status:        app.StatusError,
			serverVersion: "",
			errorMsg:      strings.Repeat("a", 200),
			width:         50,
			contains:      []string{"..."},
			notContains:   []string{},
		},
		{
			name:          "status disconnected - empty message",
			status:        app.StatusDisconnected,
			serverVersion: "",
			errorMsg:      "",
			width:         100,
			contains:      []string{},
			notContains:   []string{"Connected", "Connecting", "Connection failed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildConnectionStatusMessage(tt.status, tt.serverVersion, tt.errorMsg, tt.width)

			// Check that all expected strings are present
			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("buildConnectionStatusMessage() should contain %q", substr)
				}
			}

			// Check that unexpected strings are not present
			for _, substr := range tt.notContains {
				if strings.Contains(result, substr) {
					t.Errorf("buildConnectionStatusMessage() should not contain %q", substr)
				}
			}
		})
	}
}

func TestRenderWithBorder(t *testing.T) {
	tests := []struct {
		name          string
		width         int
		height        int
		content       string
		statusMessage string
		helpText      string
		contains      []string
	}{
		{
			name:          "basic render with border",
			width:         80,
			height:        25,
			content:       "Test content",
			statusMessage: "Status message",
			helpText:      "Help text",
			contains:      []string{"Dito", "Test content", "Status message", "Help text", "╭", "╮", "╰", "╯"},
		},
		{
			name:          "render with empty status",
			width:         80,
			height:        25,
			content:       "Test content",
			statusMessage: "",
			helpText:      "Help text",
			contains:      []string{"Dito", "Test content", "Help text", "╭", "╮", "╰", "╯"},
		},
		{
			name:          "render with multiline content",
			width:         100,
			height:        30,
			content:       "Line 1\nLine 2\nLine 3",
			statusMessage: "Status",
			helpText:      "Help",
			contains:      []string{"Dito", "Line 1", "Line 2", "Line 3", "Status", "Help"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderWithBorder(tt.width, tt.height, tt.content, tt.statusMessage, tt.helpText)

			// Check that all expected strings are present
			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("renderWithBorder() should contain %q", substr)
				}
			}

			// Check that result is not empty
			if result == "" {
				t.Error("renderWithBorder() should not return empty string")
			}

			// Check that all border characters are present
			if !strings.Contains(result, "╭") || !strings.Contains(result, "╮") ||
				!strings.Contains(result, "╰") || !strings.Contains(result, "╯") ||
				!strings.Contains(result, "│") {
				t.Error("renderWithBorder() should contain all border characters")
			}
		})
	}
}
