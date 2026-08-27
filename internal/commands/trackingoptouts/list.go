package trackingoptouts

import (
	"context"
	"fmt"
	"net/url"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type trackingOptOutsListResponse struct {
	Data   []TrackingOptOut `json:"data"`
	LastID string           `json:"last_id"`
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

			var resp trackingOptOutsListResponse
			if err := c.Get(context.Background(), client.BaseGeneral, trackingOptOutsPath, query, &resp); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			if err := output.Print(f.IOStreams.Out, format, resp.Data, trackingOptOutColumns); err != nil {
				return err
			}
			if format != output.FormatJSON && resp.LastID != "" {
				fmt.Fprintf(f.IOStreams.Out, "\nNext page: --last-id %s\n", resp.LastID)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&email, "email", "", "Filter by email address")
	cmd.Flags().StringVar(&startTime, "start-time", "", "Filter by start time")
	cmd.Flags().StringVar(&endTime, "end-time", "", "Filter by end time")
	cmd.Flags().StringVar(&lastID, "last-id", "", "Pagination cursor (last_id from previous response)")

	return cmd
}
