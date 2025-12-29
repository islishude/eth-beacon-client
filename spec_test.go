package beacon

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSpec_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/eth/v1/config/spec" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": {
				"CONFIG_NAME": "mainnet",
				"PRESET_BASE": "mainnet",
				"DEPOSIT_CHAIN_ID": "1",
				"DEPOSIT_NETWORK_ID": "1",
				"DEPOSIT_CONTRACT_ADDRESS": "0x00000000219ab540356cBB839Cbe05303d7705Fa",
				"SECONDS_PER_SLOT": "12",
				"SLOTS_PER_EPOCH": "32",
				"GENESIS_DELAY": "604800",
				"MIN_GENESIS_TIME": "1606824000",
				"SECONDS_PER_ETH1_BLOCK": "14",
				"DENEB_FORK_EPOCH": "269568",
				"ELECTRA_FORK_EPOCH": "364032",
				"MAX_BLOBS_PER_BLOCK": "6",
				"BLOB_SCHEDULE": [
					{"EPOCH": "269568", "MAX_BLOBS_PER_BLOCK": "6"},
					{"EPOCH": "364032", "MAX_BLOBS_PER_BLOCK": "9"}
				]
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	spec, err := client.GetSpec(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if spec.ConfigName != "mainnet" {
		t.Errorf("expected CONFIG_NAME mainnet, got %s", spec.ConfigName)
	}
	if spec.PresetBase != "mainnet" {
		t.Errorf("expected PRESET_BASE mainnet, got %s", spec.PresetBase)
	}
	if spec.DepositChainID != 1 {
		t.Errorf("expected DEPOSIT_CHAIN_ID 1, got %d", spec.DepositChainID)
	}
	if spec.SecondsPerSlot != 12 {
		t.Errorf("expected SECONDS_PER_SLOT 12, got %d", spec.SecondsPerSlot)
	}
	if spec.SlotsPerEpoch != 32 {
		t.Errorf("expected SLOTS_PER_EPOCH 32, got %d", spec.SlotsPerEpoch)
	}
	if spec.DenebForkEpoch != 269568 {
		t.Errorf("expected DENEB_FORK_EPOCH 269568, got %d", spec.DenebForkEpoch)
	}
	if len(spec.BlobSchedule) != 2 {
		t.Errorf("expected 2 BLOB_SCHEDULE entries, got %d", len(spec.BlobSchedule))
	}
	if spec.BlobSchedule[0].Epoch != 269568 {
		t.Errorf("expected first BLOB_SCHEDULE epoch 269568, got %d", spec.BlobSchedule[0].Epoch)
	}
	if spec.BlobSchedule[0].MaxBlobsPerBlock != 6 {
		t.Errorf("expected first BLOB_SCHEDULE max_blobs 6, got %d", spec.BlobSchedule[0].MaxBlobsPerBlock)
	}
}

func TestGetSpec_DepositContractAddress(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": {
				"DEPOSIT_CONTRACT_ADDRESS": "0x00000000219ab540356cBB839Cbe05303d7705Fa"
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	spec, err := client.GetSpec(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedAddr := "0x00000000219ab540356cBB839Cbe05303d7705Fa"
	if spec.DepositContractAddress.Hex() != expectedAddr {
		t.Errorf("expected DEPOSIT_CONTRACT_ADDRESS %s, got %s", expectedAddr, spec.DepositContractAddress.Hex())
	}
}

func TestGetSpec_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"code": 500, "message": "Internal server error"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetSpec(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Code != 500 {
		t.Errorf("expected code 500, got %d", apiErr.Code)
	}
}

func TestGetSpec_EmptyBlobSchedule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": {
				"CONFIG_NAME": "testnet",
				"BLOB_SCHEDULE": []
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	spec, err := client.GetSpec(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(spec.BlobSchedule) != 0 {
		t.Errorf("expected empty BLOB_SCHEDULE, got %d entries", len(spec.BlobSchedule))
	}
}
