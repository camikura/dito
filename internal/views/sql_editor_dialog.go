package views

import (
	"strings"

	"github.com/camikura/dito/internal/ui"
	"github.com/charmbracelet/lipgloss"
)

// SQLEditorDialogViewModel represents the data needed to render the SQL editor dialog
type SQLEditorDialogViewModel struct {
	SQL       string
	CursorPos int
	Width     int
	Height    int
}

// RenderSQLEditorDialog renders the SQL editor dialog
func RenderSQLEditorDialog(m SQLEditorDialogViewModel) string {
	// ダイアログのサイズを計算
	dialogWidth := m.Width - 10
	if dialogWidth < 60 {
		dialogWidth = 60
	}
	if dialogWidth > 100 {
		dialogWidth = 100
	}

	dialogHeight := m.Height - 10
	if dialogHeight < 10 {
		dialogHeight = 10
	}
	if dialogHeight > 20 {
		dialogHeight = 20
	}

	// タイトル
	title := ui.StyleTitle.Render("SQL Editor")

	// ヘルプテキスト
	help := ui.StyleLabel.Render("Execute: ctrl+r | Cancel: <esc>")

	// SQL入力エリア
	sqlAreaHeight := dialogHeight - 4 // タイトル、ヘルプ、ボーダー分を引く
	sqlLines := strings.Split(m.SQL, "\n")

	// カーソル位置を計算（行と列）
	cursorLine := 0
	cursorCol := m.CursorPos
	currentPos := 0
	for i, line := range sqlLines {
		lineLen := len(line) + 1 // +1 for newline
		if currentPos+lineLen > m.CursorPos {
			cursorLine = i
			cursorCol = m.CursorPos - currentPos
			break
		}
		currentPos += lineLen
	}

	// SQL表示エリアを構築
	var sqlContent strings.Builder
	for i := 0; i < sqlAreaHeight; i++ {
		if i < len(sqlLines) {
			line := sqlLines[i]
			// カーソルを表示
			if i == cursorLine {
				if cursorCol <= len(line) {
					displayLine := line[:cursorCol] + "_" + line[cursorCol:]
					sqlContent.WriteString(displayLine + "\n")
				} else {
					sqlContent.WriteString(line + "_\n")
				}
			} else {
				sqlContent.WriteString(line + "\n")
			}
		} else {
			sqlContent.WriteString("\n")
		}
	}

	// ダイアログ内容を組み立て
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		sqlContent.String(),
		"",
		help,
	)

	// ボーダースタイル
	dialogStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.ColorPrimary).
		Padding(1, 2).
		Width(dialogWidth).
		Height(dialogHeight)

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
