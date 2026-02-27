package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/lornest/whoop-cli/internal/tui/style"
)

// TabNames defines the available tabs.
var TabNames = []string{"Dashboard", "Recovery", "Sleep", "Workouts", "Profile"}

// RenderTabs renders the tab bar with the given active index.
func RenderTabs(active int, width int) string {
	var tabs []string
	for i, name := range TabNames {
		if i == active {
			tabs = append(tabs, style.ActiveTabStyle.Render(name))
		} else {
			tabs = append(tabs, style.InactiveTabStyle.Render(name))
		}
	}
	row := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)

	gapWidth := width - lipgloss.Width(row)
	if gapWidth < 0 {
		gapWidth = 0
	}

	gap := style.TabGapStyle.
		Width(gapWidth).
		Render("")

	return lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
}
