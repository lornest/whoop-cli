package cli

import (
	"fmt"
	"time"

	"github.com/lornest/whoop-cli/pkg/whoop"
	"github.com/spf13/cobra"
)

var cyclesCmd = &cobra.Command{
	Use:   "cycles",
	Short: "Show recent physiological cycles",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		start := time.Now().AddDate(0, 0, -daysFlag).Format("2006-01-02T15:04:05.000Z")
		resp, err := client.GetCycles(&whoop.QueryParams{Start: start, Limit: daysFlag})
		if err != nil {
			return fmt.Errorf("fetch cycles: %w", err)
		}

		return formatOutput(formatFlag, "cycles", resp.Records)
	},
}

func init() {
	rootCmd.AddCommand(cyclesCmd)
}
