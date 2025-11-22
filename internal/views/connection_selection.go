package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/camikura/dito/internal/ui"
)

// ConnectionSelectionModel represents the data needed to render the connection selection view.
type ConnectionSelectionModel struct {
	Choices []string
	Cursor  int
}

// RenderConnectionSelection renders the connection selection screen.
// This is a pure rendering function that takes model data and returns HTML.
func RenderConnectionSelection(m ConnectionSelectionModel) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		MarginBottom(1)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	var s strings.Builder
	s.WriteString(titleStyle.Render("Select Connection") + "\n")

	for i, choice := range m.Choices {
		if m.Cursor == i {
			s.WriteString(ui.StyleSelected.Render(choice) + "\n")
		} else {
			s.WriteString(normalStyle.Render(choice) + "\n")
		}
	}

	return s.String()
}
