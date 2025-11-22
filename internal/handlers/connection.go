package handlers

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/camikura/dito/internal/app"
	"github.com/camikura/dito/internal/db"
)

// HandleCloudConfig handles the cloud connection configuration screen input
func HandleCloudConfig(m app.Model, msg tea.KeyMsg) (app.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		// エディション選択画面に戻る
		m.Screen = app.ScreenSelection
		return m, nil
	case "tab":
		// 次のフィールドへ
		m.CloudConfig.Focus = (m.CloudConfig.Focus + 1) % 8
		// テキスト入力フィールドの場合、カーソルを末尾に
		if m.CloudConfig.Focus == 0 {
			m.CloudConfig.CursorPos = len(m.CloudConfig.Region)
		} else if m.CloudConfig.Focus == 1 {
			m.CloudConfig.CursorPos = len(m.CloudConfig.Compartment)
		} else if m.CloudConfig.Focus == 5 {
			m.CloudConfig.CursorPos = len(m.CloudConfig.ConfigFile)
		}
		return m, nil
	case "shift+tab":
		// 前のフィールドへ
		m.CloudConfig.Focus--
		if m.CloudConfig.Focus < 0 {
			m.CloudConfig.Focus = 7
		}
		// テキスト入力フィールドの場合、カーソルを末尾に
		if m.CloudConfig.Focus == 0 {
			m.CloudConfig.CursorPos = len(m.CloudConfig.Region)
		} else if m.CloudConfig.Focus == 1 {
			m.CloudConfig.CursorPos = len(m.CloudConfig.Compartment)
		} else if m.CloudConfig.Focus == 5 {
			m.CloudConfig.CursorPos = len(m.CloudConfig.ConfigFile)
		}
		return m, nil
	case "enter":
		// ボタンが選択されている場合
		if m.CloudConfig.Focus == 6 {
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
	case "left":
		// カーソルを左に移動
		if m.CloudConfig.CursorPos > 0 {
			m.CloudConfig.CursorPos--
		}
		return m, nil
	case "right":
		// カーソルを右に移動
		var maxPos int
		if m.CloudConfig.Focus == 0 {
			maxPos = len(m.CloudConfig.Region)
		} else if m.CloudConfig.Focus == 1 {
			maxPos = len(m.CloudConfig.Compartment)
		} else if m.CloudConfig.Focus == 5 {
			maxPos = len(m.CloudConfig.ConfigFile)
		}
		if m.CloudConfig.CursorPos < maxPos {
			m.CloudConfig.CursorPos++
		}
		return m, nil
	case "backspace":
		// テキストフィールドの入力削除
		if m.CloudConfig.Focus == 0 && m.CloudConfig.CursorPos > 0 {
			m.CloudConfig.Region = m.CloudConfig.Region[:m.CloudConfig.CursorPos-1] + m.CloudConfig.Region[m.CloudConfig.CursorPos:]
			m.CloudConfig.CursorPos--
		} else if m.CloudConfig.Focus == 1 && m.CloudConfig.CursorPos > 0 {
			m.CloudConfig.Compartment = m.CloudConfig.Compartment[:m.CloudConfig.CursorPos-1] + m.CloudConfig.Compartment[m.CloudConfig.CursorPos:]
			m.CloudConfig.CursorPos--
		} else if m.CloudConfig.Focus == 5 && m.CloudConfig.CursorPos > 0 {
			m.CloudConfig.ConfigFile = m.CloudConfig.ConfigFile[:m.CloudConfig.CursorPos-1] + m.CloudConfig.ConfigFile[m.CloudConfig.CursorPos:]
			m.CloudConfig.CursorPos--
		}
		return m, nil
	default:
		// テキスト入力
		if len(msg.String()) == 1 {
			if m.CloudConfig.Focus == 0 {
				m.CloudConfig.Region = m.CloudConfig.Region[:m.CloudConfig.CursorPos] + msg.String() + m.CloudConfig.Region[m.CloudConfig.CursorPos:]
				m.CloudConfig.CursorPos++
			} else if m.CloudConfig.Focus == 1 {
				m.CloudConfig.Compartment = m.CloudConfig.Compartment[:m.CloudConfig.CursorPos] + msg.String() + m.CloudConfig.Compartment[m.CloudConfig.CursorPos:]
				m.CloudConfig.CursorPos++
			} else if m.CloudConfig.Focus == 5 {
				m.CloudConfig.ConfigFile = m.CloudConfig.ConfigFile[:m.CloudConfig.CursorPos] + msg.String() + m.CloudConfig.ConfigFile[m.CloudConfig.CursorPos:]
				m.CloudConfig.CursorPos++
			}
		}
		return m, nil
	}
}

