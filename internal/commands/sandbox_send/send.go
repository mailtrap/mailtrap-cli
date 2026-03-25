package sandbox_send

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type sendRequest struct {
	From         emailAddr   `json:"from"`
	To           []emailAddr `json:"to"`
	Subject      string      `json:"subject"`
	Text         string      `json:"text,omitempty"`
	HTML         string      `json:"html,omitempty"`
	CC           []emailAddr `json:"cc,omitempty"`
	BCC          []emailAddr `json:"bcc,omitempty"`
	Category     string      `json:"category,omitempty"`
	TemplateUUID string      `json:"template_uuid,omitempty"`
	ReplyTo      *emailAddr  `json:"reply_to,omitempty"`
}

type sendResponse struct {
	Success    bool     `json:"success"`
	MessageIDs []string `json:"message_ids"`
}

func NewCmdSingle(f *cmdutil.Factory) *cobra.Command {
	var (
		inboxID      string
		from         string
		to           []string
		subject      string
		text         string
		html         string
		cc           []string
		bcc          []string
		category     string
		templateUUID string
		replyTo      string
	)

	cmd := &cobra.Command{
		Use:   "single",
		Short: "Send a single email to a sandbox inbox",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			fromAddr, err := parseEmailAddr(from)
			if err != nil {
				return fmt.Errorf("invalid --from address: %w", err)
			}

			toAddrs, err := parseEmailAddrs(to)
			if err != nil {
				return fmt.Errorf("invalid --to address: %w", err)
			}

			req := sendRequest{
				From:         fromAddr,
				To:           toAddrs,
				Subject:      subject,
				Text:         text,
				HTML:         html,
				Category:     category,
				TemplateUUID: templateUUID,
			}

			if len(cc) > 0 {
				ccAddrs, err := parseEmailAddrs(cc)
				if err != nil {
					return fmt.Errorf("invalid --cc address: %w", err)
				}
				req.CC = ccAddrs
			}

			if len(bcc) > 0 {
				bccAddrs, err := parseEmailAddrs(bcc)
				if err != nil {
					return fmt.Errorf("invalid --bcc address: %w", err)
				}
				req.BCC = bccAddrs
			}

			if replyTo != "" {
				addr, err := parseEmailAddr(replyTo)
				if err != nil {
					return fmt.Errorf("invalid --reply-to address: %w", err)
				}
				req.ReplyTo = &addr
			}

			path := fmt.Sprintf("/api/send/%s", inboxID)

			var resp sendResponse
			if err := c.Post(context.Background(), client.BaseSandbox, path, req, &resp); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, resp, []output.Column{
				{Header: "SUCCESS", Field: "success"},
				{Header: "MESSAGE IDS", Field: "message_ids"},
			})
		},
	}

	cmd.Flags().StringVar(&inboxID, "inbox-id", "", "Sandbox inbox ID")
	cmd.Flags().StringVar(&from, "from", "", "Sender email address (e.g. 'Name <email>' or 'email')")
	cmd.Flags().StringSliceVar(&to, "to", nil, "Recipient email address (can be repeated)")
	cmd.Flags().StringVar(&subject, "subject", "", "Email subject")
	cmd.Flags().StringVar(&text, "text", "", "Plain text body")
	cmd.Flags().StringVar(&html, "html", "", "HTML body")
	cmd.Flags().StringSliceVar(&cc, "cc", nil, "CC recipient (can be repeated)")
	cmd.Flags().StringSliceVar(&bcc, "bcc", nil, "BCC recipient (can be repeated)")
	cmd.Flags().StringVar(&category, "category", "", "Email category")
	cmd.Flags().StringVar(&templateUUID, "template-uuid", "", "Template UUID")
	cmd.Flags().StringVar(&replyTo, "reply-to", "", "Reply-to email address")

	_ = cmd.MarkFlagRequired("inbox-id")
	_ = cmd.MarkFlagRequired("from")
	_ = cmd.MarkFlagRequired("to")
	_ = cmd.MarkFlagRequired("subject")

	return cmd
}
