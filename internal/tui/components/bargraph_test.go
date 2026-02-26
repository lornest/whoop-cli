package components

import (
	"testing"

	"github.com/lornest/whoop-cli/internal/tui/style"
	"github.com/stretchr/testify/assert"
)

func TestHorizontalBarChart_Empty(t *testing.T) {
	result := HorizontalBarChart(nil, 60)
	assert.Equal(t, "", result)
}

func TestHorizontalBarChart_SingleItem(t *testing.T) {
	items := []BarItem{
		{Label: "Mon", Value: 78, Color: style.ColorGreen},
	}
	result := HorizontalBarChart(items, 60)
	assert.Contains(t, result, "Mon")
	assert.Contains(t, result, "78")
	assert.Contains(t, result, "█")
}

func TestHorizontalBarChart_MultipleItems(t *testing.T) {
	items := []BarItem{
		{Label: "Mon", Value: 78, Color: style.ColorGreen},
		{Label: "Tue", Value: 45, Color: style.ColorYellow},
		{Label: "Wed", Value: 25, Color: style.ColorRed},
	}
	result := HorizontalBarChart(items, 60)
	assert.Contains(t, result, "Mon")
	assert.Contains(t, result, "Tue")
	assert.Contains(t, result, "Wed")
}

func TestHorizontalBarChart_ZeroValues(t *testing.T) {
	items := []BarItem{
		{Label: "Mon", Value: 0, Color: style.ColorGreen},
	}
	result := HorizontalBarChart(items, 60)
	assert.Contains(t, result, "Mon")
}
