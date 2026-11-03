package sandboxes

import (
	"context"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var sandboxID string
	var name string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a sandbox",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", sandboxID); err != nil {
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

			path := cmdutil.AccountPath("inboxes", sandboxID)
			body := map[string]interface{}{
				"inbox": map[string]string{"name": name},
			}

			var sandbox Sandbox
			if err := c.Patch(context.Background(), client.BaseGeneral, path, body, &sandbox); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, sandbox, sandboxColumns)
		},
	}

	cmd.Flags().StringVar(&sandboxID, "id", "", "Sandbox ID")
	cmd.Flags().StringVar(&name, "name", "", "Sandbox name")

	return cmd
}
