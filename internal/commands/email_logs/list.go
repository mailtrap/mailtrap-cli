package emaillogs

import (
	"context"
	"fmt"
	"net/url"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type EmailLog struct {
	ID        string `json:"id"`
	Subject   string `json:"subject"`
	FromEmail string `json:"from_email"`
	ToEmail   string `json:"to_email"`
	Status    string `json:"status"`
	SentAt    string `json:"sent_at"`
}

var emailLogColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "SUBJECT", Field: "subject"},
	{Header: "FROM", Field: "from_email"},
	{Header: "TO", Field: "to_email"},
	{Header: "STATUS", Field: "status"},
	{Header: "SENT AT", Field: "sent_at"},
}

type ListOptions struct {
	Page    int
	PerPage int
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List email logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			if _, err := config.RequireAccountID(); err != nil {
				return err
			}

			path := cmdutil.AccountPath("email_logs")

			params := url.Values{}
			if cmd.Flags().Changed("page") {
				params.Set("page", fmt.Sprintf("%d", opts.Page))
			}
			if cmd.Flags().Changed("per-page") {
				params.Set("per_page", fmt.Sprintf("%d", opts.PerPage))
			}

			if len(params) > 0 {
				path = fmt.Sprintf("%s?%s", path, params.Encode())
			}

			var result []EmailLog
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &result); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			output.Print(f.IOStreams.Out, format, result, emailLogColumns)

			return nil
		},
	}

	cmd.Flags().IntVar(&opts.Page, "page", 0, "Page number")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "Results per page")

	return cmd
}
