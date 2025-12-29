package beacon

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
)

// SpecData contains specification configuration
// Most values are strings (hex or quoted integers), but some like BLOB_SCHEDULE are arrays
type SpecData struct {
	// Fork epochs
	AltairForkEpoch    uint64 `json:"ALTAIR_FORK_EPOCH,string"`
	BellatrixForkEpoch uint64 `json:"BELLATRIX_FORK_EPOCH,string"`
	CapellaForkEpoch   uint64 `json:"CAPELLA_FORK_EPOCH,string"`
	DenebForkEpoch     uint64 `json:"DENEB_FORK_EPOCH,string"`
	ElectraForkEpoch   uint64 `json:"ELECTRA_FORK_EPOCH,string"`
	FuluForkEpoch      uint64 `json:"FULU_FORK_EPOCH,string"`
	GloasForkEpoch     uint64 `json:"GLOAS_FORK_EPOCH,string"`

	// Fork versions
	GenesisForkVersion   string `json:"GENESIS_FORK_VERSION"`
	AltairForkVersion    string `json:"ALTAIR_FORK_VERSION"`
	BellatrixForkVersion string `json:"BELLATRIX_FORK_VERSION"`
	CapellaForkVersion   string `json:"CAPELLA_FORK_VERSION"`
	DenebForkVersion     string `json:"DENEB_FORK_VERSION"`
	ElectraForkVersion   string `json:"ELECTRA_FORK_VERSION"`
	FuluForkVersion      string `json:"FULU_FORK_VERSION"`
	GloasForkVersion     string `json:"GLOAS_FORK_VERSION"`

	// Config
	ConfigName             string         `json:"CONFIG_NAME"`
	PresetBase             string         `json:"PRESET_BASE"`
	DepositChainID         uint64         `json:"DEPOSIT_CHAIN_ID,string"`
	DepositNetworkID       uint64         `json:"DEPOSIT_NETWORK_ID,string"`
	DepositContractAddress common.Address `json:"DEPOSIT_CONTRACT_ADDRESS"`

	// Timing
	SecondsPerSlot                   uint64 `json:"SECONDS_PER_SLOT,string"`
	SlotsPerEpoch                    uint64 `json:"SLOTS_PER_EPOCH,string"`
	GenesisDelay                     uint64 `json:"GENESIS_DELAY,string"`
	MinGenesisTime                   uint64 `json:"MIN_GENESIS_TIME,string"`
	SecondsPerEth1Block              uint64 `json:"SECONDS_PER_ETH1_BLOCK,string"`
	MinValidatorWithdrawabilityDelay uint64 `json:"MIN_VALIDATOR_WITHDRAWABILITY_DELAY,string"`
	ShardCommitteePeriod             uint64 `json:"SHARD_COMMITTEE_PERIOD,string"`
	Eth1FollowDistance               uint64 `json:"ETH1_FOLLOW_DISTANCE,string"`
	MinGenesisActiveValidatorCount   uint64 `json:"MIN_GENESIS_ACTIVE_VALIDATOR_COUNT,string"`

	// Blob related
	FieldElementsPerBlob             uint64              `json:"FIELD_ELEMENTS_PER_BLOB,string"`
	MaxBlobsPerBlock                 uint64              `json:"MAX_BLOBS_PER_BLOCK,string"`
	MaxBlobsPerBlockElectra          uint64              `json:"MAX_BLOBS_PER_BLOCK_ELECTRA,string"`
	MaxBlobCommitmentsPerBlock       uint64              `json:"MAX_BLOB_COMMITMENTS_PER_BLOCK,string"`
	BlobSidecarSubnetCount           uint64              `json:"BLOB_SIDECAR_SUBNET_COUNT,string"`
	BlobSidecarSubnetCountElectra    uint64              `json:"BLOB_SIDECAR_SUBNET_COUNT_ELECTRA,string"`
	MinEpochsForBlobSidecarsRequests uint64              `json:"MIN_EPOCHS_FOR_BLOB_SIDECARS_REQUESTS,string"`
	MaxRequestBlobSidecars           uint64              `json:"MAX_REQUEST_BLOB_SIDECARS,string"`
	MaxRequestBlobSidecarsElectra    uint64              `json:"MAX_REQUEST_BLOB_SIDECARS_ELECTRA,string"`
	BlobSchedule                     []BlobScheduleEntry `json:"BLOB_SCHEDULE"`

	// Validator
	MaxEffectiveBalance        uint64 `json:"MAX_EFFECTIVE_BALANCE,string"`
	MaxEffectiveBalanceElectra uint64 `json:"MAX_EFFECTIVE_BALANCE_ELECTRA,string"`
	MinActivationBalance       uint64 `json:"MIN_ACTIVATION_BALANCE,string"`
	EjectionBalance            uint64 `json:"EJECTION_BALANCE,string"`
	EffectiveBalanceIncrement  uint64 `json:"EFFECTIVE_BALANCE_INCREMENT,string"`
	MinDepositAmount           uint64 `json:"MIN_DEPOSIT_AMOUNT,string"`

	// Committee
	MaxCommitteesPerSlot      uint64 `json:"MAX_COMMITTEES_PER_SLOT,string"`
	TargetCommitteeSize       uint64 `json:"TARGET_COMMITTEE_SIZE,string"`
	MaxValidatorsPerCommittee uint64 `json:"MAX_VALIDATORS_PER_COMMITTEE,string"`
	ShuffleRoundCount         uint64 `json:"SHUFFLE_ROUND_COUNT,string"`

	// Rewards and penalties
	BaseRewardFactor            uint64 `json:"BASE_REWARD_FACTOR,string"`
	WhistleblowerRewardQuotient uint64 `json:"WHISTLEBLOWER_REWARD_QUOTIENT,string"`
	ProposerRewardQuotient      uint64 `json:"PROPOSER_REWARD_QUOTIENT,string"`
	InactivityPenaltyQuotient   uint64 `json:"INACTIVITY_PENALTY_QUOTIENT,string"`
	MinSlashingPenaltyQuotient  uint64 `json:"MIN_SLASHING_PENALTY_QUOTIENT,string"`

	// Max operations per block
	MaxProposerSlashings     uint64 `json:"MAX_PROPOSER_SLASHINGS,string"`
	MaxAttesterSlashings     uint64 `json:"MAX_ATTESTER_SLASHINGS,string"`
	MaxAttestations          uint64 `json:"MAX_ATTESTATIONS,string"`
	MaxDeposits              uint64 `json:"MAX_DEPOSITS,string"`
	MaxVoluntaryExits        uint64 `json:"MAX_VOLUNTARY_EXITS,string"`
	MaxBlsToExecutionChanges uint64 `json:"MAX_BLS_TO_EXECUTION_CHANGES,string"`
	MaxWithdrawalsPerPayload uint64 `json:"MAX_WITHDRAWALS_PER_PAYLOAD,string"`

	// Sync committee
	SyncCommitteeSize            uint64 `json:"SYNC_COMMITTEE_SIZE,string"`
	EpochsPerSyncCommitteePeriod uint64 `json:"EPOCHS_PER_SYNC_COMMITTEE_PERIOD,string"`
	MinSyncCommitteeParticipants uint64 `json:"MIN_SYNC_COMMITTEE_PARTICIPANTS,string"`
	SyncCommitteeSubnetCount     uint64 `json:"SYNC_COMMITTEE_SUBNET_COUNT,string"`

	// Networking
	SubnetsPerNode                  uint64 `json:"SUBNETS_PER_NODE,string"`
	AttestationPropagationSlotRange uint64 `json:"ATTESTATION_PROPAGATION_SLOT_RANGE,string"`
	MaxRequestBlocksDeneb           uint64 `json:"MAX_REQUEST_BLOCKS_DENEB,string"`
}

// BlobScheduleEntry represents an entry in BLOB_SCHEDULE
type BlobScheduleEntry struct {
	Epoch            uint64 `json:"EPOCH,string"`
	MaxBlobsPerBlock uint64 `json:"MAX_BLOBS_PER_BLOCK,string"`
}

// specResponse represents the full response from /eth/v1/config/spec
type specResponse struct {
	Data SpecData `json:"data"`
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
func (c *Client) GetSpec(ctx context.Context) (*SpecData, error) {
	body, err := c.doRequest(ctx, http.MethodGet, "/eth/v1/config/spec", nil)
	if err != nil {
		return nil, err
	}

	var resp specResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
