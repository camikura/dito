package new_ui

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/camikura/dito/internal/db"
	"github.com/camikura/dito/internal/views"
)

// Update handles messages and updates the model
func Update(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		// Calculate pane heights dynamically
		// Use actual connection pane height from model, or default to 5 if not yet set
		connectionPaneHeight := m.ConnectionPaneHeight
		if connectionPaneHeight == 0 {
			connectionPaneHeight = 5 // Default for cloud connection
		}
		// Available height for Tables, Schema, and SQL content (2:2:1 ratio)
		// Total: m.Height = leftPanes + footer
		// leftPanes = Connection + Tables(+2) + Schema(+2) + SQL(+2)
		// So: availableHeight = m.Height - 1(footer) - connectionPaneHeight - 6(borders)
		availableHeight := m.Height - 1 - connectionPaneHeight - 6

		// Split available height in 2:2:1 ratio (Tables:Schema:SQL)
		totalParts := 5 // 2+2+1
		partHeight := availableHeight / totalParts
		remainder := availableHeight % totalParts

		m.TablesHeight = partHeight * 2
		m.SchemaHeight = partHeight * 2
		m.SQLHeight = partHeight

		// Distribute remainder to maximize space usage (may be up to 4)
		for remainder > 0 {
			if remainder >= 1 {
				m.TablesHeight++
				remainder--
			}
			if remainder >= 1 {
				m.SchemaHeight++
				remainder--
			}
			if remainder >= 1 {
				m.SQLHeight++
				remainder--
			}
		}

		// Ensure minimum heights
		if m.TablesHeight < 3 {
			m.TablesHeight = 3
		}
		if m.SchemaHeight < 3 {
			m.SchemaHeight = 3
		}
		if m.SQLHeight < 2 {
			m.SQLHeight = 2
		}

		// After applying minimum heights, check if we have unused space
		usedHeight := m.TablesHeight + m.SchemaHeight + m.SQLHeight
		if usedHeight < availableHeight {
			// Distribute unused space in 2:2:1 ratio again
			extraSpace := availableHeight - usedHeight
			for extraSpace > 0 {
				if extraSpace >= 1 {
					m.TablesHeight++
					extraSpace--
				}
				if extraSpace >= 1 {
					m.SchemaHeight++
					extraSpace--
				}
				if extraSpace >= 1 {
					m.SQLHeight++
					extraSpace--
				}
			}
		}

		return m, nil

	case tea.KeyMsg:
		return handleKeyPress(m, msg)

	case db.ConnectionResult:
		return handleConnectionResult(m, msg)

	case db.TableListResult:
		return handleTableListResult(m, msg)

	case db.TableDetailsResult:
		return handleTableDetailsResult(m, msg)

	case db.TableDataResult:
		return handleTableDataResult(m, msg)
	}

	return m, nil
}

func handleKeyPress(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "tab":
		m = m.NextPane()
		return m, nil

	case "shift+tab":
		m = m.PrevPane()
		return m, nil
	}

	// Pane-specific keys
	switch m.CurrentPane {
	case FocusPaneTables:
		return handleTablesKeys(m, msg)
	case FocusPaneSchema:
		return handleSchemaKeys(m, msg)
	case FocusPaneData:
		return handleDataKeys(m, msg)
	}

	return m, nil
}

func handleTablesKeys(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	visibleLines := m.TablesHeight // Tables pane visible height (dynamic)

	switch msg.String() {
	case "up", "k":
		if m.CursorTable > 0 {
			m.CursorTable--
			m.SelectedTable = m.CursorTable // Move selection with cursor
			m.SchemaScrollOffset = 0 // Reset scroll when changing tables

			// Adjust scroll offset to keep cursor visible
			if m.CursorTable < m.TablesScrollOffset {
				m.TablesScrollOffset = m.CursorTable
			}

			// Auto-update schema for table under cursor
			if m.CursorTable < len(m.Tables) {
				tableName := m.Tables[m.CursorTable]
				return m, db.FetchTableDetails(m.NosqlClient, tableName)
			}
		}
		return m, nil

	case "down", "j":
		if m.CursorTable < len(m.Tables)-1 {
			m.CursorTable++
			m.SelectedTable = m.CursorTable // Move selection with cursor
			m.SchemaScrollOffset = 0 // Reset scroll when changing tables

			// Adjust scroll offset to keep cursor visible
			if m.CursorTable >= m.TablesScrollOffset+visibleLines {
				m.TablesScrollOffset = m.CursorTable - visibleLines + 1
			}

			// Auto-update schema for table under cursor
			if m.CursorTable < len(m.Tables) {
				tableName := m.Tables[m.CursorTable]
				return m, db.FetchTableDetails(m.NosqlClient, tableName)
			}
		}
		return m, nil

	case "enter":
		// Select table and load data
		if m.CursorTable < len(m.Tables) {
			m.SelectedTable = m.CursorTable
			tableName := m.Tables[m.SelectedTable]

			// Generate SQL query
			m.CurrentSQL = "SELECT * FROM " + tableName
			m.CustomSQL = false

			// Reset data scrolling state
			m.SelectedDataRow = 0
			m.ViewportOffset = 0
			m.HorizontalOffset = 0

			// Get primary keys from schema if available
			var primaryKeys []string
			if details, exists := m.TableDetails[tableName]; exists && details != nil && details.Schema != nil {
				// We'll parse from DDL in a moment - for now just pass empty slice
				primaryKeys = []string{}
			}

			// Load table data
			return m, db.FetchTableData(m.NosqlClient, tableName, 1000, primaryKeys)
		}
		return m, nil
	}

	return m, nil
}

