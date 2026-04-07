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
