package beaconclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		expected string
	}{
		{"without trailing slash", "http://localhost:5052", "http://localhost:5052"},
		{"with trailing slash", "http://localhost:5052/", "http://localhost:5052"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.baseURL)
			if client.baseURL != tt.expected {
				t.Errorf("expected baseURL %q, got %q", tt.expected, client.baseURL)
			}
		})
	}
}

func TestAPIError_Error(t *testing.T) {
	err := &APIError{Code: 404, Message: "Block not found"}
	expected := "beacon API error (code 404): Block not found"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestDoRequest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/eth/v1/beacon/genesis" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("unexpected Accept header: %s", r.Header.Get("Accept"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data": {"test": "value"}}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	body, err := client.doRequest(context.Background(), http.MethodGet, "/eth/v1/beacon/genesis", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// doRequest now returns raw body, each endpoint parses it independently
	var result struct {
		Data map[string]string `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("failed to unmarshal body: %v", err)
	}
	if result.Data["test"] != "value" {
		t.Errorf("expected test=value, got %v", result.Data)
	}
}

func TestDoRequest_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"code": 400, "message": "Invalid block ID: current"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.doRequest(context.Background(), http.MethodGet, "/eth/v1/beacon/blobs/current", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Code != 400 {
		t.Errorf("expected code 400, got %d", apiErr.Code)
	}
	if apiErr.Message != "Invalid block ID: current" {
		t.Errorf("expected message 'Invalid block ID: current', got %q", apiErr.Message)
	}
}

func TestDoRequest_NotFoundError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code": 404, "message": "Block not found"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.doRequest(context.Background(), http.MethodGet, "/eth/v1/beacon/blobs/999999999", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Code != 404 {
		t.Errorf("expected code 404, got %d", apiErr.Code)
	}
}

func TestDoRequest_WithQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("foo") != "bar" {
			t.Errorf("expected query param foo=bar, got %s", r.URL.Query().Get("foo"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data": {}}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	query := make(map[string][]string)
	query["foo"] = []string{"bar"}
	_, err := client.doRequest(context.Background(), http.MethodGet, "/test", query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDoRequest_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`not valid json`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	body, err := client.doRequest(context.Background(), http.MethodGet, "/test", nil)
	// doRequest now returns raw body without parsing, so no error here
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// The body should be the invalid JSON string
	if string(body) != "not valid json" {
		t.Errorf("expected 'not valid json', got %q", string(body))
	}
}

func TestDoRequest_NonJSONError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`Internal Server Error`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.doRequest(context.Background(), http.MethodGet, "/test", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	// Should contain the status code and body
	if _, ok := err.(*APIError); ok {
		t.Error("should not be APIError for non-JSON response")
	}
}
