package ui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

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
