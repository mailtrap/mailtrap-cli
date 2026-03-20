package accounts

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdAccounts(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accounts",
		Short: "Manage accounts",
	}

	cmd.AddCommand(NewCmdList(f))

	return cmd
}
