package screens

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lornest/whoop-cli/pkg/whoop"
	"github.com/lornest/whoop-cli/internal/tui/style"
	"github.com/lornest/whoop-cli/internal/util"
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
	err      error
	loading  bool
	selected int
	width    int
	height   int
}

func NewWorkoutsModel(client *whoop.Client) WorkoutsModel {
	return WorkoutsModel{client: client, loading: true}
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
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if m.workouts != nil && len(m.workouts.Records) > 0 {
			switch msg.String() {
			case "j", "down":
				if m.selected < len(m.workouts.Records)-1 {
					m.selected++
				}
			case "k", "up":
				if m.selected > 0 {
					m.selected--
				}
			}
		}
	case WorkoutsMsg:
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.workouts = msg.Data
		}
	}
	return m, nil
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

	// Table header
	header := lipgloss.NewStyle().Foreground(style.ColorDim).Bold(true).Render(
		fmt.Sprintf("  %-14s %-20s %-10s %-10s %-10s %-10s", "Date", "Activity", "Strain", "Avg HR", "Max HR", "Duration"))

	var rows string
	for i, w := range m.workouts.Records {
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

		row := fmt.Sprintf("  %-14s %-20s %-10s %-10s %-10s %-10s", date, sport, strain, avgHR, maxHR, duration)
		if i == m.selected {
			row = lipgloss.NewStyle().Foreground(style.ColorBlue).Bold(true).Render(row)
		}
		rows += "\n" + row
	}

	// Detail pane for selected workout
	var detail string
	if m.selected < len(m.workouts.Records) {
		w := m.workouts.Records[m.selected]
		if w.Score != nil {
			detail = lipgloss.JoinVertical(lipgloss.Left,
				"",
				style.LabelStyle.Render("─── Detail ───"),
				fmt.Sprintf("  Calories: %.0f kcal", util.KilojoulesToCalories(w.Score.Kilojoule)),
				fmt.Sprintf("  Distance: %.2f mi", util.MetersToMiles(w.Score.DistanceMeter)),
				fmt.Sprintf("  Altitude Gain: %.0f m", w.Score.AltitudeGainMeter),
			)
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left, title, "", header, rows, detail)
	return lipgloss.Place(m.width, m.height-4, lipgloss.Center, lipgloss.Center, content)
}
