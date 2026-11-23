package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"
)

type APIError struct {
	StatusCode int
	Message    string
}

// newAPIError builds the error for a failed response. Its message is the body's
// "error" text followed by its "errors" (see errorMessages), joined with "; ".
// A body that isn't a JSON object of this shape, or carries no message in it,
// becomes the message as is.
func newAPIError(statusCode int, body []byte) *APIError {
	var parsed struct {
		Error  string          `json:"error"`
		Errors json.RawMessage `json:"errors"`
	}
	if json.Unmarshal(body, &parsed) == nil {
		messages := errorMessages(parsed.Errors)
		if parsed.Error != "" {
			messages = append([]string{parsed.Error}, messages...)
		}
		if len(messages) > 0 {
			return &APIError{StatusCode: statusCode, Message: strings.Join(messages, "; ")}
		}
	}
	return &APIError{StatusCode: statusCode, Message: string(body)}
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("API error %d", e.StatusCode)
}

// errorMessages flattens an "errors" value into readable messages: a string as
// is, an array element by element, and an object as "field message" pairs
// sorted by field ("base" errors, which belong to no field, without one).
func errorMessages(raw json.RawMessage) []string {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 || bytes.Equal(raw, []byte("null")) {
		return nil
	}

	var text string
	if json.Unmarshal(raw, &text) == nil {
		if text == "" {
			return nil
		}
		return []string{text}
	}

	var list []json.RawMessage
	if json.Unmarshal(raw, &list) == nil {
		var messages []string
		for _, item := range list {
			messages = append(messages, errorMessages(item)...)
		}
		return messages
	}

	var fields map[string]json.RawMessage
	if json.Unmarshal(raw, &fields) == nil {
		var messages []string
		for _, name := range slices.Sorted(maps.Keys(fields)) {
			for _, message := range errorMessages(fields[name]) {
				if name != "base" {
					message = name + " " + message
				}
				messages = append(messages, message)
			}
		}
		return messages
	}

	// A number or a boolean prints as sent.
	return []string{string(raw)}
}
