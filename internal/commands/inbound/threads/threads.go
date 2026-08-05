package threads

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdThreads creates the `inbound threads` command group.
func NewCmdThreads(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "threads",
		Short: "Manage inbound conversation threads",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdGet(f))
	cmd.AddCommand(NewCmdDelete(f))

	return cmd
}
