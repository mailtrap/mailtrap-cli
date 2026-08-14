package emailcampaigns

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdTerminate(f *cmdutil.Factory) *cobra.Command {
	return newLifecycleCmd(f, "terminate", "Terminate a sending email campaign")
}
