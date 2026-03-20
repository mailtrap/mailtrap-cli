package contact_fields

import (
	"context"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type ContactField struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	FieldType string `json:"field_type"`
	CreatedAt string `json:"created_at"`
}

var contactFieldColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "NAME", Field: "name"},
	{Header: "FIELD_TYPE", Field: "field_type"},
	{Header: "CREATED_AT", Field: "created_at"},
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all contact fields",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			_, err = config.RequireAccountID()
			if err != nil {
				return err
			}

			path := cmdutil.AccountPath("contacts", "fields")

			var fields []ContactField
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &fields); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, fields, contactFieldColumns)
		},
	}

	return cmd
}
