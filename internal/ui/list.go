package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// SelectableList represents a list of items with single selection support.
type SelectableList struct {
	Title         string   // List title (e.g., "Tables (5)")
	Items         []string // List items
	SelectedIndex int      // Currently selected item index
	Focused       bool     // Whether the list is focused (affects styling)
}

// Render renders the selectable list with title and items.
// When focused, items use primary colors. When not focused, items are grayed out.
func (sl *SelectableList) Render() string {
	var result strings.Builder

	// Determine colors based on focus state
	var titleStyle, selectedStyle, normalStyle lipgloss.Style
	if sl.Focused {
		// Focused: normal colors
		titleStyle = StyleTitle
		selectedStyle = StyleSelected
		normalStyle = StyleNormal
	} else {
		// Not focused: grayed out
		titleStyle = StyleLabel
		selectedStyle = StyleLabel
		normalStyle = lipgloss.NewStyle().Foreground(ColorGrayMid)
	}

	// Render title
	if sl.Title != "" {
		result.WriteString(titleStyle.Render(sl.Title) + "\n\n")
	}

	// Render items
	for i, item := range sl.Items {
		if i == sl.SelectedIndex {
			result.WriteString(selectedStyle.Render(fmt.Sprintf("> %s", item)) + "\n")
		} else {
			result.WriteString(normalStyle.Render(fmt.Sprintf("  %s", item)) + "\n")
		}
	}

	// Remove trailing newline
	return strings.TrimSuffix(result.String(), "\n")
}
