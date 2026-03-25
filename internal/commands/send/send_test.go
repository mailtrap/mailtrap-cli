package send_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/commands/send"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/spf13/viper"
)

func setupTest(handler http.HandlerFunc) (*cmdutil.Factory, *bytes.Buffer, func()) {
	server := httptest.NewServer(handler)

	c := client.New("test-token")
	c.SetBaseURL(client.BaseTransactional, server.URL)
	c.SetBaseURL(client.BaseBulk, server.URL)
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

func TestSendTransactional(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/send" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("expected Authorization header 'Bearer test-token', got %q", r.Header.Get("Authorization"))
		}

		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]interface{}
		json.Unmarshal(body, &reqBody)

		// Verify from
		from, ok := reqBody["from"].(map[string]interface{})
		if !ok {
			t.Fatal("expected 'from' field in request body")
		}
		if from["email"] != "sender@example.com" {
			t.Errorf("expected from email 'sender@example.com', got %v", from["email"])
		}

		// Verify to
		to, ok := reqBody["to"].([]interface{})
		if !ok || len(to) == 0 {
			t.Fatal("expected 'to' field in request body")
		}
		toFirst := to[0].(map[string]interface{})
		if toFirst["email"] != "recipient@example.com" {
			t.Errorf("expected to email 'recipient@example.com', got %v", toFirst["email"])
		}

		// Verify subject
		if reqBody["subject"] != "Test Subject" {
			t.Errorf("expected subject 'Test Subject', got %v", reqBody["subject"])
		}

		// Verify text body
		if reqBody["text"] != "Hello World" {
			t.Errorf("expected text 'Hello World', got %v", reqBody["text"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":     true,
			"message_ids": []string{"msg-001"},
		})
	})
	defer cleanup()

	cmd := send.NewCmdSend(f)
	cmd.SetArgs([]string{
		"transactional",
		"--from", "sender@example.com",
		"--to", "recipient@example.com",
		"--subject", "Test Subject",
		"--text", "Hello World",
	})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "true") {
		t.Errorf("expected output to contain 'true', got:\n%s", output)
	}
}

func TestSendTransactionalWithNamedFrom(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]interface{}
		json.Unmarshal(body, &reqBody)

		from := reqBody["from"].(map[string]interface{})
		if from["email"] != "sender@example.com" {
			t.Errorf("expected from email 'sender@example.com', got %v", from["email"])
		}
		if from["name"] != "Sender Name" {
			t.Errorf("expected from name 'Sender Name', got %v", from["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":     true,
			"message_ids": []string{"msg-002"},
		})
	})
	defer cleanup()

	cmd := send.NewCmdSend(f)
	cmd.SetArgs([]string{
		"transactional",
		"--from", "Sender Name <sender@example.com>",
		"--to", "recipient@example.com",
		"--subject", "Test Subject",
		"--text", "Hello",
	})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSendTransactionalMissingRequiredFlags(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	buf := &bytes.Buffer{}
	f.IOStreams.Out = buf

	cmd := send.NewCmdSend(f)
	cmd.SetArgs([]string{"transactional"})
	cmd.SetOut(buf)
	cmd.SetErr(&bytes.Buffer{})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when required flags are missing")
	}
}

func TestSendTransactionalJSON(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":     true,
			"message_ids": []string{"msg-003"},
		})
	})
	defer cleanup()

	viper.Set("output", "json")

	cmd := send.NewCmdSend(f)
	cmd.SetArgs([]string{
		"transactional",
		"--from", "sender@example.com",
		"--to", "recipient@example.com",
		"--subject", "Test",
		"--text", "Hello",
	})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput:\n%s", err, output)
	}
	if result["success"] != true {
		t.Errorf("expected success true, got %v", result["success"])
	}
}
