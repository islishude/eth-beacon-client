package beaconclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
)

// BlobsData represents the blobs data
type BlobsData struct {
	ExecutionOptimistic bool           `json:"execution_optimistic"`
	Finalized           bool           `json:"finalized"`
	Data                []kzg4844.Blob `json:"data"`
}

// GetBlobs retrieves blobs for a given block id
// Endpoint: GET /eth/v1/beacon/blobs/{block_id}
//
// block_id can be: "head", "genesis", "finalized", <slot>, <hex encoded blockRoot with 0x prefix>
// versionedHashes is optional - if provided, only blobs for specified versioned hashes are returned
func (c *Client) GetBlobs(ctx context.Context, blockID string, versionedHashes ...common.Hash) (*BlobsData, error) {
	var query url.Values
	if len(versionedHashes) > 0 {
		query = url.Values{}
		for _, hash := range versionedHashes {
			query.Add("versioned_hashes", hash.Hex())
		}
	}

	body, err := c.doRequest(ctx, http.MethodGet, "/eth/v1/beacon/blobs/"+blockID, query)
	if err != nil {
		return nil, err
	}

	var resp BlobsData
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
