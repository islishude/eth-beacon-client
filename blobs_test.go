package beaconclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestGetBlobs_Success(t *testing.T) {
	// Create a sample blob hex string (131072 bytes = 262144 hex chars)
	// For testing, we use a shorter representation
	blobHex := "0x" + strings.Repeat("00", 131072)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/eth/v1/beacon/blobs/head" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"execution_optimistic": false,
			"finalized": true,
			"data": ["` + blobHex + `"]
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	blobs, err := client.GetBlobs(context.Background(), "head")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if blobs.ExecutionOptimistic != false {
		t.Errorf("expected execution_optimistic false, got %v", blobs.ExecutionOptimistic)
	}
	if blobs.Finalized != true {
		t.Errorf("expected finalized true, got %v", blobs.Finalized)
	}
	if len(blobs.Data) != 1 {
		t.Errorf("expected 1 blob, got %d", len(blobs.Data))
	}
}

func TestGetBlobs_WithVersionedHashes(t *testing.T) {
	hash1 := common.HexToHash("0xabc123")
	hash2 := common.HexToHash("0xdef456")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hashes := r.URL.Query()["versioned_hashes"]
		if len(hashes) != 2 {
			t.Errorf("expected 2 versioned_hashes, got %d", len(hashes))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"execution_optimistic": false,
			"finalized": false,
			"data": []
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	blobs, err := client.GetBlobs(context.Background(), "12345", hash1, hash2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(blobs.Data) != 0 {
		t.Errorf("expected 0 blobs, got %d", len(blobs.Data))
	}
}

func TestGetBlobs_BlockNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code": 404, "message": "Block not found"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetBlobs(context.Background(), "999999999")
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
	if apiErr.Message != "Block not found" {
		t.Errorf("expected message 'Block not found', got %q", apiErr.Message)
	}
}

func TestGetBlobs_InvalidBlockID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"code": 400, "message": "Invalid block ID: current"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetBlobs(context.Background(), "current")
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
}

func TestGetBlobs_EmptyBlobs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"execution_optimistic": false,
			"finalized": true,
			"data": []
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	blobs, err := client.GetBlobs(context.Background(), "head")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(blobs.Data) != 0 {
		t.Errorf("expected 0 blobs, got %d", len(blobs.Data))
	}
}
