package messages

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdForward(f *cmdutil.Factory) *cobra.Command {
	var (
		inboxID   string
		messageID string
		send      sendFlags
	)

	cmd := &cobra.Command{
		Use:   "forward",
		Short: "Forward an inbound message to new recipients (sends a real email)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("inbox-id", inboxID); err != nil {
				return err
			}
			if err := cmdutil.RequireFlag("id", messageID); err != nil {
				return err
			}
			if len(send.to) == 0 {
				return fmt.Errorf("--to is required")
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			body, err := send.body()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/api/inbound/inboxes/%s/messages/%s/forward", inboxID, messageID)

			var resp SendMessageResult
			if err := c.Post(context.Background(), client.BaseGeneral, path, body, &resp); err != nil {
				return err
			}

			return output.Print(f.IOStreams.Out, cmdutil.GetOutputFormat(), resp, sendResultColumns)
		},
	}

	cmd.Flags().StringVar(&inboxID, "inbox-id", "", "Inbox ID (required)")
	cmd.Flags().StringVar(&messageID, "id", "", "Message ID (required)")
	send.register(cmd)

	return cmd
}
