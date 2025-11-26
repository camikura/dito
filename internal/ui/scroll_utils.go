package ui

// ScrollStrategy defines how scrolling behavior works.
type ScrollStrategy int

const (
	// ScrollLinear scrolls the viewport linearly (cursor moves to edge before scrolling).
	// Used by Old UI.
	ScrollLinear ScrollStrategy = iota

	// ScrollCentered keeps the cursor centered in the viewport when possible.
	// Used by New UI.
	ScrollCentered
)

// ScrollState holds the current scroll state.
type ScrollState struct {
	SelectedRow   int // Current selected row index
	TotalRows     int // Total number of rows
	VisibleRows   int // Number of visible rows in the viewport
	CurrentOffset int // Current viewport offset
}

// ScrollResult holds the result of a scroll calculation.
type ScrollResult struct {
	NewSelection int // New selected row index
	NewOffset    int // New viewport offset
}

// CalculateViewportOffset calculates the viewport offset for the current selection.
// This is useful for computing the correct offset after navigation.
func CalculateViewportOffset(state ScrollState, strategy ScrollStrategy) int {
	maxOffset := state.TotalRows - state.VisibleRows
	if maxOffset < 0 {
		maxOffset = 0
	}

	switch strategy {
	case ScrollCentered:
		// Keep cursor centered when possible
		middlePosition := state.VisibleRows / 2
		if state.SelectedRow <= middlePosition {
			// Near top - no offset needed
			return 0
		}
		if state.SelectedRow >= state.TotalRows-state.VisibleRows+middlePosition {
			// Near bottom - use max offset
			return maxOffset
		}
		// In middle - center the cursor
		return state.SelectedRow - middlePosition

	case ScrollLinear:
		fallthrough
	default:
		// Linear scrolling - ensure selected row is visible
		if state.SelectedRow < state.CurrentOffset {
			// Selection is above viewport
			return state.SelectedRow
		}
		if state.SelectedRow >= state.CurrentOffset+state.VisibleRows {
			// Selection is below viewport
			return state.SelectedRow - state.VisibleRows + 1
		}
		// Selection is within viewport - keep current offset
		if state.CurrentOffset > maxOffset {
			return maxOffset
		}
		return state.CurrentOffset
	}
}

// ScrollUp calculates the new state when scrolling up by one row.
func ScrollUp(state ScrollState, strategy ScrollStrategy) ScrollResult {
	result := ScrollResult{
		NewSelection: state.SelectedRow,
		NewOffset:    state.CurrentOffset,
	}

	if state.SelectedRow <= 0 {
		return result
	}

	result.NewSelection = state.SelectedRow - 1

	// Calculate new offset based on strategy
	newState := ScrollState{
		SelectedRow:   result.NewSelection,
		TotalRows:     state.TotalRows,
		VisibleRows:   state.VisibleRows,
		CurrentOffset: state.CurrentOffset,
	}
	result.NewOffset = CalculateViewportOffset(newState, strategy)

	return result
}

// ScrollDown calculates the new state when scrolling down by one row.
func ScrollDown(state ScrollState, strategy ScrollStrategy) ScrollResult {
	result := ScrollResult{
		NewSelection: state.SelectedRow,
		NewOffset:    state.CurrentOffset,
	}

	if state.SelectedRow >= state.TotalRows-1 {
		return result
	}

	result.NewSelection = state.SelectedRow + 1

	// Calculate new offset based on strategy
	newState := ScrollState{
		SelectedRow:   result.NewSelection,
		TotalRows:     state.TotalRows,
		VisibleRows:   state.VisibleRows,
		CurrentOffset: state.CurrentOffset,
	}
	result.NewOffset = CalculateViewportOffset(newState, strategy)

	return result
}

// ScrollPageDown calculates the new state when scrolling down by a page.
func ScrollPageDown(state ScrollState, strategy ScrollStrategy) ScrollResult {
	result := ScrollResult{
		NewSelection: state.SelectedRow,
		NewOffset:    state.CurrentOffset,
	}

	// Move by half a page
	pageSize := state.VisibleRows / 2
	if pageSize < 1 {
		pageSize = 1
	}

	newSelection := state.SelectedRow + pageSize
	if newSelection >= state.TotalRows {
		newSelection = state.TotalRows - 1
	}
	if newSelection < 0 {
		newSelection = 0
	}

	result.NewSelection = newSelection

	// Calculate new offset
	newState := ScrollState{
		SelectedRow:   result.NewSelection,
		TotalRows:     state.TotalRows,
		VisibleRows:   state.VisibleRows,
		CurrentOffset: state.CurrentOffset,
	}
	result.NewOffset = CalculateViewportOffset(newState, strategy)

	return result
}

// ScrollPageUp calculates the new state when scrolling up by a page.
func ScrollPageUp(state ScrollState, strategy ScrollStrategy) ScrollResult {
	result := ScrollResult{
		NewSelection: state.SelectedRow,
		NewOffset:    state.CurrentOffset,
	}

	// Move by half a page
	pageSize := state.VisibleRows / 2
	if pageSize < 1 {
		pageSize = 1
	}

	newSelection := state.SelectedRow - pageSize
	if newSelection < 0 {
		newSelection = 0
	}

	result.NewSelection = newSelection

	// Calculate new offset
	newState := ScrollState{
		SelectedRow:   result.NewSelection,
		TotalRows:     state.TotalRows,
		VisibleRows:   state.VisibleRows,
		CurrentOffset: state.CurrentOffset,
	}
	result.NewOffset = CalculateViewportOffset(newState, strategy)

	return result
}

// ScrollToTop calculates the state for jumping to the top.
func ScrollToTop(state ScrollState) ScrollResult {
	return ScrollResult{
		NewSelection: 0,
		NewOffset:    0,
	}
}

// ScrollToBottom calculates the state for jumping to the bottom.
func ScrollToBottom(state ScrollState) ScrollResult {
	maxOffset := state.TotalRows - state.VisibleRows
	if maxOffset < 0 {
		maxOffset = 0
	}

	return ScrollResult{
		NewSelection: state.TotalRows - 1,
		NewOffset:    maxOffset,
	}
}

// CalculateMaxHorizontalOffset calculates the maximum horizontal scroll offset.
func CalculateMaxHorizontalOffset(totalWidth, viewportWidth int) int {
	maxOffset := totalWidth - viewportWidth
	if maxOffset < 0 {
		return 0
	}
	return maxOffset
}

// ScrollHorizontal calculates the new horizontal offset.
func ScrollHorizontal(currentOffset, maxOffset, step int) int {
	newOffset := currentOffset + step
	if newOffset < 0 {
		return 0
	}
	if newOffset > maxOffset {
		return maxOffset
	}
	return newOffset
}
