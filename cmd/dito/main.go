package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/camikura/dito/internal/app"
	"github.com/camikura/dito/internal/db"
	"github.com/camikura/dito/internal/handlers"
	"github.com/camikura/dito/internal/ui"
	"github.com/camikura/dito/internal/views"
)

// Message type aliases for db package types
type connectionResultMsg = db.ConnectionResult
type tableListResultMsg = db.TableListResult
type tableDetailsResultMsg = db.TableDetailsResult
type tableDataResultMsg = db.TableDataResult

// model wraps app.Model to allow methods in main package
type model struct {
	app.Model
}

// Initメソッド
func (m model) Init() tea.Cmd {
	return nil
}

// Updateメソッド
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Model = handlers.HandleWindowSize(m.Model, msg)
	case tea.KeyMsg:
		m.Model, cmd = handlers.HandleKeyPress(m.Model, msg)
	case connectionResultMsg:
		m.Model, cmd = handlers.HandleConnectionResult(m.Model, msg)
	case tableListResultMsg:
		m.Model, cmd = handlers.HandleTableListResult(m.Model, msg)
	case tableDetailsResultMsg:
		m.Model, cmd = handlers.HandleTableDetailsResult(m.Model, msg)
	case tableDataResultMsg:
		m.Model, cmd = handlers.HandleTableDataResult(m.Model, msg)
	}

	return m, cmd
}


