package screens

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lornest/whoop-cli/internal/tui/style"
)

// LoginMsg signals that the user wants to start the login flow.
type LoginMsg struct{}

// LoginSuccessMsg signals that login completed successfully.
type LoginSuccessMsg struct{}

// LoginErrorMsg signals that login failed.
type LoginErrorMsg struct {
	Err error
}

// LoginModel is the login screen.
type LoginModel struct {
	spinner  spinner.Model
	loading  bool
	err      error
	width    int
	height   int
}

// NewLoginModel creates a new login screen.
func NewLoginModel() LoginModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(style.ColorBlue)
	return LoginModel{spinner: s}
}

func (m LoginModel) Init() tea.Cmd {
	return nil
}

func (m LoginModel) Update(msg tea.Msg) (LoginModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case LoginMsg:
		m.loading = true
		m.err = nil
		return m, m.spinner.Tick
	case LoginSuccessMsg:
		m.loading = false
	case LoginErrorMsg:
		m.loading = false
		m.err = msg.Err
	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m LoginModel) View() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(style.ColorGreen).
		Bold(true).
		MarginBottom(1)
		
	dialogBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(style.ColorBorder).
		Padding(1, 4)

	var content string

	if m.loading {
		content = lipgloss.JoinVertical(lipgloss.Center,
			titleStyle.Render("WHOOP CLI"),
			"",
			m.spinner.View()+" Waiting for browser authentication...",
			"",
			lipgloss.NewStyle().Foreground(style.ColorDim).Render("Complete the login in your browser"),
		)
	} else if m.err != nil {
		content = lipgloss.JoinVertical(lipgloss.Center,
			titleStyle.Render("WHOOP CLI"),
			"",
			style.ErrorStyle.Render("Login failed: "+m.err.Error()),
			"",
			lipgloss.NewStyle().Foreground(style.ColorDim).Render("Press 'l' to try again or 'q' to quit"),
		)
	} else {
		content = lipgloss.JoinVertical(lipgloss.Center,
			titleStyle.Render("WHOOP CLI"),
			"",
			lipgloss.NewStyle().Foreground(style.ColorWhite).Render("Welcome! You need to log in to continue."),
			"",
			lipgloss.NewStyle().Foreground(style.ColorDim).Render("Press 'l' to login or 'q' to quit"),
		)
	}

	dialog := dialogBoxStyle.Render(content)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, dialog)
}

