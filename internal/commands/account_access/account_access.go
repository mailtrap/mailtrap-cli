package account_access

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdAccountAccess(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account-access",
		Short: "Manage account access",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdRemove(f))

	return cmd
}
