package v0_9

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkparams "github.com/cosmos/cosmos-sdk/x/params"
)

// const
const (
	ModuleName = sdkparams.ModuleName
)

type (
	// Params is the struct of the parameters in this module
	Params struct {
		// DexList proposal params
		// Maximum period for okb holders to deposit on a dex list proposal. Initial value: 2 days
		MaxDepositPeriod time.Duration `json:"max_deposit_period"`
		// Minimum deposit for a critical dex list proposal to enter voting period
		MinDeposit sdk.DecCoins `json:"min_deposit"`
		// Length of the critical voting period for dex list proposal
		VotingPeriod time.Duration `json:"voting_period"`
		// block height for dex list can not be greater than DexListMaxBlockHeight
		MaxBlockHeight uint64 `json:"max_block_height"`
	}

	// GenesisState is the struct of the genesis state in this module
	GenesisState struct {
		Params Params `json:"params" yaml:"params"`
	}
)
