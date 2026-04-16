package stats

import (
	"context"
	"encoding/json"
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

// Grouped response types — the API wraps stats under a "stats" key alongside the grouping field.
// MarshalJSON on each type flattens the nested stats for table/text output.

type DomainStats struct {
	SendingDomainID int   `json:"sending_domain_id"`
	Stats           Stats `json:"stats"`
}

func (d DomainStats) MarshalJSON() ([]byte, error) {
	type flat struct {
		SendingDomainID int     `json:"sending_domain_id"`
		DeliveryCount   int     `json:"delivery_count"`
		DeliveryRate    float64 `json:"delivery_rate"`
		BounceCount     int     `json:"bounce_count"`
		BounceRate      float64 `json:"bounce_rate"`
		OpenCount       int     `json:"open_count"`
		OpenRate        float64 `json:"open_rate"`
		ClickCount      int     `json:"click_count"`
		ClickRate       float64 `json:"click_rate"`
		SpamCount       int     `json:"spam_count"`
		SpamRate        float64 `json:"spam_rate"`
	}
	return json.Marshal(flat{d.SendingDomainID, d.Stats.DeliveryCount, d.Stats.DeliveryRate, d.Stats.BounceCount, d.Stats.BounceRate, d.Stats.OpenCount, d.Stats.OpenRate, d.Stats.ClickCount, d.Stats.ClickRate, d.Stats.SpamCount, d.Stats.SpamRate})
}

var domainStatsColumns = []output.Column{
	{Header: "DOMAIN ID", Field: "sending_domain_id"},
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

type CategoryStats struct {
	Category string `json:"category"`
	Stats    Stats  `json:"stats"`
}

func (c CategoryStats) MarshalJSON() ([]byte, error) {
	type flat struct {
		Category      string  `json:"category"`
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
	return json.Marshal(flat{c.Category, c.Stats.DeliveryCount, c.Stats.DeliveryRate, c.Stats.BounceCount, c.Stats.BounceRate, c.Stats.OpenCount, c.Stats.OpenRate, c.Stats.ClickCount, c.Stats.ClickRate, c.Stats.SpamCount, c.Stats.SpamRate})
}

var categoryStatsColumns = []output.Column{
	{Header: "CATEGORY", Field: "category"},
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

type ESPStats struct {
	EmailServiceProvider string `json:"email_service_provider"`
	Stats                Stats  `json:"stats"`
}

func (e ESPStats) MarshalJSON() ([]byte, error) {
	type flat struct {
		EmailServiceProvider string  `json:"email_service_provider"`
		DeliveryCount        int     `json:"delivery_count"`
		DeliveryRate         float64 `json:"delivery_rate"`
		BounceCount          int     `json:"bounce_count"`
		BounceRate           float64 `json:"bounce_rate"`
		OpenCount            int     `json:"open_count"`
		OpenRate             float64 `json:"open_rate"`
		ClickCount           int     `json:"click_count"`
		ClickRate            float64 `json:"click_rate"`
		SpamCount            int     `json:"spam_count"`
		SpamRate             float64 `json:"spam_rate"`
	}
	return json.Marshal(flat{e.EmailServiceProvider, e.Stats.DeliveryCount, e.Stats.DeliveryRate, e.Stats.BounceCount, e.Stats.BounceRate, e.Stats.OpenCount, e.Stats.OpenRate, e.Stats.ClickCount, e.Stats.ClickRate, e.Stats.SpamCount, e.Stats.SpamRate})
}

var espStatsColumns = []output.Column{
	{Header: "ESP", Field: "email_service_provider"},
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

type DateStats struct {
	Date  string `json:"date"`
	Stats Stats  `json:"stats"`
}

func (d DateStats) MarshalJSON() ([]byte, error) {
	type flat struct {
		Date          string  `json:"date"`
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
	return json.Marshal(flat{d.Date, d.Stats.DeliveryCount, d.Stats.DeliveryRate, d.Stats.BounceCount, d.Stats.BounceRate, d.Stats.OpenCount, d.Stats.OpenRate, d.Stats.ClickCount, d.Stats.ClickRate, d.Stats.SpamCount, d.Stats.SpamRate})
}

var dateStatsColumns = []output.Column{
	{Header: "DATE", Field: "date"},
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
				params.Add("sending_domain_ids[]", d)
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
