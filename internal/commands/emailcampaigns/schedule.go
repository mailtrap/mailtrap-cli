package emailcampaigns

import (
	"context"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdSchedule(f *cmdutil.Factory) *cobra.Command {
	var (
		campaignID string
		datetime   string
	)

	cmd := &cobra.Command{
		Use:   "schedule",
		Short: "Schedule a draft email campaign",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", campaignID); err != nil {
				return err
			}
			if err := cmdutil.RequireFlag("datetime", datetime); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			body := map[string]interface{}{"datetime": datetime}

			var resp campaignResponse
			if err := c.Post(context.Background(), client.BaseGeneral, campaignPath(campaignID, "schedule"), body, &resp); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, resp.Data, campaignColumns)
		},
	}

	cmd.Flags().StringVar(&campaignID, "id", "", "Email campaign ID (required)")
	cmd.Flags().StringVar(&datetime, "datetime", "", "When to send the campaign, ISO 8601 (required)")

	return cmd
}
