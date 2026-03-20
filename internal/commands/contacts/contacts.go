package contacts

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdContacts(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contacts",
		Short: "Manage contacts",
	}

	cmd.AddCommand(NewCmdGet(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdUpdate(f))
	cmd.AddCommand(NewCmdDelete(f))
	cmd.AddCommand(NewCmdImport(f))
	cmd.AddCommand(NewCmdImportStatus(f))
	cmd.AddCommand(NewCmdExport(f))
	cmd.AddCommand(NewCmdExportStatus(f))
	cmd.AddCommand(NewCmdCreateEvent(f))

	return cmd
}
