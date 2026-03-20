package accounts

import (
	"context"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type Account struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	AccessLevels []int  `json:"access_levels"`
}

var accountColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "NAME", Field: "name"},
	{Header: "ACCESS_LEVELS", Field: "access_levels"},
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			path := "/api/accounts"

			var accounts []Account
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &accounts); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, accounts, accountColumns)
		},
	}

	return cmd
}
