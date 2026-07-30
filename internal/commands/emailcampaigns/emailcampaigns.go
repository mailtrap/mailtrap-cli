package emailcampaigns

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// Email campaigns are token-scoped: the API path is /api/email_campaigns with no
// account_id segment, so commands in this package do not call config.RequireAccountID().
const basePath = "/api/email_campaigns"

func campaignPath(segments ...string) string {
	path := basePath
	for _, s := range segments {
		path += "/" + s
	}
	return path
}

func NewCmdEmailCampaigns(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "email-campaigns",
		Short: "Manage email campaigns",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdGet(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdUpdate(f))
	cmd.AddCommand(NewCmdDelete(f))
	cmd.AddCommand(NewCmdStart(f))
	cmd.AddCommand(NewCmdSchedule(f))
	cmd.AddCommand(NewCmdCancel(f))
	cmd.AddCommand(NewCmdTerminate(f))
	cmd.AddCommand(NewCmdReset(f))
	cmd.AddCommand(NewCmdStats(f))

	return cmd
}
