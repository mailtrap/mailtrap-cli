package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type APIError struct {
	StatusCode int
	Message    string `json:"error"`
	// Errors is the response's "errors" field as sent: a string, an array of
	// messages, or an object mapping fields to a message or an array of them.
	Errors json.RawMessage `json:"errors,omitempty"`
}

// newAPIError builds the error for a failed response. A body that isn't a JSON
// object of this shape, or carries no message in it, becomes the message as is.
func newAPIError(statusCode int, body []byte) *APIError {
	apiErr := &APIError{StatusCode: statusCode}
	if json.Unmarshal(body, apiErr) != nil || (apiErr.Message == "" && len(errorMessages(apiErr.Errors)) == 0) {
		return &APIError{StatusCode: statusCode, Message: string(body)}
	}
	return apiErr
}

func (e *APIError) Error() string {
	messages := errorMessages(e.Errors)
	if e.Message != "" {
		messages = append([]string{e.Message}, messages...)
	}
	if len(messages) == 0 {
		return fmt.Sprintf("API error %d", e.StatusCode)
	}
	return fmt.Sprintf("API error %d: %s", e.StatusCode, strings.Join(messages, "; "))
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
		names := make([]string, 0, len(fields))
		for name := range fields {
			names = append(names, name)
		}
		sort.Strings(names)

		var messages []string
		for _, name := range names {
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
