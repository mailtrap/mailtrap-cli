package inboxes

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var (
		folderID string
		name     string
		domainID int
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an inbound inbox",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("folder-id", folderID); err != nil {
				return err
			}
			if err := cmdutil.RequireFlag("name", name); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/api/inbound/folders/%s/inboxes", folderID)

			// Omit domain-id for a Mailtrap-hosted inbox; set it to create a custom-domain (catch-all) inbox.
			body := map[string]interface{}{"name": name}
			if cmd.Flags().Changed("domain-id") {
				body["domain_id"] = domainID
			}

			var resp InboundInbox
			if err := c.Post(context.Background(), client.BaseGeneral, path, body, &resp); err != nil {
				return err
			}

			return output.Print(f.IOStreams.Out, cmdutil.GetOutputFormat(), resp, inboxColumns)
		},
	}

	cmd.Flags().StringVar(&folderID, "folder-id", "", "Folder ID (required)")
	cmd.Flags().StringVar(&name, "name", "", "Inbox name (required)")
	cmd.Flags().IntVar(&domainID, "domain-id", 0, "Sending domain ID for a custom-domain (catch-all) inbox")

	return cmd
}
