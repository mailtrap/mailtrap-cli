package send

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdSend(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "Send emails via Mailtrap",
		Long:  "Send transactional or bulk emails via Mailtrap Email API.",
	}

	cmd.AddCommand(NewCmdTransactional(f))
	cmd.AddCommand(NewCmdBulk(f))
	cmd.AddCommand(NewCmdBatchTransactional(f))
	cmd.AddCommand(NewCmdBatchBulk(f))

	return cmd
}
