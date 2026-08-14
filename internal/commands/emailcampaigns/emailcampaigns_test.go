package emailcampaigns_test

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
	"github.com/mailtrap/mailtrap-cli/internal/commands/emailcampaigns"
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

func sampleCampaign() map[string]interface{} {
	return map[string]interface{}{
		"id":                    4567,
		"domain_id":             4321,
		"domain_name":           "acme.com",
		"name":                  "Spring Sale",
		"from_local_part":       "news",
		"from_display_name":     "Acme Marketing",
		"current_state":         "draft",
		"created_at":            "2026-05-01T10:15:00.000Z",
		"updated_at":            "2026-05-02T09:00:00.000Z",
		"recipient_total_count": 1500,
		"contact_list_ids":      []int{55, 56},
		"delivery_mode":         "rapid",
		"template": map[string]interface{}{
			"id":      789,
			"subject": "Spring is here",
		},
	}
}

func TestEmailCampaignsList(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/email_campaigns" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Api-Token") != "test-token" {
			t.Errorf("expected Api-Token header 'test-token', got %q", r.Header.Get("Api-Token"))
		}
		if r.URL.Query().Get("per_page") != "10" {
			t.Errorf("expected per_page=10, got %q", r.URL.Query().Get("per_page"))
		}
		if r.URL.Query().Get("search") != "spring" {
			t.Errorf("expected search=spring, got %q", r.URL.Query().Get("search"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data":       []map[string]interface{}{sampleCampaign()},
			"pagination": map[string]interface{}{"token": 1, "next_token": nil},
		})
	})
	defer cleanup()

	cmd := emailcampaigns.NewCmdEmailCampaigns(f)
	cmd.SetArgs([]string{"list", "--per-page", "10", "--search", "spring"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Spring Sale") {
		t.Errorf("expected output to contain 'Spring Sale', got:\n%s", output)
	}
	if !strings.Contains(output, "draft") {
		t.Errorf("expected output to contain 'draft', got:\n%s", output)
	}
	if !strings.Contains(output, "acme.com") {
		t.Errorf("expected output to contain 'acme.com', got:\n%s", output)
	}
}

func TestEmailCampaignsListJSON(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{sampleCampaign()},
		})
	})
	defer cleanup()

	viper.Set("output", "json")

	cmd := emailcampaigns.NewCmdEmailCampaigns(f)
	cmd.SetArgs([]string{"list"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput:\n%s", err, buf.String())
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 campaign, got %d", len(result))
	}
	if result[0]["name"] != "Spring Sale" {
		t.Errorf("expected name 'Spring Sale', got %v", result[0]["name"])
	}
	if result[0]["domain_id"] != float64(4321) {
		t.Errorf("expected domain_id 4321, got %v", result[0]["domain_id"])
	}
}

func TestEmailCampaignsGet(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/email_campaigns/4567" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": sampleCampaign()})
	})
	defer cleanup()

	cmd := emailcampaigns.NewCmdEmailCampaigns(f)
	cmd.SetArgs([]string{"get", "--id", "4567"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Spring Sale") {
		t.Errorf("expected output to contain 'Spring Sale', got:\n%s", output)
	}
	if !strings.Contains(output, "1500") {
		t.Errorf("expected output to contain recipient count, got:\n%s", output)
	}
}