// HandleOnPremiseConfig handles the on-premise connection configuration screen input
func HandleOnPremiseConfig(m app.Model, msg tea.KeyMsg) (app.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		// エディション選択画面に戻る
		m.Screen = app.ScreenSelection
		return m, nil
	case "tab":
		// 次のフィールドへ
		m.OnPremiseConfig.Focus = (m.OnPremiseConfig.Focus + 1) % 5
		// テキスト入力フィールドの場合、カーソルを末尾に
		if m.OnPremiseConfig.Focus == 0 {
			m.OnPremiseConfig.CursorPos = len(m.OnPremiseConfig.Endpoint)
		} else if m.OnPremiseConfig.Focus == 1 {
			m.OnPremiseConfig.CursorPos = len(m.OnPremiseConfig.Port)
		}
		return m, nil
	case "shift+tab":
		// 前のフィールドへ
		m.OnPremiseConfig.Focus--
		if m.OnPremiseConfig.Focus < 0 {
			m.OnPremiseConfig.Focus = 4
		}
		// テキスト入力フィールドの場合、カーソルを末尾に
		if m.OnPremiseConfig.Focus == 0 {
			m.OnPremiseConfig.CursorPos = len(m.OnPremiseConfig.Endpoint)
		} else if m.OnPremiseConfig.Focus == 1 {
			m.OnPremiseConfig.CursorPos = len(m.OnPremiseConfig.Port)
		}
		return m, nil
	case "enter":
		// ボタンが選択されている場合
		if m.OnPremiseConfig.Focus == 3 {
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
	case "left":
		// カーソルを左に移動
		if m.OnPremiseConfig.CursorPos > 0 {
			m.OnPremiseConfig.CursorPos--
		}
		return m, nil
	case "right":
		// カーソルを右に移動
		var maxPos int
		if m.OnPremiseConfig.Focus == 0 {
			maxPos = len(m.OnPremiseConfig.Endpoint)
		} else if m.OnPremiseConfig.Focus == 1 {
			maxPos = len(m.OnPremiseConfig.Port)
		}
		if m.OnPremiseConfig.CursorPos < maxPos {
			m.OnPremiseConfig.CursorPos++
		}
		return m, nil
	case "backspace":
		// テキストフィールドの入力削除
		if m.OnPremiseConfig.Focus == 0 && m.OnPremiseConfig.CursorPos > 0 {
			m.OnPremiseConfig.Endpoint = m.OnPremiseConfig.Endpoint[:m.OnPremiseConfig.CursorPos-1] + m.OnPremiseConfig.Endpoint[m.OnPremiseConfig.CursorPos:]
			m.OnPremiseConfig.CursorPos--
		} else if m.OnPremiseConfig.Focus == 1 && m.OnPremiseConfig.CursorPos > 0 {
			m.OnPremiseConfig.Port = m.OnPremiseConfig.Port[:m.OnPremiseConfig.CursorPos-1] + m.OnPremiseConfig.Port[m.OnPremiseConfig.CursorPos:]
			m.OnPremiseConfig.CursorPos--
		}
		return m, nil
	default:
		// テキスト入力
		if len(msg.String()) == 1 {
			if m.OnPremiseConfig.Focus == 0 {
				m.OnPremiseConfig.Endpoint = m.OnPremiseConfig.Endpoint[:m.OnPremiseConfig.CursorPos] + msg.String() + m.OnPremiseConfig.Endpoint[m.OnPremiseConfig.CursorPos:]
				m.OnPremiseConfig.CursorPos++
			} else if m.OnPremiseConfig.Focus == 1 {
				m.OnPremiseConfig.Port = m.OnPremiseConfig.Port[:m.OnPremiseConfig.CursorPos] + msg.String() + m.OnPremiseConfig.Port[m.OnPremiseConfig.CursorPos:]
				m.OnPremiseConfig.CursorPos++
			}
		}
		return m, nil
	}
}
