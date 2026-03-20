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

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var projectID string
	var name string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an inbox",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("project-id", projectID); err != nil {
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

			path := cmdutil.AccountPath("projects", fmt.Sprintf("%s", projectID), "inboxes")
			body := map[string]interface{}{
				"inbox": map[string]string{"name": name},
			}

			var inbox Inbox
			if err := c.Post(context.Background(), client.BaseGeneral, path, body, &inbox); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, inbox, inboxColumns)
		},
	}

	cmd.Flags().StringVar(&projectID, "project-id", "", "Project ID")
	cmd.Flags().StringVar(&name, "name", "", "Inbox name")

	return cmd
}
