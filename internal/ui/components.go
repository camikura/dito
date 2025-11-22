package ui

import (
	"fmt"
)

// TextField renders a text input field with optional cursor support.
// When focused, displays a cursor at the specified position with background color highlighting.
// Returns a formatted string like "[ value__ ]" with proper width.
func TextField(value string, width int, focused bool, cursorPos int) string {
	// カーソル位置が範囲内であることを確認
	if cursorPos > len(value) {
		cursorPos = len(value)
	}
	if cursorPos < 0 {
		cursorPos = 0
	}

	var displayValue string
	if focused {
		// カーソル位置にアンダースコアを挿入
		valueWithCursor := value[:cursorPos] + "_" + value[cursorPos:]

		if len(valueWithCursor) > width {
			// カーソル位置に応じてスクロール
			// 表示可能な文字数（"..."を除く）
			visibleWidth := width - 3

			// 表示開始位置を計算
			var start int
			if cursorPos < visibleWidth {
				// カーソルが左端近くにある場合、先頭から表示
				start = 0
				displayValue = valueWithCursor[:width-3] + "..."
			} else {
				// カーソルが右側にある場合、カーソルが見えるように右側を表示
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
		// フォーカスが外れている時は先頭から表示
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
