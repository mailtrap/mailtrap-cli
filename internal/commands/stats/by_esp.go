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

type ByESPOptions struct {
	StartDate  string
	EndDate    string
	DomainIDs  []string
	Streams    []string
	Categories []string
}

func NewCmdByESP(f *cmdutil.Factory) *cobra.Command {
	opts := &ByESPOptions{}

	cmd := &cobra.Command{
		Use:   "by-esp",
		Short: "Get email sending statistics grouped by email service provider",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			if _, err := config.RequireAccountID(); err != nil {
				return err
			}

			path := cmdutil.AccountPath("stats", "email_service_providers")

			params := url.Values{}
			params.Set("start_date", opts.StartDate)
			params.Set("end_date", opts.EndDate)
			for _, d := range opts.DomainIDs {
				params.Add("domain_ids[]", d)
			}
			for _, s := range opts.Streams {
				params.Add("streams[]", s)
			}
			for _, cat := range opts.Categories {
				params.Add("categories[]", cat)
			}

			fullPath := fmt.Sprintf("%s?%s", path, params.Encode())

			var result []Stats
			if err := c.Get(context.Background(), client.BaseGeneral, fullPath, nil, &result); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			output.Print(f.IOStreams.Out, format, result, statsColumns)

			return nil
		},
	}

	cmd.Flags().StringVar(&opts.StartDate, "start-date", "", "Start date (required)")
	cmd.Flags().StringVar(&opts.EndDate, "end-date", "", "End date (required)")
	cmd.Flags().StringSliceVar(&opts.DomainIDs, "domain-ids", nil, "Filter by domain IDs")
	cmd.Flags().StringSliceVar(&opts.Streams, "streams", nil, "Filter by streams")
	cmd.Flags().StringSliceVar(&opts.Categories, "categories", nil, "Filter by categories")

	_ = cmd.MarkFlagRequired("start-date")
	_ = cmd.MarkFlagRequired("end-date")

	return cmd
}
