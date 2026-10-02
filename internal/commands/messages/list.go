package messages

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

var messageColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "SUBJECT", Field: "subject"},
	{Header: "FROM_EMAIL", Field: "from_email"},
	{Header: "TO_EMAIL", Field: "to_email"},
	{Header: "IS_READ", Field: "is_read"},
	{Header: "CREATED_AT", Field: "created_at"},
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var (
		sandboxID string
		lastID    string
		page      int
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all messages in a sandbox",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("sandbox-id", sandboxID); err != nil {
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

			path := cmdutil.AccountPath("inboxes", fmt.Sprintf("%s", sandboxID), "messages")

			query := url.Values{}
			if lastID != "" {
				query.Set("last_id", lastID)
			}
			if cmd.Flags().Changed("page") {
				query.Set("page", strconv.Itoa(page))
			}

			var messages json.RawMessage
			if err := c.Get(context.Background(), client.BaseGeneral, path, query, &messages); err != nil {
				return err
			}

			return output.Print(f.IOStreams.Out, cmdutil.GetOutputFormat(), messages, messageColumns)
		},
	}

	cmd.Flags().StringVar(&sandboxID, "sandbox-id", "", "Sandbox ID")
	cmd.Flags().StringVar(&lastID, "last-id", "", "Pagination cursor (id of the last message from the previous response)")
	cmd.Flags().IntVar(&page, "page", 1, "Page number to retrieve (ignored when --last-id is set)")

	return cmd
}
