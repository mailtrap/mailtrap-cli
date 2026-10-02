package messages

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

var messageColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "FROM", Field: "from"},
	{Header: "SUBJECT", Field: "subject"},
	{Header: "RECEIVED AT", Field: "received_at"},
	{Header: "THREAD ID", Field: "thread_id"},
}

var messagesPage = output.Page{
	Items:      "data",
	Total:      "total_count",
	Cursor:     []string{"last_id"},
	CursorFlag: "last-id",
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

			var resp json.RawMessage
			if err := c.Get(context.Background(), client.BaseGeneral, path, params, &resp); err != nil {
				return err
			}

			return output.PrintPage(f.IOStreams.Out, cmdutil.GetOutputFormat(), resp, messagesPage, messageColumns)
		},
	}

	cmd.Flags().StringVar(&inboxID, "inbox-id", "", "Inbox ID (required)")
	cmd.Flags().StringVar(&lastID, "last-id", "", "Pagination cursor (last_id from previous response)")

	return cmd
}
