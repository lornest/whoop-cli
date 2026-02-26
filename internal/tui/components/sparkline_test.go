package components

import (
	"testing"

	"github.com/lornest/whoop-cli/internal/tui/style"
	"github.com/stretchr/testify/assert"
)

func TestSparkline_Empty(t *testing.T) {
	result := Sparkline(nil, style.ColorGreen)
	assert.Equal(t, "", result)
}

func TestSparkline_SingleValue(t *testing.T) {
	result := Sparkline([]float64{50}, style.ColorGreen)
	assert.NotEmpty(t, result)
}

func TestSparkline_Ascending(t *testing.T) {
	result := Sparkline([]float64{1, 2, 3, 4, 5, 6, 7, 8}, style.ColorGreen)
	// Should contain ascending block characters
	assert.Contains(t, result, "▁")
	assert.Contains(t, result, "█")
}

func TestSparkline_Constant(t *testing.T) {
	result := Sparkline([]float64{5, 5, 5, 5}, style.ColorBlue)
	assert.NotEmpty(t, result)
}

func TestSparkline_Length(t *testing.T) {
	values := []float64{10, 20, 30, 40, 50}
	result := Sparkline(values, style.ColorCyan)
	// Result contains ANSI codes, but the spark characters should be there
	assert.NotEmpty(t, result)
}
