package forwardrules_test

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
	"github.com/mailtrap/mailtrap-cli/internal/commands/inbound/forwardrules"
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

func sampleRule(id int, name string) map[string]interface{} {
	return map[string]interface{}{
		"id":         id,
		"name":       name,
		"created_at": "2026-05-08T10:30:00.000Z",
		"updated_at": "2026-05-08T10:30:00.000Z",
		"conditions": []map[string]interface{}{
			{"match_type": "sender", "operator": "ends_with", "value": "@billing.example.com", "header_key": nil},
		},
		"destinations": []map[string]interface{}{
			{"email": "finance@example.com"},
		},
	}
}

func readBody(t *testing.T, r *http.Request) map[string]interface{} {
	t.Helper()
	body, _ := io.ReadAll(r.Body)
	var reqBody map[string]interface{}
	if err := json.Unmarshal(body, &reqBody); err != nil {
		t.Fatalf("request body is not valid JSON: %v\nbody: %s", err, body)
	}
	return reqBody
}

func TestForwardRulesList(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/inbound/inboxes/201/forward_rules") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{sampleRule(7, "Copy billing mail to finance")},
		})
	})
	defer cleanup()

	cmd := forwardrules.NewCmdForwardRules(f)
	cmd.SetArgs([]string{"list", "--inbox-id", "201"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"CONDITIONS", "DESTINATIONS", "Copy billing mail to finance", "finance@example.com", "ends_with"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestForwardRulesListJSON(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{sampleRule(7, "Copy billing mail to finance")},
		})
	})
	defer cleanup()

	viper.Set("output", "json")

	cmd := forwardrules.NewCmdForwardRules(f)
	cmd.SetArgs([]string{"list", "--inbox-id", "201"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result []forwardrules.ForwardRule
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput:\n%s", err, buf.String())
	}
	if len(result) != 1 || result[0].ID != 7 {
		t.Fatalf("unexpected JSON result: %v", result)
	}
	if len(result[0].Conditions) != 1 || result[0].Conditions[0].Operator != "ends_with" {
		t.Errorf("unexpected conditions: %+v", result[0].Conditions)
	}
	if len(result[0].Destinations) != 1 || result[0].Destinations[0].Email != "finance@example.com" {
		t.Errorf("unexpected destinations: %+v", result[0].Destinations)
	}
}

func TestForwardRulesListMissingInboxID(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	cmd := forwardrules.NewCmdForwardRules(f)
	cmd.SetArgs([]string{"list"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when --inbox-id is missing")
	}
	if !strings.Contains(err.Error(), "--inbox-id is required") {
		t.Errorf("expected '--inbox-id is required' error, got: %v", err)
	}
}

func TestForwardRulesGet(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/api/inbound/inboxes/201/forward_rules/7") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": sampleRule(7, "Copy billing mail to finance")})
	})
	defer cleanup()

	cmd := forwardrules.NewCmdForwardRules(f)
	cmd.SetArgs([]string{"get", "--inbox-id", "201", "--id", "7"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "Copy billing mail to finance") {
		t.Errorf("expected output to contain rule name, got:\n%s", buf.String())
	}
}

func TestForwardRulesGetMissingID(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	cmd := forwardrules.NewCmdForwardRules(f)
	cmd.SetArgs([]string{"get", "--inbox-id", "201"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when --id is missing")
	}
	if !strings.Contains(err.Error(), "--id is required") {
		t.Errorf("expected '--id is required' error, got: %v", err)
	}
}

func TestForwardRulesCreate(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/inbound/inboxes/201/forward_rules") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		reqBody := readBody(t, r)
		if _, wrapped := reqBody["forward_rule"]; wrapped {
			t.Errorf("request body must not be wrapped, got %v", reqBody)
		}
		if reqBody["name"] != "Copy billing mail to finance" {
			t.Errorf("unexpected name: %v", reqBody["name"])
		}

		conditions, _ := reqBody["conditions"].([]interface{})
		if len(conditions) != 2 {
			t.Fatalf("expected 2 conditions, got %v", reqBody["conditions"])
		}
		sender := conditions[0].(map[string]interface{})
		if sender["match_type"] != "sender" || sender["operator"] != "ends_with" || sender["value"] != "@billing.example.com" {
			t.Errorf("unexpected sender condition: %v", sender)
		}
		if _, ok := sender["header_key"]; ok {
			t.Errorf("header_key must be omitted when not given, got %v", sender)
		}
		header := conditions[1].(map[string]interface{})
		if header["header_key"] != "X-Priority" || header["operator"] != "not_empty" {
			t.Errorf("unexpected header condition: %v", header)
		}
		if _, ok := header["value"]; ok {
			t.Errorf("value must be omitted when not given, got %v", header)
		}

		destinations, _ := reqBody["destinations"].([]interface{})
		if len(destinations) != 2 {
			t.Fatalf("expected 2 destinations, got %v", reqBody["destinations"])
		}
		if destinations[1].(map[string]interface{})["email"] != "accounting@example.com" {
			t.Errorf("unexpected destinations: %v", destinations)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{"data": sampleRule(7, "Copy billing mail to finance")})
	})
	defer cleanup()

	cmd := forwardrules.NewCmdForwardRules(f)
	cmd.SetArgs([]string{
		"create", "--inbox-id", "201", "--name", "Copy billing mail to finance",
		"--conditions", `[{"match_type":"sender","operator":"ends_with","value":"@billing.example.com"},{"match_type":"header","operator":"not_empty","header_key":"X-Priority"}]`,
		"--destinations", "finance@example.com",
		"--destinations", "accounting@example.com",
	})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "Copy billing mail to finance") {
		t.Errorf("expected output to contain rule name, got:\n%s", buf.String())
	}
}

