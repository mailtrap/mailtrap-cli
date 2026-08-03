package inbound

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/commands/inbound/folders"
	"github.com/mailtrap/mailtrap-cli/internal/commands/inbound/inboxes"
	"github.com/mailtrap/mailtrap-cli/internal/commands/inbound/messages"
	"github.com/mailtrap/mailtrap-cli/internal/commands/inbound/threads"
	"github.com/spf13/cobra"
)

// NewCmdInbound creates the `inbound` command group.
func NewCmdInbound(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inbound",
		Short: "Manage inbound email folders, inboxes, messages, and threads",
	}

	cmd.AddCommand(folders.NewCmdFolders(f))
	cmd.AddCommand(inboxes.NewCmdInboxes(f))
	cmd.AddCommand(messages.NewCmdMessages(f))
	cmd.AddCommand(threads.NewCmdThreads(f))

	return cmd
}
