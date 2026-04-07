package account_access

import (
	"context"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type Specifier struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Resource struct {
	ResourceID   int    `json:"resource_id"`
	ResourceType string `json:"resource_type"`
	AccessLevel  int    `json:"access_level"`
}

type AccountAccess struct {
	ID            int        `json:"id"`
	SpecifierType string     `json:"specifier_type"`
	Specifier     Specifier  `json:"specifier"`
	Resources     []Resource `json:"resources"`
}

type accountAccessRow struct {
	ID            int    `json:"id"`
	SpecifierType string `json:"specifier_type"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	AccessLevel   int    `json:"access_level"`
	ResourceType  string `json:"resource_type"`
}

var accountAccessColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "TYPE", Field: "specifier_type"},
	{Header: "EMAIL", Field: "email"},
	{Header: "NAME", Field: "name"},
	{Header: "ACCESS_LEVEL", Field: "access_level"},
	{Header: "RESOURCE_TYPE", Field: "resource_type"},
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all account accesses",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			_, err = config.RequireAccountID()
			if err != nil {
				return err
			}

			path := cmdutil.AccountPath("account_accesses")

			var accesses []AccountAccess
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &accesses); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()

			if format == output.FormatJSON {
				return output.Print(f.IOStreams.Out, format, accesses, nil)
			}

			var rows []accountAccessRow
			for _, a := range accesses {
				topLevel := 0
				topResource := ""
				if len(a.Resources) > 0 {
					topLevel = a.Resources[0].AccessLevel
					topResource = a.Resources[0].ResourceType
				}
				rows = append(rows, accountAccessRow{
					ID:            a.ID,
					SpecifierType: a.SpecifierType,
					Email:         a.Specifier.Email,
					Name:          a.Specifier.Name,
					AccessLevel:   topLevel,
					ResourceType:  topResource,
				})
			}
			return output.Print(f.IOStreams.Out, format, rows, accountAccessColumns)
		},
	}

	return cmd
}
