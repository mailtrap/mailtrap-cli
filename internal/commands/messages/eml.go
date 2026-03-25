package messages

import (
	"context"
	"fmt"
	"os"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdEml(f *cmdutil.Factory) *cobra.Command {
	var inboxID string
	var messageID string
	var outputFile string

	cmd := &cobra.Command{
		Use:   "eml",
		Short: "Download EML file of a message",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("inbox-id", inboxID); err != nil {
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

			path := cmdutil.AccountPath("inboxes", fmt.Sprintf("%s", inboxID), "messages", fmt.Sprintf("%s", messageID), "body.eml")

			data, err := c.GetRaw(context.Background(), client.BaseGeneral, path, nil)
			if err != nil {
				return err
			}

			if outputFile != "" {
				if err := os.WriteFile(outputFile, data, 0644); err != nil {
					return fmt.Errorf("write file: %w", err)
				}
				fmt.Fprintf(f.IOStreams.Out, "EML saved to %s\n", outputFile)
				return nil
			}

			return output.PrintRaw(f.IOStreams.Out, data)
		},
	}

	cmd.Flags().StringVar(&inboxID, "inbox-id", "", "Inbox ID")
	cmd.Flags().StringVar(&messageID, "id", "", "Message ID")
	cmd.Flags().StringVar(&outputFile, "output-file", "", "Output file path (default: stdout)")

	return cmd
}
