package messages_test

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/mailtrap/mailtrap-cli/internal/commands/messages"
	"github.com/spf13/viper"
)

func TestMessagesUpdate(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/inboxes/1/messages/1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]interface{}
		json.Unmarshal(body, &payload)
		msg, ok := payload["message"].(map[string]interface{})
		if !ok {
			t.Errorf("expected 'message' key in body, got: %s", string(body))
		}
		if msg["is_read"] != true {
			t.Errorf("expected is_read to be true, got: %v", msg["is_read"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": 1, "subject": "Welcome", "from_email": "a@b.com", "to_email": "c@d.com", "is_read": true, "created_at": "2024-01-01",
		})
	})
	defer cleanup()

	cmd := messages.NewCmdMessages(f)
	cmd.SetArgs([]string{"update", "--sandbox-id", "1", "--id", "1", "--is-read"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Welcome") {
		t.Errorf("expected output to contain 'Welcome', got:\n%s", output)
	}
}

func TestMessagesDelete(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/inboxes/1/messages/1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
	defer cleanup()

	cmd := messages.NewCmdMessages(f)
	cmd.SetArgs([]string{"delete", "--sandbox-id", "1", "--id", "1"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "deleted successfully") {
		t.Errorf("expected 'deleted successfully' in output, got:\n%s", output)
	}
}

func TestMessagesForward(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/inboxes/1/messages/1/forward") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]interface{}
		json.Unmarshal(body, &payload)
		if payload["email"] != "fwd@test.com" {
			t.Errorf("expected email 'fwd@test.com', got: %v", payload["email"])
		}

		w.WriteHeader(http.StatusOK)
	})
	defer cleanup()

	cmd := messages.NewCmdMessages(f)
	cmd.SetArgs([]string{"forward", "--sandbox-id", "1", "--id", "1", "--email", "fwd@test.com"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "forwarded successfully") {
		t.Errorf("expected 'forwarded successfully' in output, got:\n%s", output)
	}
}

func TestMessagesSpamScore(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/inboxes/1/messages/1/spam_report") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"score":5.2,"details":[]}`))
	})
	defer cleanup()

	cmd := messages.NewCmdMessages(f)
	cmd.SetArgs([]string{"spam-score", "--sandbox-id", "1", "--id", "1"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "score") {
		t.Errorf("expected output to contain 'score', got:\n%s", output)
	}
}

func TestMessagesHTMLAnalysis(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/inboxes/1/messages/1/analyze") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})
	defer cleanup()

	// html-analysis uses GetOutputFormat, so set output to json to ensure clean output
	viper.Set("output", "json")

	cmd := messages.NewCmdMessages(f)
	cmd.SetArgs([]string{"html-analysis", "--sandbox-id", "1", "--id", "1"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "status") {
		t.Errorf("expected output to contain 'status', got:\n%s", output)
	}
}

func TestMessagesHeaders(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/inboxes/1/messages/1/mail_headers") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"headers":{"Subject":"Test","From":"a@b.com"}}`))
	})
	defer cleanup()

	cmd := messages.NewCmdMessages(f)
	cmd.SetArgs([]string{"headers", "--sandbox-id", "1", "--id", "1"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Subject") {
		t.Errorf("expected output to contain 'Subject', got:\n%s", output)
	}
}

func TestMessagesHTML(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/inboxes/1/messages/1/body.html") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<h1>Hello</h1>"))
	})
	defer cleanup()

	cmd := messages.NewCmdMessages(f)
	cmd.SetArgs([]string{"html", "--sandbox-id", "1", "--id", "1"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "<h1>Hello</h1>") {
		t.Errorf("expected output to contain '<h1>Hello</h1>', got:\n%s", output)
	}
}

func TestMessagesText(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/inboxes/1/messages/1/body.txt") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Hello plain text"))
	})
	defer cleanup()

	cmd := messages.NewCmdMessages(f)
	cmd.SetArgs([]string{"text", "--sandbox-id", "1", "--id", "1"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Hello plain text") {
		t.Errorf("expected output to contain 'Hello plain text', got:\n%s", output)
	}
}

func TestMessagesSource(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/inboxes/1/messages/1/body.htmlsource") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<html><body>Source</body></html>"))
	})
	defer cleanup()

	cmd := messages.NewCmdMessages(f)
	cmd.SetArgs([]string{"source", "--sandbox-id", "1", "--id", "1"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Source") {
		t.Errorf("expected output to contain 'Source', got:\n%s", output)
	}
}

func TestMessagesRaw(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/inboxes/1/messages/1/body.raw") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "message/rfc822")
		w.Write([]byte("MIME-Version: 1.0\r\nSubject: Test\r\n\r\nRaw body"))
	})
	defer cleanup()

	cmd := messages.NewCmdMessages(f)
	cmd.SetArgs([]string{"raw", "--sandbox-id", "1", "--id", "1"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Raw body") {
		t.Errorf("expected output to contain 'Raw body', got:\n%s", output)
	}
}

func TestMessagesEml(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/accounts/123/inboxes/1/messages/1/body.eml") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "message/rfc822")
		w.Write([]byte("From: a@b.com\r\nTo: c@d.com\r\nSubject: Test EML\r\n\r\nEML body content"))
	})
	defer cleanup()

	cmd := messages.NewCmdMessages(f)
	cmd.SetArgs([]string{"eml", "--sandbox-id", "1", "--id", "1"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "EML body content") {
		t.Errorf("expected output to contain 'EML body content', got:\n%s", output)
	}
}
