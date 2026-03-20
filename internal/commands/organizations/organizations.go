package organizations

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdOrganizations(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "organizations",
		Short: "Manage organizations",
	}

	cmd.AddCommand(NewCmdListSubAccounts(f))
	cmd.AddCommand(NewCmdCreateSubAccount(f))

	return cmd
}
