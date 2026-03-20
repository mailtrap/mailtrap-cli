package client

import "fmt"

type APIError struct {
	StatusCode int
	Message    string            `json:"error"`
	Errors     map[string]string `json:"errors,omitempty"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("API error %d", e.StatusCode)
}
