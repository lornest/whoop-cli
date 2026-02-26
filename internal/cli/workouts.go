package cli

import (
	"fmt"
	"time"

	"github.com/lornest/whoop-cli/pkg/whoop"
	"github.com/spf13/cobra"
)

var workoutsCmd = &cobra.Command{
	Use:   "workouts",
	Short: "Show recent workouts",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		start := time.Now().AddDate(0, 0, -daysFlag).Format("2006-01-02T15:04:05.000Z")
		resp, err := client.GetWorkouts(&whoop.QueryParams{Start: start, Limit: 25})
		if err != nil {
			return fmt.Errorf("fetch workouts: %w", err)
		}

		return formatOutput(formatFlag, "workouts", resp.Records)
	},
}

func init() {
	rootCmd.AddCommand(workoutsCmd)
}
