package contact_fields

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
	var fieldType string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a contact field",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("name", name); err != nil {
				return err
			}
			if err := cmdutil.RequireFlag("field-type", fieldType); err != nil {
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

			path := cmdutil.AccountPath("contacts", "fields")

			body := map[string]interface{}{
				"contact_field": map[string]string{
					"name":       name,
					"field_type": fieldType,
				},
			}

			var field ContactField
			if err := c.Post(context.Background(), client.BaseGeneral, path, body, &field); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, field, contactFieldColumns)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Contact field name")
	cmd.Flags().StringVar(&fieldType, "field-type", "", "Contact field type")

	return cmd
}
