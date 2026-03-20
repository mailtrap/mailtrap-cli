package contacts

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdCreateEvent(f *cmdutil.Factory) *cobra.Command {
	var contactID string
	var eventType string
	var data string

	cmd := &cobra.Command{
		Use:   "create-event",
		Short: "Create an event for a contact",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", contactID); err != nil {
				return err
			}
			if err := cmdutil.RequireFlag("type", eventType); err != nil {
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

			path := cmdutil.AccountPath("contacts", fmt.Sprintf("%s", contactID), "events")

			body := map[string]interface{}{
				"type": eventType,
			}

			if data != "" {
				var dataObj interface{}
				if err := json.Unmarshal([]byte(data), &dataObj); err != nil {
					return fmt.Errorf("parse data JSON: %w", err)
				}
				body["data"] = dataObj
			}

			var resp interface{}
			if err := c.Post(context.Background(), client.BaseGeneral, path, body, &resp); err != nil {
				return err
			}

			respJSON, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				return err
			}

			fmt.Fprintln(f.IOStreams.Out, string(respJSON))
			return nil
		},
	}

	cmd.Flags().StringVar(&contactID, "id", "", "Contact ID")
	cmd.Flags().StringVar(&eventType, "type", "", "Event type (required)")
	cmd.Flags().StringVar(&data, "data", "", "Event data as JSON string")

	return cmd
}
