package sandbox_send

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdSandboxSend(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sandbox-send",
		Short: "Send emails to a Mailtrap sandbox inbox",
		Long:  "Send single or batch emails to a Mailtrap sandbox inbox for testing.",
	}

	cmd.AddCommand(NewCmdSingle(f))
	cmd.AddCommand(NewCmdBatch(f))

	return cmd
}

// emailAddr represents an email address with an optional name.
type emailAddr struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

// parseEmailAddr parses an email address string in either "email" or "Name <email>" format.
func parseEmailAddr(s string) (emailAddr, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return emailAddr{}, fmt.Errorf("empty email address")
	}

	// Try "Name <email>" format.
	re := regexp.MustCompile(`^(.+?)\s*<([^>]+)>$`)
	if matches := re.FindStringSubmatch(s); matches != nil {
		return emailAddr{
			Name:  strings.TrimSpace(matches[1]),
			Email: strings.TrimSpace(matches[2]),
		}, nil
	}

	// Plain email address.
	return emailAddr{Email: s}, nil
}

// parseEmailAddrs parses a slice of email address strings.
func parseEmailAddrs(addrs []string) ([]emailAddr, error) {
	result := make([]emailAddr, 0, len(addrs))
	for _, s := range addrs {
		addr, err := parseEmailAddr(s)
		if err != nil {
			return nil, fmt.Errorf("invalid email address %q: %w", s, err)
		}
		result = append(result, addr)
	}
	return result, nil
}
