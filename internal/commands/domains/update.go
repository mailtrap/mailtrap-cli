package domains

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var (
		domainID            string
		openTracking        bool
		clickTracking       bool
		trackingOptOut      bool
		autoUnsubscribeLink bool
		inboundEnabled      bool
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a sending domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", domainID); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			if _, err := config.RequireAccountID(); err != nil {
				return err
			}

			path := cmdutil.AccountPath("sending_domains", domainID)

			domainFields := map[string]interface{}{}
			if cmd.Flags().Changed("open-tracking") {
				domainFields["open_tracking_enabled"] = openTracking
			}
			if cmd.Flags().Changed("click-tracking") {
				domainFields["click_tracking_enabled"] = clickTracking
			}
			if cmd.Flags().Changed("tracking-opt-out") {
				domainFields["tracking_opt_out_enabled"] = trackingOptOut
			}
			if cmd.Flags().Changed("auto-unsubscribe-link") {
				domainFields["auto_unsubscribe_link_enabled"] = autoUnsubscribeLink
			}
			if cmd.Flags().Changed("inbound-enabled") {
				domainFields["inbound_enabled"] = inboundEnabled
			}

			if len(domainFields) == 0 {
				return fmt.Errorf("at least one attribute flag is required")
			}

			body := map[string]interface{}{"sending_domain": domainFields}

			var domain Domain
			if err := c.Patch(context.Background(), client.BaseGeneral, path, body, &domain); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, domain, domainSettingsColumns)
		},
	}

	cmd.Flags().StringVar(&domainID, "id", "", "Domain ID (required)")
	cmd.Flags().BoolVar(&openTracking, "open-tracking", false, "Enable open tracking for emails sent from this domain")
	cmd.Flags().BoolVar(&clickTracking, "click-tracking", false, "Enable click tracking for links in emails sent from this domain")
	cmd.Flags().BoolVar(&trackingOptOut, "tracking-opt-out", false, "Enable the tracking opt-out link in tracked emails, requires open or click tracking")
	cmd.Flags().BoolVar(&autoUnsubscribeLink, "auto-unsubscribe-link", false, "Automatically add an unsubscribe link to emails")
	cmd.Flags().BoolVar(&inboundEnabled, "inbound-enabled", false, "Enable inbound email so the domain can be attached to an inbound inbox as a catch-all")

	return cmd
}
