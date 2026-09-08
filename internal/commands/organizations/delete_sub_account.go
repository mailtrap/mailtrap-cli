package organizations

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdDeleteSubAccount(f *cmdutil.Factory) *cobra.Command {
	var orgID string
	var subAccountID string

	cmd := &cobra.Command{
		Use:   "delete-sub-account",
		Short: "Delete a sub-account from an organization",
		Long: `Delete a sub-account from an organization. Requires sub-account management permissions for the organization.

The deletion is permanent: the sub-account and all of its data are removed and cannot be restored.
Deleting the organization's last sub-account also deletes the organization.

A repeated call for the same sub-account returns 404. Rate limit: 10 requests per minute per organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("org-id", orgID); err != nil {
				return err
			}
			if err := cmdutil.RequireFlag("sub-account-id", subAccountID); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/api/organizations/%s/sub_accounts/%s", orgID, subAccountID)

			if err := c.Delete(context.Background(), client.BaseGeneral, path, nil); err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "Sub-account %s deleted successfully\n", subAccountID)

			return nil
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID")
	cmd.Flags().StringVar(&subAccountID, "sub-account-id", "", "Sub-account ID")

	return cmd
}
