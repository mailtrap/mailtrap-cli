package tokens

import (
	"context"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type APIToken struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Last4Digits string `json:"last_4_digits"`
	CreatedBy   string `json:"created_by"`
	ExpiresAt   string `json:"expires_at"`
	Token       string `json:"token,omitempty"`
}

var tokenColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "NAME", Field: "name"},
	{Header: "LAST_4_DIGITS", Field: "last_4_digits"},
	{Header: "CREATED_BY", Field: "created_by"},
	{Header: "EXPIRES_AT", Field: "expires_at"},
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all API tokens",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			_, err = config.RequireAccountID()
			if err != nil {
				return err
			}

			path := cmdutil.AccountPath("api_tokens")

			var tokens []APIToken
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &tokens); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, tokens, tokenColumns)
		},
	}

	return cmd
}
