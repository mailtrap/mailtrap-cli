package messages

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdMessages(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "messages",
		Short: "Manage messages",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdGet(f))
	cmd.AddCommand(NewCmdUpdate(f))
	cmd.AddCommand(NewCmdDelete(f))
	cmd.AddCommand(NewCmdForward(f))
	cmd.AddCommand(NewCmdSpamScore(f))
	cmd.AddCommand(NewCmdHTMLAnalysis(f))
	cmd.AddCommand(NewCmdHeaders(f))
	cmd.AddCommand(NewCmdHTML(f))
	cmd.AddCommand(NewCmdText(f))
	cmd.AddCommand(NewCmdSource(f))
	cmd.AddCommand(NewCmdRaw(f))
	cmd.AddCommand(NewCmdEml(f))

	return cmd
}
