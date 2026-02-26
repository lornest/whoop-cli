package components

import (
	"testing"

	"github.com/lornest/whoop-cli/internal/tui/style"
	"github.com/stretchr/testify/assert"
)

func TestMetricCard_ContainsValueAndLabel(t *testing.T) {
	result := MetricCard("Recovery", "78%", style.ColorGreen, 30)
	assert.Contains(t, result, "78%")
	assert.Contains(t, result, "Recovery")
}

func TestMetricCard_DifferentColors(t *testing.T) {
	// Just verify it doesn't panic with different colors
	MetricCard("Strain", "12.5", style.ColorBlue, 30)
	MetricCard("Sleep", "92%", style.ColorCyan, 30)
	MetricCard("Warning", "45%", style.ColorYellow, 30)
}
