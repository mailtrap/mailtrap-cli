package threads_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/commands/inbound/threads"
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

func TestThreadsList(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/api/inbound/inboxes/201/threads") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "thr_1", "subject": "Support request", "message_count": 3},
			},
			"total_count": 1,
			"last_id":     "thr_1",
		})
	})
	defer cleanup()

	cmd := threads.NewCmdThreads(f)
	cmd.SetArgs([]string{"list", "--inbox-id", "201"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "Support request") {
		t.Errorf("expected output to contain subject, got:\n%s", buf.String())
	}
}

func TestThreadsGet(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/api/inbound/inboxes/201/threads/thr_1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": "thr_1", "subject": "Support request",
			"messages": []map[string]interface{}{
				{"id": "msg_1", "direction": "inbound", "visibility_status": "available"},
			},
		})
	})
	defer cleanup()

	cmd := threads.NewCmdThreads(f)
	cmd.SetArgs([]string{"get", "--inbox-id", "201", "--id", "thr_1"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "Support request") {
		t.Errorf("expected output to contain subject, got:\n%s", buf.String())
	}
}

func TestThreadsDelete(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/inbound/inboxes/201/threads/thr_1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer cleanup()

	cmd := threads.NewCmdThreads(f)
	cmd.SetArgs([]string{"delete", "--inbox-id", "201", "--id", "thr_1"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "deleted successfully") {
		t.Errorf("expected success message, got:\n%s", buf.String())
	}
}

func TestThreadsListJSON(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data":        []map[string]interface{}{{"id": "thr_1", "subject": "Support request"}},
			"total_count": 1,
			"last_id":     "thr_1",
		})
	})
	defer cleanup()

	viper.Set("output", "json")

	cmd := threads.NewCmdThreads(f)
	cmd.SetArgs([]string{"list", "--inbox-id", "201"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput:\n%s", err, buf.String())
	}
	if len(result) != 1 || result[0]["id"] != "thr_1" {
		t.Errorf("unexpected JSON result: %v", result)
	}
}
