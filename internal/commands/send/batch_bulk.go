package send

import (
	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdBatchBulk(f *cmdutil.Factory) *cobra.Command {
	return newBatchCmd(f, "batch-bulk", "Send a batch of bulk emails", client.BaseBulk)
}
