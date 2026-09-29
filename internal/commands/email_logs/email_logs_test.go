package emaillogs_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	emaillogs "github.com/mailtrap/mailtrap-cli/internal/commands/email_logs"
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

func TestEmailLogsList(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/email_logs") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"messages": []map[string]interface{}{
				{"message_id": "msg-1", "subject": "Test", "from": "a@b.com", "to": "c@d.com", "status": "delivered", "sent_at": "2024-01-01"},
			},
			"total_count":      1,
			"next_page_cursor": "cursor-abc",
		})
	})
	defer cleanup()

	cmd := emaillogs.NewCmdEmailLogs(f)
	cmd.SetArgs([]string{"list"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Test") {
		t.Errorf("expected output to contain 'Test', got:\n%s", output)
	}
	if !strings.Contains(output, "delivered") {
		t.Errorf("expected output to contain 'delivered', got:\n%s", output)
	}
	if !strings.Contains(output, "--cursor cursor-abc") {
		t.Errorf("expected output to surface the next-page cursor, got:\n%s", output)
	}
}

func TestEmailLogsListJSON(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"messages": []map[string]interface{}{
				{"message_id": "msg-1", "subject": "Test", "from": "a@b.com", "to": "c@d.com", "status": "delivered", "sent_at": "2024-01-01"},
			},
			"total_count":      1,
			"next_page_cursor": "cursor-abc",
		})
	})
	defer cleanup()

	viper.Set("output", "json")

	cmd := emaillogs.NewCmdEmailLogs(f)
	cmd.SetArgs([]string{"list"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	var result struct {
		Messages       []map[string]interface{} `json:"messages"`
		TotalCount     int                      `json:"total_count"`
		NextPageCursor string                   `json:"next_page_cursor"`
	}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput:\n%s", err, output)
	}
	if len(result.Messages) != 1 {
		t.Fatalf("expected 1 email log, got %d", len(result.Messages))
	}
	if result.TotalCount != 1 {
		t.Errorf("expected total_count 1, got %d", result.TotalCount)
	}
	if result.NextPageCursor != "cursor-abc" {
		t.Errorf("expected next_page_cursor 'cursor-abc', got %q", result.NextPageCursor)
	}
}

func TestEmailLogsListWithFilters(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		q := r.URL.Query()
		if q.Get("filters[to][value]") != "user@test.com" {
			t.Errorf("expected to filter 'user@test.com', got %s", q.Get("filters[to][value]"))
		}
		if q.Get("filters[status][value]") != "delivered" {
			t.Errorf("expected status filter 'delivered', got %s", q.Get("filters[status][value]"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"messages": []map[string]interface{}{
				{"message_id": "msg-1", "subject": "Filtered", "from": "a@b.com", "to": "user@test.com", "status": "delivered", "sent_at": "2024-01-01"},
			},
			"total_count":      1,
			"next_page_cursor": "",
		})
	})
	defer cleanup()

	cmd := emaillogs.NewCmdEmailLogs(f)
	cmd.SetArgs([]string{"list", "--to", "user@test.com", "--status", "delivered"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Filtered") {
		t.Errorf("expected output to contain 'Filtered', got:\n%s", output)
	}
}

func TestEmailLogsGet(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/accounts/123/email_logs/msg-1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message_id": "msg-1",
			"subject":    "Test",
			"from":       "a@b.com",
			"to":         "c@d.com",
			"status":     "delivered",
			"sent_at":    "2024-01-01",
		})
	})
	defer cleanup()

	cmd := emaillogs.NewCmdEmailLogs(f)
	cmd.SetArgs([]string{"get", "--id", "msg-1"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Test") {
		t.Errorf("expected output to contain 'Test', got:\n%s", output)
	}
	if !strings.Contains(output, "delivered") {
		t.Errorf("expected output to contain 'delivered', got:\n%s", output)
	}
}

func TestEmailLogsGetJSONKeepsResponseAsIs(t *testing.T) {
	body := `{"message_id":"msg-1","status":"delivered","client_ip":null,"category":"Welcome",` +
		`"custom_variables":{"user_id":"42"},"sending_stream":"transactional","domain_id":3,` +
		`"template_id":null,"references":[],"opens_count":2,"clicks_count":0,` +
		`"events":[{"event_type":"delivery","created_at":"2024-01-01T00:00:01Z"}]}`

	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, body)
	})
	defer cleanup()

	viper.Set("output", "json")

	cmd := emaillogs.NewCmdEmailLogs(f)
	cmd.SetArgs([]string{"get", "--id", "msg-1"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got, want map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput:\n%s", err, buf.String())
	}
	json.Unmarshal([]byte(body), &want)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected the API response unchanged\nwant: %v\ngot:  %v", want, got)
	}
}
