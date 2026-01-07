package beaconclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetNodeIdentity_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/eth/v1/node/identity" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": {
				"peer_id": "QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N",
				"enr": "enr:-IS4QHCYrYZbAKWCBRlAy5zzaDZXJBGkcnh4MHcBFZntXNFrdvJjX04jRzjzCBOonrkTfj499SZuOh8R33Ls8RRcy5wBgmlkgnY0gmlwhH8AAAGJc2VjcDI1NmsxoQPKY0yuDUmstAHYpMa2_oxVtw0RW_QAdpzBQA8yWM0xOIN1ZHCCdl8",
				"p2p_addresses": ["/ip4/7.7.7.7/tcp/4242/p2p/QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N"],
				"discovery_addresses": ["/ip4/7.7.7.7/udp/30303/p2p/QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N"],
				"metadata": {
					"seq_number": "1",
					"attnets": "0x0000000000000000",
					"syncnets": "0x0f"
				}
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	identity, err := client.GetNodeIdentity(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if identity.PeerID != "QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N" {
		t.Errorf("unexpected peer_id: %s", identity.PeerID)
	}
	if len(identity.P2PAddresses) != 1 {
		t.Errorf("expected 1 p2p_address, got %d", len(identity.P2PAddresses))
	}
	if identity.Metadata.SeqNumber != "1" {
		t.Errorf("unexpected seq_number: %s", identity.Metadata.SeqNumber)
	}
	if identity.Metadata.Syncnets != "0x0f" {
		t.Errorf("unexpected syncnets: %s", identity.Metadata.Syncnets)
	}
}

func TestGetPeers_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/eth/v1/node/peers" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": [
				{
					"peer_id": "QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N",
					"enr": "enr:-test",
					"last_seen_p2p_address": "/ip4/7.7.7.7/tcp/4242",
					"state": "connected",
					"direction": "inbound"
				}
			],
			"meta": {
				"count": 1
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	resp, err := client.GetPeers(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 peer, got %d", len(resp.Data))
	}
	if resp.Data[0].State != PeerStateConnected {
		t.Errorf("unexpected state: %s", resp.Data[0].State)
	}
	if resp.Data[0].Direction != PeerDirectionInbound {
		t.Errorf("unexpected direction: %s", resp.Data[0].Direction)
	}
	if resp.Meta.Count != 1 {
		t.Errorf("unexpected count: %d", resp.Meta.Count)
	}
}

func TestGetPeers_WithFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		states := r.URL.Query()["state"]
		if len(states) != 1 || states[0] != "connected" {
			t.Errorf("unexpected state filter: %v", states)
		}
		directions := r.URL.Query()["direction"]
		if len(directions) != 1 || directions[0] != "inbound" {
			t.Errorf("unexpected direction filter: %v", directions)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data": [], "meta": {"count": 0}}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetPeers(context.Background(), &GetPeersOption{
		State:     []PeerState{PeerStateConnected},
		Direction: []PeerDirection{PeerDirectionInbound},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetPeer_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/eth/v1/node/peers/QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": {
				"peer_id": "QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N",
				"enr": null,
				"last_seen_p2p_address": "/ip4/7.7.7.7/tcp/4242",
				"state": "connected",
				"direction": "outbound"
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	peer, err := client.GetPeer(context.Background(), "QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if peer.PeerID != "QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N" {
		t.Errorf("unexpected peer_id: %s", peer.PeerID)
	}
	if peer.ENR != nil {
		t.Errorf("expected nil enr, got %v", peer.ENR)
	}
	if peer.Direction != PeerDirectionOutbound {
		t.Errorf("unexpected direction: %s", peer.Direction)
	}
}

func TestGetPeer_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code": 404, "message": "Peer not found"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetPeer(context.Background(), "invalid-peer")
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

func TestGetPeerCount_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/eth/v1/node/peer_count" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": {
				"disconnected": "10",
				"connecting": "2",
				"connected": "50",
				"disconnecting": "1"
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	count, err := client.GetPeerCount(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if count.Disconnected != 10 {
		t.Errorf("expected disconnected 10, got %d", count.Disconnected)
	}
	if count.Connecting != 2 {
		t.Errorf("expected connecting 2, got %d", count.Connecting)
	}
	if count.Connected != 50 {
		t.Errorf("expected connected 50, got %d", count.Connected)
	}
	if count.Disconnecting != 1 {
		t.Errorf("expected disconnecting 1, got %d", count.Disconnecting)
	}
}

func TestGetNodeVersion_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/eth/v1/node/version" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": {
				"version": "Lighthouse/v0.1.5 (Linux x86_64)"
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	version, err := client.GetNodeVersion(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if version.Version != "Lighthouse/v0.1.5 (Linux x86_64)" {
		t.Errorf("unexpected version: %s", version.Version)
	}
}

func TestGetSyncingStatus_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/eth/v1/node/syncing" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": {
				"head_slot": "12345",
				"sync_distance": "100",
				"is_syncing": true,
				"is_optimistic": false,
				"el_offline": false
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	status, err := client.GetSyncingStatus(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if status.HeadSlot != 12345 {
		t.Errorf("expected head_slot 12345, got %d", status.HeadSlot)
	}
	if status.SyncDistance != 100 {
		t.Errorf("expected sync_distance 100, got %d", status.SyncDistance)
	}
	if !status.IsSyncing {
		t.Error("expected is_syncing true")
	}
	if status.IsOptimistic {
		t.Error("expected is_optimistic false")
	}
	if status.ELOffline {
		t.Error("expected el_offline false")
	}
}

func TestGetSyncingStatus_Synced(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": {
				"head_slot": "12345",
				"sync_distance": "0",
				"is_syncing": false,
				"is_optimistic": false,
				"el_offline": false
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	status, err := client.GetSyncingStatus(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if status.SyncDistance != 0 {
		t.Errorf("expected sync_distance 0, got %d", status.SyncDistance)
	}
	if status.IsSyncing {
		t.Error("expected is_syncing false")
	}
}

func TestGetHealth_Ready(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/eth/v1/node/health" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	status, err := client.GetHealth(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if status != HealthStatusReady {
		t.Errorf("expected status 200, got %d", status)
	}
}

func TestGetHealth_Syncing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusPartialContent) // 206
	}))
	defer server.Close()

	client := NewClient(server.URL)
	status, err := client.GetHealth(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if status != HealthStatusSyncing {
		t.Errorf("expected status 206, got %d", status)
	}
}

func TestGetHealth_NotInitialized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable) // 503
	}))
	defer server.Close()

	client := NewClient(server.URL)
	status, err := client.GetHealth(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if status != HealthStatusNotInitialized {
		t.Errorf("expected status 503, got %d", status)
	}
}
