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

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	var fieldID string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a contact field",
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

			var field ContactField
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &field); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, field, contactFieldColumns)
		},
	}

	cmd.Flags().StringVar(&fieldID, "id", "", "Contact field ID")

	return cmd
}
