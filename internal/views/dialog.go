package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/camikura/dito/internal/app"
)

// DialogModel represents the model for dialog rendering
type DialogModel struct {
	Type    app.DialogType
	Title   string
	Message string
}

// RenderDialog renders a dialog box overlay
func RenderDialog(width, height int, dm DialogModel) string {
	// ダイアログのサイズ
	dialogWidth := 60
	if width < 70 {
		dialogWidth = width - 10
	}

	// 背景スタイル（半透明効果は難しいので、シンプルなボックス）
	var borderColor, titleColor lipgloss.Color
	if dm.Type == app.DialogTypeSuccess {
		borderColor = lipgloss.Color("#00FF00")
		titleColor = lipgloss.Color("#00FF00")
	} else {
		borderColor = lipgloss.Color("#FF0000")
		titleColor = lipgloss.Color("#FF0000")
	}

	dialogStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(dialogWidth).
		Padding(1, 2).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)

	titleStyle := lipgloss.NewStyle().
		Foreground(titleColor).
		Bold(true).
		AlignHorizontal(lipgloss.Center).
		Width(dialogWidth - 4)

	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		AlignHorizontal(lipgloss.Center).
		Width(dialogWidth - 4).
		MaxWidth(dialogWidth - 4)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		AlignHorizontal(lipgloss.Center).
		Width(dialogWidth - 4)

	// メッセージを折り返し
	wrappedMessage := wrapText(dm.Message, dialogWidth-8)

	// コンテンツを構築
	var content strings.Builder
	content.WriteString(titleStyle.Render(dm.Title))
	content.WriteString("\n\n")
	content.WriteString(messageStyle.Render(wrappedMessage))
	content.WriteString("\n\n")
	content.WriteString(helpStyle.Render("Press Enter or Esc to close"))

	dialog := dialogStyle.Render(content.String())

	// ダイアログを画面の中央に配置
	dialogPlacement := lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		dialog,
	)

	return dialogPlacement
}

// wrapText wraps text to fit within the specified width
func wrapText(text string, width int) string {
	if len(text) <= width {
		return text
	}

	var result strings.Builder
	words := strings.Fields(text)
	lineLen := 0

	for i, word := range words {
		wordLen := len(word)
		if lineLen+wordLen+1 > width {
			result.WriteString("\n")
			result.WriteString(word)
			lineLen = wordLen
		} else {
			if i > 0 {
				result.WriteString(" ")
				lineLen++
			}
			result.WriteString(word)
			lineLen += wordLen
		}
	}

	return result.String()
}
