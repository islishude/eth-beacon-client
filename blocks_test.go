package beaconclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/protolambda/zrnt/eth2/beacon/altair"
	"github.com/protolambda/zrnt/eth2/beacon/bellatrix"
	"github.com/protolambda/zrnt/eth2/beacon/capella"
	"github.com/protolambda/zrnt/eth2/beacon/deneb"
	"github.com/protolambda/zrnt/eth2/beacon/electra"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
)

func TestGetBlock_ByVersion(t *testing.T) {
	tests := []struct {
		name         string
		blockID      string
		testdataFile string
		wantVersion  ConsensusVersion
		wantSlot     uint64
		wantProposer uint64
		extraCheck   func(t *testing.T, block any)
	}{
		{
			name:         "fulu",
			blockID:      "head",
			testdataFile: "testdata/fulu.block.json",
			wantVersion:  ConsensusVersionFulu,
			wantSlot:     13410020,
			wantProposer: 1797581,
			extraCheck: func(t *testing.T, block any) {
				if _, ok := block.(*electra.BeaconBlock); !ok {
					t.Errorf("expected *electra.BeaconBlock, got %T", block)
				}
			},
		},
		{
			name:         "electra",
			blockID:      "11982020",
			testdataFile: "testdata/electra.block.json",
			wantVersion:  ConsensusVersionElectra,
			wantSlot:     11982020,
			wantProposer: 1605697,
			extraCheck: func(t *testing.T, block any) {
				if _, ok := block.(*electra.BeaconBlock); !ok {
					t.Errorf("expected *electra.BeaconBlock, got %T", block)
				}
			},
		},
		{
			name:         "deneb",
			blockID:      "11511320",
			testdataFile: "testdata/deneb.block.json",
			wantVersion:  ConsensusVersionDeneb,
			wantSlot:     11511320,
			wantProposer: 1667419,
			extraCheck: func(t *testing.T, block any) {
				denebBlock, ok := block.(*deneb.BeaconBlock)
				if !ok {
					t.Errorf("expected *deneb.BeaconBlock, got %T", block)
					return
				}
				if len(denebBlock.Body.BlobKZGCommitments) != 5 {
					t.Errorf("BlobKZGCommitments length is not 5")
				}
			},
		},
		{
			name:         "capella",
			blockID:      "11511320",
			testdataFile: "testdata/cepella.block.json",
			wantVersion:  ConsensusVersionCapella,
			wantSlot:     0, // skip slot check
			wantProposer: 0, // skip proposer check
			extraCheck: func(t *testing.T, block any) {
				if _, ok := block.(*capella.BeaconBlock); !ok {
					t.Errorf("expected *capella.BeaconBlock, got %T", block)
				}
			},
		},
		{
			name:         "bellatrix",
			blockID:      "6155220",
			testdataFile: "testdata/bellatrix.block.json",
			wantVersion:  ConsensusVersionBellatrix,
			wantSlot:     6155220,
			wantProposer: 218470,
			extraCheck: func(t *testing.T, block any) {
				if _, ok := block.(*bellatrix.BeaconBlock); !ok {
					t.Errorf("expected *bellatrix.BeaconBlock, got %T", block)
				}
			},
		},
		{
			name:         "altair",
			blockID:      "3199220",
			testdataFile: "testdata/altair.block.json",
			wantVersion:  ConsensusVersionAltair,
			wantSlot:     3199220,
			wantProposer: 66269,
			extraCheck: func(t *testing.T, block any) {
				if _, ok := block.(altair.BeaconBlock); !ok {
					t.Errorf("expected altair.BeaconBlock, got %T", block)
				}
			},
		},
		{
			name:         "phase0",
			blockID:      "1511320",
			testdataFile: "testdata/phase0.block.json",
			wantVersion:  ConsensusVersionPhase0,
			wantSlot:     1511320,
			wantProposer: 145572,
			extraCheck: func(t *testing.T, block any) {
				if _, ok := block.(*phase0.BeaconBlock); !ok {
					t.Errorf("expected *phase0.BeaconBlock, got %T", block)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := "/eth/v2/beacon/blocks/" + tt.blockID
				if r.URL.Path != expectedPath {
					t.Errorf("unexpected path: got %s, want %s", r.URL.Path, expectedPath)
				}
				if r.Method != http.MethodGet {
					t.Errorf("unexpected method: %s", r.Method)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				data, err := os.ReadFile(tt.testdataFile)
				if err != nil {
					t.Fatalf("failed to read test data: %v", err)
				}
				_, _ = w.Write(data)
			}))
			defer server.Close()

			client := NewClient(server.URL)
			resp, err := client.GetBlock(context.Background(), tt.blockID)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp.Version != tt.wantVersion {
				t.Errorf("version: got %s, want %s", resp.Version, tt.wantVersion)
			}
			if resp.ExecutionOptimistic {
				t.Error("expected execution_optimistic false")
			}
			if !resp.Finalized {
				t.Error("expected finalized true")
			}

			block, err := resp.ParseBlock()
			if err != nil {
				t.Fatalf("unexpected error parsing block: %v", err)
			}

			// Check slot and proposer if specified
			if tt.wantSlot != 0 || tt.wantProposer != 0 {
				switch b := block.(type) {
				case *phase0.BeaconBlock:
					if tt.wantSlot != 0 && uint64(b.Slot) != tt.wantSlot {
						t.Errorf("slot: got %d, want %d", b.Slot, tt.wantSlot)
					}
					if tt.wantProposer != 0 && uint64(b.ProposerIndex) != tt.wantProposer {
						t.Errorf("proposer_index: got %d, want %d", b.ProposerIndex, tt.wantProposer)
					}
				case altair.BeaconBlock:
					if tt.wantSlot != 0 && uint64(b.Slot) != tt.wantSlot {
						t.Errorf("slot: got %d, want %d", b.Slot, tt.wantSlot)
					}
					if tt.wantProposer != 0 && uint64(b.ProposerIndex) != tt.wantProposer {
						t.Errorf("proposer_index: got %d, want %d", b.ProposerIndex, tt.wantProposer)
					}
				case *bellatrix.BeaconBlock:
					if tt.wantSlot != 0 && uint64(b.Slot) != tt.wantSlot {
						t.Errorf("slot: got %d, want %d", b.Slot, tt.wantSlot)
					}
					if tt.wantProposer != 0 && uint64(b.ProposerIndex) != tt.wantProposer {
						t.Errorf("proposer_index: got %d, want %d", b.ProposerIndex, tt.wantProposer)
					}
				case *capella.BeaconBlock:
					if tt.wantSlot != 0 && uint64(b.Slot) != tt.wantSlot {
						t.Errorf("slot: got %d, want %d", b.Slot, tt.wantSlot)
					}
					if tt.wantProposer != 0 && uint64(b.ProposerIndex) != tt.wantProposer {
						t.Errorf("proposer_index: got %d, want %d", b.ProposerIndex, tt.wantProposer)
					}
				case *deneb.BeaconBlock:
					if tt.wantSlot != 0 && uint64(b.Slot) != tt.wantSlot {
						t.Errorf("slot: got %d, want %d", b.Slot, tt.wantSlot)
					}
					if tt.wantProposer != 0 && uint64(b.ProposerIndex) != tt.wantProposer {
						t.Errorf("proposer_index: got %d, want %d", b.ProposerIndex, tt.wantProposer)
					}
				case *electra.BeaconBlock:
					if tt.wantSlot != 0 && uint64(b.Slot) != tt.wantSlot {
						t.Errorf("slot: got %d, want %d", b.Slot, tt.wantSlot)
					}
					if tt.wantProposer != 0 && uint64(b.ProposerIndex) != tt.wantProposer {
						t.Errorf("proposer_index: got %d, want %d", b.ProposerIndex, tt.wantProposer)
					}
				}
			}

			if tt.extraCheck != nil {
				tt.extraCheck(t, block)
			}
		})
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

func TestParseBlock_Error(t *testing.T) {
	tests := []struct {
		name      string
		block     *BlockResponse
		wantError string
	}{
		{
			name:      "nil block response",
			block:     nil,
			wantError: "block response is nil",
		},
		{
			name: "unsupported version",
			block: &BlockResponse{
				Version: "unknown_version",
				Data:    SignedBeaconBlock{Message: []byte(`{}`)},
			},
			wantError: "unsupported consensus version: unknown_version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.block.ParseBlock()
			if tt.wantError == "" {
				// Skip error check for cases where we just want to ensure it doesn't panic
				return
			}
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantError) {
				t.Errorf("error %q does not contain %q", err.Error(), tt.wantError)
			}
		})
	}
}
