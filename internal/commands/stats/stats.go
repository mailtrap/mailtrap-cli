package stats

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdStats(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Manage email sending statistics",
	}

	cmd.AddCommand(NewCmdGet(f))
	cmd.AddCommand(NewCmdByDomain(f))
	cmd.AddCommand(NewCmdByCategory(f))
	cmd.AddCommand(NewCmdByESP(f))
	cmd.AddCommand(NewCmdByDate(f))

	return cmd
}
