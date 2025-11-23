package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/camikura/dito/internal/app"
	"github.com/camikura/dito/internal/ui"
)

// ScreenViewModel represents the model for full screen rendering
type ScreenViewModel struct {
	Width  int
	Height int
	Model  app.Model
}

// RenderSelectionScreen renders the connection selection screen with border
func RenderSelectionScreen(vm ScreenViewModel) string {
	// メインコンテンツ
	content := RenderConnectionSelection(ConnectionSelectionModel{
		Choices: vm.Model.Choices,
		Cursor:  vm.Model.Cursor,
	})

	// ヘルプテキスト
	helpText := "Tab/Shift+Tab or ↑/↓: Navigate  Enter: Select  q: Quit"

	baseScreen := renderWithBorder(vm.Width, vm.Height, content, helpText)

	// ダイアログが表示されている場合は重ねて表示
	if vm.Model.DialogVisible {
		dialog := RenderDialog(vm.Width, vm.Height, DialogModel{
			Type:    vm.Model.DialogType,
			Title:   vm.Model.DialogTitle,
			Message: vm.Model.DialogMessage,
		})
		// ベース画面の上にダイアログを重ねる（簡易的な実装）
		return overlayDialog(baseScreen, dialog)
	}

	return baseScreen
}

// RenderOnPremiseConfigScreen renders the on-premise configuration screen with border
func RenderOnPremiseConfigScreen(vm ScreenViewModel) string {
	// メインコンテンツ
	content := RenderOnPremiseForm(OnPremiseFormModel{
		Endpoint:  vm.Model.OnPremiseConfig.Endpoint,
		Port:      vm.Model.OnPremiseConfig.Port,
		Secure:    vm.Model.OnPremiseConfig.Secure,
		Focus:     vm.Model.OnPremiseConfig.Focus,
		CursorPos: vm.Model.OnPremiseConfig.CursorPos,
	})

	// ヘルプテキスト
	helpText := "Tab/Shift+Tab: Navigate  Space: Toggle  Enter: Execute  Esc: Back  Ctrl+C: Quit"

	baseScreen := renderWithBorder(vm.Width, vm.Height, content, helpText)

	// ダイアログが表示されている場合は重ねて表示
	if vm.Model.DialogVisible {
		dialog := RenderDialog(vm.Width, vm.Height, DialogModel{
			Type:    vm.Model.DialogType,
			Title:   vm.Model.DialogTitle,
			Message: vm.Model.DialogMessage,
		})
		// ベース画面の上にダイアログを重ねる（簡易的な実装）
		return overlayDialog(baseScreen, dialog)
	}

	return baseScreen
}

// RenderCloudConfigScreen renders the cloud configuration screen with border
func RenderCloudConfigScreen(vm ScreenViewModel) string {
	// メインコンテンツ
	content := RenderCloudForm(CloudFormModel{
		Region:      vm.Model.CloudConfig.Region,
		Compartment: vm.Model.CloudConfig.Compartment,
		AuthMethod:  vm.Model.CloudConfig.AuthMethod,
		ConfigFile:  vm.Model.CloudConfig.ConfigFile,
		Focus:       vm.Model.CloudConfig.Focus,
		CursorPos:   vm.Model.CloudConfig.CursorPos,
	})

	// ヘルプテキスト
	helpText := "Tab/Shift+Tab: Navigate  Space: Toggle  Enter: Execute  Esc: Back  Ctrl+C: Quit"

	baseScreen := renderWithBorder(vm.Width, vm.Height, content, helpText)

	// ダイアログが表示されている場合は重ねて表示
	if vm.Model.DialogVisible {
		dialog := RenderDialog(vm.Width, vm.Height, DialogModel{
			Type:    vm.Model.DialogType,
			Title:   vm.Model.DialogTitle,
			Message: vm.Model.DialogMessage,
		})
		// ベース画面の上にダイアログを重ねる（簡易的な実装）
		return overlayDialog(baseScreen, dialog)
	}

	return baseScreen
}

// overlayDialog overlays a dialog on top of a base screen (simple implementation)
func overlayDialog(baseScreen, dialog string) string {
	// ダイアログは既に中央配置されているので、そのまま返す
	return dialog
}

// renderWithBorder renders the content with border and footer
func renderWithBorder(width, height int, content, helpText string) string {
	// 共通スタイル
	statusBarStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Padding(0, 1).
		Width(width - 2)

	// コンテンツを左寄せ
	contentHeight := height - 5 // タイトル行、空行、セパレーター×1、フッターを除く
	contentStyle := lipgloss.NewStyle().
		Width(width - 2).
		Height(contentHeight).
		AlignVertical(lipgloss.Top).
		AlignHorizontal(lipgloss.Left).
		PaddingLeft(1)

	leftAlignedContent := contentStyle.Render(content)

	// セパレーター
	separator := ui.Separator(width - 2)

	// フッター
	footer := statusBarStyle.Render(helpText)

	// 全体を組み立て（手動でボーダーを描画）
	borderStyleColor := lipgloss.NewStyle().Foreground(lipgloss.Color("#00D9FF"))

	// 上部ボーダー: ╭── Dito ─────...╮
	title := " Dito "
	// 全体の幅 = width
	// "╭──" = 3文字, title = 6文字, "╮" = 1文字
	// 残りの "─" = width - 3 - 6 - 1 = width - 10
	topBorder := borderStyleColor.Render("╭──" + title + strings.Repeat("─", width-10) + "╮")

	// 左右のボーダー文字
	leftBorder := borderStyleColor.Render("│")
	rightBorder := borderStyleColor.Render("│")

	// コンテンツの各行にボーダーを追加
	var result strings.Builder
	result.WriteString(topBorder + "\n")

	// タイトル行の下に空行を追加
	emptyLine := strings.Repeat(" ", width-2)
	result.WriteString(leftBorder + emptyLine + rightBorder + "\n")

	// コンテンツを行ごとに分割してボーダーを追加
	lines := []string{
		leftAlignedContent,
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
	bottomBorder := borderStyleColor.Render("╰" + strings.Repeat("─", width-2) + "╯")
	result.WriteString(bottomBorder)

	return result.String()
}
