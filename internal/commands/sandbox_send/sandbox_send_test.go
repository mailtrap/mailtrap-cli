package sandbox_send_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/commands/sandbox_send"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/spf13/viper"
)

func setupTest(handler http.HandlerFunc) (*cmdutil.Factory, *bytes.Buffer, func()) {
	server := httptest.NewServer(handler)

	c := client.New("test-token")
	c.SetBaseURL(client.BaseSandbox, server.URL)

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

func TestSandboxSendSingle(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/send/1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]interface{}
		json.Unmarshal(body, &reqBody)

		from, ok := reqBody["from"].(map[string]interface{})
		if !ok {
			t.Fatal("expected 'from' field in request body")
		}
		if from["email"] != "test@example.com" {
			t.Errorf("expected from email 'test@example.com', got %v", from["email"])
		}

		to, ok := reqBody["to"].([]interface{})
		if !ok || len(to) == 0 {
			t.Fatal("expected 'to' field in request body")
		}
		toFirst := to[0].(map[string]interface{})
		if toFirst["email"] != "user@example.com" {
			t.Errorf("expected to email 'user@example.com', got %v", toFirst["email"])
		}

		if reqBody["subject"] != "Test" {
			t.Errorf("expected subject 'Test', got %v", reqBody["subject"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":     true,
			"message_ids": []string{"msg-1"},
		})
	})
	defer cleanup()

	cmd := sandbox_send.NewCmdSandboxSend(f)
	cmd.SetArgs([]string{
		"single",
		"--sandbox-id", "1",
		"--from", "test@example.com",
		"--to", "user@example.com",
		"--subject", "Test",
		"--text", "Hello",
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

func TestSandboxSendSingleMissingFlags(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	buf := &bytes.Buffer{}
	f.IOStreams.Out = buf

	cmd := sandbox_send.NewCmdSandboxSend(f)
	cmd.SetArgs([]string{"single"})
	cmd.SetOut(buf)
	cmd.SetErr(&bytes.Buffer{})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when required flags are missing")
	}
}

func TestSandboxSendBatch(t *testing.T) {
	batchData := `[{"from":{"email":"test@example.com"},"to":[{"email":"user@example.com"}],"subject":"Batch Test","text":"Hello"}]`

	tmpFile, err := os.CreateTemp("", "sandbox-batch-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(batchData); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/batch/1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":   true,
			"responses": []interface{}{},
		})
	})
	defer cleanup()

	cmd := sandbox_send.NewCmdSandboxSend(f)
	cmd.SetArgs([]string{"batch", "--sandbox-id", "1", "--file", tmpFile.Name()})
	cmd.SetOut(buf)

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "true") {
		t.Errorf("expected output to contain 'true', got:\n%s", output)
	}
}
