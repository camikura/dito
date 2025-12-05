package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/camikura/dito/internal/db"
	"github.com/camikura/dito/internal/ui"
)

// Update handles messages and updates the model
func Update(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		return handleMouseClick(m, msg)

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		// Calculate pane heights dynamically
		// Use actual connection pane height from model, or default to 5 if not yet set
		connectionPaneHeight := m.ConnectionPaneHeight
		if connectionPaneHeight == 0 {
			connectionPaneHeight = 5 // Default for cloud connection
		}
		// Available height for Tables, Schema, and SQL content (2:2:1 ratio)
		// Total: m.Height = leftPanes + footer
		// leftPanes = Connection + Tables(+2) + Schema(+2) + SQL(+2)
		// So: availableHeight = m.Height - 1(footer) - connectionPaneHeight - 6(borders)
		availableHeight := m.Height - 1 - connectionPaneHeight - 6

		// Split available height in 2:2:1 ratio (Tables:Schema:SQL)
		partHeight := availableHeight / ui.PaneHeightTotalParts
		remainder := availableHeight % ui.PaneHeightTotalParts

		m.TablesHeight = partHeight * ui.PaneHeightTablesParts
		m.SchemaHeight = partHeight * ui.PaneHeightSchemaParts
		m.SQLHeight = partHeight * ui.PaneHeightSQLParts

		// Distribute remainder
		ui.DistributeSpace(remainder, &m.TablesHeight, &m.SchemaHeight, &m.SQLHeight)

		// Ensure minimum heights
		if m.TablesHeight < 3 {
			m.TablesHeight = 3
		}
		if m.SchemaHeight < 3 {
			m.SchemaHeight = 3
		}
		if m.SQLHeight < 2 {
			m.SQLHeight = 2
		}

		// After applying minimum heights, redistribute unused space
		usedHeight := m.TablesHeight + m.SchemaHeight + m.SQLHeight
		if usedHeight < availableHeight {
			ui.DistributeSpace(availableHeight-usedHeight, &m.TablesHeight, &m.SchemaHeight, &m.SQLHeight)
		}

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
