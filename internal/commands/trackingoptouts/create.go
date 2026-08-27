package trackingoptouts

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type trackingOptOutResponse struct {
	Data TrackingOptOut `json:"data"`
}

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var (
		email    string
		domainID int64
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Opt an email address out of open and click tracking",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("email", email); err != nil {
				return err
			}
			if !cmd.Flags().Changed("domain-id") {
				return fmt.Errorf("--domain-id is required")
			}
			if domainID <= 0 {
				return fmt.Errorf("--domain-id must be greater than 0")
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			body := map[string]interface{}{
				"email":     email,
				"domain_id": domainID,
			}

			var resp trackingOptOutResponse
			if err := c.Post(context.Background(), client.BaseGeneral, trackingOptOutsPath, body, &resp); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, resp.Data, trackingOptOutColumns)
		},
	}

	cmd.Flags().StringVar(&email, "email", "", "Email address to opt out of tracking (required)")
	cmd.Flags().Int64Var(&domainID, "domain-id", 0, "ID of the sending domain the opt-out applies to (required)")

	return cmd
}
