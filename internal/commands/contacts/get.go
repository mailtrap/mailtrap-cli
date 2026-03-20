package contacts

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type Contact struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

var contactColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "EMAIL", Field: "email"},
	{Header: "FIRST_NAME", Field: "first_name"},
	{Header: "LAST_NAME", Field: "last_name"},
	{Header: "CREATED_AT", Field: "created_at"},
	{Header: "UPDATED_AT", Field: "updated_at"},
}

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	var contactID string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a contact",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", contactID); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			_, err = config.RequireAccountID()
			if err != nil {
				return err
			}

			path := cmdutil.AccountPath("contacts", fmt.Sprintf("%s", contactID))

			var contact Contact
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &contact); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, contact, contactColumns)
		},
	}

	cmd.Flags().StringVar(&contactID, "id", "", "Contact ID")

	return cmd
}
