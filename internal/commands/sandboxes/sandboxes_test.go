package sandboxes_test

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
	"github.com/mailtrap/mailtrap-cli/internal/commands/sandboxes"
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

func TestSandboxesList(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/inboxes") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": 1, "name": "Test Sandbox", "email_username": "abc123", "status": "active", "max_size": 50},
		})
	})
	defer cleanup()

	cmd := sandboxes.NewCmdSandboxes(f)
	cmd.SetArgs([]string{"list"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Test Sandbox") {
		t.Errorf("expected output to contain 'Test Sandbox', got:\n%s", output)
	}
	if !strings.Contains(output, "abc123") {
		t.Errorf("expected output to contain 'abc123', got:\n%s", output)
	}
	if !strings.Contains(output, "active") {
		t.Errorf("expected output to contain 'active', got:\n%s", output)
	}
}

func TestSandboxesListJSON(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": 1, "name": "Test Sandbox", "email_username": "abc123", "status": "active", "max_size": 50},
		})
	})
	defer cleanup()

	viper.Set("output", "json")

	cmd := sandboxes.NewCmdSandboxes(f)
	cmd.SetArgs([]string{"list"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	var result []map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput:\n%s", err, output)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 sandbox, got %d", len(result))
	}
	if result[0]["name"] != "Test Sandbox" {
		t.Errorf("expected name 'Test Sandbox', got %v", result[0]["name"])
	}
}

func TestSandboxesGet(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/inboxes/1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 1, "name": "Test Sandbox", "email_username": "abc123", "status": "active", "max_size": 50,
		})
	})
	defer cleanup()

	cmd := sandboxes.NewCmdSandboxes(f)
	cmd.SetArgs([]string{"get", "--id", "1"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Test Sandbox") {
		t.Errorf("expected output to contain 'Test Sandbox', got:\n%s", output)
	}
}

func TestSandboxesGetMissingID(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	buf := &bytes.Buffer{}
	f.IOStreams.Out = buf

	cmd := sandboxes.NewCmdSandboxes(f)
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

func TestSandboxesCreate(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/projects/1/inboxes") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]interface{}
		json.Unmarshal(body, &payload)
		inbox, ok := payload["inbox"].(map[string]interface{})
		if !ok {
			t.Errorf("expected 'inbox' key in body, got: %s", string(body))
		}
		if inbox["name"] != "New Sandbox" {
			t.Errorf("expected name 'New Sandbox', got: %v", inbox["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 2, "name": "New Sandbox", "email_username": "def456", "status": "active", "max_size": 50,
		})
	})
	defer cleanup()

	cmd := sandboxes.NewCmdSandboxes(f)
	cmd.SetArgs([]string{"create", "--project-id", "1", "--name", "New Sandbox"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "New Sandbox") {
		t.Errorf("expected output to contain 'New Sandbox', got:\n%s", output)
	}
	if !strings.Contains(output, "def456") {
		t.Errorf("expected output to contain 'def456', got:\n%s", output)
	}
}

func TestSandboxesCreateMissingFlags(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	buf := &bytes.Buffer{}
	f.IOStreams.Out = buf

	cmd := sandboxes.NewCmdSandboxes(f)
	cmd.SetArgs([]string{"create", "--name", "Test"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when --project-id is missing")
	}
	if !strings.Contains(err.Error(), "--project-id is required") {
		t.Errorf("expected '--project-id is required' error, got: %v", err)
	}
}

func TestSandboxesUpdate(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/inboxes/1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]interface{}
		json.Unmarshal(body, &payload)
		inbox, ok := payload["inbox"].(map[string]interface{})
		if !ok {
			t.Errorf("expected 'inbox' key in body, got: %s", string(body))
		}
		if inbox["name"] != "Updated" {
			t.Errorf("expected name 'Updated', got: %v", inbox["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 1, "name": "Updated", "email_username": "abc123", "status": "active", "max_size": 50,
		})
	})
	defer cleanup()

	cmd := sandboxes.NewCmdSandboxes(f)
	cmd.SetArgs([]string{"update", "--id", "1", "--name", "Updated"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Updated") {
		t.Errorf("expected output to contain 'Updated', got:\n%s", output)
	}
}

func TestSandboxesDelete(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/inboxes/1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
	defer cleanup()

	cmd := sandboxes.NewCmdSandboxes(f)
	cmd.SetArgs([]string{"delete", "--id", "1"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "deleted successfully") {
		t.Errorf("expected 'deleted successfully' in output, got:\n%s", output)
	}
}

func TestSandboxesClean(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/inboxes/1/clean") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
	defer cleanup()

	cmd := sandboxes.NewCmdSandboxes(f)
	cmd.SetArgs([]string{"clean", "--id", "1"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "cleaned successfully") {
		t.Errorf("expected 'cleaned successfully' in output, got:\n%s", output)
	}
}

func TestSandboxesMarkRead(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/inboxes/1/all_read") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 1, "name": "Test Sandbox", "email_username": "abc123", "status": "active", "max_size": 50,
		})
	})
	defer cleanup()

	cmd := sandboxes.NewCmdSandboxes(f)
	cmd.SetArgs([]string{"mark-read", "--id", "1"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Test Sandbox") {
		t.Errorf("expected output to contain 'Test Sandbox', got:\n%s", output)
	}
}

func TestSandboxesResetCredentials(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/inboxes/1/reset_credentials") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 1, "name": "Test Sandbox", "email_username": "abc123", "status": "active", "max_size": 50,
		})
	})
	defer cleanup()

	cmd := sandboxes.NewCmdSandboxes(f)
	cmd.SetArgs([]string{"reset-credentials", "--id", "1"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Test Sandbox") {
		t.Errorf("expected output to contain 'Test Sandbox', got:\n%s", output)
	}
}

func TestSandboxesToggleEmail(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/inboxes/1/toggle_email_username") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 1, "name": "Test Sandbox", "email_username": "abc123", "status": "active", "max_size": 50,
		})
	})
	defer cleanup()

	cmd := sandboxes.NewCmdSandboxes(f)
	cmd.SetArgs([]string{"toggle-email", "--id", "1"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Test Sandbox") {
		t.Errorf("expected output to contain 'Test Sandbox', got:\n%s", output)
	}
}

func TestSandboxesResetEmail(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/inboxes/1/reset_email_username") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 1, "name": "Test Sandbox", "email_username": "abc123", "status": "active", "max_size": 50,
		})
	})
	defer cleanup()

	cmd := sandboxes.NewCmdSandboxes(f)
	cmd.SetArgs([]string{"reset-email", "--id", "1"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Test Sandbox") {
		t.Errorf("expected output to contain 'Test Sandbox', got:\n%s", output)
	}
}
