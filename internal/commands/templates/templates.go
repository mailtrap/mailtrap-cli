package templates

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdTemplates(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "templates",
		Short: "Manage email templates",
		Long:  "Manage email templates.\n\nUses the experimental /api/templates endpoints; their request and response shapes may change before general availability.",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdGet(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdUpdate(f))
	cmd.AddCommand(NewCmdDelete(f))

	return cmd
}
