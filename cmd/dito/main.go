package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/camikura/dito/internal/app/new_ui"
)

// model wraps new_ui.Model to allow methods in main package
type model struct {
	new_ui.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Model, cmd = new_ui.Update(m.Model, msg)
	return m, cmd
}

func (m model) View() string {
	return new_ui.RenderView(m.Model)
}

func main() {
	p := tea.NewProgram(
		model{Model: new_ui.InitialModel()},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
