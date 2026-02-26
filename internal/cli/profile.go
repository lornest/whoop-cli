package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Show body measurements",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		body, err := client.GetBodyMeasurement()
		if err != nil {
			return fmt.Errorf("fetch profile: %w", err)
		}

		return formatOutput(formatFlag, "profile", body)
	},
}

func init() {
	rootCmd.AddCommand(profileCmd)
}