// Viewメソッド
func (m model) View() string {
	if m.Width == 0 {
		return "Loading..."
	}

	// 共通スタイル
	statusBarStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Padding(0, 1).
		Width(m.Width - 2)

	// メインコンテンツ
	var content string
	switch m.Screen {
	case app.ScreenSelection:
		content = views.RenderConnectionSelection(views.ConnectionSelectionModel{
			Choices: m.Choices,
			Cursor:  m.Cursor,
		})
	case app.ScreenOnPremiseConfig:
		content = views.RenderOnPremiseForm(views.OnPremiseFormModel{
			Endpoint:  m.OnPremiseConfig.Endpoint,
			Port:      m.OnPremiseConfig.Port,
			Secure:    m.OnPremiseConfig.Secure,
			Focus:     m.OnPremiseConfig.Focus,
			CursorPos: m.OnPremiseConfig.CursorPos,
		})
	case app.ScreenCloudConfig:
		content = views.RenderCloudForm(views.CloudFormModel{
			Region:      m.CloudConfig.Region,
			Compartment: m.CloudConfig.Compartment,
			AuthMethod:  m.CloudConfig.AuthMethod,
			ConfigFile:  m.CloudConfig.ConfigFile,
			Focus:       m.CloudConfig.Focus,
			CursorPos:   m.CloudConfig.CursorPos,
		})
	case app.ScreenTableList:
		return views.RenderTableListView(views.TableListViewModel{
			Width:            m.Width,
			Height:           m.Height,
			Endpoint:         m.Endpoint,
			Tables:           m.Tables,
			SelectedTable:    m.SelectedTable,
			RightPaneMode:    m.RightPaneMode,
			TableData:        m.TableData,
			TableDetails:     m.TableDetails,
			LoadingDetails:   m.LoadingDetails,
			LoadingData:      m.LoadingData,
			SelectedDataRow:  m.SelectedDataRow,
			HorizontalOffset: m.HorizontalOffset,
			ViewportOffset:   m.ViewportOffset,
		})
	default:
		content = "Unknown screen"
	}

	// コンテンツを左寄せ
	contentHeight := m.Height - 7 // タイトル行、空行、セパレーター×3、ステータスエリア、フッターを除く
	contentStyle := lipgloss.NewStyle().
		Width(m.Width - 2).
		Height(contentHeight).
		AlignVertical(lipgloss.Top).
		AlignHorizontal(lipgloss.Left).
		PaddingLeft(1)

	leftAlignedContent := contentStyle.Render(content)

	// セパレーター
	separator := ui.Separator(m.Width - 2)

	// ステータス表示エリア（1行）
	var statusMessage string
	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00"))
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000"))

	if m.Screen == app.ScreenOnPremiseConfig {
		switch m.OnPremiseConfig.Status {
		case app.StatusConnecting:
			statusMessage = statusStyle.Render("Connecting...")
		case app.StatusConnected:
			msg := "Connected"
			if m.OnPremiseConfig.ServerVersion != "" {
				msg = m.OnPremiseConfig.ServerVersion
			}
			statusMessage = statusStyle.Render(msg)
		case app.StatusError:
			msg := "Connection failed"
			if m.OnPremiseConfig.ErrorMsg != "" {
				errMsg := m.OnPremiseConfig.ErrorMsg
				maxWidth := m.Width - 10
				if len(errMsg) > maxWidth {
					errMsg = errMsg[:maxWidth] + "..."
				}
				msg = errMsg
			}
			statusMessage = errorStyle.Render(msg)
		}
	} else if m.Screen == app.ScreenCloudConfig {
		switch m.CloudConfig.Status {
		case app.StatusConnecting:
			statusMessage = statusStyle.Render("Connecting...")
		case app.StatusConnected:
			msg := "Connected"
			if m.CloudConfig.ServerVersion != "" {
				msg = m.CloudConfig.ServerVersion
			}
			statusMessage = statusStyle.Render(msg)
		case app.StatusError:
			msg := "Connection failed"
			if m.CloudConfig.ErrorMsg != "" {
				errMsg := m.CloudConfig.ErrorMsg
				maxWidth := m.Width - 10
				if len(errMsg) > maxWidth {
					errMsg = errMsg[:maxWidth] + "..."
				}
				msg = errMsg
			}
			statusMessage = errorStyle.Render(msg)
		}
	}

	statusAreaStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Width(m.Width - 2)
	statusArea := statusAreaStyle.Render(statusMessage)

	// フッター（ヘルプテキスト）
	var helpText string
	switch m.Screen {
	case app.ScreenSelection:
		helpText = "Tab/Shift+Tab or ↑/↓: Navigate  Enter: Select  q: Quit"
	case app.ScreenOnPremiseConfig:
		helpText = "Tab/Shift+Tab: Navigate  Space: Toggle  Enter: Execute  Esc: Back  Ctrl+C: Quit"
	case app.ScreenCloudConfig:
		helpText = "Tab/Shift+Tab: Navigate  Space: Toggle  Enter: Execute  Esc: Back  Ctrl+C: Quit"
	}
	footer := statusBarStyle.Render(helpText)

	// 全体を組み立て（手動でボーダーを描画）
	borderStyleColor := lipgloss.NewStyle().Foreground(lipgloss.Color("#00D9FF"))

	// 上部ボーダー: ╭── Dito ─────...╮
	title := " Dito "
	// 全体の幅 = m.Width
	// "╭──" = 3文字, title = 6文字, "╮" = 1文字
	// 残りの "─" = m.Width - 3 - 6 - 1 = m.Width - 10
	topBorder := borderStyleColor.Render("╭──" + title + strings.Repeat("─", m.Width-10) + "╮")

	// 左右のボーダー文字
	leftBorder := borderStyleColor.Render("│")
	rightBorder := borderStyleColor.Render("│")

	// コンテンツの各行にボーダーを追加
	var result strings.Builder
	result.WriteString(topBorder + "\n")

	// タイトル行の下に空行を追加
	emptyLine := strings.Repeat(" ", m.Width-2)
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
	bottomBorder := borderStyleColor.Render("╰" + strings.Repeat("─", m.Width-2) + "╯")
	result.WriteString(bottomBorder)

	return result.String()
}

// テーブル一覧画面のView
func main() {
	p := tea.NewProgram(
		model{Model: app.InitialModel()},
		tea.WithAltScreen(),       // 全画面モード
		tea.WithMouseCellMotion(), // マウスサポート（オプション）
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
