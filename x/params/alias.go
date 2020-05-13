package params

import (
	sdkparams "github.com/cosmos/cosmos-sdk/x/params"
	"github.com/okex/okchain/x/params/types"
)

// const
const (
	ModuleName        = sdkparams.ModuleName
	DefaultCodespace  = sdkparams.DefaultCodespace
	DefaultParamspace = sdkparams.ModuleName
	StoreKey          = sdkparams.StoreKey
	TStoreKey         = sdkparams.TStoreKey
	RouterKey         = sdkparams.RouterKey
)

type (
	// KeyTable is the type alias of the one in cmsdk
	KeyTable = sdkparams.KeyTable
	// ParamSetPairs is the type alias of the one in cmsdk
	ParamSetPairs = sdkparams.ParamSetPairs
	// Subspace is the type alias of the one in cmsdk
	Subspace = sdkparams.Subspace
	// ParamSet is the type alias of the one in cmsdk
	ParamSet = sdkparams.ParamSet
	// ParamChange is the type alias of the one in cmsdk
	ParamChange = sdkparams.ParamChange
	// ParameterChangeProposal is alias of ParameterChangeProposal in types
	ParameterChangeProposal = types.ParameterChangeProposal
)

var (
	// nolint
	NewKeyTable    = sdkparams.NewKeyTable
	NewParamChange = sdkparams.NewParamChange
	DefaultParams  = types.DefaultParams
)
