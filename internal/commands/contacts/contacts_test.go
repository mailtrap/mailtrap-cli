package contacts_test

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
	"github.com/mailtrap/mailtrap-cli/internal/commands/contacts"
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

func TestContactsGet(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/accounts/123/contacts/abc-123") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Api-Token") != "test-token" {
			t.Errorf("expected Api-Token header 'test-token', got %q", r.Header.Get("Api-Token"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{"id":"abc-123","email":"user@test.com","status":"subscribed","created_at":"2024-01-01","updated_at":"2024-01-02"}}`))
	})
	defer cleanup()

	cmd := contacts.NewCmdContacts(f)
	cmd.SetArgs([]string{"get", "--id", "abc-123"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "user@test.com") {
		t.Errorf("expected output to contain 'user@test.com', got:\n%s", output)
	}
}

func TestContactsGetMissingID(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	cmd := contacts.NewCmdContacts(f)
	cmd.SetArgs([]string{"get"})
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "--id is required") {
		t.Errorf("expected '--id is required' error, got: %v", err)
	}
}

func TestContactsCreate(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/accounts/123/contacts") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal body: %v", err)
		}

		contact, ok := payload["contact"].(map[string]interface{})
		if !ok {
			t.Fatal("expected 'contact' key in body")
		}
		if contact["email"] != "new@test.com" {
			t.Errorf("expected email 'new@test.com', got %v", contact["email"])
		}
		fields, ok := contact["fields"].(map[string]interface{})
		if !ok {
			t.Fatal("expected 'fields' key in contact")
		}
		if fields["first_name"] != "John" {
			t.Errorf("expected first_name 'John', got %v", fields["first_name"])
		}
		if fields["last_name"] != "Doe" {
			t.Errorf("expected last_name 'Doe', got %v", fields["last_name"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{"id":"new-id","email":"new@test.com","status":"subscribed","created_at":"2024-01-01","updated_at":"2024-01-01"}}`))
	})
	defer cleanup()

	cmd := contacts.NewCmdContacts(f)
	cmd.SetArgs([]string{"create", "--email", "new@test.com", "--first-name", "John", "--last-name", "Doe"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "new@test.com") {
		t.Errorf("expected output to contain 'new@test.com', got:\n%s", output)
	}
}

func TestContactsCreateWithListIDs(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal body: %v", err)
		}

		contact, ok := payload["contact"].(map[string]interface{})
		if !ok {
			t.Fatal("expected 'contact' key in body")
		}
		listIDs, ok := contact["list_ids"].([]interface{})
		if !ok {
			t.Fatal("expected 'list_ids' key in contact")
		}
		if len(listIDs) != 2 {
			t.Errorf("expected 2 list_ids, got %d", len(listIDs))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{"id":"new-id","email":"x@y.com","status":"subscribed","created_at":"2024-01-01","updated_at":"2024-01-01"}}`))
	})
	defer cleanup()

	cmd := contacts.NewCmdContacts(f)
	cmd.SetArgs([]string{"create", "--email", "x@y.com", "--list-ids", "1,2"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestContactsCreateMissingEmail(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	cmd := contacts.NewCmdContacts(f)
	cmd.SetArgs([]string{"create"})
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "--email is required") {
		t.Errorf("expected '--email is required' error, got: %v", err)
	}
}

func TestContactsUpdate(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/accounts/123/contacts/abc-123") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal body: %v", err)
		}

		contact, ok := payload["contact"].(map[string]interface{})
		if !ok {
			t.Fatal("expected 'contact' key in body")
		}
		if contact["email"] != "updated@test.com" {
			t.Errorf("expected email 'updated@test.com', got %v", contact["email"])
		}
		fields, ok := contact["fields"].(map[string]interface{})
		if !ok {
			t.Fatal("expected 'fields' key in contact")
		}
		if fields["first_name"] != "Jane" {
			t.Errorf("expected first_name 'Jane', got %v", fields["first_name"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{"id":"abc-123","email":"updated@test.com","status":"subscribed","created_at":"2024-01-01","updated_at":"2024-01-02"}}`))
	})
	defer cleanup()

	cmd := contacts.NewCmdContacts(f)
	cmd.SetArgs([]string{"update", "--id", "abc-123", "--email", "updated@test.com", "--first-name", "Jane"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "updated@test.com") {
		t.Errorf("expected output to contain 'updated@test.com', got:\n%s", output)
	}
}

