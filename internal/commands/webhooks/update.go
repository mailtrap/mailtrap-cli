package webhooks

import (
	"context"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var (
		webhookID     string
		webhookURL    string
		active        bool
		payloadFormat string
		eventTypes    []string
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a webhook",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", webhookID); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			if _, err := config.RequireAccountID(); err != nil {
				return err
			}

			path := cmdutil.AccountPath("webhooks", webhookID)

			webhookFields := map[string]interface{}{}
			if cmd.Flags().Changed("url") {
				webhookFields["url"] = webhookURL
			}
			if cmd.Flags().Changed("active") {
				webhookFields["active"] = active
			}
			if cmd.Flags().Changed("payload-format") {
				webhookFields["payload_format"] = payloadFormat
			}
			if cmd.Flags().Changed("event-types") {
				webhookFields["event_types"] = eventTypes
			}

			body := map[string]interface{}{"webhook": webhookFields}

			var resp webhookResponse
			if err := c.Patch(context.Background(), client.BaseGeneral, path, body, &resp); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, resp.Data, webhookColumns)
		},
	}

	cmd.Flags().StringVar(&webhookID, "id", "", "Webhook ID (required)")
	cmd.Flags().StringVar(&webhookURL, "url", "", "Webhook URL")
	cmd.Flags().BoolVar(&active, "active", true, "Whether the webhook is active")
	cmd.Flags().StringVar(&payloadFormat, "payload-format", "", "Payload format: json, jsonlines")
	cmd.Flags().StringSliceVar(&eventTypes, "event-types", nil, "Event types (comma-separated): delivery, soft_bounce, bounce, suspension, unsubscribe, open, spam_complaint, click, reject")

	return cmd
}