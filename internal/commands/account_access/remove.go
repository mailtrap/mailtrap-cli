package account_access

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdRemove(f *cmdutil.Factory) *cobra.Command {
	var accessID string

	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove an account access",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", accessID); err != nil {
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

			path := cmdutil.AccountPath("account_accesses", fmt.Sprintf("%s", accessID))

			if err := c.Delete(context.Background(), client.BaseGeneral, path, nil); err != nil {
				return err
			}

			fmt.Fprintln(f.IOStreams.Out, "Account access removed successfully.")
			return nil
		},
	}

	cmd.Flags().StringVar(&accessID, "id", "", "Account access ID")

	return cmd
}