func TestContactsUpdateWithListFlags(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal body: %v", err)
		}

		contact, ok := payload["contact"].(map[string]interface{})
		if !ok {
			t.Fatal("expected 'contact' key in body")
		}

		included, ok := contact["list_ids_included"].([]interface{})
		if !ok {
			t.Fatal("expected 'list_ids_included' key in contact")
		}
		if len(included) != 2 {
			t.Errorf("expected 2 list_ids_included, got %d", len(included))
		}

		excluded, ok := contact["list_ids_excluded"].([]interface{})
		if !ok {
			t.Fatal("expected 'list_ids_excluded' key in contact")
		}
		if len(excluded) != 1 {
			t.Errorf("expected 1 list_ids_excluded, got %d", len(excluded))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{"id":"abc-123","email":"x@y.com","status":"subscribed","created_at":"2024-01-01","updated_at":"2024-01-02"}}`))
	})
	defer cleanup()

	cmd := contacts.NewCmdContacts(f)
	cmd.SetArgs([]string{"update", "--id", "abc-123", "--email", "x@y.com", "--list-ids-included", "1,2", "--list-ids-excluded", "3"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestContactsDelete(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/accounts/123/contacts/abc-123") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
	defer cleanup()

	cmd := contacts.NewCmdContacts(f)
	cmd.SetArgs([]string{"delete", "--id", "abc-123"})
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

func TestContactsDeleteMissingID(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	cmd := contacts.NewCmdContacts(f)
	cmd.SetArgs([]string{"delete"})
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "--id is required") {
		t.Errorf("expected '--id is required' error, got: %v", err)
	}
}

func TestContactsCreateEvent(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/accounts/123/contacts/abc-123/events") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal body: %v", err)
		}

		if payload["name"] != "purchase" {
			t.Errorf("expected name 'purchase', got %v", payload["name"])
		}
		params, ok := payload["params"].(map[string]interface{})
		if !ok {
			t.Fatal("expected 'params' to be an object")
		}
		if params["amount"] != float64(99) {
			t.Errorf("expected amount 99, got %v", params["amount"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})
	defer cleanup()

	cmd := contacts.NewCmdContacts(f)
	cmd.SetArgs([]string{"create-event", "--id", "abc-123", "--name", "purchase", "--params", `{"amount":99}`})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestContactsCreateEventMissingFlags(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	cmd := contacts.NewCmdContacts(f)
	cmd.SetArgs([]string{"create-event"})
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "is required") {
		t.Errorf("expected 'is required' error, got: %v", err)
	}
}

func TestContactsExport(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/accounts/123/contacts/exports") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal body: %v", err)
		}

		filters, ok := payload["filters"].([]interface{})
		if !ok {
			t.Fatal("expected 'filters' key in body")
		}
		if len(filters) != 2 {
			t.Errorf("expected 2 filters, got %d", len(filters))
		}

		filter0, _ := filters[0].(map[string]interface{})
		if filter0["name"] != "list_id" {
			t.Errorf("expected first filter name 'list_id', got %v", filter0["name"])
		}
		filter1, _ := filters[1].(map[string]interface{})
		if filter1["name"] != "subscription_status" {
			t.Errorf("expected second filter name 'subscription_status', got %v", filter1["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":"exp-1","status":"pending"}`))
	})
	defer cleanup()

	cmd := contacts.NewCmdContacts(f)
	cmd.SetArgs([]string{"export", "--list-ids", "1,2", "--status", "subscribed"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestContactsImportStatus(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/accounts/123/contacts/imports/imp-123") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"completed","total":100}`))
	})
	defer cleanup()

	cmd := contacts.NewCmdContacts(f)
	cmd.SetArgs([]string{"import-status", "--id", "imp-123"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "completed") {
		t.Errorf("expected output to contain 'completed', got:\n%s", output)
	}
}

func TestContactsExportStatus(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/accounts/123/contacts/exports/exp-123") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"completed"}`))
	})
	defer cleanup()

	cmd := contacts.NewCmdContacts(f)
	cmd.SetArgs([]string{"export-status", "--id", "exp-123"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "completed") {
		t.Errorf("expected output to contain 'completed', got:\n%s", output)
	}
}

