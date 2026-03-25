package sandboxes

import (
	"context"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type Inbox struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	EmailUsername string `json:"email_username"`
	Status        string `json:"status"`
	MaxSize       int    `json:"max_size"`
}

var inboxColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "NAME", Field: "name"},
	{Header: "EMAIL_USERNAME", Field: "email_username"},
	{Header: "STATUS", Field: "status"},
	{Header: "MAX_SIZE", Field: "max_size"},
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all sandboxes",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			_, err = config.RequireAccountID()
			if err != nil {
				return err
			}

			path := cmdutil.AccountPath("inboxes")

			var inboxes []Inbox
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &inboxes); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, inboxes, inboxColumns)
		},
	}

	return cmd
}
