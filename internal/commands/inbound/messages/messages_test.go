package messages_test

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
	"github.com/mailtrap/mailtrap-cli/internal/commands/inbound/messages"
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

func TestMessagesList(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/api/inbound/inboxes/201/messages") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "msg_1", "from": "customer@example.com", "subject": "Support request", "thread_id": "thr_1"},
			},
			"total_count": 1,
			"last_id":     "msg_1",
		})
	})
	defer cleanup()

	cmd := messages.NewCmdMessages(f)
	cmd.SetArgs([]string{"list", "--inbox-id", "201"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "Support request") {
		t.Errorf("expected output to contain subject, got:\n%s", buf.String())
	}
}

func TestMessagesListWithCursor(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("last_id"); got != "msg_2" {
			t.Errorf("expected last_id=msg_2, got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data":        []map[string]interface{}{{"id": "msg_3", "subject": "Next page"}},
			"total_count": 3,
			"last_id":     "msg_3",
		})
	})
	defer cleanup()

	cmd := messages.NewCmdMessages(f)
	cmd.SetArgs([]string{"list", "--inbox-id", "201", "--last-id", "msg_2"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMessagesGet(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/api/inbound/inboxes/201/messages/msg_1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": "msg_1", "subject": "Support request", "html_body": "<p>Hello</p>",
		})
	})
	defer cleanup()

	cmd := messages.NewCmdMessages(f)
	cmd.SetArgs([]string{"get", "--inbox-id", "201", "--id", "msg_1"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "Support request") {
		t.Errorf("expected output to contain subject, got:\n%s", buf.String())
	}
}

func TestMessagesReply(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/inbound/inboxes/201/messages/msg_1/reply") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]interface{}
		json.Unmarshal(body, &reqBody)
		if reqBody["text"] != "Thanks!" {
			t.Errorf("unexpected text: %v", reqBody["text"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"message_ids": []string{"0000000000000001"}})
	})
	defer cleanup()

	cmd := messages.NewCmdMessages(f)
	cmd.SetArgs([]string{"reply", "--inbox-id", "201", "--id", "msg_1", "--text", "Thanks!"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "0000000000000001") {
		t.Errorf("expected output to contain message id, got:\n%s", buf.String())
	}
}

func TestMessagesForward(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/api/inbound/inboxes/201/messages/msg_1/forward") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]interface{}
		json.Unmarshal(body, &reqBody)
		to, ok := reqBody["to"].([]interface{})
		if !ok || len(to) != 1 {
			t.Fatalf("expected one 'to' recipient, got %v", reqBody["to"])
		}
		if to[0].(map[string]interface{})["email"] != "colleague@example.com" {
			t.Errorf("unexpected to: %v", to[0])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"message_ids": []string{"0000000000000002"}})
	})
	defer cleanup()

	cmd := messages.NewCmdMessages(f)
	cmd.SetArgs([]string{"forward", "--inbox-id", "201", "--id", "msg_1", "--to", "colleague@example.com", "--text", "FYI"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "0000000000000002") {
		t.Errorf("expected output to contain message id, got:\n%s", buf.String())
	}
}

func TestMessagesForwardMissingTo(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	cmd := messages.NewCmdMessages(f)
	cmd.SetArgs([]string{"forward", "--inbox-id", "201", "--id", "msg_1", "--text", "FYI"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when --to is missing")
	}
	if !strings.Contains(err.Error(), "--to is required") {
		t.Errorf("expected '--to is required' error, got: %v", err)
	}
}

func TestMessagesDelete(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/inbound/inboxes/201/messages/msg_1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer cleanup()

	cmd := messages.NewCmdMessages(f)
	cmd.SetArgs([]string{"delete", "--inbox-id", "201", "--id", "msg_1"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "deleted successfully") {
		t.Errorf("expected success message, got:\n%s", buf.String())
	}
}

func TestMessagesListJSON(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data":        []map[string]interface{}{{"id": "msg_1", "subject": "Support request"}},
			"total_count": 1,
			"last_id":     "msg_1",
		})
	})
	defer cleanup()

	viper.Set("output", "json")

	cmd := messages.NewCmdMessages(f)
	cmd.SetArgs([]string{"list", "--inbox-id", "201"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput:\n%s", err, buf.String())
	}
	if len(result) != 1 || result[0]["id"] != "msg_1" {
		t.Errorf("unexpected JSON result: %v", result)
	}
}
