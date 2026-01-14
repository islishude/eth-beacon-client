package beaconclient

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
)

func TestGetBlobs_WithVersionedHashes(t *testing.T) {
	reqHashes := []common.Hash{
		common.HexToHash("0x019064de4167396f5930109e6027c41b5194597507f218a8c864dd4ed1593ac0"),
		common.HexToHash("0x015aac334f58c4e289af77d6b20ca015df30fec70988c450bc904a34301da3d9"),
	}
	var blockId int64 = 9388204
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != fmt.Sprintf("/eth/v1/beacon/blobs/%d", blockId) {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}

		hashes := r.URL.Query()["versioned_hashes"]
		if len(hashes) != len(reqHashes) {
			t.Errorf("expected %d versioned_hashes, got %d", len(reqHashes), len(hashes))
		}

		for i, h := range hashes {
			if reqHashes[i].Hex() != h {
				t.Errorf("unexpected versioned_hash[%d]: %s", i, h)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		data, err := os.ReadFile("testdata/blob_success.json")
		if err != nil {
			t.Fatalf("failed to read test data: %v", err)
		}
		_, _ = w.Write(data)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	blobs, err := client.GetBlobs(t.Context(),
		strconv.FormatInt(blockId, 10), reqHashes...)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if blobs.ExecutionOptimistic != false {
		t.Errorf("expected execution_optimistic false, got %v", blobs.ExecutionOptimistic)
	}
	if blobs.Finalized != true {
		t.Errorf("expected finalized true, got %v", blobs.Finalized)
	}

	if len(blobs.Data) != len(reqHashes) {
		t.Errorf("expected %d blobs, got %d", len(reqHashes), len(blobs.Data))
	}

	for i, blob := range blobs.Data {
		commitment, err := kzg4844.BlobToCommitment(&blob)
		if err != nil {
			t.Errorf("commitment error: %v", err)
		}
		blobHash := common.Hash(kzg4844.CalcBlobHashV1(sha256.New(), &commitment))
		if blobHash != reqHashes[i] {
			t.Errorf("unexpected blob hash: got %s, want %s", blobHash.Hex(), reqHashes[i].Hex())
		}
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
