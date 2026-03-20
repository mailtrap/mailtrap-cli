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

type Stats struct {
	DeliveryCount int     `json:"delivery_count"`
	DeliveryRate  float64 `json:"delivery_rate"`
	BounceCount   int     `json:"bounce_count"`
	BounceRate    float64 `json:"bounce_rate"`
	OpenCount     int     `json:"open_count"`
	OpenRate      float64 `json:"open_rate"`
	ClickCount    int     `json:"click_count"`
	ClickRate     float64 `json:"click_rate"`
	SpamCount     int     `json:"spam_count"`
	SpamRate      float64 `json:"spam_rate"`
}

var statsColumns = []output.Column{
	{Header: "DELIVERY COUNT", Field: "delivery_count"},
	{Header: "DELIVERY RATE", Field: "delivery_rate"},
	{Header: "BOUNCE COUNT", Field: "bounce_count"},
	{Header: "BOUNCE RATE", Field: "bounce_rate"},
	{Header: "OPEN COUNT", Field: "open_count"},
	{Header: "OPEN RATE", Field: "open_rate"},
	{Header: "CLICK COUNT", Field: "click_count"},
	{Header: "CLICK RATE", Field: "click_rate"},
	{Header: "SPAM COUNT", Field: "spam_count"},
	{Header: "SPAM RATE", Field: "spam_rate"},
}

type GetOptions struct {
	StartDate string
	EndDate   string
	DomainIDs []string
	Streams   []string
	Categories []string
}

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	opts := &GetOptions{}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get aggregated email sending statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			if _, err := config.RequireAccountID(); err != nil {
				return err
			}

			path := cmdutil.AccountPath("stats")

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

			var result Stats
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
