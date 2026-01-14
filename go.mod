module github.com/islishude/eth-beacon-client

go 1.25

require (
	github.com/ethereum/go-ethereum v1.16.8
	github.com/protolambda/zrnt v0.34.1
	github.com/protolambda/ztyp v0.2.2
)

require (
	github.com/bits-and-blooms/bitset v1.24.4 // indirect
	github.com/consensys/gnark-crypto v0.19.2 // indirect
	github.com/crate-crypto/go-eth-kzg v1.4.0 // indirect
	github.com/ethereum/c-kzg-4844/v2 v2.1.5 // indirect
	github.com/holiman/uint256 v1.3.2 // indirect
	github.com/kilic/bls12-381 v0.1.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.9 // indirect
	github.com/minio/sha256-simd v1.0.1 // indirect
	github.com/protolambda/bls12-381-util v0.1.0 // indirect
	github.com/supranational/blst v0.3.16 // indirect
	golang.org/x/crypto v0.47.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract (
	v0.0.2 // old package name
	v0.0.1 // include unnessary log
)
