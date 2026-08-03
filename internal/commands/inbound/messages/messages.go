package messages

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdMessages creates the `inbound messages` command group.
func NewCmdMessages(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "messages",
		Short: "Manage inbound messages",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdGet(f))
	cmd.AddCommand(NewCmdDelete(f))
	cmd.AddCommand(NewCmdReply(f))
	cmd.AddCommand(NewCmdReplyAll(f))
	cmd.AddCommand(NewCmdForward(f))

	return cmd
}
