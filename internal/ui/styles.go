package ui

import "github.com/charmbracelet/lipgloss"

// Color palette
var (
	ColorPrimary   = lipgloss.Color("#00D9FF") // Cyan - used for selected items, borders
	ColorWhite     = lipgloss.Color("#FFFFFF") // White - normal text
	ColorGray      = lipgloss.Color("#888888") // Gray - labels, footer
	ColorGrayMid   = lipgloss.Color("#666666") // Mid Gray - grayed out items
	ColorGrayDark  = lipgloss.Color("#555555") // Dark Gray - separators
	ColorGrayLight = lipgloss.Color("#CCCCCC") // Light Gray - SQL display
	ColorSuccess   = lipgloss.Color("#00FF00") // Green - success messages
	ColorError     = lipgloss.Color("#FF0000") // Red - error messages
)

// Common text styles
var (
	StyleTitle    = lipgloss.NewStyle().Foreground(ColorWhite).Bold(true)
	StyleNormal   = lipgloss.NewStyle().Foreground(ColorWhite)
	StyleFocused  = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
	StyleSelected = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
	StyleLabel    = lipgloss.NewStyle().Foreground(ColorGray)
	StyleSuccess  = lipgloss.NewStyle().Foreground(ColorSuccess)
	StyleError    = lipgloss.NewStyle().Foreground(ColorError)
)

// Layout styles
var (
	StyleSeparator = lipgloss.NewStyle().Foreground(ColorGrayDark)
	StyleBorder    = lipgloss.NewStyle().Foreground(ColorPrimary)
)
