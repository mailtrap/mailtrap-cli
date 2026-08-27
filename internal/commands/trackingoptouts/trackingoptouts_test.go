package trackingoptouts_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/commands/trackingoptouts"
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

func TestTrackingOptOutsList(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/tracking_opt_outs" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Api-Token") != "test-token" {
			t.Errorf("expected Api-Token header 'test-token', got %q", r.Header.Get("Api-Token"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "uuid-1", "email": "test@example.com", "domain_name": "example.com", "created_at": "2024-01-01"},
			},
			"last_id": nil,
		})
	})
	defer cleanup()

	cmd := trackingoptouts.NewCmdTrackingOptOuts(f)
	cmd.SetArgs([]string{"list"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "test@example.com") {
		t.Errorf("expected output to contain 'test@example.com', got:\n%s", output)
	}
	if strings.Contains(output, "Next page") {
		t.Errorf("expected no pagination hint without a cursor, got:\n%s", output)
	}
}

func TestTrackingOptOutsListJSON(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "uuid-1", "email": "test@example.com"},
			},
			"last_id": "uuid-1",
		})
	})
	defer cleanup()

	viper.Set("output", "json")

	cmd := trackingoptouts.NewCmdTrackingOptOuts(f)
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
		t.Fatalf("expected 1 tracking opt-out, got %d", len(result))
	}
	if result[0]["id"] != "uuid-1" {
		t.Errorf("expected id 'uuid-1', got %v", result[0]["id"])
	}
}

func TestTrackingOptOutsListWithFilters(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if got := query.Get("email"); got != "test@example.com" {
			t.Errorf("expected email 'test@example.com', got %q", got)
		}
		if got := query.Get("start_time"); got != "2024-01-01" {
			t.Errorf("expected start_time '2024-01-01', got %q", got)
		}
		if got := query.Get("end_time"); got != "2024-01-31" {
			t.Errorf("expected end_time '2024-01-31', got %q", got)
		}
		if got := query.Get("last_id"); got != "uuid-1" {
			t.Errorf("expected last_id 'uuid-1', got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "uuid-2", "email": "test@example.com"},
			},
			"last_id": "uuid-2",
		})
	})
	defer cleanup()

	cmd := trackingoptouts.NewCmdTrackingOptOuts(f)
	cmd.SetArgs([]string{
		"list",
		"--email", "test@example.com",
		"--start-time", "2024-01-01",
		"--end-time", "2024-01-31",
		"--last-id", "uuid-1",
	})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "--last-id uuid-2") {
		t.Errorf("expected output to hint the next page cursor, got:\n%s", buf.String())
	}
}

func TestTrackingOptOutsCreate(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/tracking_opt_outs" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("could not decode request body: %v", err)
		}
		if body["email"] != "test@example.com" {
			t.Errorf("expected email 'test@example.com', got %v", body["email"])
		}
		if body["domain_id"] != float64(4321) {
			t.Errorf("expected domain_id 4321, got %v", body["domain_id"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"id":          "uuid-1",
				"email":       "test@example.com",
				"domain_name": "example.com",
				"created_at":  "2024-01-01",
			},
		})
	})
	defer cleanup()

	cmd := trackingoptouts.NewCmdTrackingOptOuts(f)
	cmd.SetArgs([]string{
		"create",
		"--email", "test@example.com",
		"--domain-id", "4321",
	})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "test@example.com") {
		t.Errorf("expected output to contain 'test@example.com', got:\n%s", buf.String())
	}
}

func TestTrackingOptOutsCreateMissingFlags(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "missing email",
			args: []string{"create", "--domain-id", "4321"},
			want: "--email is required",
		},
		{
			name: "missing domain id",
			args: []string{"create", "--email", "test@example.com"},
			want: "--domain-id is required",
		},
		{
			name: "non-positive domain id",
			args: []string{"create", "--email", "test@example.com", "--domain-id", "0"},
			want: "--domain-id must be greater than 0",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
				t.Error("expected no request to be sent")
			})
			defer cleanup()

			cmd := trackingoptouts.NewCmdTrackingOptOuts(f)
			cmd.SetArgs(tc.args)

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

func TestTrackingOptOutsDelete(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/tracking_opt_outs/uuid-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	})
	defer cleanup()

	cmd := trackingoptouts.NewCmdTrackingOptOuts(f)
	cmd.SetArgs([]string{"delete", "--id", "uuid-1"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "deleted successfully") {
		t.Errorf("expected output to contain 'deleted successfully', got:\n%s", buf.String())
	}
}

func TestTrackingOptOutsDeleteMissingID(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		t.Error("expected no request to be sent")
	})
	defer cleanup()

	cmd := trackingoptouts.NewCmdTrackingOptOuts(f)
	cmd.SetArgs([]string{"delete"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "--id is required") {
		t.Errorf("expected error to contain '--id is required', got: %v", err)
	}
}
