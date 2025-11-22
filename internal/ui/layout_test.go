package ui

import (
	"strings"
	"testing"
)

func TestSeparator(t *testing.T) {
	tests := []struct {
		name  string
		width int
	}{
		{"width 10", 10},
		{"width 20", 20},
		{"width 50", 50},
		{"width 1", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Separator(tt.width)
			// Check that result contains horizontal line characters
			if !strings.Contains(result, "─") {
				t.Errorf("Separator() should contain ─ character")
			}
		})
	}
}

func TestBorderedBox(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		title    string
		width    int
		contains []string // strings that should be in the output
	}{
		{
			name:     "simple content with title",
			content:  "Hello World",
			title:    "Test",
			width:    40,
			contains: []string{"╭", "╮", "╰", "╯", "│", "Test", "Hello World"},
		},
		{
			name:     "simple content without title",
			content:  "Hello World",
			title:    "",
			width:    40,
			contains: []string{"╭", "╮", "╰", "╯", "│", "Hello World"},
		},
		{
			name:     "multiline content",
			content:  "Line 1\nLine 2\nLine 3",
			title:    "",
			width:    40,
			contains: []string{"Line 1", "Line 2", "Line 3"},
		},
		{
			name:     "empty content",
			content:  "",
			title:    "",
			width:    20,
			contains: []string{"╭", "╮", "╰", "╯"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			if tt.title != "" {
				result = BorderedBox(tt.content, tt.width, tt.title)
			} else {
				result = BorderedBox(tt.content, tt.width)
			}

			// Check for required substrings
			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("BorderedBox() = %q, should contain %q", result, substr)
				}
			}

			// Check that result has top border (╭)
			if !strings.HasPrefix(result, StyleBorder.Render("╭")) {
				t.Errorf("BorderedBox() should start with top border (╭)")
			}

			// Check that result has bottom border (╯)
			if !strings.HasSuffix(result, StyleBorder.Render("╯")) {
				t.Errorf("BorderedBox() should end with bottom border (╯)")
			}
		})
	}
}

func TestBorderedBoxWithTitle(t *testing.T) {
	content := "Test content"
	title := "My Title"
	width := 40

	result := BorderedBox(content, width, title)

	// Should contain title
	if !strings.Contains(result, title) {
		t.Errorf("BorderedBox() should contain title %q", title)
	}

	// Should have at least 3 lines (top border + empty line + content + bottom border)
	lines := strings.Split(result, "\n")
	if len(lines) < 4 {
		t.Errorf("BorderedBox() with title should have at least 4 lines, got %d", len(lines))
	}

	// First line should contain title
	if !strings.Contains(lines[0], title) {
		t.Errorf("First line should contain title, got %q", lines[0])
	}
}

func TestBorderedBoxWithoutTitle(t *testing.T) {
	content := "Test content"
	width := 40

	result := BorderedBox(content, width)

	// Should not have empty line after top border
	lines := strings.Split(result, "\n")

	// Count non-empty lines
	nonEmptyCount := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmptyCount++
		}
	}

	// Without title: top border + content line + bottom border = 3 lines minimum
	if nonEmptyCount < 3 {
		t.Errorf("BorderedBox() without title should have at least 3 non-empty lines, got %d", nonEmptyCount)
	}
}

func TestBorderedBoxBorderCharacters(t *testing.T) {
	content := "test"
	result := BorderedBox(content, 20)

	// Check for all border characters
	borderChars := []string{"╭", "╮", "╰", "╯", "│"}
	for _, char := range borderChars {
		if !strings.Contains(result, char) {
			t.Errorf("BorderedBox() should contain border character %q", char)
		}
	}
}
