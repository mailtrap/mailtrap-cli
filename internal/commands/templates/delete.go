package templates

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/spf13/cobra"
)

type DeleteOptions struct {
	ID int
}

func NewCmdDelete(f *cmdutil.Factory) *cobra.Command {
	opts := &DeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an email template",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			if _, err := config.RequireAccountID(); err != nil {
				return err
			}

			path := cmdutil.AccountPath("email_templates", fmt.Sprintf("%d", opts.ID))

			if err := c.Delete(context.Background(), client.BaseGeneral, path, nil); err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "Template %d deleted successfully\n", opts.ID)

			return nil
		},
	}

	cmd.Flags().IntVar(&opts.ID, "id", 0, "Template ID (required)")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}
