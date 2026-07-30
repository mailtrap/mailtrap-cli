package emailcampaigns

import (
	"context"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

// newLifecycleCmd builds a body-less lifecycle action subcommand
// (POST /api/email_campaigns/{id}/{action}) shared by start, cancel, terminate, and reset.
func newLifecycleCmd(f *cmdutil.Factory, action, short string) *cobra.Command {
	var campaignID string

	cmd := &cobra.Command{
		Use:   action,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", campaignID); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			var resp campaignResponse
			if err := c.Post(context.Background(), client.BaseGeneral, campaignPath(campaignID, action), nil, &resp); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, resp.Data, campaignColumns)
		},
	}

	cmd.Flags().StringVar(&campaignID, "id", "", "Email campaign ID (required)")

	return cmd
}

func NewCmdStart(f *cmdutil.Factory) *cobra.Command {
	return newLifecycleCmd(f, "start", "Start a draft email campaign")
}
