package ui

import "github.com/charmbracelet/lipgloss"

// Color palette
var (
	ColorPrimary     = lipgloss.Color("#00D9FF") // Cyan - used for selected items, borders
	ColorPrimaryBg   = lipgloss.Color("#004466") // Dark Cyan - background for selected items
	ColorWhite       = lipgloss.Color("#FFFFFF") // White - normal text
	ColorBlack       = lipgloss.Color("#000000") // Black - text on light backgrounds
	ColorGray        = lipgloss.Color("#888888") // Gray - labels, footer
	ColorGrayMid     = lipgloss.Color("#666666") // Mid Gray - grayed out items
	ColorGrayDark    = lipgloss.Color("#555555") // Dark Gray - separators
	ColorGrayLight   = lipgloss.Color("#CCCCCC") // Light Gray - SQL display
	ColorGrayLightBg = lipgloss.Color("#333333") // Light Gray Bg - background for unfocused selected items
	ColorHeaderBg    = lipgloss.Color("#AAAAAA") // Medium Gray - table header background
	ColorHeaderText  = lipgloss.Color("#00AA00") // Dark Green - table header text
	ColorSuccess     = lipgloss.Color("#00FF00") // Green - success messages
	ColorError       = lipgloss.Color("#FF0000") // Red - error messages
)

// Common text styles
var (
	StyleTitle    = lipgloss.NewStyle().Foreground(ColorWhite).Bold(true)
	StyleNormal   = lipgloss.NewStyle().Foreground(ColorWhite)
	StyleFocused  = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
	StyleSelected = lipgloss.NewStyle().Foreground(ColorWhite).Background(ColorPrimaryBg) // Background highlight
	StyleHeader   = lipgloss.NewStyle().Foreground(ColorHeaderText).Bold(true).Underline(true) // Table header with underline (dark green)
	StyleLabel    = lipgloss.NewStyle().Foreground(ColorGray)
	StyleDim      = lipgloss.NewStyle().Foreground(ColorGrayMid) // Dim text for null values
	StyleSuccess  = lipgloss.NewStyle().Foreground(ColorSuccess)
	StyleError    = lipgloss.NewStyle().Foreground(ColorError)
)

// Layout styles
var (
	StyleSeparator = lipgloss.NewStyle().Foreground(ColorGrayDark)
	StyleBorder    = lipgloss.NewStyle().Foreground(ColorPrimary)
)
