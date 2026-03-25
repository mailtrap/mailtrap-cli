package sandboxes

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdClean(f *cmdutil.Factory) *cobra.Command {
	var inboxID string

	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Clean a sandbox",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", inboxID); err != nil {
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

			path := cmdutil.AccountPath("inboxes", fmt.Sprintf("%s", inboxID), "clean")

			if err := c.Patch(context.Background(), client.BaseGeneral, path, nil, nil); err != nil {
				return err
			}

			fmt.Fprintln(f.IOStreams.Out, "Sandbox cleaned successfully.")
			return nil
		},
	}

	cmd.Flags().StringVar(&inboxID, "id", "", "Sandbox ID")

	return cmd
}
