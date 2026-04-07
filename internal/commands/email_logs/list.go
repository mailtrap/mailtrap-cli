package emaillogs

import (
	"context"
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

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var (
		searchAfter string
		sentAfter   string
		sentBefore  string
		toFilter    string
		toOperator  string
		fromFilter  string
		fromOp      string
		subjectVal  string
		subjectOp   string
		statusVal   string
		eventsVal   string
		category    string
	)

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
			if searchAfter != "" {
				params.Set("search_after", searchAfter)
			}
			if sentAfter != "" {
				params.Set("filters[sent_after]", sentAfter)
			}
			if sentBefore != "" {
				params.Set("filters[sent_before]", sentBefore)
			}
			if toFilter != "" {
				params.Set("filters[to][value]", toFilter)
				if toOperator != "" {
					params.Set("filters[to][operator]", toOperator)
				}
			}
			if fromFilter != "" {
				params.Set("filters[from][value]", fromFilter)
				if fromOp != "" {
					params.Set("filters[from][operator]", fromOp)
				}
			}
			if subjectVal != "" {
				params.Set("filters[subject][value]", subjectVal)
				if subjectOp != "" {
					params.Set("filters[subject][operator]", subjectOp)
				}
			}
			if statusVal != "" {
				params.Set("filters[status][value]", statusVal)
			}
			if eventsVal != "" {
				params.Set("filters[events][value]", eventsVal)
			}
			if category != "" {
				params.Set("filters[category]", category)
			}

			var resp emailLogListResponse
			if err := c.Get(context.Background(), client.BaseGeneral, path, params, &resp); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, resp.Messages, emailLogColumns)
		},
	}

	cmd.Flags().StringVar(&searchAfter, "cursor", "", "Pagination cursor (next_page_cursor from previous response)")
	cmd.Flags().StringVar(&sentAfter, "sent-after", "", "Filter: sent after (ISO 8601)")
	cmd.Flags().StringVar(&sentBefore, "sent-before", "", "Filter: sent before (ISO 8601)")
	cmd.Flags().StringVar(&toFilter, "to", "", "Filter by recipient email")
	cmd.Flags().StringVar(&toOperator, "to-operator", "ci_equal", "Operator for --to (ci_equal|ci_not_equal|ci_contain|ci_not_contain)")
	cmd.Flags().StringVar(&fromFilter, "from", "", "Filter by sender email")
	cmd.Flags().StringVar(&fromOp, "from-operator", "ci_equal", "Operator for --from")
	cmd.Flags().StringVar(&subjectVal, "subject", "", "Filter by subject")
	cmd.Flags().StringVar(&subjectOp, "subject-operator", "ci_equal", "Operator for --subject")
	cmd.Flags().StringVar(&statusVal, "status", "", "Filter by status (delivered|not_delivered|enqueued|opted_out)")
	cmd.Flags().StringVar(&eventsVal, "event", "", "Filter by event (delivery|open|click|bounce|spam|unsubscribe)")
	cmd.Flags().StringVar(&category, "category", "", "Filter by category")

	return cmd
}
