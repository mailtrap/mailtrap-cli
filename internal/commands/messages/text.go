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

func NewCmdText(f *cmdutil.Factory) *cobra.Command {
	var inboxID string
	var messageID string

	cmd := &cobra.Command{
		Use:   "text",
		Short: "Get text body of a message",
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

			path := cmdutil.AccountPath("inboxes", fmt.Sprintf("%s", inboxID), "messages", fmt.Sprintf("%s", messageID), "body.txt")

			data, err := c.GetRaw(context.Background(), client.BaseGeneral, path, nil)
			if err != nil {
				return err
			}

			return output.PrintRaw(f.IOStreams.Out, data)
		},
	}

	cmd.Flags().StringVar(&inboxID, "inbox-id", "", "Inbox ID")
	cmd.Flags().StringVar(&messageID, "id", "", "Message ID")

	return cmd
}
