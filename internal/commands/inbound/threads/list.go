package threads

import (
	"context"
	"fmt"
	"net/url"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

// InboundThreadMessage represents a message inside a thread. Populated on get-by-id.
type InboundThreadMessage struct {
	ID               string `json:"id,omitempty"`
	VisibilityStatus string `json:"visibility_status,omitempty"`
	Direction        string `json:"direction,omitempty"`
	Subject          string `json:"subject,omitempty"`
	From             string `json:"from,omitempty"`
	CreatedAt        string `json:"created_at,omitempty"`
	DeliveryStatus   string `json:"delivery_status,omitempty"`
}

// InboundThread represents a conversation thread. Messages are populated on get-by-id.
type InboundThread struct {
	ID             string                 `json:"id"`
	Subject        string                 `json:"subject,omitempty"`
	MessageCount   *int                   `json:"message_count,omitempty"`
	Size           *int                   `json:"size,omitempty"`
	LastActivityAt string                 `json:"last_activity_at,omitempty"`
	Senders        []string               `json:"senders,omitempty"`
	Recipients     []string               `json:"recipients,omitempty"`
	Messages       []InboundThreadMessage `json:"messages,omitempty"`
}

type threadsListResponse struct {
	Data       []InboundThread `json:"data"`
	TotalCount int             `json:"total_count"`
	LastID     string          `json:"last_id"`
}

var threadColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "SUBJECT", Field: "subject"},
	{Header: "MESSAGES", Field: "message_count"},
	{Header: "LAST ACTIVITY", Field: "last_activity_at"},
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var (
		inboxID string
		lastID  string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List conversation threads in an inbox",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("inbox-id", inboxID); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/api/inbound/inboxes/%s/threads", inboxID)

			var params url.Values
			if lastID != "" {
				params = url.Values{}
				params.Set("last_id", lastID)
			}

			var resp threadsListResponse
			if err := c.Get(context.Background(), client.BaseGeneral, path, params, &resp); err != nil {
				return err
			}

			return output.Print(f.IOStreams.Out, cmdutil.GetOutputFormat(), resp.Data, threadColumns)
		},
	}

	cmd.Flags().StringVar(&inboxID, "inbox-id", "", "Inbox ID (required)")
	cmd.Flags().StringVar(&lastID, "last-id", "", "Pagination cursor (last_id from previous response)")

	return cmd
}
