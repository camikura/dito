package new_ui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestGetFooterHelp(t *testing.T) {
	tests := []struct {
		name      string
		model     Model
		expected  string
	}{
		{
			name:     "Connection pane not connected",
			model:    Model{CurrentPane: FocusPaneConnection, Connected: false},
			expected: "Setup: <enter>",
		},
		{
			name:     "Connection pane connected",
			model:    Model{CurrentPane: FocusPaneConnection, Connected: true},
			expected: "Switch Pane: tab | Disconnect: ctrl+d",
		},
		{
			name:     "Tables pane",
			model:    Model{CurrentPane: FocusPaneTables},
			expected: "Select: <enter>",
		},
		{
			name:     "SQL pane",
			model:    Model{CurrentPane: FocusPaneSQL},
			expected: "Execute: ctrl+r",
		},
		{
			name:     "Schema pane",
			model:    Model{CurrentPane: FocusPaneSchema},
			expected: "",
		},
		{
			name:     "Data pane normal",
			model:    Model{CurrentPane: FocusPaneData, CustomSQL: false},
			expected: "Detail: <enter>",
		},
		{
			name:     "Data pane custom SQL",
			model:    Model{CurrentPane: FocusPaneData, CustomSQL: true},
			expected: "Detail: <enter> | Reset: esc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getFooterHelp(tt.model)
			if result != tt.expected {
				t.Errorf("getFooterHelp() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestBuildFooterContent(t *testing.T) {
	tests := []struct {
		name       string
		footerHelp string
		width      int
	}{
		{
			name:       "Connection not connected",
			footerHelp: "Setup: <enter>",
			width:      120,
		},
		{
			name:       "Connection connected",
			footerHelp: "Switch Pane: tab | Disconnect: ctrl+d",
			width:      120,
		},
		{
			name:       "Tables pane",
			footerHelp: "Select: <enter>",
			width:      120,
		},
		{
			name:       "Data pane custom SQL",
			footerHelp: "Detail: <enter> | Reset: esc",
			width:      120,
		},
		{
			name:       "Narrow width",
			footerHelp: "Setup: <enter>",
			width:      80,
		},
		{
			name:       "Very narrow width",
			footerHelp: "Setup: <enter>",
			width:      50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildFooterContent(tt.footerHelp, tt.width)
			resultWidth := lipgloss.Width(result)

			// Footer width should match the specified width
			if resultWidth != tt.width {
				t.Errorf("buildFooterContent() width = %d, want %d", resultWidth, tt.width)
			}

			// Footer should start with space
			if len(result) > 0 && result[0] != ' ' {
				t.Errorf("buildFooterContent() should start with space, got %q", result[0:1])
			}

			// Footer should end with " Dito "
			if len(result) >= 6 {
				suffix := result[len(result)-6:]
				// Check for "Dito " (app name + trailing space)
				if suffix[len(suffix)-5:] != "Dito " {
					t.Errorf("buildFooterContent() should end with 'Dito ', got %q", suffix)
				}
			}
		})
	}
}

func TestBuildFooterContentNarrowWidth(t *testing.T) {
	// When width is too narrow, padding should be 0 (not negative)
	footerHelp := "Setup: <enter>"
	width := 15 // Too narrow for the content

	result := buildFooterContent(footerHelp, width)

	// Should not panic and should produce some output
	if len(result) == 0 {
		t.Error("buildFooterContent() should produce output even with narrow width")
	}

	// The result should be at least as wide as the minimum content
	// " " + footerHelp + " " + "Dito" + " " = 1 + len(footerHelp) + 1 + 4 + 1
	minWidth := 1 + lipgloss.Width(footerHelp) + 1 + 4 + 1
	resultWidth := lipgloss.Width(result)

	if resultWidth < minWidth {
		t.Errorf("buildFooterContent() width = %d, want at least %d", resultWidth, minWidth)
	}
}

func TestFooterWidthConsistency(t *testing.T) {
	// Test that all pane states produce consistent footer widths
	width := 120

	paneStates := []Model{
		{CurrentPane: FocusPaneConnection, Connected: false},
		{CurrentPane: FocusPaneConnection, Connected: true},
		{CurrentPane: FocusPaneTables},
		{CurrentPane: FocusPaneSQL},
		{CurrentPane: FocusPaneData, CustomSQL: false},
		{CurrentPane: FocusPaneData, CustomSQL: true},
	}

	for _, m := range paneStates {
		footerHelp := getFooterHelp(m)
		if footerHelp == "" {
			continue // Skip panes with no footer help
		}

		result := buildFooterContent(footerHelp, width)
		resultWidth := lipgloss.Width(result)

		if resultWidth != width {
			t.Errorf("Footer width for pane %v = %d, want %d", m.CurrentPane, resultWidth, width)
		}
	}
}
