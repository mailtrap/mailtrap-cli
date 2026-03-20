package projects

import (
	"context"
	"encoding/json"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type Project struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	ShareLinks  json.RawMessage `json:"share_links"`
	Permissions json.RawMessage `json:"permissions"`
}

var projectColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "NAME", Field: "name"},
	{Header: "SHARE_LINKS", Field: "share_links"},
	{Header: "PERMISSIONS", Field: "permissions"},
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			_, err = config.RequireAccountID()
			if err != nil {
				return err
			}

			path := cmdutil.AccountPath("projects")

			var projects []Project
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &projects); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, projects, projectColumns)
		},
	}

	return cmd
}
