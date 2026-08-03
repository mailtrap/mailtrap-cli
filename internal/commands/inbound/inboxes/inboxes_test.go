package inboxes_test

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
	"github.com/mailtrap/mailtrap-cli/internal/commands/inbound/inboxes"
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
			return &config.Config{APIToken: "test-token"}
		},
		IOStreams: &cmdutil.IOStreams{
			Out:    buf,
			ErrOut: &bytes.Buffer{},
		},
		ClientOverride: c,
	}

	viper.Set("api-token", "test-token")
	viper.Set("output", "table")

	return f, buf, func() {
		server.Close()
		viper.Reset()
	}
}

func TestInboxesList(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/inbound/folders/90/inboxes") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": 201, "name": "Support inbox", "address": "support@inbound-mailtrap.io"},
			{"id": 202, "name": "Catch-all", "address": "catch-all@example.com", "domain_id": 6},
		})
	})
	defer cleanup()

	cmd := inboxes.NewCmdInboxes(f)
	cmd.SetArgs([]string{"list", "--folder-id", "90"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "support@inbound-mailtrap.io") {
		t.Errorf("expected output to contain inbox address, got:\n%s", buf.String())
	}
}

func TestInboxesListMissingFolderID(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	cmd := inboxes.NewCmdInboxes(f)
	cmd.SetArgs([]string{"list"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when --folder-id is missing")
	}
	if !strings.Contains(err.Error(), "--folder-id is required") {
		t.Errorf("expected '--folder-id is required' error, got: %v", err)
	}
}

func TestInboxesGet(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/api/inbound/folders/90/inboxes/201") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 201, "name": "Support inbox", "address": "support@inbound-mailtrap.io",
		})
	})
	defer cleanup()

	cmd := inboxes.NewCmdInboxes(f)
	cmd.SetArgs([]string{"get", "--folder-id", "90", "--id", "201"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "Support inbox") {
		t.Errorf("expected output to contain 'Support inbox', got:\n%s", buf.String())
	}
}

func TestInboxesCreateWithDomainID(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/inbound/folders/90/inboxes") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]interface{}
		json.Unmarshal(body, &reqBody)
		if reqBody["name"] != "Custom domain inbox" {
			t.Errorf("unexpected name: %v", reqBody["name"])
		}
		if reqBody["domain_id"] != float64(6) {
			t.Errorf("expected domain_id 6, got %v", reqBody["domain_id"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 203, "name": "Custom domain inbox", "address": "catch-all@example.com", "domain_id": 6,
		})
	})
	defer cleanup()

	cmd := inboxes.NewCmdInboxes(f)
	cmd.SetArgs([]string{"create", "--folder-id", "90", "--name", "Custom domain inbox", "--domain-id", "6"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "Custom domain inbox") {
		t.Errorf("expected output to contain inbox name, got:\n%s", buf.String())
	}
}

func TestInboxesDelete(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/inbound/folders/90/inboxes/201") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer cleanup()

	cmd := inboxes.NewCmdInboxes(f)
	cmd.SetArgs([]string{"delete", "--folder-id", "90", "--id", "201"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "deleted successfully") {
		t.Errorf("expected success message, got:\n%s", buf.String())
	}
}

func TestInboxesListJSON(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": 201, "name": "Support inbox", "address": "support@inbound-mailtrap.io"},
		})
	})
	defer cleanup()

	viper.Set("output", "json")

	cmd := inboxes.NewCmdInboxes(f)
	cmd.SetArgs([]string{"list", "--folder-id", "90"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput:\n%s", err, buf.String())
	}
	if len(result) != 1 || result[0]["address"] != "support@inbound-mailtrap.io" {
		t.Errorf("unexpected JSON result: %v", result)
	}
}

func TestInboxesCreateMissingName(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	cmd := inboxes.NewCmdInboxes(f)
	cmd.SetArgs([]string{"create", "--folder-id", "90"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when --name is missing")
	}
	if !strings.Contains(err.Error(), "--name is required") {
		t.Errorf("expected '--name is required' error, got: %v", err)
	}
}

func TestInboxesUpdate(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/inbound/folders/90/inboxes/201") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]interface{}
		json.Unmarshal(body, &reqBody)
		if reqBody["name"] != "Renamed inbox" {
			t.Errorf("unexpected name: %v", reqBody["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 201, "name": "Renamed inbox", "address": "support@inbound-mailtrap.io",
		})
	})
	defer cleanup()

	cmd := inboxes.NewCmdInboxes(f)
	cmd.SetArgs([]string{"update", "--folder-id", "90", "--id", "201", "--name", "Renamed inbox"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "Renamed inbox") {
		t.Errorf("expected output to contain 'Renamed inbox', got:\n%s", buf.String())
	}
}
