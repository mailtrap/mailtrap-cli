package organizations

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type SubAccount struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

var subAccountColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "NAME", Field: "name"},
	{Header: "CREATED_AT", Field: "created_at"},
}

func NewCmdListSubAccounts(f *cmdutil.Factory) *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "list-sub-accounts",
		Short: "List sub-accounts for an organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("org-id", orgID); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/api/organizations/%s/sub_accounts", orgID)

			var subAccounts []SubAccount
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &subAccounts); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, subAccounts, subAccountColumns)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID")

	return cmd
}
