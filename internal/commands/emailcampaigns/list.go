package emailcampaigns

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

// campaignResponse unwraps the data envelope of a single-campaign response and
// keeps the campaign as the API returned it.
type campaignResponse struct {
	Data json.RawMessage `json:"data"`
}

var campaignColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "NAME", Field: "name"},
	{Header: "STATE", Field: "current_state"},
	{Header: "DOMAIN", Field: "domain_name"},
	{Header: "RECIPIENTS", Field: "recipient_total_count"},
	{Header: "CREATED", Field: "created_at"},
}

var campaignsPage = output.Page{
	Items:      "data",
	Cursor:     []string{"pagination", "next_token"},
	CursorFlag: "token",
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var (
		perPage int
		search  string
		token   int
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all email campaigns",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			query := url.Values{}
			if cmd.Flags().Changed("per-page") {
				query.Set("per_page", fmt.Sprintf("%d", perPage))
			}
			if search != "" {
				query.Set("search", search)
			}
			if cmd.Flags().Changed("token") {
				query.Set("token", fmt.Sprintf("%d", token))
			}

			var resp json.RawMessage
			if err := c.Get(context.Background(), client.BaseGeneral, basePath, query, &resp); err != nil {
				return err
			}

			return output.PrintPage(f.IOStreams.Out, cmdutil.GetOutputFormat(), resp, campaignsPage, campaignColumns)
		},
	}

	cmd.Flags().IntVar(&perPage, "per-page", 50, "Number of campaigns per page (max 100)")
	cmd.Flags().StringVar(&search, "search", "", "Filter campaigns by name")
	cmd.Flags().IntVar(&token, "token", 0, "Page number to retrieve (page-token pagination)")

	return cmd
}
