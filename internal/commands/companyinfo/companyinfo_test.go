package companyinfo_test

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
	"github.com/mailtrap/mailtrap-cli/internal/commands/companyinfo"
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

func companyInfoPayload() map[string]interface{} {
	return map[string]interface{}{
		"data": map[string]interface{}{
			"name":        "Mailtrap",
			"address":     "123 Main St",
			"city":        "San Francisco",
			"country":     "US",
			"phone":       "+1-555-0100",
			"zip_code":    "94105",
			"website_url": "https://mailtrap.io",
			"info_level":  "business",
		},
	}
}

func TestCompanyInfoGet(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/domains/435/company_info" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Api-Token") != "test-token" {
			t.Errorf("expected Api-Token header 'test-token', got %q", r.Header.Get("Api-Token"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(companyInfoPayload())
	})
	defer cleanup()

	cmd := companyinfo.NewCmdCompanyInfo(f)
	cmd.SetArgs([]string{"get", "--domain-id", "435"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Mailtrap") {
		t.Errorf("expected output to contain 'Mailtrap', got:\n%s", output)
	}
	if !strings.Contains(output, "business") {
		t.Errorf("expected output to contain 'business', got:\n%s", output)
	}
}

func TestCompanyInfoGetJSON(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(companyInfoPayload())
	})
	defer cleanup()

	viper.Set("output", "json")

	cmd := companyinfo.NewCmdCompanyInfo(f)
	cmd.SetArgs([]string{"get", "--domain-id", "435"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}
	if result["name"] != "Mailtrap" {
		t.Errorf("expected name 'Mailtrap', got %v", result["name"])
	}
	if result["zip_code"] != "94105" {
		t.Errorf("expected zip_code '94105', got %v", result["zip_code"])
	}
}

func TestCompanyInfoGetMissingDomainID(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		t.Error("expected no request to be sent")
	})
	defer cleanup()

	cmd := companyinfo.NewCmdCompanyInfo(f)
	cmd.SetArgs([]string{"get"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing --domain-id")
	}
	if !strings.Contains(err.Error(), "--domain-id is required") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCompanyInfoCreate(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/domains/435/company_info" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]map[string]interface{}
		json.Unmarshal(body, &reqBody)

		fields := reqBody["company_info"]
		if fields["name"] != "Mailtrap" {
			t.Errorf("expected name 'Mailtrap', got %v", fields["name"])
		}
		if fields["website_url"] != "https://mailtrap.io" {
			t.Errorf("expected website_url 'https://mailtrap.io', got %v", fields["website_url"])
		}
		if _, ok := fields["phone"]; ok {
			t.Error("expected unset phone flag to be omitted from the body")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(companyInfoPayload())
	})
	defer cleanup()

	cmd := companyinfo.NewCmdCompanyInfo(f)
	cmd.SetArgs([]string{
		"create", "--domain-id", "435",
		"--name", "Mailtrap",
		"--address", "123 Main St",
		"--city", "San Francisco",
		"--country", "US",
		"--zip-code", "94105",
		"--website-url", "https://mailtrap.io",
	})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "Mailtrap") {
		t.Errorf("expected output to contain 'Mailtrap', got:\n%s", buf.String())
	}
}

func TestCompanyInfoCreateMissingRequiredFlag(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		t.Error("expected no request to be sent")
	})
	defer cleanup()

	cmd := companyinfo.NewCmdCompanyInfo(f)
	cmd.SetArgs([]string{"create", "--domain-id", "435", "--name", "Mailtrap"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing --address")
	}
	if !strings.Contains(err.Error(), "--address is required") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCompanyInfoUpdate(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/domains/435/company_info" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]map[string]interface{}
		json.Unmarshal(body, &reqBody)

		fields := reqBody["company_info"]
		if len(fields) != 2 {
			t.Errorf("expected only the changed flags in the body, got %v", fields)
		}
		if fields["city"] != "New York" {
			t.Errorf("expected city 'New York', got %v", fields["city"])
		}
		if fields["zip_code"] != "10001" {
			t.Errorf("expected zip_code '10001', got %v", fields["zip_code"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(companyInfoPayload())
	})
	defer cleanup()

	cmd := companyinfo.NewCmdCompanyInfo(f)
	cmd.SetArgs([]string{"update", "--domain-id", "435", "--city", "New York", "--zip-code", "10001"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCompanyInfoUpdateNoFields(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		t.Error("expected no request to be sent")
	})
	defer cleanup()

	cmd := companyinfo.NewCmdCompanyInfo(f)
	cmd.SetArgs([]string{"update", "--domain-id", "435"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when no attribute flag is passed")
	}
	if !strings.Contains(err.Error(), "at least one attribute flag is required") {
		t.Errorf("unexpected error: %v", err)
	}
}
