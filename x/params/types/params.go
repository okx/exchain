package types

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/x/params/subspace"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkparams "github.com/cosmos/cosmos-sdk/x/params"
)

const (
	QueryParams = "params"
)

// ParamKeyTable returns the key declaration for parameters
func ParamKeyTable() sdkparams.KeyTable {
	return sdkparams.NewKeyTable().RegisterParamSet(&Params{})
}

// Params is the struct of the parameters in this module
type Params struct {
	// DexList proposal params
	// Maximum period for okb holders to deposit on a dex list proposal. Initial value: 2 days
	MaxDepositPeriod time.Duration `json:"max_deposit_period"`
	// Minimum deposit for a critical dex list proposal to enter voting period
	MinDeposit sdk.SysCoins `json:"min_deposit"`
	// Length of the critical voting period for dex list proposal
	VotingPeriod time.Duration `json:"voting_period"`
	// block height for dex list can not be greater than DexListMaxBlockHeight
	MaxBlockHeight uint64 `json:"max_block_height"`
}

// DefaultParams returns the instance of Params with default value
func DefaultParams() Params {
	minDeposit := sdk.SysCoins{sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(100))}
	return Params{
		MaxDepositPeriod: time.Hour * 24,
		MinDeposit:       minDeposit,
		VotingPeriod:     time.Hour * 72,
		MaxBlockHeight:   100000,
	}
}

func (p Params) String() string {
	return fmt.Sprintf(`
MaxDepositPeriod: %s,
MinDeposit:       %s,
VotingPeriod:     %s,
MaxBlockHeight:   %d,
`, p.MaxDepositPeriod, p.MinDeposit, p.VotingPeriod, p.MaxBlockHeight)
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{KeyMaxDepositPeriod, &p.MaxDepositPeriod},
		{KeyMinDeposit, &p.MinDeposit},
		{KeyVotingPeriod, &p.VotingPeriod},
		{KeyMaxBlockHeight, &p.MaxBlockHeight},
	}
}
