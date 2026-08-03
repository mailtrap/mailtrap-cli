package messages

import (
	"context"
	"fmt"
	"net/url"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

// InboundMessage represents a received inbound message. Body fields
// (html_body, text_body) are populated only on get-by-id.
type InboundMessage struct {
	ID         string   `json:"id"`
	InboxID    *int     `json:"inbox_id,omitempty"`
	From       string   `json:"from,omitempty"`
	To         []string `json:"to,omitempty"`
	Cc         []string `json:"cc,omitempty"`
	Subject    string   `json:"subject,omitempty"`
	Size       *int     `json:"size,omitempty"`
	ReceivedAt string   `json:"received_at,omitempty"`
	ThreadID   string   `json:"thread_id,omitempty"`
	HTMLBody   string   `json:"html_body,omitempty"`
	TextBody   string   `json:"text_body,omitempty"`
}

type messagesListResponse struct {
	Data       []InboundMessage `json:"data"`
	TotalCount int              `json:"total_count"`
	LastID     string           `json:"last_id"`
}

var messageColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "FROM", Field: "from"},
	{Header: "SUBJECT", Field: "subject"},
	{Header: "RECEIVED AT", Field: "received_at"},
	{Header: "THREAD ID", Field: "thread_id"},
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var (
		inboxID string
		lastID  string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List received messages in an inbox",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("inbox-id", inboxID); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/api/inbound/inboxes/%s/messages", inboxID)

			var params url.Values
			if lastID != "" {
				params = url.Values{}
				params.Set("last_id", lastID)
			}

			var resp messagesListResponse
			if err := c.Get(context.Background(), client.BaseGeneral, path, params, &resp); err != nil {
				return err
			}

			return output.Print(f.IOStreams.Out, cmdutil.GetOutputFormat(), resp.Data, messageColumns)
		},
	}

	cmd.Flags().StringVar(&inboxID, "inbox-id", "", "Inbox ID (required)")
	cmd.Flags().StringVar(&lastID, "last-id", "", "Pagination cursor (last_id from previous response)")

	return cmd
}
