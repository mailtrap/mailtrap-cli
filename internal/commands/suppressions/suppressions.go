package suppressions

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdSuppressions(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "suppressions",
		Short: "Manage suppressions",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdDelete(f))

	return cmd
}
