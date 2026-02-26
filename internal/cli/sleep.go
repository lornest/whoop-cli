package cli

import (
	"fmt"
	"time"

	"github.com/lornest/whoop-cli/pkg/whoop"
	"github.com/spf13/cobra"
)

var sleepCmd = &cobra.Command{
	Use:   "sleep",
	Short: "Show recent sleep data",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		start := time.Now().AddDate(0, 0, -daysFlag).Format("2006-01-02T15:04:05.000Z")
		resp, err := client.GetSleep(&whoop.QueryParams{Start: start, Limit: daysFlag})
		if err != nil {
			return fmt.Errorf("fetch sleep: %w", err)
		}

		return formatOutput(formatFlag, "sleep", resp.Records)
	},
}

func init() {
	rootCmd.AddCommand(sleepCmd)
}
