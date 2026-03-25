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
	MessageID string `json:"message_id"`
	Subject   string `json:"subject"`
	From      string `json:"from"`
	To        string `json:"to"`
	Status    string `json:"status"`
	SentAt    string `json:"sent_at"`
}

type emailLogListResponse struct {
	Messages       []EmailLog `json:"messages"`
	TotalCount     int        `json:"total_count"`
	NextPageCursor string     `json:"next_page_cursor"`
}

var emailLogColumns = []output.Column{
	{Header: "MESSAGE ID", Field: "message_id"},
	{Header: "SUBJECT", Field: "subject"},
	{Header: "FROM", Field: "from"},
	{Header: "TO", Field: "to"},
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

			var resp emailLogListResponse
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &resp); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			output.Print(f.IOStreams.Out, format, resp.Messages, emailLogColumns)

			return nil
		},
	}

	cmd.Flags().IntVar(&opts.Page, "page", 0, "Page number")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "Results per page")

	return cmd
}
