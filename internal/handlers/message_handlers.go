package handlers

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/camikura/dito/internal/app"
	"github.com/camikura/dito/internal/db"
	"github.com/camikura/dito/internal/views"
)

// HandleWindowSize handles window size change messages
func HandleWindowSize(m app.Model, msg tea.WindowSizeMsg) app.Model {
	m.Width = msg.Width
	m.Height = msg.Height

	// ビューポートサイズを画面の高さから計算
	// 右ペインの高さ (m.Height - 8) からヘッダー等を引く
	rightPaneHeight := m.Height - 8
	m.ViewportSize = rightPaneHeight - 3 // SQLエリア+ボーダー: 2行 + カラムヘッダー: 1行
	if m.ViewportSize < 1 {
		m.ViewportSize = 1
	}

	return m
}

// HandleKeyPress handles keyboard input by routing to appropriate screen handlers
func HandleKeyPress(m app.Model, msg tea.KeyMsg) (app.Model, tea.Cmd) {
	// ダイアログが表示されている場合はEnterまたはEscで閉じる
	if m.DialogVisible {
		if msg.String() == "enter" || msg.String() == "esc" {
			m.DialogVisible = false
			m.DialogTitle = ""
			m.DialogMessage = ""
		}
		return m, nil
	}

	switch m.Screen {
	case app.ScreenSelection:
		return HandleSelection(m, msg)
	case app.ScreenOnPremiseConfig:
		return HandleOnPremiseConfig(m, msg)
	case app.ScreenCloudConfig:
		return HandleCloudConfig(m, msg)
	case app.ScreenTableList:
		return HandleTableList(m, msg)
	default:
		return m, nil
	}
}

// HandleConnectionResult handles database connection result messages
func HandleConnectionResult(m app.Model, msg db.ConnectionResult) (app.Model, tea.Cmd) {
	// Statusをリセット
	m.OnPremiseConfig.Status = app.StatusDisconnected

	if msg.Err != nil {
		// エラーダイアログを表示
		m.DialogVisible = true
		m.DialogType = app.DialogTypeError
		m.DialogTitle = "Connection Failed"
		m.DialogMessage = msg.Err.Error()
	} else {
		// テスト接続の場合は成功ダイアログを表示
		if msg.IsTest {
			m.DialogVisible = true
			m.DialogType = app.DialogTypeSuccess
			m.DialogTitle = "Connection Successful"
			m.DialogMessage = fmt.Sprintf("Successfully connected to %s\nServer version: %s", msg.Endpoint, msg.Version)
			// テスト接続のクライアントは閉じる
			if msg.Client != nil {
				msg.Client.Close()
			}
		} else {
			// 実接続の場合はテーブル一覧を取得して画面遷移
			m.NosqlClient = msg.Client
			m.Endpoint = msg.Endpoint
			return m, db.FetchTables(msg.Client)
		}
	}
	return m, nil
}

// HandleTableListResult handles table list fetch result messages
func HandleTableListResult(m app.Model, msg db.TableListResult) (app.Model, tea.Cmd) {
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
}

// HandleTableDetailsResult handles table details fetch result messages
func HandleTableDetailsResult(m app.Model, msg db.TableDetailsResult) (app.Model, tea.Cmd) {
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
}

// HandleTableDataResult handles table data fetch result messages
func HandleTableDataResult(m app.Model, msg db.TableDataResult) (app.Model, tea.Cmd) {
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
