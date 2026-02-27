package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lornest/whoop-cli/internal/tui/components"
	"github.com/lornest/whoop-cli/internal/tui/screens"
	"github.com/lornest/whoop-cli/pkg/whoop"
)

const (
	TabDashboard = iota
	TabRecovery
	TabSleep
	TabWorkouts
	TabProfile
)

// AppModel is the root TUI model.
type AppModel struct {
	client    *whoop.Client
	keys      KeyMap
	activeTab int
	width     int
	height    int

	// Screens
	login     screens.LoginModel
	dashboard screens.DashboardModel
	recovery  screens.RecoveryModel
	sleep     screens.SleepModel
	workouts  screens.WorkoutsModel
	profile   screens.ProfileModel

	authenticated bool
	initialized   [5]bool // track which screens have been initialized
	statusMsg     string
}

// NewAppModel creates the root model.
func NewAppModel(client *whoop.Client, authenticated bool) AppModel {
	m := AppModel{
		client:        client,
		keys:          DefaultKeyMap(),
		authenticated: authenticated,
		login:         screens.NewLoginModel(),
		dashboard:     screens.NewDashboardModel(client),
		recovery:      screens.NewRecoveryModel(client),
		sleep:         screens.NewSleepModel(client),
		workouts:      screens.NewWorkoutsModel(client),
		profile:       screens.NewProfileModel(client),
		statusMsg:     "Press ? for help  │  Tab to switch  │  q to quit",
	}
	return m
}

func (m AppModel) Init() tea.Cmd {
	if !m.authenticated {
		return m.login.Init()
	}
	m.initialized[TabDashboard] = true
	return m.dashboard.Init()
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Forward to all screens
		m.login, _ = m.login.Update(msg)
		m.dashboard, _ = m.dashboard.Update(msg)
		m.recovery, _ = m.recovery.Update(msg)
		m.sleep, _ = m.sleep.Update(msg)
		m.workouts, _ = m.workouts.Update(msg)
		m.profile, _ = m.profile.Update(msg)
		return m, nil

	case tea.KeyMsg:
		// Global keys
		if key.Matches(msg, m.keys.Quit) {
			return m, tea.Quit
		}

		if !m.authenticated {
			if msg.String() == "l" {
				m.login, _ = m.login.Update(screens.LoginMsg{})
				return m, nil // In real app, this would trigger the auth flow
			}
			return m, nil
		}

		// Tab navigation
		var cmd tea.Cmd
		switch {
		case key.Matches(msg, m.keys.NextTab):
			m.activeTab = (m.activeTab + 1) % 5
			cmd = m.initScreenIfNeeded()
		case key.Matches(msg, m.keys.PrevTab):
			m.activeTab = (m.activeTab + 4) % 5
			cmd = m.initScreenIfNeeded()
		case key.Matches(msg, m.keys.Tab1):
			m.activeTab = TabDashboard
			cmd = m.initScreenIfNeeded()
		case key.Matches(msg, m.keys.Tab2):
			m.activeTab = TabRecovery
			cmd = m.initScreenIfNeeded()
		case key.Matches(msg, m.keys.Tab3):
			m.activeTab = TabSleep
			cmd = m.initScreenIfNeeded()
		case key.Matches(msg, m.keys.Tab4):
			m.activeTab = TabWorkouts
			cmd = m.initScreenIfNeeded()
		case key.Matches(msg, m.keys.Tab5):
			m.activeTab = TabProfile
			cmd = m.initScreenIfNeeded()
		case key.Matches(msg, m.keys.Refresh):
			m.initialized[m.activeTab] = false
			cmd = m.initScreenIfNeeded()
		default:
			// Forward to active screen
			cmd = m.updateActiveScreen(msg)
		}
		return m, cmd

	case screens.LoginSuccessMsg:
		m.authenticated = true
		m.login, _ = m.login.Update(msg)
		m.initialized[TabDashboard] = true
		return m, m.dashboard.Init()
	}

	// Forward data messages to appropriate screens
	cmd := m.updateActiveScreen(msg)
	return m, cmd
}

func (m *AppModel) initScreenIfNeeded() tea.Cmd {
	if m.initialized[m.activeTab] {
		return nil
	}
	m.initialized[m.activeTab] = true
	switch m.activeTab {
	case TabDashboard:
		return m.dashboard.Init()
	case TabRecovery:
		return m.recovery.Init()
	case TabSleep:
		return m.sleep.Init()
	case TabWorkouts:
		return m.workouts.Init()
	case TabProfile:
		return m.profile.Init()
	}
	return nil
}

func (m *AppModel) updateActiveScreen(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch m.activeTab {
	case TabDashboard:
		m.dashboard, cmd = m.dashboard.Update(msg)
	case TabRecovery:
		m.recovery, cmd = m.recovery.Update(msg)
	case TabSleep:
		m.sleep, cmd = m.sleep.Update(msg)
	case TabWorkouts:
		m.workouts, cmd = m.workouts.Update(msg)
	case TabProfile:
		m.profile, cmd = m.profile.Update(msg)
	}
	return cmd
}

func (m AppModel) View() string {
	if !m.authenticated {
		return m.login.View()
	}

	tabs := components.RenderTabs(m.activeTab, m.width)

	var content string
	switch m.activeTab {
	case TabDashboard:
		content = m.dashboard.View()
	case TabRecovery:
		content = m.recovery.View()
	case TabSleep:
		content = m.sleep.View()
	case TabWorkouts:
		content = m.workouts.View()
	case TabProfile:
		content = m.profile.View()
	}

	statusBar := components.StatusBar(m.statusMsg, m.width)

	return lipgloss.JoinVertical(lipgloss.Left, tabs, "", content, statusBar)
}
