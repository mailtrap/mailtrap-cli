package configure

import (
	"context"
	"fmt"
	"os"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type configFile struct {
	APIToken  string `yaml:"api-token,omitempty"`
	AccountID string `yaml:"account-id,omitempty"`
}

type account struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewCmdConfigure(f *cmdutil.Factory) *cobra.Command {
	var apiToken string

	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Configure Mailtrap CLI with your API token",
		Long:  "Validates the API token, auto-detects your account ID, and saves both to ~/.mailtrap.yaml.",
		Example: `  mailtrap configure --api-token your-token-here`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if apiToken == "" {
				return fmt.Errorf("--api-token is required")
			}

			// Validate the token by fetching accounts
			c := client.New(apiToken)
			var accounts []account
			if err := c.Get(context.Background(), client.BaseGeneral, "/api/accounts", nil, &accounts); err != nil {
				return fmt.Errorf("invalid API token: %w", err)
			}
			if len(accounts) == 0 {
				return fmt.Errorf("no accounts found for this API token")
			}

			accountID := fmt.Sprintf("%d", accounts[0].ID)
			accountName := accounts[0].Name

			path := config.ConfigFilePath()
			cfg := configFile{
				APIToken:  apiToken,
				AccountID: accountID,
			}

			out, err := yaml.Marshal(&cfg)
			if err != nil {
				return fmt.Errorf("failed to marshal config: %w", err)
			}

			if err := os.WriteFile(path, out, 0600); err != nil {
				return fmt.Errorf("failed to write config file: %w", err)
			}

			fmt.Fprintf(f.IOStreams.Out, "Authenticated as %q (account ID: %s)\n", accountName, accountID)
			fmt.Fprintf(f.IOStreams.Out, "Configuration saved to %s\n", path)
			return nil
		},
	}

	cmd.Flags().StringVar(&apiToken, "api-token", "", "Mailtrap API token")

	return cmd
}
