package account_access

import (
	"context"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type AccountAccess struct {
	ID          int    `json:"id"`
	UserID      int    `json:"specifier_id"`
	UserEmail   string `json:"specifier_email"`
	AccessLevel int    `json:"access_level"`
	CreatedAt   string `json:"created_at"`
}

var accountAccessColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "USER_ID", Field: "specifier_id"},
	{Header: "USER_EMAIL", Field: "specifier_email"},
	{Header: "ACCESS_LEVEL", Field: "access_level"},
	{Header: "CREATED_AT", Field: "created_at"},
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all account accesses",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			_, err = config.RequireAccountID()
			if err != nil {
				return err
			}

			path := cmdutil.AccountPath("account_accesses")

			var accesses []AccountAccess
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &accesses); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, accesses, accountAccessColumns)
		},
	}

	return cmd
}
