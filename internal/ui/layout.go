package ui

import (
	"strings"
)

// Separator renders a horizontal separator line with the specified width.
// Uses the StyleSeparator for consistent styling.
func Separator(width int) string {
	return StyleSeparator.Render(strings.Repeat("─", width))
}

// BorderedBox renders content wrapped in a border with an optional title.
// The border uses box-drawing characters (╭╮╰╯│).
// If title is non-empty, it appears in the top border with a blank line below it.
//
// Example with title:
//   ╭── Title ────╮
//   │            │
//   │  content   │
//   ╰────────────╯
//
// Example without title:
//   ╭────────────╮
//   │  content   │
//   ╰────────────╯
func BorderedBox(content string, title string, width int) string {
	var result strings.Builder

	// Top border
	if title != "" {
		// ╭── Title ─────╮
		// Calculate: "╭──" (3) + title + "╮" (1) + padding
		titleWithPadding := " " + title + " "
		remainingWidth := width - 3 - len(titleWithPadding) - 1
		if remainingWidth < 0 {
			remainingWidth = 0
		}
		topBorder := StyleBorder.Render("╭──" + titleWithPadding + strings.Repeat("─", remainingWidth) + "╮")
		result.WriteString(topBorder + "\n")

		// Empty line after title
		emptyLine := strings.Repeat(" ", width-2)
		leftBorder := StyleBorder.Render("│")
		rightBorder := StyleBorder.Render("│")
		result.WriteString(leftBorder + emptyLine + rightBorder + "\n")
	} else {
		// ╭──────────╮
		topBorder := StyleBorder.Render("╭" + strings.Repeat("─", width-2) + "╮")
		result.WriteString(topBorder + "\n")
	}

	// Content with left and right borders
	leftBorder := StyleBorder.Render("│")
	rightBorder := StyleBorder.Render("│")

	for _, line := range strings.Split(content, "\n") {
		if line != "" {
			result.WriteString(leftBorder + line + rightBorder + "\n")
		}
	}

	// Bottom border
	bottomBorder := StyleBorder.Render("╰" + strings.Repeat("─", width-2) + "╯")
	result.WriteString(bottomBorder)

	return result.String()
}
