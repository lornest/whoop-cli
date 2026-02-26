package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestAppModel_UnauthenticatedShowsLogin(t *testing.T) {
	m := NewAppModel(nil, false)
	m.width = 80
	m.height = 24
	// Forward window size to login
	m.login.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	view := m.View()
	assert.Contains(t, view, "WHOOP CLI")
}

func TestAppModel_AuthenticatedShowsDashboard(t *testing.T) {
	m := NewAppModel(nil, true)
	m.width = 80
	m.height = 24
	view := m.View()
	assert.Contains(t, view, "Dashboard")
}

func TestAppModel_TabNavigation(t *testing.T) {
	m := NewAppModel(nil, true)
	m.width = 80
	m.height = 24

	// Start on dashboard (tab 0)
	assert.Equal(t, TabDashboard, m.activeTab)

	// Press Tab to go to Recovery
	model, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = model.(AppModel)
	assert.Equal(t, TabRecovery, m.activeTab)

	// Press 4 to jump to Workouts
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'4'}})
	m = model.(AppModel)
	assert.Equal(t, TabWorkouts, m.activeTab)

	// Press Shift+Tab to go back
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	m = model.(AppModel)
	assert.Equal(t, TabSleep, m.activeTab)
}

func TestAppModel_QuitKey(t *testing.T) {
	m := NewAppModel(nil, true)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	// tea.Quit returns a special command
	assert.NotNil(t, cmd)
}

func TestAppModel_WindowResize(t *testing.T) {
	m := NewAppModel(nil, true)
	model, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = model.(AppModel)
	assert.Equal(t, 120, m.width)
	assert.Equal(t, 40, m.height)
}