func TestForwardRulesCreateNameOnly(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		reqBody := readBody(t, r)
		if len(reqBody) != 1 || reqBody["name"] != "Catch all" {
			t.Errorf("expected only name in body, got %v", reqBody)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{"data": sampleRule(8, "Catch all")})
	})
	defer cleanup()

	cmd := forwardrules.NewCmdForwardRules(f)
	cmd.SetArgs([]string{"create", "--inbox-id", "201", "--name", "Catch all"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestForwardRulesCreateMissingName(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		t.Error("no request expected")
	})
	defer cleanup()

	cmd := forwardrules.NewCmdForwardRules(f)
	cmd.SetArgs([]string{"create", "--inbox-id", "201"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when --name is missing")
	}
	if !strings.Contains(err.Error(), "--name is required") {
		t.Errorf("expected '--name is required' error, got: %v", err)
	}
}

func TestForwardRulesCreateInvalidConditions(t *testing.T) {
	cases := map[string]struct {
		conditions string
		wantErr    string
	}{
		"not an array":     {`{"match_type":"sender"}`, "must be a JSON array"},
		"malformed":        {`[{"match_type":`, "invalid --conditions JSON"},
		"unknown key":      {`[{"match_type":"sender","operator":"equal","value":"a","typo":1}]`, "unknown field"},
		"bad match_type":   {`[{"match_type":"subject","operator":"equal","value":"a"}]`, "match_type must be one of"},
		"bad operator":     {`[{"match_type":"sender","operator":"regex","value":"a"}]`, "operator must be one of"},
		"missing operator": {`[{"match_type":"sender","value":"a"}]`, "operator must be one of"},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
				t.Error("no request expected")
			})
			defer cleanup()

			cmd := forwardrules.NewCmdForwardRules(f)
			cmd.SetArgs([]string{"create", "--inbox-id", "201", "--name", "Rule", "--conditions", tc.conditions})
			cmd.SetOut(buf)

			err := cmd.Execute()
			if err == nil {
				t.Fatal("expected error for invalid --conditions")
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("expected error containing %q, got: %v", tc.wantErr, err)
			}
		})
	}
}

func TestForwardRulesUpdatePartial(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/inbound/inboxes/201/forward_rules/7") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		reqBody := readBody(t, r)
		if len(reqBody) != 1 {
			t.Errorf("expected only destinations in body, got %v", reqBody)
		}
		destinations, _ := reqBody["destinations"].([]interface{})
		if len(destinations) != 2 {
			t.Errorf("expected 2 destinations, got %v", reqBody["destinations"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": sampleRule(7, "Copy billing mail to finance")})
	})
	defer cleanup()

	cmd := forwardrules.NewCmdForwardRules(f)
	cmd.SetArgs([]string{"update", "--inbox-id", "201", "--id", "7", "--destinations", "finance@example.com,accounting@example.com"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "Copy billing mail to finance") {
		t.Errorf("expected output to contain rule name, got:\n%s", buf.String())
	}
}

func TestForwardRulesUpdateClearsSets(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)

		var reqBody map[string]json.RawMessage
		if err := json.Unmarshal(body, &reqBody); err != nil {
			t.Fatalf("request body is not valid JSON: %v", err)
		}
		if string(reqBody["conditions"]) != "[]" {
			t.Errorf("expected conditions to be sent as [], got %s", reqBody["conditions"])
		}
		if string(reqBody["destinations"]) != "[]" {
			t.Errorf("expected destinations to be sent as [], got %s", reqBody["destinations"])
		}
		if string(reqBody["name"]) != `"Renamed"` {
			t.Errorf("unexpected name: %s", reqBody["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{
			"id": 7, "name": "Renamed", "conditions": []interface{}{}, "destinations": []interface{}{},
		}})
	})
	defer cleanup()

	cmd := forwardrules.NewCmdForwardRules(f)
	cmd.SetArgs([]string{"update", "--inbox-id", "201", "--id", "7", "--name", "Renamed", "--conditions", "[]", "--destinations", ""})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestForwardRulesUpdateMissingID(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	cmd := forwardrules.NewCmdForwardRules(f)
	cmd.SetArgs([]string{"update", "--inbox-id", "201", "--name", "Renamed"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when --id is missing")
	}
	if !strings.Contains(err.Error(), "--id is required") {
		t.Errorf("expected '--id is required' error, got: %v", err)
	}
}

func TestForwardRulesDelete(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/inbound/inboxes/201/forward_rules/7") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer cleanup()

	cmd := forwardrules.NewCmdForwardRules(f)
	cmd.SetArgs([]string{"delete", "--inbox-id", "201", "--id", "7"})
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "deleted successfully") {
		t.Errorf("expected success message, got:\n%s", buf.String())
	}
}
