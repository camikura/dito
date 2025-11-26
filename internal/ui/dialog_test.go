package ui

import (
	"strings"
	"testing"
)

func TestNewDialog(t *testing.T) {
	tests := []struct {
		name          string
		config        DialogConfig
		expectedColor string
	}{
		{
			name: "default info dialog",
			config: DialogConfig{
				Title:   "Test",
				Content: "Test content",
				Type:    DialogTypeInfo,
			},
			expectedColor: string(ColorPrimary),
		},
		{
			name: "success dialog",
			config: DialogConfig{
				Title: "Success",
				Type:  DialogTypeSuccess,
			},
			expectedColor: "#00FF00",
		},
		{
			name: "error dialog",
			config: DialogConfig{
				Title: "Error",
				Type:  DialogTypeError,
			},
			expectedColor: "#FF0000",
		},
		{
			name: "custom color",
			config: DialogConfig{
				Title:       "Custom",
				BorderColor: "#AABBCC",
			},
			expectedColor: "#AABBCC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDialog(tt.config)
			if d.config.BorderColor != tt.expectedColor {
				t.Errorf("NewDialog() border color = %q, want %q", d.config.BorderColor, tt.expectedColor)
			}
		})
	}
}

func TestDialog_Render(t *testing.T) {
	config := DialogConfig{
		Title:    "Test Dialog",
		Content:  "This is test content",
		HelpText: "Press Enter to close",
		Width:    50,
	}

	d := NewDialog(config)
	result := d.Render()

	// Should contain borders
	if !strings.Contains(result, "╭") {
		t.Error("Result should contain top-left corner")
	}
	if !strings.Contains(result, "╯") {
		t.Error("Result should contain bottom-right corner")
	}

	// Should contain title
	if !strings.Contains(result, "Test Dialog") {
		t.Error("Result should contain title")
	}

	// Should contain content
	if !strings.Contains(result, "test content") {
		t.Error("Result should contain content")
	}

	// Should contain help text
	if !strings.Contains(result, "Press Enter") {
		t.Error("Result should contain help text")
	}
}

func TestDialog_RenderCentered(t *testing.T) {
	config := DialogConfig{
		Title: "Centered",
		Width: 30,
	}

	d := NewDialog(config)
	result := d.RenderCentered(80, 24)

	lines := strings.Split(result, "\n")
	if len(lines) < 1 {
		t.Error("RenderCentered should produce output")
	}
}

func TestWrapText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		width    int
		expected int // expected number of lines
	}{
		{
			name:     "no wrap needed",
			text:     "short text",
			width:    50,
			expected: 1,
		},
		{
			name:     "wrap needed",
			text:     "this is a longer text that needs to be wrapped",
			width:    20,
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapText(tt.text, tt.width)
			lines := strings.Split(result, "\n")
			if len(lines) != tt.expected {
				t.Errorf("WrapText() produced %d lines, want %d", len(lines), tt.expected)
			}
		})
	}
}

func TestCalculateDialogSize(t *testing.T) {
	tests := []struct {
		name                     string
		screenWidth, screenHeight int
		minWidth, maxWidth       int
		minHeight, maxHeight     int
		expectedWidth            int
		expectedHeight           int
	}{
		{
			name:           "normal size",
			screenWidth:    100,
			screenHeight:   50,
			minWidth:       40,
			maxWidth:       80,
			minHeight:      10,
			maxHeight:      30,
			expectedWidth:  80, // 100-10=90, clamped to 80
			expectedHeight: 30, // 50-10=40, clamped to 30
		},
		{
			name:           "small screen",
			screenWidth:    50,
			screenHeight:   20,
			minWidth:       40,
			maxWidth:       80,
			minHeight:      10,
			maxHeight:      30,
			expectedWidth:  40, // 50-10=40
			expectedHeight: 10, // 20-10=10
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width, height := CalculateDialogSize(
				tt.screenWidth, tt.screenHeight,
				tt.minWidth, tt.maxWidth,
				tt.minHeight, tt.maxHeight,
			)
			if width != tt.expectedWidth {
				t.Errorf("CalculateDialogSize() width = %d, want %d", width, tt.expectedWidth)
			}
			if height != tt.expectedHeight {
				t.Errorf("CalculateDialogSize() height = %d, want %d", height, tt.expectedHeight)
			}
		})
	}
}
