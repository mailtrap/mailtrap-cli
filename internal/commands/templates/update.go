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

type UpdateOptions struct {
	ID       int
	Name     string
	Subject  string
	BodyHTML string
	BodyText string
	Category string
}

func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	opts := &UpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update an existing email template",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			if _, err := config.RequireAccountID(); err != nil {
				return err
			}

			path := cmdutil.AccountPath("email_templates", fmt.Sprintf("%d", opts.ID))

			templateFields := map[string]interface{}{}
			if cmd.Flags().Changed("name") {
				templateFields["name"] = opts.Name
			}
			if cmd.Flags().Changed("subject") {
				templateFields["subject"] = opts.Subject
			}
			if cmd.Flags().Changed("body-html") {
				templateFields["body_html"] = opts.BodyHTML
			}
			if cmd.Flags().Changed("body-text") {
				templateFields["body_text"] = opts.BodyText
			}
			if cmd.Flags().Changed("category") {
				templateFields["category"] = opts.Category
			}

			body := map[string]interface{}{
				"email_template": templateFields,
			}

			var result Template
			if err := c.Patch(context.Background(), client.BaseGeneral, path, body, &result); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			output.Print(f.IOStreams.Out, format, result, templateColumns)

			return nil
		},
	}

	cmd.Flags().IntVar(&opts.ID, "id", 0, "Template ID (required)")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Template name")
	cmd.Flags().StringVar(&opts.Subject, "subject", "", "Template subject")
	cmd.Flags().StringVar(&opts.BodyHTML, "body-html", "", "HTML body content")
	cmd.Flags().StringVar(&opts.BodyText, "body-text", "", "Plain text body content")
	cmd.Flags().StringVar(&opts.Category, "category", "", "Template category")

	_ = cmd.MarkFlagRequired("id")

	return cmd
}
