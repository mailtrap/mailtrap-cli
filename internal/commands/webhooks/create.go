package webhooks

import (
	"context"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type WebhookWithSecret struct {
	Webhook
	SigningSecret string `json:"signing_secret,omitempty"`
}

type webhookCreateResponse struct {
	Data WebhookWithSecret `json:"data"`
}

var webhookCreateColumns = append(append([]output.Column{}, webhookColumns...), output.Column{Header: "SIGNING SECRET", Field: "signing_secret"})

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var (
		webhookURL     string
		webhookType    string
		active         bool
		payloadFormat  string
		sendingStream  string
		eventTypes     []string
		domainID       int
		inboundInboxID int
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a webhook",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("url", webhookURL); err != nil {
				return err
			}
			if err := cmdutil.RequireFlag("type", webhookType); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			if _, err := config.RequireAccountID(); err != nil {
				return err
			}

			path := cmdutil.AccountPath("webhooks")

			webhookFields := map[string]interface{}{
				"url":          webhookURL,
				"webhook_type": webhookType,
			}
			if cmd.Flags().Changed("active") {
				webhookFields["active"] = active
			}
			if cmd.Flags().Changed("payload-format") {
				webhookFields["payload_format"] = payloadFormat
			}
			if cmd.Flags().Changed("sending-stream") {
				webhookFields["sending_stream"] = sendingStream
			}
			if cmd.Flags().Changed("event-types") {
				webhookFields["event_types"] = eventTypes
			}
			if cmd.Flags().Changed("domain-id") {
				webhookFields["domain_id"] = domainID
			}
			if cmd.Flags().Changed("inbound-inbox-id") {
				webhookFields["inbound_inbox_id"] = inboundInboxID
			}

			body := map[string]interface{}{"webhook": webhookFields}

			var resp webhookCreateResponse
			if err := c.Post(context.Background(), client.BaseGeneral, path, body, &resp); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, resp.Data, webhookCreateColumns)
		},
	}

	cmd.Flags().StringVar(&webhookURL, "url", "", "Webhook URL (required)")
	cmd.Flags().StringVar(&webhookType, "type", "", "Webhook type: email_sending, audit_log, inbound_receiving (required)")
	cmd.Flags().BoolVar(&active, "active", true, "Whether the webhook is active")
	cmd.Flags().StringVar(&payloadFormat, "payload-format", "", "Payload format: json, jsonlines")
	cmd.Flags().StringVar(&sendingStream, "sending-stream", "", "Sending stream: transactional, bulk")
	cmd.Flags().StringSliceVar(&eventTypes, "event-types", nil, "Event types (comma-separated): delivery, soft_bounce, bounce, suspension, unsubscribe, open, spam_complaint, click, reject")
	cmd.Flags().IntVar(&domainID, "domain-id", 0, "Domain ID to scope the webhook to (email_sending only)")
	cmd.Flags().IntVar(&inboundInboxID, "inbound-inbox-id", 0, "Inbox ID to scope the webhook to (inbound_receiving only; omit to apply to all inboxes)")

	return cmd
}
