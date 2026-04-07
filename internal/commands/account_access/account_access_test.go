package account_access_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/commands/account_access"
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

func TestAccountAccessList(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/account_accesses") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{
				"id":             1,
				"specifier_type": "User",
				"specifier": map[string]interface{}{
					"id":    10,
					"email": "user@test.com",
					"name":  "Test User",
				},
				"resources": []map[string]interface{}{
					{
						"resource_id":   123,
						"resource_type": "account",
						"access_level":  100,
					},
				},
			},
		})
	})
	defer cleanup()

	cmd := account_access.NewCmdAccountAccess(f)
	cmd.SetArgs([]string{"list"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "user@test.com") {
		t.Errorf("expected output to contain 'user@test.com', got:\n%s", output)
	}
	if !strings.Contains(output, "Test User") {
		t.Errorf("expected output to contain 'Test User', got:\n%s", output)
	}
	if !strings.Contains(output, "account") {
		t.Errorf("expected output to contain 'account', got:\n%s", output)
	}
}

func TestAccountAccessListJSON(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{
				"id":             1,
				"specifier_type": "User",
				"specifier": map[string]interface{}{
					"id":    10,
					"email": "user@test.com",
					"name":  "Test User",
				},
				"resources": []map[string]interface{}{
					{
						"resource_id":   123,
						"resource_type": "account",
						"access_level":  100,
					},
				},
			},
		})
	})
	defer cleanup()

	viper.Set("output", "json")

	cmd := account_access.NewCmdAccountAccess(f)
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
		t.Fatalf("expected 1 access, got %d", len(result))
	}
	// Verify nested structure is preserved in JSON mode
	specifier, ok := result[0]["specifier"].(map[string]interface{})
	if !ok {
		t.Fatal("expected 'specifier' to be a nested object")
	}
	if specifier["email"] != "user@test.com" {
		t.Errorf("expected specifier email 'user@test.com', got %v", specifier["email"])
	}
}

func TestAccountAccessRemove(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/account_accesses/1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
	defer cleanup()

	cmd := account_access.NewCmdAccountAccess(f)
	cmd.SetArgs([]string{"remove", "--id", "1"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "removed successfully") {
		t.Errorf("expected 'removed successfully' in output, got:\n%s", output)
	}
}

func TestAccountAccessRemoveMissingID(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	buf := &bytes.Buffer{}
	f.IOStreams.Out = buf

	cmd := account_access.NewCmdAccountAccess(f)
	cmd.SetArgs([]string{"remove"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when --id is missing")
	}
	if !strings.Contains(err.Error(), "--id is required") {
		t.Errorf("expected '--id is required' error, got: %v", err)
	}
}
