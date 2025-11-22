package ui

import (
	"strings"
	"testing"
)

func TestTextField(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		width     int
		focused   bool
		cursorPos int
		contains  []string // strings that should be in the output
	}{
		{
			name:      "simple value not focused",
			value:     "localhost",
			width:     20,
			focused:   false,
			cursorPos: 0,
			contains:  []string{"localhost", "[", "]"},
		},
		{
			name:      "simple value focused with cursor",
			value:     "localhost",
			width:     20,
			focused:   true,
			cursorPos: 5,
			contains:  []string{"local_host", "[", "]"},
		},
		{
			name:      "cursor at start",
			value:     "test",
			width:     20,
			focused:   true,
			cursorPos: 0,
			contains:  []string{"_test", "[", "]"},
		},
		{
			name:      "cursor at end",
			value:     "test",
			width:     20,
			focused:   true,
			cursorPos: 4,
			contains:  []string{"test_", "[", "]"},
		},
		{
			name:      "long value not focused truncated",
			value:     "very-long-hostname-that-exceeds-width",
			width:     10,
			focused:   false,
			cursorPos: 0,
			contains:  []string{"...", "[", "]"},
		},
		{
			name:      "long value focused with scrolling",
			value:     "very-long-hostname-that-exceeds-width",
			width:     10,
			focused:   true,
			cursorPos: 20,
			contains:  []string{"...", "[", "]"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TextField(tt.value, tt.width, tt.focused, tt.cursorPos)
			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("TextField() = %q, should contain %q", result, substr)
				}
			}
			// Check that result starts with [ and ends with ]
			if !strings.HasPrefix(result, "[") || !strings.HasSuffix(result, "]") {
				t.Errorf("TextField() = %q, should be wrapped in [ ]", result)
			}
		})
	}
}

func TestButton(t *testing.T) {
	tests := []struct {
		name    string
		label   string
		focused bool
		want    string // substring to check for
	}{
		{
			name:    "not focused",
			label:   "Connect",
			focused: false,
			want:    "Connect",
		},
		{
			name:    "focused",
			label:   "Test Connection",
			focused: true,
			want:    "Test Connection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Button(tt.label, tt.focused)
			if !strings.Contains(result, tt.want) {
				t.Errorf("Button() = %q, should contain %q", result, tt.want)
			}
		})
	}
}

func TestCheckbox(t *testing.T) {
	tests := []struct {
		name    string
		label   string
		checked bool
		focused bool
		want    string
	}{
		{
			name:    "unchecked not focused",
			label:   "HTTPS/TLS",
			checked: false,
			focused: false,
			want:    "[ ] HTTPS/TLS",
		},
		{
			name:    "checked not focused",
			label:   "HTTPS/TLS",
			checked: true,
			focused: false,
			want:    "[x] HTTPS/TLS",
		},
		{
			name:    "unchecked focused",
			label:   "Option",
			checked: false,
			focused: true,
			want:    "[ ] Option",
		},
		{
			name:    "checked focused",
			label:   "Option",
			checked: true,
			focused: true,
			want:    "[x] Option",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Checkbox(tt.label, tt.checked, tt.focused)
			if !strings.Contains(result, tt.want) {
				t.Errorf("Checkbox() = %q, should contain %q", result, tt.want)
			}
		})
	}
}

func TestRadioButton(t *testing.T) {
	tests := []struct {
		name     string
		label    string
		selected bool
		focused  bool
		want     string
	}{
		{
			name:     "not selected not focused",
			label:    "Option A",
			selected: false,
			focused:  false,
			want:     "( ) Option A",
		},
		{
			name:     "selected not focused",
			label:    "Option B",
			selected: true,
			focused:  false,
			want:     "(*) Option B",
		},
		{
			name:     "not selected focused",
			label:    "Option C",
			selected: false,
			focused:  true,
			want:     "( ) Option C",
		},
		{
			name:     "selected focused",
			label:    "Option D",
			selected: true,
			focused:  true,
			want:     "(*) Option D",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RadioButton(tt.label, tt.selected, tt.focused)
			if !strings.Contains(result, tt.want) {
				t.Errorf("RadioButton() = %q, should contain %q", result, tt.want)
			}
		})
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{
			name:   "short string unchanged",
			input:  "hello",
			maxLen: 10,
			want:   "hello",
		},
		{
			name:   "exact length unchanged",
			input:  "hello",
			maxLen: 5,
			want:   "hello",
		},
		{
			name:   "long string truncated",
			input:  "very long string",
			maxLen: 10,
			want:   "very long…",
		},
		{
			name:   "maxLen 1",
			input:  "test",
			maxLen: 1,
			want:   "…",
		},
		{
			name:   "maxLen 2",
			input:  "test",
			maxLen: 2,
			want:   "t…",
		},
		{
			name:   "empty string",
			input:  "",
			maxLen: 5,
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TruncateString(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("TruncateString(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}
