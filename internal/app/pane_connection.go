package app

import (
	"strings"

	"github.com/camikura/dito/internal/ui"
)

func renderConnectionPane(m Model, width int) string {
	borderStyle := ui.StyleBorderInactive
	titleStyle := ui.StyleTitleInactive
	if m.CurrentPane == FocusPaneConnection {
		borderStyle = ui.StyleBorderActive
		titleStyle = ui.StyleTitleActive
	}

	var titleText string
	var titleDisplayWidth int
	if m.Connection.Connected {
		checkmark := ui.StyleCheckmark.Render("✓")
		titleText = titleStyle.Render(" Connection ") + checkmark + " "
		titleDisplayWidth = 14
	} else {
		titleText = titleStyle.Render(" Connection ")
		titleDisplayWidth = 12
	}

	// Title line: ╭─ + titleText + ─...─ + ╮
	// Total width = width, so dashes = width - 2(╭╮) - 1(─ after ╭) - titleDisplayWidth
	dashesLen := width - 3 - titleDisplayWidth
	if dashesLen < 0 {
		dashesLen = 0
	}
	title := borderStyle.Render("╭─") + titleText + borderStyle.Render(strings.Repeat("─", dashesLen)+"╮")

	content := "(not configured)"
	if m.Connection.Message != "" {
		// Show error message if connection failed
		content = m.Connection.Message
		if len(content) > width-4 {
			content = content[:width-7] + "..."
		}
	} else if m.Connection.Endpoint != "" {
		content = m.Connection.Endpoint
	}

	// Pad content to width (no left/right padding)
	paddingLen := width - len(content) - 2
	if paddingLen < 0 {
		paddingLen = 0
	}
	contentPadded := content + strings.Repeat(" ", paddingLen)

	leftBorder := borderStyle.Render("│")
	rightBorder := borderStyle.Render("│")
	bottomBorder := borderStyle.Render("╰" + strings.Repeat("─", width-2) + "╯")

	var result strings.Builder
	result.WriteString(title + "\n")
	result.WriteString(leftBorder + contentPadded + rightBorder + "\n")
	result.WriteString(bottomBorder)

	return result.String()
}
