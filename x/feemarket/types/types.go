package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

// GenesisState defines the feemarket module's genesis state.
type GenesisState struct {
	// params defines all the paramaters of the module.
	Params Params `json:"params"`
	// base fee is the exported value from previous software version.
	// Zero by default.
	BaseFee sdk.Int `json:"base_fee"`
	// block gas is the amount of gas used on the last block before the upgrade.
	// Zero by default.
	BlockGas uint64 `json:"block_gas,omitempty"`
}
