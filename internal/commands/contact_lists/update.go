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

func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var listID string
	var name string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a contact list",
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

			body := map[string]interface{}{
				"name": name,
			}

			var list ContactList
			if err := c.Patch(context.Background(), client.BaseGeneral, path, body, &list); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, list, contactListColumns)
		},
	}

	cmd.Flags().StringVar(&listID, "id", "", "Contact list ID")
	cmd.Flags().StringVar(&name, "name", "", "Contact list name")

	return cmd
}
