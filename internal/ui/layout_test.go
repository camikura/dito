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

func TestCalculatePaneLayout(t *testing.T) {
	tests := []struct {
		name        string
		config      LayoutConfig
		checkLayout func(t *testing.T, layout PaneLayout)
	}{
		{
			name: "default values applied",
			config: LayoutConfig{
				TotalWidth:  100,
				TotalHeight: 50,
			},
			checkLayout: func(t *testing.T, layout PaneLayout) {
				if layout.LeftPaneContentWidth != 50 {
					t.Errorf("LeftPaneContentWidth = %d, want 50", layout.LeftPaneContentWidth)
				}
				if layout.LeftPaneActualWidth != 52 {
					t.Errorf("LeftPaneActualWidth = %d, want 52", layout.LeftPaneActualWidth)
				}
			},
		},
		{
			name: "heights in 2:2:1 ratio",
			config: LayoutConfig{
				TotalWidth:           100,
				TotalHeight:          50,
				ConnectionPaneHeight: 3,
				LeftPaneContentWidth: 50,
			},
			checkLayout: func(t *testing.T, layout PaneLayout) {
				// Tables and Schema should be roughly equal
				diff := layout.TablesHeight - layout.SchemaHeight
				if diff < -1 || diff > 1 {
					t.Errorf("Tables (%d) and Schema (%d) heights differ too much", layout.TablesHeight, layout.SchemaHeight)
				}
				// SQL should be roughly half of Tables
				expectedSQLRatio := float64(layout.SQLHeight) / float64(layout.TablesHeight)
				if expectedSQLRatio < 0.3 || expectedSQLRatio > 0.7 {
					t.Errorf("SQL height ratio = %.2f, want ~0.5", expectedSQLRatio)
				}
			},
		},
		{
			name: "minimum heights enforced",
			config: LayoutConfig{
				TotalWidth:           100,
				TotalHeight:          15, // Very small
				ConnectionPaneHeight: 3,
				MinTablesHeight:      3,
				MinSchemaHeight:      3,
				MinSQLHeight:         2,
			},
			checkLayout: func(t *testing.T, layout PaneLayout) {
				if layout.TablesHeight < 3 {
					t.Errorf("TablesHeight = %d, want >= 3", layout.TablesHeight)
				}
				if layout.SchemaHeight < 3 {
					t.Errorf("SchemaHeight = %d, want >= 3", layout.SchemaHeight)
				}
				if layout.SQLHeight < 2 {
					t.Errorf("SQLHeight = %d, want >= 2", layout.SQLHeight)
				}
			},
		},
		{
			name: "right pane width calculation",
			config: LayoutConfig{
				TotalWidth:           100,
				TotalHeight:          50,
				LeftPaneContentWidth: 40,
			},
			checkLayout: func(t *testing.T, layout PaneLayout) {
				// LeftPaneActualWidth = 40 + 2 = 42
				// RightPaneActualWidth = 100 - 42 + 1 = 59
				if layout.RightPaneActualWidth != 59 {
					t.Errorf("RightPaneActualWidth = %d, want 59", layout.RightPaneActualWidth)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			layout := CalculatePaneLayout(tt.config)
			tt.checkLayout(t, layout)
		})
	}
}

func TestCalculateSplitLayout(t *testing.T) {
	tests := []struct {
		name          string
		totalWidth    int
		totalHeight   int
		ratio         float64
		expectedLeft  int
		expectedRight int
	}{
		{
			name:          "50-50 split",
			totalWidth:    100,
			totalHeight:   50,
			ratio:         0.5,
			expectedLeft:  50,
			expectedRight: 50,
		},
		{
			name:          "30-70 split",
			totalWidth:    100,
			totalHeight:   50,
			ratio:         0.3,
			expectedLeft:  30,
			expectedRight: 70,
		},
		{
			name:          "clamp negative ratio",
			totalWidth:    100,
			totalHeight:   50,
			ratio:         -0.5,
			expectedLeft:  0,
			expectedRight: 100,
		},
		{
			name:          "clamp ratio over 1",
			totalWidth:    100,
			totalHeight:   50,
			ratio:         1.5,
			expectedLeft:  100,
			expectedRight: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateSplitLayout(tt.totalWidth, tt.totalHeight, tt.ratio)
			if result.LeftWidth != tt.expectedLeft {
				t.Errorf("LeftWidth = %d, want %d", result.LeftWidth, tt.expectedLeft)
			}
			if result.RightWidth != tt.expectedRight {
				t.Errorf("RightWidth = %d, want %d", result.RightWidth, tt.expectedRight)
			}
			if result.Height != tt.totalHeight {
				t.Errorf("Height = %d, want %d", result.Height, tt.totalHeight)
			}
		})
	}
}

