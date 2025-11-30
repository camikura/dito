package ui

import "github.com/charmbracelet/lipgloss"

// Color hex values as strings (for use in new_ui and other packages that need string colors)
const (
	ColorPrimaryHex   = "#00D9FF" // Cyan for active borders
	ColorInactiveHex  = "#AAAAAA" // Light gray for inactive borders
	ColorGreenHex     = "#00FF00" // Green for connection status
	ColorLabelHex     = "#00D9FF" // Cyan for section labels
	ColorSecondaryHex = "#C47D7D" // Muted reddish for schema section labels
	ColorTertiaryHex  = "#7AA2F7" // Soft blue for data types
	ColorPKHex        = "#7FBA7A" // Muted green for primary key marker
	ColorIndexHex     = "#E5C07B" // Warm yellow/beige for index field names
	ColorHelpHex      = "#888888" // Gray for help text
	ColorErrorHex     = "#FF0000" // Red for error messages
	ColorSuccessHex   = "#00FF00" // Green for success messages
)

// Color palette as lipgloss.Color (for use in styles)
var (
	ColorPrimary     = lipgloss.Color(ColorPrimaryHex)   // Cyan - used for selected items, borders
	ColorPrimaryBg   = lipgloss.Color("#004466")         // Dark Cyan - background for selected items
	ColorWhite       = lipgloss.Color("#FFFFFF")         // White - normal text
	ColorBlack       = lipgloss.Color("#000000")         // Black - text on light backgrounds
	ColorGray        = lipgloss.Color(ColorHelpHex)      // Gray - labels, footer
	ColorGrayMid     = lipgloss.Color("#666666")         // Mid Gray - grayed out items
	ColorGrayDark    = lipgloss.Color("#555555")         // Dark Gray - separators
	ColorGrayLight   = lipgloss.Color("#CCCCCC")         // Light Gray - SQL display
	ColorGrayLightBg = lipgloss.Color("#333333")         // Light Gray Bg - background for unfocused selected items
	ColorHeaderBg    = lipgloss.Color(ColorInactiveHex)  // Medium Gray - table header background
	ColorHeaderText  = lipgloss.Color("#00AA00")         // Dark Green - table header text
	ColorSuccess     = lipgloss.Color(ColorSuccessHex)   // Green - success messages
	ColorError       = lipgloss.Color(ColorErrorHex)     // Red - error messages
	ColorInactive    = lipgloss.Color(ColorInactiveHex)  // Light gray for inactive borders
	ColorSecondary   = lipgloss.Color(ColorSecondaryHex) // Muted reddish for schema labels
	ColorTertiary    = lipgloss.Color(ColorTertiaryHex)  // Soft blue for data types
	ColorPK          = lipgloss.Color(ColorPKHex)        // Muted green for primary key marker
	ColorIndex       = lipgloss.Color(ColorIndexHex)     // Warm yellow/beige for index field names
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

// Text input cursor styles (unified across the app)
// Uses reverse video (white background, black text) for visibility
var (
	// CursorNarrow is for single-width characters (ASCII, half-width)
	CursorNarrow = lipgloss.NewStyle().Reverse(true)
	// CursorWide is for double-width characters (CJK, full-width)
	CursorWide = lipgloss.NewStyle().Reverse(true)
)
