package beacon

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestGetDepositContract_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/eth/v1/config/deposit_contract" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": {
				"chain_id": "1",
				"address": "0x00000000219ab540356cBB839Cbe05303d7705Fa"
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	contract, err := client.GetDepositContract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if contract.ChainID != 1 {
		t.Errorf("expected chain_id 1, got %d", contract.ChainID)
	}

	expectedAddr := common.HexToAddress("0x00000000219ab540356cBB839Cbe05303d7705Fa")
	if contract.Address != expectedAddr {
		t.Errorf("expected address %s, got %s", expectedAddr.Hex(), contract.Address.Hex())
	}
}

func TestGetDepositContract_Goerli(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": {
				"chain_id": "5",
				"address": "0xff50ed3d0ec03aC01D4C79aAd74928BFF48a7b2b"
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	contract, err := client.GetDepositContract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if contract.ChainID != 5 {
		t.Errorf("expected chain_id 5, got %d", contract.ChainID)
	}

	expectedAddr := common.HexToAddress("0xff50ed3d0ec03aC01D4C79aAd74928BFF48a7b2b")
	if contract.Address != expectedAddr {
		t.Errorf("expected address %s, got %s", expectedAddr.Hex(), contract.Address.Hex())
	}
}

func TestGetDepositContract_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{
			"code": 500,
			"message": "Internal server error"
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetDepositContract(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.Code != 500 {
		t.Errorf("expected error code 500, got %d", apiErr.Code)
	}
}
