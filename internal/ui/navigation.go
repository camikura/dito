package ui

// NavigationDirection represents the direction of navigation.
type NavigationDirection int

const (
	// NavUp represents upward navigation.
	NavUp NavigationDirection = iota
	// NavDown represents downward navigation.
	NavDown
	// NavLeft represents leftward navigation.
	NavLeft
	// NavRight represents rightward navigation.
	NavRight
	// NavPageUp represents page up navigation.
	NavPageUp
	// NavPageDown represents page down navigation.
	NavPageDown
	// NavHome represents jump to beginning.
	NavHome
	// NavEnd represents jump to end.
	NavEnd
)

// ListNavigationResult holds the result of list navigation.
type ListNavigationResult struct {
	NewIndex    int  // New selected index
	Changed     bool // Whether selection changed
	NeedsFetch  bool // Whether more data needs to be fetched
}

// NavigateList handles navigation in a selectable list.
func NavigateList(direction NavigationDirection, currentIndex, totalItems int) ListNavigationResult {
	result := ListNavigationResult{
		NewIndex: currentIndex,
		Changed:  false,
	}

	if totalItems <= 0 {
		return result
	}

	switch direction {
	case NavUp:
		if currentIndex > 0 {
			result.NewIndex = currentIndex - 1
			result.Changed = true
		}
	case NavDown:
		if currentIndex < totalItems-1 {
			result.NewIndex = currentIndex + 1
			result.Changed = true
		}
	case NavHome:
		if currentIndex != 0 {
			result.NewIndex = 0
			result.Changed = true
		}
	case NavEnd:
		lastIndex := totalItems - 1
		if currentIndex != lastIndex {
			result.NewIndex = lastIndex
			result.Changed = true
		}
	}

	return result
}

// NavigateListWithPageSize handles navigation with page size support.
func NavigateListWithPageSize(direction NavigationDirection, currentIndex, totalItems, pageSize int) ListNavigationResult {
	result := ListNavigationResult{
		NewIndex: currentIndex,
		Changed:  false,
	}

	if totalItems <= 0 {
		return result
	}

	switch direction {
	case NavUp, NavDown, NavHome, NavEnd:
		return NavigateList(direction, currentIndex, totalItems)

	case NavPageUp:
		newIndex := currentIndex - pageSize
		if newIndex < 0 {
			newIndex = 0
		}
		if newIndex != currentIndex {
			result.NewIndex = newIndex
			result.Changed = true
		}

	case NavPageDown:
		newIndex := currentIndex + pageSize
		lastIndex := totalItems - 1
		if newIndex > lastIndex {
			newIndex = lastIndex
		}
		if newIndex != currentIndex {
			result.NewIndex = newIndex
			result.Changed = true
		}
	}

	return result
}

// NavigateWithFetchThreshold navigates and signals when to fetch more data.
// threshold is the number of remaining items that triggers a fetch.
func NavigateWithFetchThreshold(direction NavigationDirection, currentIndex, totalItems, threshold int, hasMore bool) ListNavigationResult {
	result := NavigateList(direction, currentIndex, totalItems)

	// Check if we need to fetch more data
	if result.Changed && hasMore {
		remainingItems := totalItems - result.NewIndex - 1
		if remainingItems <= threshold {
			result.NeedsFetch = true
		}
	}

	return result
}

// CycleNavigation handles cyclic navigation (wraps around).
func CycleNavigation(direction NavigationDirection, currentIndex, totalItems int) ListNavigationResult {
	result := ListNavigationResult{
		NewIndex: currentIndex,
		Changed:  false,
	}

	if totalItems <= 0 {
		return result
	}

	switch direction {
	case NavUp:
		if currentIndex > 0 {
			result.NewIndex = currentIndex - 1
		} else {
			result.NewIndex = totalItems - 1
		}
		result.Changed = true
	case NavDown:
		if currentIndex < totalItems-1 {
			result.NewIndex = currentIndex + 1
		} else {
			result.NewIndex = 0
		}
		result.Changed = true
	case NavHome:
		result.NewIndex = 0
		result.Changed = currentIndex != 0
	case NavEnd:
		result.NewIndex = totalItems - 1
		result.Changed = currentIndex != totalItems-1
	}

	return result
}

// KeyToDirection maps common key strings to navigation directions.
// Returns the direction and whether the key was recognized.
func KeyToDirection(key string) (NavigationDirection, bool) {
	switch key {
	case "up", "k":
		return NavUp, true
	case "down", "j":
		return NavDown, true
	case "left", "h":
		return NavLeft, true
	case "right", "l":
		return NavRight, true
	case "pgup", "ctrl+u":
		return NavPageUp, true
	case "pgdown", "ctrl+d":
		return NavPageDown, true
	case "home", "g":
		return NavHome, true
	case "end", "G":
		return NavEnd, true
	default:
		return 0, false
	}
}

// ClampIndex ensures an index is within valid bounds.
func ClampIndex(index, totalItems int) int {
	if index < 0 {
		return 0
	}
	if totalItems <= 0 {
		return 0
	}
	if index >= totalItems {
		return totalItems - 1
	}
	return index
}
