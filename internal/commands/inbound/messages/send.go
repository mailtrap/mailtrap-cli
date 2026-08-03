package messages

import (
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

// SendMessageResult is the result of a reply, reply-all, or forward (each sends a real email).
type SendMessageResult struct {
	MessageIDs []string `json:"message_ids"`
}

var sendResultColumns = []output.Column{
	{Header: "MESSAGE IDS", Field: "message_ids"},
}

// sendFlags holds the fields shared by the reply, reply-all, and forward commands.
type sendFlags struct {
	from     string
	to       []string
	cc       []string
	bcc      []string
	replyTo  string
	text     string
	html     string
	category string
}

func (s *sendFlags) register(cmd *cobra.Command) {
	cmd.Flags().StringVar(&s.from, "from", "", "Sender address, 'Name <email>' or 'email' (custom-domain inboxes only)")
	cmd.Flags().StringSliceVar(&s.to, "to", nil, "Recipient address, 'Name <email>' or 'email' (can be repeated)")
	cmd.Flags().StringSliceVar(&s.cc, "cc", nil, "CC recipient (can be repeated)")
	cmd.Flags().StringSliceVar(&s.bcc, "bcc", nil, "BCC recipient (can be repeated)")
	cmd.Flags().StringVar(&s.replyTo, "reply-to", "", "Reply-To address")
	cmd.Flags().StringVar(&s.text, "text", "", "Plain-text body")
	cmd.Flags().StringVar(&s.html, "html", "", "HTML body")
	cmd.Flags().StringVar(&s.category, "category", "", "Email API category for the sent message")
}

func (s *sendFlags) body() (map[string]interface{}, error) {
	body := map[string]interface{}{}

	if s.from != "" {
		addr, err := cmdutil.ParseEmailAddr(s.from)
		if err != nil {
			return nil, fmt.Errorf("invalid --from address: %w", err)
		}
		body["from"] = addr
	}
	if len(s.to) > 0 {
		addrs, err := cmdutil.ParseEmailAddrs(s.to)
		if err != nil {
			return nil, fmt.Errorf("invalid --to address: %w", err)
		}
		body["to"] = addrs
	}
	if len(s.cc) > 0 {
		addrs, err := cmdutil.ParseEmailAddrs(s.cc)
		if err != nil {
			return nil, fmt.Errorf("invalid --cc address: %w", err)
		}
		body["cc"] = addrs
	}
	if len(s.bcc) > 0 {
		addrs, err := cmdutil.ParseEmailAddrs(s.bcc)
		if err != nil {
			return nil, fmt.Errorf("invalid --bcc address: %w", err)
		}
		body["bcc"] = addrs
	}
	if s.replyTo != "" {
		addr, err := cmdutil.ParseEmailAddr(s.replyTo)
		if err != nil {
			return nil, fmt.Errorf("invalid --reply-to address: %w", err)
		}
		body["reply_to"] = addr
	}
	if s.text != "" {
		body["text"] = s.text
	}
	if s.html != "" {
		body["html"] = s.html
	}
	if s.category != "" {
		body["category"] = s.category
	}

	return body, nil
}
