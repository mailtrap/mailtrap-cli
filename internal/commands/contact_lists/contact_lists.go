package contact_lists

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdContactLists(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contact-lists",
		Short: "Manage contact lists",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdGet(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdUpdate(f))
	cmd.AddCommand(NewCmdDelete(f))

	return cmd
}
