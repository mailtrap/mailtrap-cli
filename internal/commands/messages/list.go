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

type Message struct {
	ID        int    `json:"id"`
	Subject   string `json:"subject"`
	FromEmail string `json:"from_email"`
	ToEmail   string `json:"to_email"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at"`
}

var messageColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "SUBJECT", Field: "subject"},
	{Header: "FROM_EMAIL", Field: "from_email"},
	{Header: "TO_EMAIL", Field: "to_email"},
	{Header: "IS_READ", Field: "is_read"},
	{Header: "CREATED_AT", Field: "created_at"},
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var inboxID string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all messages in an inbox",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("inbox-id", inboxID); err != nil {
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

			path := cmdutil.AccountPath("inboxes", fmt.Sprintf("%s", inboxID), "messages")

			var messages []Message
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &messages); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, messages, messageColumns)
		},
	}

	cmd.Flags().StringVar(&inboxID, "inbox-id", "", "Inbox ID")

	return cmd
}
