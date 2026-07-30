package emailcampaigns

import (
	"context"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var campaignID string
	attrs := &campaignAttrs{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a draft email campaign",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", campaignID); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			body := buildAttributesBody(cmd, attrs)

			var resp campaignResponse
			if err := c.Patch(context.Background(), client.BaseGeneral, campaignPath(campaignID), body, &resp); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, resp.Data, campaignColumns)
		},
	}

	cmd.Flags().StringVar(&campaignID, "id", "", "Email campaign ID (required)")
	addAttributeFlags(cmd, attrs)

	return cmd
}
