package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/lornest/whoop-cli/internal/tui/style"
)

// StatusBar renders the bottom status bar.
func StatusBar(message string, width int) string {
	s := lipgloss.NewStyle().
		Foreground(style.ColorDim).
		Background(style.ColorBg).
		Width(width).
		Padding(0, 1)

	return s.Render(message)
}

// ErrorBar renders an error status bar.
func ErrorBar(message string, width int) string {
	s := lipgloss.NewStyle().
		Foreground(style.ColorRed).
		Background(style.ColorBg).
		Width(width).
		Padding(0, 1)

	return s.Render("Error: " + message)
}
