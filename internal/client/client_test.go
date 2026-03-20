package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// helper: creates a Client whose base URL points at the given test server.
// Since the Client uses BaseURL as a prefix, we return the server URL as a BaseURL.
func setupTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Client, BaseURL) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	c := New("test-api-token")
	return srv, c, BaseURL(srv.URL)
}

// --- 1. Client creation ---

func TestNew(t *testing.T) {
	c := New("my-token")
	if c == nil {
		t.Fatal("expected non-nil client")
	}
	if c.apiToken != "my-token" {
		t.Fatalf("expected apiToken %q, got %q", "my-token", c.apiToken)
	}
	if c.httpClient == nil {
		t.Fatal("expected non-nil httpClient")
	}
}

// --- 2. GET request - success with JSON response ---

func TestGet_Success(t *testing.T) {
	type respBody struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	_, c, base := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/items/1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(respBody{ID: 1, Name: "item-one"})
	})

	var result respBody
	err := c.Get(context.Background(), base, "/api/items/1", nil, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}
	if result.Name != "item-one" {
		t.Errorf("expected Name %q, got %q", "item-one", result.Name)
	}
}

// --- 3. GET request - with query params ---

func TestGet_WithQueryParams(t *testing.T) {
	_, c, base := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		qPage := r.URL.Query().Get("page")
		qLimit := r.URL.Query().Get("limit")
		if qPage != "2" {
			t.Errorf("expected page=2, got %q", qPage)
		}
		if qLimit != "10" {
			t.Errorf("expected limit=10, got %q", qLimit)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"total":100}`))
	})

	type listResp struct {
		Total int `json:"total"`
	}

	q := url.Values{}
	q.Set("page", "2")
	q.Set("limit", "10")

	var result listResp
	err := c.Get(context.Background(), base, "/api/items", q, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 100 {
		t.Errorf("expected Total 100, got %d", result.Total)
	}
}

// --- 4. POST request - success with JSON body and response ---

func TestPost_Success(t *testing.T) {
	type reqBody struct {
		Name string `json:"name"`
	}
	type respBody struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	_, c, base := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		var req reqBody
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}
		if req.Name != "new-item" {
			t.Errorf("expected name %q, got %q", "new-item", req.Name)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(respBody{ID: 42, Name: req.Name})
	})

	var result respBody
	err := c.Post(context.Background(), base, "/api/items", reqBody{Name: "new-item"}, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != 42 {
		t.Errorf("expected ID 42, got %d", result.ID)
	}
	if result.Name != "new-item" {
		t.Errorf("expected Name %q, got %q", "new-item", result.Name)
	}
}

// --- 5. PATCH request ---

func TestPatch_Success(t *testing.T) {
	type patchReq struct {
		Name string `json:"name"`
	}
	type patchResp struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	_, c, base := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/items/5" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var req patchReq
		json.Unmarshal(body, &req)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(patchResp{ID: 5, Name: req.Name})
	})

	var result patchResp
	err := c.Patch(context.Background(), base, "/api/items/5", patchReq{Name: "updated"}, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != 5 || result.Name != "updated" {
		t.Errorf("unexpected result: %+v", result)
	}
}

// --- 6. PUT request ---

func TestPut_Success(t *testing.T) {
	type putReq struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	type putResp struct {
		OK bool `json:"ok"`
	}

	_, c, base := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		var req putReq
		json.Unmarshal(body, &req)
		if req.Name != "replaced" || req.Value != 99 {
			t.Errorf("unexpected request body: %+v", req)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(putResp{OK: true})
	})

	var result putResp
	err := c.Put(context.Background(), base, "/api/items/7", putReq{Name: "replaced", Value: 99}, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.OK {
		t.Error("expected OK true")
	}
}

// --- 7. DELETE request ---

func TestDelete_Success(t *testing.T) {
	_, c, base := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/items/3" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.Delete(context.Background(), base, "/api/items/3", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDelete_WithResponseBody(t *testing.T) {
	type deleteResp struct {
		Deleted bool `json:"deleted"`
	}

	_, c, base := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(deleteResp{Deleted: true})
	})

	var result deleteResp
	err := c.Delete(context.Background(), base, "/api/items/3", &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Deleted {
		t.Error("expected Deleted true")
	}
}

// --- 8. API error handling (4xx, 5xx status codes) ---

func TestAPIError_4xx_PlainText(t *testing.T) {
	_, c, base := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	})

	err := c.Get(context.Background(), base, "/api/missing", nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("expected status 404, got %d", apiErr.StatusCode)
	}
	// When body is not valid JSON, Message should contain the raw body text
	if apiErr.Message != "not found" {
		t.Errorf("expected message %q, got %q", "not found", apiErr.Message)
	}
}

func TestAPIError_5xx(t *testing.T) {
	_, c, base := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	})

	err := c.Post(context.Background(), base, "/api/items", map[string]string{"a": "b"}, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 500 {
		t.Errorf("expected status 500, got %d", apiErr.StatusCode)
	}
}

// --- 9. API error with JSON error body ---

func TestAPIError_JSONBody(t *testing.T) {
	_, c, base := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  "validation failed",
			"errors": map[string]string{"name": "is required"},
		})
	})

	err := c.Post(context.Background(), base, "/api/items", map[string]string{}, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 422 {
		t.Errorf("expected status 422, got %d", apiErr.StatusCode)
	}
	if apiErr.Message != "validation failed" {
		t.Errorf("expected message %q, got %q", "validation failed", apiErr.Message)
	}
	if v, ok := apiErr.Errors["name"]; !ok || v != "is required" {
		t.Errorf("expected errors[name]=%q, got %v", "is required", apiErr.Errors)
	}
	// Check Error() string output
	errStr := apiErr.Error()
	if !strings.Contains(errStr, "422") || !strings.Contains(errStr, "validation failed") {
		t.Errorf("unexpected Error() output: %s", errStr)
	}
}

func TestAPIError_EmptyMessage(t *testing.T) {
	apiErr := &APIError{StatusCode: 403}
	expected := "API error 403"
	if apiErr.Error() != expected {
		t.Errorf("expected %q, got %q", expected, apiErr.Error())
	}
}

// --- 10. GetRaw - success returning raw bytes ---

func TestGetRaw_Success(t *testing.T) {
	rawContent := "raw email content\nSubject: Hello\n\nBody here"

	_, c, base := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(rawContent))
	})

	data, err := c.GetRaw(context.Background(), base, "/api/messages/1/raw", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != rawContent {
		t.Errorf("expected %q, got %q", rawContent, string(data))
	}
}

func TestGetRaw_WithQueryParams(t *testing.T) {
	_, c, base := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("format") != "eml" {
			t.Errorf("expected format=eml, got %q", r.URL.Query().Get("format"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("eml-data"))
	})

	q := url.Values{}
	q.Set("format", "eml")
	data, err := c.GetRaw(context.Background(), base, "/api/messages/1", q)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != "eml-data" {
		t.Errorf("expected %q, got %q", "eml-data", string(data))
	}
}

// --- 11. GetRaw - error handling ---

func TestGetRaw_Error(t *testing.T) {
	_, c, base := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("access denied"))
	})

	data, err := c.GetRaw(context.Background(), base, "/api/messages/99/raw", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if data != nil {
		t.Errorf("expected nil data on error, got %v", data)
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 403 {
		t.Errorf("expected status 403, got %d", apiErr.StatusCode)
	}
	if apiErr.Message != "access denied" {
		t.Errorf("expected message %q, got %q", "access denied", apiErr.Message)
	}
}

// --- 12. Auth header is set correctly (Api-Token) ---

func TestAuthHeader(t *testing.T) {
	const token = "secret-token-12345"

	_, c, base := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		got := r.Header.Get("Api-Token")
		if got != token {
			t.Errorf("expected Api-Token header %q, got %q", token, got)
		}
		w.WriteHeader(http.StatusOK)
	})

	c.apiToken = token
	err := c.Get(context.Background(), base, "/api/check", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAuthHeader_GetRaw(t *testing.T) {
	const token = "raw-token-xyz"

	_, c, base := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		got := r.Header.Get("Api-Token")
		if got != token {
			t.Errorf("expected Api-Token header %q, got %q", token, got)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	c.apiToken = token
	_, err := c.GetRaw(context.Background(), base, "/api/raw", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- 13. Content-Type header is set correctly ---

func TestContentTypeHeader(t *testing.T) {
	_, c, base := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("expected Content-Type %q, got %q", "application/json", ct)
		}
		accept := r.Header.Get("Accept")
		if accept != "application/json" {
			t.Errorf("expected Accept %q, got %q", "application/json", accept)
		}
		w.WriteHeader(http.StatusOK)
	})

	err := c.Post(context.Background(), base, "/api/test", map[string]string{"k": "v"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- Additional edge cases ---

func TestGet_NilResult(t *testing.T) {
	_, c, base := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":1}`))
	})

	// Passing nil result should not error even when body is returned
	err := c.Get(context.Background(), base, "/api/items/1", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGet_EmptyResponseBody(t *testing.T) {
	_, c, base := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	type resp struct {
		ID int `json:"id"`
	}
	var result resp
	err := c.Get(context.Background(), base, "/api/items/1", nil, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// result should remain zero value since body is empty
	if result.ID != 0 {
		t.Errorf("expected zero value, got %+v", result)
	}
}

func TestPost_NilBody(t *testing.T) {
	_, c, base := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if len(body) != 0 {
			t.Errorf("expected empty body for nil input, got %q", string(body))
		}
		w.WriteHeader(http.StatusOK)
	})

	err := c.Post(context.Background(), base, "/api/trigger", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAPIError_401(t *testing.T) {
	_, c, base := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
	})

	err := c.Get(context.Background(), base, "/api/protected", nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 401 {
		t.Errorf("expected status 401, got %d", apiErr.StatusCode)
	}
	if apiErr.Message != "unauthorized" {
		t.Errorf("expected message %q, got %q", "unauthorized", apiErr.Message)
	}
}

func TestCancelledContext(t *testing.T) {
	_, c, base := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	err := c.Get(ctx, base, "/api/items", nil, nil)
	if err == nil {
		t.Fatal("expected error for cancelled context, got nil")
	}
}
