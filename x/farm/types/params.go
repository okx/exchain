package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/params"
	"time"
)

// Default parameter namespace
const (
	DefaultParamspace = ModuleName
	// TODO: Define your default parameters
)

// Parameter store keys
var (
	KeyMaxDepositPeriod = []byte("MaxDepositPeriod")
	KeyMinDeposit       = []byte("MinDeposit")
	KeyVotingPeriod     = []byte("VotingPeriod")
	KeyQuoteToken       = []byte("QuoteToken")
)

// ParamKeyTable for farm module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Params - used for initializing default parameter for farm at genesis
type Params struct {
	MaxDepositPeriod time.Duration `json:"max_deposit_period"`
	MinDeposit       sdk.DecCoins  `json:"min_deposit"`
	VotingPeriod     time.Duration `json:"voting_period"`

	QuoteToken string `json:"quote_token"`
}

// NewParams creates a new Params object
func NewParams( /* TODO: Pass in the paramters*/) Params {
	return Params{
		// TODO: Create your Params Type
	}
}

// String implements the stringer interface for Params
func (p Params) String() string {
	return fmt.Sprintf(`
	// TODO: Return all the params as a string
	`)
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyMaxDepositPeriod, Value: &p.MaxDepositPeriod},
		{Key: KeyMinDeposit, Value: &p.MinDeposit},
		{Key: KeyVotingPeriod, Value: &p.VotingPeriod},

		{Key: KeyQuoteToken, Value: &p.QuoteToken},
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams( /* TODO: Pass in your default Params */)
}
