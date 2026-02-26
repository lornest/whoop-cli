package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lornest/whoop-cli/internal/tui/style"
)

// TabNames defines the available tabs.
var TabNames = []string{"Dashboard", "Recovery", "Sleep", "Workouts", "Profile"}

// RenderTabs renders the tab bar with the given active index.
func RenderTabs(active int, width int) string {
	var tabs []string
	for i, name := range TabNames {
		label := name
		if i == active {
			tabs = append(tabs, style.ActiveTabStyle.Render(label))
		} else {
			tabs = append(tabs, style.InactiveTabStyle.Render(label))
		}
	}
	bar := strings.Join(tabs, "  │  ")
	return lipgloss.NewStyle().Width(width).Render(bar)
}
