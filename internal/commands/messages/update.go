package messages

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var inboxID string
	var messageID string
	var isRead bool

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a message",
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

			path := cmdutil.AccountPath("inboxes", fmt.Sprintf("%s", inboxID), "messages", fmt.Sprintf("%s", messageID))
			body := map[string]interface{}{
				"message": map[string]interface{}{"is_read": isRead},
			}

			var message Message
			if err := c.Patch(context.Background(), client.BaseGeneral, path, body, &message); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, message, messageColumns)
		},
	}

	cmd.Flags().StringVar(&inboxID, "inbox-id", "", "Inbox ID")
	cmd.Flags().StringVar(&messageID, "id", "", "Message ID")
	cmd.Flags().BoolVar(&isRead, "is-read", false, "Mark message as read")

	return cmd
}
