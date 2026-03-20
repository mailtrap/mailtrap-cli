package tokens

import (
	"context"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an API token",
		RunE: func(cmd *cobra.Command, args []string) error {
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

			path := cmdutil.AccountPath("api_tokens")

			body := map[string]interface{}{
				"api_token": map[string]string{
					"name": name,
				},
			}

			var token APIToken
			if err := c.Post(context.Background(), client.BaseGeneral, path, body, &token); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, token, tokenColumns)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "API token name")

	return cmd
}
