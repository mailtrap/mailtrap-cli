package projects

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
	var projectID string
	var name string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a project",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", projectID); err != nil {
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

			path := cmdutil.AccountPath("projects", fmt.Sprintf("%s", projectID))
			body := map[string]interface{}{
				"project": map[string]string{"name": name},
			}

			var project Project
			if err := c.Patch(context.Background(), client.BaseGeneral, path, body, &project); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, project, projectColumns)
		},
	}

	cmd.Flags().StringVar(&projectID, "id", "", "Project ID")
	cmd.Flags().StringVar(&name, "name", "", "Project name")

	return cmd
}
