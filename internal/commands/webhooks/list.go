package webhooks

import (
	"context"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type Webhook struct {
	ID            int      `json:"id"`
	URL           string   `json:"url"`
	Active        bool     `json:"active"`
	WebhookType   string   `json:"webhook_type"`
	PayloadFormat string   `json:"payload_format"`
	SendingStream *string  `json:"sending_stream,omitempty"`
	DomainID      *int     `json:"domain_id,omitempty"`
	EventTypes    []string `json:"event_types,omitempty"`
}

type webhookListResponse struct {
	Data []Webhook `json:"data"`
}

type webhookResponse struct {
	Data Webhook `json:"data"`
}

var webhookColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "URL", Field: "url"},
	{Header: "TYPE", Field: "webhook_type"},
	{Header: "ACTIVE", Field: "active"},
	{Header: "FORMAT", Field: "payload_format"},
	{Header: "STREAM", Field: "sending_stream"},
	{Header: "DOMAIN ID", Field: "domain_id"},
	{Header: "EVENTS", Field: "event_types"},
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all webhooks",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			if _, err := config.RequireAccountID(); err != nil {
				return err
			}

			path := cmdutil.AccountPath("webhooks")

			var resp webhookListResponse
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &resp); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, resp.Data, webhookColumns)
		},
	}

	return cmd
}