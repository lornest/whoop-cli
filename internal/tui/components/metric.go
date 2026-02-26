package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/lornest/whoop-cli/internal/tui/style"
)

// MetricCard renders a styled metric card with value, label, and color.
func MetricCard(label, value string, color lipgloss.Color, width int) string {
	valueStyle := lipgloss.NewStyle().
		Foreground(color).
		Background(style.ColorCardBg).
		Bold(true).
		Width(width - 6). // account for padding
		Align(lipgloss.Center)

	labelStyle := lipgloss.NewStyle().
		Foreground(style.ColorDim).
		Background(style.ColorCardBg).
		Width(width - 6).
		Align(lipgloss.Center)

	content := lipgloss.JoinVertical(lipgloss.Center,
		valueStyle.Render(value),
		labelStyle.Render(label),
	)

	return style.CardStyle.Width(width).Render(content)
}
