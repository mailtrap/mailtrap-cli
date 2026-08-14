package emailcampaigns

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdCancel(f *cmdutil.Factory) *cobra.Command {
	return newLifecycleCmd(f, "cancel", "Cancel a scheduled email campaign")
}
