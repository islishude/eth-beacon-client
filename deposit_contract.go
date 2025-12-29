package beacon

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
)

// DepositContractData contains deposit contract information
type DepositContractData struct {
	// ChainID is the Id of Eth1 chain on which contract is deployed
	ChainID uint64 `json:"chain_id,string"`
	// Address is the hex encoded deposit contract address with 0x prefix
	Address common.Address `json:"address"`
}

// depositContractResponse represents the full response from /eth/v1/config/deposit_contract
type depositContractResponse struct {
	Data DepositContractData `json:"data"`
}

// GetDepositContract retrieves deposit contract address and chain ID
// Endpoint: GET /eth/v1/config/deposit_contract
//
// Returns:
//   - chain_id: Id of Eth1 chain on which contract is deployed
//   - address: Hex encoded deposit contract address with 0x prefix
func (c *Client) GetDepositContract(ctx context.Context) (*DepositContractData, error) {
	body, err := c.doRequest(ctx, http.MethodGet, "/eth/v1/config/deposit_contract", nil)
	if err != nil {
		return nil, err
	}

	var resp depositContractResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
