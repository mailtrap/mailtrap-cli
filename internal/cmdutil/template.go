package cmdutil

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// ValidateTemplateFlags checks send flags against the two request shapes:
// a template send takes its subject, bodies and category from the template,
// any other send needs --subject.
func ValidateTemplateFlags(cmd *cobra.Command, templateUUID, templateVariables string) error {
	if templateUUID == "" {
		if templateVariables != "" {
			return errors.New("--template-variables requires --template-uuid")
		}
		subject, _ := cmd.Flags().GetString("subject")
		return RequireFlag("subject", subject)
	}

	for _, name := range []string{"subject", "text", "html", "category"} {
		if cmd.Flags().Changed(name) {
			return fmt.Errorf("--%s cannot be used with --template-uuid", name)
		}
	}
	return nil
}

// ParseTemplateVariables parses the --template-variables JSON object. Values
// stay raw so they are sent exactly as given, including large integers.
func ParseTemplateVariables(raw string) (map[string]json.RawMessage, error) {
	var vars map[string]json.RawMessage
	if err := json.Unmarshal([]byte(raw), &vars); err != nil {
		return nil, fmt.Errorf("invalid --template-variables JSON object: %w", err)
	}
	if vars == nil {
		return nil, errors.New("invalid --template-variables JSON object: got null")
	}
	return vars, nil
}
