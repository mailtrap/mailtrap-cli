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

			path := cmdutil.AccountPath("templates", fmt.Sprintf("%d", opts.ID))

			body := map[string]interface{}{}
			if cmd.Flags().Changed("name") {
				body["name"] = opts.Name
			}
			if cmd.Flags().Changed("subject") {
				body["subject"] = opts.Subject
			}
			if cmd.Flags().Changed("body-html") {
				body["body_html"] = opts.BodyHTML
			}
			if cmd.Flags().Changed("body-text") {
				body["body_text"] = opts.BodyText
			}
			if cmd.Flags().Changed("category") {
				body["category"] = opts.Category
			}
			if len(body) == 0 {
				return fmt.Errorf("at least one attribute flag is required")
			}

			var resp templateResponse
			if err := c.Patch(context.Background(), client.BaseGeneral, path, body, &resp); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, resp.Data, templateColumns)
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
