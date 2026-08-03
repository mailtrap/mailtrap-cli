package inboxes

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

// InboundInbox represents an inbound inbox.
type InboundInbox struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	DomainID *int   `json:"domain_id,omitempty"`
}

var inboxColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "NAME", Field: "name"},
	{Header: "ADDRESS", Field: "address"},
	{Header: "DOMAIN ID", Field: "domain_id"},
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var folderID string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List inboxes in a folder",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("folder-id", folderID); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/api/inbound/folders/%s/inboxes", folderID)

			var resp []InboundInbox
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &resp); err != nil {
				return err
			}

			return output.Print(f.IOStreams.Out, cmdutil.GetOutputFormat(), resp, inboxColumns)
		},
	}

	cmd.Flags().StringVar(&folderID, "folder-id", "", "Folder ID (required)")

	return cmd
}
