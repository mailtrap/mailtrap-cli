package templates

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

// templateResponse keeps the template as raw JSON so --output json prints every field.
type templateResponse struct {
	Data json.RawMessage `json:"data"`
}

var templateColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "UUID", Field: "uuid"},
	{Header: "NAME", Field: "name"},
	{Header: "SUBJECT", Field: "subject"},
	{Header: "CATEGORY", Field: "category"},
	{Header: "CREATED AT", Field: "created_at"},
}

var templatesPage = output.Page{
	Items:      "data",
	Cursor:     []string{"pagination", "next_token"},
	CursorFlag: "token",
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var (
		perPage int
		token   int
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List email templates, one page at a time",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			if _, err := config.RequireAccountID(); err != nil {
				return err
			}

			query := url.Values{}
			page := templatesPage
			if cmd.Flags().Changed("per-page") {
				query.Set("per_page", fmt.Sprintf("%d", perPage))
				// Without per_page the API reads the next token at its default page size.
				page.NextArgs = fmt.Sprintf("--per-page %d", perPage)
			}
			if cmd.Flags().Changed("token") {
				query.Set("token", fmt.Sprintf("%d", token))
			}

			var resp json.RawMessage
			if err := c.Get(context.Background(), client.BaseGeneral, cmdutil.AccountPath("templates"), query, &resp); err != nil {
				return err
			}

			return output.PrintPage(f.IOStreams.Out, cmdutil.GetOutputFormat(), resp, page, templateColumns)
		},
	}

	cmd.Flags().IntVar(&perPage, "per-page", 50, "Number of templates per page (max 100)")
	cmd.Flags().IntVar(&token, "token", 0, "Page number to retrieve (page-token pagination)")

	return cmd
}
