package contact_lists

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	var listID string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a contact list",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", listID); err != nil {
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

			path := cmdutil.AccountPath("contacts", "lists", fmt.Sprintf("%s", listID))

			var list ContactList
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &list); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, list, contactListColumns)
		},
	}

	cmd.Flags().StringVar(&listID, "id", "", "Contact list ID")

	return cmd
}