func TestContactsSendSubscriptionTopics(t *testing.T) {
	topics := `[{"subscription_topic_id":7,"domain_id":4321,"status":"unsubscribed"},{"subscription_topic_id":8,"domain_id":4322,"status":"subscribed"}]`

	// A topic-only update: the CLI still requires --email.
	for _, args := range [][]string{
		{"create", "--email", "x@y.com", "--subscription-topics", topics},
		{"update", "--id", "abc-123", "--email", "x@y.com", "--subscription-topics", topics},
	} {
		t.Run(args[0], func(t *testing.T) {
			f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
				var payload map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					t.Fatalf("failed to decode body: %v", err)
				}

				want := map[string]interface{}{"contact": map[string]interface{}{
					"email": "x@y.com",
					"subscription_topics": []interface{}{
						map[string]interface{}{"subscription_topic_id": float64(7), "domain_id": float64(4321), "status": "unsubscribed"},
						map[string]interface{}{"subscription_topic_id": float64(8), "domain_id": float64(4322), "status": "subscribed"},
					},
				}}
				if !reflect.DeepEqual(payload, want) {
					t.Errorf("expected body %v, got %v", want, payload)
				}

				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"data":{"id":"abc-123","email":"x@y.com","status":"subscribed"}}`))
			})
			defer cleanup()

			cmd := contacts.NewCmdContacts(f)
			cmd.SetArgs(args)
			cmd.SetOut(buf)

			if err := cmd.Execute(); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestContactsOmitSubscriptionTopicsWhenUnset(t *testing.T) {
	for _, args := range [][]string{
		{"create", "--email", "x@y.com"},
		{"update", "--id", "abc-123", "--email", "x@y.com"},
	} {
		t.Run(args[0], func(t *testing.T) {
			f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
				var payload map[string]map[string]interface{}
				json.NewDecoder(r.Body).Decode(&payload)
				if _, ok := payload["contact"]["subscription_topics"]; ok {
					t.Errorf("did not expect subscription_topics, got %v", payload["contact"])
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"data":{"id":"abc-123","email":"x@y.com"}}`))
			})
			defer cleanup()

			cmd := contacts.NewCmdContacts(f)
			cmd.SetArgs(args)
			cmd.SetOut(buf)

			if err := cmd.Execute(); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestContactsInvalidSubscriptionTopics(t *testing.T) {
	cases := []struct {
		name  string
		value string
		want  string
	}{
		{"not an array", `{"subscription_topic_id":7,"domain_id":4321,"status":"subscribed"}`, "must be a JSON array"},
		{"malformed", `[{"subscription_topic_id":7`, "invalid --subscription-topics JSON"},
		{"unknown key", `[{"subscription_topic_id":7,"domain_id":4321,"status":"subscribed","name":"News"}]`, `unknown field "name"`},
		{"trailing data", `[{"subscription_topic_id":7,"domain_id":4321,"status":"subscribed"}] []`, "unexpected data after the array"},
		{"string id", `[{"subscription_topic_id":"7","domain_id":4321,"status":"subscribed"}]`, "invalid --subscription-topics JSON"},
		{"missing topic", `[{"domain_id":4321,"status":"subscribed"}]`, "entry 1: subscription_topic_id must be greater than 0"},
		{"missing domain", `[{"subscription_topic_id":7,"domain_id":4321,"status":"subscribed"},{"subscription_topic_id":7,"status":"subscribed"}]`, "entry 2: domain_id must be greater than 0"},
		{"bad status", `[{"subscription_topic_id":7,"domain_id":4321,"status":"active"}]`, "entry 1: status must be one of subscribed, unsubscribed"},
		{"null entry", `[null]`, "entry 1: subscription_topic_id must be greater than 0"},
	}

	run := func(t *testing.T, args []string, want string) {
		f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
			t.Error("unexpected request")
		})
		defer cleanup()

		cmd := contacts.NewCmdContacts(f)
		cmd.SetArgs(args)
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true

		err := cmd.Execute()
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Fatalf("expected error containing %q, got: %v", want, err)
		}
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, []string{"create", "--email", "x@y.com", "--subscription-topics", tc.value}, tc.want)
		})
	}
	// update shares the parser; one case shows it is wired in.
	t.Run("update", func(t *testing.T) {
		run(t, []string{"update", "--id", "abc-123", "--email", "x@y.com", "--subscription-topics", `[{"subscription_topic_id":7}]`},
			"entry 1: domain_id must be greater than 0")
	})
}

