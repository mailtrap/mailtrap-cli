package contacts

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdImportStatus(f *cmdutil.Factory) *cobra.Command {
	var importID string

	cmd := &cobra.Command{
		Use:   "import-status",
		Short: "Get status of a contact import",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", importID); err != nil {
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

			path := cmdutil.AccountPath("contacts", "imports", fmt.Sprintf("%s", importID))

			var resp interface{}
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &resp); err != nil {
				return err
			}

			respJSON, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				return err
			}

			fmt.Fprintln(f.IOStreams.Out, string(respJSON))
			return nil
		},
	}

	cmd.Flags().StringVar(&importID, "id", "", "Import ID")

	return cmd
}
