package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Grid represents a data grid component with columns, rows, and scroll support.
// It handles:
// - Column width calculation based on content
// - Horizontal scrolling (character-based, not column-based)
// - Vertical scrolling (row-based)
// - Cell truncation with ellipsis
// - Row selection highlighting
// - Numeric column right-alignment
type Grid struct {
	// Data
	Columns []GridColumn           // Column definitions
	Rows    []map[string]interface{} // Row data

	// Scroll state
	HorizontalOffset int // Character offset for horizontal scroll
	VerticalOffset   int // Row offset for vertical scroll
	SelectedRow      int // Currently selected row index (absolute)

	// Display dimensions
	Width  int // Available width for rendering
	Height int // Available height (number of rows including header)

	// Loading state
	ShowLoading bool // Show "Loading..." at the bottom when fetching more data
	HasMore     bool // Whether there are more rows to fetch

	// Focus state
	IsFocused bool // Whether the grid has focus (affects selection style)
}

// GridColumn represents a column definition.
type GridColumn struct {
	Name      string // Column name (header)
	Type      string // Column type (INTEGER, STRING, etc.)
	Width     int    // Calculated width for this column
}

// NewGrid creates a new Grid with the given columns and data.
func NewGrid(columns []string, columnTypes map[string]string, rows []map[string]interface{}) *Grid {
	g := &Grid{
		Rows: rows,
	}

	// Create column definitions
	for _, name := range columns {
		colType := ""
		if columnTypes != nil {
			colType = columnTypes[name]
		}
		g.Columns = append(g.Columns, GridColumn{
			Name: name,
			Type: colType,
		})
	}

	// Calculate column widths
	g.calculateColumnWidths()

	return g
}

// calculateColumnWidths calculates the optimal width for each column.
// Width is based on max(header length, max data length), capped at 50.
func (g *Grid) calculateColumnWidths() {
	const maxWidth = 50
	const minWidth = 3

	for i := range g.Columns {
		col := &g.Columns[i]

		// Start with header width
		width := len([]rune(col.Name))
		if width < minWidth {
			width = minWidth
		}

		// Check data widths (sample first 100 rows for performance)
		sampleSize := len(g.Rows)
		if sampleSize > 100 {
			sampleSize = 100
		}

		for j := 0; j < sampleSize; j++ {
			if val, exists := g.Rows[j][col.Name]; exists {
				valStr := FormatValue(val)
				valWidth := len([]rune(valStr))
				if valWidth > width {
					width = valWidth
				}
			}
		}

		// Cap at maximum
		if width > maxWidth {
			width = maxWidth
		}

		col.Width = width
	}
}

// TotalContentWidth returns the total width of all columns plus separators.
func (g *Grid) TotalContentWidth() int {
	total := 0
	for _, col := range g.Columns {
		total += col.Width
	}
	// Add separators (1 space between each column)
	if len(g.Columns) > 1 {
		total += len(g.Columns) - 1
	}
	return total
}

// Render renders the grid to a string.
// The output fits within the specified Width and Height.
func (g *Grid) Render() string {
	if len(g.Rows) == 0 {
		return "No data"
	}
	if g.Width <= 0 || g.Height <= 0 {
		return ""
	}

	var lines []string

	// Render header
	headerLine := g.renderHeader()
	lines = append(lines, headerLine)

	// Render separator
	separatorLine := g.renderSeparator()
	lines = append(lines, separatorLine)

	// Render data rows
	dataHeight := g.Height - 2 // Subtract header and separator
	for i := 0; i < dataHeight; i++ {
		rowIndex := g.VerticalOffset + i
		if rowIndex < len(g.Rows) {
			isSelected := rowIndex == g.SelectedRow
			rowLine := g.renderRow(g.Rows[rowIndex], isSelected)
			lines = append(lines, rowLine)
		} else if rowIndex == len(g.Rows) && g.ShowLoading && g.HasMore {
			// Show loading indicator at the position after last row
			loadingText := "Loading..."
			grayStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
			styledLoading := grayStyle.Render(loadingText)
			padding := g.Width - len(loadingText)
			if padding < 0 {
				padding = 0
			}
			lines = append(lines, styledLoading+strings.Repeat(" ", padding))
		} else {
			// Empty line to fill height
			lines = append(lines, strings.Repeat(" ", g.Width))
		}
	}

	return strings.Join(lines, "\n")
}

