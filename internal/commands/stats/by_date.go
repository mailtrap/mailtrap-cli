package stats

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

type ByDateOptions struct {
	StartDate  string
	EndDate    string
	DomainIDs  []string
	Streams    []string
	Categories []string
	ESPs       []string
}

func NewCmdByDate(f *cmdutil.Factory) *cobra.Command {
	opts := &ByDateOptions{}

	cmd := &cobra.Command{
		Use:   "by-date",
		Short: "Get email sending statistics grouped by date",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			if _, err := config.RequireAccountID(); err != nil {
				return err
			}

			path := cmdutil.AccountPath("stats", "date")

			params := url.Values{}
			params.Set("start_date", opts.StartDate)
			params.Set("end_date", opts.EndDate)
			for _, d := range opts.DomainIDs {
				params.Add("sending_domain_ids[]", d)
			}
			for _, s := range opts.Streams {
				params.Add("sending_streams[]", s)
			}
			for _, cat := range opts.Categories {
				params.Add("categories[]", cat)
			}
			for _, esp := range opts.ESPs {
				params.Add("email_service_providers[]", esp)
			}

			fullPath := fmt.Sprintf("%s?%s", path, params.Encode())

			var result []DateStats
			if err := c.Get(context.Background(), client.BaseGeneral, fullPath, nil, &result); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			output.Print(f.IOStreams.Out, format, result, dateStatsColumns)

			return nil
		},
	}

	cmd.Flags().StringVar(&opts.StartDate, "start-date", "", "Start date (required)")
	cmd.Flags().StringVar(&opts.EndDate, "end-date", "", "End date (required)")
	cmd.Flags().StringSliceVar(&opts.DomainIDs, "domain-ids", nil, "Filter by domain IDs")
	cmd.Flags().StringSliceVar(&opts.Streams, "streams", nil, "Filter by sending streams (e.g. transactional, bulk)")
	cmd.Flags().StringSliceVar(&opts.Categories, "categories", nil, "Filter by categories")
	cmd.Flags().StringSliceVar(&opts.ESPs, "esps", nil, "Filter by email service providers (e.g. Google, Yahoo)")

	_ = cmd.MarkFlagRequired("start-date")
	_ = cmd.MarkFlagRequired("end-date")

	return cmd
}
