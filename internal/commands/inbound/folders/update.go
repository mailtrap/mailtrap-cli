package folders

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var (
		folderID string
		name     string
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update an inbound folder",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", folderID); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/api/inbound/folders/%s", folderID)

			body := map[string]interface{}{}
			if cmd.Flags().Changed("name") {
				body["name"] = name
			}

			var resp InboundFolder
			if err := c.Patch(context.Background(), client.BaseGeneral, path, body, &resp); err != nil {
				return err
			}

			return output.Print(f.IOStreams.Out, cmdutil.GetOutputFormat(), resp, folderColumns)
		},
	}

	cmd.Flags().StringVar(&folderID, "id", "", "Folder ID (required)")
	cmd.Flags().StringVar(&name, "name", "", "Folder name")

	return cmd
}
