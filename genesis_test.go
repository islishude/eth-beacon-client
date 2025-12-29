package beacon

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetGenesis_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/eth/v1/beacon/genesis" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": {
				"genesis_time": "1606824023",
				"genesis_validators_root": "0x4b363db94e286120d76eb905340fdd4e54bfe9f06bf33ff6cf5ad27f511bfe95",
				"genesis_fork_version": "0"
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	genesis, err := client.GetGenesis(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if genesis.GenesisTime != 1606824023 {
		t.Errorf("expected genesis_time 1606824023, got %d", genesis.GenesisTime)
	}
	if genesis.GenesisValidatorsRoot != "0x4b363db94e286120d76eb905340fdd4e54bfe9f06bf33ff6cf5ad27f511bfe95" {
		t.Errorf("unexpected genesis_validators_root: %s", genesis.GenesisValidatorsRoot)
	}
}

func TestGetGenesis_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"code": 500, "message": "Internal server error"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetGenesis(context.Background())
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
