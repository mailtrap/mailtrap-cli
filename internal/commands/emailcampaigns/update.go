package emailcampaigns

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var campaignID string
	var clearContactLists, clearContactSegments bool
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
			// The API treats audience ids as the full set, so an explicit `[]`
			// clears them — inexpressible via the ids flags (pflag rejects an
			// empty slice value), hence the dedicated flags.
			if clearContactLists {
				body["contact_list_ids"] = []int64{}
			}
			if clearContactSegments {
				body["contact_segment_ids"] = []int64{}
			}
			if len(body) == 0 {
				return fmt.Errorf("at least one attribute flag is required")
			}

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
	cmd.Flags().BoolVar(&clearContactLists, "clear-contact-lists", false, "Remove all contact lists from the audience")
	cmd.Flags().BoolVar(&clearContactSegments, "clear-contact-segments", false, "Remove all contact segments from the audience")
	cmd.MarkFlagsMutuallyExclusive("contact-list-ids", "clear-contact-lists")
	cmd.MarkFlagsMutuallyExclusive("contact-segment-ids", "clear-contact-segments")

	return cmd
}
