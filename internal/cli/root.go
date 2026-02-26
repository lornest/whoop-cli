package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lornest/whoop-cli/internal/tui"
	"github.com/lornest/whoop-cli/pkg/whoop"
	"github.com/spf13/cobra"
)

var (
	formatFlag string
	daysFlag   int
)

var rootCmd = &cobra.Command{
	Use:   "whoop",
	Short: "Whoop CLI — view your health data from the terminal",
	Long:  "A command-line interface for the Whoop API. Run without arguments to launch the interactive TUI.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTUI()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&formatFlag, "format", "table", "Output format: table, json, or text")
	rootCmd.PersistentFlags().IntVar(&daysFlag, "days", 7, "Number of days of data to fetch")
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newClient() (*whoop.Client, error) {
	storage := whoop.NewStorage()
	_, err := storage.Load()
	if err != nil {
		return nil, fmt.Errorf("not logged in — run 'whoop login' first")
	}
	return whoop.NewClient(storage), nil
}

func runTUI() error {
	storage := whoop.NewStorage()

	authenticated := true
	_, err := storage.Load()
	if err != nil {
		authenticated = false
	}

	client := whoop.NewClient(storage)
	model := tui.NewAppModel(client, authenticated)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}
	return nil
}
