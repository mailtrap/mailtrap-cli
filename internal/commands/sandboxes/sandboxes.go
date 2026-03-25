package sandboxes

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdSandboxes(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sandboxes",
		Short: "Manage sandboxes",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdGet(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdUpdate(f))
	cmd.AddCommand(NewCmdDelete(f))
	cmd.AddCommand(NewCmdClean(f))
	cmd.AddCommand(NewCmdMarkRead(f))
	cmd.AddCommand(NewCmdResetCredentials(f))
	cmd.AddCommand(NewCmdToggleEmail(f))
	cmd.AddCommand(NewCmdResetEmail(f))

	return cmd
}
