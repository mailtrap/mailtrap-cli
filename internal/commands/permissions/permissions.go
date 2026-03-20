package permissions

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdPermissions(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "permissions",
		Short: "Manage permissions",
	}

	cmd.AddCommand(NewCmdBulkUpdate(f))
	cmd.AddCommand(NewCmdResources(f))

	return cmd
}
