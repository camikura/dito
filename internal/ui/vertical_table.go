package ui

import (
	"fmt"
	"strings"
)

// VerticalTable represents a vertical table view where each row shows a key-value pair.
// Used for displaying a single record's details vertically.
type VerticalTable struct {
	Data     map[string]interface{} // Column name -> value
	Keys     []string               // Display order of column names
	MaxWidth int                    // Maximum width for wrapping (0 = no wrapping)
}

// Render renders the vertical table with each column name and value on a separate line.
// Format: "column_name  value"
func (vt *VerticalTable) Render() string {
	if len(vt.Keys) == 0 || vt.Data == nil {
		return "No data"
	}

	var result strings.Builder

	// Calculate maximum column name width
	maxKeyWidth := 0
	for _, key := range vt.Keys {
		if len(key) > maxKeyWidth {
			maxKeyWidth = len(key)
		}
	}

	// Calculate available width for value (label + 2 spaces + value)
	valueWidth := 0
	if vt.MaxWidth > 0 {
		valueWidth = vt.MaxWidth - maxKeyWidth - 2
		if valueWidth < 20 {
			valueWidth = 20 // Minimum value width
		}
	}

	// Render each key-value pair
	for _, key := range vt.Keys {
		value := FormatValuePretty(vt.Data[key])
		// Left-align the key with padding and use header style (without underline) for labels
		// This matches the grid view column headers but without underline
		labelStyle := StyleHeader.Copy().Underline(false)
		label := labelStyle.Render(fmt.Sprintf("%-*s", maxKeyWidth, key))

		// Choose style based on value (dim for null)
		valueStyle := StyleNormal
		if value == "(null)" {
			valueStyle = StyleDim
		}

		// Handle multi-line values (e.g., formatted JSON)
		// Add indentation to continuation lines to align with the value column
		lines := strings.Split(value, "\n")
		indent := strings.Repeat(" ", maxKeyWidth+2)

		for i, line := range lines {
			// Wrap long lines if MaxWidth is set
			wrappedLines := wrapText(line, valueWidth)

			for j, wrappedLine := range wrappedLines {
				if i == 0 && j == 0 {
					result.WriteString(label + "  " + valueStyle.Render(wrappedLine) + "\n")
				} else {
					result.WriteString(indent + valueStyle.Render(wrappedLine) + "\n")
				}
			}
		}
	}

	// Remove trailing newline
	return strings.TrimSuffix(result.String(), "\n")
}

// wrapText wraps text to fit within maxWidth characters.
// If maxWidth is 0 or negative, returns the original text as a single-element slice.
func wrapText(text string, maxWidth int) []string {
	if maxWidth <= 0 || len(text) <= maxWidth {
		return []string{text}
	}

	var lines []string
	runes := []rune(text)

	for len(runes) > 0 {
		if len(runes) <= maxWidth {
			lines = append(lines, string(runes))
			break
		}

		// Find a good break point (prefer space)
		breakPoint := maxWidth
		for i := maxWidth; i > maxWidth/2; i-- {
			if runes[i] == ' ' {
				breakPoint = i
				break
			}
		}

		lines = append(lines, string(runes[:breakPoint]))
		runes = runes[breakPoint:]

		// Skip leading space on next line
		if len(runes) > 0 && runes[0] == ' ' {
			runes = runes[1:]
		}
	}

	return lines
}