func TestEmailCampaignsGetMissingID(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	buf := &bytes.Buffer{}
	f.IOStreams.Out = buf

	cmd := emailcampaigns.NewCmdEmailCampaigns(f)
	cmd.SetArgs([]string{"get"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when --id is missing")
	}
	if !strings.Contains(err.Error(), "--id is required") {
		t.Errorf("expected '--id is required' error, got: %v", err)
	}
}

func TestEmailCampaignsCreate(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/email_campaigns" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]interface{}
		json.Unmarshal(body, &reqBody)

		// The body is flat — no {"email_campaign": ...} wrapper.
		if _, ok := reqBody["email_campaign"]; ok {
			t.Errorf("did not expect an 'email_campaign' wrapper, got %v", reqBody)
		}
		if reqBody["name"] != "Spring Sale" {
			t.Errorf("unexpected name: %v", reqBody["name"])
		}
		if reqBody["domain_id"] != float64(4321) {
			t.Errorf("unexpected domain_id: %v", reqBody["domain_id"])
		}
		if reqBody["from_local_part"] != "news" {
			t.Errorf("unexpected from_local_part: %v", reqBody["from_local_part"])
		}
		if reqBody["delivery_mode"] != "gradual" {
			t.Errorf("unexpected delivery_mode: %v", reqBody["delivery_mode"])
		}
		templateAttrs, ok := reqBody["template_attributes"].(map[string]interface{})
		if !ok || templateAttrs["subject"] != "Spring is here" {
			t.Errorf("unexpected template_attributes: %v", reqBody["template_attributes"])
		}
		if templateAttrs["body_html"] != "<h1>Hi</h1>" {
			t.Errorf("unexpected body_html: %v", templateAttrs["body_html"])
		}
		replyTo, ok := reqBody["reply_to"].(map[string]interface{})
		if !ok || replyTo["local_part"] != "support" || replyTo["domain"] != "acme.com" {
			t.Errorf("unexpected reply_to: %v", reqBody["reply_to"])
		}
		deliveryOptions, ok := reqBody["delivery_options"].(map[string]interface{})
		if !ok || deliveryOptions["emails_per_hour"] != float64(1000) {
			t.Errorf("unexpected delivery_options: %v", reqBody["delivery_options"])
		}
		if lists, ok := reqBody["contact_list_ids"].([]interface{}); !ok || len(lists) != 2 {
			t.Errorf("expected 2 contact_list_ids, got %v", reqBody["contact_list_ids"])
		}
		if segments, ok := reqBody["contact_segment_ids"].([]interface{}); !ok || len(segments) != 1 {
			t.Errorf("expected 1 contact_segment_id, got %v", reqBody["contact_segment_ids"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{"data": sampleCampaign()})
	})
	defer cleanup()

	cmd := emailcampaigns.NewCmdEmailCampaigns(f)
	cmd.SetArgs([]string{
		"create",
		"--name", "Spring Sale",
		"--domain-id", "4321",
		"--from-local-part", "news",
		"--subject", "Spring is here",
		"--body-html", "<h1>Hi</h1>",
		"--reply-to-local-part", "support",
		"--reply-to-domain", "acme.com",
		"--delivery-mode", "gradual",
		"--emails-per-hour", "1000",
		"--contact-list-ids", "55,56",
		"--contact-segment-ids", "12",
	})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Spring Sale") {
		t.Errorf("expected output to contain 'Spring Sale', got:\n%s", output)
	}
	if !strings.Contains(output, "draft") {
		t.Errorf("expected output to contain 'draft', got:\n%s", output)
	}
}

func TestEmailCampaignsCreateMissingName(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	buf := &bytes.Buffer{}
	f.IOStreams.Out = buf

	cmd := emailcampaigns.NewCmdEmailCampaigns(f)
	cmd.SetArgs([]string{
		"create",
		"--domain-id", "4321",
		"--from-local-part", "news",
		"--subject", "Spring is here",
	})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when --name is missing")
	}
	if !strings.Contains(err.Error(), "--name is required") {
		t.Errorf("expected '--name is required' error, got: %v", err)
	}
}

func TestEmailCampaignsUpdate(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/email_campaigns/4567" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]interface{}
		json.Unmarshal(body, &reqBody)

		if reqBody["name"] != "Spring Sale (updated)" {
			t.Errorf("unexpected name: %v", reqBody["name"])
		}
		templateAttrs, ok := reqBody["template_attributes"].(map[string]interface{})
		if !ok || templateAttrs["subject"] != "New subject" {
			t.Errorf("unexpected template_attributes: %v", reqBody["template_attributes"])
		}
		// PATCH semantics: unchanged flags stay out of the body.
		if _, ok := templateAttrs["body_html"]; ok {
			t.Errorf("did not expect body_html to be set, got %v", templateAttrs["body_html"])
		}
		if _, ok := reqBody["domain_id"]; ok {
			t.Errorf("did not expect domain_id to be set, got %v", reqBody["domain_id"])
		}
		if _, ok := reqBody["reply_to"]; ok {
			t.Errorf("did not expect reply_to to be set, got %v", reqBody["reply_to"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": sampleCampaign()})
	})
	defer cleanup()

	cmd := emailcampaigns.NewCmdEmailCampaigns(f)
	cmd.SetArgs([]string{
		"update",
		"--id", "4567",
		"--name", "Spring Sale (updated)",
		"--subject", "New subject",
	})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "Spring Sale") {
		t.Errorf("expected output to contain 'Spring Sale', got:\n%s", buf.String())
	}
}

func TestEmailCampaignsUpdateClearAudience(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]interface{}
		json.Unmarshal(body, &reqBody)

		// An explicit empty array clears the audience; the API treats the ids
		// as the full set of included lists.
		lists, ok := reqBody["contact_list_ids"].([]interface{})
		if !ok || len(lists) != 0 {
			t.Errorf("expected contact_list_ids to be [], got %v", reqBody["contact_list_ids"])
		}
		if _, ok := reqBody["contact_segment_ids"]; ok {
			t.Errorf("did not expect contact_segment_ids to be set, got %v", reqBody["contact_segment_ids"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": sampleCampaign()})
	})
	defer cleanup()

	cmd := emailcampaigns.NewCmdEmailCampaigns(f)
	cmd.SetArgs([]string{
		"update",
		"--id", "4567",
		"--clear-contact-lists",
	})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEmailCampaignsUpdateClearConflictsWithIDs(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		t.Error("request should not be sent when flags conflict")
	})
	defer cleanup()

	cmd := emailcampaigns.NewCmdEmailCampaigns(f)
	cmd.SetArgs([]string{
		"update",
		"--id", "4567",
		"--contact-list-ids", "55,56",
		"--clear-contact-lists",
	})
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected an error for mutually exclusive flags")
	}
}

