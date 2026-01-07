package beaconclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/protolambda/zrnt/eth2/beacon/altair"
	"github.com/protolambda/zrnt/eth2/beacon/bellatrix"
	"github.com/protolambda/zrnt/eth2/beacon/capella"
	"github.com/protolambda/zrnt/eth2/beacon/deneb"
	"github.com/protolambda/zrnt/eth2/beacon/electra"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
)

// ConsensusVersion represents the consensus version/fork name
type ConsensusVersion string

const (
	ConsensusVersionPhase0    ConsensusVersion = "phase0"
	ConsensusVersionAltair    ConsensusVersion = "altair"
	ConsensusVersionBellatrix ConsensusVersion = "bellatrix"
	ConsensusVersionCapella   ConsensusVersion = "capella"
	ConsensusVersionDeneb     ConsensusVersion = "deneb"
	ConsensusVersionElectra   ConsensusVersion = "electra"
	ConsensusVersionFulu      ConsensusVersion = "fulu"
)

// BeaconBlockHeader represents the header of a beacon block
type BeaconBlockHeader struct {
	Slot          uint64      `json:"slot,string"`
	ProposerIndex uint64      `json:"proposer_index,string"`
	ParentRoot    common.Hash `json:"parent_root"`
	StateRoot     common.Hash `json:"state_root"`
	BodyRoot      common.Hash `json:"body_root"`
}

// SignedBeaconBlockHeader represents a signed beacon block header
type SignedBeaconBlockHeader struct {
	Message   BeaconBlockHeader `json:"message"`
	Signature string            `json:"signature"`
}

// Eth1Data represents the Eth1 data in a beacon block
type Eth1Data struct {
	DepositRoot  common.Hash `json:"deposit_root"`
	DepositCount uint64      `json:"deposit_count,string"`
	BlockHash    common.Hash `json:"block_hash"`
}

// Checkpoint represents a checkpoint in the beacon chain
type Checkpoint struct {
	Epoch uint64      `json:"epoch,string"`
	Root  common.Hash `json:"root"`
}

// AttestationData represents attestation data
type AttestationData struct {
	Slot            uint64      `json:"slot,string"`
	Index           uint64      `json:"index,string"`
	BeaconBlockRoot common.Hash `json:"beacon_block_root"`
	Source          Checkpoint  `json:"source"`
	Target          Checkpoint  `json:"target"`
}

// SignedBeaconBlock represents a signed beacon block
type SignedBeaconBlock struct {
	Message   json.RawMessage `json:"message"`
	Signature string      `json:"signature"`
}

// BlockResponse represents the response from /eth/v2/beacon/blocks/{block_id}
type BlockResponse struct {
	// Version is the consensus version (phase0, altair, bellatrix, capella, deneb, electra, fulu)
	Version ConsensusVersion `json:"version"`
	// ExecutionOptimistic is true if the response references an unverified execution payload
	ExecutionOptimistic bool `json:"execution_optimistic"`
	// Finalized is true if the response references the finalized history of the chain
	Finalized bool `json:"finalized"`
	// Data contains the signed beacon block (structure varies by version)
	Data SignedBeaconBlock `json:"data"`
}

// ParseBlock parses the block data into the appropriate beacon block structure based on the version
func (block *BlockResponse) ParseBlock() (any, error) {
	switch block.Version {
	case ConsensusVersionPhase0:
		var body phase0.BeaconBlock
		if err := json.Unmarshal(block.Data.Message, &body); err != nil {
			return nil, err
		}
		return &body, nil
	case ConsensusVersionAltair:
		var body altair.BeaconBlock
		if err := json.Unmarshal(block.Data.Message, &body); err != nil {
			return nil, err
		}
		return body, nil
	case ConsensusVersionBellatrix:
		var body bellatrix.BeaconBlock
		if err := json.Unmarshal(block.Data.Message, &body); err != nil {
			return nil, err
		}
		return &body, nil
	case ConsensusVersionCapella:
		var body capella.BeaconBlock
		if err := json.Unmarshal(block.Data.Message, &body); err != nil {
			return nil, err
		}
		return &body, nil
	case ConsensusVersionDeneb:
		var body deneb.BeaconBlock
		if err := json.Unmarshal(block.Data.Message, &body); err != nil {
			return nil, err
		}
		return &body, nil
	case ConsensusVersionElectra,ConsensusVersionFulu:
		var body electra.BeaconBlock
		if err := json.Unmarshal(block.Data.Message, &body); err != nil {
			return nil, err
		}
		return &body, nil
	default:
		return nil, fmt.Errorf("unsupported consensus version: %s", block.Version)
	}
}

// GetBlock retrieves block details for a given block id
// Endpoint: GET /eth/v2/beacon/blocks/{block_id}
//
// block_id can be: "head", "genesis", "finalized", <slot>, <hex encoded blockRoot with 0x prefix>
//
// Note: The block body structure varies by consensus version. Use the Version field to determine
// the appropriate structure for parsing the Body field.
func (c *Client) GetBlock(ctx context.Context, blockID string) (*BlockResponse, error) {
	body, err := c.doRequest(ctx, http.MethodGet, "/eth/v2/beacon/blocks/"+blockID, nil)
	if err != nil {
		return nil, err
	}

	var resp BlockResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BlockRootData represents the block root response data
type BlockRootData struct {
	Root common.Hash `json:"root"`
}

// BlockRootResponse represents the response from /eth/v1/beacon/blocks/{block_id}/root
type BlockRootResponse struct {
	ExecutionOptimistic bool          `json:"execution_optimistic"`
	Finalized           bool          `json:"finalized"`
	Data                BlockRootData `json:"data"`
}

// GetBlockRoot retrieves hashTreeRoot of BeaconBlock/BeaconBlockHeader
// Endpoint: GET /eth/v1/beacon/blocks/{block_id}/root
//
// block_id can be: "head", "genesis", "finalized", <slot>, <hex encoded blockRoot with 0x prefix>
func (c *Client) GetBlockRoot(ctx context.Context, blockID string) (*BlockRootResponse, error) {
	body, err := c.doRequest(ctx, http.MethodGet, "/eth/v1/beacon/blocks/"+blockID+"/root", nil)
	if err != nil {
		return nil, err
	}

	var resp BlockRootResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BlockHeaderData represents a beacon block header with additional metadata
type BlockHeaderData struct {
	Root      common.Hash             `json:"root"`
	Canonical bool                    `json:"canonical"`
	Header    SignedBeaconBlockHeader `json:"header"`
}

// BlockHeaderResponse represents the response from /eth/v1/beacon/headers/{block_id}
type BlockHeaderResponse struct {
	ExecutionOptimistic bool            `json:"execution_optimistic"`
	Finalized           bool            `json:"finalized"`
	Data                BlockHeaderData `json:"data"`
}

// GetBlockHeader retrieves block header for a given block id
// Endpoint: GET /eth/v1/beacon/headers/{block_id}
//
// block_id can be: "head", "genesis", "finalized", <slot>, <hex encoded blockRoot with 0x prefix>
func (c *Client) GetBlockHeader(ctx context.Context, blockID string) (*BlockHeaderResponse, error) {
	body, err := c.doRequest(ctx, http.MethodGet, "/eth/v1/beacon/headers/"+blockID, nil)
	if err != nil {
		return nil, err
	}

	var resp BlockHeaderResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
