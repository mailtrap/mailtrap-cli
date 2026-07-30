package emailcampaigns

import (
	"context"
	"fmt"
	"net/url"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type ReplyTo struct {
	DisplayName string `json:"display_name,omitempty"`
	LocalPart   string `json:"local_part,omitempty"`
	Domain      string `json:"domain,omitempty"`
}

type RecipientError struct {
	Message   string `json:"message"`
	RcptIndex int    `json:"rcpt_index"`
}

type StateMetadata struct {
	Reason      *string          `json:"reason,omitempty"`
	Error       *string          `json:"error,omitempty"`
	Errors      []RecipientError `json:"errors,omitempty"`
	ScheduledAt *string          `json:"scheduled_at,omitempty"`
}

type DeliveryOptions struct {
	EmailsPerHour *int64 `json:"emails_per_hour,omitempty"`
}

type Template struct {
	ID        int64    `json:"id"`
	Subject   string   `json:"subject"`
	MergeTags []string `json:"merge_tags,omitempty"`
	BodyHTML  *string  `json:"body_html,omitempty"`
	BodyText  *string  `json:"body_text,omitempty"`
}

type EmailCampaign struct {
	ID                   int64            `json:"id"`
	DomainID             int64            `json:"domain_id"`
	DomainName           string           `json:"domain_name"`
	Name                 string           `json:"name"`
	FromLocalPart        string           `json:"from_local_part"`
	FromDisplayName      string           `json:"from_display_name"`
	ReplyTo              *ReplyTo         `json:"reply_to,omitempty"`
	CurrentState         string           `json:"current_state"`
	CurrentStateMetadata *StateMetadata   `json:"current_state_metadata,omitempty"`
	CreatedAt            string           `json:"created_at"`
	UpdatedAt            string           `json:"updated_at"`
	LastStartedAt        *string          `json:"last_started_at,omitempty"`
	LastStartedAtDate    *string          `json:"last_started_at_date,omitempty"`
	RecipientTotalCount  *int64           `json:"recipient_total_count,omitempty"`
	ContactListIDs       []int64          `json:"contact_list_ids,omitempty"`
	ContactSegmentIDs    []int64          `json:"contact_segment_ids,omitempty"`
	DeliveryMode         string           `json:"delivery_mode"`
	DeliveryOptions      *DeliveryOptions `json:"delivery_options,omitempty"`
	Template             *Template        `json:"template,omitempty"`
}

type campaignListResponse struct {
	Data []EmailCampaign `json:"data"`
}

type campaignResponse struct {
	Data EmailCampaign `json:"data"`
}

var campaignColumns = []output.Column{
	{Header: "ID", Field: "id"},
	{Header: "NAME", Field: "name"},
	{Header: "STATE", Field: "current_state"},
	{Header: "DOMAIN", Field: "domain_name"},
	{Header: "RECIPIENTS", Field: "recipient_total_count"},
	{Header: "CREATED", Field: "created_at"},
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var (
		perPage int
		search  string
		token   int
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all email campaigns",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			query := url.Values{}
			if cmd.Flags().Changed("per-page") {
				query.Set("per_page", fmt.Sprintf("%d", perPage))
			}
			if search != "" {
				query.Set("search", search)
			}
			if cmd.Flags().Changed("token") {
				query.Set("token", fmt.Sprintf("%d", token))
			}

			var resp campaignListResponse
			if err := c.Get(context.Background(), client.BaseGeneral, basePath, query, &resp); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, resp.Data, campaignColumns)
		},
	}

	cmd.Flags().IntVar(&perPage, "per-page", 0, "Number of campaigns per page (default 50, max 100)")
	cmd.Flags().StringVar(&search, "search", "", "Filter campaigns by name")
	cmd.Flags().IntVar(&token, "token", 0, "Page number to retrieve (page-token pagination)")

	return cmd
}
