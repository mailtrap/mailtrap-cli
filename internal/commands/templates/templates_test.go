package templates_test

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
	"github.com/mailtrap/mailtrap-cli/internal/commands/templates"
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

func TestTemplatesList(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/accounts/123/email_templates" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Api-Token") != "test-token" {
			t.Errorf("expected Api-Token header 'test-token', got %q", r.Header.Get("Api-Token"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": 1, "uuid": "abc-123", "name": "Welcome", "subject": "Hello", "category": "transactional", "created_at": "2024-01-01"},
		})
	})
	defer cleanup()

	cmd := templates.NewCmdTemplates(f)
	cmd.SetArgs([]string{"list"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Welcome") {
		t.Errorf("expected output to contain 'Welcome', got:\n%s", output)
	}
}

func TestTemplatesListJSON(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": 1, "uuid": "abc-123", "name": "Welcome", "subject": "Hello", "category": "transactional", "created_at": "2024-01-01"},
		})
	})
	defer cleanup()

	viper.Set("output", "json")

	cmd := templates.NewCmdTemplates(f)
	cmd.SetArgs([]string{"list"})
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
		t.Fatalf("expected 1 template, got %d", len(result))
	}
	if result[0]["name"] != "Welcome" {
		t.Errorf("expected name 'Welcome', got %v", result[0]["name"])
	}
}

func TestTemplatesGet(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/accounts/123/email_templates/1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 1, "uuid": "abc-123", "name": "Welcome", "subject": "Hello", "category": "transactional", "created_at": "2024-01-01",
		})
	})
	defer cleanup()

	cmd := templates.NewCmdTemplates(f)
	cmd.SetArgs([]string{"get", "--id", "1"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Welcome") {
		t.Errorf("expected output to contain 'Welcome', got:\n%s", output)
	}
}

func TestTemplatesCreate(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/accounts/123/email_templates" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal body: %v", err)
		}

		tmpl, ok := payload["email_template"].(map[string]interface{})
		if !ok {
			t.Fatal("expected 'email_template' key in body")
		}
		if tmpl["name"] != "New" {
			t.Errorf("expected name 'New', got %v", tmpl["name"])
		}
		if tmpl["subject"] != "Hello {{name}}" {
			t.Errorf("expected subject 'Hello {{name}}', got %v", tmpl["subject"])
		}
		if tmpl["body_html"] != "<h1>Hi</h1>" {
			t.Errorf("expected body_html '<h1>Hi</h1>', got %v", tmpl["body_html"])
		}
		if tmpl["body_text"] != "" {
			t.Errorf("expected body_text '', got %v", tmpl["body_text"])
		}
		if tmpl["category"] != "" {
			t.Errorf("expected category '', got %v", tmpl["category"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 2, "uuid": "def-456", "name": "New", "subject": "Hello {{name}}", "category": "", "created_at": "2024-01-01",
		})
	})
	defer cleanup()

	cmd := templates.NewCmdTemplates(f)
	cmd.SetArgs([]string{"create", "--name", "New", "--subject", "Hello {{name}}", "--body-html", "<h1>Hi</h1>"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "New") {
		t.Errorf("expected output to contain 'New', got:\n%s", output)
	}
}

func TestTemplatesCreateMissingRequired(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	cmd := templates.NewCmdTemplates(f)
	cmd.SetArgs([]string{"create"})
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "required flag") {
		t.Errorf("expected error about required flags, got: %v", err)
	}
}

func TestTemplatesUpdate(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/accounts/123/email_templates/1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal body: %v", err)
		}

		tmpl, ok := payload["email_template"].(map[string]interface{})
		if !ok {
			t.Fatal("expected 'email_template' key in body")
		}
		if tmpl["name"] != "Updated" {
			t.Errorf("expected name 'Updated', got %v", tmpl["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 1, "uuid": "abc-123", "name": "Updated", "subject": "Hello", "category": "transactional", "created_at": "2024-01-01",
		})
	})
	defer cleanup()

	cmd := templates.NewCmdTemplates(f)
	cmd.SetArgs([]string{"update", "--id", "1", "--name", "Updated"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Updated") {
		t.Errorf("expected output to contain 'Updated', got:\n%s", output)
	}
}

func TestTemplatesDelete(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/accounts/123/email_templates/1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
	defer cleanup()

	cmd := templates.NewCmdTemplates(f)
	cmd.SetArgs([]string{"delete", "--id", "1"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "deleted successfully") {
		t.Errorf("expected output to contain 'deleted successfully', got:\n%s", output)
	}
}
