package contact_fields

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdDelete(f *cmdutil.Factory) *cobra.Command {
	var fieldID string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a contact field",
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

			if err := c.Delete(context.Background(), client.BaseGeneral, path, nil); err != nil {
				return err
			}

			fmt.Fprintln(f.IOStreams.Out, "Contact field deleted successfully.")
			return nil
		},
	}

	cmd.Flags().StringVar(&fieldID, "id", "", "Contact field ID")

	return cmd
}
