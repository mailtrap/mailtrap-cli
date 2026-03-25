package messages

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdRaw(f *cmdutil.Factory) *cobra.Command {
	var sandboxID string
	var messageID string

	cmd := &cobra.Command{
		Use:   "raw",
		Short: "Get raw email of a message",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("sandbox-id", sandboxID); err != nil {
				return err
			}
			if err := cmdutil.RequireFlag("id", messageID); err != nil {
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

			path := cmdutil.AccountPath("inboxes", fmt.Sprintf("%s", sandboxID), "messages", fmt.Sprintf("%s", messageID), "body.raw")

			data, err := c.GetRaw(context.Background(), client.BaseGeneral, path, nil)
			if err != nil {
				return err
			}

			return output.PrintRaw(f.IOStreams.Out, data)
		},
	}

	cmd.Flags().StringVar(&sandboxID, "sandbox-id", "", "Sandbox ID")
	cmd.Flags().StringVar(&messageID, "id", "", "Message ID")

	return cmd
}
