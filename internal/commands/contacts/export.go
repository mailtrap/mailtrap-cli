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

func NewCmdExport(f *cmdutil.Factory) *cobra.Command {
	var listIDs []int
	var status string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export contacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			_, err = config.RequireAccountID()
			if err != nil {
				return err
			}

			path := cmdutil.AccountPath("contacts", "exports")

			var filters []map[string]interface{}
			if len(listIDs) > 0 {
				filters = append(filters, map[string]interface{}{
					"name":     "list_id",
					"operator": "equal",
					"value":    listIDs,
				})
			}
			if status != "" {
				filters = append(filters, map[string]interface{}{
					"name":     "subscription_status",
					"operator": "equal",
					"value":    status,
				})
			}

			body := map[string]interface{}{}
			if len(filters) > 0 {
				body["filters"] = filters
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

	cmd.Flags().IntSliceVar(&listIDs, "list-ids", nil, "Filter by list IDs")
	cmd.Flags().StringVar(&status, "status", "", "Filter by subscription status (subscribed|unsubscribed)")

	return cmd
}
