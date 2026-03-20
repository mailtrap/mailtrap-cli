package contact_lists

import (
	"context"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a contact list",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("name", name); err != nil {
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

			path := cmdutil.AccountPath("contacts", "lists")

			body := map[string]interface{}{
				"contact_list": map[string]string{
					"name": name,
				},
			}

			var list ContactList
			if err := c.Post(context.Background(), client.BaseGeneral, path, body, &list); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, list, contactListColumns)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Contact list name")

	return cmd
}
