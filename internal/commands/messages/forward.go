package messages

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdForward(f *cmdutil.Factory) *cobra.Command {
	var sandboxID string
	var messageID string
	var email string

	cmd := &cobra.Command{
		Use:   "forward",
		Short: "Forward a message to an email address",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("sandbox-id", sandboxID); err != nil {
				return err
			}
			if err := cmdutil.RequireFlag("id", messageID); err != nil {
				return err
			}
			if err := cmdutil.RequireFlag("email", email); err != nil {
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

			path := cmdutil.AccountPath("inboxes", fmt.Sprintf("%s", sandboxID), "messages", fmt.Sprintf("%s", messageID), "forward")
			body := map[string]string{"email": email}

			if err := c.Post(context.Background(), client.BaseGeneral, path, body, nil); err != nil {
				return err
			}

			fmt.Fprintln(f.IOStreams.Out, "Message forwarded successfully.")
			return nil
		},
	}

	cmd.Flags().StringVar(&sandboxID, "sandbox-id", "", "Sandbox ID")
	cmd.Flags().StringVar(&messageID, "id", "", "Message ID")
	cmd.Flags().StringVar(&email, "email", "", "Email address to forward to")

	return cmd
}
