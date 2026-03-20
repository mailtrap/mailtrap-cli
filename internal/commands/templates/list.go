package templates

import (
	"context"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type Template struct {
	ID        int    `json:"id"`
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Subject   string `json:"subject"`
	Category  string `json:"category"`
	CreatedAt string `json:"created_at"`
}

var templateColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "UUID", Field: "uuid"},
	{Header: "NAME", Field: "name"},
	{Header: "SUBJECT", Field: "subject"},
	{Header: "CATEGORY", Field: "category"},
	{Header: "CREATED AT", Field: "created_at"},
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all email templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			if _, err := config.RequireAccountID(); err != nil {
				return err
			}

			path := cmdutil.AccountPath("email_templates")

			var result []Template
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &result); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			output.Print(f.IOStreams.Out, format, result, templateColumns)

			return nil
		},
	}

	return cmd
}
