package emaillogs

import (
	"context"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type GetOptions struct {
	ID string
}

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	opts := &GetOptions{}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get an email log by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			if _, err := config.RequireAccountID(); err != nil {
				return err
			}

			path := cmdutil.AccountPath("email_logs", opts.ID)

			var result EmailLog
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &result); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			output.Print(f.IOStreams.Out, format, result, emailLogColumns)

			return nil
		},
	}

	cmd.Flags().StringVar(&opts.ID, "id", "", "Email log ID (required)")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}
