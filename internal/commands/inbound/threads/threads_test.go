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
	if !strings.Contains(buf.String(), "--last-id thr_1") {
		t.Errorf("expected output to surface the next-page cursor, got:\n%s", buf.String())
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

	var result struct {
		Data       []map[string]interface{} `json:"data"`
		TotalCount int                      `json:"total_count"`
		LastID     string                   `json:"last_id"`
	}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput:\n%s", err, buf.String())
	}
	if len(result.Data) != 1 || result.Data[0]["id"] != "thr_1" {
		t.Errorf("unexpected JSON data: %v", result.Data)
	}
	if result.TotalCount != 1 {
		t.Errorf("expected total_count 1, got %d", result.TotalCount)
	}
	if result.LastID != "thr_1" {
		t.Errorf("expected last_id 'thr_1', got %q", result.LastID)
	}
}

func TestThreadsListWithSearch(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("search"); got != "acme corp" {
			t.Errorf("expected search=acme corp, got %q", got)
		}
		if got := r.URL.Query().Get("last_id"); got != "WzE3NzgyNDE5MDAwMDAsIjE3MDAwMDAwMDAwMDAxMjMiXQ==" {
			t.Errorf("expected opaque last_id cursor, got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data":        []map[string]interface{}{{"id": "thr_2", "subject": "ACME order"}},
			"total_count": 2,
			"last_id":     "WzE3NzgyNDE5MDAwMDAsIjE3MDAwMDAwMDAwMDAxMjQiXQ==",
		})
	})
	defer cleanup()

	cmd := threads.NewCmdThreads(f)
	cmd.SetArgs([]string{"list", "--inbox-id", "201", "--search", "acme corp", "--last-id", "WzE3NzgyNDE5MDAwMDAsIjE3MDAwMDAwMDAwMDAxMjMiXQ=="})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `--last-id WzE3NzgyNDE5MDAwMDAsIjE3MDAwMDAwMDAwMDAxMjQiXQ== --search "acme corp"`
	if !strings.Contains(buf.String(), want) {
		t.Errorf("expected next-page hint to carry the search, got:\n%s", buf.String())
	}
}

func TestThreadsListWithoutSearchSendsNoQuery(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Errorf("expected no query string, got %q", r.URL.RawQuery)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": []interface{}{}, "total_count": 0, "last_id": nil})
	})
	defer cleanup()

	cmd := threads.NewCmdThreads(f)
	cmd.SetArgs([]string{"list", "--inbox-id", "201"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestThreadsGetDeliveryAndForwardsJSON(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": "thr_1", "subject": "Support request",
			"messages": []map[string]interface{}{
				{"visibility_status": "placeholder", "direction": "inbound"},
				{
					"id": "msg_1", "direction": "inbound", "visibility_status": "available",
					"forwards": []map[string]interface{}{
						{
							"rule_id": 7, "rule_name": "Copy to support team", "destination": "team@example.com",
							"status": "forwarded", "reason": nil, "message_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
						},
					},
				},
				{
					"id": "1a2b3c4d", "direction": "outbound", "visibility_status": "available",
					"delivery": map[string]interface{}{
						"to": "customer@example.com", "status": "delivered",
						"delivered_at": "2026-05-08T11:40:05.000Z", "bounced_at": nil,
					},
				},
			},
		})
	})
	defer cleanup()

	viper.Set("output", "json")

	cmd := threads.NewCmdThreads(f)
	cmd.SetArgs([]string{"get", "--inbox-id", "201", "--id", "thr_1"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result struct {
		Messages []struct {
			Delivery *struct {
				To          string  `json:"to"`
				Status      string  `json:"status"`
				DeliveredAt *string `json:"delivered_at"`
				BouncedAt   *string `json:"bounced_at"`
			} `json:"delivery"`
			Forwards []struct {
				Destination string `json:"destination"`
				Status      string `json:"status"`
			} `json:"forwards"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput:\n%s", err, buf.String())
	}
	if len(result.Messages) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(result.Messages))
	}

	placeholder := result.Messages[0]
	if placeholder.Delivery != nil || placeholder.Forwards != nil {
		t.Errorf("placeholder must carry neither delivery nor forwards: %+v", placeholder)
	}

	inbound := result.Messages[1]
	if inbound.Delivery != nil {
		t.Errorf("inbound message must not carry delivery: %+v", inbound.Delivery)
	}
	if len(inbound.Forwards) != 1 || inbound.Forwards[0].Destination != "team@example.com" || inbound.Forwards[0].Status != "forwarded" {
		t.Errorf("unexpected forwards: %+v", inbound.Forwards)
	}

	outbound := result.Messages[2]
	if outbound.Delivery == nil {
		t.Fatal("expected delivery on the outbound message")
	}
	if outbound.Delivery.To != "customer@example.com" || outbound.Delivery.Status != "delivered" {
		t.Errorf("unexpected delivery: %+v", outbound.Delivery)
	}
	if outbound.Delivery.DeliveredAt == nil || *outbound.Delivery.DeliveredAt != "2026-05-08T11:40:05.000Z" {
		t.Errorf("unexpected delivered_at: %v", outbound.Delivery.DeliveredAt)
	}
	if outbound.Delivery.BouncedAt != nil {
		t.Errorf("expected nil bounced_at, got %v", *outbound.Delivery.BouncedAt)
	}
}
