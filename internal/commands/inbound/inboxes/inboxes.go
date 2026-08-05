package inboxes

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdInboxes creates the `inbound inboxes` command group.
func NewCmdInboxes(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inboxes",
		Short: "Manage inbound inboxes",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdGet(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdUpdate(f))
	cmd.AddCommand(NewCmdDelete(f))

	return cmd
}
