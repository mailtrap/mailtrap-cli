package suppressions

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type suppressionResponse struct {
	Data Suppression `json:"data"`
}

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var (
		email         string
		domainID      int64
		sendingStream string
		suppressType  string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Add an email address to the suppression list",
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
			if err := cmdutil.RequireFlag("sending-stream", sendingStream); err != nil {
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

			body := map[string]interface{}{
				"email":          email,
				"domain_id":      domainID,
				"sending_stream": sendingStream,
			}
			if cmd.Flags().Changed("type") {
				body["type"] = suppressType
			}

			var resp suppressionResponse
			if err := c.Post(context.Background(), client.BaseGeneral, cmdutil.AccountPath("suppressions"), body, &resp); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, resp.Data, suppressionColumns)
		},
	}

	cmd.Flags().StringVar(&email, "email", "", "Email address to suppress (required)")
	cmd.Flags().Int64Var(&domainID, "domain-id", 0, "ID of the sending domain the suppression applies to (required)")
	cmd.Flags().StringVar(&sendingStream, "sending-stream", "", "Sending stream to suppress for: transactional, bulk (required)")
	cmd.Flags().StringVar(&suppressType, "type", "", "Suppression reason: hard bounce, spam complaint, unsubscription, manual import")

	return cmd
}
