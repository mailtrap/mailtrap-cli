package tokens

import (
	"strings"

	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

const expiresAtUsage = "Token expiration as an ISO 8601 date-time; pass 'never' for a token that never expires; omit for the server default (a 1-year default is being rolled out)"

// expiresAtValue maps the --expires-at flag value to the request body value:
// "never" (case-insensitive) becomes an explicit JSON null, anything else is
// passed through verbatim for the server to validate.
func expiresAtValue(expiresAt string) interface{} {
	if strings.EqualFold(expiresAt, "never") {
		return nil
	}
	return expiresAt
}

func NewCmdTokens(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokens",
		Short: "Manage API tokens",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdGet(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdDelete(f))
	cmd.AddCommand(NewCmdReset(f))

	return cmd
}
