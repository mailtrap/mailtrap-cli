package contacts

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdImport(f *cmdutil.Factory) *cobra.Command {
	var file string
	var listID int

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import contacts from a JSON file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("file", file); err != nil {
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

			data, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("read file %s: %w", file, err)
			}

			var body interface{}
			if err := json.Unmarshal(data, &body); err != nil {
				return fmt.Errorf("parse JSON from %s: %w", file, err)
			}

			path := cmdutil.AccountPath("contacts", "imports")

			var resp interface{}
			if err := c.Post(context.Background(), client.BaseGeneral, path, body, &resp); err != nil {
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

	cmd.Flags().StringVar(&file, "file", "", "Path to JSON file with import data")
	cmd.Flags().IntVar(&listID, "list-id", 0, "List ID to import contacts into")

	return cmd
}
