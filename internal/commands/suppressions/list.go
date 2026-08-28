package suppressions

import (
	"context"
	"net/url"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type Suppression struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Type          string `json:"type"`
	SendingStream string `json:"sending_stream,omitempty"`
	DomainName    string `json:"domain_name,omitempty"`
	CreatedAt     string `json:"created_at"`
}

var suppressionColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "EMAIL", Field: "email"},
	{Header: "TYPE", Field: "type"},
	{Header: "SENDING_STREAM", Field: "sending_stream"},
	{Header: "DOMAIN_NAME", Field: "domain_name"},
	{Header: "CREATED_AT", Field: "created_at"},
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var email string
	var startTime string
	var endTime string
	var lastID string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List suppressions",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			_, err = config.RequireAccountID()
			if err != nil {
				return err
			}

			path := cmdutil.AccountPath("suppressions")

			query := url.Values{}
			if email != "" {
				query.Set("email", email)
			}
			if startTime != "" {
				query.Set("start_time", startTime)
			}
			if endTime != "" {
				query.Set("end_time", endTime)
			}
			if lastID != "" {
				query.Set("last_id", lastID)
			}

			var suppressions []Suppression
			if err := c.Get(context.Background(), client.BaseGeneral, path, query, &suppressions); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, suppressions, suppressionColumns)
		},
	}

	cmd.Flags().StringVar(&email, "email", "", "Filter by email address")
	cmd.Flags().StringVar(&startTime, "start-time", "", "Filter by start time")
	cmd.Flags().StringVar(&endTime, "end-time", "", "Filter by end time")
	cmd.Flags().StringVar(&lastID, "last-id", "", "Pagination cursor (id of the last record from the previous response)")

	return cmd
}
