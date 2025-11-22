package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oracle/nosql-go-sdk/nosqldb"

	"github.com/camikura/dito/internal/app"
	"github.com/camikura/dito/internal/db"
	"github.com/camikura/dito/internal/handlers"
	"github.com/camikura/dito/internal/ui"
	"github.com/camikura/dito/internal/views"
)

// Message type aliases for db package types
type connectionResultMsg = db.ConnectionResult
type tableListResultMsg = db.TableListResult
type tableDetailsResultMsg = db.TableDetailsResult
type tableDataResultMsg = db.TableDataResult

// model wraps app.Model to allow methods in main package
type model struct {
	app.Model
}

// Initメソッド
func (m model) Init() tea.Cmd {
	return nil
}

// Updateメソッド
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		// ビューポートサイズを画面の高さから計算
		// 右ペインの高さ (m.Height - 8) からヘッダー等を引く
		rightPaneHeight := m.Height - 8
		m.ViewportSize = rightPaneHeight - 3 // SQLエリア+ボーダー: 2行 + カラムヘッダー: 1行
		if m.ViewportSize < 1 {
			m.ViewportSize = 1
		}
		return m, nil
	case tea.KeyMsg:
		switch m.Screen {
		case app.ScreenSelection:
			newModel, cmd := handlers.HandleSelection(m.Model, msg)
			m.Model = newModel
			return m, cmd
		case app.ScreenOnPremiseConfig:
			newModel, cmd := handlers.HandleOnPremiseConfig(m.Model, msg)
			m.Model = newModel
			return m, cmd
		case app.ScreenCloudConfig:
			newModel, cmd := handlers.HandleCloudConfig(m.Model, msg)
			m.Model = newModel
			return m, cmd
		case app.ScreenTableList:
			newModel, cmd := handlers.HandleTableList(m.Model, msg)
			m.Model = newModel
			return m, cmd
		}
	case connectionResultMsg:
		// 接続結果を処理
		if msg.Err != nil {
			m.OnPremiseConfig.Status = app.StatusError
			m.OnPremiseConfig.ErrorMsg = msg.Err.Error()
		} else {
			m.OnPremiseConfig.Status = app.StatusConnected
			m.OnPremiseConfig.ServerVersion = msg.Version
			m.OnPremiseConfig.ErrorMsg = ""

			// テスト接続でない場合のみテーブル一覧を取得して画面遷移
			if !msg.IsTest {
				// クライアントとエンドポイントを保存
				m.NosqlClient = msg.Client
				m.Endpoint = msg.Endpoint
				// テーブル一覧を取得
				return m, db.FetchTables(msg.Client)
			}
		}
		return m, nil
	case tableListResultMsg:
		// テーブル一覧取得結果を処理
		if msg.Err != nil {
			m.OnPremiseConfig.Status = app.StatusError
			m.OnPremiseConfig.ErrorMsg = fmt.Sprintf("Failed to fetch tables: %v", msg.Err)
		} else {
			m.Tables = msg.Tables
			m.SelectedTable = 0
			// テーブル一覧画面に遷移
			m.Screen = app.ScreenTableList
			// 最初のテーブルの詳細を取得
			if len(m.Tables) > 0 {
				return m, db.FetchTableDetails(m.NosqlClient, m.Tables[0])
			}
		}
		return m, nil
	case tableDetailsResultMsg:
		// テーブル詳細取得結果を処理
		if msg.Err == nil {
			m.TableDetails[msg.TableName] = &msg
		}
		m.LoadingDetails = false

		// グリッドビューモードで、このテーブルのデータがまだ取得されていない場合は取得
		if m.RightPaneMode == app.RightPaneModeList && msg.Err == nil {
			if _, exists := m.TableData[msg.TableName]; !exists {
				m.LoadingData = true
				primaryKeys := views.ParsePrimaryKeysFromDDL(msg.Schema.DDL)
				return m, db.FetchTableData(m.NosqlClient, msg.TableName, m.FetchSize, primaryKeys)
			}
		}
		return m, nil
	case tableDataResultMsg:
		// テーブルデータ取得結果を処理
		if msg.Err == nil {
			if msg.IsAppend {
				// 既存データに追加（SQLは更新しない）
				if existingData, exists := m.TableData[msg.TableName]; exists {
					existingData.Rows = append(existingData.Rows, msg.Rows...)
					existingData.LastPKValues = msg.LastPKValues
					existingData.HasMore = msg.HasMore
				}
			} else {
				// 新規データとして設定
				m.TableData[msg.TableName] = &msg
			}
		} else {
			// エラーの場合もデータを保存（エラーメッセージとSQLを表示するため）
			m.TableData[msg.TableName] = &msg
		}
		m.LoadingData = false
		return m, nil
	}

	return m, nil
}


