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
					},
				},
			},
			contains: []string{"Dito", "On-Premise Connection", "localhost", "8080"},
		},
		// Note: Status message tests removed as status is now shown in dialogs
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
					},
				},
			},
			contains: []string{"Dito", "Cloud Connection", "us-ashburn-1"},
		},
		// Note: Status message tests removed as status is now shown in dialogs
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

// Note: buildConnectionStatusMessage and renderWithBorder tests were removed
// as these functions have been refactored. Status messages are now shown in dialogs.
