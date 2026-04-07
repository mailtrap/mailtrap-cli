package domains

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdDomains(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domains",
		Short: "Manage sending domains",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdGet(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdDelete(f))
	cmd.AddCommand(NewCmdSendSetupInstructions(f))

	return cmd
}