func handleSchemaKeys(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Calculate max scroll offset based on content
	var maxScroll int
	if len(m.Tables) > 0 && m.CursorTable < len(m.Tables) {
		tableName := m.Tables[m.CursorTable]
		if details, exists := m.TableDetails[tableName]; exists && details != nil {
			// Count content lines
			lineCount := 1 // "Columns:"
			if details.Schema.DDL != "" {
				primaryKeys := views.ParsePrimaryKeysFromDDL(details.Schema.DDL)
				columns := views.ParseColumnsFromDDL(details.Schema.DDL, primaryKeys)
				lineCount += len(columns)
			}
			lineCount += 2 // Empty line + "Indexes:"
			lineCount += len(details.Indexes)
			if len(details.Indexes) == 0 {
				lineCount++ // "(none)" line
			}

			// Max scroll = total lines - visible lines (dynamic)
			maxScroll = lineCount - m.SchemaHeight
			if maxScroll < 0 {
				maxScroll = 0
			}
		}
	}

	switch msg.String() {
	case "down", "j": // Scroll down
		if m.SchemaScrollOffset < maxScroll {
			m.SchemaScrollOffset++
		}
		return m, nil

	case "up", "k": // Scroll up
		if m.SchemaScrollOffset > 0 {
			m.SchemaScrollOffset--
		}
		return m, nil
	}

	return m, nil
}

func handleDataKeys(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Get total row count for current table
	var totalRows int
	if m.SelectedTable >= 0 && m.SelectedTable < len(m.Tables) {
		tableName := m.Tables[m.SelectedTable]
		if data, exists := m.TableData[tableName]; exists && data != nil {
			totalRows = len(data.Rows)
		}
	}

	// Calculate visible lines (Data pane height - 2 for header and separator)
	visibleLines := m.Height - m.ConnectionPaneHeight - m.TablesHeight - m.SchemaHeight - m.SQLHeight - 8 // 8 = borders + footer
	if visibleLines < 1 {
		visibleLines = 1
	}
	// Subtract 2 for header and separator in data pane
	dataVisibleLines := visibleLines - 2
	if dataVisibleLines < 1 {
		dataVisibleLines = 1
	}

	// Calculate max horizontal offset
	maxHorizontalOffset := calculateMaxHorizontalOffset(m)

	switch msg.String() {
	case "up", "k":
		if m.SelectedDataRow > 0 {
			m.SelectedDataRow--

			// Adjust viewport offset if cursor goes above visible area
			if m.SelectedDataRow < m.ViewportOffset {
				m.ViewportOffset = m.SelectedDataRow
			}
		}
		return m, nil

	case "down", "j":
		if totalRows > 0 && m.SelectedDataRow < totalRows-1 {
			m.SelectedDataRow++

			// Adjust viewport offset if cursor goes below visible area
			if m.SelectedDataRow >= m.ViewportOffset+dataVisibleLines {
				m.ViewportOffset = m.SelectedDataRow - dataVisibleLines + 1
			}

			// Limit viewport offset so we don't scroll past the last screen
			maxViewportOffset := totalRows - dataVisibleLines
			if maxViewportOffset < 0 {
				maxViewportOffset = 0
			}
			if m.ViewportOffset > maxViewportOffset {
				m.ViewportOffset = maxViewportOffset
			}
		}
		return m, nil

	case "left", "h":
		// Scroll left
		if m.HorizontalOffset > 0 {
			m.HorizontalOffset -= 5 // Scroll by 5 characters
			if m.HorizontalOffset < 0 {
				m.HorizontalOffset = 0
			}
		}
		return m, nil

	case "right", "l":
		// Scroll right (limited to max offset)
		if m.HorizontalOffset < maxHorizontalOffset {
			m.HorizontalOffset += 5 // Scroll by 5 characters
			if m.HorizontalOffset > maxHorizontalOffset {
				m.HorizontalOffset = maxHorizontalOffset
			}
		}
		return m, nil
	}

	return m, nil
}