func TestContentDimensions(t *testing.T) {
	tests := []struct {
		name           string
		totalWidth     int
		totalHeight    int
		borderWidth    int
		padding        int
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "standard box",
			totalWidth:     50,
			totalHeight:    20,
			borderWidth:    1,
			padding:        1,
			expectedWidth:  46, // 50 - 2*1 - 2*1
			expectedHeight: 16, // 20 - 2*1 - 2*1
		},
		{
			name:           "no padding",
			totalWidth:     50,
			totalHeight:    20,
			borderWidth:    1,
			padding:        0,
			expectedWidth:  48, // 50 - 2*1
			expectedHeight: 18, // 20 - 2*1
		},
		{
			name:           "clamp to zero",
			totalWidth:     4,
			totalHeight:    4,
			borderWidth:    2,
			padding:        1,
			expectedWidth:  0,
			expectedHeight: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, h := ContentDimensions(tt.totalWidth, tt.totalHeight, tt.borderWidth, tt.padding)
			if w != tt.expectedWidth {
				t.Errorf("ContentDimensions() width = %d, want %d", w, tt.expectedWidth)
			}
			if h != tt.expectedHeight {
				t.Errorf("ContentDimensions() height = %d, want %d", h, tt.expectedHeight)
			}
		})
	}
}

func TestCenterPosition(t *testing.T) {
	tests := []struct {
		name          string
		containerSize int
		itemSize      int
		expected      int
	}{
		{"centered", 100, 20, 40},
		{"item larger", 20, 100, 0},
		{"exact fit", 50, 50, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CenterPosition(tt.containerSize, tt.itemSize)
			if result != tt.expected {
				t.Errorf("CenterPosition(%d, %d) = %d, want %d", tt.containerSize, tt.itemSize, result, tt.expected)
			}
		})
	}
}

func TestDistributeSpace(t *testing.T) {
	tests := []struct {
		name     string
		amount   int
		initial  []int
		expected []int
	}{
		{
			name:     "distribute 5 to 3 targets",
			amount:   5,
			initial:  []int{0, 0, 0},
			expected: []int{2, 2, 1},
		},
		{
			name:     "distribute 3 to 3 targets",
			amount:   3,
			initial:  []int{0, 0, 0},
			expected: []int{1, 1, 1},
		},
		{
			name:     "distribute 0",
			amount:   0,
			initial:  []int{5, 5, 5},
			expected: []int{5, 5, 5},
		},
		{
			name:     "distribute 1",
			amount:   1,
			initial:  []int{0, 0, 0},
			expected: []int{1, 0, 0},
		},
		{
			name:     "distribute 7 to 3 targets",
			amount:   7,
			initial:  []int{10, 10, 5},
			expected: []int{13, 12, 7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := make([]int, len(tt.initial))
			copy(values, tt.initial)

			ptrs := make([]*int, len(values))
			for i := range values {
				ptrs[i] = &values[i]
			}

			DistributeSpace(tt.amount, ptrs...)

			for i, v := range values {
				if v != tt.expected[i] {
					t.Errorf("DistributeSpace() target[%d] = %d, want %d", i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestDistributeSpaceEmptyTargets(t *testing.T) {
	// Should not panic with empty targets
	DistributeSpace(5)
}