// renderHeader renders the header row with horizontal scrolling applied.
func (g *Grid) renderHeader() string {
	// Build full header line
	var parts []string
	for _, col := range g.Columns {
		cell := g.formatCell(col.Name, col.Width, false)
		parts = append(parts, cell)
	}
	fullLine := strings.Join(parts, " ")

	// Apply horizontal scroll and width constraint
	return g.applyHorizontalScroll(fullLine)
}

// renderSeparator renders the separator line (─── ─── ───).
func (g *Grid) renderSeparator() string {
	// Build full separator line
	var parts []string
	for _, col := range g.Columns {
		parts = append(parts, strings.Repeat("─", col.Width))
	}
	fullLine := strings.Join(parts, " ")

	// Apply horizontal scroll and width constraint
	return g.applyHorizontalScroll(fullLine)
}

// renderRow renders a single data row with optional selection highlighting.
func (g *Grid) renderRow(row map[string]interface{}, isSelected bool) string {
	// Build full row line WITHOUT styles first (for correct width calculation)
	var parts []string
	var nullPositions []nullRegion // Track null value positions

	currentPos := 0
	for _, col := range g.Columns {
		val := FormatValue(row[col.Name])
		isNull := row[col.Name] == nil
		isNumeric := isNumericType(col.Type)

		cell := g.formatCellWithAlignment(val, col.Width, isNumeric)
		cellLen := len([]rune(cell))

		// Track null positions for later styling
		if isNull {
			nullPositions = append(nullPositions, nullRegion{
				start: currentPos,
				end:   currentPos + cellLen,
			})
		}

		parts = append(parts, cell)
		currentPos += cellLen + 1 // +1 for separator space
	}

	fullLine := strings.Join(parts, " ")

	// Apply horizontal scroll and width constraint (on unstyled text)
	scrolledLine := g.applyHorizontalScroll(fullLine)

	// Apply null styling after scrolling
	if len(nullPositions) > 0 && !isSelected {
		scrolledLine = g.applyNullStyling(scrolledLine, nullPositions)
	}

	// Apply selection highlighting
	if isSelected {
		if len(nullPositions) > 0 {
			scrolledLine = g.applyNullStylingWithSelection(scrolledLine, nullPositions)
		} else {
			// Use different background color based on focus state
			bgColor := ColorPrimaryBg
			if !g.IsFocused {
				bgColor = ColorGrayLightBg
			}
			scrolledLine = lipgloss.NewStyle().
				Background(bgColor).
				Foreground(ColorWhite).
				Render(scrolledLine)
		}
	}

	return scrolledLine
}

// nullRegion represents a region in the row that contains a null value
type nullRegion struct {
	start int
	end   int
}

// applyNullStyling applies dim styling to null value regions in the visible part
func (g *Grid) applyNullStyling(line string, nullRegions []nullRegion) string {
	runes := []rune(line)
	var result strings.Builder

	i := 0
	for i < len(runes) {
		// Calculate absolute position (accounting for horizontal offset)
		absPos := g.HorizontalOffset + i

		// Check if current position is in a null region
		inNullRegion := false
		var regionEnd int
		for _, region := range nullRegions {
			if absPos >= region.start && absPos < region.end {
				inNullRegion = true
				regionEnd = region.end
				break
			}
		}

		if inNullRegion {
			// Find the end of the null region within visible part
			j := i
			for j < len(runes) && (g.HorizontalOffset+j) < regionEnd {
				j++
			}
			// Apply dim style to this segment
			segment := string(runes[i:j])
			result.WriteString(StyleDim.Render(segment))
			i = j
		} else {
			// Normal character
			result.WriteRune(runes[i])
			i++
		}
	}

	return result.String()
}

