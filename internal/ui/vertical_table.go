package ui

import (
	"fmt"
	"strings"
)

// VerticalTable represents a vertical table view where each row shows a key-value pair.
// Used for displaying a single record's details vertically.
type VerticalTable struct {
	Data map[string]interface{} // Column name -> value
	Keys []string                // Display order of column names
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
		if len(lines) > 1 {
			// Multi-line value: indent continuation lines
			indent := strings.Repeat(" ", maxKeyWidth+2)
			for i, line := range lines {
				if i == 0 {
					result.WriteString(label + "  " + valueStyle.Render(line) + "\n")
				} else {
					result.WriteString(indent + valueStyle.Render(line) + "\n")
				}
			}
		} else {
			// Single-line value
			result.WriteString(label + "  " + valueStyle.Render(value) + "\n")
		}
	}

	// Remove trailing newline
	return strings.TrimSuffix(result.String(), "\n")
}
