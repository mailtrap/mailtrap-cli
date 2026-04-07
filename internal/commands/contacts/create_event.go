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
	var eventName string
	var params string

	cmd := &cobra.Command{
		Use:   "create-event",
		Short: "Create an event for a contact",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", contactID); err != nil {
				return err
			}
			if err := cmdutil.RequireFlag("name", eventName); err != nil {
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
				"name": eventName,
			}

			if params != "" {
				var paramsObj interface{}
				if err := json.Unmarshal([]byte(params), &paramsObj); err != nil {
					return fmt.Errorf("parse params JSON: %w", err)
				}
				body["params"] = paramsObj
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

	cmd.Flags().StringVar(&contactID, "id", "", "Contact ID (required)")
	cmd.Flags().StringVar(&eventName, "name", "", "Event name (required, max 255 chars)")
	cmd.Flags().StringVar(&params, "params", "", "Event params as JSON string")

	return cmd
}
