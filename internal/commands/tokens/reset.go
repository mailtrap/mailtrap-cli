package tokens

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdReset(f *cmdutil.Factory) *cobra.Command {
	var tokenID string

	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset an API token",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", tokenID); err != nil {
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

			path := cmdutil.AccountPath("api_tokens", fmt.Sprintf("%s", tokenID), "reset")

			var resp interface{}
			if err := c.Post(context.Background(), client.BaseGeneral, path, nil, &resp); err != nil {
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

	cmd.Flags().StringVar(&tokenID, "id", "", "API token ID")

	return cmd
}
