package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lornest/whoop-cli/pkg/whoop"
	"github.com/lornest/whoop-cli/internal/tui/style"
	"github.com/lornest/whoop-cli/internal/tui/components"
	"github.com/lornest/whoop-cli/internal/util"
)

// ProfileModel shows body measurements.
type ProfileModel struct {
	client *whoop.Client
	body   *whoop.BodyMeasurement
	err    error
	loading bool
	width   int
	height  int
}

func NewProfileModel(client *whoop.Client) ProfileModel {
	return ProfileModel{client: client, loading: true}
}

func (m ProfileModel) Init() tea.Cmd {
	return m.fetchBody()
}

func (m ProfileModel) fetchBody() tea.Cmd {
	return func() tea.Msg {
		data, err := m.client.GetBodyMeasurement()
		return BodyMsg{Data: data, Err: err}
	}
}

func (m ProfileModel) Update(msg tea.Msg) (ProfileModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case BodyMsg:
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.body = msg.Data
		}
	}
	return m, nil
}

func (m ProfileModel) View() string {
	if m.loading {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Foreground(style.ColorDim).Render("Loading profile..."))
	}

	if m.err != nil {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			style.ErrorStyle.Render("Error: "+m.err.Error()+"\n\nPress 'r' to retry"))
	}

	if m.body == nil {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Foreground(style.ColorDim).Render("No profile data available"))
	}

	title := style.TitleStyle.Render("Profile — Body Measurements")

	heightStr := fmt.Sprintf("%s (%.2f m)", util.MetersToFeetInches(m.body.HeightMeter), m.body.HeightMeter)
	weightStr := fmt.Sprintf("%.0f lbs (%.1f kg)", util.KgToLbs(m.body.WeightKilogram), m.body.WeightKilogram)
	hrStr := fmt.Sprintf("%d bpm", m.body.MaxHeartRate)

	cards := lipgloss.JoinHorizontal(lipgloss.Top,
		components.MetricCard("Height", heightStr, style.ColorBlue, 28),
		"  ",
		components.MetricCard("Weight", weightStr, style.ColorGreen, 28),
		"  ",
		components.MetricCard("Max HR", hrStr, style.ColorRed, 28),
	)

	content := lipgloss.JoinVertical(lipgloss.Left, title, "", cards)
	return lipgloss.Place(m.width, m.height-4, lipgloss.Center, lipgloss.Center, content)
}
