package messages

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

var messageDetailColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "FROM", Field: "from"},
	{Header: "TO", Field: "to"},
	{Header: "CC", Field: "cc"},
	{Header: "SUBJECT", Field: "subject"},
	{Header: "SIZE", Field: "size"},
	{Header: "RECEIVED AT", Field: "received_at"},
	{Header: "THREAD ID", Field: "thread_id"},
}

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	var (
		inboxID   string
		messageID string
	)

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get an inbound message with its body and attachment download URLs",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("inbox-id", inboxID); err != nil {
				return err
			}
			if err := cmdutil.RequireFlag("id", messageID); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/api/inbound/inboxes/%s/messages/%s", inboxID, messageID)

			var resp json.RawMessage
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &resp); err != nil {
				return err
			}

			return output.Print(f.IOStreams.Out, cmdutil.GetOutputFormat(), resp, messageDetailColumns)
		},
	}

	cmd.Flags().StringVar(&inboxID, "inbox-id", "", "Inbox ID (required)")
	cmd.Flags().StringVar(&messageID, "id", "", "Message ID (required)")

	return cmd
}
