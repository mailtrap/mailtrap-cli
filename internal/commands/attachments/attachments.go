package attachments

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdAttachments(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attachments",
		Short: "Manage message attachments",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdGet(f))

	return cmd
}
