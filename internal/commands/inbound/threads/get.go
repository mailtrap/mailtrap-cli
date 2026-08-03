package threads

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
		inboxID  string
		threadID string
	)

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a conversation thread with its messages",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("inbox-id", inboxID); err != nil {
				return err
			}
			if err := cmdutil.RequireFlag("id", threadID); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/api/inbound/inboxes/%s/threads/%s", inboxID, threadID)

			var resp InboundThread
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &resp); err != nil {
				return err
			}

			return output.Print(f.IOStreams.Out, cmdutil.GetOutputFormat(), resp, threadColumns)
		},
	}

	cmd.Flags().StringVar(&inboxID, "inbox-id", "", "Inbox ID (required)")
	cmd.Flags().StringVar(&threadID, "id", "", "Thread ID (required)")

	return cmd
}
