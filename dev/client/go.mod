module github.com/okex/exchain/dev/client

go 1.14

require (
	github.com/cespare/cp v1.1.1 // indirect
	github.com/cosmos/cosmos-sdk v0.39.2
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/ethereum/go-ethereum v1.10.8
	github.com/kr/pretty v0.2.0 // indirect
	github.com/okex/exchain-ethereum-compatible v1.0.2
	github.com/prometheus/tsdb v0.9.1 // indirect
)

replace github.com/cosmos/cosmos-sdk => github.com/okex/cosmos-sdk v0.39.2-exchain1
