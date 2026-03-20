package attachments

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	var inboxID string
	var messageID string
	var attachmentID string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get an attachment",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("inbox-id", inboxID); err != nil {
				return err
			}
			if err := cmdutil.RequireFlag("message-id", messageID); err != nil {
				return err
			}
			if err := cmdutil.RequireFlag("id", attachmentID); err != nil {
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

			path := cmdutil.AccountPath("inboxes", fmt.Sprintf("%s", inboxID), "messages", fmt.Sprintf("%s", messageID), "attachments", fmt.Sprintf("%s", attachmentID))

			var attachment Attachment
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &attachment); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, attachment, attachmentColumns)
		},
	}

	cmd.Flags().StringVar(&inboxID, "inbox-id", "", "Inbox ID")
	cmd.Flags().StringVar(&messageID, "message-id", "", "Message ID")
	cmd.Flags().StringVar(&attachmentID, "id", "", "Attachment ID")

	return cmd
}
