package ui

import (
	"fmt"
	"strings"
)

// DataGrid represents a data grid/table view with scrolling and selection support.
type DataGrid struct {
	Rows             []map[string]interface{} // Data rows to display
	Columns          []string                 // Column names in display order
	SelectedRow      int                      // Currently selected row index (absolute position)
	HorizontalOffset int                      // Horizontal scroll offset (column index)
	ViewportOffset   int                      // Vertical scroll offset (row index)
}

// Render renders the data grid with header, separator, and rows.
// Returns the rendered grid as a string.
func (dg *DataGrid) Render(maxWidth, maxHeight int) string {
	if len(dg.Rows) == 0 {
		return "No data"
	}

	var content string

	// Calculate viewport size
	// Header (1 line) + Separator (1 line) = 2 lines
	viewportSize := maxHeight - 2
	if viewportSize < 1 {
		viewportSize = 1
	}

	// Get viewport rows
	viewportRows := dg.getViewportRows(viewportSize)
	if len(viewportRows) == 0 {
		return "No data in viewport"
	}

	// Get visible columns after horizontal scrolling
	visibleColumns := dg.getVisibleColumns()
	if len(visibleColumns) == 0 {
		return "No visible columns"
	}

	// Calculate column widths
	columnWidths := dg.calculateColumnWidths(visibleColumns)

	// Calculate available width for content
	// Account for padding and cursor space
	availableWidth := maxWidth - 4
	if availableWidth < 10 {
		availableWidth = 10 // Minimum width
	}

	// Render header
	headerParts, headerWidths := dg.renderHeader(visibleColumns, columnWidths, availableWidth)
	headerLine := "　" + strings.Join(headerParts, "  ") // Use full-width space
	content += StyleNormal.Render(headerLine) + "\n"

	// Render separator
	sepParts := make([]string, len(headerWidths))
	for i, width := range headerWidths {
		sepParts[i] = strings.Repeat("─", width)
	}
	sepLine := "　" + strings.Join(sepParts, "  ")
	content += StyleLabel.Render(sepLine) + "\n"

	// Render data rows
	for i, row := range viewportRows {
		rowParts := dg.renderRow(row, visibleColumns, columnWidths, availableWidth, headerWidths)
		rowContent := strings.Join(rowParts, "  ")

		// Apply selection style
		absoluteRowIndex := dg.ViewportOffset + i
		if absoluteRowIndex == dg.SelectedRow {
			content += StyleSelected.Render("> "+rowContent) + "\n"
		} else {
			content += StyleNormal.Render("  "+rowContent) + "\n"
		}
	}

	// Remove trailing newline
	content = strings.TrimSuffix(content, "\n")

	return content
}

// calculateColumnWidths calculates the width for each column based on content.
// Maximum width is capped at 32 characters.
func (dg *DataGrid) calculateColumnWidths(columns []string) map[string]int {
	columnWidths := make(map[string]int)

	for _, colName := range columns {
		// Start with column name length
		maxWidth := len(colName)

		// Check all rows for maximum data width
		for _, row := range dg.Rows {
			if value, exists := row[colName]; exists {
				valueStr := fmt.Sprintf("%v", value)
				if len(valueStr) > maxWidth {
					maxWidth = len(valueStr)
				}
			}
		}

		// Cap at 32 characters
		if maxWidth > 32 {
			maxWidth = 32
		}

		columnWidths[colName] = maxWidth
	}

	return columnWidths
}

// getVisibleColumns returns columns after applying horizontal offset.
func (dg *DataGrid) getVisibleColumns() []string {
	if dg.HorizontalOffset >= len(dg.Columns) {
		return dg.Columns
	}
	return dg.Columns[dg.HorizontalOffset:]
}

// getViewportRows returns the rows visible in the current viewport.
func (dg *DataGrid) getViewportRows(viewportSize int) []map[string]interface{} {
	totalRows := len(dg.Rows)
	start := dg.ViewportOffset
	end := start + viewportSize

	if start >= totalRows {
		return nil
	}
	if end > totalRows {
		end = totalRows
	}

	return dg.Rows[start:end]
}

// renderHeader renders the header row with width constraints.
// Returns the header parts and their actual widths.
func (dg *DataGrid) renderHeader(columns []string, columnWidths map[string]int, availableWidth int) ([]string, []int) {
	var headerParts []string
	var headerWidths []int
	currentWidth := 0

	for _, colName := range columns {
		width := columnWidths[colName]
		truncated := TruncateString(colName, width)
		part := fmt.Sprintf("%-*s", width, truncated)

		nextWidth := currentWidth + len(part)
		if len(headerParts) > 0 {
			nextWidth += 2 // Space between columns
		}

		if nextWidth > availableWidth {
			remaining := availableWidth - currentWidth
			if len(headerParts) > 0 {
				remaining -= 2
			}
			if remaining > 0 {
				headerParts = append(headerParts, part[:remaining])
				headerWidths = append(headerWidths, remaining)
			}
			break
		}

		headerParts = append(headerParts, part)
		headerWidths = append(headerWidths, len(part))
		currentWidth = nextWidth
	}

	return headerParts, headerWidths
}

// renderRow renders a single data row with width constraints.
// Returns the row parts that fit within the available width.
func (dg *DataGrid) renderRow(row map[string]interface{}, columns []string, columnWidths map[string]int, availableWidth int, headerWidths []int) []string {
	var rowParts []string
	currentWidth := 0

	for i, colName := range columns {
		if i >= len(headerWidths) {
			break // Don't render more columns than in header
		}

		width := columnWidths[colName]
		value := fmt.Sprintf("%v", row[colName])
		truncated := TruncateString(value, width)
		part := fmt.Sprintf("%-*s", width, truncated)

		nextWidth := currentWidth + len(part)
		if len(rowParts) > 0 {
			nextWidth += 2 // Space between columns
		}

		if nextWidth > availableWidth {
			remaining := availableWidth - currentWidth
			if len(rowParts) > 0 {
				remaining -= 2
			}
			if remaining > 0 {
				rowParts = append(rowParts, part[:remaining])
			}
			break
		}

		rowParts = append(rowParts, part)
		currentWidth = nextWidth
	}

	return rowParts
}
