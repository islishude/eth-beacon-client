package beaconclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// PeerState represents the connection state of a peer
type PeerState string

const (
	PeerStateDisconnected  PeerState = "disconnected"
	PeerStateConnecting    PeerState = "connecting"
	PeerStateConnected     PeerState = "connected"
	PeerStateDisconnecting PeerState = "disconnecting"
)

// PeerDirection represents the direction of a peer connection
type PeerDirection string

const (
	PeerDirectionInbound  PeerDirection = "inbound"
	PeerDirectionOutbound PeerDirection = "outbound"
)

// NodeIdentityMetadata contains node metadata
type NodeIdentityMetadata struct {
	// SeqNumber is uint64 starting at 0 used to version the node's metadata
	SeqNumber string `json:"seq_number"`
	// Attnets is bitvector representing the node's persistent attestation subnet subscriptions
	Attnets string `json:"attnets"`
	// Syncnets is bitvector representing the node's sync committee subnet subscriptions (present from Altair)
	Syncnets string `json:"syncnets,omitempty"`
	// CustodyGroupCount is uint64 representing the node's custody group count (present from Fulu)
	CustodyGroupCount string `json:"custody_group_count,omitempty"`
}

// NodeIdentity contains node network identity information
type NodeIdentity struct {
	// PeerID is cryptographic hash of a peer's public key
	PeerID string `json:"peer_id"`
	// ENR is Ethereum node record
	ENR string `json:"enr"`
	// P2PAddresses are node's addresses on which eth2 RPC requests are served
	P2PAddresses []string `json:"p2p_addresses"`
	// DiscoveryAddresses are node's addresses on which is listening for discv5 requests
	DiscoveryAddresses []string `json:"discovery_addresses"`
	// Metadata contains node metadata
	Metadata NodeIdentityMetadata `json:"metadata"`
}

// nodeIdentityResponse represents the response from /eth/v1/node/identity
type nodeIdentityResponse struct {
	Data NodeIdentity `json:"data"`
}

