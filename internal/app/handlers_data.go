package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/camikura/dito/internal/db"
	"github.com/camikura/dito/internal/ui"
)

func handleDataKeys(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Get total row count for current table
	var totalRows int
	if data := m.GetSelectedTableData(); data != nil {
		totalRows = len(data.Rows)
	}

	// Calculate visible lines for data rows
	// Data pane structure: title(1) + content lines + bottom(1)
	// Content lines = header(1) + separator(1) + data rows
	contentLines := m.Window.Height - ui.DataPaneTitleAndBorderLines
	if contentLines < ui.MinContentLines {
		contentLines = ui.MinContentLines
	}
	// Data visible lines = content lines - header lines
	dataVisibleLines := contentLines - ui.DataPaneHeaderLines
	if dataVisibleLines < 1 {
		dataVisibleLines = 1
	}

	// Calculate max horizontal offset
	maxHorizontalOffset := calculateMaxHorizontalOffset(m)

	// Calculate maximum viewport offset
	maxViewportOffset := totalRows - dataVisibleLines
	if maxViewportOffset < 0 {
		maxViewportOffset = 0
	}
	_ = maxViewportOffset // Currently unused but kept for future use

	// Handle M-< and M-> (Alt+Shift+, and Alt+Shift+.)
	// On Mac, these produce special characters: ¯ (175) and ˘ (728)
	switch msg.String() {
	case "alt+<", "¯":
		// Jump to first row
		m.Data.SelectedDataRow = 0
		m.Data.ViewportOffset = 0
		return m, nil

	case "alt+>", "˘":
		// Jump to last row, keeping cursor at center of screen (VS Code style)
		if totalRows > 0 {
			m.Data.SelectedDataRow = totalRows - 1

			// Calculate middle position of visible area
			middlePosition := dataVisibleLines / 2

			// Set viewport so cursor appears at center
			// This leaves empty space below if at end of data
			m.Data.ViewportOffset = m.Data.SelectedDataRow - middlePosition
			if m.Data.ViewportOffset < 0 {
				m.Data.ViewportOffset = 0
			}

			// Check if we need to fetch more data
			if cmd := fetchMoreDataIfNeeded(m, false); cmd != nil {
				m.Data.LoadingData = true
				return m, cmd
			}
		}
		return m, nil
	}

	switch msg.Type {
	case tea.KeyUp, tea.KeyCtrlP:
		if m.Data.SelectedDataRow > 0 {
			m.Data.SelectedDataRow--

			// Calculate middle position of visible area
			middlePosition := dataVisibleLines / 2

			// Scrolling logic: keep cursor at middle of screen (VS Code style)
			if m.Data.SelectedDataRow > middlePosition {
				m.Data.ViewportOffset = m.Data.SelectedDataRow - middlePosition
			} else {
				m.Data.ViewportOffset = 0
			}
		}
		return m, nil

	case tea.KeyDown, tea.KeyCtrlN:
		if totalRows > 0 && m.Data.SelectedDataRow < totalRows-1 {
			m.Data.SelectedDataRow++

			// Calculate middle position of visible area
			middlePosition := dataVisibleLines / 2

			// Scrolling logic: keep cursor at middle of screen (VS Code style)
			// This allows scrolling past the last row to show empty space below
			if m.Data.SelectedDataRow > middlePosition {
				m.Data.ViewportOffset = m.Data.SelectedDataRow - middlePosition
			}

			// Check if we need to fetch more data
			remainingRows := totalRows - m.Data.SelectedDataRow - 1
			if cmd := fetchMoreDataIfNeeded(m, remainingRows <= ui.FetchMoreThreshold); cmd != nil {
				m.Data.LoadingData = true
				return m, cmd
			}
		}
		return m, nil

	case tea.KeyLeft, tea.KeyCtrlB:
		if m.Data.HorizontalOffset > 0 {
			m.Data.HorizontalOffset--
		}
		return m, nil

	case tea.KeyRight, tea.KeyCtrlF:
		if m.Data.HorizontalOffset < maxHorizontalOffset {
			m.Data.HorizontalOffset++
		}
		return m, nil

	case tea.KeyEnter:
		// Show record detail dialog
		if totalRows > 0 && m.Data.SelectedDataRow < totalRows {
			m.RecordDetail.Visible = true
			m.RecordDetail.ScrollOffset = 0
		}
		return m, nil

	case tea.KeyEscape:
		// Reset to default SQL (only if custom SQL is active)
		if m.SQL.CustomSQL {
			m.SQL.CustomSQL = false
			m.SQL.ColumnOrder = nil
			m.Data.SelectedDataRow = 0
			m.Data.ViewportOffset = 0
			m.Data.HorizontalOffset = 0
			m.Schema.ErrorMsg = ""
			m.Data.ErrorMsg = ""

			// Reload data with default SQL if a table is selected
			tableName := m.SelectedTableName()
			if tableName != "" {
				var ddl string
				var primaryKeys []string
				if details := m.GetSelectedTableDetails(); details != nil && details.Schema != nil {
					ddl = details.Schema.DDL
					primaryKeys = ui.ParsePrimaryKeysFromDDL(ddl)
				}

				m.SQL.CurrentSQL = buildDefaultSQL(tableName, ddl)
				m.SQL.CursorPos = ui.RuneLen(m.SQL.CurrentSQL)
				return m, db.FetchTableData(m.Connection.NosqlClient, tableName, ui.DefaultFetchSize, primaryKeys)
			}
			m.SQL.CurrentSQL = ""
			m.SQL.CursorPos = 0
		}
		return m, nil

	case tea.KeyCtrlA:
		// Scroll to leftmost
		m.Data.HorizontalOffset = 0
		return m, nil

	case tea.KeyCtrlE:
		// Scroll to rightmost
		m.Data.HorizontalOffset = maxHorizontalOffset
		return m, nil
	}

	return m, nil
}