// calculateMaxHorizontalOffset calculates the maximum horizontal scroll offset
// so that the rightmost column is visible at the right edge of the pane
func calculateMaxHorizontalOffset(m Model) int {
	if m.SelectedTable < 0 || m.SelectedTable >= len(m.Tables) {
		return 0
	}

	tableName := m.Tables[m.SelectedTable]
	data, exists := m.TableData[tableName]
	if !exists || data == nil || len(data.Rows) == 0 {
		return 0
	}

	// Get columns in schema order
	var columns []string
	for col := range data.Rows[0] {
		columns = append(columns, col)
	}

	// Calculate natural column widths
	widths := make([]int, len(columns))
	for i, col := range columns {
		widths[i] = len(col)
		if widths[i] < 3 {
			widths[i] = 3
		}
	}

	// Check data widths (sample first 100 rows)
	sampleSize := len(data.Rows)
	if sampleSize > 100 {
		sampleSize = 100
	}
	for i := 0; i < sampleSize; i++ {
		row := data.Rows[i]
		for j, col := range columns {
			if val, exists := row[col]; exists && val != nil {
				valStr := fmt.Sprintf("%v", val)
				valLen := len([]rune(valStr))
				if valLen > widths[j] {
					widths[j] = valLen
				}
			}
		}
	}

	// Cap at 50 characters per column
	for i := range widths {
		if widths[i] > 50 {
			widths[i] = 50
		}
	}

	// Calculate total width (columns + separators)
	totalWidth := 0
	for _, w := range widths {
		totalWidth += w
	}
	// Add separators (1 space between columns)
	if len(widths) > 0 {
		totalWidth += len(widths) - 1
	}

	// Calculate data pane visible width
	dataWidth := m.Width - 2 // Subtract borders

	// Max offset = total width - visible width
	maxOffset := totalWidth - dataWidth
	if maxOffset < 0 {
		maxOffset = 0
	}

	return maxOffset
}

func handleConnectionResult(m Model, msg db.ConnectionResult) (Model, tea.Cmd) {
	if msg.Err != nil {
		// Connection failed
		m.Connected = false
		m.ConnectionMsg = msg.Err.Error()
		return m, nil
	}

	// Connection successful
	m.Connected = true
	m.NosqlClient = msg.Client
	m.Endpoint = msg.Endpoint
	m.ConnectionMsg = ""

	// Fetch table list
	return m, db.FetchTables(msg.Client)
}

func handleTableListResult(m Model, msg db.TableListResult) (Model, tea.Cmd) {
	if msg.Err != nil {
		// TODO: Show error
		return m, nil
	}

	// Sort tables for tree display (parents before children)
	m.Tables = sortTablesForTree(msg.Tables)
	if len(m.Tables) > 0 {
		m.CursorTable = 0
		m.SelectedTable = 0 // Initialize selection to first table
		// Fetch details for first table
		return m, db.FetchTableDetails(m.NosqlClient, m.Tables[0])
	}

	return m, nil
}

// sortTablesForTree sorts table names so parent tables appear before their children
// e.g., ["users.phones", "users", "products", "users.addresses"] ->
//       ["products", "users", "users.addresses", "users.phones"]
func sortTablesForTree(tables []string) []string {
	sorted := make([]string, len(tables))
	copy(sorted, tables)

	sort.Slice(sorted, func(i, j int) bool {
		a, b := sorted[i], sorted[j]

		// Get parent names
		parentA := a
		if dotIndex := strings.LastIndex(a, "."); dotIndex != -1 {
			parentA = a[:dotIndex]
		}
		parentB := b
		if dotIndex := strings.LastIndex(b, "."); dotIndex != -1 {
			parentB = b[:dotIndex]
		}

		// If one is parent of the other, parent comes first
		if a == parentB {
			return true // a is parent of b
		}
		if b == parentA {
			return false // b is parent of a
		}

		// If they have the same parent, sort alphabetically
		if parentA == parentB {
			return a < b
		}

		// Different parents - sort by parent name, then by full name
		if parentA != a && parentB != b {
			// Both are children - compare parents first
			if parentA != parentB {
				return parentA < parentB
			}
		}

		// One is parent, one is not - sort by full name
		return a < b
	})

	return sorted
}

func handleTableDetailsResult(m Model, msg db.TableDetailsResult) (Model, tea.Cmd) {
	if msg.Err != nil {
		// TODO: Show error
		return m, nil
	}

	m.TableDetails[msg.TableName] = &msg
	return m, nil
}

func handleTableDataResult(m Model, msg db.TableDataResult) (Model, tea.Cmd) {
	if msg.Err != nil {
		// TODO: Show error
		return m, nil
	}

	m.TableData[msg.TableName] = &msg
	m.LoadingData = false
	return m, nil
}
