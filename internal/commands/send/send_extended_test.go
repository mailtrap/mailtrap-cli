package send_test

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/mailtrap/mailtrap-cli/internal/commands/send"
)

func TestSendBulk(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/send" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]interface{}
		json.Unmarshal(body, &reqBody)

		from, ok := reqBody["from"].(map[string]interface{})
		if !ok {
			t.Fatal("expected 'from' field in request body")
		}
		if from["email"] != "test@example.com" {
			t.Errorf("expected from email 'test@example.com', got %v", from["email"])
		}

		to, ok := reqBody["to"].([]interface{})
		if !ok || len(to) == 0 {
			t.Fatal("expected 'to' field in request body")
		}
		toFirst := to[0].(map[string]interface{})
		if toFirst["email"] != "user@example.com" {
			t.Errorf("expected to email 'user@example.com', got %v", toFirst["email"])
		}

		if reqBody["subject"] != "Update" {
			t.Errorf("expected subject 'Update', got %v", reqBody["subject"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":     true,
			"message_ids": []string{"msg-1"},
		})
	})
	defer cleanup()

	cmd := send.NewCmdSend(f)
	cmd.SetArgs([]string{
		"bulk",
		"--from", "test@example.com",
		"--to", "user@example.com",
		"--subject", "Update",
		"--text", "News",
	})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "true") {
		t.Errorf("expected output to contain 'true', got:\n%s", output)
	}
}

func TestSendBatchTransactional(t *testing.T) {
	batchData := `[{"from":{"email":"test@example.com"},"to":[{"email":"user@example.com"}],"subject":"Batch","text":"Hello"}]`

	tmpFile, err := os.CreateTemp("", "batch-transactional-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(batchData); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/batch" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":   true,
			"responses": []interface{}{},
		})
	})
	defer cleanup()

	cmd := send.NewCmdSend(f)
	cmd.SetArgs([]string{"batch-transactional", "--file", tmpFile.Name()})
	cmd.SetOut(buf)

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "true") {
		t.Errorf("expected output to contain 'true', got:\n%s", output)
	}
}

func TestSendBatchBulk(t *testing.T) {
	batchData := `[{"from":{"email":"test@example.com"},"to":[{"email":"user@example.com"}],"subject":"Bulk Batch","text":"Hello"}]`

	tmpFile, err := os.CreateTemp("", "batch-bulk-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(batchData); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/batch" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":   true,
			"responses": []interface{}{},
		})
	})
	defer cleanup()

	cmd := send.NewCmdSend(f)
	cmd.SetArgs([]string{"batch-bulk", "--file", tmpFile.Name()})
	cmd.SetOut(buf)

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "true") {
		t.Errorf("expected output to contain 'true', got:\n%s", output)
	}
}
