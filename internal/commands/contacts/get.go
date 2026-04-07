package contacts

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type Contact struct {
	ID        string                 `json:"id"`
	Email     string                 `json:"email"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	ListIDs   []int                  `json:"list_ids,omitempty"`
	Status    string                 `json:"status,omitempty"`
	CreatedAt interface{}            `json:"created_at"`
	UpdatedAt interface{}            `json:"updated_at"`
}

type contactResponse struct {
	Data Contact `json:"data"`
}

var contactColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "EMAIL", Field: "email"},
	{Header: "STATUS", Field: "status"},
	{Header: "CREATED_AT", Field: "created_at"},
	{Header: "UPDATED_AT", Field: "updated_at"},
}

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	var contactID string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a contact",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", contactID); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			_, err = config.RequireAccountID()
			if err != nil {
				return err
			}

			path := cmdutil.AccountPath("contacts", fmt.Sprintf("%s", contactID))

			var resp contactResponse
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &resp); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, resp.Data, contactColumns)
		},
	}

	cmd.Flags().StringVar(&contactID, "id", "", "Contact ID")

	return cmd
}
