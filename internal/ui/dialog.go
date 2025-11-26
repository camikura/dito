package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// DialogType represents the type of dialog.
type DialogType int

const (
	// DialogTypeInfo represents an informational dialog.
	DialogTypeInfo DialogType = iota
	// DialogTypeSuccess represents a success dialog.
	DialogTypeSuccess
	// DialogTypeError represents an error dialog.
	DialogTypeError
	// DialogTypeInput represents an input dialog.
	DialogTypeInput
)

// DialogConfig holds configuration for a dialog.
type DialogConfig struct {
	Title       string
	Content     string
	HelpText    string
	Width       int    // 0 = auto
	Height      int    // 0 = auto
	BorderColor string // Default: ColorPrimary
	Type        DialogType
}

// Dialog represents a modal dialog box.
type Dialog struct {
	config DialogConfig
}

// NewDialog creates a new Dialog with the given configuration.
func NewDialog(config DialogConfig) *Dialog {
	if config.BorderColor == "" {
		switch config.Type {
		case DialogTypeSuccess:
			config.BorderColor = "#00FF00"
		case DialogTypeError:
			config.BorderColor = "#FF0000"
		default:
			config.BorderColor = string(ColorPrimary)
		}
	}
	return &Dialog{config: config}
}

// Render renders the dialog as a string.
func (d *Dialog) Render() string {
	width := d.config.Width
	if width <= 0 {
		width = 60
	}

	// Build styles
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(d.config.BorderColor)).
		Padding(1, 2).
		Width(width)

	if d.config.Height > 0 {
		borderStyle = borderStyle.Height(d.config.Height)
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(d.config.BorderColor)).
		Bold(true).
		Width(width - 6)

	contentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Width(width - 6)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Width(width - 6)

	// Build content
	var parts []string

	if d.config.Title != "" {
		parts = append(parts, titleStyle.Render(d.config.Title))
		parts = append(parts, "")
	}

	if d.config.Content != "" {
		// Wrap text if needed
		wrapped := WrapText(d.config.Content, width-8)
		parts = append(parts, contentStyle.Render(wrapped))
	}

	if d.config.HelpText != "" {
		parts = append(parts, "")
		parts = append(parts, helpStyle.Render(d.config.HelpText))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, parts...)
	return borderStyle.Render(content)
}

// RenderCentered renders the dialog centered on screen.
func (d *Dialog) RenderCentered(screenWidth, screenHeight int) string {
	return lipgloss.Place(
		screenWidth,
		screenHeight,
		lipgloss.Center,
		lipgloss.Center,
		d.Render(),
	)
}

// WrapText wraps text to fit within the specified width.
func WrapText(text string, width int) string {
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

// CalculateDialogSize calculates appropriate dialog size based on screen size.
func CalculateDialogSize(screenWidth, screenHeight int, minWidth, maxWidth, minHeight, maxHeight int) (width, height int) {
	// Width calculation
	width = screenWidth - 10
	if minWidth > 0 && width < minWidth {
		width = minWidth
	}
	if maxWidth > 0 && width > maxWidth {
		width = maxWidth
	}

	// Height calculation
	height = screenHeight - 10
	if minHeight > 0 && height < minHeight {
		height = minHeight
	}
	if maxHeight > 0 && height > maxHeight {
		height = maxHeight
	}

	return width, height
}
