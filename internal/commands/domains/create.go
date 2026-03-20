package domains

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
		Short: "Create a sending domain",
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

			path := cmdutil.AccountPath("sending_domains")
			body := map[string]string{"domain": name}

			var domain Domain
			if err := c.Post(context.Background(), client.BaseGeneral, path, body, &domain); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, domain, domainColumns)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Domain name")

	return cmd
}
