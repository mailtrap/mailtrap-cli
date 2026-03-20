package templates

import (
	"context"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type CreateOptions struct {
	Name     string
	Subject  string
	BodyHTML string
	BodyText string
	Category string
}

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	opts := &CreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new email template",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			if _, err := config.RequireAccountID(); err != nil {
				return err
			}

			path := cmdutil.AccountPath("email_templates")

			body := map[string]interface{}{
				"email_template": map[string]interface{}{
					"name":      opts.Name,
					"subject":   opts.Subject,
					"body_html": opts.BodyHTML,
					"body_text": opts.BodyText,
					"category":  opts.Category,
				},
			}

			var result Template
			if err := c.Post(context.Background(), client.BaseGeneral, path, body, &result); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			output.Print(f.IOStreams.Out, format, result, templateColumns)

			return nil
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "Template name (required)")
	cmd.Flags().StringVar(&opts.Subject, "subject", "", "Template subject (required)")
	cmd.Flags().StringVar(&opts.BodyHTML, "body-html", "", "HTML body content")
	cmd.Flags().StringVar(&opts.BodyText, "body-text", "", "Plain text body content")
	cmd.Flags().StringVar(&opts.Category, "category", "", "Template category")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("subject")

	return cmd
}
