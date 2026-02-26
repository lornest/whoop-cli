package cli

import (
	"fmt"
	"time"

	"github.com/lornest/whoop-cli/pkg/whoop"
	"github.com/spf13/cobra"
)

var recoveryCmd = &cobra.Command{
	Use:   "recovery",
	Short: "Show recent recovery scores",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		start := time.Now().AddDate(0, 0, -daysFlag).Format("2006-01-02T15:04:05.000Z")
		resp, err := client.GetRecovery(&whoop.QueryParams{Start: start, Limit: daysFlag})
		if err != nil {
			return fmt.Errorf("fetch recovery: %w", err)
		}

		return formatOutput(formatFlag, "recovery", resp.Records)
	},
}

func init() {
	rootCmd.AddCommand(recoveryCmd)
}
