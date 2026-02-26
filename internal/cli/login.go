package cli

import (
	"fmt"
	"os"

	"github.com/lornest/whoop-cli/pkg/whoop"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Whoop via OAuth2",
	Long:  "Opens a browser to complete the Whoop OAuth2 login flow. Requires WHOOP_CLIENT_ID and WHOOP_CLIENT_SECRET environment variables.",
	RunE: func(cmd *cobra.Command, args []string) error {
		clientID := os.Getenv("WHOOP_CLIENT_ID")
		clientSecret := os.Getenv("WHOOP_CLIENT_SECRET")

		if clientID == "" || clientSecret == "" {
			return fmt.Errorf("set WHOOP_CLIENT_ID and WHOOP_CLIENT_SECRET environment variables")
		}

		storage := whoop.NewStorage()
		_, err := whoop.Login(clientID, clientSecret, storage)
		if err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		fmt.Println("Login successful!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
