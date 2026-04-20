package contact_fields_test

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
	"github.com/mailtrap/mailtrap-cli/internal/commands/contact_fields"
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

func TestContactFieldsList(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/accounts/123/contacts/fields" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Api-Token") != "test-token" {
			t.Errorf("expected Api-Token header 'test-token', got %q", r.Header.Get("Api-Token"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": 1, "name": "Company", "data_type": "text", "merge_tag": "{{company}}"},
		})
	})
	defer cleanup()

	cmd := contact_fields.NewCmdContactFields(f)
	cmd.SetArgs([]string{"list"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Company") {
		t.Errorf("expected output to contain 'Company', got:\n%s", output)
	}
	if !strings.Contains(output, "text") {
		t.Errorf("expected output to contain 'text', got:\n%s", output)
	}
}

func TestContactFieldsListJSON(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": 1, "name": "Company", "data_type": "text", "merge_tag": "{{company}}"},
		})
	})
	defer cleanup()

	viper.Set("output", "json")

	cmd := contact_fields.NewCmdContactFields(f)
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
		t.Fatalf("expected 1 contact field, got %d", len(result))
	}
	if result[0]["name"] != "Company" {
		t.Errorf("expected name 'Company', got %v", result[0]["name"])
	}
}

func TestContactFieldsGet(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/contacts/fields/1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 1, "name": "Company", "data_type": "text", "merge_tag": "{{company}}",
		})
	})
	defer cleanup()

	cmd := contact_fields.NewCmdContactFields(f)
	cmd.SetArgs([]string{"get", "--id", "1"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Company") {
		t.Errorf("expected output to contain 'Company', got:\n%s", output)
	}
}

func TestContactFieldsGetMissingID(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	cmd := contact_fields.NewCmdContactFields(f)
	cmd.SetArgs([]string{"get"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "--id is required") {
		t.Errorf("expected error to contain '--id is required', got: %v", err)
	}
}

func TestContactFieldsCreate(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/contacts/fields") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to parse request body: %v", err)
		}
		if payload["name"] != "Company" {
			t.Errorf("expected name 'Company', got %v", payload["name"])
		}
		if payload["data_type"] != "text" {
			t.Errorf("expected data_type 'text', got %v", payload["data_type"])
		}
		if payload["merge_tag"] != "{{company}}" {
			t.Errorf("expected merge_tag '{{company}}', got %v", payload["merge_tag"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 1, "name": "Company", "data_type": "text", "merge_tag": "{{company}}",
		})
	})
	defer cleanup()

	cmd := contact_fields.NewCmdContactFields(f)
	cmd.SetArgs([]string{"create", "--name", "Company", "--data-type", "text", "--merge-tag", "{{company}}"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Company") {
		t.Errorf("expected output to contain 'Company', got:\n%s", output)
	}
}

func TestContactFieldsCreateMissingName(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	cmd := contact_fields.NewCmdContactFields(f)
	cmd.SetArgs([]string{"create", "--data-type", "text"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "--name is required") {
		t.Errorf("expected error to contain '--name is required', got: %v", err)
	}
}

func TestContactFieldsCreateMissingDataType(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	cmd := contact_fields.NewCmdContactFields(f)
	cmd.SetArgs([]string{"create", "--name", "Company"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "--data-type is required") {
		t.Errorf("expected error to contain '--data-type is required', got: %v", err)
	}
}

func TestContactFieldsCreateMissingMergeTag(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	cmd := contact_fields.NewCmdContactFields(f)
	cmd.SetArgs([]string{"create", "--name", "Company", "--data-type", "text"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "--merge-tag is required") {
		t.Errorf("expected error to contain '--merge-tag is required', got: %v", err)
	}
}

func TestContactFieldsUpdate(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/contacts/fields/1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to parse request body: %v", err)
		}
		if _, ok := payload["name"]; !ok {
			t.Error("expected body to contain 'name' key")
		}
		if payload["name"] != "Updated" {
			t.Errorf("expected name 'Updated', got %v", payload["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 1, "name": "Updated", "data_type": "text", "merge_tag": "{{company}}",
		})
	})
	defer cleanup()

	cmd := contact_fields.NewCmdContactFields(f)
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

func TestContactFieldsUpdateWithMergeTag(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to parse request body: %v", err)
		}
		if _, ok := payload["merge_tag"]; !ok {
			t.Error("expected body to contain 'merge_tag' key")
		}
		if payload["merge_tag"] != "{{new_tag}}" {
			t.Errorf("expected merge_tag '{{new_tag}}', got %v", payload["merge_tag"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 1, "name": "Company", "data_type": "text", "merge_tag": "{{new_tag}}",
		})
	})
	defer cleanup()

	cmd := contact_fields.NewCmdContactFields(f)
	cmd.SetArgs([]string{"update", "--id", "1", "--merge-tag", "{{new_tag}}"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "{{new_tag}}") {
		t.Errorf("expected output to contain '{{new_tag}}', got:\n%s", output)
	}
}

func TestContactFieldsDelete(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/contacts/fields/1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	})
	defer cleanup()

	cmd := contact_fields.NewCmdContactFields(f)
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
