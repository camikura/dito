package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/camikura/dito/internal/db"
	"github.com/camikura/dito/internal/ui"
)

func renderDataPane(m Model, width int, totalHeight int) string {
	borderStyle := ui.StyleBorderInactive
	titleStyle := ui.StyleTitleInactive
	if m.CurrentPane == FocusPaneData {
		borderStyle = ui.StyleBorderActive
		titleStyle = ui.StyleTitleActive
	}

	// Build title with table name if available
	// For custom SQL, show extracted table name even if not in tables list
	var dataTableName string
	if m.CustomSQL && m.CurrentSQL != "" {
		// Use extracted name directly (may include child table like orders.addresses)
		dataTableName = ui.ExtractTableNameFromSQL(m.CurrentSQL)
	} else if m.SelectedTable >= 0 && m.SelectedTable < len(m.Tables) {
		dataTableName = m.Tables[m.SelectedTable]
	}

	var titleText string
	if dataTableName != "" {
		titleText = fmt.Sprintf(" Data (%s) ", dataTableName)
	} else {
		titleText = " Data "
	}
	styledTitle := titleStyle.Render(titleText)
	title := borderStyle.Render("╭─") + styledTitle + borderStyle.Render(strings.Repeat("─", width-len(titleText)-3) + "╮")

	// Data pane should match left panes height
	// Total height = totalHeight (m.Height passed in)
	// Data pane takes: totalHeight - 1 (footer)
	// Data pane structure: title(1) + content lines(?) + bottom(1)
	// So: content lines = totalHeight - 1 - 2 = totalHeight - 3
	contentLines := totalHeight - 3
	if contentLines < 5 {
		contentLines = 5
	}

	leftBorder := borderStyle.Render("│")
	rightBorder := borderStyle.Render("│")

	var result strings.Builder
	result.WriteString(title + "\n")

	// Track scroll info for bottom border
	var totalContentWidth, viewportWidth int

	// Determine which table's data to display
	// For custom SQL, use the extracted table name; otherwise use SelectedTable
	var dataLookupTableName string
	if m.CustomSQL && m.CurrentSQL != "" {
		dataLookupTableName = ui.ExtractTableNameFromSQL(m.CurrentSQL)
	} else if m.SelectedTable >= 0 && m.SelectedTable < len(m.Tables) {
		dataLookupTableName = m.Tables[m.SelectedTable]
	}

	// Prepare content
	if dataLookupTableName == "" {
		// No table selected
		for i := 0; i < contentLines; i++ {
			line := "Select a table"
			if i > 0 {
				line = ""
			}
			paddingLen := width - len(line) - 2
			if paddingLen < 0 {
				paddingLen = 0
			}
			styledLine := line
			if line != "" {
				styledLine = ui.StyleGrayText.Render(line)
			}
			result.WriteString(leftBorder + styledLine + strings.Repeat(" ", paddingLen) + rightBorder + "\n")
		}
	} else {
		data, exists := m.TableData[dataLookupTableName]
		if m.DataErrorMsg != "" {
			// Show error message
			errMsg := m.DataErrorMsg
			maxLen := width - 4
			if len(errMsg) > maxLen {
				errMsg = errMsg[:maxLen-3] + "..."
			}
			for i := 0; i < contentLines; i++ {
				line := ""
				styledLine := ""
				if i == 0 {
					line = errMsg
					styledLine = ui.StyleErrorLight.Render(line)
				}
				paddingLen := width - len(line) - 2
				if paddingLen < 0 {
					paddingLen = 0
				}
				result.WriteString(leftBorder + styledLine + strings.Repeat(" ", paddingLen) + rightBorder + "\n")
			}
		} else if !exists || data == nil {
			// No data loaded yet
			message := "No data"
			if m.LoadingData {
				message = "Loading..."
			}
			for i := 0; i < contentLines; i++ {
				line := ""
				styledLine := ""
				if i == 0 {
					line = message
					styledLine = ui.StyleGrayText.Render(line)
				}
				paddingLen := width - len(line) - 2
				if paddingLen < 0 {
					paddingLen = 0
				}
				result.WriteString(leftBorder + styledLine + strings.Repeat(" ", paddingLen) + rightBorder + "\n")
			}
		} else if len(data.Rows) == 0 {
			// No rows in result
			for i := 0; i < contentLines; i++ {
				line := ""
				styledLine := ""
				if i == 0 {
					line = "No rows"
					styledLine = ui.StyleGrayText.Render(line)
				}
				paddingLen := width - len(line) - 2
				if paddingLen < 0 {
					paddingLen = 0
				}
				result.WriteString(leftBorder + styledLine + strings.Repeat(" ", paddingLen) + rightBorder + "\n")
			}
		} else {
			// Render grid view and get scroll info
			gridContent, scrollInfo := renderGridViewWithScrollInfo(m, dataLookupTableName, data, width, contentLines, leftBorder, borderStyle)
			result.WriteString(gridContent)
			totalContentWidth = scrollInfo.totalWidth
			viewportWidth = scrollInfo.viewportWidth
		}
	}

	// Render bottom border with scrollbar
	bottomBorder := renderBottomBorderWithScrollbar(borderStyle, width, totalContentWidth, viewportWidth, m.HorizontalOffset)
	result.WriteString(bottomBorder)

	return result.String()
}

