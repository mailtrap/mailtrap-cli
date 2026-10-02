package forwardrules

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	var (
		inboxID string
		ruleID  string
	)

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a forward rule",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("inbox-id", inboxID); err != nil {
				return err
			}
			if err := cmdutil.RequireFlag("id", ruleID); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/api/inbound/inboxes/%s/forward_rules/%s", inboxID, ruleID)

			var resp forwardRuleResponse
			if err := c.Get(context.Background(), client.BaseGeneral, path, nil, &resp); err != nil {
				return err
			}

			return output.Print(f.IOStreams.Out, cmdutil.GetOutputFormat(), resp.Data, forwardRuleColumns)
		},
	}

	cmd.Flags().StringVar(&inboxID, "inbox-id", "", "Inbox ID (required)")
	cmd.Flags().StringVar(&ruleID, "id", "", "Forward rule ID (required)")

	return cmd
}
