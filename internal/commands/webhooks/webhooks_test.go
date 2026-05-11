package webhooks_test

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
	"github.com/mailtrap/mailtrap-cli/internal/commands/webhooks"
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

func TestWebhooksList(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/accounts/123/webhooks") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Api-Token") != "test-token" {
			t.Errorf("expected Api-Token header 'test-token', got %q", r.Header.Get("Api-Token"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{
					"id":             1,
					"url":            "https://example.com/mailtrap/webhooks",
					"active":         true,
					"webhook_type":   "email_sending",
					"payload_format": "json",
					"sending_stream": "transactional",
					"domain_id":      435,
					"event_types":    []string{"delivery", "bounce"},
				},
				{
					"id":             2,
					"url":            "https://example.com/audit",
					"active":         true,
					"webhook_type":   "audit_log",
					"payload_format": "json",
				},
			},
		})
	})
	defer cleanup()

	cmd := webhooks.NewCmdWebhooks(f)
	cmd.SetArgs([]string{"list"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "email_sending") {
		t.Errorf("expected output to contain 'email_sending', got:\n%s", output)
	}
	if !strings.Contains(output, "audit_log") {
		t.Errorf("expected output to contain 'audit_log', got:\n%s", output)
	}
	if !strings.Contains(output, "transactional") {
		t.Errorf("expected output to contain 'transactional', got:\n%s", output)
	}
}

func TestWebhooksListJSON(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{
					"id":             1,
					"url":            "https://example.com/mailtrap/webhooks",
					"active":         true,
					"webhook_type":   "email_sending",
					"payload_format": "json",
				},
			},
		})
	})
	defer cleanup()

	viper.Set("output", "json")

	cmd := webhooks.NewCmdWebhooks(f)
	cmd.SetArgs([]string{"list"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput:\n%s", err, buf.String())
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 webhook, got %d", len(result))
	}
	if result[0]["webhook_type"] != "email_sending" {
		t.Errorf("expected webhook_type 'email_sending', got %v", result[0]["webhook_type"])
	}
}

func TestWebhooksGet(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/accounts/123/webhooks/1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"id":             1,
				"url":            "https://example.com/mailtrap/webhooks",
				"active":         true,
				"webhook_type":   "email_sending",
				"payload_format": "json",
				"sending_stream": "transactional",
				"domain_id":      435,
				"event_types":    []string{"delivery", "bounce"},
			},
		})
	})
	defer cleanup()

	cmd := webhooks.NewCmdWebhooks(f)
	cmd.SetArgs([]string{"get", "--id", "1"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "email_sending") {
		t.Errorf("expected output to contain 'email_sending', got:\n%s", output)
	}
	if !strings.Contains(output, "https://example.com/mailtrap/webhooks") {
		t.Errorf("expected output to contain webhook URL, got:\n%s", output)
	}
}

func TestWebhooksGetMissingID(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	buf := &bytes.Buffer{}
	f.IOStreams.Out = buf

	cmd := webhooks.NewCmdWebhooks(f)
	cmd.SetArgs([]string{"get"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when --id is missing")
	}
	if !strings.Contains(err.Error(), "--id is required") {
		t.Errorf("expected '--id is required' error, got: %v", err)
	}
}

func TestWebhooksCreate(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/accounts/123/webhooks") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]map[string]interface{}
		json.Unmarshal(body, &reqBody)

		webhook := reqBody["webhook"]
		if webhook["url"] != "https://example.com/hooks" {
			t.Errorf("unexpected url: %v", webhook["url"])
		}
		if webhook["webhook_type"] != "email_sending" {
			t.Errorf("unexpected webhook_type: %v", webhook["webhook_type"])
		}
		if webhook["sending_stream"] != "transactional" {
			t.Errorf("unexpected sending_stream: %v", webhook["sending_stream"])
		}
		if events, ok := webhook["event_types"].([]interface{}); !ok || len(events) != 2 {
			t.Errorf("expected 2 event_types, got %v", webhook["event_types"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"id":             10,
				"url":            "https://example.com/hooks",
				"active":         true,
				"webhook_type":   "email_sending",
				"payload_format": "json",
				"sending_stream": "transactional",
				"event_types":    []string{"delivery", "bounce"},
				"signing_secret": "abc123secret",
			},
		})
	})
	defer cleanup()

	cmd := webhooks.NewCmdWebhooks(f)
	cmd.SetArgs([]string{
		"create",
		"--url", "https://example.com/hooks",
		"--type", "email_sending",
		"--sending-stream", "transactional",
		"--event-types", "delivery,bounce",
	})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "abc123secret") {
		t.Errorf("expected output to contain signing secret, got:\n%s", output)
	}
	if !strings.Contains(output, "email_sending") {
		t.Errorf("expected output to contain 'email_sending', got:\n%s", output)
	}
}

func TestWebhooksCreateMissingURL(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	buf := &bytes.Buffer{}
	f.IOStreams.Out = buf

	cmd := webhooks.NewCmdWebhooks(f)
	cmd.SetArgs([]string{"create", "--type", "email_sending"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when --url is missing")
	}
	if !strings.Contains(err.Error(), "--url is required") {
		t.Errorf("expected '--url is required' error, got: %v", err)
	}
}

func TestWebhooksUpdate(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/accounts/123/webhooks/1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]map[string]interface{}
		json.Unmarshal(body, &reqBody)

		webhook := reqBody["webhook"]
		if webhook["active"] != false {
			t.Errorf("expected active=false, got %v", webhook["active"])
		}
		if _, ok := webhook["url"]; ok {
			t.Errorf("did not expect url to be set, got %v", webhook["url"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"id":             1,
				"url":            "https://example.com/mailtrap/webhooks",
				"active":         false,
				"webhook_type":   "email_sending",
				"payload_format": "json",
			},
		})
	})
	defer cleanup()

	cmd := webhooks.NewCmdWebhooks(f)
	cmd.SetArgs([]string{"update", "--id", "1", "--active=false"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "email_sending") {
		t.Errorf("expected output to contain 'email_sending', got:\n%s", output)
	}
	if !strings.Contains(output, "false") {
		t.Errorf("expected output to contain 'false' for active, got:\n%s", output)
	}
}

func TestWebhooksDelete(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/accounts/123/webhooks/1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
	defer cleanup()

	cmd := webhooks.NewCmdWebhooks(f)
	cmd.SetArgs([]string{"delete", "--id", "1"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "Webhook deleted successfully") {
		t.Errorf("expected success message, got:\n%s", buf.String())
	}
}
