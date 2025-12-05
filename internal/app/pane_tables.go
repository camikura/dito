package app

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/camikura/dito/internal/ui"
)

func renderTablesPane(m Model, width int) string {
	return renderTablesPaneWithHeight(m, width, 12)
}

func renderTablesPaneWithHeight(m Model, width int, height int) string {
	borderStyle := ui.StyleBorderInactive
	titleStyle := ui.StyleTitleInactive
	if m.CurrentPane == FocusPaneTables {
		borderStyle = ui.StyleBorderActive
		titleStyle = ui.StyleTitleActive
	}

	titleText := " Tables"
	if len(m.Tables) > 0 {
		titleText += fmt.Sprintf(" (%d)", len(m.Tables))
	}
	titleText += " "

	dashCount := width - ui.RuneLen(titleText) - 3
	if dashCount < 0 {
		dashCount = 0
	}
	styledTitle := titleStyle.Render(titleText)
	title := borderStyle.Render("╭─") + styledTitle + borderStyle.Render(strings.Repeat("─", dashCount) + "╮")

	// Prepare content lines with tree structure
	type tableLineInfo struct {
		text       string
		isSelected bool // * marker (Enter pressed)
		isCursor   bool // cursor position (up/down navigation)
	}

	// Determine if selection marker should be shown
	// Hide * when custom SQL targets a table not in the list
	showSelectionMarker := true
	if m.CustomSQL && m.CurrentSQL != "" {
		extractedName := ui.ExtractTableNameFromSQL(m.CurrentSQL)
		if extractedName != "" && m.FindTableName(extractedName) == "" {
			// Custom SQL targets a table not in the list
			showSelectionMarker = false
		}
	}

	var contentLines []tableLineInfo
	if len(m.Tables) == 0 {
		contentLines = []tableLineInfo{{text: "No tables", isSelected: false, isCursor: false}}
	} else {
		// Render each table with tree structure
		// Calculate available width for table name (excluding borders)
		availableWidth := width - 2 // -2 for left and right borders

		for i, tableName := range m.Tables {
			// Determine indentation based on nesting level (count of '.' separators)
			nestLevel := strings.Count(tableName, ".")
			indent := strings.Repeat(" ", nestLevel)
			displayName := tableName
			if dotIndex := strings.LastIndex(tableName, "."); dotIndex != -1 {
				// Child table - show only the last part of the name
				displayName = tableName[dotIndex+1:]
			}

			// Add selection marker (* for selected table via Enter)
			var prefix string
			isSelected := showSelectionMarker && i == m.SelectedTable
			if isSelected {
				prefix = "* "
			} else {
				prefix = "  "
			}

			// Truncate if too long
			fullText := prefix + indent + displayName
			maxTextWidth := availableWidth
			if ui.RuneLen(fullText) > maxTextWidth {
				// Truncate with ellipsis
				fullText = ui.TruncateString(fullText, maxTextWidth)
			}

			contentLines = append(contentLines, tableLineInfo{
				text:       fullText,
				isSelected: isSelected,
				isCursor:   i == m.CursorTable,
			})
		}
	}

	leftBorder := borderStyle.Render("│")
	rightBorder := borderStyle.Render("│")
	bottomBorder := borderStyle.Render("╰" + strings.Repeat("─", width-2) + "╯")

	var result strings.Builder
	result.WriteString(title + "\n")

	// Render content lines (fill allocated height with content or empty lines)
	isFocused := m.CurrentPane == FocusPaneTables
	for i := 0; i < height; i++ {
		contentIndex := i + m.TablesScrollOffset
		if contentIndex < len(contentLines) {
			lineInfo := contentLines[contentIndex]
			// Apply color based on state
			var styledText string
			if isFocused && lineInfo.isCursor {
				styledText = ui.StyleTableCursor.Render(lineInfo.text)
			} else if lineInfo.isSelected {
				styledText = ui.StyleTableSelected.Render(lineInfo.text)
			} else {
				styledText = ui.StyleTableNormal.Render(lineInfo.text)
			}
			// Calculate padding (based on rune length for correct display width)
			paddingLen := width - ui.RuneLen(lineInfo.text) - 2
			if paddingLen < 0 {
				paddingLen = 0
			}
			line := styledText + strings.Repeat(" ", paddingLen)
			result.WriteString(leftBorder + line + rightBorder + "\n")
		} else {
			// Empty line for remaining allocated height
			emptyLine := strings.Repeat(" ", width-2)
			result.WriteString(leftBorder + emptyLine + rightBorder + "\n")
		}
	}

	result.WriteString(bottomBorder)

	return result.String()
}

func renderSchemaPane(m Model, width int) string {
	return renderSchemaPaneWithHeight(m, width, 12)
}

