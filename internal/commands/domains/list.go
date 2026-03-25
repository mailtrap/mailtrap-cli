package domains

import (
	"context"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type Domain struct {
	ID               int    `json:"id"`
	DomainName       string `json:"domain_name"`
	DNSVerified      bool   `json:"dns_verified"`
	ComplianceStatus string `json:"compliance_status"`
}

type domainListResponse struct {
	Data []Domain `json:"data"`
}

var domainColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "DOMAIN", Field: "domain_name"},
	{Header: "DNS VERIFIED", Field: "dns_verified"},
	{Header: "COMPLIANCE", Field: "compliance_status"},
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all sending domains",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			_, err = config.RequireAccountID()
			if err != nil {
				return err
			}

			path := cmdutil.AccountPath("sending_domains")

			var resp domainListResponse
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &resp); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, resp.Data, domainColumns)
		},
	}

	return cmd
}
