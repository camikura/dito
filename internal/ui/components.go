package ui

import (
	"fmt"
)

// TextField renders a text input field with optional cursor support.
// When focused, displays a cursor at the specified position with background color highlighting.
// Returns a formatted string like "[ value__ ]" with proper width.
func TextField(value string, width int, focused bool, cursorPos int) string {
	// Ensure cursor position is within bounds
	if cursorPos > len(value) {
		cursorPos = len(value)
	}
	if cursorPos < 0 {
		cursorPos = 0
	}

	var displayValue string
	if focused {
		// Insert underscore at cursor position
		valueWithCursor := value[:cursorPos] + "_" + value[cursorPos:]

		if len(valueWithCursor) > width {
			// Scroll based on cursor position
			// Visible character count (excluding "...")
			visibleWidth := width - 3

			// Calculate display start position
			var start int
			if cursorPos < visibleWidth {
				// When cursor is near left edge, display from the beginning
				start = 0
				displayValue = valueWithCursor[:width-3] + "..."
			} else {
				// When cursor is on the right side, scroll to keep cursor visible
				start = cursorPos - visibleWidth + 1
				end := start + visibleWidth
				if end > len(valueWithCursor) {
					end = len(valueWithCursor)
				}
				displayValue = "..." + valueWithCursor[start:end]
			}
		} else {
			displayValue = valueWithCursor
		}
	} else {
		// When not focused, display from the beginning
		if len(value) > width {
			displayValue = value[:width-3] + "..."
		} else {
			displayValue = value
		}
	}

	formattedText := fmt.Sprintf("[ %-*s ]", width, displayValue)

	// Apply background color highlighting when focused
	if focused {
		return StyleSelected.Render(formattedText)
	}
	return StyleNormal.Render(formattedText)
}

// Button renders a button with focus indicator.
// When focused, uses background color highlighting.
func Button(label string, focused bool) string {
	if focused {
		return StyleSelected.Render(label)
	}
	return StyleNormal.Render(label)
}

// Checkbox renders a checkbox with label.
// Displays "[x] Label" when checked, "[ ] Label" when unchecked.
// When focused, uses background color highlighting.
func Checkbox(label string, checked bool, focused bool) string {
	checkbox := "[ ]"
	if checked {
		checkbox = "[x]"
	}
	text := checkbox + " " + label
	if focused {
		return StyleSelected.Render(text)
	}
	return StyleNormal.Render(text)
}

// RadioButton renders a radio button with label.
// Displays "(*) Label" when selected, "( ) Label" when not selected.
// When focused, uses background color highlighting.
func RadioButton(label string, selected bool, focused bool) string {
	radio := "( )"
	if selected {
		radio = "(*)"
	}
	text := radio + " " + label
	if focused {
		return StyleSelected.Render(text)
	}
	return StyleNormal.Render(text)
}

// TruncateString truncates a string to maxLen characters with an ellipsis.
// If the string is shorter than maxLen, returns it unchanged.
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 1 {
		return "…"
	}
	return s[:maxLen-1] + "…"
}
