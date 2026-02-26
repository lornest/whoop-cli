package style

import "github.com/charmbracelet/lipgloss"

// Whoop-inspired color palette.
var (
	ColorGreen  = lipgloss.Color("#00D68F") // High recovery
	ColorYellow = lipgloss.Color("#FFBE0B") // Medium recovery
	ColorRed    = lipgloss.Color("#FF006E") // Low recovery
	ColorBlue   = lipgloss.Color("#3A86FF") // Strain / accent
	ColorPurple = lipgloss.Color("#8338EC") // REM sleep
	ColorCyan   = lipgloss.Color("#00B4D8") // Deep sleep
	ColorWhite  = lipgloss.Color("#E0E0E0")
	ColorDim    = lipgloss.Color("#888888")

	ColorBg     = lipgloss.Color("#1A1A2E")
	ColorCardBg = lipgloss.Color("#16213E")
	ColorBorder = lipgloss.Color("#2D3561")
)

// RecoveryColor returns the appropriate color for a recovery percentage.
func RecoveryColor(pct float64) lipgloss.Color {
	switch {
	case pct >= 67:
		return ColorGreen
	case pct >= 34:
		return ColorYellow
	default:
		return ColorRed
	}
}

// Common styles.
var (
	CardStyle = lipgloss.NewStyle().
			Background(ColorCardBg).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(1, 2)

	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorWhite).
			Bold(true)

	LabelStyle = lipgloss.NewStyle().
			Foreground(ColorDim)

	ValueStyle = lipgloss.NewStyle().
			Foreground(ColorWhite).
			Bold(true)

	ActiveTabStyle = lipgloss.NewStyle().
			Foreground(ColorBlue).
			Bold(true).
			Underline(true)

	InactiveTabStyle = lipgloss.NewStyle().
				Foreground(ColorDim)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(ColorDim).
			Background(ColorBg)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorRed).
			Bold(true)
)
