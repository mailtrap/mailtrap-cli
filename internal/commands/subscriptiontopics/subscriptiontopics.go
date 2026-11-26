package subscriptiontopics

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// Subscription topics are token-scoped: the account is resolved from the API token, so
// commands in this package do not call config.RequireAccountID().
const basePath = "/api/subscription_topics"

func NewCmdSubscriptionTopics(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscription-topics",
		Short: "List and inspect subscription topics",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdGet(f))

	return cmd
}
