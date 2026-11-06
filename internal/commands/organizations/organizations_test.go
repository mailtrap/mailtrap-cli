package organizations_test

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
	"github.com/mailtrap/mailtrap-cli/internal/commands/organizations"
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

func TestListSubAccounts(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/organizations/456/sub_accounts") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": 1, "name": "Sub Account 1", "created_at": "2024-01-01"},
		})
	})
	defer cleanup()

	cmd := organizations.NewCmdOrganizations(f)
	cmd.SetArgs([]string{"list-sub-accounts", "--org-id", "456"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Sub Account 1") {
		t.Errorf("expected output to contain 'Sub Account 1', got:\n%s", output)
	}
}

func TestListSubAccountsMissingOrgID(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	buf := &bytes.Buffer{}
	f.IOStreams.Out = buf

	cmd := organizations.NewCmdOrganizations(f)
	cmd.SetArgs([]string{"list-sub-accounts"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when --org-id is missing")
	}
	if !strings.Contains(err.Error(), "--org-id is required") {
		t.Errorf("expected '--org-id is required' error, got: %v", err)
	}
}

func TestCreateSubAccount(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/organizations/456/sub_accounts") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]map[string]string
		json.Unmarshal(body, &reqBody)

		if reqBody["account"]["name"] != "New Sub" {
			t.Errorf("expected account name 'New Sub', got %q", reqBody["account"]["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 2, "name": "New Sub", "created_at": "2024-01-01",
		})
	})
	defer cleanup()

	cmd := organizations.NewCmdOrganizations(f)
	cmd.SetArgs([]string{"create-sub-account", "--org-id", "456", "--name", "New Sub"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "New Sub") {
		t.Errorf("expected output to contain 'New Sub', got:\n%s", output)
	}
}

func TestCreateSubAccountMissingFlags(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	buf := &bytes.Buffer{}
	f.IOStreams.Out = buf

	cmd := organizations.NewCmdOrganizations(f)
	cmd.SetArgs([]string{"create-sub-account"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when required flags are missing")
	}
}

func TestDeleteSubAccount(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/organizations/456/sub_accounts/789") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	})
	defer cleanup()

	cmd := organizations.NewCmdOrganizations(f)
	cmd.SetArgs([]string{"delete-sub-account", "--org-id", "456", "--sub-account-id", "789"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Sub-account 789 deleted successfully") {
		t.Errorf("expected output to contain 'Sub-account 789 deleted successfully', got:\n%s", output)
	}
}

func TestDeleteSubAccountAPIError(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Not Found"})
	})
	defer cleanup()

	cmd := organizations.NewCmdOrganizations(f)
	cmd.SetArgs([]string{"delete-sub-account", "--org-id", "456", "--sub-account-id", "789"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "API error 404") {
		t.Errorf("expected error to contain 'API error 404', got: %v", err)
	}
	if strings.Contains(buf.String(), "deleted successfully") {
		t.Errorf("expected no success message, got:\n%s", buf.String())
	}
}

func TestDeleteSubAccountInvalidFlags(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "missing org id",
			args: []string{"delete-sub-account", "--sub-account-id", "789"},
			want: "--org-id is required",
		},
		{
			name: "missing sub-account id",
			args: []string{"delete-sub-account", "--org-id", "456"},
			want: "--sub-account-id is required",
		},
		{
			name: "list of sub-account ids",
			args: []string{"delete-sub-account", "--org-id", "456", "--sub-account-id", "12347,12348"},
			want: `invalid argument "12347,12348" for "--sub-account-id"`,
		},
		{
			name: "non-numeric org id",
			args: []string{"delete-sub-account", "--org-id", "abc", "--sub-account-id", "789"},
			want: `invalid argument "abc" for "--org-id"`,
		},
		{
			name: "non-positive sub-account id",
			args: []string{"delete-sub-account", "--org-id", "456", "--sub-account-id", "0"},
			want: "--sub-account-id must be greater than 0",
		},
		{
			name: "non-positive org id",
			args: []string{"delete-sub-account", "--org-id", "0", "--sub-account-id", "789"},
			want: "--org-id must be greater than 0",
		},
		{
			name: "extra positional argument",
			args: []string{"delete-sub-account", "--org-id", "456", "--sub-account-id", "789", "790"},
			want: `unknown command "790"`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
				t.Error("expected no request to be sent")
			})
			defer cleanup()

			cmd := organizations.NewCmdOrganizations(f)
			cmd.SetArgs(tc.args)
			cmd.SetOut(buf)

			err := cmd.Execute()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("expected error to contain %q, got: %v", tc.want, err)
			}
		})
	}
}
