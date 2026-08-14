package emailcampaigns

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	attrs := &campaignAttrs{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an email campaign",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("name", attrs.Name); err != nil {
				return err
			}
			if !cmd.Flags().Changed("domain-id") {
				return fmt.Errorf("--domain-id is required")
			}
			if err := cmdutil.RequireFlag("from-local-part", attrs.FromLocalPart); err != nil {
				return err
			}
			if err := cmdutil.RequireFlag("subject", attrs.Subject); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			body := buildAttributesBody(cmd, attrs)

			var resp campaignResponse
			if err := c.Post(context.Background(), client.BaseGeneral, basePath, body, &resp); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, resp.Data, campaignColumns)
		},
	}

	addAttributeFlags(cmd, attrs)
	// The attribute flags are shared with update, where all of them are optional.
	for _, name := range []string{"name", "domain-id", "from-local-part", "subject"} {
		cmd.Flags().Lookup(name).Usage += " (required)"
	}

	return cmd
}
