package screens

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lornest/whoop-cli/internal/tui/components"
	"github.com/lornest/whoop-cli/internal/tui/style"
	"github.com/lornest/whoop-cli/pkg/whoop"
)

// RecoveryModel shows recovery trends.
type RecoveryModel struct {
	client   *whoop.Client
	recovery *whoop.RecoveryResponse
	err      error
	loading  bool
	width    int
	height   int
}

func NewRecoveryModel(client *whoop.Client) RecoveryModel {
	return RecoveryModel{client: client, loading: true}
}

func (m RecoveryModel) Init() tea.Cmd {
	return m.fetchRecovery()
}

func (m RecoveryModel) fetchRecovery() tea.Cmd {
	return func() tea.Msg {
		start := time.Now().AddDate(0, 0, -7).Format("2006-01-02T15:04:05.000Z")
		resp, err := m.client.GetRecovery(&whoop.QueryParams{Start: start, Limit: 7})
		return RecoveryMsg{Data: resp, Err: err}
	}
}

func (m RecoveryModel) Update(msg tea.Msg) (RecoveryModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case RecoveryMsg:
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.recovery = msg.Data
		}
	}
	return m, nil
}

func (m RecoveryModel) View() string {
	if m.loading {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Foreground(style.ColorDim).Render("Loading recovery data..."))
	}

	if m.err != nil {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			style.ErrorStyle.Render("Error: "+m.err.Error()+"\n\nPress 'r' to retry"))
	}

	if m.recovery == nil || len(m.recovery.Records) == 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Foreground(style.ColorDim).Render("No recovery data available"))
	}

	title := style.TitleStyle.Render("Recovery — 7 Day Trend")

	// Build bar chart items (reverse order so oldest is first)
	var barItems []components.BarItem
	var hrvValues, rhrValues []float64

	for i := len(m.recovery.Records) - 1; i >= 0; i-- {
		r := m.recovery.Records[i]
		if r.Score == nil {
			continue
		}
		label := fmt.Sprintf("Day %d", len(m.recovery.Records)-i)
		pct := r.Score.RecoveryScore
		barItems = append(barItems, components.BarItem{
			Label: label,
			Value: pct,
			Color: style.RecoveryColor(pct),
		})
		hrvValues = append(hrvValues, r.Score.HRVRmssdMilli)
		rhrValues = append(rhrValues, r.Score.RestingHeartRate)
	}

	chart := components.HorizontalBarChart(barItems, m.width-4)

	// Latest metrics
	latest := m.recovery.Records[0]
	var metricsStr string
	if latest.Score != nil {
		metricsStr = lipgloss.JoinHorizontal(lipgloss.Top,
			components.MetricCard("Recovery", fmt.Sprintf("%.0f%%", latest.Score.RecoveryScore),
				style.RecoveryColor(latest.Score.RecoveryScore), 20),
			"  ",
			components.MetricCard("RHR", fmt.Sprintf("%.0f bpm", latest.Score.RestingHeartRate),
				style.ColorBlue, 20),
			"  ",
			components.MetricCard("HRV", fmt.Sprintf("%.1f ms", latest.Score.HRVRmssdMilli),
				style.ColorGreen, 20),
		)
	}

	// Sparklines
	sparkSection := ""
	if len(hrvValues) > 1 {
		hrvSpark := components.Sparkline(hrvValues, style.ColorGreen)
		rhrSpark := components.Sparkline(rhrValues, style.ColorBlue)
		sparkSection = lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.NewStyle().Foreground(style.ColorDim).Render("HRV Trend: ")+hrvSpark,
			lipgloss.NewStyle().Foreground(style.ColorDim).Render("RHR Trend: ")+rhrSpark,
		)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, title, "", chart, "", metricsStr, "", sparkSection)
	return lipgloss.Place(m.width, m.height-4, lipgloss.Center, lipgloss.Center, content)
}
