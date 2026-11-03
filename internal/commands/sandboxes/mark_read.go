package sandboxes

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdMarkRead(f *cmdutil.Factory) *cobra.Command {
	var sandboxID string

	cmd := &cobra.Command{
		Use:   "mark-read",
		Short: "Mark all messages in a sandbox as read",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", sandboxID); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			_, err = config.RequireAccountID()
			if err != nil {
				return err
			}

			path := cmdutil.AccountPath("inboxes", fmt.Sprintf("%s", sandboxID), "all_read")

			var sandbox Sandbox
			if err := c.Patch(context.Background(), client.BaseGeneral, path, nil, &sandbox); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, sandbox, sandboxColumns)
		},
	}

	cmd.Flags().StringVar(&sandboxID, "id", "", "Sandbox ID")

	return cmd
}
