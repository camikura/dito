package ui

import (
	"strings"
	"testing"
)

// =====================
// Horizontal ScrollBar Tests
// =====================

func TestScrollBar_NeedsScrollBar(t *testing.T) {
	tests := []struct {
		name          string
		totalWidth    int
		viewportWidth int
		expected      bool
	}{
		{
			name:          "needs scrollbar when content wider",
			totalWidth:    100,
			viewportWidth: 50,
			expected:      true,
		},
		{
			name:          "no scrollbar when content fits",
			totalWidth:    50,
			viewportWidth: 100,
			expected:      false,
		},
		{
			name:          "no scrollbar when exact fit",
			totalWidth:    50,
			viewportWidth: 50,
			expected:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewScrollBar(tt.totalWidth, tt.viewportWidth, 0, 50)
			if got := sb.NeedsScrollBar(); got != tt.expected {
				t.Errorf("NeedsScrollBar() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestScrollBar_Render_NoScroll(t *testing.T) {
	// When no scrolling is needed, should return normal line
	sb := NewScrollBar(50, 100, 0, 20)
	result := sb.Render()

	// Should be all normal characters
	expected := strings.Repeat(ScrollBarNormal, 20)
	if result != expected {
		t.Errorf("Render() = %q, expected %q", result, expected)
	}
}

func TestScrollBar_Render_AtStart(t *testing.T) {
	// At the start (offset=0), thumb should be at the left
	sb := NewScrollBar(100, 50, 0, 20)
	result := sb.Render()

	// Should start with thumb
	if !strings.HasPrefix(result, ScrollBarThumb) {
		t.Errorf("At start, scrollbar should begin with thumb: %q", result)
	}

	// Should end with normal line
	if !strings.HasSuffix(result, ScrollBarNormal) {
		t.Errorf("At start, scrollbar should end with normal: %q", result)
	}

	// Check length
	if len([]rune(result)) != 20 {
		t.Errorf("Render() length = %d, expected 20", len([]rune(result)))
	}
}

func TestScrollBar_Render_AtEnd(t *testing.T) {
	// At the end (offset=maxOffset), thumb should be at the right
	totalWidth := 100
	viewportWidth := 50
	maxOffset := totalWidth - viewportWidth // 50
	sb := NewScrollBar(totalWidth, viewportWidth, maxOffset, 20)
	result := sb.Render()

	// Should end with thumb
	if !strings.HasSuffix(result, ScrollBarThumb) {
		t.Errorf("At end, scrollbar should end with thumb: %q", result)
	}

	// Should start with normal line
	if !strings.HasPrefix(result, ScrollBarNormal) {
		t.Errorf("At end, scrollbar should start with normal: %q", result)
	}

	// Check length
	if len([]rune(result)) != 20 {
		t.Errorf("Render() length = %d, expected 20", len([]rune(result)))
	}
}

func TestScrollBar_Render_AtMiddle(t *testing.T) {
	// At middle position
	totalWidth := 100
	viewportWidth := 50
	maxOffset := totalWidth - viewportWidth // 50
	offset := maxOffset / 2                 // 25
	sb := NewScrollBar(totalWidth, viewportWidth, offset, 20)
	result := sb.Render()

	// Should have normal on both sides of thumb
	hasNormalBefore := strings.HasPrefix(result, ScrollBarNormal)
	hasNormalAfter := strings.HasSuffix(result, ScrollBarNormal)

	if !hasNormalBefore || !hasNormalAfter {
		t.Errorf("At middle, scrollbar should have normal on both sides: %q", result)
	}

	// Check length
	if len([]rune(result)) != 20 {
		t.Errorf("Render() length = %d, expected 20", len([]rune(result)))
	}
}

func TestScrollBar_Render_ThumbSize(t *testing.T) {
	// Test that thumb size is proportional to viewport/total ratio
	tests := []struct {
		name          string
		totalWidth    int
		viewportWidth int
		renderWidth   int
		minThumbSize  int
		maxThumbSize  int
	}{
		{
			name:          "large viewport ratio (50%)",
			totalWidth:    100,
			viewportWidth: 50,
			renderWidth:   20,
			minThumbSize:  10, // 50% of 20
			maxThumbSize:  10,
		},
		{
			name:          "small viewport ratio (10%)",
			totalWidth:    100,
			viewportWidth: 10,
			renderWidth:   20,
			minThumbSize:  3, // Minimum is 3
			maxThumbSize:  3,
		},
		{
			name:          "medium viewport ratio (25%)",
			totalWidth:    100,
			viewportWidth: 25,
			renderWidth:   20,
			minThumbSize:  5, // 25% of 20
			maxThumbSize:  5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewScrollBar(tt.totalWidth, tt.viewportWidth, 0, tt.renderWidth)
			result := sb.Render()

			// Count thumb characters
			thumbCount := strings.Count(result, ScrollBarThumb)

			if thumbCount < tt.minThumbSize || thumbCount > tt.maxThumbSize {
				t.Errorf("Thumb size = %d, expected between %d and %d. Result: %q",
					thumbCount, tt.minThumbSize, tt.maxThumbSize, result)
			}
		})
	}
}

func TestScrollBar_Render_ExactWidth(t *testing.T) {
	// Test that rendered scrollbar is exactly the specified width
	widths := []int{10, 20, 50, 100}

	for _, width := range widths {
		sb := NewScrollBar(200, 50, 25, width)
		result := sb.Render()
		resultLen := len([]rune(result))

		if resultLen != width {
			t.Errorf("Render() with width=%d returned length=%d", width, resultLen)
		}
	}
}

func TestScrollBar_Render_ZeroWidth(t *testing.T) {
	sb := NewScrollBar(100, 50, 0, 0)
	result := sb.Render()

	if result != "" {
		t.Errorf("Render() with width=0 should return empty string, got %q", result)
	}
}

func TestScrollBar_ThumbMovement(t *testing.T) {
	// Test that thumb moves smoothly from start to end
	totalWidth := 200
	viewportWidth := 50
	maxOffset := totalWidth - viewportWidth
	renderWidth := 40

	var lastThumbStart int = -1

	for offset := 0; offset <= maxOffset; offset += 10 {
		sb := NewScrollBar(totalWidth, viewportWidth, offset, renderWidth)
		result := sb.Render()

		// Find thumb start position
		thumbStart := strings.Index(result, ScrollBarThumb)
		if thumbStart == -1 {
			t.Errorf("No thumb found at offset=%d", offset)
			continue
		}

		// Thumb should move right as offset increases
		if lastThumbStart != -1 && thumbStart < lastThumbStart {
			t.Errorf("Thumb moved backwards: offset=%d, thumbStart=%d, lastThumbStart=%d",
				offset, thumbStart, lastThumbStart)
		}
		lastThumbStart = thumbStart
	}
}

// =====================
// Vertical ScrollBar Tests
// =====================

func TestVerticalScrollBar_NeedsScrollBar(t *testing.T) {
	tests := []struct {
		name         string
		totalRows    int
		viewportRows int
		expected     bool
	}{
		{
			name:         "needs scrollbar when more rows than viewport",
			totalRows:    100,
			viewportRows: 20,
			expected:     true,
		},
		{
			name:         "no scrollbar when rows fit",
			totalRows:    10,
			viewportRows: 20,
			expected:     false,
		},
		{
			name:         "no scrollbar when exact fit",
			totalRows:    20,
			viewportRows: 20,
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vsb := NewVerticalScrollBar(tt.totalRows, tt.viewportRows, 0, 10)
			if got := vsb.NeedsScrollBar(); got != tt.expected {
				t.Errorf("NeedsScrollBar() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestVerticalScrollBar_GetThumbRange_AtStart(t *testing.T) {
	// 100 total rows, 20 visible, at start
	vsb := NewVerticalScrollBar(100, 20, 0, 10)
	start, end := vsb.GetThumbRange()

	// Thumb should start at 0
	if start != 0 {
		t.Errorf("At start, thumb should start at 0, got %d", start)
	}

	// Thumb should have reasonable size (20% of 10 = 2, but min is 1)
	thumbSize := end - start
	if thumbSize < 1 {
		t.Errorf("Thumb size should be at least 1, got %d", thumbSize)
	}
}

func TestVerticalScrollBar_GetThumbRange_AtEnd(t *testing.T) {
	totalRows := 100
	viewportRows := 20
	maxOffset := totalRows - viewportRows // 80
	vsb := NewVerticalScrollBar(totalRows, viewportRows, maxOffset, 10)
	start, end := vsb.GetThumbRange()

	// Thumb should end at the bottom (10)
	if end != 10 {
		t.Errorf("At end, thumb should end at 10, got %d", end)
	}

	// Thumb should not start at 0
	if start == 0 {
		t.Errorf("At end, thumb should not start at 0")
	}
}

func TestVerticalScrollBar_IsThumbAt(t *testing.T) {
	// 50 total rows, 10 visible, at start, render height 10
	vsb := NewVerticalScrollBar(50, 10, 0, 10)

	// At start, thumb should be at the top
	if !vsb.IsThumbAt(0) {
		t.Error("At start, IsThumbAt(0) should be true")
	}

	// Line 9 should not have thumb (at start)
	if vsb.IsThumbAt(9) {
		t.Error("At start, IsThumbAt(9) should be false")
	}
}

func TestVerticalScrollBar_GetCharAt(t *testing.T) {
	vsb := NewVerticalScrollBar(100, 20, 0, 10)

	// At thumb position, should return thick character
	char := vsb.GetCharAt(0)
	if char != ScrollBarVThumb {
		t.Errorf("GetCharAt(0) at start should return thumb char, got %q", char)
	}

	// At non-thumb position, should return normal character
	char = vsb.GetCharAt(9)
	if char != ScrollBarVNormal {
		t.Errorf("GetCharAt(9) at start should return normal char, got %q", char)
	}
}

func TestVerticalScrollBar_NoScrollNeeded(t *testing.T) {
	// When no scroll needed, IsThumbAt should always return false
	vsb := NewVerticalScrollBar(10, 20, 0, 10)

	for i := 0; i < 10; i++ {
		if vsb.IsThumbAt(i) {
			t.Errorf("No scroll needed, IsThumbAt(%d) should be false", i)
		}
	}
}

func TestVerticalScrollBar_ThumbMovement(t *testing.T) {
	totalRows := 100
	viewportRows := 20
	maxOffset := totalRows - viewportRows
	renderHeight := 10

	var lastThumbStart int = -1

	for offset := 0; offset <= maxOffset; offset += 10 {
		vsb := NewVerticalScrollBar(totalRows, viewportRows, offset, renderHeight)
		start, _ := vsb.GetThumbRange()

		// Thumb should move down as offset increases
		if lastThumbStart != -1 && start < lastThumbStart {
			t.Errorf("Thumb moved backwards: offset=%d, start=%d, lastStart=%d",
				offset, start, lastThumbStart)
		}
		lastThumbStart = start
	}
}
