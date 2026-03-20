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

type Attachment struct {
	ID             int    `json:"id"`
	Filename       string `json:"filename"`
	AttachmentType string `json:"attachment_type"`
	ContentType    string `json:"content_type"`
	AttachmentSize int    `json:"attachment_size"`
	DownloadPath   string `json:"download_path"`
}

var attachmentColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "FILENAME", Field: "filename"},
	{Header: "ATTACHMENT_TYPE", Field: "attachment_type"},
	{Header: "CONTENT_TYPE", Field: "content_type"},
	{Header: "ATTACHMENT_SIZE", Field: "attachment_size"},
	{Header: "DOWNLOAD_PATH", Field: "download_path"},
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var inboxID string
	var messageID string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all attachments of a message",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("inbox-id", inboxID); err != nil {
				return err
			}
			if err := cmdutil.RequireFlag("message-id", messageID); err != nil {
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

			path := cmdutil.AccountPath("inboxes", fmt.Sprintf("%s", inboxID), "messages", fmt.Sprintf("%s", messageID), "attachments")

			var attachments []Attachment
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &attachments); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, attachments, attachmentColumns)
		},
	}

	cmd.Flags().StringVar(&inboxID, "inbox-id", "", "Inbox ID")
	cmd.Flags().StringVar(&messageID, "message-id", "", "Message ID")

	return cmd
}
