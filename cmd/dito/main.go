package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/camikura/dito/internal/app"
	"github.com/camikura/dito/internal/version"
)

// model wraps app.Model to allow methods in main package
type model struct {
	app.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Model, cmd = app.Update(m.Model, msg)
	return m, cmd
}

func (m model) View() string {
	return app.RenderView(m.Model)
}

func main() {
	showVersion := flag.Bool("version", false, "Show version information")
	flag.BoolVar(showVersion, "v", false, "Show version information (shorthand)")
	flag.Parse()

	if *showVersion {
		fmt.Printf("dito version %s\n", version.Version)
		if version.GitCommit != "unknown" {
			fmt.Printf("  commit: %s\n", version.GitCommit)
		}
		if version.BuildDate != "unknown" {
			fmt.Printf("  built:  %s\n", version.BuildDate)
		}
		return
	}

	p := tea.NewProgram(
		model{Model: app.InitialModel()},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
