package organizations

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdCreateSubAccount(f *cmdutil.Factory) *cobra.Command {
	var orgID string
	var name string

	cmd := &cobra.Command{
		Use:   "create-sub-account",
		Short: "Create a sub-account for an organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("org-id", orgID); err != nil {
				return err
			}
			if err := cmdutil.RequireFlag("name", name); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/api/organizations/%s/sub_accounts", orgID)

			body := map[string]interface{}{
				"account": map[string]string{
					"name": name,
				},
			}

			var subAccount SubAccount
			if err := c.Post(context.Background(), client.BaseGeneral, path, body, &subAccount); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, subAccount, subAccountColumns)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID")
	cmd.Flags().StringVar(&name, "name", "", "Sub-account name")

	return cmd
}
