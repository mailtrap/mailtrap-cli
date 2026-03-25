package messages_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/commands/messages"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/spf13/viper"
)

func setupTest(handler http.HandlerFunc) (*cmdutil.Factory, *bytes.Buffer, func()) {
	server := httptest.NewServer(handler)

	c := client.New("test-token")
	c.SetBaseURL(client.BaseGeneral, server.URL)

	buf := &bytes.Buffer{}
	f := &cmdutil.Factory{
		Config: func() *config.Config {
			return &config.Config{APIToken: "test-token", AccountID: "123"}
		},
		IOStreams: &cmdutil.IOStreams{
			Out:    buf,
			ErrOut: &bytes.Buffer{},
		},
		ClientOverride: c,
	}

	viper.Set("api-token", "test-token")
	viper.Set("account-id", "123")
	viper.Set("output", "table")

	return f, buf, func() {
		server.Close()
		viper.Reset()
	}
}

func TestMessagesList(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/inboxes/1/messages") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("expected Authorization header 'Bearer test-token', got %q", r.Header.Get("Authorization"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{
				"id":         1,
				"subject":    "Welcome",
				"from_email": "noreply@example.com",
				"to_email":   "user@example.com",
				"is_read":    false,
				"created_at": "2024-01-01",
			},
			{
				"id":         2,
				"subject":    "Update",
				"from_email": "info@example.com",
				"to_email":   "user@example.com",
				"is_read":    true,
				"created_at": "2024-01-02",
			},
		})
	})
	defer cleanup()

	cmd := messages.NewCmdMessages(f)
	cmd.SetArgs([]string{"list", "--inbox-id", "1"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Welcome") {
		t.Errorf("expected output to contain 'Welcome', got:\n%s", output)
	}
	if !strings.Contains(output, "Update") {
		t.Errorf("expected output to contain 'Update', got:\n%s", output)
	}
	if !strings.Contains(output, "noreply@example.com") {
		t.Errorf("expected output to contain 'noreply@example.com', got:\n%s", output)
	}
	if !strings.Contains(output, "SUBJECT") {
		t.Errorf("expected output to contain header 'SUBJECT', got:\n%s", output)
	}
}

func TestMessagesListMissingInboxID(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	buf := &bytes.Buffer{}
	f.IOStreams.Out = buf

	cmd := messages.NewCmdMessages(f)
	cmd.SetArgs([]string{"list"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when --inbox-id is missing")
	}
	if !strings.Contains(err.Error(), "--inbox-id is required") {
		t.Errorf("expected '--inbox-id is required' error, got: %v", err)
	}
}

func TestMessagesGet(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/inboxes/1/messages/1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":         1,
			"subject":    "Welcome Email",
			"from_email": "noreply@example.com",
			"to_email":   "user@example.com",
			"is_read":    false,
			"created_at": "2024-01-01",
		})
	})
	defer cleanup()

	cmd := messages.NewCmdMessages(f)
	cmd.SetArgs([]string{"get", "--inbox-id", "1", "--id", "1"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Welcome Email") {
		t.Errorf("expected output to contain 'Welcome Email', got:\n%s", output)
	}
	if !strings.Contains(output, "noreply@example.com") {
		t.Errorf("expected output to contain 'noreply@example.com', got:\n%s", output)
	}
}

func TestMessagesGetMissingFlags(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	buf := &bytes.Buffer{}
	f.IOStreams.Out = buf

	// Missing both flags
	cmd := messages.NewCmdMessages(f)
	cmd.SetArgs([]string{"get"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when flags are missing")
	}

	// Missing --id
	cmd2 := messages.NewCmdMessages(f)
	cmd2.SetArgs([]string{"get", "--inbox-id", "1"})
	cmd2.SetOut(buf)

	err = cmd2.Execute()
	if err == nil {
		t.Fatal("expected error when --id is missing")
	}
	if !strings.Contains(err.Error(), "--id is required") {
		t.Errorf("expected '--id is required' error, got: %v", err)
	}
}

func TestMessagesListJSON(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{
				"id":         1,
				"subject":    "Test",
				"from_email": "a@b.com",
				"to_email":   "c@d.com",
				"is_read":    false,
				"created_at": "2024-01-01",
			},
		})
	})
	defer cleanup()

	viper.Set("output", "json")

	cmd := messages.NewCmdMessages(f)
	cmd.SetArgs([]string{"list", "--inbox-id", "1"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	var result []map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput:\n%s", err, output)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 message, got %d", len(result))
	}
	if result[0]["subject"] != "Test" {
		t.Errorf("expected subject 'Test', got %v", result[0]["subject"])
	}
}
