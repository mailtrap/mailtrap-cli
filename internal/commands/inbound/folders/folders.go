package folders

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdFolders creates the `inbound folders` command group.
func NewCmdFolders(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "folders",
		Short: "Manage inbound folders",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdGet(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdUpdate(f))
	cmd.AddCommand(NewCmdDelete(f))

	return cmd
}
