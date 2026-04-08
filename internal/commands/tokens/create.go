package tokens

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var name string
	var permissions string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an API token",
		Long: `Create an API token with permissions.

Permissions must be provided as a JSON array. Each entry requires:
  - resource_type: "account", "project", "inbox", or "sending_domain"
  - resource_id: the resource ID (integer)
  - access_level: 10 (viewer) or 100 (admin)

Example:
  mailtrap tokens create --name "my-token" --permissions '[{"resource_type":"account","resource_id":12345,"access_level":100}]'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("name", name); err != nil {
				return err
			}
			if err := cmdutil.RequireFlag("permissions", permissions); err != nil {
				return err
			}

			var resources []map[string]interface{}
			if err := json.Unmarshal([]byte(permissions), &resources); err != nil {
				return fmt.Errorf("invalid permissions JSON: %w", err)
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			_, err = config.RequireAccountID()
			if err != nil {
				return err
			}

			path := cmdutil.AccountPath("api_tokens")

			body := map[string]interface{}{
				"name":      name,
				"resources": resources,
			}

			var token APIToken
			if err := c.Post(context.Background(), client.BaseGeneral, path, body, &token); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, token, tokenColumns)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "API token name (required)")
	cmd.Flags().StringVar(&permissions, "permissions", "", `Permissions JSON array (required), e.g. '[{"resource_type":"account","resource_id":123,"access_level":100}]'`)

	return cmd
}
