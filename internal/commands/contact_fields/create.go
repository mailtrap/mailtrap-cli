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
	var dataType string
	var mergeTag string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a contact field",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("name", name); err != nil {
				return err
			}
			if err := cmdutil.RequireFlag("data-type", dataType); err != nil {
				return err
			}
			if err := cmdutil.RequireFlag("merge-tag", mergeTag); err != nil {
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
				"name":      name,
				"data_type": dataType,
				"merge_tag": mergeTag,
			}

			var field ContactField
			if err := c.Post(context.Background(), client.BaseGeneral, path, body, &field); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, field, contactFieldColumns)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Contact field name (required)")
	cmd.Flags().StringVar(&dataType, "data-type", "", "Data type: text, integer, float, boolean, date (required)")
	cmd.Flags().StringVar(&mergeTag, "merge-tag", "", "Merge tag for the field (required)")

	return cmd
}
