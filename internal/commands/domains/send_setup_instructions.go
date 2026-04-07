package domains

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdSendSetupInstructions(f *cmdutil.Factory) *cobra.Command {
	var domainID string
	var email string

	cmd := &cobra.Command{
		Use:   "send-setup-instructions",
		Short: "Send DNS setup instructions for a domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", domainID); err != nil {
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

			path := cmdutil.AccountPath("sending_domains", domainID, "send_setup_instructions")
			body := map[string]string{"email": email}

			if err := c.Post(context.Background(), client.BaseGeneral, path, body, nil); err != nil {
				return err
			}

			fmt.Fprintln(f.IOStreams.Out, "Setup instructions sent successfully.")
			return nil
		},
	}

	cmd.Flags().StringVar(&domainID, "id", "", "Domain ID (required)")
	cmd.Flags().StringVar(&email, "email", "", "Email to send instructions to (required)")

	return cmd
}
