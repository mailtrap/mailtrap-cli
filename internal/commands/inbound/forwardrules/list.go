package forwardrules

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

// ForwardRule represents an inbound forward rule.
type ForwardRule struct {
	ID           int                      `json:"id"`
	Name         string                   `json:"name"`
	CreatedAt    string                   `json:"created_at"`
	UpdatedAt    string                   `json:"updated_at"`
	Conditions   []ForwardRuleCondition   `json:"conditions"`
	Destinations []ForwardRuleDestination `json:"destinations"`
}

// ForwardRuleCondition represents a forward rule condition.
type ForwardRuleCondition struct {
	MatchType string  `json:"match_type"`
	Operator  string  `json:"operator"`
	Value     *string `json:"value"`
	HeaderKey *string `json:"header_key"`
}

// ForwardRuleDestination represents a forward rule destination.
type ForwardRuleDestination struct {
	Email string `json:"email"`
}

type forwardRuleResponse struct {
	Data ForwardRule `json:"data"`
}

type forwardRulesListResponse struct {
	Data []ForwardRule `json:"data"`
}

var forwardRuleColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "NAME", Field: "name"},
	{Header: "CONDITIONS", Field: "conditions"},
	{Header: "DESTINATIONS", Field: "destinations"},
	{Header: "UPDATED AT", Field: "updated_at"},
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var inboxID string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List forward rules in an inbox",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("inbox-id", inboxID); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/api/inbound/inboxes/%s/forward_rules", inboxID)

			var resp forwardRulesListResponse
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &resp); err != nil {
				return err
			}

			return output.Print(f.IOStreams.Out, cmdutil.GetOutputFormat(), resp.Data, forwardRuleColumns)
		},
	}

	cmd.Flags().StringVar(&inboxID, "inbox-id", "", "Inbox ID (required)")

	return cmd
}
