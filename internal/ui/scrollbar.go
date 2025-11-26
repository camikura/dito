package ui

import (
	"strings"
)

// ScrollBar represents a horizontal scroll indicator using border line thickness.
// It renders a line where the "thumb" (current viewport) is shown with a thicker line.
//
// Example with total=100, viewport=30, offset=20, width=50:
//
//	─────────━━━━━━━━━━━━━━━─────────────────────────
//	         ↑ thumb showing current viewport position
//
// The thumb width is proportional to (viewport / total).
// The thumb position is proportional to (offset / total).
type ScrollBar struct {
	TotalWidth    int // Total content width (all columns)
	ViewportWidth int // Visible viewport width
	Offset        int // Current horizontal scroll offset
	RenderWidth   int // Width to render the scrollbar
}

// Characters for scrollbar rendering
const (
	// Horizontal scrollbar
	ScrollBarHNormal = "─" // Normal horizontal (not in viewport)
	ScrollBarHThumb  = "━" // Thumb horizontal (current viewport indicator)

	// Vertical scrollbar
	ScrollBarVNormal = "│" // Normal vertical (not in viewport)
	ScrollBarVThumb  = "┃" // Thumb vertical (current viewport indicator)
)

// For backward compatibility
const (
	ScrollBarNormal = ScrollBarHNormal
	ScrollBarThumb  = ScrollBarHThumb
)

// NeedsScrollBar returns true if scrolling is possible (content wider than viewport).
func (s *ScrollBar) NeedsScrollBar() bool {
	return s.TotalWidth > s.ViewportWidth
}

// Render renders the scrollbar as a string.
// If no scrolling is needed, returns a normal line.
func (s *ScrollBar) Render() string {
	if s.RenderWidth <= 0 {
		return ""
	}

	// If no scrolling needed, return normal line
	if !s.NeedsScrollBar() {
		return strings.Repeat(ScrollBarNormal, s.RenderWidth)
	}

	// Calculate thumb size (proportional to viewport/total ratio)
	// Minimum thumb size is 3 characters for visibility
	thumbRatio := float64(s.ViewportWidth) / float64(s.TotalWidth)
	thumbSize := int(float64(s.RenderWidth) * thumbRatio)
	if thumbSize < 3 {
		thumbSize = 3
	}
	if thumbSize > s.RenderWidth {
		thumbSize = s.RenderWidth
	}

	// Calculate thumb position (proportional to offset/maxOffset ratio)
	maxOffset := s.TotalWidth - s.ViewportWidth
	if maxOffset <= 0 {
		maxOffset = 1
	}

	// Available space for thumb to move
	availableSpace := s.RenderWidth - thumbSize

	// Calculate thumb start position
	var thumbStart int
	if s.Offset >= maxOffset {
		// At the end
		thumbStart = availableSpace
	} else if s.Offset <= 0 {
		// At the beginning
		thumbStart = 0
	} else {
		// Proportional position
		thumbStart = int(float64(s.Offset) / float64(maxOffset) * float64(availableSpace))
	}

	// Ensure thumb stays within bounds
	if thumbStart < 0 {
		thumbStart = 0
	}
	if thumbStart+thumbSize > s.RenderWidth {
		thumbStart = s.RenderWidth - thumbSize
	}

	// Build the scrollbar string
	var result strings.Builder

	// Left part (before thumb)
	if thumbStart > 0 {
		result.WriteString(strings.Repeat(ScrollBarNormal, thumbStart))
	}

	// Thumb
	result.WriteString(strings.Repeat(ScrollBarThumb, thumbSize))

	// Right part (after thumb)
	rightSize := s.RenderWidth - thumbStart - thumbSize
	if rightSize > 0 {
		result.WriteString(strings.Repeat(ScrollBarNormal, rightSize))
	}

	return result.String()
}

// NewScrollBar creates a new ScrollBar with the given parameters.
func NewScrollBar(totalWidth, viewportWidth, offset, renderWidth int) *ScrollBar {
	return &ScrollBar{
		TotalWidth:    totalWidth,
		ViewportWidth: viewportWidth,
		Offset:        offset,
		RenderWidth:   renderWidth,
	}
}

// VerticalScrollBar represents a vertical scroll indicator.
// It provides information about which line positions should show the thumb.
type VerticalScrollBar struct {
	TotalRows    int // Total number of rows
	ViewportRows int // Number of visible rows
	Offset       int // Current vertical scroll offset (first visible row)
	RenderHeight int // Height to render the scrollbar
}

// NewVerticalScrollBar creates a new VerticalScrollBar.
func NewVerticalScrollBar(totalRows, viewportRows, offset, renderHeight int) *VerticalScrollBar {
	return &VerticalScrollBar{
		TotalRows:    totalRows,
		ViewportRows: viewportRows,
		Offset:       offset,
		RenderHeight: renderHeight,
	}
}

// NeedsScrollBar returns true if scrolling is possible.
func (v *VerticalScrollBar) NeedsScrollBar() bool {
	return v.TotalRows > v.ViewportRows
}

// IsThumbAt returns true if the thumb should be displayed at the given line index.
func (v *VerticalScrollBar) IsThumbAt(lineIndex int) bool {
	if !v.NeedsScrollBar() || v.RenderHeight <= 0 {
		return false
	}

	thumbStart, thumbEnd := v.GetThumbRange()
	return lineIndex >= thumbStart && lineIndex < thumbEnd
}

// GetThumbRange returns the start and end indices of the thumb (end is exclusive).
func (v *VerticalScrollBar) GetThumbRange() (start, end int) {
	if !v.NeedsScrollBar() || v.RenderHeight <= 0 {
		return 0, 0
	}

	// Calculate thumb size (proportional to viewport/total ratio)
	// Minimum thumb size is 1 for visibility
	thumbRatio := float64(v.ViewportRows) / float64(v.TotalRows)
	thumbSize := int(float64(v.RenderHeight) * thumbRatio)
	if thumbSize < 1 {
		thumbSize = 1
	}
	if thumbSize > v.RenderHeight {
		thumbSize = v.RenderHeight
	}

	// Calculate thumb position
	maxOffset := v.TotalRows - v.ViewportRows
	if maxOffset <= 0 {
		return 0, thumbSize
	}

	// Available space for thumb to move
	availableSpace := v.RenderHeight - thumbSize

	// Calculate thumb start position
	var thumbStart int
	if v.Offset >= maxOffset {
		thumbStart = availableSpace
	} else if v.Offset <= 0 {
		thumbStart = 0
	} else {
		thumbStart = int(float64(v.Offset) / float64(maxOffset) * float64(availableSpace))
	}

	// Ensure thumb stays within bounds
	if thumbStart < 0 {
		thumbStart = 0
	}
	if thumbStart+thumbSize > v.RenderHeight {
		thumbStart = v.RenderHeight - thumbSize
	}

	return thumbStart, thumbStart + thumbSize
}

// GetCharAt returns the appropriate border character for the given line index.
func (v *VerticalScrollBar) GetCharAt(lineIndex int) string {
	if v.IsThumbAt(lineIndex) {
		return ScrollBarVThumb
	}
	return ScrollBarVNormal
}
