package templates

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type GetOptions struct {
	ID int
}

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	opts := &GetOptions{}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get an email template by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			if _, err := config.RequireAccountID(); err != nil {
				return err
			}

			path := cmdutil.AccountPath("email_templates", fmt.Sprintf("%d", opts.ID))

			var result Template
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &result); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			output.Print(f.IOStreams.Out, format, result, templateColumns)

			return nil
		},
	}

	cmd.Flags().IntVar(&opts.ID, "id", 0, "Template ID (required)")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}
