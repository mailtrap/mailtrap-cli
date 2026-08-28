package trackingoptouts

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdDelete(f *cmdutil.Factory) *cobra.Command {
	var trackingOptOutID string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Remove an email address from the tracking opt-out list",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", trackingOptOutID); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			path := trackingOptOutsPath + "/" + trackingOptOutID

			if err := c.Delete(context.Background(), client.BaseGeneral, path, nil); err != nil {
				return err
			}

			fmt.Fprintln(f.IOStreams.Out, "Tracking opt-out deleted successfully.")
			return nil
		},
	}

	cmd.Flags().StringVar(&trackingOptOutID, "id", "", "Tracking opt-out ID (required)")

	return cmd
}
