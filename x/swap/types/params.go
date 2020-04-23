package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/params"
)

var (
	FeeRate = sdk.NewDecWithPrec(3, 3)
)

// Default parameter namespace
const (
	DefaultParamspace = ModuleName

	GenerateTokenType = 2
)

// Parameter store keys
var (
// TODO: Define your keys for the parameter store
// KeyParamName          = []byte("ParamName")
)

// ParamKeyTable for swap module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Params - used for initializing default parameter for swap at genesis
type Params struct {
	// TODO: Add your Paramaters to the Paramter struct
	// KeyParamName string `json:"key_param_name"`
}

// NewParams creates a new Params object
func NewParams( /* TODO: Pass in the paramters*/ ) Params {
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
		// TODO: Pair your key with the param
		// params.NewParamSetPair(KeyParamName, &p.ParamName),
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams( /* TODO: Pass in your default Params */ )
}
