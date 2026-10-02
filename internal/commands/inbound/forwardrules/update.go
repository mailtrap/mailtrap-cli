package forwardrules

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var (
		inboxID      string
		ruleID       string
		name         string
		conditions   string
		destinations []string
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a forward rule",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("inbox-id", inboxID); err != nil {
				return err
			}
			if err := cmdutil.RequireFlag("id", ruleID); err != nil {
				return err
			}

			body := map[string]interface{}{}
			if cmd.Flags().Changed("name") {
				body["name"] = name
			}
			if cmd.Flags().Changed("conditions") {
				parsed, err := parseConditions(conditions)
				if err != nil {
					return err
				}
				body["conditions"] = parsed
			}
			if cmd.Flags().Changed("destinations") {
				body["destinations"] = buildDestinations(destinations)
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/api/inbound/inboxes/%s/forward_rules/%s", inboxID, ruleID)

			var resp forwardRuleResponse
			if err := c.Patch(context.Background(), client.BaseGeneral, path, body, &resp); err != nil {
				return err
			}

			return output.Print(f.IOStreams.Out, cmdutil.GetOutputFormat(), resp.Data, forwardRuleColumns)
		},
	}

	cmd.Flags().StringVar(&inboxID, "inbox-id", "", "Inbox ID (required)")
	cmd.Flags().StringVar(&ruleID, "id", "", "Forward rule ID (required)")
	cmd.Flags().StringVar(&name, "name", "", "Rule name")
	cmd.Flags().StringVar(&conditions, "conditions", "", conditionsUsage)
	cmd.Flags().StringSliceVar(&destinations, "destinations", nil, destinationsUsage)

	return cmd
}
