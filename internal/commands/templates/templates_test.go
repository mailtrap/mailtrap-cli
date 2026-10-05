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

func listBody() map[string]interface{} {
	return map[string]interface{}{
		"data": []map[string]interface{}{
			{"id": 1, "uuid": "abc-123", "name": "Welcome", "subject": "Hello", "category": "transactional", "created_at": "2024-01-01"},
		},
		"pagination": map[string]interface{}{"token": 1, "prev_token": nil, "next_token": 2},
	}
}

func TestTemplatesList(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/accounts/123/templates" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Api-Token") != "test-token" {
			t.Errorf("expected Api-Token header 'test-token', got %q", r.Header.Get("Api-Token"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(listBody())
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
		json.NewEncoder(w).Encode(listBody())
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
	var result struct {
		Data       []map[string]interface{} `json:"data"`
		Pagination map[string]interface{}   `json:"pagination"`
	}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("output is not a JSON object: %v\noutput:\n%s", err, output)
	}
	if len(result.Data) != 1 {
		t.Fatalf("expected 1 template, got %d", len(result.Data))
	}
	if result.Data[0]["name"] != "Welcome" {
		t.Errorf("expected name 'Welcome', got %v", result.Data[0]["name"])
	}
	if result.Pagination["next_token"] != float64(2) {
		t.Errorf("expected pagination.next_token 2, got %v", result.Pagination["next_token"])
	}
}

func TestTemplatesListNextPage(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(listBody())
	})
	defer cleanup()

	cmd := templates.NewCmdTemplates(f)
	cmd.SetArgs([]string{"list"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "Next page: --token 2") {
		t.Errorf("expected next page footer, got:\n%s", buf.String())
	}
}

func TestTemplatesListQuery(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("per_page"); got != "10" {
			t.Errorf("expected per_page=10, got %q", got)
		}
		if got := r.URL.Query().Get("token"); got != "2" {
			t.Errorf("expected token=2, got %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(listBody())
	})
	defer cleanup()

	cmd := templates.NewCmdTemplates(f)
	cmd.SetArgs([]string{"list", "--per-page", "10", "--token", "2"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTemplatesListOmitsUnsetQuery(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Errorf("expected no query, got %q", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(listBody())
	})
	defer cleanup()

	cmd := templates.NewCmdTemplates(f)
	cmd.SetArgs([]string{"list"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTemplatesGet(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/accounts/123/templates/1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{
			"id": 1, "uuid": "abc-123", "name": "Welcome", "subject": "Hello", "category": "transactional", "created_at": "2024-01-01",
		}})
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
		if r.URL.Path != "/api/accounts/123/templates" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal body: %v", err)
		}

		if _, ok := payload["email_template"]; ok {
			t.Error("expected a flat body without 'email_template' key")
		}
		if payload["name"] != "New" {
			t.Errorf("expected name 'New', got %v", payload["name"])
		}
		if payload["subject"] != "Hello {{name}}" {
			t.Errorf("expected subject 'Hello {{name}}', got %v", payload["subject"])
		}
		if payload["body_html"] != "<h1>Hi</h1>" {
			t.Errorf("expected body_html '<h1>Hi</h1>', got %v", payload["body_html"])
		}
		if payload["body_text"] != "" {
			t.Errorf("expected body_text '', got %v", payload["body_text"])
		}
		if payload["category"] != "General" {
			t.Errorf("expected category 'General', got %v", payload["category"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{
			"id": 2, "uuid": "def-456", "name": "New", "subject": "Hello {{name}}", "category": "General", "created_at": "2024-01-01",
		}})
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
		if !strings.HasSuffix(r.URL.Path, "/api/accounts/123/templates/1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal body: %v", err)
		}

		if _, ok := payload["email_template"]; ok {
			t.Error("expected a flat body without 'email_template' key")
		}
		if payload["name"] != "Updated" {
			t.Errorf("expected name 'Updated', got %v", payload["name"])
		}
		if len(payload) != 1 {
			t.Errorf("expected only changed flags in body, got %v", payload)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{
			"id": 1, "uuid": "abc-123", "name": "Updated", "subject": "Hello", "category": "transactional", "created_at": "2024-01-01",
		}})
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
		if !strings.HasSuffix(r.URL.Path, "/api/accounts/123/templates/1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
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

func TestTemplatesCreateKeepsExplicitCategory(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&payload)
		if payload["category"] != "Promo" {
			t.Errorf("expected category 'Promo', got %v", payload["category"])
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"data":{"id":2,"name":"New"}}`))
	})
	defer cleanup()

	cmd := templates.NewCmdTemplates(f)
	cmd.SetArgs([]string{"create", "--name", "New", "--subject", "Hi", "--category", "Promo"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTemplatesUpdateRequiresAttribute(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		t.Error("unexpected request")
	})
	defer cleanup()

	cmd := templates.NewCmdTemplates(f)
	cmd.SetArgs([]string{"update", "--id", "1"})
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "at least one attribute flag is required") {
		t.Fatalf("expected attribute flag error, got: %v", err)
	}
}
