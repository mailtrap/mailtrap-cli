package suppressions

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdDelete(f *cmdutil.Factory) *cobra.Command {
	var suppressionID string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a suppression",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", suppressionID); err != nil {
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

			path := cmdutil.AccountPath("suppressions", fmt.Sprintf("%s", suppressionID))

			if err := c.Delete(context.Background(), client.BaseGeneral, path, nil); err != nil {
				return err
			}

			fmt.Fprintln(f.IOStreams.Out, "Suppression deleted successfully.")
			return nil
		},
	}

	cmd.Flags().StringVar(&suppressionID, "id", "", "Suppression ID")

	return cmd
}
