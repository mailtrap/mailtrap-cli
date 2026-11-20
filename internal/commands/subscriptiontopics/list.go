package subscriptiontopics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

var topicColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "NAME", Field: "name"},
	{Header: "DESCRIPTION", Field: "description"},
	{Header: "DOMAIN ID", Field: "domain_id"},
	{Header: "CREATED AT", Field: "created_at"},
}

var topicsPage = output.Page{
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
		Short: "List subscription topics, one page at a time",
		Long:  "List the account's subscription topics, one page at a time. Deleted topics are left out. An empty DOMAIN ID means the topic covers every domain.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			query := url.Values{}
			page := topicsPage
			if cmd.Flags().Changed("per-page") {
				query.Set("per_page", fmt.Sprintf("%d", perPage))
				// Without per_page the API reads the next token at its default page size.
				page.NextArgs = fmt.Sprintf("--per-page %d", perPage)
			}
			if cmd.Flags().Changed("token") {
				query.Set("token", fmt.Sprintf("%d", token))
			}

			var resp json.RawMessage
			if err := c.Get(context.Background(), client.BaseGeneral, basePath, query, &resp); err != nil {
				return err
			}

			return output.PrintPage(f.IOStreams.Out, cmdutil.GetOutputFormat(), resp, page, topicColumns)
		},
	}

	cmd.Flags().IntVar(&perPage, "per-page", 50, "Number of topics per page (max 100)")
	cmd.Flags().IntVar(&token, "token", 0, "Page number to retrieve (page-token pagination)")

	return cmd
}
