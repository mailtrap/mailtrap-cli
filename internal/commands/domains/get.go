package domains

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
	var domainID string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a sending domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", domainID); err != nil {
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

			path := cmdutil.AccountPath("sending_domains", fmt.Sprintf("%s", domainID))

			var domain Domain
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &domain); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, domain, domainColumns)
		},
	}

	cmd.Flags().StringVar(&domainID, "id", "", "Domain ID")

	return cmd
}
