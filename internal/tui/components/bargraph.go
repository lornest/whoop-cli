package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lornest/whoop-cli/internal/tui/style"
)

// BarItem represents a single bar in the chart.
type BarItem struct {
	Label string
	Value float64
	Color lipgloss.Color
}

// HorizontalBarChart renders a horizontal bar chart.
func HorizontalBarChart(items []BarItem, maxWidth int) string {
	if len(items) == 0 {
		return ""
	}

	// Find max value for scaling
	maxVal := 0.0
	for _, item := range items {
		if item.Value > maxVal {
			maxVal = item.Value
		}
	}
	if maxVal == 0 {
		maxVal = 1
	}

	barWidth := maxWidth - 20 // leave room for label + value
	if barWidth < 10 {
		barWidth = 10
	}

	var lines []string
	for _, item := range items {
		width := int(float64(barWidth) * (item.Value / maxVal))
		if width < 1 && item.Value > 0 {
			width = 1
		}

		bar := lipgloss.NewStyle().
			Foreground(item.Color).
			Render(strings.Repeat("█", width))

		label := lipgloss.NewStyle().
			Foreground(style.ColorDim).
			Width(8).
			Align(lipgloss.Right).
			Render(item.Label)

		value := lipgloss.NewStyle().
			Foreground(item.Color).
			Width(8).
			Render(fmt.Sprintf(" %.0f", item.Value))

		lines = append(lines, label+" "+bar+value)
	}

	return strings.Join(lines, "\n")
}
