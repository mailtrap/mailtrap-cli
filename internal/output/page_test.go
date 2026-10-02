package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

var pageCols = []Column{{Header: "ID", Field: "id"}}

func TestPrintPage_JSONIsResponseBody(t *testing.T) {
	var buf bytes.Buffer
	body := json.RawMessage(`{"data":[{"id":"a","cc":[],"reply_to":null}],"total_count":7,"last_id":"a"}`)
	page := Page{Items: "data", Total: "total_count", Cursor: []string{"last_id"}, CursorFlag: "last-id"}

	if err := PrintPage(&buf, FormatJSON, body, page, pageCols); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got, want interface{}
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput:\n%s", err, buf.String())
	}
	json.Unmarshal(body, &want)
	if got, want := mustMarshal(t, got), mustMarshal(t, want); got != want {
		t.Errorf("JSON output = %s, want %s", got, want)
	}
}

func TestPrintPage_TableShowsTotalAndCursor(t *testing.T) {
	var buf bytes.Buffer
	body := json.RawMessage(`{"data":[{"id":"a"}],"total_count":7,"last_id":"a"}`)
	page := Page{Items: "data", Total: "total_count", Cursor: []string{"last_id"}, CursorFlag: "last-id"}

	if err := PrintPage(&buf, FormatTable, body, page, pageCols); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"ID", "a", "Total: 7", "Next page: --last-id a"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestPrintPage_NextArgsFollowCursor(t *testing.T) {
	var buf bytes.Buffer
	body := json.RawMessage(`{"data":[{"id":"a"}],"last_id":"a"}`)
	page := Page{Items: "data", Cursor: []string{"last_id"}, CursorFlag: "last-id", NextArgs: `--search "acme"`}

	if err := PrintPage(&buf, FormatTable, body, page, pageCols); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), `Next page: --last-id a --search "acme"`) {
		t.Errorf("expected next-page hint to repeat NextArgs, got:\n%s", buf.String())
	}
}

func TestPrintPage_NestedNumericCursor(t *testing.T) {
	var buf bytes.Buffer
	body := json.RawMessage(`{"data":[{"id":1}],"pagination":{"token":1,"next_token":2}}`)
	page := Page{Items: "data", Cursor: []string{"pagination", "next_token"}, CursorFlag: "token"}

	if err := PrintPage(&buf, FormatText, body, page, pageCols); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "Next page: --token 2") {
		t.Errorf("expected next-page token, got:\n%s", buf.String())
	}
}

func TestPrintPage_NoFooterOnLastPage(t *testing.T) {
	var buf bytes.Buffer
	body := json.RawMessage(`{"data":[],"last_id":null}`)
	page := Page{Items: "data", Cursor: []string{"last_id"}, CursorFlag: "last-id"}

	if err := PrintPage(&buf, FormatTable, body, page, pageCols); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(buf.String(), "Next page") {
		t.Errorf("expected no next-page hint without a cursor, got:\n%s", buf.String())
	}
}

func TestPrintPage_NoFooterWithoutCursorOrTotal(t *testing.T) {
	var buf bytes.Buffer
	body := json.RawMessage(`{"data":[{"id":"a"}],"total_count":7,"last_id":"a"}`)
	page := Page{Items: "data"}

	if err := PrintPage(&buf, FormatTable, body, page, pageCols); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "Next page") || strings.Contains(out, "Total") {
		t.Errorf("expected no footer for a page without Cursor or Total, got:\n%s", out)
	}
}

func mustMarshal(t *testing.T, v interface{}) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(b)
}
