package tokens

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	var tokenID string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get an API token",
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

			path := cmdutil.AccountPath("api_tokens", fmt.Sprintf("%s", tokenID))

			var token APIToken
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &token); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, token, tokenColumns)
		},
	}

	cmd.Flags().StringVar(&tokenID, "id", "", "API token ID")

	return cmd
}
