package ui

import (
	"testing"
)

func TestNavigateList(t *testing.T) {
	tests := []struct {
		name         string
		direction    NavigationDirection
		currentIndex int
		totalItems   int
		expectedIdx  int
		expectedChg  bool
	}{
		{
			name:         "move down",
			direction:    NavDown,
			currentIndex: 5,
			totalItems:   10,
			expectedIdx:  6,
			expectedChg:  true,
		},
		{
			name:         "move up",
			direction:    NavUp,
			currentIndex: 5,
			totalItems:   10,
			expectedIdx:  4,
			expectedChg:  true,
		},
		{
			name:         "at top, move up",
			direction:    NavUp,
			currentIndex: 0,
			totalItems:   10,
			expectedIdx:  0,
			expectedChg:  false,
		},
		{
			name:         "at bottom, move down",
			direction:    NavDown,
			currentIndex: 9,
			totalItems:   10,
			expectedIdx:  9,
			expectedChg:  false,
		},
		{
			name:         "jump to home",
			direction:    NavHome,
			currentIndex: 5,
			totalItems:   10,
			expectedIdx:  0,
			expectedChg:  true,
		},
		{
			name:         "jump to end",
			direction:    NavEnd,
			currentIndex: 5,
			totalItems:   10,
			expectedIdx:  9,
			expectedChg:  true,
		},
		{
			name:         "empty list",
			direction:    NavDown,
			currentIndex: 0,
			totalItems:   0,
			expectedIdx:  0,
			expectedChg:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NavigateList(tt.direction, tt.currentIndex, tt.totalItems)
			if result.NewIndex != tt.expectedIdx {
				t.Errorf("NavigateList().NewIndex = %d, want %d", result.NewIndex, tt.expectedIdx)
			}
			if result.Changed != tt.expectedChg {
				t.Errorf("NavigateList().Changed = %v, want %v", result.Changed, tt.expectedChg)
			}
		})
	}
}

func TestNavigateListWithPageSize(t *testing.T) {
	tests := []struct {
		name         string
		direction    NavigationDirection
		currentIndex int
		totalItems   int
		pageSize     int
		expectedIdx  int
	}{
		{
			name:         "page down",
			direction:    NavPageDown,
			currentIndex: 0,
			totalItems:   100,
			pageSize:     10,
			expectedIdx:  10,
		},
		{
			name:         "page up",
			direction:    NavPageUp,
			currentIndex: 50,
			totalItems:   100,
			pageSize:     10,
			expectedIdx:  40,
		},
		{
			name:         "page down clamp to end",
			direction:    NavPageDown,
			currentIndex: 95,
			totalItems:   100,
			pageSize:     10,
			expectedIdx:  99,
		},
		{
			name:         "page up clamp to start",
			direction:    NavPageUp,
			currentIndex: 5,
			totalItems:   100,
			pageSize:     10,
			expectedIdx:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NavigateListWithPageSize(tt.direction, tt.currentIndex, tt.totalItems, tt.pageSize)
			if result.NewIndex != tt.expectedIdx {
				t.Errorf("NavigateListWithPageSize().NewIndex = %d, want %d", result.NewIndex, tt.expectedIdx)
			}
		})
	}
}

func TestNavigateWithFetchThreshold(t *testing.T) {
	tests := []struct {
		name          string
		direction     NavigationDirection
		currentIndex  int
		totalItems    int
		threshold     int
		hasMore       bool
		expectedFetch bool
	}{
		{
			name:          "fetch needed",
			direction:     NavDown,
			currentIndex:  95,
			totalItems:    100,
			threshold:     10,
			hasMore:       true,
			expectedFetch: true,
		},
		{
			name:          "no fetch needed - not near end",
			direction:     NavDown,
			currentIndex:  50,
			totalItems:    100,
			threshold:     10,
			hasMore:       true,
			expectedFetch: false,
		},
		{
			name:          "no fetch needed - no more data",
			direction:     NavDown,
			currentIndex:  95,
			totalItems:    100,
			threshold:     10,
			hasMore:       false,
			expectedFetch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NavigateWithFetchThreshold(tt.direction, tt.currentIndex, tt.totalItems, tt.threshold, tt.hasMore)
			if result.NeedsFetch != tt.expectedFetch {
				t.Errorf("NavigateWithFetchThreshold().NeedsFetch = %v, want %v", result.NeedsFetch, tt.expectedFetch)
			}
		})
	}
}

func TestCycleNavigation(t *testing.T) {
	tests := []struct {
		name         string
		direction    NavigationDirection
		currentIndex int
		totalItems   int
		expectedIdx  int
	}{
		{
			name:         "wrap to end",
			direction:    NavUp,
			currentIndex: 0,
			totalItems:   10,
			expectedIdx:  9,
		},
		{
			name:         "wrap to start",
			direction:    NavDown,
			currentIndex: 9,
			totalItems:   10,
			expectedIdx:  0,
		},
		{
			name:         "normal up",
			direction:    NavUp,
			currentIndex: 5,
			totalItems:   10,
			expectedIdx:  4,
		},
		{
			name:         "normal down",
			direction:    NavDown,
			currentIndex: 5,
			totalItems:   10,
			expectedIdx:  6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CycleNavigation(tt.direction, tt.currentIndex, tt.totalItems)
			if result.NewIndex != tt.expectedIdx {
				t.Errorf("CycleNavigation().NewIndex = %d, want %d", result.NewIndex, tt.expectedIdx)
			}
		})
	}
}

func TestKeyToDirection(t *testing.T) {
	tests := []struct {
		key       string
		expected  NavigationDirection
		expectOk  bool
	}{
		{"up", NavUp, true},
		{"k", NavUp, true},
		{"down", NavDown, true},
		{"j", NavDown, true},
		{"left", NavLeft, true},
		{"h", NavLeft, true},
		{"right", NavRight, true},
		{"l", NavRight, true},
		{"pgup", NavPageUp, true},
		{"ctrl+u", NavPageUp, true},
		{"pgdown", NavPageDown, true},
		{"ctrl+d", NavPageDown, true},
		{"home", NavHome, true},
		{"g", NavHome, true},
		{"end", NavEnd, true},
		{"G", NavEnd, true},
		{"unknown", 0, false},
		{"", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			dir, ok := KeyToDirection(tt.key)
			if ok != tt.expectOk {
				t.Errorf("KeyToDirection(%q) ok = %v, want %v", tt.key, ok, tt.expectOk)
			}
			if ok && dir != tt.expected {
				t.Errorf("KeyToDirection(%q) = %v, want %v", tt.key, dir, tt.expected)
			}
		})
	}
}

func TestClampIndex(t *testing.T) {
	tests := []struct {
		name       string
		index      int
		totalItems int
		expected   int
	}{
		{"normal", 5, 10, 5},
		{"negative", -1, 10, 0},
		{"over max", 15, 10, 9},
		{"empty list", 5, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ClampIndex(tt.index, tt.totalItems)
			if result != tt.expected {
				t.Errorf("ClampIndex(%d, %d) = %d, want %d", tt.index, tt.totalItems, result, tt.expected)
			}
		})
	}
}
