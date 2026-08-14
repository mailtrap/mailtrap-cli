package emailcampaigns

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdDelete(f *cmdutil.Factory) *cobra.Command {
	var campaignID string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an email campaign",
		Long:  "Delete an email campaign. Only a campaign in the draft state can be deleted.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", campaignID); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			// The API returns 204 No Content with no body.
			if err := c.Delete(context.Background(), client.BaseGeneral, campaignPath(campaignID), nil); err != nil {
				return err
			}

			fmt.Fprintln(f.IOStreams.Out, "Email campaign deleted successfully.")
			return nil
		},
	}

	cmd.Flags().StringVar(&campaignID, "id", "", "Email campaign ID (required)")

	return cmd
}
