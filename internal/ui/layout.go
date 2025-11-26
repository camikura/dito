package ui

import (
	"strings"
)

// PaneLayout holds the calculated layout dimensions for a multi-pane UI.
type PaneLayout struct {
	TotalWidth  int
	TotalHeight int

	// Left pane dimensions
	LeftPaneContentWidth int
	LeftPaneActualWidth  int

	// Right pane dimensions
	RightPaneActualWidth int

	// Individual pane heights
	ConnectionHeight int
	TablesHeight     int
	SchemaHeight     int
	SQLHeight        int

	// Footer
	FooterHeight int
}

// LayoutConfig holds configuration for layout calculation.
type LayoutConfig struct {
	// TotalWidth is the total available width.
	TotalWidth int
	// TotalHeight is the total available height.
	TotalHeight int
	// ConnectionPaneHeight is the height of the connection pane (can vary).
	ConnectionPaneHeight int
	// LeftPaneContentWidth is the content width for left panes (default: 50).
	LeftPaneContentWidth int
	// MinTablesHeight is the minimum height for tables pane (default: 3).
	MinTablesHeight int
	// MinSchemaHeight is the minimum height for schema pane (default: 3).
	MinSchemaHeight int
	// MinSQLHeight is the minimum height for SQL pane (default: 2).
	MinSQLHeight int
	// FooterHeight is the height of the footer (default: 1).
	FooterHeight int
}

// CalculatePaneLayout calculates the layout for a multi-pane UI.
// The layout distributes available height in a 2:2:1 ratio among Tables, Schema, and SQL panes.
func CalculatePaneLayout(config LayoutConfig) PaneLayout {
	// Apply defaults
	if config.LeftPaneContentWidth <= 0 {
		config.LeftPaneContentWidth = 50
	}
	if config.MinTablesHeight <= 0 {
		config.MinTablesHeight = 3
	}
	if config.MinSchemaHeight <= 0 {
		config.MinSchemaHeight = 3
	}
	if config.MinSQLHeight <= 0 {
		config.MinSQLHeight = 2
	}
	if config.FooterHeight <= 0 {
		config.FooterHeight = 1
	}

	layout := PaneLayout{
		TotalWidth:           config.TotalWidth,
		TotalHeight:          config.TotalHeight,
		LeftPaneContentWidth: config.LeftPaneContentWidth,
		LeftPaneActualWidth:  config.LeftPaneContentWidth + 2, // +2 for borders
		ConnectionHeight:     config.ConnectionPaneHeight,
		FooterHeight:         config.FooterHeight,
	}

	// Calculate right pane width
	layout.RightPaneActualWidth = config.TotalWidth - layout.LeftPaneActualWidth + 1

	// Calculate available height for Tables, Schema, SQL panes
	// Total - Footer - ConnectionPane - borders (approximately 6 lines for separators etc.)
	availableHeight := config.TotalHeight - config.FooterHeight - config.ConnectionPaneHeight - 6

	// Distribute in 2:2:1 ratio
	partHeight := availableHeight / PaneHeightTotalParts
	remainder := availableHeight % PaneHeightTotalParts

	tablesHeight := partHeight * PaneHeightTablesParts
	schemaHeight := partHeight * PaneHeightSchemaParts
	sqlHeight := partHeight * PaneHeightSQLParts

	// Distribute remainder
	DistributeSpace(remainder, &tablesHeight, &schemaHeight, &sqlHeight)

	// Ensure minimum heights
	if tablesHeight < config.MinTablesHeight {
		tablesHeight = config.MinTablesHeight
	}
	if schemaHeight < config.MinSchemaHeight {
		schemaHeight = config.MinSchemaHeight
	}
	if sqlHeight < config.MinSQLHeight {
		sqlHeight = config.MinSQLHeight
	}

	// If we're under the available height after minimums, redistribute
	usedHeight := tablesHeight + schemaHeight + sqlHeight
	if usedHeight < availableHeight {
		DistributeSpace(availableHeight-usedHeight, &tablesHeight, &schemaHeight, &sqlHeight)
	}

	layout.TablesHeight = tablesHeight
	layout.SchemaHeight = schemaHeight
	layout.SQLHeight = sqlHeight

	return layout
}

