package inboxes

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var (
		folderID string
		inboxID  string
		name     string
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update an inbound inbox",
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

			body := map[string]interface{}{}
			if cmd.Flags().Changed("name") {
				body["name"] = name
			}

			var resp InboundInbox
			if err := c.Patch(context.Background(), client.BaseGeneral, path, body, &resp); err != nil {
				return err
			}

			return output.Print(f.IOStreams.Out, cmdutil.GetOutputFormat(), resp, inboxColumns)
		},
	}

	cmd.Flags().StringVar(&folderID, "folder-id", "", "Folder ID (required)")
	cmd.Flags().StringVar(&inboxID, "id", "", "Inbox ID (required)")
	cmd.Flags().StringVar(&name, "name", "", "Inbox name")

	return cmd
}
