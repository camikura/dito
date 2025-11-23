package handlers

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/camikura/dito/internal/app"
)

func TestHandleSelection(t *testing.T) {
	tests := []struct {
		name            string
		initialModel    app.Model
		key             string
		expectedScreen  app.Screen
		expectedCursor  int
		expectQuitCmd   bool
	}{
		{
			name: "quit with ctrl+c",
			initialModel: app.Model{
				Screen:  app.ScreenSelection,
				Choices: []string{"Cloud", "On-Premise"},
				Cursor:  0,
			},
			key:           "ctrl+c",
			expectedScreen: app.ScreenSelection,
			expectedCursor: 0,
			expectQuitCmd:  true,
		},
		{
			name: "cursor down with down key",
			initialModel: app.Model{
				Screen:  app.ScreenSelection,
				Choices: []string{"Cloud", "On-Premise"},
				Cursor:  0,
			},
			key:           "down",
			expectedScreen: app.ScreenSelection,
			expectedCursor: 1,
			expectQuitCmd:  false,
		},
		{
			name: "cursor down with tab key",
			initialModel: app.Model{
				Screen:  app.ScreenSelection,
				Choices: []string{"Cloud", "On-Premise"},
				Cursor:  0,
			},
			key:           "tab",
			expectedScreen: app.ScreenSelection,
			expectedCursor: 1,
			expectQuitCmd:  false,
		},
		{
			name: "cursor wraps to beginning when at end",
			initialModel: app.Model{
				Screen:  app.ScreenSelection,
				Choices: []string{"Cloud", "On-Premise"},
				Cursor:  1,
			},
			key:           "down",
			expectedScreen: app.ScreenSelection,
			expectedCursor: 0,
			expectQuitCmd:  false,
		},
		{
			name: "cursor up with up key",
			initialModel: app.Model{
				Screen:  app.ScreenSelection,
				Choices: []string{"Cloud", "On-Premise"},
				Cursor:  1,
			},
			key:           "up",
			expectedScreen: app.ScreenSelection,
			expectedCursor: 0,
			expectQuitCmd:  false,
		},
		{
			name: "cursor up with shift+tab key",
			initialModel: app.Model{
				Screen:  app.ScreenSelection,
				Choices: []string{"Cloud", "On-Premise"},
				Cursor:  1,
			},
			key:           "shift+tab",
			expectedScreen: app.ScreenSelection,
			expectedCursor: 0,
			expectQuitCmd:  false,
		},
		{
			name: "cursor wraps to end when at beginning",
			initialModel: app.Model{
				Screen:  app.ScreenSelection,
				Choices: []string{"Cloud", "On-Premise"},
				Cursor:  0,
			},
			key:           "up",
			expectedScreen: app.ScreenSelection,
			expectedCursor: 1,
			expectQuitCmd:  false,
		},
		{
			name: "enter on Cloud selection",
			initialModel: app.Model{
				Screen:  app.ScreenSelection,
				Choices: []string{"Cloud", "On-Premise"},
				Cursor:  0,
			},
			key:           "enter",
			expectedScreen: app.ScreenCloudConfig,
			expectedCursor: 0,
			expectQuitCmd:  false,
		},
		{
			name: "enter on On-Premise selection",
			initialModel: app.Model{
				Screen:  app.ScreenSelection,
				Choices: []string{"Cloud", "On-Premise"},
				Cursor:  1,
			},
			key:           "enter",
			expectedScreen: app.ScreenOnPremiseConfig,
			expectedCursor: 1,
			expectQuitCmd:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tea.KeyMsg{Type: tea.KeyRunes}
			switch tt.key {
			case "ctrl+c":
				msg.Type = tea.KeyCtrlC
			case "up":
				msg.Type = tea.KeyUp
			case "down":
				msg.Type = tea.KeyDown
			case "tab":
				msg.Type = tea.KeyTab
			case "shift+tab":
				msg.Type = tea.KeyShiftTab
			case "enter":
				msg.Type = tea.KeyEnter
			}

			resultModel, resultCmd := HandleSelection(tt.initialModel, msg)

			// Check screen
			if resultModel.Screen != tt.expectedScreen {
				t.Errorf("HandleSelection() Screen = %v, want %v", resultModel.Screen, tt.expectedScreen)
			}

			// Check cursor
			if resultModel.Cursor != tt.expectedCursor {
				t.Errorf("HandleSelection() Cursor = %v, want %v", resultModel.Cursor, tt.expectedCursor)
			}

			// Check quit command
			if tt.expectQuitCmd {
				if resultCmd == nil {
					t.Error("HandleSelection() should return tea.Quit command")
				}
			} else {
				if tt.key != "enter" && resultCmd != nil {
					t.Error("HandleSelection() should not return a command for non-enter keys")
				}
			}
		})
	}
}
