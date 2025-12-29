package beacon

import (
	"context"
	"encoding/json"
	"net/http"
)

// genesisResponse represents the full response from /eth/v1/beacon/genesis
type genesisResponse struct {
	Data GenesisData `json:"data"`
}

// GenesisData contains the genesis information
type GenesisData struct {
	GenesisTime           uint64 `json:"genesis_time,string"`
	GenesisValidatorsRoot string `json:"genesis_validators_root"`
	GenesisForkVersion    string `json:"genesis_fork_version"`
}

// GetGenesis retrieves details of the chain's genesis
// Endpoint: GET /eth/v1/beacon/genesis
func (c *Client) GetGenesis(ctx context.Context) (*GenesisData, error) {
	body, err := c.doRequest(ctx, http.MethodGet, "/eth/v1/beacon/genesis", nil)
	if err != nil {
		return nil, err
	}

	var resp genesisResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
