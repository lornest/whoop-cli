package screens

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lornest/whoop-cli/internal/tui/components"
	"github.com/lornest/whoop-cli/internal/tui/style"
	"github.com/lornest/whoop-cli/internal/util"
	"github.com/lornest/whoop-cli/pkg/whoop"
)

// DashboardModel shows today's overview.
type DashboardModel struct {
	client   *whoop.Client
	recovery *whoop.RecoveryResponse
	cycles   *whoop.CycleResponse
	sleep    *whoop.SleepResponse
	workouts *whoop.WorkoutResponse
	err      error
	loading  bool
	width    int
	height   int
}

// NewDashboardModel creates a new dashboard screen.
func NewDashboardModel(client *whoop.Client) DashboardModel {
	return DashboardModel{client: client, loading: true}
}

func (m DashboardModel) Init() tea.Cmd {
	return tea.Batch(
		m.fetchRecovery(),
		m.fetchCycles(),
		m.fetchSleep(),
		m.fetchWorkouts(),
	)
}

func (m DashboardModel) fetchRecovery() tea.Cmd {
	return func() tea.Msg {
		start := time.Now().AddDate(0, 0, -1).Format("2006-01-02T15:04:05.000Z")
		resp, err := m.client.GetRecovery(&whoop.QueryParams{Start: start, Limit: 1})
		return RecoveryMsg{Data: resp, Err: err}
	}
}

func (m DashboardModel) fetchCycles() tea.Cmd {
	return func() tea.Msg {
		start := time.Now().AddDate(0, 0, -1).Format("2006-01-02T15:04:05.000Z")
		resp, err := m.client.GetCycles(&whoop.QueryParams{Start: start, Limit: 1})
		return CyclesMsg{Data: resp, Err: err}
	}
}

func (m DashboardModel) fetchSleep() tea.Cmd {
	return func() tea.Msg {
		start := time.Now().AddDate(0, 0, -1).Format("2006-01-02T15:04:05.000Z")
		resp, err := m.client.GetSleep(&whoop.QueryParams{Start: start, Limit: 1})
		return SleepMsg{Data: resp, Err: err}
	}
}

func (m DashboardModel) fetchWorkouts() tea.Cmd {
	return func() tea.Msg {
		start := time.Now().AddDate(0, 0, -1).Format("2006-01-02T15:04:05.000Z")
		resp, err := m.client.GetWorkouts(&whoop.QueryParams{Start: start, Limit: 1})
		return WorkoutsMsg{Data: resp, Err: err}
	}
}

func (m DashboardModel) Update(msg tea.Msg) (DashboardModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case RecoveryMsg:
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.recovery = msg.Data
		}
		m.checkLoading()
	case CyclesMsg:
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.cycles = msg.Data
		}
		m.checkLoading()
	case SleepMsg:
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.sleep = msg.Data
		}
		m.checkLoading()
	case WorkoutsMsg:
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.workouts = msg.Data
		}
		m.checkLoading()
	}
	return m, nil
}

func (m *DashboardModel) checkLoading() {
	if m.recovery != nil || m.cycles != nil || m.sleep != nil || m.workouts != nil || m.err != nil {
		m.loading = false
	}
}

func (m DashboardModel) View() string {
	if m.loading {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Foreground(style.ColorDim).Render("Loading dashboard..."))
	}

	if m.err != nil {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			style.ErrorStyle.Render("Error: "+m.err.Error()+"\n\nPress 'r' to retry"))
	}

	cardWidth := (m.width - 6) / 2
	if cardWidth < 20 {
		cardWidth = 20
	}

	// Recovery card
	recoveryVal := "--"
	recoveryColor := style.ColorDim
	if m.recovery != nil && len(m.recovery.Records) > 0 && m.recovery.Records[0].Score != nil {
		pct := m.recovery.Records[0].Score.RecoveryScore
		recoveryVal = util.FormatPercent(pct)
		recoveryColor = style.RecoveryColor(pct)
	}
	recoveryCard := components.MetricCard("Recovery", recoveryVal, recoveryColor, cardWidth)

	// Strain card
	strainVal := "--"
	if m.cycles != nil && len(m.cycles.Records) > 0 && m.cycles.Records[0].Score != nil {
		strainVal = util.FormatStrain(m.cycles.Records[0].Score.Strain)
	}
	strainCard := components.MetricCard("Day Strain", strainVal, style.ColorBlue, cardWidth)

	// Sleep card
	sleepVal := "--"
	if m.sleep != nil && len(m.sleep.Records) > 0 && m.sleep.Records[0].Score != nil {
		sleepVal = util.FormatPercent(m.sleep.Records[0].Score.SleepPerformancePercentage)
	}
	sleepCard := components.MetricCard("Sleep", sleepVal, style.ColorCyan, cardWidth)

	// Workout card
	workoutVal := "--"
	if m.workouts != nil && len(m.workouts.Records) > 0 && m.workouts.Records[0].Score != nil {
		w := m.workouts.Records[0]
		workoutVal = fmt.Sprintf("%.1f strain", w.Score.Strain)
	}
	workoutCard := components.MetricCard("Last Workout", workoutVal, style.ColorPurple, cardWidth)

	// 2x2 grid
	row1 := lipgloss.JoinHorizontal(lipgloss.Top, recoveryCard, "  ", strainCard)
	row2 := lipgloss.JoinHorizontal(lipgloss.Top, sleepCard, "  ", workoutCard)
	grid := lipgloss.JoinVertical(lipgloss.Left, row1, "", row2)

	title := style.TitleStyle.Render("Today's Overview")
	content := lipgloss.JoinVertical(lipgloss.Left, title, "", grid)

	return lipgloss.Place(m.width, m.height-4, lipgloss.Center, lipgloss.Center, content)
}
