package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/oracle/nosql-go-sdk/nosqldb"

	"github.com/camikura/dito/internal/app"
	"github.com/camikura/dito/internal/db"
	"github.com/camikura/dito/internal/ui"
)

// TableListViewModel represents the model for table list view
type TableListViewModel struct {
	Width            int
	Height           int
	Endpoint         string
	Tables           []string
	SelectedTable    int
	RightPaneMode    app.RightPaneMode
	TableData        map[string]*db.TableDataResult
	TableDetails     map[string]*db.TableDetailsResult
	LoadingDetails   bool
	LoadingData      bool
	SelectedDataRow  int
	HorizontalOffset int
	ViewportOffset   int
}

// RenderTableListView renders the table list view
func RenderTableListView(m TableListViewModel) string {
	// 2ペインレイアウト
	leftPaneWidth := 30 // 固定幅
	// rightPaneWidth = (borderの内側の幅) - (leftPaneWidth + leftPaneBorderRight)
	// = (m.Width - 2) - (30 + 1) = m.Width - 33
	rightPaneWidth := m.Width - leftPaneWidth - 3

	// ヘッダー
	// borderStyleの内側の幅 m.Width - 2 に合わせる
	// 右寄せで接続サーバ情報を表示
	rightText := "Connected to " + m.Endpoint

	// 使用可能な幅（パディング分を引く）
	availableWidth := m.Width - 4
	spaceBefore := availableWidth - len(rightText)
	if spaceBefore < 0 {
		spaceBefore = 0
	}

	headerContent := strings.Repeat(" ", spaceBefore) + rightText

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Padding(0, 1).
		Width(m.Width - 2)
	header := headerStyle.Render(headerContent)

	// 左ペイン: テーブルリスト
	// SelectableListを使用
	tableList := ui.SelectableList{
		Title:         fmt.Sprintf("Tables (%d)", len(m.Tables)),
		Items:         m.Tables,
		SelectedIndex: m.SelectedTable,
		Focused:       m.RightPaneMode == app.RightPaneModeSchema, // スキーマビューモードの時のみフォーカス
	}
	leftPaneContent := tableList.Render()

	// ボーダー色の決定
	var borderColor string
	if m.RightPaneMode == app.RightPaneModeList || m.RightPaneMode == app.RightPaneModeDetail {
		borderColor = "#666666"
	} else {
		borderColor = "#555555"
	}
	leftPaneStyle := lipgloss.NewStyle().
		Width(leftPaneWidth).
		Height(m.Height - 8). // タイトル行、ヘッダー、セパレーター×3、ステータス、フッター、ボーダー×2を除く
		BorderStyle(lipgloss.NormalBorder()).
		BorderRight(true).
		BorderForeground(lipgloss.Color(borderColor)).
		Padding(0, 1)
	leftPane := leftPaneStyle.Render(leftPaneContent)

	// 右ペイン: テーブル詳細またはデータ表示
	rightPaneContent := ""
	if len(m.Tables) > 0 && m.SelectedTable < len(m.Tables) {
		selectedTableName := m.Tables[m.SelectedTable]

		// モードに応じてヘッダーを変更
		if m.RightPaneMode == app.RightPaneModeList || m.RightPaneMode == app.RightPaneModeDetail {
			// グリッドビュー/レコードビューモード: SQLエリアを表示
			if data, exists := m.TableData[selectedTableName]; exists {
				// SQLエリアのスタイル（背景なし）
				sqlStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("#CCCCCC"))

				// SQLとセパレーターを手動で組み立て
				sqlText := sqlStyle.Render(data.DisplaySQL)
				separator := ui.Separator(rightPaneWidth - 2)

				rightPaneContent = sqlText + "\n" + separator
			}
		}

		if m.RightPaneMode == app.RightPaneModeSchema {
			// Schema表示モード
			var tableSchema *nosqldb.TableResult
			var indexes []nosqldb.IndexInfo
			if details, exists := m.TableDetails[selectedTableName]; exists && details != nil {
				tableSchema = details.Schema
				indexes = details.Indexes
			}
			rightPaneContent += RenderSchemaView(SchemaViewModel{
				TableName:      selectedTableName,
				AllTables:      m.Tables,
				TableSchema:    tableSchema,
				Indexes:        indexes,
				LoadingDetails: m.LoadingDetails,
			})
		} else if m.RightPaneMode == app.RightPaneModeList {
			// グリッドビューモード
			// rightPane全体の高さ(m.Height-8)からSQLエリア(2行)を引く
			rightPaneHeight := m.Height - 10

			// データの取得状態を確認
			data, exists := m.TableData[selectedTableName]
			var rows []map[string]interface{}
			var dataErr error
			var sql string
			if exists && data != nil {
				rows = data.Rows
				dataErr = data.Err
				sql = data.SQL
			}

			var tableSchema *nosqldb.TableResult
			if details, exists := m.TableDetails[selectedTableName]; exists && details != nil {
				tableSchema = details.Schema
			}

			rightPaneContent += RenderDataGridView(DataGridViewModel{
				Rows:             rows,
				TableSchema:      tableSchema,
				SelectedRow:      m.SelectedDataRow,
				HorizontalOffset: m.HorizontalOffset,
				ViewportOffset:   m.ViewportOffset,
				Width:            rightPaneWidth,
				Height:           rightPaneHeight,
				LoadingData:      m.LoadingData,
				Error:            dataErr,
				SQL:              sql,
			})
		} else if m.RightPaneMode == app.RightPaneModeDetail {
			// レコードビューモード
			// データの取得状態を確認
			data, exists := m.TableData[selectedTableName]
			var rows []map[string]interface{}
			var dataErr error
			if exists && data != nil {
				rows = data.Rows
				dataErr = data.Err
			}

			var tableSchema *nosqldb.TableResult
			if details, exists := m.TableDetails[selectedTableName]; exists && details != nil {
				tableSchema = details.Schema
			}

			rightPaneContent += RenderRecordView(RecordViewModel{
				Rows:        rows,
				TableSchema: tableSchema,
				SelectedRow: m.SelectedDataRow,
				LoadingData: m.LoadingData,
				Error:       dataErr,
			})
		}
	}

	rightPaneStyle := lipgloss.NewStyle().
		Width(rightPaneWidth).
		Height(m.Height - 8).
		Padding(0, 1)
	// 末尾の空行を削除
	rightPaneContent = strings.TrimSuffix(rightPaneContent, "\n")
	rightPane := rightPaneStyle.Render(rightPaneContent)

	// 2ペインを横に並べる
	panes := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)

	// ステータスバー
	statusBarStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CCCCCC")).
		Padding(0, 1).
		Width(m.Width - 2)
	var status string
	if len(m.Tables) > 0 {
		selectedTableName := m.Tables[m.SelectedTable]
		if m.RightPaneMode == app.RightPaneModeList || m.RightPaneMode == app.RightPaneModeDetail {
			// グリッドビュー/レコードビューモード: テーブル名と行数を表示
			if data, exists := m.TableData[selectedTableName]; exists {
				if data.Err != nil {
					// エラーが発生した場合は赤色で表示
					errorStyle := lipgloss.NewStyle().
						Foreground(lipgloss.Color("#FF0000")).
						Padding(0, 1)
					status = errorStyle.Render(fmt.Sprintf("Error: %v", data.Err))
				} else if len(data.Rows) > 0 {
					totalRows := len(data.Rows)
					// データがまだある場合は "+" を追加
					moreIndicator := ""
					if data.HasMore {
						moreIndicator = "+"
					}
					// テーブル名と行数のみ表示
					status = statusBarStyle.Render(fmt.Sprintf("Table: %s (%d%s rows)", selectedTableName, totalRows, moreIndicator))
				} else {
					status = statusBarStyle.Render(fmt.Sprintf("Table: %s (0 rows)", selectedTableName))
				}
			} else if m.LoadingData {
				status = statusBarStyle.Render(fmt.Sprintf("Table: %s (loading...)", selectedTableName))
			} else {
				status = statusBarStyle.Render(fmt.Sprintf("Table: %s", selectedTableName))
			}
		} else {
			// スキーマビューモード: テーブル名のみ表示
			status = statusBarStyle.Render(fmt.Sprintf("Table: %s", selectedTableName))
		}
	} else {
		status = statusBarStyle.Render("")
	}

	// フッター
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Padding(0, 1).
		Width(m.Width - 2)
	var footer string
	if m.RightPaneMode == app.RightPaneModeList {
		footer = footerStyle.Render("j/k: Scroll  h/l: Scroll Left/Right  o: Detail  u: Back  q: Quit")
	} else if m.RightPaneMode == app.RightPaneModeDetail {
		footer = footerStyle.Render("j/k: Scroll  u: Back to List  q: Quit")
	} else {
		footer = footerStyle.Render("j/k: Navigate  o: View Data  u: Back  q: Quit")
	}

	// セパレーター
	topSeparator := ui.Separator(m.Width - 2)
	statusSeparator := ui.Separator(m.Width - 2)

	// 全体を組み立て
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		topSeparator,
		panes,
		statusSeparator,
		status,
		statusSeparator,
		footer,
	)

	// 手動でボーダーを描画
	borderStyleColor := lipgloss.NewStyle().Foreground(lipgloss.Color("#00D9FF"))

	// 上部ボーダー: ╭── Dito ─────...╮
	title := " Dito "
	// 全体の幅 = m.Width
	// "╭──" = 3文字, title = 6文字, "╮" = 1文字
	// 残りの "─" = m.Width - 3 - 6 - 1 = m.Width - 10
	topBorder := borderStyleColor.Render("╭──" + title + strings.Repeat("─", m.Width-10) + "╮")

	// 左右のボーダー文字
	leftBorder := borderStyleColor.Render("│")
	rightBorder := borderStyleColor.Render("│")

	// コンテンツの各行にボーダーを追加
	var result strings.Builder
	result.WriteString(topBorder + "\n")

	for _, line := range strings.Split(content, "\n") {
		result.WriteString(leftBorder + line + rightBorder + "\n")
	}

	// 下部ボーダー: ╰─────...╯
	bottomBorder := borderStyleColor.Render("╰" + strings.Repeat("─", m.Width-2) + "╯")
	result.WriteString(bottomBorder)

	return result.String()
}
