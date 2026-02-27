package screens

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lornest/whoop-cli/pkg/whoop"
	"github.com/stretchr/testify/assert"
)

func TestWorkoutsModel_Loading(t *testing.T) {
	m := WorkoutsModel{loading: true, width: 80, height: 24}
	view := m.View()
	assert.Contains(t, view, "Loading")
}

func TestWorkoutsModel_WithData(t *testing.T) {
	m := NewWorkoutsModel(nil)
	m.width = 120
	m.height = 40

	msg := WorkoutsMsg{
		Data: &whoop.WorkoutResponse{
			Records: []whoop.Workout{
				{
					Start:      "2026-02-25T17:00:00.000Z",
					End:        "2026-02-25T18:05:00.000Z",
					SportID:    1,
					ScoreState: "SCORED",
					Score:      &whoop.WorkoutScore{Strain: 14.2, AverageHeartRate: 145, MaxHeartRate: 178, Kilojoule: 1850.5, DistanceMeter: 8045},
				},
			},
		},
	}
	m, _ = m.Update(msg)

	view := m.View()
	assert.Contains(t, view, "Workouts")
	assert.Contains(t, view, "Running")
	assert.Contains(t, view, "14.2")
}

func TestWorkoutsModel_Navigation(t *testing.T) {
	m := NewWorkoutsModel(nil)
	m.width = 120
	m.height = 40

	msg := WorkoutsMsg{
		Data: &whoop.WorkoutResponse{
			Records: []whoop.Workout{
				{Start: "2026-02-25T17:00:00.000Z", SportID: 1, Score: &whoop.WorkoutScore{Strain: 14.2}},
				{Start: "2026-02-24T07:00:00.000Z", SportID: 71, Score: &whoop.WorkoutScore{Strain: 8.7}},
			},
		},
	}

	m, _ = m.Update(msg)

	assert.Equal(t, 0, m.table.Cursor())

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	assert.Equal(t, 1, m.table.Cursor())

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	assert.Equal(t, 0, m.table.Cursor())

	// Can't go below 0
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	assert.Equal(t, 0, m.table.Cursor())
}

func TestWorkoutsModel_Empty(t *testing.T) {
	m := NewWorkoutsModel(nil)
	m.width = 80
	m.height = 24
	m, _ = m.Update(WorkoutsMsg{Data: &whoop.WorkoutResponse{Records: []whoop.Workout{}}})
	view := m.View()
	assert.Contains(t, view, "No workouts")
}

func TestGetSportName(t *testing.T) {
	assert.Equal(t, "Running", getSportName(1))
	assert.Equal(t, "Functional Fitness", getSportName(71))
	assert.Contains(t, getSportName(999), "Sport 999")
}
