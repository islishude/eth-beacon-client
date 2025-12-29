package beacon

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestGetBlock_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/eth/v2/beacon/blocks/head" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Eth-Consensus-Version", "deneb")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"version": "deneb",
			"execution_optimistic": false,
			"finalized": true,
			"data": {
				"message": {
					"slot": "12345",
					"proposer_index": "100",
					"parent_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"state_root": "0xaabbccdd9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"body": {
						"randao_reveal": "0x1234",
						"eth1_data": {
							"deposit_root": "0x0000000000000000000000000000000000000000000000000000000000000000",
							"deposit_count": "100",
							"block_hash": "0x0000000000000000000000000000000000000000000000000000000000000000"
						},
						"graffiti": "0x0000000000000000000000000000000000000000000000000000000000000000"
					}
				},
				"signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	resp, err := client.GetBlock(context.Background(), "head")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Version != ConsensusVersionDeneb {
		t.Errorf("expected version deneb, got %s", resp.Version)
	}
	if resp.ExecutionOptimistic {
		t.Error("expected execution_optimistic false")
	}
	if !resp.Finalized {
		t.Error("expected finalized true")
	}
	if resp.Data.Message.Slot != 12345 {
		t.Errorf("expected slot 12345, got %d", resp.Data.Message.Slot)
	}
	if resp.Data.Message.ProposerIndex != 100 {
		t.Errorf("expected proposer_index 100, got %d", resp.Data.Message.ProposerIndex)
	}
	expectedParentRoot := common.HexToHash("0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2")
	if resp.Data.Message.ParentRoot != expectedParentRoot {
		t.Errorf("unexpected parent_root: %s", resp.Data.Message.ParentRoot.Hex())
	}
}

func TestGetBlock_BySlot(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/eth/v2/beacon/blocks/12345" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"version": "electra",
			"execution_optimistic": true,
			"finalized": false,
			"data": {
				"message": {
					"slot": "12345",
					"proposer_index": "200",
					"parent_root": "0x0000000000000000000000000000000000000000000000000000000000000000",
					"state_root": "0x0000000000000000000000000000000000000000000000000000000000000000",
					"body": {}
				},
				"signature": "0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	resp, err := client.GetBlock(context.Background(), "12345")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Version != ConsensusVersionElectra {
		t.Errorf("expected version electra, got %s", resp.Version)
	}
	if !resp.ExecutionOptimistic {
		t.Error("expected execution_optimistic true")
	}
	if resp.Finalized {
		t.Error("expected finalized false")
	}
}

func TestGetBlock_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code": 404, "message": "Block not found"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetBlock(context.Background(), "999999999")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.Code != 404 {
		t.Errorf("expected error code 404, got %d", apiErr.Code)
	}
}

func TestGetBlockRoot_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/eth/v1/beacon/blocks/head/root" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"execution_optimistic": false,
			"finalized": true,
			"data": {
				"root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	resp, err := client.GetBlockRoot(context.Background(), "head")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ExecutionOptimistic {
		t.Error("expected execution_optimistic false")
	}
	if !resp.Finalized {
		t.Error("expected finalized true")
	}
	expectedRoot := common.HexToHash("0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2")
	if resp.Data.Root != expectedRoot {
		t.Errorf("unexpected root: %s", resp.Data.Root.Hex())
	}
}

func TestGetBlockRoot_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code": 404, "message": "Block not found"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetBlockRoot(context.Background(), "999999999")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetBlockHeader_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/eth/v1/beacon/headers/finalized" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"execution_optimistic": false,
			"finalized": true,
			"data": {
				"root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
				"canonical": true,
				"header": {
					"message": {
						"slot": "12345",
						"proposer_index": "100",
						"parent_root": "0x0000000000000000000000000000000000000000000000000000000000000001",
						"state_root": "0x0000000000000000000000000000000000000000000000000000000000000002",
						"body_root": "0x0000000000000000000000000000000000000000000000000000000000000003"
					},
					"signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
				}
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	resp, err := client.GetBlockHeader(context.Background(), "finalized")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !resp.Finalized {
		t.Error("expected finalized true")
	}
	if !resp.Data.Canonical {
		t.Error("expected canonical true")
	}
	if resp.Data.Header.Message.Slot != 12345 {
		t.Errorf("expected slot 12345, got %d", resp.Data.Header.Message.Slot)
	}
	if resp.Data.Header.Message.ProposerIndex != 100 {
		t.Errorf("expected proposer_index 100, got %d", resp.Data.Header.Message.ProposerIndex)
	}

	expectedRoot := common.HexToHash("0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2")
	if resp.Data.Root != expectedRoot {
		t.Errorf("unexpected root: %s", resp.Data.Root.Hex())
	}

	expectedBodyRoot := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000003")
	if resp.Data.Header.Message.BodyRoot != expectedBodyRoot {
		t.Errorf("unexpected body_root: %s", resp.Data.Header.Message.BodyRoot.Hex())
	}
}

func TestGetBlockHeader_ByBlockRoot(t *testing.T) {
	blockRoot := "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/eth/v1/beacon/headers/"+blockRoot {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"execution_optimistic": false,
			"finalized": false,
			"data": {
				"root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
				"canonical": true,
				"header": {
					"message": {
						"slot": "99999",
						"proposer_index": "500",
						"parent_root": "0x0000000000000000000000000000000000000000000000000000000000000000",
						"state_root": "0x0000000000000000000000000000000000000000000000000000000000000000",
						"body_root": "0x0000000000000000000000000000000000000000000000000000000000000000"
					},
					"signature": "0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
				}
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	resp, err := client.GetBlockHeader(context.Background(), blockRoot)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Data.Header.Message.Slot != 99999 {
		t.Errorf("expected slot 99999, got %d", resp.Data.Header.Message.Slot)
	}
}

func TestConsensusVersions(t *testing.T) {
	tests := []struct {
		version ConsensusVersion
		want    string
	}{
		{ConsensusVersionPhase0, "phase0"},
		{ConsensusVersionAltair, "altair"},
		{ConsensusVersionBellatrix, "bellatrix"},
		{ConsensusVersionCapella, "capella"},
		{ConsensusVersionDeneb, "deneb"},
		{ConsensusVersionElectra, "electra"},
		{ConsensusVersionFulu, "fulu"},
	}

	for _, tt := range tests {
		if string(tt.version) != tt.want {
			t.Errorf("ConsensusVersion %s != %s", tt.version, tt.want)
		}
	}
}
