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
		m.Window.Width = msg.Width
		m.Window.Height = msg.Height

		// Calculate pane heights using shared utility
		m.Window.TablesHeight, m.Window.SchemaHeight, m.Window.SQLHeight = calculatePaneHeights(m)

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
		m.UI.CopyMessage = ""
		return m, nil

	case clearQuitConfirmationMsg:
		m.UI.QuitConfirmation = false
		return m, nil
	}

	return m, nil
}