// Viewメソッド
func (m model) View() string {
	if m.Width == 0 {
		return "Loading..."
	}

	// 共通スタイル
	statusBarStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Padding(0, 1).
		Width(m.Width - 2)

	// メインコンテンツ
	var content string
	switch m.Screen {
	case app.ScreenSelection:
		content = views.RenderConnectionSelection(views.ConnectionSelectionModel{
			Choices: m.Choices,
			Cursor:  m.Cursor,
		})
	case app.ScreenOnPremiseConfig:
		content = views.RenderOnPremiseForm(views.OnPremiseFormModel{
			Endpoint:  m.OnPremiseConfig.Endpoint,
			Port:      m.OnPremiseConfig.Port,
			Secure:    m.OnPremiseConfig.Secure,
			Focus:     m.OnPremiseConfig.Focus,
			CursorPos: m.OnPremiseConfig.CursorPos,
		})
	case app.ScreenCloudConfig:
		content = views.RenderCloudForm(views.CloudFormModel{
			Region:      m.CloudConfig.Region,
			Compartment: m.CloudConfig.Compartment,
			AuthMethod:  m.CloudConfig.AuthMethod,
			ConfigFile:  m.CloudConfig.ConfigFile,
			Focus:       m.CloudConfig.Focus,
			CursorPos:   m.CloudConfig.CursorPos,
		})
	case app.ScreenTableList:
		return m.viewTableList() // テーブル一覧は独自レイアウト
	default:
		content = "Unknown screen"
	}

	// コンテンツを左寄せ
	contentHeight := m.Height - 7 // タイトル行、空行、セパレーター×3、ステータスエリア、フッターを除く
	contentStyle := lipgloss.NewStyle().
		Width(m.Width - 2).
		Height(contentHeight).
		AlignVertical(lipgloss.Top).
		AlignHorizontal(lipgloss.Left).
		PaddingLeft(1)

	leftAlignedContent := contentStyle.Render(content)

	// セパレーター
	separator := ui.Separator(m.Width - 2)

	// ステータス表示エリア（1行）
	var statusMessage string
	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00"))
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000"))

	if m.Screen == app.ScreenOnPremiseConfig {
		switch m.OnPremiseConfig.Status {
		case app.StatusConnecting:
			statusMessage = statusStyle.Render("Connecting...")
		case app.StatusConnected:
			msg := "Connected"
			if m.OnPremiseConfig.ServerVersion != "" {
				msg = m.OnPremiseConfig.ServerVersion
			}
			statusMessage = statusStyle.Render(msg)
		case app.StatusError:
			msg := "Connection failed"
			if m.OnPremiseConfig.ErrorMsg != "" {
				errMsg := m.OnPremiseConfig.ErrorMsg
				maxWidth := m.Width - 10
				if len(errMsg) > maxWidth {
					errMsg = errMsg[:maxWidth] + "..."
				}
				msg = errMsg
			}
			statusMessage = errorStyle.Render(msg)
		}
	} else if m.Screen == app.ScreenCloudConfig {
		switch m.CloudConfig.Status {
		case app.StatusConnecting:
			statusMessage = statusStyle.Render("Connecting...")
		case app.StatusConnected:
			msg := "Connected"
			if m.CloudConfig.ServerVersion != "" {
				msg = m.CloudConfig.ServerVersion
			}
			statusMessage = statusStyle.Render(msg)
		case app.StatusError:
			msg := "Connection failed"
			if m.CloudConfig.ErrorMsg != "" {
				errMsg := m.CloudConfig.ErrorMsg
				maxWidth := m.Width - 10
				if len(errMsg) > maxWidth {
					errMsg = errMsg[:maxWidth] + "..."
				}
				msg = errMsg
			}
			statusMessage = errorStyle.Render(msg)
		}
	}

	statusAreaStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Width(m.Width - 2)
	statusArea := statusAreaStyle.Render(statusMessage)

	// フッター（ヘルプテキスト）
	var helpText string
	switch m.Screen {
	case app.ScreenSelection:
		helpText = "Tab/Shift+Tab or ↑/↓: Navigate  Enter: Select  q: Quit"
	case app.ScreenOnPremiseConfig:
		helpText = "Tab/Shift+Tab: Navigate  Space: Toggle  Enter: Execute  Esc: Back  Ctrl+C: Quit"
	case app.ScreenCloudConfig:
		helpText = "Tab/Shift+Tab: Navigate  Space: Toggle  Enter: Execute  Esc: Back  Ctrl+C: Quit"
	}
	footer := statusBarStyle.Render(helpText)

	// 全体を組み立て（手動でボーダーを描画）
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

	// タイトル行の下に空行を追加
	emptyLine := strings.Repeat(" ", m.Width-2)
	result.WriteString(leftBorder + emptyLine + rightBorder + "\n")

	// コンテンツを行ごとに分割してボーダーを追加
	lines := []string{
		leftAlignedContent,
		separator,
		statusArea,
		separator,
		footer,
	}

	for _, line := range lines {
		for _, l := range strings.Split(line, "\n") {
			if l != "" {
				result.WriteString(leftBorder + l + rightBorder + "\n")
			}
		}
	}

	// 下部ボーダー: ╰─────...╯
	bottomBorder := borderStyleColor.Render("╰" + strings.Repeat("─", m.Width-2) + "╯")
	result.WriteString(bottomBorder)

	return result.String()
}

// テーブル一覧画面のView
func (m model) viewTableList() string {
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
			rightPaneContent += views.RenderSchemaView(views.SchemaViewModel{
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

			rightPaneContent += views.RenderDataGridView(views.DataGridViewModel{
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

			rightPaneContent += views.RenderRecordView(views.RecordViewModel{
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

func main() {
	p := tea.NewProgram(
		model{Model: app.InitialModel()},
		tea.WithAltScreen(),       // 全画面モード
		tea.WithMouseCellMotion(), // マウスサポート（オプション）
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
