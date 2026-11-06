package organizations

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdDeleteSubAccount(f *cmdutil.Factory) *cobra.Command {
	var orgID int64
	var subAccountID int64

	cmd := &cobra.Command{
		Use:   "delete-sub-account",
		Args:  cobra.NoArgs,
		Short: "Delete a sub-account from an organization",
		Long: `Delete a sub-account from an organization. Requires organization admin permissions.
The organization must have the Organization & Sub-accounts feature enabled; otherwise the command fails with 403.

The deletion is permanent: the sub-account and all of its data are removed and cannot be restored.
Deleting the organization's last sub-account also deletes the organization.

A repeated call for the same sub-account returns 404. A concurrent deletion of the same sub-account may fail with 409. Rate limit: 10 requests per minute per organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("org-id") {
				return fmt.Errorf("--org-id is required")
			}
			if orgID <= 0 {
				return fmt.Errorf("--org-id must be greater than 0")
			}
			if !cmd.Flags().Changed("sub-account-id") {
				return fmt.Errorf("--sub-account-id is required")
			}
			if subAccountID <= 0 {
				return fmt.Errorf("--sub-account-id must be greater than 0")
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/api/organizations/%d/sub_accounts/%d", orgID, subAccountID)

			if err := c.Delete(context.Background(), client.BaseGeneral, path, nil); err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "Sub-account %d deleted successfully\n", subAccountID)

			return nil
		},
	}

	cmd.Flags().Int64Var(&orgID, "org-id", 0, "Organization ID")
	cmd.Flags().Int64Var(&subAccountID, "sub-account-id", 0, "Sub-account ID")

	return cmd
}
