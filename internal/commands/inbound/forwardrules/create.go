package forwardrules

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

const conditionsUsage = "Conditions as a JSON array"

const destinationsUsage = "Destination email addresses (comma-separated)"

var (
	validMatchTypes = []string{"sender", "recipient", "header"}
	validOperators  = []string{"equal", "not_equal", "contains", "starts_with", "ends_with", "empty", "not_empty"}
)

type conditionInput struct {
	MatchType string  `json:"match_type"`
	Operator  string  `json:"operator"`
	Value     *string `json:"value,omitempty"`
	HeaderKey *string `json:"header_key,omitempty"`
}

func parseConditions(raw string) ([]conditionInput, error) {
	if !strings.HasPrefix(strings.TrimSpace(raw), "[") {
		return nil, fmt.Errorf("invalid --conditions: must be a JSON array")
	}

	dec := json.NewDecoder(bytes.NewReader([]byte(raw)))
	dec.DisallowUnknownFields()

	conditions := []conditionInput{}
	if err := dec.Decode(&conditions); err != nil {
		return nil, fmt.Errorf("invalid --conditions JSON: %w", err)
	}

	for i, cond := range conditions {
		if !slices.Contains(validMatchTypes, cond.MatchType) {
			return nil, fmt.Errorf("invalid --conditions: entry %d: match_type must be one of %s", i+1, strings.Join(validMatchTypes, ", "))
		}
		if !slices.Contains(validOperators, cond.Operator) {
			return nil, fmt.Errorf("invalid --conditions: entry %d: operator must be one of %s", i+1, strings.Join(validOperators, ", "))
		}
	}

	return conditions, nil
}

func buildDestinations(emails []string) []ForwardRuleDestination {
	destinations := make([]ForwardRuleDestination, 0, len(emails))
	for _, email := range emails {
		destinations = append(destinations, ForwardRuleDestination{Email: email})
	}
	return destinations
}

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var (
		inboxID      string
		name         string
		conditions   string
		destinations []string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a forward rule",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("inbox-id", inboxID); err != nil {
				return err
			}
			if err := cmdutil.RequireFlag("name", name); err != nil {
				return err
			}

			body := map[string]interface{}{"name": name}
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

			path := fmt.Sprintf("/api/inbound/inboxes/%s/forward_rules", inboxID)

			var resp forwardRuleResponse
			if err := c.Post(context.Background(), client.BaseGeneral, path, body, &resp); err != nil {
				return err
			}

			return output.Print(f.IOStreams.Out, cmdutil.GetOutputFormat(), resp.Data, forwardRuleColumns)
		},
	}

	cmd.Flags().StringVar(&inboxID, "inbox-id", "", "Inbox ID (required)")
	cmd.Flags().StringVar(&name, "name", "", "Rule name (required)")
	cmd.Flags().StringVar(&conditions, "conditions", "", conditionsUsage)
	cmd.Flags().StringSliceVar(&destinations, "destinations", nil, destinationsUsage)

	return cmd
}