func renderSchemaPaneWithHeight(m Model, width int, height int) string {
	// Determine which table to show schema for
	// Use SelectedTable, or extract from custom SQL if applicable
	var schemaTableName string
	if m.CustomSQL && m.CurrentSQL != "" {
		// Extract table name from custom SQL and find exact match from tables list
		extractedName := ui.ExtractTableNameFromSQL(m.CurrentSQL)
		schemaTableName = m.FindTableName(extractedName)
		// Use extracted name if not found in tables list
		if schemaTableName == "" && extractedName != "" {
			schemaTableName = extractedName
		}
	} else if m.SelectedTable >= 0 && m.SelectedTable < len(m.Tables) {
		schemaTableName = m.Tables[m.SelectedTable]
	}

	// Title includes table name if available
	titleText := " Schema "
	if schemaTableName != "" {
		titleText = " Schema (" + schemaTableName + ") "
	}

	// Truncate title if too long (leave room for borders: ╭─ ... ─╮)
	maxTitleLen := width - 3
	if ui.RuneLen(titleText) > maxTitleLen {
		// Truncate table name within title
		maxTableNameLen := maxTitleLen - ui.RuneLen(" Schema (...) ")
		if maxTableNameLen > 3 {
			titleText = " Schema (" + ui.TruncateString(schemaTableName, maxTableNameLen) + ") "
		} else {
			titleText = " Schema "
		}
	}

	// Schema pane can be focused for scrolling
	borderStyle := ui.StyleBorderInactive
	titleStyle := ui.StyleTitleInactive
	if m.CurrentPane == FocusPaneSchema {
		borderStyle = ui.StyleBorderActive
		titleStyle = ui.StyleTitleActive
	}
	dashCount := width - ui.RuneLen(titleText) - 3
	if dashCount < 0 {
		dashCount = 0
	}
	styledTitle := titleStyle.Render(titleText)
	title := borderStyle.Render("╭─") + styledTitle + borderStyle.Render(strings.Repeat("─", dashCount) + "╮")

	// Prepare content lines
	var contentLines []string
	var schemaError string
	if m.SchemaErrorMsg != "" {
		schemaError = m.SchemaErrorMsg
	}
	if schemaTableName == "" {
		if m.CustomSQL && m.CurrentSQL != "" {
			// Custom SQL with table not found in tables list
			contentLines = []string{"No schema"}
		} else {
			contentLines = []string{"Select a table"}
		}
	} else if schemaError != "" {
		// Show error message
		contentLines = []string{schemaError}
	} else {
		details, exists := m.TableDetails[schemaTableName]
		if !exists || details == nil {
			contentLines = []string{"Loading..."}
		} else {
			// Render schema information
			contentLines = append(contentLines, "Columns:")

			// Collect all columns including inherited from ancestors
			var allColumns []ui.ColumnInfo

			// Get ancestor table names (root to immediate parent)
			ancestors := ui.GetAncestorTableNames(schemaTableName)

			// Add inherited primary key columns from ancestors
			for _, ancestorName := range ancestors {
				ancestorDetails, ancestorExists := m.TableDetails[ancestorName]
				if ancestorExists && ancestorDetails != nil && ancestorDetails.Schema != nil && ancestorDetails.Schema.DDL != "" {
					ancestorPKs := ui.ParsePrimaryKeysFromDDL(ancestorDetails.Schema.DDL)
					ancestorCols := ui.ParseColumnsFromDDL(ancestorDetails.Schema.DDL, ancestorPKs)
					// Only add primary key columns from ancestors
					for _, col := range ancestorCols {
						if col.IsPrimaryKey {
							col.IsInherited = true
							allColumns = append(allColumns, col)
						}
					}
				}
			}

			// Add this table's own columns
			if details.Schema.DDL != "" {
				primaryKeys := ui.ParsePrimaryKeysFromDDL(details.Schema.DDL)
				columns := ui.ParseColumnsFromDDL(details.Schema.DDL, primaryKeys)
				allColumns = append(allColumns, columns...)
			}

			// Find the longest column name
			maxColNameLen := 0
			for _, col := range allColumns {
				if len(col.Name) > maxColNameLen {
					maxColNameLen = len(col.Name)
				}
			}

			// Format each column: PK|||Name|||Type|||maxLen|||IsInherited (use ||| as separator)
			for _, col := range allColumns {
				pkMarker := " " // Single space when not PK
				if col.IsPrimaryKey {
					pkMarker = "P" // Single "P" for primary key
				}
				inherited := ""
				if col.IsInherited {
					inherited = "|||inherited"
				}
				contentLines = append(contentLines, fmt.Sprintf("%s|||%s|||%s|||%d%s", pkMarker, col.Name, col.Type, maxColNameLen, inherited))
			}

			// Add indexes section
			contentLines = append(contentLines, "")
			contentLines = append(contentLines, "Indexes:")
			if len(details.Indexes) > 0 {
				for _, index := range details.Indexes {
					fields := strings.Join(index.FieldNames, ", ")
					// Format: IndexName|||Fields (use ||| as separator to apply color in rendering)
					contentLines = append(contentLines, fmt.Sprintf("IDX|||%s|||%s", index.IndexName, fields))
				}
			} else {
				contentLines = append(contentLines, "  (none)")
			}
		}
	}

	leftBorder := borderStyle.Render("│")
	rightBorder := borderStyle.Render("│")
	bottomBorder := borderStyle.Render("╰" + strings.Repeat("─", width-2) + "╯")

	var result strings.Builder
	result.WriteString(title + "\n")

	// Render content lines (fill allocated height with content or empty lines)
	for i := 0; i < height; i++ {
		var line string
		contentIndex := i + m.SchemaScrollOffset
		if contentIndex < len(contentLines) {
			content := contentLines[contentIndex]

			// Apply color to label lines (Columns:, Indexes:)
			if content == "Columns:" || content == "Indexes:" {
				paddingLen := width - len(content) - 2
				if paddingLen < 0 {
					paddingLen = 0
				}
				line = ui.StyleSchemaLabel.Render(content) + strings.Repeat(" ", paddingLen)
			} else if strings.HasPrefix(content, "IDX|||") {
				// Index line: IDX|||IndexName|||Fields
				parts := strings.Split(content, "|||")
				if len(parts) >= 3 {
					indexName := parts[1]
					fields := parts[2]

					// Format: "  indexName fields" with field names in index color and commas in white
					var fieldsDisplay string
					if strings.Contains(fields, ", ") {
						// Multiple fields: color each field name separately, keep commas white
						fieldList := strings.Split(fields, ", ")
						for i, field := range fieldList {
							if i > 0 {
								fieldsDisplay += ", " // White comma
							}
							fieldsDisplay += ui.StyleSchemaIndex.Render(field)
						}
					} else {
						// Single field
						fieldsDisplay = ui.StyleSchemaIndex.Render(fields)
					}

					displayText := "  " + indexName + " " + fieldsDisplay
					displayLen := 2 + len(indexName) + 1 + len(fields)

					availableWidth := width - 2
					rightPadding := availableWidth - displayLen
					if rightPadding < 0 {
						rightPadding = 0
					}
					line = displayText + strings.Repeat(" ", rightPadding)
				} else {
					// Fallback
					paddingLen := width - len(content) - 2
					if paddingLen < 0 {
						paddingLen = 0
					}
					line = content + strings.Repeat(" ", paddingLen)
				}
			} else if strings.Contains(content, "|||") {
				// Column line with PK, name, type, and maxColNameLen separated by |||
				// Format: PK|||Name|||Type|||maxLen or PK|||Name|||Type|||maxLen|||inherited
				parts := strings.Split(content, "|||")
				if len(parts) >= 4 {
					pkMarker := parts[0] // "P" or " "
					colName := parts[1]
					colType := parts[2]
					maxColNameLen, _ := strconv.Atoi(parts[3])
					isInherited := len(parts) >= 5 && parts[4] == "inherited"

					// Fixed column widths for alignment
					const pkColWidth = 2              // Fixed width for PK marker (1 char + 1 space)
					nameColWidth := maxColNameLen + 1 // Use actual max column name length + 1 space

					// PK marker with fixed width
					var pkField string
					if pkMarker == "P" {
						pkField = ui.StyleSchemaPK.Render(pkMarker) + " "
					} else {
						pkField = strings.Repeat(" ", pkColWidth)
					}

					// Pad column name to fixed width
					namePadding := nameColWidth - len(colName)
					if namePadding < 0 {
						namePadding = 0
					}
					nameField := colName + strings.Repeat(" ", namePadding)

					// Type field with inherited marker if applicable
					typeDisplay := colType
					typeDisplayWidth := len(colType)
					if isInherited {
						typeDisplay = colType + " (↑)"
						typeDisplayWidth = len(colType) + 4 // " (↑)" is 4 display chars (space + parens + arrow)
					}
					typeField := ui.StyleSchemaType.Render(typeDisplay)

					// Build line with fixed-width columns: PK + Name + Type
					alignedLine := pkField + nameField + typeField

					// Calculate right padding
					displayLen := pkColWidth + nameColWidth + typeDisplayWidth
					availableWidth := width - 2 // -2 for borders
					rightPadding := availableWidth - displayLen
					if rightPadding < 0 {
						rightPadding = 0
					}
					line = alignedLine + strings.Repeat(" ", rightPadding)
				} else {
					// Fallback if parsing fails
					paddingLen := width - len(content) - 2
					if paddingLen < 0 {
						paddingLen = 0
					}
					line = content + strings.Repeat(" ", paddingLen)
				}
			} else {
				// Other content (like "Select a table", "Loading...", error messages)
				if len(content) > width-2 {
					content = content[:width-5] + "..."
				}
				paddingLen := width - len(content) - 2
				if paddingLen < 0 {
					paddingLen = 0
				}
				// Use red for errors, gray for other messages
				if schemaError != "" && content == schemaError {
					line = ui.StyleErrorLight.Render(content) + strings.Repeat(" ", paddingLen)
				} else {
					line = ui.StyleGrayText.Render(content) + strings.Repeat(" ", paddingLen)
				}
			}
			result.WriteString(leftBorder + line + rightBorder + "\n")
		} else {
			// Empty line for remaining allocated height
			emptyLine := strings.Repeat(" ", width-2)
			result.WriteString(leftBorder + emptyLine + rightBorder + "\n")
		}
	}

	result.WriteString(bottomBorder)

	return result.String()
}
