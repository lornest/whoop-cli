package screens

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lornest/whoop-cli/internal/tui/components"
	"github.com/lornest/whoop-cli/internal/tui/style"
	"github.com/lornest/whoop-cli/internal/util"
	"github.com/lornest/whoop-cli/pkg/whoop"
)

// SleepModel shows sleep breakdown.
type SleepModel struct {
	client  *whoop.Client
	sleep   *whoop.SleepResponse
	table   table.Model
	err     error
	loading bool
	width   int
	height  int
}

func NewSleepModel(client *whoop.Client) SleepModel {
	columns := []table.Column{
		{Title: "Date", Width: 12},
		{Title: "In Bed", Width: 10},
		{Title: "Perf %", Width: 10},
		{Title: "Efficiency", Width: 12},
		{Title: "Nap", Width: 8},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(style.ColorBorder).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(style.ColorCyan).
		Bold(true)
	t.SetStyles(s)

	return SleepModel{client: client, loading: true, table: t}
}

func (m SleepModel) Init() tea.Cmd {
	return m.fetchSleep()
}

func (m SleepModel) fetchSleep() tea.Cmd {
	return func() tea.Msg {
		start := time.Now().AddDate(0, 0, -7).Format("2006-01-02T15:04:05.000Z")
		resp, err := m.client.GetSleep(&whoop.QueryParams{Start: start, Limit: 7})
		return SleepMsg{Data: resp, Err: err}
	}
}

func (m SleepModel) Update(msg tea.Msg) (SleepModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table.SetWidth(m.width - 8)
	case SleepMsg:
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.sleep = msg.Data
			m.populateTable()
		}
	}

	if !m.loading && m.sleep != nil && len(m.sleep.Records) > 0 {
		m.table, cmd = m.table.Update(msg)
	}

	return m, cmd
}

func (m *SleepModel) populateTable() {
	var rows []table.Row
	for _, rec := range m.sleep.Records {
		if rec.Score == nil {
			continue
		}
		t, _ := time.Parse("2006-01-02T15:04:05.000Z", rec.Start)
		date := t.Format("Jan 02")
		inBed := util.MillisToDuration(rec.Score.StageSummary.TotalInBedTimeMilli)
		perf := util.FormatPercent(rec.Score.SleepPerformancePercentage)
		eff := util.FormatPercent(rec.Score.SleepEfficiencyPercentage)
		nap := "No"
		if rec.Nap {
			nap = "Yes"
		}
		rows = append(rows, table.Row{date, inBed, perf, eff, nap})
	}
	m.table.SetRows(rows)
}

func (m SleepModel) View() string {
	if m.loading {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Foreground(style.ColorDim).Render("Loading sleep data..."))
	}

	if m.err != nil {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			style.ErrorStyle.Render("Error: "+m.err.Error()+"\n\nPress 'r' to retry"))
	}

	if m.sleep == nil || len(m.sleep.Records) == 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Foreground(style.ColorDim).Render("No sleep data available"))
	}

	title := style.TitleStyle.Render("Sleep — 7 Day Overview")

	// Selected sleep breakdown
	selectedIdx := m.table.Cursor()
	var stagesView string
	var metricsStr string

	if selectedIdx >= 0 && selectedIdx < len(m.sleep.Records) {
		selected := m.sleep.Records[selectedIdx]
		if selected.Score != nil {
			s := selected.Score.StageSummary
			total := float64(s.TotalInBedTimeMilli)
			if total == 0 {
				total = 1
			}

			stageItems := []components.BarItem{
				{Label: "Awake", Value: float64(s.TotalAwakeTimeMilli) / total * 100, Color: style.ColorRed},
				{Label: "Light", Value: float64(s.TotalLightSleepTimeMilli) / total * 100, Color: style.ColorYellow},
				{Label: "Deep", Value: float64(s.TotalSlowWaveSleepTimeMilli) / total * 100, Color: style.ColorCyan},
				{Label: "REM", Value: float64(s.TotalREMSleepTimeMilli) / total * 100, Color: style.ColorPurple},
			}

			stagesView = lipgloss.JoinVertical(lipgloss.Left,
				style.LabelStyle.Render("Sleep Stages"),
				components.HorizontalBarChart(stageItems, m.width-4),
				"",
				fmt.Sprintf("  In Bed: %s  |  Awake: %s  |  Light: %s  |  Deep: %s  |  REM: %s",
					util.MillisToDuration(s.TotalInBedTimeMilli),
					util.MillisToDuration(s.TotalAwakeTimeMilli),
					util.MillisToDuration(s.TotalLightSleepTimeMilli),
					util.MillisToDuration(s.TotalSlowWaveSleepTimeMilli),
					util.MillisToDuration(s.TotalREMSleepTimeMilli),
				),
			)

			metricsStr = lipgloss.JoinHorizontal(lipgloss.Top,
				components.MetricCard("Performance",
					util.FormatPercent(selected.Score.SleepPerformancePercentage), style.ColorGreen, 22),
				"  ",
				components.MetricCard("Efficiency",
					util.FormatPercent(selected.Score.SleepEfficiencyPercentage), style.ColorCyan, 22),
				"  ",
				components.MetricCard("Resp Rate",
					fmt.Sprintf("%.1f/min", selected.Score.RespiratoryRate), style.ColorBlue, 22),
			)
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		title, "", metricsStr, "", stagesView, "", lipgloss.NewStyle().Render(m.table.View()))
	return lipgloss.Place(m.width, m.height-4, lipgloss.Left, lipgloss.Top, lipgloss.NewStyle().Padding(1, 2).Render(content))
}
