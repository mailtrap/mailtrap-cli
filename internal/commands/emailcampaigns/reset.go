package emailcampaigns

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdReset(f *cmdutil.Factory) *cobra.Command {
	return newLifecycleCmd(f, "reset", "Reset a scheduled email campaign to draft")
}
