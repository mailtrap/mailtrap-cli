package contact_fields

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
	var fieldID string
	var name string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a contact field",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", fieldID); err != nil {
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

			path := cmdutil.AccountPath("contacts", "fields", fmt.Sprintf("%s", fieldID))

			body := map[string]interface{}{
				"contact_field": map[string]string{
					"name": name,
				},
			}

			var field ContactField
			if err := c.Patch(context.Background(), client.BaseGeneral, path, body, &field); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, field, contactFieldColumns)
		},
	}

	cmd.Flags().StringVar(&fieldID, "id", "", "Contact field ID")
	cmd.Flags().StringVar(&name, "name", "", "Contact field name")

	return cmd
}
