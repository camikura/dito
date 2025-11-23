package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/camikura/dito/internal/app"
	"github.com/camikura/dito/internal/db"
	"github.com/camikura/dito/internal/handlers"
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

	vm := views.ScreenViewModel{
		Width:  m.Width,
		Height: m.Height,
		Model:  m.Model,
	}

	switch m.Screen {
	case app.ScreenSelection:
		return views.RenderSelectionScreen(vm)
	case app.ScreenOnPremiseConfig:
		return views.RenderOnPremiseConfigScreen(vm)
	case app.ScreenCloudConfig:
		return views.RenderCloudConfigScreen(vm)
	case app.ScreenTableList:
		return views.RenderTableListView(m.ToTableListViewModel())
	default:
		return "Unknown screen"
	}
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
