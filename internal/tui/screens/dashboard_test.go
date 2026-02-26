package screens

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lornest/whoop-cli/pkg/whoop"
	"github.com/stretchr/testify/assert"
)

func TestDashboardModel_Loading(t *testing.T) {
	m := DashboardModel{loading: true, width: 80, height: 24}
	view := m.View()
	assert.Contains(t, view, "Loading")
}

func TestDashboardModel_WithData(t *testing.T) {
	score := 78.0
	m := DashboardModel{
		width:  80,
		height: 24,
		recovery: &whoop.RecoveryResponse{
			Records: []whoop.Recovery{
				{ScoreState: "SCORED", Score: &whoop.RecoveryScore{RecoveryScore: score}},
			},
		},
		cycles: &whoop.CycleResponse{
			Records: []whoop.Cycle{
				{ScoreState: "SCORED", Score: &whoop.CycleScore{Strain: 12.5}},
			},
		},
		sleep: &whoop.SleepResponse{
			Records: []whoop.Sleep{
				{ScoreState: "SCORED", Score: &whoop.SleepScore{SleepPerformancePercentage: 92}},
			},
		},
		workouts: &whoop.WorkoutResponse{
			Records: []whoop.Workout{
				{ScoreState: "SCORED", Score: &whoop.WorkoutScore{Strain: 14.2}},
			},
		},
	}

	view := m.View()
	assert.Contains(t, view, "78%")
	assert.Contains(t, view, "12.5")
	assert.Contains(t, view, "92%")
	assert.Contains(t, view, "14.2")
	assert.Contains(t, view, "Recovery")
	assert.Contains(t, view, "Day Strain")
}

func TestDashboardModel_Error(t *testing.T) {
	m := DashboardModel{
		width:  80,
		height: 24,
		err:    assert.AnError,
	}
	view := m.View()
	assert.Contains(t, view, "Error")
}

func TestDashboardModel_EmptyData(t *testing.T) {
	m := DashboardModel{
		width:    80,
		height:   24,
		recovery: &whoop.RecoveryResponse{Records: []whoop.Recovery{}},
		cycles:   &whoop.CycleResponse{Records: []whoop.Cycle{}},
	}
	view := m.View()
	assert.Contains(t, view, "--")
}

func TestDashboardModel_WindowResize(t *testing.T) {
	m := DashboardModel{}
	m, _ = m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	assert.Equal(t, 120, m.width)
	assert.Equal(t, 40, m.height)
}
