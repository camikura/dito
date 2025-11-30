package app

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/camikura/dito/internal/ui"
)

func renderSQLPane(m Model, width int) string {
	return renderSQLPaneWithHeight(m, width, 6)
}

func renderSQLPaneWithHeight(m Model, width int, height int) string {
	isFocused := m.CurrentPane == FocusPaneSQL
	borderStyle := ui.StyleBorderInactive
	titleStyle := ui.StyleTitleInactive
	if isFocused {
		borderStyle = ui.StyleBorderActive
		titleStyle = ui.StyleTitleActive
	}

	// Add [Custom] label if custom SQL is active
	titleText := " SQL "
	if m.CustomSQL {
		titleText = " SQL [Custom] "
	}
	styledTitle := titleStyle.Render(titleText)
	title := borderStyle.Render("╭─") + styledTitle + borderStyle.Render(strings.Repeat("─", width-len(titleText)-3) + "╮")

	leftBorder := borderStyle.Render("│")
	rightBorder := borderStyle.Render("│")
	bottomBorder := borderStyle.Render("╰" + strings.Repeat("─", width-2) + "╯")

	contentWidth := width - 2 // Width inside borders

	// Wrap SQL text to fit content width and track cursor position
	type wrappedLine struct {
		text      string
		cursorCol int // -1 if cursor is not on this line, >= 0 means cursor position
	}

	var wrappedLines []wrappedLine
	sqlRunes := []rune(m.CurrentSQL)

	if len(sqlRunes) == 0 {
		// Empty SQL
		if isFocused {
			wrappedLines = append(wrappedLines, wrappedLine{text: "", cursorCol: 0})
		} else {
			wrappedLines = append(wrappedLines, wrappedLine{text: "", cursorCol: -1})
		}
	} else {
		// Wrap text
		lineStart := 0
		lineWidth := 0
		for i, r := range sqlRunes {
			charWidth := lipgloss.Width(string(r))

			if r == '\n' {
				// End of logical line
				line := string(sqlRunes[lineStart:i])
				cursorCol := -1
				if isFocused && m.SQLCursorPos >= lineStart && m.SQLCursorPos <= i {
					cursorCol = m.SQLCursorPos - lineStart
				}
				wrappedLines = append(wrappedLines, wrappedLine{text: line, cursorCol: cursorCol})
				lineStart = i + 1
				lineWidth = 0
			} else if lineWidth+charWidth > contentWidth && lineWidth > 0 {
				// Wrap line
				line := string(sqlRunes[lineStart:i])
				cursorCol := -1
				if isFocused && m.SQLCursorPos >= lineStart && m.SQLCursorPos < i {
					cursorCol = m.SQLCursorPos - lineStart
				}
				wrappedLines = append(wrappedLines, wrappedLine{text: line, cursorCol: cursorCol})
				lineStart = i
				lineWidth = charWidth
			} else {
				lineWidth += charWidth
			}
		}

		// Add remaining text
		if lineStart <= len(sqlRunes) {
			line := string(sqlRunes[lineStart:])
			lineDisplayWidth := lipgloss.Width(line)
			cursorCol := -1
			if isFocused && m.SQLCursorPos >= lineStart {
				cursorCol = m.SQLCursorPos - lineStart
				// If cursor is at end and line is full width, move cursor to next line
				if cursorCol == len([]rune(line)) && lineDisplayWidth >= contentWidth {
					wrappedLines = append(wrappedLines, wrappedLine{text: line, cursorCol: -1})
					wrappedLines = append(wrappedLines, wrappedLine{text: "", cursorCol: 0})
				} else {
					wrappedLines = append(wrappedLines, wrappedLine{text: line, cursorCol: cursorCol})
				}
			} else {
				wrappedLines = append(wrappedLines, wrappedLine{text: line, cursorCol: cursorCol})
			}
		}
	}

	var result strings.Builder
	result.WriteString(title + "\n")

	// Render wrapped lines
	for i := 0; i < height; i++ {
		var lineContent string
		var lineDisplayWidth int

		if i < len(wrappedLines) {
			wl := wrappedLines[i]
			lineRunes := []rune(wl.text)

			if wl.cursorCol >= 0 {
				// This line has the cursor
				if wl.cursorCol < len(lineRunes) {
					beforeCursor := string(lineRunes[:wl.cursorCol])
					cursorChar := string(lineRunes[wl.cursorCol])
					afterCursor := string(lineRunes[wl.cursorCol+1:])

					var cursorBlock string
					if lipgloss.Width(cursorChar) > 1 {
						cursorBlock = ui.CursorWide.Render(cursorChar)
					} else {
						cursorBlock = ui.CursorNarrow.Render(cursorChar)
					}
					lineContent = beforeCursor + cursorBlock + afterCursor
					lineDisplayWidth = lipgloss.Width(wl.text)
				} else {
					// Cursor at end of line
					lineContent = wl.text + ui.CursorNarrow.Render(" ")
					lineDisplayWidth = lipgloss.Width(wl.text) + 1
				}
			} else {
				lineContent = wl.text
				lineDisplayWidth = lipgloss.Width(wl.text)
			}
		}

		paddingLen := contentWidth - lineDisplayWidth
		if paddingLen < 0 {
			paddingLen = 0
		}
		result.WriteString(leftBorder + lineContent + strings.Repeat(" ", paddingLen) + rightBorder + "\n")
	}
	result.WriteString(bottomBorder)

	return result.String()
}
