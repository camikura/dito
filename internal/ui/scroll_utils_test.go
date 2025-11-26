package ui

import (
	"testing"
)

func TestCalculateViewportOffset_Linear(t *testing.T) {
	tests := []struct {
		name     string
		state    ScrollState
		expected int
	}{
		{
			name: "at top",
			state: ScrollState{
				SelectedRow:   0,
				TotalRows:     100,
				VisibleRows:   10,
				CurrentOffset: 0,
			},
			expected: 0,
		},
		{
			name: "selection above viewport",
			state: ScrollState{
				SelectedRow:   5,
				TotalRows:     100,
				VisibleRows:   10,
				CurrentOffset: 10,
			},
			expected: 5,
		},
		{
			name: "selection below viewport",
			state: ScrollState{
				SelectedRow:   20,
				TotalRows:     100,
				VisibleRows:   10,
				CurrentOffset: 5,
			},
			expected: 11, // 20 - 10 + 1
		},
		{
			name: "selection within viewport",
			state: ScrollState{
				SelectedRow:   15,
				TotalRows:     100,
				VisibleRows:   10,
				CurrentOffset: 10,
			},
			expected: 10, // keep current
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateViewportOffset(tt.state, ScrollLinear)
			if result != tt.expected {
				t.Errorf("CalculateViewportOffset() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestCalculateViewportOffset_Centered(t *testing.T) {
	tests := []struct {
		name     string
		state    ScrollState
		expected int
	}{
		{
			name: "near top (no scroll)",
			state: ScrollState{
				SelectedRow:   2,
				TotalRows:     100,
				VisibleRows:   10,
				CurrentOffset: 0,
			},
			expected: 0,
		},
		{
			name: "in middle (centered)",
			state: ScrollState{
				SelectedRow:   50,
				TotalRows:     100,
				VisibleRows:   10,
				CurrentOffset: 0,
			},
			expected: 45, // 50 - 5 (middle = 10/2)
		},
		{
			name: "near bottom (max offset)",
			state: ScrollState{
				SelectedRow:   98,
				TotalRows:     100,
				VisibleRows:   10,
				CurrentOffset: 0,
			},
			expected: 90, // max offset
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateViewportOffset(tt.state, ScrollCentered)
			if result != tt.expected {
				t.Errorf("CalculateViewportOffset() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestScrollDown(t *testing.T) {
	state := ScrollState{
		SelectedRow:   5,
		TotalRows:     100,
		VisibleRows:   10,
		CurrentOffset: 0,
	}

	result := ScrollDown(state, ScrollLinear)

	if result.NewSelection != 6 {
		t.Errorf("ScrollDown().NewSelection = %d, want 6", result.NewSelection)
	}
}

func TestScrollDown_AtBottom(t *testing.T) {
	state := ScrollState{
		SelectedRow:   99,
		TotalRows:     100,
		VisibleRows:   10,
		CurrentOffset: 90,
	}

	result := ScrollDown(state, ScrollLinear)

	// Should not change
	if result.NewSelection != 99 {
		t.Errorf("ScrollDown() at bottom should not change selection, got %d", result.NewSelection)
	}
}

func TestScrollUp(t *testing.T) {
	state := ScrollState{
		SelectedRow:   5,
		TotalRows:     100,
		VisibleRows:   10,
		CurrentOffset: 0,
	}

	result := ScrollUp(state, ScrollLinear)

	if result.NewSelection != 4 {
		t.Errorf("ScrollUp().NewSelection = %d, want 4", result.NewSelection)
	}
}

func TestScrollUp_AtTop(t *testing.T) {
	state := ScrollState{
		SelectedRow:   0,
		TotalRows:     100,
		VisibleRows:   10,
		CurrentOffset: 0,
	}

	result := ScrollUp(state, ScrollLinear)

	// Should not change
	if result.NewSelection != 0 {
		t.Errorf("ScrollUp() at top should not change selection, got %d", result.NewSelection)
	}
}

func TestScrollPageDown(t *testing.T) {
	state := ScrollState{
		SelectedRow:   10,
		TotalRows:     100,
		VisibleRows:   10,
		CurrentOffset: 5,
	}

	result := ScrollPageDown(state, ScrollLinear)

	// Should move by half page (5)
	if result.NewSelection != 15 {
		t.Errorf("ScrollPageDown().NewSelection = %d, want 15", result.NewSelection)
	}
}

func TestScrollPageUp(t *testing.T) {
	state := ScrollState{
		SelectedRow:   20,
		TotalRows:     100,
		VisibleRows:   10,
		CurrentOffset: 15,
	}

	result := ScrollPageUp(state, ScrollLinear)

	// Should move by half page (5)
	if result.NewSelection != 15 {
		t.Errorf("ScrollPageUp().NewSelection = %d, want 15", result.NewSelection)
	}
}

func TestScrollToTop(t *testing.T) {
	state := ScrollState{
		SelectedRow:   50,
		TotalRows:     100,
		VisibleRows:   10,
		CurrentOffset: 45,
	}

	result := ScrollToTop(state)

	if result.NewSelection != 0 {
		t.Errorf("ScrollToTop().NewSelection = %d, want 0", result.NewSelection)
	}
	if result.NewOffset != 0 {
		t.Errorf("ScrollToTop().NewOffset = %d, want 0", result.NewOffset)
	}
}

func TestScrollToBottom(t *testing.T) {
	state := ScrollState{
		SelectedRow:   50,
		TotalRows:     100,
		VisibleRows:   10,
		CurrentOffset: 45,
	}

	result := ScrollToBottom(state)

	if result.NewSelection != 99 {
		t.Errorf("ScrollToBottom().NewSelection = %d, want 99", result.NewSelection)
	}
	if result.NewOffset != 90 {
		t.Errorf("ScrollToBottom().NewOffset = %d, want 90", result.NewOffset)
	}
}

func TestCalculateMaxHorizontalOffset(t *testing.T) {
	tests := []struct {
		name          string
		totalWidth    int
		viewportWidth int
		expected      int
	}{
		{
			name:          "content wider than viewport",
			totalWidth:    100,
			viewportWidth: 50,
			expected:      50,
		},
		{
			name:          "content fits viewport",
			totalWidth:    30,
			viewportWidth: 50,
			expected:      0,
		},
		{
			name:          "exact fit",
			totalWidth:    50,
			viewportWidth: 50,
			expected:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateMaxHorizontalOffset(tt.totalWidth, tt.viewportWidth)
			if result != tt.expected {
				t.Errorf("CalculateMaxHorizontalOffset() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestScrollHorizontal(t *testing.T) {
	tests := []struct {
		name          string
		currentOffset int
		maxOffset     int
		step          int
		expected      int
	}{
		{
			name:          "scroll right",
			currentOffset: 10,
			maxOffset:     50,
			step:          5,
			expected:      15,
		},
		{
			name:          "scroll left",
			currentOffset: 10,
			maxOffset:     50,
			step:          -5,
			expected:      5,
		},
		{
			name:          "clamp to max",
			currentOffset: 48,
			maxOffset:     50,
			step:          5,
			expected:      50,
		},
		{
			name:          "clamp to min",
			currentOffset: 3,
			maxOffset:     50,
			step:          -5,
			expected:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ScrollHorizontal(tt.currentOffset, tt.maxOffset, tt.step)
			if result != tt.expected {
				t.Errorf("ScrollHorizontal() = %d, want %d", result, tt.expected)
			}
		})
	}
}
