package emailcampaigns

import (
	"github.com/spf13/cobra"
)

// campaignAttrs holds the writable campaign attributes shared by create and update.
type campaignAttrs struct {
	Name               string
	DomainID           int64
	FromDisplayName    string
	FromLocalPart      string
	ReplyToDisplayName string
	ReplyToLocalPart   string
	ReplyToDomain      string
	Subject            string
	BodyHTML           string
	BodyText           string
	MergeTags          []string
	DeliveryMode       string
	EmailsPerHour      int64
	ContactListIDs     []int64
	ContactSegmentIDs  []int64
}

func addAttributeFlags(cmd *cobra.Command, attrs *campaignAttrs) {
	cmd.Flags().StringVar(&attrs.Name, "name", "", "Campaign name")
	cmd.Flags().Int64Var(&attrs.DomainID, "domain-id", 0, "ID of the verified sending domain, as returned by the Sending Domains endpoints")
	cmd.Flags().StringVar(&attrs.FromDisplayName, "from-display-name", "", "Display name shown in the From header")
	cmd.Flags().StringVar(&attrs.FromLocalPart, "from-local-part", "", "Local part (before the @) of the From address")
	cmd.Flags().StringVar(&attrs.ReplyToDisplayName, "reply-to-display-name", "", "Reply-To display name")
	cmd.Flags().StringVar(&attrs.ReplyToLocalPart, "reply-to-local-part", "", "Reply-To local part (before the @)")
	cmd.Flags().StringVar(&attrs.ReplyToDomain, "reply-to-domain", "", "Reply-To domain")
	cmd.Flags().StringVar(&attrs.Subject, "subject", "", "Email subject line")
	cmd.Flags().StringVar(&attrs.BodyHTML, "body-html", "", "HTML body of the email (the design)")
	cmd.Flags().StringVar(&attrs.BodyText, "body-text", "", "Plain-text alternative of the email body")
	cmd.Flags().StringSliceVar(&attrs.MergeTags, "merge-tags", nil, "Merge tag names used in the content (comma-separated), without {{ }}")
	cmd.Flags().StringVar(&attrs.DeliveryMode, "delivery-mode", "", "Delivery mode: rapid, gradual")
	cmd.Flags().Int64Var(&attrs.EmailsPerHour, "emails-per-hour", 0, "Emails per hour when delivery mode is gradual")
	cmd.Flags().Int64SliceVar(&attrs.ContactListIDs, "contact-list-ids", nil, "Contact list IDs to send to (comma-separated)")
	cmd.Flags().Int64SliceVar(&attrs.ContactSegmentIDs, "contact-segment-ids", nil, "Contact segment IDs to send to (comma-separated)")
}

// buildAttributesBody builds the flat request body (no wrapper key) from the flags
// that were explicitly set, assembling the nested reply_to and template_attributes objects.
func buildAttributesBody(cmd *cobra.Command, attrs *campaignAttrs) map[string]interface{} {
	body := map[string]interface{}{}

	if cmd.Flags().Changed("name") {
		body["name"] = attrs.Name
	}
	if cmd.Flags().Changed("domain-id") {
		body["domain_id"] = attrs.DomainID
	}
	if cmd.Flags().Changed("from-display-name") {
		body["from_display_name"] = attrs.FromDisplayName
	}
	if cmd.Flags().Changed("from-local-part") {
		body["from_local_part"] = attrs.FromLocalPart
	}

	replyTo := map[string]interface{}{}
	if cmd.Flags().Changed("reply-to-display-name") {
		replyTo["display_name"] = attrs.ReplyToDisplayName
	}
	if cmd.Flags().Changed("reply-to-local-part") {
		replyTo["local_part"] = attrs.ReplyToLocalPart
	}
	if cmd.Flags().Changed("reply-to-domain") {
		replyTo["domain"] = attrs.ReplyToDomain
	}
	if len(replyTo) > 0 {
		body["reply_to"] = replyTo
	}

	templateAttrs := map[string]interface{}{}
	if cmd.Flags().Changed("subject") {
		templateAttrs["subject"] = attrs.Subject
	}
	if cmd.Flags().Changed("body-html") {
		templateAttrs["body_html"] = attrs.BodyHTML
	}
	if cmd.Flags().Changed("body-text") {
		templateAttrs["body_text"] = attrs.BodyText
	}
	if cmd.Flags().Changed("merge-tags") {
		templateAttrs["merge_tags"] = attrs.MergeTags
	}
	if len(templateAttrs) > 0 {
		body["template_attributes"] = templateAttrs
	}

	if cmd.Flags().Changed("delivery-mode") {
		body["delivery_mode"] = attrs.DeliveryMode
	}
	if cmd.Flags().Changed("emails-per-hour") {
		body["delivery_options"] = map[string]interface{}{"emails_per_hour": attrs.EmailsPerHour}
	}
	if cmd.Flags().Changed("contact-list-ids") {
		body["contact_list_ids"] = attrs.ContactListIDs
	}
	if cmd.Flags().Changed("contact-segment-ids") {
		body["contact_segment_ids"] = attrs.ContactSegmentIDs
	}

	return body
}
