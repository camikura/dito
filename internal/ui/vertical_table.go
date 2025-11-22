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
		value := fmt.Sprintf("%v", vt.Data[key])
		// Left-align the key with padding
		label := StyleLabel.Render(fmt.Sprintf("%-*s", maxKeyWidth, key))
		result.WriteString(label + "  " + StyleNormal.Render(value) + "\n")
	}

	// Remove trailing newline
	return strings.TrimSuffix(result.String(), "\n")
}
