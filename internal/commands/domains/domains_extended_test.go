package domains_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/mailtrap/mailtrap-cli/internal/commands/domains"
)

func TestDomainsSendSetupInstructions(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/sending_domains/1/send_setup_instructions") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]string
		json.Unmarshal(body, &reqBody)

		if reqBody["email"] != "admin@test.com" {
			t.Errorf("expected email 'admin@test.com', got %q", reqBody["email"])
		}

		w.WriteHeader(http.StatusOK)
	})
	defer cleanup()

	cmd := domains.NewCmdDomains(f)
	cmd.SetArgs([]string{"send-setup-instructions", "--id", "1", "--email", "admin@test.com"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Setup instructions sent successfully") {
		t.Errorf("expected success message, got:\n%s", output)
	}
}

func TestDomainsSendSetupInstructionsMissingFlags(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	buf := &bytes.Buffer{}
	f.IOStreams.Out = buf

	cmd := domains.NewCmdDomains(f)
	cmd.SetArgs([]string{"send-setup-instructions"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when required flags are missing")
	}
	if !strings.Contains(err.Error(), "--id is required") {
		t.Errorf("expected '--id is required' error, got: %v", err)
	}
}

func TestDomainsUpdate(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/accounts/123/sending_domains/1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]map[string]interface{}
		json.Unmarshal(body, &reqBody)

		fields := reqBody["sending_domain"]
		if len(fields) != 2 {
			t.Errorf("expected only the changed flags in the body, got %v", fields)
		}
		if fields["tracking_opt_out_enabled"] != true {
			t.Errorf("expected tracking_opt_out_enabled true, got %v", fields["tracking_opt_out_enabled"])
		}
		if fields["auto_unsubscribe_link_enabled"] != false {
			t.Errorf("expected auto_unsubscribe_link_enabled false, got %v", fields["auto_unsubscribe_link_enabled"])
		}
		if _, ok := fields["open_tracking_enabled"]; ok {
			t.Error("expected unset open-tracking flag to be omitted from the body")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 1, "domain_name": "example.com", "dns_verified": true, "compliance_status": "compliant",
			"tracking_opt_out_enabled": true, "auto_unsubscribe_link_enabled": false,
		})
	})
	defer cleanup()

	cmd := domains.NewCmdDomains(f)
	cmd.SetArgs([]string{"update", "--id", "1", "--tracking-opt-out=true", "--auto-unsubscribe-link=false"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "TRACKING OPT OUT") {
		t.Errorf("expected output to contain the tracking opt out column, got:\n%s", output)
	}
}

func TestDomainsUpdateMissingID(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		t.Error("expected no request to be sent")
	})
	defer cleanup()

	cmd := domains.NewCmdDomains(f)
	cmd.SetArgs([]string{"update", "--tracking-opt-out=true"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing --id")
	}
	if !strings.Contains(err.Error(), "--id is required") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDomainsUpdateNoFields(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		t.Error("expected no request to be sent")
	})
	defer cleanup()

	cmd := domains.NewCmdDomains(f)
	cmd.SetArgs([]string{"update", "--id", "1"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when no attribute flag is passed")
	}
	if !strings.Contains(err.Error(), "at least one attribute flag is required") {
		t.Errorf("unexpected error: %v", err)
	}
}
