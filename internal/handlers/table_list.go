package handlers

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/camikura/dito/internal/app"
	"github.com/camikura/dito/internal/db"
	"github.com/camikura/dito/internal/views"
)

// HandleTableList handles the table list screen input
func HandleTableList(m app.Model, msg tea.KeyMsg) (app.Model, tea.Cmd) {
	oldSelection := m.SelectedTable

	switch msg.String() {
	case "ctrl+c", "q":
		// クライアントをクローズしてから終了
		if m.NosqlClient != nil {
			m.NosqlClient.Close()
		}
		return m, tea.Quit
	case "up", "k":
		if m.RightPaneMode == app.RightPaneModeList || m.RightPaneMode == app.RightPaneModeDetail {
			// グリッドビュー/レコードビュー: データ行を選択
			if m.SelectedDataRow > 0 {
				m.SelectedDataRow--
				// ビューポートを調整（選択行がビューポートの上端より上になった場合）
				if m.SelectedDataRow < m.ViewportOffset {
					m.ViewportOffset = m.SelectedDataRow
				}
			}
		} else {
			// スキーマビューモード: テーブルを選択
			if m.SelectedTable > 0 {
				m.SelectedTable--
			}
		}
	case "down", "j":
		if m.RightPaneMode == app.RightPaneModeList || m.RightPaneMode == app.RightPaneModeDetail {
			// グリッドビュー/レコードビュー: データ行を選択
			tableName := m.Tables[m.SelectedTable]
			if data, exists := m.TableData[tableName]; exists && data.Err == nil {
				totalRows := len(data.Rows)
				if m.SelectedDataRow < totalRows-1 {
					m.SelectedDataRow++
					// ビューポートを調整（選択行がビューポートの下端より下になった場合）
					if m.SelectedDataRow >= m.ViewportOffset+m.ViewportSize {
						m.ViewportOffset = m.SelectedDataRow - m.ViewportSize + 1
					}

					// 残り10行以内まで来たら、さらにデータがある場合は追加取得
					remainingRows := totalRows - m.SelectedDataRow - 1
					if remainingRows <= 10 && data.HasMore && !m.LoadingData {
						m.LoadingData = true
						// PRIMARY KEYを取得
						var primaryKeys []string
						if details, exists := m.TableDetails[tableName]; exists && details.Schema != nil && details.Schema.DDL != "" {
							primaryKeys = views.ParsePrimaryKeysFromDDL(details.Schema.DDL)
						}
						return m, db.FetchMoreTableData(m.NosqlClient, tableName, m.FetchSize, primaryKeys, data.LastPKValues)
					}
				}
			}
		} else {
			// スキーマビューモード: テーブルを選択
			if m.SelectedTable < len(m.Tables)-1 {
				m.SelectedTable++
			}
		}
	case "h", "left":
		// データビューモード: 左にスクロール
		if m.RightPaneMode == app.RightPaneModeList {
			if m.HorizontalOffset > 0 {
				m.HorizontalOffset--
			}
		}
	case "l", "right":
		// データビューモード: 右にスクロール
		if m.RightPaneMode == app.RightPaneModeList {
			tableName := m.Tables[m.SelectedTable]
			// カラム数を取得
			var totalColumns int
			if details, exists := m.TableDetails[tableName]; exists && details.Schema != nil && details.Schema.DDL != "" {
				primaryKeys := views.ParsePrimaryKeysFromDDL(details.Schema.DDL)
				columns := views.ParseColumnsFromDDL(details.Schema.DDL, primaryKeys)
				totalColumns = len(columns)
			} else if data, exists := m.TableData[tableName]; exists && len(data.Rows) > 0 {
				totalColumns = len(data.Rows[0])
			}
			// 最後のカラムまでスクロールできるが、少なくとも1カラムは表示する
			if m.HorizontalOffset < totalColumns-1 {
				m.HorizontalOffset++
			}
		}
	case "esc", "u":
		if m.RightPaneMode == app.RightPaneModeDetail {
			// レコードビュー → グリッドビュー
			m.RightPaneMode = app.RightPaneModeList
			return m, nil
		} else if m.RightPaneMode == app.RightPaneModeList {
			// グリッドビュー → スキーマビュー
			m.RightPaneMode = app.RightPaneModeSchema
			m.HorizontalOffset = 0 // 横スクロールをリセット
			return m, nil
		}
		// スキーマビュー → 接続設定画面に戻る
		// 接続状態をリセット
		m.Screen = app.ScreenOnPremiseConfig
		m.OnPremiseConfig.Status = app.StatusDisconnected
		m.OnPremiseConfig.ErrorMsg = ""
		m.OnPremiseConfig.ServerVersion = ""
		return m, nil
	case "enter", "o":
		if m.RightPaneMode == app.RightPaneModeSchema {
			// スキーマビュー → グリッドビュー
			m.RightPaneMode = app.RightPaneModeList
			m.SelectedDataRow = 0   // 行選択をリセット
			m.ViewportOffset = 0    // ビューポートをリセット
			m.HorizontalOffset = 0  // 横スクロールをリセット
			// データ表示モードに切り替えたとき、データとテーブル詳細を取得
			if len(m.Tables) > 0 {
				tableName := m.Tables[m.SelectedTable]

				// テーブル詳細がまだ取得されていない場合は取得
				var cmds []tea.Cmd
				if _, exists := m.TableDetails[tableName]; !exists {
					m.LoadingDetails = true
					cmds = append(cmds, db.FetchTableDetails(m.NosqlClient, tableName))
				}

				// データがまだ取得されていない場合は取得
				if _, exists := m.TableData[tableName]; !exists {
					m.LoadingData = true
					// PRIMARY KEYを取得（テーブル詳細があれば）
					var primaryKeys []string
					if details, exists := m.TableDetails[tableName]; exists && details.Schema != nil && details.Schema.DDL != "" {
						primaryKeys = views.ParsePrimaryKeysFromDDL(details.Schema.DDL)
					}
					cmds = append(cmds, db.FetchTableData(m.NosqlClient, tableName, m.FetchSize, primaryKeys))
				}

				if len(cmds) > 0 {
					return m, tea.Batch(cmds...)
				}
			}
		} else if m.RightPaneMode == app.RightPaneModeList {
			// グリッドビュー → レコードビュー
			m.RightPaneMode = app.RightPaneModeDetail
		}
	}

	// テーブル選択が変わった場合、詳細を取得（スキーマビューモードのみ）
	if m.RightPaneMode == app.RightPaneModeSchema && oldSelection != m.SelectedTable && len(m.Tables) > 0 {
		tableName := m.Tables[m.SelectedTable]
		// まだ取得していないテーブルの場合のみ取得
		if _, exists := m.TableDetails[tableName]; !exists {
			m.LoadingDetails = true
			return m, db.FetchTableDetails(m.NosqlClient, tableName)
		}
	}

	return m, nil
}