// handleDataCopy copies the selected row to clipboard
func handleDataCopy(m Model) (Model, tea.Cmd) {
	tableName := m.SelectedTableName()
	if tableName == "" {
		return m, nil
	}

	data := m.GetSelectedTableData()
	if data == nil || len(data.Rows) == 0 {
		return m, nil
	}

	if m.Data.SelectedDataRow < 0 || m.Data.SelectedDataRow >= len(data.Rows) {
		return m, nil
	}

	row := data.Rows[m.Data.SelectedDataRow]

	// Get column order to match display order
	columnOrder := getColumnsInSchemaOrder(m, tableName, data.Rows)

	err := ui.CopyRowToClipboard(row, columnOrder)
	if err != nil {
		m.UI.CopyMessage = "Copy failed: " + err.Error()
	} else {
		m.UI.CopyMessage = "Copied to clipboard"
	}

	// Clear message after a short delay using a timer command
	return m, tea.Tick(ui.CopyMessageDuration, func(_ time.Time) tea.Msg {
		return clearCopyMessageMsg{}
	})
}

func handleRecordDetailKeys(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Calculate max scroll for record detail
	maxScroll := calculateRecordDetailMaxScroll(m)

	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit

	case tea.KeyEscape:
		m.RecordDetail.Visible = false
		m.RecordDetail.ScrollOffset = 0
		return m, nil

	case tea.KeyUp, tea.KeyCtrlP:
		if m.RecordDetail.ScrollOffset > 0 {
			m.RecordDetail.ScrollOffset--
		}
		return m, nil

	case tea.KeyDown, tea.KeyCtrlN:
		if m.RecordDetail.ScrollOffset < maxScroll {
			m.RecordDetail.ScrollOffset++
		}
		return m, nil

	case tea.KeyHome:
		m.RecordDetail.ScrollOffset = 0
		return m, nil

	case tea.KeyEnd:
		m.RecordDetail.ScrollOffset = maxScroll
		return m, nil

	case tea.KeyPgUp:
		// Scroll up by page
		m.RecordDetail.ScrollOffset -= ui.PageScrollAmount
		if m.RecordDetail.ScrollOffset < 0 {
			m.RecordDetail.ScrollOffset = 0
		}
		return m, nil

	case tea.KeyPgDown:
		// Scroll down by page
		m.RecordDetail.ScrollOffset += ui.PageScrollAmount
		if m.RecordDetail.ScrollOffset > maxScroll {
			m.RecordDetail.ScrollOffset = maxScroll
		}
		return m, nil
	}

	return m, nil
}
