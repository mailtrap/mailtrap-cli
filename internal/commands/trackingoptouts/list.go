package trackingoptouts

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

var trackingOptOutsPage = output.Page{
	Items:      "data",
	Cursor:     []string{"last_id"},
	CursorFlag: "last-id",
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var (
		email     string
		startTime string
		endTime   string
		lastID    string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tracking opt-outs",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

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

			var resp json.RawMessage
			if err := c.Get(context.Background(), client.BaseGeneral, trackingOptOutsPath, query, &resp); err != nil {
				return err
			}

			return output.PrintPage(f.IOStreams.Out, cmdutil.GetOutputFormat(), resp, trackingOptOutsPage, trackingOptOutColumns)
		},
	}

	cmd.Flags().StringVar(&email, "email", "", "Filter by email address")
	cmd.Flags().StringVar(&startTime, "start-time", "", "Filter by start time")
	cmd.Flags().StringVar(&endTime, "end-time", "", "Filter by end time")
	cmd.Flags().StringVar(&lastID, "last-id", "", "Pagination cursor (last_id from previous response)")

	return cmd
}
