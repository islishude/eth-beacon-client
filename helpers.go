package beacon

// GENESIS_SLOT is the first slot of the beacon chain
const GENESIS_SLOT uint64 = 0

// ComputeTimestampAtSlot computes the timestamp at a given slot
// This is equivalent to the Python spec function:
//
//	def compute_timestamp_at_slot(state: BeaconState, slot: Slot) -> uint64:
//	    slots_since_genesis = slot - GENESIS_SLOT
//	    return uint64(state.genesis_time + slots_since_genesis * SECONDS_PER_SLOT)
//
// Parameters:
//   - genesisTime: the genesis time from GenesisData
//   - slot: the target slot number
//   - secondsPerSlot: seconds per slot from SpecData (typically 12)
//
// Returns the Unix timestamp at the start of the given slot
func ComputeTimestampAtSlot(genesisTime, slot, secondsPerSlot uint64) uint64 {
	slotsSinceGenesis := slot - GENESIS_SLOT
	return genesisTime + slotsSinceGenesis*secondsPerSlot
}

// ComputeSlotAtTimestamp computes the slot at a given timestamp
// This is the inverse of ComputeTimestampAtSlot
//
// Parameters:
//   - genesisTime: the genesis time from GenesisData
//   - timestamp: the Unix timestamp
//   - secondsPerSlot: seconds per slot from SpecData (typically 12)
//
// Returns the slot number at the given timestamp
func ComputeSlotAtTimestamp(genesisTime, timestamp, secondsPerSlot uint64) uint64 {
	if timestamp < genesisTime {
		return GENESIS_SLOT
	}
	return (timestamp - genesisTime) / secondsPerSlot
}