func TestContactsJSONKeepsContactAsIs(t *testing.T) {
	withTopics := `{"id":"abc-123","email":"x@y.com","created_at":1700000000000,"updated_at":1700000000000,` +
		`"list_ids":[],"status":"subscribed","fields":{"first_name":null},` +
		`"subscription_topics":[{"subscription_topic_id":7,"domain_id":4321,"status":"unsubscribed"}]}`
	unreadable := `{"id":"abc-123","email":"x@y.com","list_ids":[],"status":"subscribed","fields":{},"subscription_topics":null}`
	topicsOff := `{"id":"abc-123","email":"x@y.com","list_ids":[],"status":"subscribed","fields":{}}`

	for _, tc := range []struct {
		name    string
		contact string
		args    []string
	}{
		{"get with topics", withTopics, []string{"get", "--id", "abc-123"}},
		{"get unreadable topics", unreadable, []string{"get", "--id", "abc-123"}},
		{"get topics off", topicsOff, []string{"get", "--id", "abc-123"}},
		{"create", withTopics, []string{"create", "--email", "x@y.com"}},
		{"update", withTopics, []string{"update", "--id", "abc-123", "--email", "x@y.com"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				io.WriteString(w, `{"data":`+tc.contact+`}`)
			})
			defer cleanup()

			viper.Set("output", "json")

			cmd := contacts.NewCmdContacts(f)
			cmd.SetArgs(tc.args)
			cmd.SetOut(buf)

			if err := cmd.Execute(); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			var got, want map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
				t.Fatalf("output is not valid JSON: %v\noutput:\n%s", err, buf.String())
			}
			json.Unmarshal([]byte(tc.contact), &want)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("expected the contact unchanged\nwant: %v\ngot:  %v", want, got)
			}
		})
	}
}

func TestContactsSubscriptionTopicsErrorsKeepMessage(t *testing.T) {
	cases := []struct {
		name   string
		status int
		body   string
		want   string
	}{
		{"topics off", http.StatusUnprocessableEntity, `{"errors":{"subscription_topics":["are not available for this account"]}}`, "are not available for this account"},
		{"topics not saved", http.StatusServiceUnavailable, `{"errors":"The contact was saved, but its subscription topic preferences were not. Send them again in an update of contact abc-123."}`, "Send them again in an update of contact abc-123"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				io.WriteString(w, tc.body)
			})
			defer cleanup()

			cmd := contacts.NewCmdContacts(f)
			cmd.SetArgs([]string{"create", "--email", "x@y.com", "--subscription-topics", `[{"subscription_topic_id":7,"domain_id":4321,"status":"subscribed"}]`})
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			err := cmd.Execute()
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected error containing %q, got: %v", tc.want, err)
			}
		})
	}
}
