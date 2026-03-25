package domains_test

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
	"github.com/mailtrap/mailtrap-cli/internal/commands/domains"
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

func TestDomainsList(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/accounts/123/sending_domains") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("expected Authorization header 'Bearer test-token', got %q", r.Header.Get("Authorization"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": 1, "domain_name": "example.com", "dns_verified": true, "compliance_status": "compliant"},
				{"id": 2, "domain_name": "test.com", "dns_verified": false, "compliance_status": "pending"},
			},
		})
	})
	defer cleanup()

	cmd := domains.NewCmdDomains(f)
	cmd.SetArgs([]string{"list"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "example.com") {
		t.Errorf("expected output to contain 'example.com', got:\n%s", output)
	}
	if !strings.Contains(output, "test.com") {
		t.Errorf("expected output to contain 'test.com', got:\n%s", output)
	}
	if !strings.Contains(output, "compliant") {
		t.Errorf("expected output to contain 'compliant', got:\n%s", output)
	}
	if !strings.Contains(output, "ID") {
		t.Errorf("expected output to contain header 'ID', got:\n%s", output)
	}
}

func TestDomainsListJSON(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": 1, "domain_name": "example.com", "dns_verified": true, "compliance_status": "compliant"},
			},
		})
	})
	defer cleanup()

	viper.Set("output", "json")

	cmd := domains.NewCmdDomains(f)
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
		t.Fatalf("expected 1 domain, got %d", len(result))
	}
	if result[0]["domain_name"] != "example.com" {
		t.Errorf("expected domain_name 'example.com', got %v", result[0]["domain_name"])
	}
}

func TestDomainsGet(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/accounts/123/sending_domains/1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 1, "domain_name": "example.com", "dns_verified": true, "compliance_status": "compliant",
		})
	})
	defer cleanup()

	cmd := domains.NewCmdDomains(f)
	cmd.SetArgs([]string{"get", "--id", "1"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "example.com") {
		t.Errorf("expected output to contain 'example.com', got:\n%s", output)
	}
	if !strings.Contains(output, "compliant") {
		t.Errorf("expected output to contain 'compliant', got:\n%s", output)
	}
}

func TestDomainsCreate(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/accounts/123/sending_domains") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]string
		json.Unmarshal(body, &reqBody)

		if reqBody["domain"] != "example.com" {
			t.Errorf("expected domain 'example.com' in request body, got %q", reqBody["domain"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 10, "domain_name": "example.com", "dns_verified": false, "compliance_status": "pending",
		})
	})
	defer cleanup()

	cmd := domains.NewCmdDomains(f)
	cmd.SetArgs([]string{"create", "--name", "example.com"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "example.com") {
		t.Errorf("expected output to contain 'example.com', got:\n%s", output)
	}
	if !strings.Contains(output, "pending") {
		t.Errorf("expected output to contain 'pending', got:\n%s", output)
	}
}

func TestDomainsDelete(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/accounts/123/sending_domains/1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	})
	defer cleanup()

	cmd := domains.NewCmdDomains(f)
	cmd.SetArgs([]string{"delete", "--id", "1"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Domain deleted successfully") {
		t.Errorf("expected success message, got:\n%s", output)
	}
}

func TestDomainsGetMissingID(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	buf := &bytes.Buffer{}
	f.IOStreams.Out = buf

	cmd := domains.NewCmdDomains(f)
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
