package forwardrules

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdForwardRules creates the `inbound forward-rules` command group.
func NewCmdForwardRules(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "forward-rules",
		Short: "Manage inbound forward rules",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdGet(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdUpdate(f))
	cmd.AddCommand(NewCmdDelete(f))

	return cmd
}
