package cmdutil

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ParseJSONArrayFlag strictly decodes the JSON array given to --flag: unknown
// keys and data after the array are errors. Callers validate the entries.
func ParseJSONArrayFlag[T any](flag, raw string) ([]T, error) {
	if !strings.HasPrefix(strings.TrimSpace(raw), "[") {
		return nil, fmt.Errorf("invalid --%s: must be a JSON array", flag)
	}

	dec := json.NewDecoder(strings.NewReader(raw))
	dec.DisallowUnknownFields()

	entries := []T{}
	if err := dec.Decode(&entries); err != nil {
		return nil, fmt.Errorf("invalid --%s JSON: %w", flag, err)
	}
	if _, err := dec.Token(); err != io.EOF {
		return nil, fmt.Errorf("invalid --%s JSON: unexpected data after the array", flag)
	}

	return entries, nil
}
