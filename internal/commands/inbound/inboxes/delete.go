package inboxes

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdDelete(f *cmdutil.Factory) *cobra.Command {
	var (
		folderID string
		inboxID  string
	)

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an inbound inbox",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("folder-id", folderID); err != nil {
				return err
			}
			if err := cmdutil.RequireFlag("id", inboxID); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/api/inbound/folders/%s/inboxes/%s", folderID, inboxID)

			if err := c.Delete(context.Background(), client.BaseGeneral, path, nil); err != nil {
				return err
			}

			fmt.Fprintln(f.IOStreams.Out, "Inbound inbox deleted successfully.")
			return nil
		},
	}

	cmd.Flags().StringVar(&folderID, "folder-id", "", "Folder ID (required)")
	cmd.Flags().StringVar(&inboxID, "id", "", "Inbox ID (required)")

	return cmd
}
