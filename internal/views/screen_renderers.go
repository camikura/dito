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

	// ステータスメッセージ（この画面では空）
	statusMessage := ""

	// ヘルプテキスト
	helpText := "Tab/Shift+Tab or ↑/↓: Navigate  Enter: Select  q: Quit"

	return renderWithBorder(vm.Width, vm.Height, content, statusMessage, helpText)
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

	// ステータスメッセージ
	statusMessage := buildConnectionStatusMessage(
		vm.Model.OnPremiseConfig.Status,
		vm.Model.OnPremiseConfig.ServerVersion,
		vm.Model.OnPremiseConfig.ErrorMsg,
		vm.Width,
	)

	// ヘルプテキスト
	helpText := "Tab/Shift+Tab: Navigate  Space: Toggle  Enter: Execute  Esc: Back  Ctrl+C: Quit"

	return renderWithBorder(vm.Width, vm.Height, content, statusMessage, helpText)
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

	// ステータスメッセージ
	statusMessage := buildConnectionStatusMessage(
		vm.Model.CloudConfig.Status,
		vm.Model.CloudConfig.ServerVersion,
		vm.Model.CloudConfig.ErrorMsg,
		vm.Width,
	)

	// ヘルプテキスト
	helpText := "Tab/Shift+Tab: Navigate  Space: Toggle  Enter: Execute  Esc: Back  Ctrl+C: Quit"

	return renderWithBorder(vm.Width, vm.Height, content, statusMessage, helpText)
}

// buildConnectionStatusMessage builds the status message for connection screens
func buildConnectionStatusMessage(status app.ConnectionStatus, serverVersion string, errorMsg string, width int) string {
	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00"))
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000"))

	var statusMessage string
	switch status {
	case app.StatusConnecting:
		statusMessage = statusStyle.Render("Connecting...")
	case app.StatusConnected:
		msg := "Connected"
		if serverVersion != "" {
			msg = serverVersion
		}
		statusMessage = statusStyle.Render(msg)
	case app.StatusError:
		msg := "Connection failed"
		if errorMsg != "" {
			errMsg := errorMsg
			maxWidth := width - 10
			if len(errMsg) > maxWidth {
				errMsg = errMsg[:maxWidth] + "..."
			}
			msg = errMsg
		}
		statusMessage = errorStyle.Render(msg)
	}
	return statusMessage
}

// renderWithBorder renders the content with border, status, and footer
func renderWithBorder(width, height int, content, statusMessage, helpText string) string {
	// 共通スタイル
	statusBarStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Padding(0, 1).
		Width(width - 2)

	// コンテンツを左寄せ
	contentHeight := height - 7 // タイトル行、空行、セパレーター×3、ステータスエリア、フッターを除く
	contentStyle := lipgloss.NewStyle().
		Width(width - 2).
		Height(contentHeight).
		AlignVertical(lipgloss.Top).
		AlignHorizontal(lipgloss.Left).
		PaddingLeft(1)

	leftAlignedContent := contentStyle.Render(content)

	// セパレーター
	separator := ui.Separator(width - 2)

	// ステータスエリア
	statusAreaStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Width(width - 2)
	statusArea := statusAreaStyle.Render(statusMessage)

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
	bottomBorder := borderStyleColor.Render("╰" + strings.Repeat("─", width-2) + "╯")
	result.WriteString(bottomBorder)

	return result.String()
}
