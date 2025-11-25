package new_ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model
func Update(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return handleKeyPress(m, msg)
	}

	return m, nil
}

func handleKeyPress(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "tab":
		m = m.NextPane()
		return m, nil

	case "shift+tab":
		m = m.PrevPane()
		return m, nil

	// Pane-specific keys will be added in later phases
	case "up", "k":
		// TODO: Phase 2 - handle navigation within focused pane
		return m, nil

	case "down", "j":
		// TODO: Phase 2 - handle navigation within focused pane
		return m, nil

	case "enter":
		// TODO: Phase 2+ - handle selection/activation
		return m, nil

	case "esc":
		// TODO: Phase 2+ - handle back/cancel
		return m, nil
	}

	return m, nil
}
