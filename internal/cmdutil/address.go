package cmdutil

import (
	"fmt"
	"regexp"
	"strings"
)

// EmailAddr represents an email address with an optional display name.
type EmailAddr struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

var emailAddrRe = regexp.MustCompile(`^(.+?)\s*<([^>]+)>$`)

// ParseEmailAddr parses an email address string in either "email" or "Name <email>" format.
func ParseEmailAddr(s string) (EmailAddr, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return EmailAddr{}, fmt.Errorf("empty email address")
	}

	// Try "Name <email>" format.
	if matches := emailAddrRe.FindStringSubmatch(s); matches != nil {
		return EmailAddr{
			Name:  strings.TrimSpace(matches[1]),
			Email: strings.TrimSpace(matches[2]),
		}, nil
	}

	// Plain email address.
	return EmailAddr{Email: s}, nil
}

// ParseEmailAddrs parses a slice of email address strings.
func ParseEmailAddrs(addrs []string) ([]EmailAddr, error) {
	result := make([]EmailAddr, 0, len(addrs))
	for _, s := range addrs {
		addr, err := ParseEmailAddr(s)
		if err != nil {
			return nil, fmt.Errorf("invalid email address %q: %w", s, err)
		}
		result = append(result, addr)
	}
	return result, nil
}
