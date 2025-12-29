package beacon

import "testing"

func TestComputeTimestampAtSlot(t *testing.T) {
	tests := []struct {
		name           string
		genesisTime    uint64
		slot           uint64
		secondsPerSlot uint64
		want           uint64
	}{
		{
			name:           "genesis slot",
			genesisTime:    1606824023,
			slot:           0,
			secondsPerSlot: 12,
			want:           1606824023,
		},
		{
			name:           "slot 1",
			genesisTime:    1606824023,
			slot:           1,
			secondsPerSlot: 12,
			want:           1606824035,
		},
		{
			name:           "slot 100",
			genesisTime:    1606824023,
			slot:           100,
			secondsPerSlot: 12,
			want:           1606825223,
		},
		{
			name:           "mainnet example - slot 269568 (Deneb fork)",
			genesisTime:    1606824023,
			slot:           269568 * 32, // epoch 269568 * slots_per_epoch
			secondsPerSlot: 12,
			want:           1606824023 + 269568*32*12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComputeTimestampAtSlot(tt.genesisTime, tt.slot, tt.secondsPerSlot)
			if got != tt.want {
				t.Errorf("ComputeTimestampAtSlot() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestComputeSlotAtTimestamp(t *testing.T) {
	tests := []struct {
		name           string
		genesisTime    uint64
		timestamp      uint64
		secondsPerSlot uint64
		want           uint64
	}{
		{
			name:           "at genesis",
			genesisTime:    1606824023,
			timestamp:      1606824023,
			secondsPerSlot: 12,
			want:           0,
		},
		{
			name:           "before genesis",
			genesisTime:    1606824023,
			timestamp:      1606824000,
			secondsPerSlot: 12,
			want:           0,
		},
		{
			name:           "12 seconds after genesis",
			genesisTime:    1606824023,
			timestamp:      1606824035,
			secondsPerSlot: 12,
			want:           1,
		},
		{
			name:           "mid slot",
			genesisTime:    1606824023,
			timestamp:      1606824030, // 7 seconds after genesis
			secondsPerSlot: 12,
			want:           0,
		},
		{
			name:           "slot 100",
			genesisTime:    1606824023,
			timestamp:      1606825223,
			secondsPerSlot: 12,
			want:           100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComputeSlotAtTimestamp(tt.genesisTime, tt.timestamp, tt.secondsPerSlot)
			if got != tt.want {
				t.Errorf("ComputeSlotAtTimestamp() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestComputeTimestampAndSlotRoundTrip(t *testing.T) {
	genesisTime := uint64(1606824023)
	secondsPerSlot := uint64(12)

	// Test round trip: slot -> timestamp -> slot
	for slot := uint64(0); slot < 1000; slot += 100 {
		timestamp := ComputeTimestampAtSlot(genesisTime, slot, secondsPerSlot)
		gotSlot := ComputeSlotAtTimestamp(genesisTime, timestamp, secondsPerSlot)
		if gotSlot != slot {
			t.Errorf("round trip failed for slot %d: got %d", slot, gotSlot)
		}
	}
}
