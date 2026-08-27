package trackingoptouts

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

const trackingOptOutsPath = "/api/tracking_opt_outs"

type TrackingOptOut struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	DomainName string `json:"domain_name,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
}

var trackingOptOutColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "EMAIL", Field: "email"},
	{Header: "DOMAIN_NAME", Field: "domain_name"},
	{Header: "CREATED_AT", Field: "created_at"},
}

func NewCmdTrackingOptOuts(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tracking-opt-outs",
		Short: "Manage tracking opt-outs",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdDelete(f))

	return cmd
}
