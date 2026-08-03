package inboxes

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	var (
		folderID string
		inboxID  string
	)

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get an inbound inbox",
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

			var resp InboundInbox
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &resp); err != nil {
				return err
			}

			return output.Print(f.IOStreams.Out, cmdutil.GetOutputFormat(), resp, inboxColumns)
		},
	}

	cmd.Flags().StringVar(&folderID, "folder-id", "", "Folder ID (required)")
	cmd.Flags().StringVar(&inboxID, "id", "", "Inbox ID (required)")

	return cmd
}
