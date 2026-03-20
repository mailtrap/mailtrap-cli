package contact_fields

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdContactFields(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contact-fields",
		Short: "Manage contact fields",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdGet(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdUpdate(f))
	cmd.AddCommand(NewCmdDelete(f))

	return cmd
}