func TestEmailCampaignsUpdateNoAttributes(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		t.Error("request should not be sent when no attribute flags are provided")
	})
	defer cleanup()

	cmd := emailcampaigns.NewCmdEmailCampaigns(f)
	cmd.SetArgs([]string{"update", "--id", "4567"})
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when no attribute flags are provided")
	}
	if !strings.Contains(err.Error(), "at least one attribute flag is required") {
		t.Errorf("expected 'at least one attribute flag is required' error, got: %v", err)
	}
}

func TestEmailCampaignsDelete(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/email_campaigns/4567" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer cleanup()

	cmd := emailcampaigns.NewCmdEmailCampaigns(f)
	cmd.SetArgs([]string{"delete", "--id", "4567"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "Email campaign deleted successfully") {
		t.Errorf("expected success message, got:\n%s", buf.String())
	}
}

func testLifecycleAction(t *testing.T, action, resultState string) {
	t.Helper()

	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/email_campaigns/4567/"+action {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		campaign := sampleCampaign()
		campaign["current_state"] = resultState
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": campaign})
	})
	defer cleanup()

	cmd := emailcampaigns.NewCmdEmailCampaigns(f)
	cmd.SetArgs([]string{action, "--id", "4567"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), resultState) {
		t.Errorf("expected output to contain %q, got:\n%s", resultState, buf.String())
	}
}

func TestEmailCampaignsStart(t *testing.T) {
	testLifecycleAction(t, "start", "started")
}

func TestEmailCampaignsCancel(t *testing.T) {
	testLifecycleAction(t, "cancel", "draft")
}

func TestEmailCampaignsTerminate(t *testing.T) {
	testLifecycleAction(t, "terminate", "terminating")
}

func TestEmailCampaignsReset(t *testing.T) {
	testLifecycleAction(t, "reset", "draft")
}

func TestEmailCampaignsSchedule(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/email_campaigns/4567/schedule" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]interface{}
		json.Unmarshal(body, &reqBody)
		if reqBody["datetime"] != "2026-06-01T09:00:00.000Z" {
			t.Errorf("unexpected datetime: %v", reqBody["datetime"])
		}

		campaign := sampleCampaign()
		campaign["current_state"] = "scheduled"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": campaign})
	})
	defer cleanup()

	cmd := emailcampaigns.NewCmdEmailCampaigns(f)
	cmd.SetArgs([]string{"schedule", "--id", "4567", "--datetime", "2026-06-01T09:00:00.000Z"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "scheduled") {
		t.Errorf("expected output to contain 'scheduled', got:\n%s", buf.String())
	}
}

func TestEmailCampaignsScheduleMissingDatetime(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	buf := &bytes.Buffer{}
	f.IOStreams.Out = buf

	cmd := emailcampaigns.NewCmdEmailCampaigns(f)
	cmd.SetArgs([]string{"schedule", "--id", "4567"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when --datetime is missing")
	}
	if !strings.Contains(err.Error(), "--datetime is required") {
		t.Errorf("expected '--datetime is required' error, got: %v", err)
	}
}

func TestEmailCampaignsStats(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/email_campaigns/4567/stats" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("start_date") != "2026-05-01" {
			t.Errorf("expected start_date=2026-05-01, got %q", r.URL.Query().Get("start_date"))
		}
		if r.URL.Query().Get("end_date") != "2026-05-31" {
			t.Errorf("expected end_date=2026-05-31, got %q", r.URL.Query().Get("end_date"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"delivery_count":       1450,
				"open_count":           820,
				"click_count":          310,
				"bounce_count":         30,
				"unsubscription_count": 12,
				"sent_count":           1500,
				"spam_count":           5,
				"delivery_rate":        0.9667,
				"open_rate":            0.5655,
				"click_rate":           0.2138,
				"bounce_rate":          0.02,
				"spam_rate":            0.0033,
				"unsubscription_rate":  0.0083,
			},
		})
	})
	defer cleanup()

	cmd := emailcampaigns.NewCmdEmailCampaigns(f)
	cmd.SetArgs([]string{"stats", "--id", "4567", "--start-date", "2026-05-01", "--end-date", "2026-05-31"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "1450") {
		t.Errorf("expected output to contain delivery count, got:\n%s", output)
	}
	if !strings.Contains(output, "0.9667") {
		t.Errorf("expected output to contain delivery rate, got:\n%s", output)
	}
}