// scrollInfo holds scroll-related information for the grid
type scrollInfo struct {
	totalWidth     int
	viewportWidth  int
	totalRows      int
	viewportRows   int
	verticalOffset int
}

// renderBottomBorderWithScrollbar renders the bottom border with an integrated scrollbar
func renderBottomBorderWithScrollbar(borderStyle lipgloss.Style, width int, totalContentWidth int, viewportWidth int, offset int) string {
	// Border structure: ╰ + content + ╯
	// Content width = width - 2 (for ╰ and ╯)
	contentWidth := width - 2
	if contentWidth < 1 {
		return borderStyle.Render("╰╯")
	}

	// Create scrollbar
	scrollbar := ui.NewScrollBar(totalContentWidth, viewportWidth, offset, contentWidth)
	scrollbarLine := scrollbar.Render()

	return borderStyle.Render("╰") + borderStyle.Render(scrollbarLine) + borderStyle.Render("╯")
}

// renderGridViewWithScrollInfo renders the data grid and returns scroll information
func renderGridViewWithScrollInfo(m Model, tableName string, data *db.TableDataResult, width int, contentLines int, leftBorder string, borderStyle lipgloss.Style) (string, scrollInfo) {
	// Get column names in schema definition order
	columns := getColumnsInSchemaOrder(m, tableName, data.Rows)

	// Get column types from schema
	columnTypes := getColumnTypes(m, tableName, columns)

	// Create Grid component
	contentWidth := width - 2 // Subtract borders
	grid := ui.NewGrid(columns, columnTypes, data.Rows)
	grid.Width = contentWidth
	grid.Height = contentLines
	grid.HorizontalOffset = m.HorizontalOffset
	grid.VerticalOffset = m.ViewportOffset
	grid.SelectedRow = m.SelectedDataRow
	grid.ShowLoading = m.LoadingData
	grid.HasMore = data.HasMore

	// Calculate viewport rows (contentLines minus header and separator)
	viewportRows := contentLines - 2
	if viewportRows < 1 {
		viewportRows = 1
	}

	// Get scroll info before rendering
	info := scrollInfo{
		totalWidth:     grid.TotalContentWidth(),
		viewportWidth:  contentWidth,
		totalRows:      len(data.Rows),
		viewportRows:   viewportRows,
		verticalOffset: m.ViewportOffset,
	}

	// Create vertical scrollbar
	vScrollBar := ui.NewVerticalScrollBar(info.totalRows, info.viewportRows, info.verticalOffset, contentLines)

	// Render grid content
	gridContent := grid.Render()

	// Add borders to each line with vertical scrollbar on right
	var result strings.Builder
	lines := strings.Split(gridContent, "\n")
	for i := 0; i < contentLines; i++ {
		// Get right border character (with scrollbar indicator)
		rightBorderChar := vScrollBar.GetCharAt(i)
		rightBorder := borderStyle.Render(rightBorderChar)

		if i < len(lines) {
			result.WriteString(leftBorder + lines[i] + rightBorder + "\n")
		} else {
			// Empty line to fill height
			emptyLine := strings.Repeat(" ", contentWidth)
			result.WriteString(leftBorder + emptyLine + rightBorder + "\n")
		}
	}

	return result.String(), info
}

// getColumnTypes extracts column types from schema information
func getColumnTypes(m Model, tableName string, columns []string) map[string]string {
	// Get DDL from table details
	var ddl string
	if details, exists := m.TableDetails[tableName]; exists && details != nil && details.Schema != nil {
		ddl = details.Schema.DDL
	}

	return ui.GetColumnTypes(ddl)
}
