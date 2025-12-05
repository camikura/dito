package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/camikura/dito/internal/db"
)

// Update handles messages and updates the model
func Update(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		return handleMouseClick(m, msg)

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		// Calculate pane heights using shared utility
		m.TablesHeight, m.SchemaHeight, m.SQLHeight = calculatePaneHeights(m)

		return m, nil

	case tea.KeyMsg:
		return handleKeyPress(m, msg)

	case db.ConnectionResult:
		return handleConnectionResult(m, msg)

	case db.TableListResult:
		return handleTableListResult(m, msg)

	case db.TableDetailsResult:
		return handleTableDetailsResult(m, msg)

	case db.TableDataResult:
		return handleTableDataResult(m, msg)

	case clearCopyMessageMsg:
		m.CopyMessage = ""
		return m, nil

	case clearQuitConfirmationMsg:
		m.QuitConfirmation = false
		return m, nil
	}

	return m, nil
}
