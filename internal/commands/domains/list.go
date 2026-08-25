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
	ID                         int    `json:"id"`
	DomainName                 string `json:"domain_name"`
	DNSVerified                bool   `json:"dns_verified"`
	ComplianceStatus           string `json:"compliance_status"`
	InboundEnabled             bool   `json:"inbound_enabled"`
	InboundVerified            bool   `json:"inbound_verified"`
	OpenTrackingEnabled        bool   `json:"open_tracking_enabled"`
	ClickTrackingEnabled       bool   `json:"click_tracking_enabled"`
	TrackingOptOutEnabled      bool   `json:"tracking_opt_out_enabled"`
	AutoUnsubscribeLinkEnabled bool   `json:"auto_unsubscribe_link_enabled"`
}

type domainListResponse struct {
	Data []Domain `json:"data"`
}

var domainColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "DOMAIN", Field: "domain_name"},
	{Header: "DNS VERIFIED", Field: "dns_verified"},
	{Header: "COMPLIANCE", Field: "compliance_status"},
	{Header: "INBOUND ENABLED", Field: "inbound_enabled"},
	{Header: "INBOUND VERIFIED", Field: "inbound_verified"},
}

var domainSettingsColumns = append(
	append([]output.Column{}, domainColumns...),
	output.Column{Header: "OPEN TRACKING", Field: "open_tracking_enabled"},
	output.Column{Header: "CLICK TRACKING", Field: "click_tracking_enabled"},
	output.Column{Header: "TRACKING OPT OUT", Field: "tracking_opt_out_enabled"},
	output.Column{Header: "AUTO UNSUBSCRIBE", Field: "auto_unsubscribe_link_enabled"},
)

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
