package handlers

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/camikura/dito/internal/app"
)

// HandleSelection handles the connection selection screen input
func HandleSelection(m app.Model, msg tea.KeyMsg) (app.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		m.Cursor--
		if m.Cursor < 0 {
			m.Cursor = len(m.Choices) - 1
		}
	case "down", "j":
		m.Cursor = (m.Cursor + 1) % len(m.Choices)
	case "enter":
		// 0: Cloud, 1: On-Premise
		switch m.Cursor {
		case 0:
			// Cloud: 接続設定画面に遷移
			m.Screen = app.ScreenCloudConfig
			return m, nil
		case 1:
			// On-Premise: 接続設定画面に遷移
			m.Screen = app.ScreenOnPremiseConfig
			return m, nil
		}
	}
	return m, nil
}
