package screens

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lornest/whoop-cli/internal/tui/style"
	"github.com/lornest/whoop-cli/internal/util"
	"github.com/lornest/whoop-cli/pkg/whoop"
)

// SportName maps sport IDs to names.
var SportName = map[int]string{
	0: "Activity", 1: "Running", 44: "Weightlifting", 71: "Functional Fitness",
	63: "Cycling", 48: "Swimming", 52: "Yoga", 64: "Walking",
	-1: "Other",
}

func getSportName(id int) string {
	if name, ok := SportName[id]; ok {
		return name
	}
	return fmt.Sprintf("Sport %d", id)
}

// WorkoutsModel shows workout list with detail.
type WorkoutsModel struct {
	client   *whoop.Client
	workouts *whoop.WorkoutResponse
	table    table.Model
	err      error
	loading  bool
	width    int
	height   int
}

func NewWorkoutsModel(client *whoop.Client) WorkoutsModel {
	columns := []table.Column{
		{Title: "Date", Width: 16},
		{Title: "Activity", Width: 20},
		{Title: "Strain", Width: 8},
		{Title: "Avg HR", Width: 8},
		{Title: "Max HR", Width: 8},
		{Title: "Duration", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(style.ColorBorder).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(style.ColorBlue).
		Bold(true)
	t.SetStyles(s)

	return WorkoutsModel{client: client, loading: true, table: t}
}

func (m WorkoutsModel) Init() tea.Cmd {
	return m.fetchWorkouts()
}

func (m WorkoutsModel) fetchWorkouts() tea.Cmd {
	return func() tea.Msg {
		start := time.Now().AddDate(0, 0, -7).Format("2006-01-02T15:04:05.000Z")
		resp, err := m.client.GetWorkouts(&whoop.QueryParams{Start: start, Limit: 25})
		return WorkoutsMsg{Data: resp, Err: err}
	}
}

func (m WorkoutsModel) Update(msg tea.Msg) (WorkoutsModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table.SetWidth(m.width - 8) // adjust for padding
		// Adjust height to leave room for title and detail pane
		tableHeight := m.height - 16
		if tableHeight > 0 {
			m.table.SetHeight(tableHeight)
		}
	case WorkoutsMsg:
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.workouts = msg.Data
			m.populateTable()
		}
	}

	if !m.loading && m.workouts != nil && len(m.workouts.Records) > 0 {
		m.table, cmd = m.table.Update(msg)
	}

	return m, cmd
}

func (m *WorkoutsModel) populateTable() {
	var rows []table.Row
	for _, w := range m.workouts.Records {
		t, _ := time.Parse("2006-01-02T15:04:05.000Z", w.Start)
		date := t.Format("Jan 02 15:04")
		sport := w.SportName
		if sport == "" {
			sport = getSportName(w.SportID)
		}

		strain, avgHR, maxHR, duration := "--", "--", "--", "--"
		if w.Score != nil {
			strain = util.FormatStrain(w.Score.Strain)
			avgHR = fmt.Sprintf("%d", w.Score.AverageHeartRate)
			maxHR = fmt.Sprintf("%d", w.Score.MaxHeartRate)
		}
		if w.End != "" {
			start, _ := time.Parse("2006-01-02T15:04:05.000Z", w.Start)
			end, _ := time.Parse("2006-01-02T15:04:05.000Z", w.End)
			duration = util.MillisToDuration(int(end.Sub(start).Milliseconds()))
		}

		rows = append(rows, table.Row{date, sport, strain, avgHR, maxHR, duration})
	}
	m.table.SetRows(rows)
}

func (m WorkoutsModel) View() string {
	if m.loading {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Foreground(style.ColorDim).Render("Loading workouts..."))
	}

	if m.err != nil {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			style.ErrorStyle.Render("Error: "+m.err.Error()+"\n\nPress 'r' to retry"))
	}

	if m.workouts == nil || len(m.workouts.Records) == 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Foreground(style.ColorDim).Render("No workouts found"))
	}

	title := style.TitleStyle.Render("Workouts — Recent")

	// Detail pane for selected workout
	var detail string
	selectedIdx := m.table.Cursor()
	if selectedIdx >= 0 && selectedIdx < len(m.workouts.Records) {
		w := m.workouts.Records[selectedIdx]
		if w.Score != nil {
			detailWidth := m.width - 8
			if detailWidth < 20 {
				detailWidth = 20
			}
			
			innerStyle := lipgloss.NewStyle().
				Background(style.ColorCardBg).
				Width(detailWidth - 4) // adjust for CardStyle padding

			detailContent := lipgloss.JoinVertical(lipgloss.Left,
				innerStyle.Render(style.LabelStyle.Copy().Background(style.ColorCardBg).Render("Detail")),
				innerStyle.Render(fmt.Sprintf("Calories: %.0f kcal  |  Distance: %.2f mi  |  Altitude Gain: %.0f m",
					util.KilojoulesToCalories(w.Score.Kilojoule),
					util.MetersToMiles(w.Score.DistanceMeter),
					w.Score.AltitudeGainMeter)),
			)

			detail = style.CardStyle.Width(detailWidth).Render(detailContent)
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		lipgloss.NewStyle().Render(m.table.View()),
		"",
		lipgloss.NewStyle().Render(detail))

	return lipgloss.Place(m.width, m.height-4, lipgloss.Left, lipgloss.Top, lipgloss.NewStyle().Padding(1, 2).Render(content))
}