// GetNodeIdentity retrieves data about the node's network presence
// Endpoint: GET /eth/v1/node/identity
func (c *Client) GetNodeIdentity(ctx context.Context) (*NodeIdentity, error) {
	body, err := c.doRequest(ctx, http.MethodGet, "/eth/v1/node/identity", nil)
	if err != nil {
		return nil, err
	}

	var resp nodeIdentityResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Peer contains information about a connected peer
type Peer struct {
	// PeerID is cryptographic hash of a peer's public key
	PeerID string `json:"peer_id"`
	// ENR is Ethereum node record (can be null)
	ENR *string `json:"enr"`
	// LastSeenP2PAddress is multiaddrs used in last peer connection
	LastSeenP2PAddress string `json:"last_seen_p2p_address"`
	// State is the connection state
	State PeerState `json:"state"`
	// Direction is the connection direction
	Direction PeerDirection `json:"direction"`
}

// PeersResponse contains peers data with metadata
type PeersResponse struct {
	Data []Peer `json:"data"`
	Meta struct {
		Count int `json:"count"`
	} `json:"meta"`
}

// GetPeersOption represents options for GetPeers
type GetPeersOption struct {
	State     []PeerState
	Direction []PeerDirection
}

// GetPeers retrieves data about the node's network peers
// Endpoint: GET /eth/v1/node/peers
func (c *Client) GetPeers(ctx context.Context, opts *GetPeersOption) (*PeersResponse, error) {
	var query url.Values
	if opts != nil {
		query = url.Values{}
		for _, s := range opts.State {
			query.Add("state", string(s))
		}
		for _, d := range opts.Direction {
			query.Add("direction", string(d))
		}
	}

	body, err := c.doRequest(ctx, http.MethodGet, "/eth/v1/node/peers", query)
	if err != nil {
		return nil, err
	}

	var resp PeersResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// peerResponse represents the response from /eth/v1/node/peers/{peer_id}
type peerResponse struct {
	Data Peer `json:"data"`
}

// GetPeer retrieves data about a specific peer
// Endpoint: GET /eth/v1/node/peers/{peer_id}
func (c *Client) GetPeer(ctx context.Context, peerID string) (*Peer, error) {
	body, err := c.doRequest(ctx, http.MethodGet, "/eth/v1/node/peers/"+peerID, nil)
	if err != nil {
		return nil, err
	}

	var resp peerResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// PeerCount contains peer counts by connection state
type PeerCount struct {
	// Disconnected is the number of disconnected peers
	Disconnected uint64 `json:"disconnected,string"`
	// Connecting is the number of connecting peers
	Connecting uint64 `json:"connecting,string"`
	// Connected is the number of connected peers
	Connected uint64 `json:"connected,string"`
	// Disconnecting is the number of disconnecting peers
	Disconnecting uint64 `json:"disconnecting,string"`
}

// peerCountResponse represents the response from /eth/v1/node/peer_count
type peerCountResponse struct {
	Data PeerCount `json:"data"`
}

// GetPeerCount retrieves number of known peers
// Endpoint: GET /eth/v1/node/peer_count
func (c *Client) GetPeerCount(ctx context.Context) (*PeerCount, error) {
	body, err := c.doRequest(ctx, http.MethodGet, "/eth/v1/node/peer_count", nil)
	if err != nil {
		return nil, err
	}

	var resp peerCountResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// NodeVersion contains node version information
type NodeVersion struct {
	// Version is a string which uniquely identifies the client implementation and its version
	Version string `json:"version"`
}

// nodeVersionResponse represents the response from /eth/v1/node/version
type nodeVersionResponse struct {
	Data NodeVersion `json:"data"`
}

// GetNodeVersion retrieves version string of the running beacon node
// Endpoint: GET /eth/v1/node/version
func (c *Client) GetNodeVersion(ctx context.Context) (*NodeVersion, error) {
	body, err := c.doRequest(ctx, http.MethodGet, "/eth/v1/node/version", nil)
	if err != nil {
		return nil, err
	}

	var resp nodeVersionResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// SyncingStatus contains node syncing status
type SyncingStatus struct {
	// HeadSlot is the head slot node is trying to reach
	HeadSlot uint64 `json:"head_slot,string"`
	// SyncDistance is how many slots node needs to process to reach head (0 if synced)
	SyncDistance uint64 `json:"sync_distance,string"`
	// IsSyncing is true if the node is syncing
	IsSyncing bool `json:"is_syncing"`
	// IsOptimistic is true if the node is optimistically tracking head
	IsOptimistic bool `json:"is_optimistic"`
	// ELOffline is true if the execution client is offline
	ELOffline bool `json:"el_offline"`
}

// syncingStatusResponse represents the response from /eth/v1/node/syncing
type syncingStatusResponse struct {
	Data SyncingStatus `json:"data"`
}

// GetSyncingStatus retrieves node syncing status
// Endpoint: GET /eth/v1/node/syncing
func (c *Client) GetSyncingStatus(ctx context.Context) (*SyncingStatus, error) {
	body, err := c.doRequest(ctx, http.MethodGet, "/eth/v1/node/syncing", nil)
	if err != nil {
		return nil, err
	}

	var resp syncingStatusResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// HealthStatus represents the health status of the node
type HealthStatus int

const (
	// HealthStatusReady indicates node is ready
	HealthStatusReady HealthStatus = 200
	// HealthStatusSyncing indicates node is syncing or execution node is optimistic/offline
	HealthStatusSyncing HealthStatus = 206
	// HealthStatusNotInitialized indicates node is not initialized or having issues
	HealthStatusNotInitialized HealthStatus = 503
)

// GetHealth checks node health status
// Endpoint: GET /eth/v1/node/health
// Returns the health status code (200 = ready, 206 = syncing, 503 = not initialized)
func (c *Client) GetHealth(ctx context.Context) (HealthStatus, error) {
	fullURL := c.baseURL + "/eth/v1/node/health"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return 0, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	return HealthStatus(resp.StatusCode), nil
}