// SplitPaneLayout holds dimensions for a two-pane split layout.
type SplitPaneLayout struct {
	LeftWidth  int
	RightWidth int
	Height     int
}

// CalculateSplitLayout calculates a horizontal split layout.
// ratio is the proportion of width for the left pane (0.0-1.0).
func CalculateSplitLayout(totalWidth, totalHeight int, ratio float64) SplitPaneLayout {
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	leftWidth := int(float64(totalWidth) * ratio)
	rightWidth := totalWidth - leftWidth

	return SplitPaneLayout{
		LeftWidth:  leftWidth,
		RightWidth: rightWidth,
		Height:     totalHeight,
	}
}

// ContentDimensions calculates content dimensions after removing borders and padding.
func ContentDimensions(totalWidth, totalHeight, borderWidth, padding int) (contentWidth, contentHeight int) {
	// Each side has border and padding
	contentWidth = totalWidth - (borderWidth * 2) - (padding * 2)
	contentHeight = totalHeight - (borderWidth * 2) - (padding * 2)

	if contentWidth < 0 {
		contentWidth = 0
	}
	if contentHeight < 0 {
		contentHeight = 0
	}

	return contentWidth, contentHeight
}

// CenterPosition calculates the position to center an item.
func CenterPosition(containerSize, itemSize int) int {
	pos := (containerSize - itemSize) / 2
	if pos < 0 {
		return 0
	}
	return pos
}

// DistributeSpace distributes a given amount of space across multiple targets
// in a round-robin fashion. Each target gets incremented by 1 until all space
// is distributed.
func DistributeSpace(amount int, targets ...*int) {
	if len(targets) == 0 {
		return
	}
	for amount > 0 {
		for _, target := range targets {
			if amount <= 0 {
				return
			}
			*target++
			amount--
		}
	}
}

// Separator renders a horizontal separator line with the specified width.
// Uses the StyleSeparator for consistent styling.
func Separator(width int) string {
	return StyleSeparator.Render(strings.Repeat("─", width))
}

// BorderedBox renders content wrapped in a border with an optional title.
// The border uses box-drawing characters (╭╮╰╯│).
// If title is provided, it appears in the top border with a blank line below it.
//
// Example with title:
//   BorderedBox("content", 40, "Title")
//   ╭── Title ────╮
//   │            │
//   │  content   │
//   ╰────────────╯
//
// Example without title:
//   BorderedBox("content", 40)
//   ╭────────────╮
//   │  content   │
//   ╰────────────╯
func BorderedBox(content string, width int, title ...string) string {
	var result strings.Builder

	// Top border
	if len(title) > 0 && title[0] != "" {
		// ╭── Title ─────╮
		// Calculate: "╭──" (3) + title + "╮" (1) + padding
		titleWithPadding := " " + title[0] + " "
		remainingWidth := width - 3 - len(titleWithPadding) - 1
		if remainingWidth < 0 {
			remainingWidth = 0
		}
		topBorder := StyleBorder.Render("╭──" + titleWithPadding + strings.Repeat("─", remainingWidth) + "╮")
		result.WriteString(topBorder + "\n")

		// Empty line after title
		emptyLine := strings.Repeat(" ", width-2)
		leftBorder := StyleBorder.Render("│")
		rightBorder := StyleBorder.Render("│")
		result.WriteString(leftBorder + emptyLine + rightBorder + "\n")
	} else {
		// ╭──────────╮
		topBorder := StyleBorder.Render("╭" + strings.Repeat("─", width-2) + "╮")
		result.WriteString(topBorder + "\n")
	}

	// Content with left and right borders
	leftBorder := StyleBorder.Render("│")
	rightBorder := StyleBorder.Render("│")

	for _, line := range strings.Split(content, "\n") {
		if line != "" {
			result.WriteString(leftBorder + line + rightBorder + "\n")
		}
	}

	// Bottom border
	bottomBorder := StyleBorder.Render("╰" + strings.Repeat("─", width-2) + "╯")
	result.WriteString(bottomBorder)

	return result.String()
}
