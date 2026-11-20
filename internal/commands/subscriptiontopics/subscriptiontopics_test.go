package subscriptiontopics_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/commands/subscriptiontopics"
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

func listBody() map[string]interface{} {
	return map[string]interface{}{
		"data": []map[string]interface{}{
			{"id": 7, "name": "Product updates", "description": "News about the product", "domain_id": 4321, "deleted_at": nil, "created_at": "2026-11-01T10:00:00Z", "updated_at": "2026-11-02T10:00:00Z"},
			{"id": 8, "name": "Newsletter", "description": nil, "domain_id": nil, "deleted_at": nil, "created_at": "2026-11-03T10:00:00Z", "updated_at": "2026-11-03T10:00:00Z"},
		},
		"pagination": map[string]interface{}{"token": 1, "prev_token": nil, "next_token": 2},
	}
}

func TestSubscriptionTopicsList(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		// Token-scoped: no /api/accounts/{id} prefix.
		if r.URL.Path != "/api/subscription_topics" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Api-Token") != "test-token" {
			t.Errorf("expected Api-Token header 'test-token', got %q", r.Header.Get("Api-Token"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(listBody())
	})
	defer cleanup()

	cmd := subscriptiontopics.NewCmdSubscriptionTopics(f)
	cmd.SetArgs([]string{"list"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	for _, want := range []string{"DOMAIN ID", "Product updates", "News about the product", "4321", "Newsletter"} {
		if !strings.Contains(output, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, output)
		}
	}
}

func TestSubscriptionTopicsListJSON(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(listBody())
	})
	defer cleanup()

	viper.Set("output", "json")

	cmd := subscriptiontopics.NewCmdSubscriptionTopics(f)
	cmd.SetArgs([]string{"list"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result struct {
		Data       []map[string]interface{} `json:"data"`
		Pagination map[string]interface{}   `json:"pagination"`
	}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not a JSON object: %v\noutput:\n%s", err, buf.String())
	}
	if len(result.Data) != 2 {
		t.Fatalf("expected 2 topics, got %d", len(result.Data))
	}
	if result.Data[0]["name"] != "Product updates" {
		t.Errorf("expected name 'Product updates', got %v", result.Data[0]["name"])
	}
	if domainID, ok := result.Data[1]["domain_id"]; !ok || domainID != nil {
		t.Errorf("expected domain_id null to be kept, got %v (present: %v)", domainID, ok)
	}
	if result.Pagination["next_token"] != float64(2) {
		t.Errorf("expected pagination.next_token 2, got %v", result.Pagination["next_token"])
	}
}

func TestSubscriptionTopicsListNextPage(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(listBody())
	})
	defer cleanup()

	cmd := subscriptiontopics.NewCmdSubscriptionTopics(f)
	cmd.SetArgs([]string{"list"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "Next page: --token 2") {
		t.Errorf("expected next page footer, got:\n%s", buf.String())
	}
}

func TestSubscriptionTopicsListNextPageRunsWithSamePerPage(t *testing.T) {
	var requests []string
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.RawQuery)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(listBody())
	})
	defer cleanup()

	cmd := subscriptiontopics.NewCmdSubscriptionTopics(f)
	cmd.SetArgs([]string{"list", "--per-page", "1"})
	cmd.SetOut(buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, hint, found := strings.Cut(buf.String(), "Next page: ")
	if !found {
		t.Fatalf("expected next page footer, got:\n%s", buf.String())
	}

	cmd = subscriptiontopics.NewCmdSubscriptionTopics(f)
	cmd.SetArgs(append([]string{"list"}, strings.Fields(hint)...))
	cmd.SetOut(buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error running the footer's command: %v", err)
	}

	if want := []string{"per_page=1", "per_page=1&token=2"}; !reflect.DeepEqual(requests, want) {
		t.Errorf("expected queries %v, got %v", want, requests)
	}
}

func TestSubscriptionTopicsListQuery(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("per_page"); got != "10" {
			t.Errorf("expected per_page=10, got %q", got)
		}
		if got := r.URL.Query().Get("token"); got != "2" {
			t.Errorf("expected token=2, got %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(listBody())
	})
	defer cleanup()

	cmd := subscriptiontopics.NewCmdSubscriptionTopics(f)
	cmd.SetArgs([]string{"list", "--per-page", "10", "--token", "2"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSubscriptionTopicsListOmitsUnsetQuery(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Errorf("expected no query, got %q", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(listBody())
	})
	defer cleanup()

	cmd := subscriptiontopics.NewCmdSubscriptionTopics(f)
	cmd.SetArgs([]string{"list"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSubscriptionTopicsListNoAccountID(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(listBody())
	})
	defer cleanup()

	viper.Set("account-id", "")

	cmd := subscriptiontopics.NewCmdSubscriptionTopics(f)
	cmd.SetArgs([]string{"list"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected list to work without an account ID, got: %v", err)
	}
}

// Accounts without subscription topics get a 403 with an empty body.
func TestSubscriptionTopicsListForbidden(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})
	defer cleanup()

	cmd := subscriptiontopics.NewCmdSubscriptionTopics(f)
	cmd.SetArgs([]string{"list"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	var apiErr *client.APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusForbidden {
		t.Fatalf("expected an API error 403, got: %v", err)
	}
}

func TestSubscriptionTopicsGet(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/subscription_topics/7" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"data":{"id":7,"name":"Product updates","description":null,"domain_id":null,`+
			`"deleted_at":"2026-11-10T10:00:00Z","created_at":"2026-11-01T10:00:00Z","updated_at":"2026-11-10T10:00:00Z"}}`)
	})
	defer cleanup()

	cmd := subscriptiontopics.NewCmdSubscriptionTopics(f)
	cmd.SetArgs([]string{"get", "--id", "7"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	for _, want := range []string{"DELETED AT", "Product updates", "2026-11-10T10:00:00Z"} {
		if !strings.Contains(output, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, output)
		}
	}
}

func TestSubscriptionTopicsGetJSONKeepsTopicAsIs(t *testing.T) {
	topic := `{"id":7,"name":"Product updates","description":null,"domain_id":4321,` +
		`"deleted_at":null,"created_at":"2026-11-01T10:00:00Z","updated_at":"2026-11-02T10:00:00Z"}`

	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"data":`+topic+`}`)
	})
	defer cleanup()

	viper.Set("output", "json")

	cmd := subscriptiontopics.NewCmdSubscriptionTopics(f)
	cmd.SetArgs([]string{"get", "--id", "7"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got, want map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput:\n%s", err, buf.String())
	}
	json.Unmarshal([]byte(topic), &want)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected the topic unchanged\nwant: %v\ngot:  %v", want, got)
	}
}

func TestSubscriptionTopicsGetMissingID(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		t.Error("unexpected request")
	})
	defer cleanup()

	cmd := subscriptiontopics.NewCmdSubscriptionTopics(f)
	cmd.SetArgs([]string{"get"})
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), `required flag(s) "id" not set`) {
		t.Fatalf("expected a required --id error, got: %v", err)
	}
}
