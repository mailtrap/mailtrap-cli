package emailcampaigns

import (
	"context"
	"net/url"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type CampaignStats struct {
	DeliveryCount       int     `json:"delivery_count"`
	OpenCount           int     `json:"open_count"`
	ClickCount          int     `json:"click_count"`
	BounceCount         int     `json:"bounce_count"`
	UnsubscriptionCount int     `json:"unsubscription_count"`
	SentCount           int     `json:"sent_count"`
	SpamCount           int     `json:"spam_count"`
	DeliveryRate        float64 `json:"delivery_rate"`
	OpenRate            float64 `json:"open_rate"`
	ClickRate           float64 `json:"click_rate"`
	BounceRate          float64 `json:"bounce_rate"`
	SpamRate            float64 `json:"spam_rate"`
	UnsubscriptionRate  float64 `json:"unsubscription_rate"`
}

type campaignStatsResponse struct {
	Data CampaignStats `json:"data"`
}

var campaignStatsColumns = []output.Column{
	{Header: "SENT", Field: "sent_count"},
	{Header: "DELIVERED", Field: "delivery_count"},
	{Header: "OPENS", Field: "open_count"},
	{Header: "CLICKS", Field: "click_count"},
	{Header: "BOUNCES", Field: "bounce_count"},
	{Header: "SPAM", Field: "spam_count"},
	{Header: "UNSUBSCRIBES", Field: "unsubscription_count"},
	{Header: "DELIVERY RATE", Field: "delivery_rate"},
	{Header: "OPEN RATE", Field: "open_rate"},
	{Header: "CLICK RATE", Field: "click_rate"},
	{Header: "BOUNCE RATE", Field: "bounce_rate"},
	{Header: "SPAM RATE", Field: "spam_rate"},
	{Header: "UNSUBSCRIPTION RATE", Field: "unsubscription_rate"},
}

func NewCmdStats(f *cmdutil.Factory) *cobra.Command {
	var (
		campaignID string
		startDate  string
		endDate    string
	)

	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Get email campaign statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", campaignID); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			query := url.Values{}
			if startDate != "" {
				query.Set("start_date", startDate)
			}
			if endDate != "" {
				query.Set("end_date", endDate)
			}

			var resp campaignStatsResponse
			if err := c.Get(context.Background(), client.BaseGeneral, campaignPath(campaignID, "stats"), query, &resp); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, resp.Data, campaignStatsColumns)
		},
	}

	cmd.Flags().StringVar(&campaignID, "id", "", "Email campaign ID (required)")
	cmd.Flags().StringVar(&startDate, "start-date", "", "Start of the aggregation window (YYYY-MM-DD)")
	cmd.Flags().StringVar(&endDate, "end-date", "", "End of the aggregation window (YYYY-MM-DD)")

	return cmd
}
