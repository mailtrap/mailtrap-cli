package companyinfo

import (
	"context"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	var domainID string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get company info for a sending domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("domain-id", domainID); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			var resp companyInfoResponse
			if err := c.Get(context.Background(), client.BaseGeneral, companyInfoPath(domainID), nil, &resp); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, resp.Data, companyInfoColumns)
		},
	}

	cmd.Flags().StringVar(&domainID, "domain-id", "", "Sending domain ID (required)")

	return cmd
}
