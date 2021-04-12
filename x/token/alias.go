package token

import "github.com/okex/exchain/x/token/types"

const (
	// ModuleName is the module name constant used in many places
	ModuleName = types.ModuleName
	// StoreKey is the store key string for token module
	StoreKey = types.StoreKey
	// RouterKey is the message route for token module
	RouterKey = types.RouterKey
	// QuerierRoute is the querier route for token module
	QuerierRoute = types.QuerierRoute
	// DefaultParamspace is the param space for token module
	DefaultParamspace = types.DefaultParamspace
	// DefaultCodespace is the code space for token module
	DefaultCodespace = types.DefaultCodespace
	// KeyLock key for token lock store
	KeyLock = types.KeyLock
	// KeyMint key for token mint store
	KeyMint = types.KeyMint

	// CodeInvalidAsset error code of invalid asset
	CodeInvalidAsset = types.CodeInvalidAsset
)

type (
	// Params token params
	Params = types.Params
	// MsgSend send token message
	MsgSend = types.MsgSend
	// AccountResponse response for query account
	AccountResponse = types.AccountResponse
	// CoinInfo coin info for query token
	CoinInfo = types.CoinInfo
	// nolint
	FeeDetail = types.FeeDetail
	CoinsInfo = types.CoinsInfo
	Token     = types.Token
)

var (
	// RegisterCodec register module codec
	RegisterCodec = types.RegisterCodec
)
