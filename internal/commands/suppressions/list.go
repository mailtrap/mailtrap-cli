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
	Type          string `json:"type"`
	Email         string `json:"email"`
	DomainName    string `json:"domain_name"`
	SendingStream string `json:"sending_stream"`
	CreatedAt     string `json:"created_at"`
}

var suppressionColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "TYPE", Field: "type"},
	{Header: "EMAIL", Field: "email"},
	{Header: "DOMAIN", Field: "domain_name"},
	{Header: "STREAM", Field: "sending_stream"},
	{Header: "CREATED_AT", Field: "created_at"},
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var email string
	var startTime string
	var endTime string

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

	return cmd
}
