package emaillogs

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdEmailLogs(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "email-logs",
		Short: "Manage email logs",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdGet(f))

	return cmd
}
