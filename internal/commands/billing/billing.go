package billing

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdBilling(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "billing",
		Short: "Manage billing",
	}

	cmd.AddCommand(NewCmdUsage(f))

	return cmd
}
