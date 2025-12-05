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
	dashCount := width - ui.RuneLen(titleText) - 3
	if dashCount < 0 {
		dashCount = 0
	}
	styledTitle := titleStyle.Render(titleText)
	title := borderStyle.Render("╭─") + styledTitle + borderStyle.Render(strings.Repeat("─", dashCount) + "╮")

	leftBorder := borderStyle.Render("│")
	bottomBorder := borderStyle.Render("╰" + strings.Repeat("─", width-2) + "╯")

	contentWidth := width - 2 // Width inside borders

	// Wrap SQL text to fit content width and track cursor position
	type wrappedLine struct {
		text       string
		cursorCol  int // -1 if cursor is not on this line, >= 0 means cursor position
		selectFrom int // Start of selection in this line (-1 if none)
		selectTo   int // End of selection in this line (-1 if none)
	}

	var wrappedLines []wrappedLine
	sqlRunes := []rune(m.CurrentSQL)
	cursorLineIndex := 0 // Track which wrapped line has the cursor

	// Normalize selection range
	selectStart, selectEnd := m.SQLSelectStart, m.SQLSelectEnd
	if selectStart > selectEnd && selectStart >= 0 && selectEnd >= 0 {
		selectStart, selectEnd = selectEnd, selectStart
	}
	hasSelection := selectStart >= 0 && selectEnd >= 0 && selectStart != selectEnd

	if len(sqlRunes) == 0 {
		// Empty SQL
		if isFocused {
			wrappedLines = append(wrappedLines, wrappedLine{text: "", cursorCol: 0, selectFrom: -1, selectTo: -1})
		} else {
			wrappedLines = append(wrappedLines, wrappedLine{text: "", cursorCol: -1, selectFrom: -1, selectTo: -1})
		}
	} else {
		// Helper to calculate selection range for a line segment
		calcSelection := func(lineStart, lineEnd int) (int, int) {
			if !hasSelection {
				return -1, -1
			}
			// Line range in absolute positions: [lineStart, lineEnd)
			if selectEnd <= lineStart || selectStart >= lineEnd {
				return -1, -1 // No overlap
			}
			from := 0
			if selectStart > lineStart {
				from = selectStart - lineStart
			}
			to := lineEnd - lineStart
			if selectEnd < lineEnd {
				to = selectEnd - lineStart
			}
			return from, to
		}

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
					cursorLineIndex = len(wrappedLines)
				}
				selFrom, selTo := calcSelection(lineStart, i)
				wrappedLines = append(wrappedLines, wrappedLine{text: line, cursorCol: cursorCol, selectFrom: selFrom, selectTo: selTo})
				lineStart = i + 1
				lineWidth = 0
			} else if lineWidth+charWidth > contentWidth && lineWidth > 0 {
				// Wrap line
				line := string(sqlRunes[lineStart:i])
				cursorCol := -1
				if isFocused && m.SQLCursorPos >= lineStart && m.SQLCursorPos < i {
					cursorCol = m.SQLCursorPos - lineStart
					cursorLineIndex = len(wrappedLines)
				}
				selFrom, selTo := calcSelection(lineStart, i)
				wrappedLines = append(wrappedLines, wrappedLine{text: line, cursorCol: cursorCol, selectFrom: selFrom, selectTo: selTo})
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
			lineEnd := len(sqlRunes)
			cursorCol := -1
			if isFocused && m.SQLCursorPos >= lineStart {
				cursorCol = m.SQLCursorPos - lineStart
				cursorLineIndex = len(wrappedLines)
				// If cursor is at end and line is full width, move cursor to next line
				if cursorCol == len([]rune(line)) && lineDisplayWidth >= contentWidth {
					selFrom, selTo := calcSelection(lineStart, lineEnd)
					wrappedLines = append(wrappedLines, wrappedLine{text: line, cursorCol: -1, selectFrom: selFrom, selectTo: selTo})
					wrappedLines = append(wrappedLines, wrappedLine{text: "", cursorCol: 0, selectFrom: -1, selectTo: -1})
					cursorLineIndex = len(wrappedLines) - 1
				} else {
					selFrom, selTo := calcSelection(lineStart, lineEnd)
					wrappedLines = append(wrappedLines, wrappedLine{text: line, cursorCol: cursorCol, selectFrom: selFrom, selectTo: selTo})
				}
			} else {
				selFrom, selTo := calcSelection(lineStart, lineEnd)
				wrappedLines = append(wrappedLines, wrappedLine{text: line, cursorCol: cursorCol, selectFrom: selFrom, selectTo: selTo})
			}
		}
	}

	// Calculate scroll offset to keep cursor visible
	scrollOffset := m.SQLScrollOffset
	if cursorLineIndex < scrollOffset {
		scrollOffset = cursorLineIndex
	} else if cursorLineIndex >= scrollOffset+height {
		scrollOffset = cursorLineIndex - height + 1
	}

	// Create vertical scrollbar
	vScrollBar := ui.NewVerticalScrollBar(len(wrappedLines), height, scrollOffset, height)

	var result strings.Builder
	result.WriteString(title + "\n")

	// Render wrapped lines with scroll offset
	for i := 0; i < height; i++ {
		var lineContent string
		var lineDisplayWidth int

		lineIndex := i + scrollOffset
		if lineIndex < len(wrappedLines) {
			wl := wrappedLines[lineIndex]
			lineRunes := []rune(wl.text)

			// Check if this line has selection
			hasLineSelection := wl.selectFrom >= 0 && wl.selectTo >= 0 && wl.selectFrom < wl.selectTo

			if hasLineSelection && wl.cursorCol < 0 {
				// Selection only (no cursor on this line)
				before := string(lineRunes[:wl.selectFrom])
				selected := string(lineRunes[wl.selectFrom:wl.selectTo])
				after := ""
				if wl.selectTo < len(lineRunes) {
					after = string(lineRunes[wl.selectTo:])
				}
				lineContent = before + ui.StyleSelection.Render(selected) + after
				lineDisplayWidth = lipgloss.Width(wl.text)
			} else if wl.cursorCol >= 0 {
				// This line has the cursor (selection is shown as cursor style)
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
					textWidth := lipgloss.Width(wl.text)
					if textWidth >= contentWidth {
						// No room for cursor block, show text only (cursor handled by next line)
						lineContent = wl.text
						lineDisplayWidth = textWidth
					} else {
						lineContent = wl.text + ui.CursorNarrow.Render(" ")
						lineDisplayWidth = textWidth + 1
					}
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

		// Get right border character (with scrollbar indicator)
		rightBorderChar := vScrollBar.GetCharAt(i)
		rightBorder := borderStyle.Render(rightBorderChar)
		result.WriteString(leftBorder + lineContent + strings.Repeat(" ", paddingLen) + rightBorder + "\n")
	}
	result.WriteString(bottomBorder)

	return result.String()
}
