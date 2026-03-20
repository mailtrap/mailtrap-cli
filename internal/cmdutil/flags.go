package cmdutil

import (
	"fmt"
	"strings"

	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/viper"
)

func AccountPath(segments ...string) string {
	parts := append([]string{"/api/accounts", viper.GetString("account-id")}, segments...)
	return strings.Join(parts, "/")
}

func GetOutputFormat() output.Format {
	f := viper.GetString("output")
	switch output.Format(f) {
	case output.FormatJSON:
		return output.FormatJSON
	case output.FormatText:
		return output.FormatText
	default:
		return output.FormatTable
	}
}

func RequireFlag(name, value string) error {
	if value == "" {
		return fmt.Errorf("--%s is required", name)
	}
	return nil
}
