package handlers

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/camikura/dito/internal/app"
	"github.com/camikura/dito/internal/db"
)

// HandleCloudConfig handles the cloud connection configuration screen input
func HandleCloudConfig(m app.Model, msg tea.KeyMsg) (app.Model, tea.Cmd) {
	// テキスト入力ダイアログが表示されている場合は専用ハンドラーを呼び出す
	if m.TextInputVisible {
		return handleTextInputDialog(m, msg)
	}

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		// エディション選択画面に戻る
		m.Screen = app.ScreenSelection
		return m, nil
	case "up", "k":
		// 前のフィールドへ
		m.CloudConfig.Focus--
		if m.CloudConfig.Focus < 0 {
			m.CloudConfig.Focus = 7
		}
		return m, nil
	case "down", "j":
		// 次のフィールドへ
		m.CloudConfig.Focus = (m.CloudConfig.Focus + 1) % 8
		return m, nil
	case "enter":
		// テキストフィールドの場合はダイアログを開く
		if m.CloudConfig.Focus == 0 {
			m.TextInputVisible = true
			m.TextInputLabel = "Region"
			m.TextInputValue = m.CloudConfig.Region
			m.TextInputCursorPos = len(m.CloudConfig.Region)
			return m, nil
		} else if m.CloudConfig.Focus == 1 {
			m.TextInputVisible = true
			m.TextInputLabel = "Compartment"
			m.TextInputValue = m.CloudConfig.Compartment
			m.TextInputCursorPos = len(m.CloudConfig.Compartment)
			return m, nil
		} else if m.CloudConfig.Focus == 5 {
			m.TextInputVisible = true
			m.TextInputLabel = "Config File"
			m.TextInputValue = m.CloudConfig.ConfigFile
			m.TextInputCursorPos = len(m.CloudConfig.ConfigFile)
			return m, nil
		} else if m.CloudConfig.Focus == 6 {
			// 接続テスト - TODO: Cloud接続実装
			m.CloudConfig.Status = app.StatusConnecting
			m.CloudConfig.ErrorMsg = ""
			return m, nil
		} else if m.CloudConfig.Focus == 7 {
			// 接続する - TODO: Cloud接続実装
			m.CloudConfig.Status = app.StatusConnecting
			m.CloudConfig.ErrorMsg = ""
			return m, nil
		}
		return m, nil
	case " ":
		// ラジオボタンの選択
		if m.CloudConfig.Focus >= 2 && m.CloudConfig.Focus <= 4 {
			m.CloudConfig.AuthMethod = m.CloudConfig.Focus - 2
		}
		return m, nil
	}
	return m, nil
}

// HandleOnPremiseConfig handles the on-premise connection configuration screen input
func HandleOnPremiseConfig(m app.Model, msg tea.KeyMsg) (app.Model, tea.Cmd) {
	// テキスト入力ダイアログが表示されている場合は専用ハンドラーを呼び出す
	if m.TextInputVisible {
		return handleTextInputDialog(m, msg)
	}

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		// エディション選択画面に戻る
		m.Screen = app.ScreenSelection
		return m, nil
	case "up", "k":
		// 前のフィールドへ
		m.OnPremiseConfig.Focus--
		if m.OnPremiseConfig.Focus < 0 {
			m.OnPremiseConfig.Focus = 4
		}
		return m, nil
	case "down", "j":
		// 次のフィールドへ
		m.OnPremiseConfig.Focus = (m.OnPremiseConfig.Focus + 1) % 5
		return m, nil
	case "enter":
		// テキストフィールドの場合はダイアログを開く
		if m.OnPremiseConfig.Focus == 0 {
			m.TextInputVisible = true
			m.TextInputLabel = "Endpoint"
			m.TextInputValue = m.OnPremiseConfig.Endpoint
			m.TextInputCursorPos = len(m.OnPremiseConfig.Endpoint)
			return m, nil
		} else if m.OnPremiseConfig.Focus == 1 {
			m.TextInputVisible = true
			m.TextInputLabel = "Port"
			m.TextInputValue = m.OnPremiseConfig.Port
			m.TextInputCursorPos = len(m.OnPremiseConfig.Port)
			return m, nil
		} else if m.OnPremiseConfig.Focus == 3 {
			// 接続テスト（テスト接続なので画面遷移しない）
			m.OnPremiseConfig.Status = app.StatusConnecting
			m.OnPremiseConfig.ErrorMsg = ""
			return m, db.Connect(m.OnPremiseConfig.Endpoint, m.OnPremiseConfig.Port, true)
		} else if m.OnPremiseConfig.Focus == 4 {
			// 接続する（実接続なのでテーブル一覧画面に遷移）
			m.OnPremiseConfig.Status = app.StatusConnecting
			m.OnPremiseConfig.ErrorMsg = ""
			return m, db.Connect(m.OnPremiseConfig.Endpoint, m.OnPremiseConfig.Port, false)
		}
		return m, nil
	case " ":
		// セキュアチェックボックスのトグル
		if m.OnPremiseConfig.Focus == 2 {
			m.OnPremiseConfig.Secure = !m.OnPremiseConfig.Secure
		}
		return m, nil
	}
	return m, nil
}

