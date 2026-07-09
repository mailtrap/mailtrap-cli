package contact_lists

import (
	"context"
	"net/url"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type ContactList struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var contactListColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "NAME", Field: "name"},
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var search string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all contact lists",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			_, err = config.RequireAccountID()
			if err != nil {
				return err
			}

			path := cmdutil.AccountPath("contacts", "lists")

			query := url.Values{}
			if search != "" {
				query.Set("search", search)
			}

			var lists []ContactList
			if err := c.Get(context.Background(), client.BaseGeneral, path, query, &lists); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, lists, contactListColumns)
		},
	}

	cmd.Flags().StringVar(&search, "search", "", "Filter by name (case-insensitive prefix match)")

	return cmd
}
