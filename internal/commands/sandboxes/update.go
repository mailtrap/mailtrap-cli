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

func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var inboxID string
	var name string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a sandbox",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", inboxID); err != nil {
				return err
			}
			if err := cmdutil.RequireFlag("name", name); err != nil {
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

			path := cmdutil.AccountPath("inboxes", fmt.Sprintf("%s", inboxID))
			body := map[string]interface{}{
				"inbox": map[string]string{"name": name},
			}

			var inbox Inbox
			if err := c.Patch(context.Background(), client.BaseGeneral, path, body, &inbox); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, inbox, inboxColumns)
		},
	}

	cmd.Flags().StringVar(&inboxID, "id", "", "Sandbox ID")
	cmd.Flags().StringVar(&name, "name", "", "Sandbox name")

	return cmd
}