// handleTextInputDialog handles text input dialog input for connection forms
func handleTextInputDialog(m app.Model, msg tea.KeyMsg) (app.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit

	case tea.KeyEsc:
		// ダイアログを閉じる（変更を破棄）
		m.TextInputVisible = false
		return m, nil

	case tea.KeyEnter:
		// 変更を保存してダイアログを閉じる
		if m.Screen == app.ScreenOnPremiseConfig {
			if m.OnPremiseConfig.Focus == 0 {
				m.OnPremiseConfig.Endpoint = m.TextInputValue
			} else if m.OnPremiseConfig.Focus == 1 {
				m.OnPremiseConfig.Port = m.TextInputValue
			}
		} else if m.Screen == app.ScreenCloudConfig {
			if m.CloudConfig.Focus == 0 {
				m.CloudConfig.Region = m.TextInputValue
			} else if m.CloudConfig.Focus == 1 {
				m.CloudConfig.Compartment = m.TextInputValue
			} else if m.CloudConfig.Focus == 5 {
				m.CloudConfig.ConfigFile = m.TextInputValue
			}
		}
		m.TextInputVisible = false
		return m, nil

	case tea.KeyBackspace:
		if m.TextInputCursorPos > 0 {
			m.TextInputValue = m.TextInputValue[:m.TextInputCursorPos-1] + m.TextInputValue[m.TextInputCursorPos:]
			m.TextInputCursorPos--
		}
		return m, nil

	case tea.KeyDelete:
		if m.TextInputCursorPos < len(m.TextInputValue) {
			m.TextInputValue = m.TextInputValue[:m.TextInputCursorPos] + m.TextInputValue[m.TextInputCursorPos+1:]
		}
		return m, nil

	case tea.KeyLeft:
		if m.TextInputCursorPos > 0 {
			m.TextInputCursorPos--
		}
		return m, nil

	case tea.KeyRight:
		if m.TextInputCursorPos < len(m.TextInputValue) {
			m.TextInputCursorPos++
		}
		return m, nil

	case tea.KeyHome, tea.KeyCtrlA:
		m.TextInputCursorPos = 0
		return m, nil

	case tea.KeyEnd, tea.KeyCtrlE:
		m.TextInputCursorPos = len(m.TextInputValue)
		return m, nil

	case tea.KeySpace:
		// スペース入力
		m.TextInputValue = m.TextInputValue[:m.TextInputCursorPos] + " " + m.TextInputValue[m.TextInputCursorPos:]
		m.TextInputCursorPos++
		return m, nil

	case tea.KeyRunes:
		// 通常の文字入力
		runes := msg.Runes
		for _, r := range runes {
			m.TextInputValue = m.TextInputValue[:m.TextInputCursorPos] + string(r) + m.TextInputValue[m.TextInputCursorPos:]
			m.TextInputCursorPos++
		}
		return m, nil
	}

	return m, nil
}
