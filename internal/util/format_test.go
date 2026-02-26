package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMillisToDuration(t *testing.T) {
	tests := []struct {
		name string
		ms   int
		want string
	}{
		{"zero", 0, "0m"},
		{"negative", -100, "0m"},
		{"30 minutes", 1800000, "30m"},
		{"1 hour", 3600000, "1h 0m"},
		{"8h 15m", 29700000, "8h 15m"},
		{"2h 30m", 9000000, "2h 30m"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, MillisToDuration(tt.ms))
		})
	}
}

func TestMetersToFeet(t *testing.T) {
	assert.InDelta(t, 6.0, MetersToFeet(1.8288), 0.01)
}

func TestMetersToFeetInches(t *testing.T) {
	assert.Equal(t, "6'0\"", MetersToFeetInches(1.8288))
	assert.Equal(t, "5'10\"", MetersToFeetInches(1.778))
}

func TestKgToLbs(t *testing.T) {
	assert.InDelta(t, 180.0, KgToLbs(81.6466), 0.2)
}

func TestCelsiusToFahrenheit(t *testing.T) {
	assert.InDelta(t, 91.76, CelsiusToFahrenheit(33.2), 0.01)
	assert.InDelta(t, 32.0, CelsiusToFahrenheit(0), 0.01)
}

func TestMetersToMiles(t *testing.T) {
	assert.InDelta(t, 5.0, MetersToMiles(8046.72), 0.01)
}

func TestKilojoulesToCalories(t *testing.T) {
	assert.InDelta(t, 442.16, KilojoulesToCalories(1850.5), 0.5)
}

func TestFormatStrain(t *testing.T) {
	assert.Equal(t, "12.5", FormatStrain(12.5))
	assert.Equal(t, "0.0", FormatStrain(0))
}

func TestFormatPercent(t *testing.T) {
	assert.Equal(t, "78%", FormatPercent(78.0))
	assert.Equal(t, "92%", FormatPercent(92.0))
}
