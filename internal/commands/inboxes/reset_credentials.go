package inboxes

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdResetCredentials(f *cmdutil.Factory) *cobra.Command {
	var inboxID string

	cmd := &cobra.Command{
		Use:   "reset-credentials",
		Short: "Reset SMTP credentials of an inbox",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", inboxID); err != nil {
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

			path := cmdutil.AccountPath("inboxes", fmt.Sprintf("%s", inboxID), "reset_credentials")

			var inbox Inbox
			if err := c.Patch(context.Background(), client.BaseGeneral, path, nil, &inbox); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, inbox, inboxColumns)
		},
	}

	cmd.Flags().StringVar(&inboxID, "id", "", "Inbox ID")

	return cmd
}
