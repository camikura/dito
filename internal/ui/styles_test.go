package ui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestColorHexConstants(t *testing.T) {
	tests := []struct {
		name     string
		hex      string
		expected string
	}{
		{"ColorPrimaryHex", ColorPrimaryHex, "#00D9FF"},
		{"ColorInactiveHex", ColorInactiveHex, "#AAAAAA"},
		{"ColorGreenHex", ColorGreenHex, "#00FF00"},
		{"ColorLabelHex", ColorLabelHex, "#00D9FF"},
		{"ColorSecondaryHex", ColorSecondaryHex, "#C47D7D"},
		{"ColorTertiaryHex", ColorTertiaryHex, "#7AA2F7"},
		{"ColorPKHex", ColorPKHex, "#7FBA7A"},
		{"ColorIndexHex", ColorIndexHex, "#E5C07B"},
		{"ColorHelpHex", ColorHelpHex, "#888888"},
		{"ColorErrorHex", ColorErrorHex, "#FF0000"},
		{"ColorSuccessHex", ColorSuccessHex, "#00FF00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.hex != tt.expected {
				t.Errorf("%s: expected %s, got %s", tt.name, tt.expected, tt.hex)
			}
		})
	}
}

func TestColorPalette(t *testing.T) {
	tests := []struct {
		name     string
		color    lipgloss.Color
		expected string
	}{
		{"ColorPrimary", ColorPrimary, "#00D9FF"},
		{"ColorWhite", ColorWhite, "#FFFFFF"},
		{"ColorBlack", ColorBlack, "#000000"},
		{"ColorGray", ColorGray, "#888888"},
		{"ColorGrayMid", ColorGrayMid, "#666666"},
		{"ColorGrayDark", ColorGrayDark, "#555555"},
		{"ColorGrayLight", ColorGrayLight, "#CCCCCC"},
		{"ColorHeaderBg", ColorHeaderBg, "#AAAAAA"},
		{"ColorHeaderText", ColorHeaderText, "#00AA00"},
		{"ColorSuccess", ColorSuccess, "#00FF00"},
		{"ColorError", ColorError, "#FF0000"},
		{"ColorInactive", ColorInactive, "#AAAAAA"},
		{"ColorSecondary", ColorSecondary, "#C47D7D"},
		{"ColorTertiary", ColorTertiary, "#7AA2F7"},
		{"ColorPK", ColorPK, "#7FBA7A"},
		{"ColorIndex", ColorIndex, "#E5C07B"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.color) != tt.expected {
				t.Errorf("%s: expected %s, got %s", tt.name, tt.expected, tt.color)
			}
		})
	}
}

func TestStylesNotNil(t *testing.T) {
	// Verify that all style variables are initialized
	styles := []struct {
		name  string
		style lipgloss.Style
	}{
		{"StyleTitle", StyleTitle},
		{"StyleNormal", StyleNormal},
		{"StyleFocused", StyleFocused},
		{"StyleSelected", StyleSelected},
		{"StyleHeader", StyleHeader},
		{"StyleLabel", StyleLabel},
		{"StyleSuccess", StyleSuccess},
		{"StyleError", StyleError},
		{"StyleSeparator", StyleSeparator},
		{"StyleBorder", StyleBorder},
	}

	for _, s := range styles {
		t.Run(s.name, func(t *testing.T) {
			// Render an empty string to ensure the style can be applied
			result := s.style.Render("")
			_ = result // Just verify no panic occurs
		})
	}
}

func TestStyleRendering(t *testing.T) {
	// Test that styles can render text without panicking
	tests := []struct {
		name  string
		style lipgloss.Style
		text  string
	}{
		{"StyleTitle renders", StyleTitle, "Title"},
		{"StyleNormal renders", StyleNormal, "Normal"},
		{"StyleFocused renders", StyleFocused, "Focused"},
		{"StyleError renders", StyleError, "Error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.style.Render(tt.text)
			if result == "" {
				t.Errorf("%s: rendered result should not be empty", tt.name)
			}
		})
	}
}