// applyNullStylingWithSelection applies styling for selected rows with null values
func (g *Grid) applyNullStylingWithSelection(line string, nullRegions []nullRegion) string {
	runes := []rune(line)
	var result strings.Builder

	// Use different background color based on focus state
	bgColor := ColorPrimaryBg
	if !g.IsFocused {
		bgColor = ColorGrayLightBg
	}
	selectedStyle := lipgloss.NewStyle().Background(bgColor).Foreground(ColorWhite)
	selectedNullStyle := lipgloss.NewStyle().Background(bgColor).Foreground(ColorGrayMid)

	i := 0
	for i < len(runes) {
		// Calculate absolute position (accounting for horizontal offset)
		absPos := g.HorizontalOffset + i

		// Check if current position is in a null region
		inNullRegion := false
		var regionEnd int
		for _, region := range nullRegions {
			if absPos >= region.start && absPos < region.end {
				inNullRegion = true
				regionEnd = region.end
				break
			}
		}

		if inNullRegion {
			// Find the end of the null region within visible part
			j := i
			for j < len(runes) && (g.HorizontalOffset+j) < regionEnd {
				j++
			}
			// Apply selected null style to this segment
			segment := string(runes[i:j])
			result.WriteString(selectedNullStyle.Render(segment))
			i = j
		} else {
			// Normal character with selection background
			result.WriteString(selectedStyle.Render(string(runes[i])))
			i++
		}
	}

	return result.String()
}

// formatCell formats a cell value to fit the specified width.
// Truncates with ellipsis if too long, pads with spaces if too short.
func (g *Grid) formatCell(value string, width int, rightAlign bool) string {
	runes := []rune(value)

	if len(runes) > width {
		// Truncate with ellipsis
		if width <= 1 {
			return "…"
		}
		return string(runes[:width-1]) + "…"
	}

	// Pad to width
	padding := width - len(runes)
	if rightAlign {
		return strings.Repeat(" ", padding) + value
	}
	return value + strings.Repeat(" ", padding)
}

// formatCellWithAlignment formats a cell with appropriate alignment.
func (g *Grid) formatCellWithAlignment(value string, width int, isNumeric bool) string {
	return g.formatCell(value, width, isNumeric)
}

// applyHorizontalScroll applies horizontal offset and ensures output is exactly Width chars.
func (g *Grid) applyHorizontalScroll(line string) string {
	runes := []rune(line)

	// Apply offset
	if g.HorizontalOffset >= len(runes) {
		return strings.Repeat(" ", g.Width)
	}

	visible := runes[g.HorizontalOffset:]

	// Truncate or pad to exact width
	if len(visible) > g.Width {
		// Truncate without ellipsis (scrollbar indicates more content)
		return string(visible[:g.Width])
	}

	// Pad to exact width
	return string(visible) + strings.Repeat(" ", g.Width-len(visible))
}

// isNumericType checks if a column type is numeric.
func isNumericType(colType string) bool {
	upper := strings.ToUpper(colType)
	return strings.Contains(upper, "INTEGER") ||
		strings.Contains(upper, "LONG") ||
		strings.Contains(upper, "DOUBLE") ||
		strings.Contains(upper, "FLOAT") ||
		strings.Contains(upper, "NUMBER")
}

// MaxHorizontalOffset returns the maximum horizontal scroll offset.
func (g *Grid) MaxHorizontalOffset() int {
	totalWidth := g.TotalContentWidth()
	if totalWidth <= g.Width {
		return 0
	}
	return totalWidth - g.Width
}

// MaxVerticalOffset returns the maximum vertical scroll offset.
func (g *Grid) MaxVerticalOffset(visibleRows int) int {
	if len(g.Rows) <= visibleRows {
		return 0
	}
	return len(g.Rows) - visibleRows
}
