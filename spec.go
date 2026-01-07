package beaconclient

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/ztyp/view"
)

// BlobScheduleEntry represents an entry in BLOB_SCHEDULE
type BlobScheduleEntry struct {
	Epoch            view.Uint64View `json:"EPOCH"`
	MaxBlobsPerBlock view.Uint64View `json:"MAX_BLOBS_PER_BLOCK"`
}

type Spec struct {
	common.Config
	BLOB_SCHEDULE []BlobScheduleEntry `json:"BLOB_SCHEDULE"`
}

// specResponse represents the full response from /eth/v1/config/spec
type specResponse struct {
	Data *Spec `json:"data"`
}

// GetSpec retrieves specification configuration used on this node
// Endpoint: GET /eth/v1/config/spec
//
// The configuration includes:
//   - Constants for all hard forks known by the beacon node
//   - Presets for all hard forks supplied to the beacon node
//   - Configuration for the beacon node
//
// Values are returned with following format:
//   - any value starting with 0x in the spec is returned as a hex string
//   - numeric values are returned as a quoted integer
//   - array values are returned as a JSON array
func (c *Client) GetSpec(ctx context.Context) (*Spec, error) {
	body, err := c.doRequest(ctx, http.MethodGet, "/eth/v1/config/spec", nil)
	if err != nil {
		return nil, err
	}

	var resp specResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}
