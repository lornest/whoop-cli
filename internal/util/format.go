package util

import "fmt"

const (
	metersToFeet   = 3.28084
	kgToLbs        = 2.20462
	celsiusToFBase = 1.8
	celsiusToFAdd  = 32.0
)

// MillisToDuration converts milliseconds to a human-readable "Xh Ym" string.
func MillisToDuration(ms int) string {
	if ms <= 0 {
		return "0m"
	}
	totalMinutes := ms / 60000
	hours := totalMinutes / 60
	minutes := totalMinutes % 60
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

// MetersToFeet converts meters to feet.
func MetersToFeet(m float64) float64 {
	return m * metersToFeet
}

// MetersToFeetInches returns a formatted "X'Y\"" string.
func MetersToFeetInches(m float64) string {
	totalInches := m * metersToFeet * 12
	feet := int(totalInches) / 12
	inches := int(totalInches) % 12
	return fmt.Sprintf("%d'%d\"", feet, inches)
}

// KgToLbs converts kilograms to pounds.
func KgToLbs(kg float64) float64 {
	return kg * kgToLbs
}

// CelsiusToFahrenheit converts Celsius to Fahrenheit.
func CelsiusToFahrenheit(c float64) float64 {
	return c*celsiusToFBase + celsiusToFAdd
}

// MetersToMiles converts meters to miles.
func MetersToMiles(m float64) float64 {
	return m / 1609.344
}

// KilojoulesToCalories converts kilojoules to kilocalories.
func KilojoulesToCalories(kj float64) float64 {
	return kj / 4.184
}

// FormatStrain formats a strain value to one decimal place.
func FormatStrain(s float64) string {
	return fmt.Sprintf("%.1f", s)
}

// FormatPercent formats a percentage value with the % symbol.
func FormatPercent(p float64) string {
	return fmt.Sprintf("%.0f%%", p)
}
