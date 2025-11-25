package views

import (
	"github.com/camikura/dito/internal/ui"
	"github.com/charmbracelet/lipgloss"
)

// TextInputDialogViewModel represents the data needed to render the text input dialog
type TextInputDialogViewModel struct {
	Label     string
	Value     string
	CursorPos int
	Width     int
	Height    int
}

// RenderTextInputDialog renders the text input dialog
func RenderTextInputDialog(m TextInputDialogViewModel) string {
	// ダイアログのサイズを計算
	dialogWidth := m.Width - 20
	if dialogWidth < 40 {
		dialogWidth = 40
	}
	if dialogWidth > 80 {
		dialogWidth = 80
	}

	// タイトル
	title := ui.StyleTitle.Render(m.Label)

	// カーソルを表示した入力値
	displayValue := m.Value
	if m.CursorPos <= len(m.Value) {
		displayValue = m.Value[:m.CursorPos] + "_" + m.Value[m.CursorPos:]
	} else {
		displayValue = m.Value + "_"
	}

	// 入力フィールド
	inputStyle := lipgloss.NewStyle().
		Foreground(ui.ColorWhite).
		Background(lipgloss.Color("#444444")).
		Padding(0, 1).
		Width(dialogWidth - 4)

	inputField := inputStyle.Render(displayValue)

	// ヘルプテキスト
	help := ui.StyleLabel.Render("Enter: Save  ESC: Cancel")

	// ダイアログ内容を組み立て
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		inputField,
		"",
		help,
	)

	// ボーダースタイル
	dialogStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.ColorPrimary).
		Padding(1, 2).
		Width(dialogWidth)

	dialog := dialogStyle.Render(content)

	// 画面中央に配置
	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		dialog,
	)
}
