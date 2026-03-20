package messages

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdSpamScore(f *cmdutil.Factory) *cobra.Command {
	var inboxID string
	var messageID string

	cmd := &cobra.Command{
		Use:   "spam-score",
		Short: "Get the spam score of a message",
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

			_, err = config.RequireAccountID()
			if err != nil {
				return err
			}

			path := cmdutil.AccountPath("inboxes", fmt.Sprintf("%s", inboxID), "messages", fmt.Sprintf("%s", messageID), "spam_score")

			var result json.RawMessage
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &result); err != nil {
				return err
			}

			return output.Print(f.IOStreams.Out, output.FormatJSON, result, nil)
		},
	}

	cmd.Flags().StringVar(&inboxID, "inbox-id", "", "Inbox ID")
	cmd.Flags().StringVar(&messageID, "id", "", "Message ID")

	return cmd
}
