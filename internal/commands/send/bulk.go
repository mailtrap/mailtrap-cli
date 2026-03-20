package send

import (
	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdBulk(f *cmdutil.Factory) *cobra.Command {
	return newSendCmd(f, "bulk", "Send a bulk email", client.BaseBulk)
}
