package app

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/camikura/dito/internal/ui"
)

// RenderView renders the new UI
func RenderView(m Model) string {
	if m.Width == 0 {
		return "Loading..."
	}

	// Minimum width check to prevent crashes
	minWidth := ui.LeftPaneContentWidth + 10 // Left pane + minimum right pane
	if m.Width < minWidth {
		return "Window too narrow"
	}

	// Minimum height check
	if m.Height < 20 {
		return "Window too short"
	}

	// Layout configuration
	// Left pane renders with borders included in leftPaneContentWidth
	leftPaneContentWidth := ui.LeftPaneContentWidth
	rightPaneActualWidth := m.Width - leftPaneContentWidth

	// Render connection pane first to get its actual height
	connectionPane := renderConnectionPane(m, leftPaneContentWidth)
	connectionPaneHeight := strings.Count(connectionPane, "\n") + 1 // Count actual lines

	// Calculate pane heights based on actual connection pane height
	// This ensures heights are always correct even if connection pane height changes
	availableHeight := m.Height - 1 - connectionPaneHeight - 6
	partHeight := availableHeight / ui.PaneHeightTotalParts
	remainder := availableHeight % ui.PaneHeightTotalParts

	tablesHeight := partHeight * ui.PaneHeightTablesParts
	schemaHeight := partHeight * ui.PaneHeightSchemaParts
	sqlHeight := partHeight * ui.PaneHeightSQLParts

	// Distribute remainder
	ui.DistributeSpace(remainder, &tablesHeight, &schemaHeight, &sqlHeight)

	// Ensure minimum heights
	if tablesHeight < 3 {
		tablesHeight = 3
	}
	if schemaHeight < 3 {
		schemaHeight = 3
	}
	if sqlHeight < 2 {
		sqlHeight = 2
	}

	// After applying minimum heights, redistribute unused space
	usedHeight := tablesHeight + schemaHeight + sqlHeight
	if usedHeight < availableHeight {
		ui.DistributeSpace(availableHeight-usedHeight, &tablesHeight, &schemaHeight, &sqlHeight)
	}

	// Render remaining panes with calculated heights
	tablesPane := renderTablesPaneWithHeight(m, leftPaneContentWidth, tablesHeight)
	schemaPane := renderSchemaPaneWithHeight(m, leftPaneContentWidth, schemaHeight)
	sqlPane := renderSQLPaneWithHeight(m, leftPaneContentWidth, sqlHeight)
	dataPane := renderDataPane(m, rightPaneActualWidth, m.Height)

	// Join left panes vertically
	leftPanes := lipgloss.JoinVertical(
		lipgloss.Left,
		connectionPane,
		tablesPane,
		schemaPane,
		sqlPane,
	)

	// Join left and right panes horizontally
	panes := lipgloss.JoinHorizontal(lipgloss.Top, leftPanes, dataPane)

	// Footer
	footerHelp := getFooterHelp(m)
	footerContent := buildFooterContent(footerHelp, m.Width)

	// Assemble final output
	var result strings.Builder
	result.WriteString(panes + "\n")
	result.WriteString(footerContent)

	baseView := result.String()

	// Overlay connection dialog if visible
	if m.ConnectionDialogVisible {
		return renderConnectionDialog(m)
	}

	// Overlay record detail dialog if visible
	if m.RecordDetailVisible {
		return renderRecordDetailDialog(m)
	}

	return baseView
}

// getFooterHelp returns the footer help text based on the current pane and state
func getFooterHelp(m Model) string {
	switch m.CurrentPane {
	case FocusPaneConnection:
		if m.Connected {
			return "Switch Pane: tab | Disconnect: ctrl+d"
		}
		return "Setup: <enter>"
	case FocusPaneTables:
		return "Select: <enter>"
	case FocusPaneSQL:
		return "Execute: ctrl+r"
	case FocusPaneData:
		if m.CustomSQL {
			return "Detail: <enter> | Reset: esc"
		}
		return "Detail: <enter>"
	}
	return ""
}

// buildFooterContent builds the footer content string with proper padding
// Format: " {help} {padding} Dito "
func buildFooterContent(footerHelp string, width int) string {
	appName := "Dito"
	footerHelpWidth := lipgloss.Width(footerHelp)
	// 1 left + 1 right of help + 1 right of Dito = 3
	footerPadding := width - footerHelpWidth - len(appName) - 3
	if footerPadding < 0 {
		footerPadding = 0
	}
	return " " + footerHelp + " " + strings.Repeat(" ", footerPadding) + appName + " "
}
