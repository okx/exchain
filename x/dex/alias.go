// nolint
// aliases generated for the following subdirectories:
// ALIASGEN: github.com/okex/exchain/x/dex/keeper
// ALIASGEN: github.com/okex/exchain/x/dex/types
package dex

import (
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/x/common/version"
	"github.com/okex/exchain/x/dex/keeper"
	"github.com/okex/exchain/x/dex/types"
)

const (
	ModuleName        = types.ModuleName
	DefaultCodespace  = types.DefaultCodespace
	DefaultParamspace = types.DefaultParamspace
	TokenPairStoreKey = types.TokenPairStoreKey
	QuerierRoute      = types.QuerierRoute
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey

	DefaultMaxPriceDigitSize    = types.DefaultMaxPriceDigitSize
	DefaultMaxQuantityDigitSize = types.DefaultMaxQuantityDigitSize

	AuthFeeCollector = auth.FeeCollectorName
)

type (
	// Keepers
	Keeper              = keeper.Keeper
	IKeeper             = keeper.IKeeper
	SupplyKeeper        = keeper.SupplyKeeper
	TokenKeeper         = keeper.TokenKeeper
	StakingKeeper       = keeper.StakingKeeper
	BankKeeper          = keeper.BankKeeper
	ProtocolVersionType = version.ProtocolVersionType
	StreamKeeper        = keeper.StreamKeeper

	// Messages
	MsgList              = types.MsgList
	MsgDeposit           = types.MsgDeposit
	MsgWithdraw          = types.MsgWithdraw
	MsgTransferOwnership = types.MsgTransferOwnership
	MsgConfirmOwnership  = types.MsgConfirmOwnership
	MsgUpdateOperator    = types.MsgUpdateOperator
	MsgCreateOperator    = types.MsgCreateOperator

	TokenPair     = types.TokenPair
	Params        = types.Params
	WithdrawInfo  = types.WithdrawInfo
	WithdrawInfos = types.WithdrawInfos
	DEXOperator   = types.DEXOperator
	DEXOperators  = types.DEXOperators
)

var (
	ModuleCdc               = types.ModuleCdc
	DefaultTokenPairDeposit = types.DefaultTokenPairDeposit

	RegisterCodec       = types.RegisterCodec
	NewQuerier          = keeper.NewQuerier
	NewKeeper           = keeper.NewKeeper
	GetBuiltInTokenPair = keeper.GetBuiltInTokenPair
	DefaultParams       = types.DefaultParams

	NewMsgList     = types.NewMsgList
	NewMsgDeposit  = types.NewMsgDeposit
	NewMsgWithdraw = types.NewMsgWithdraw

	ErrTokenPairNotFound = types.ErrTokenPairNotFound
)
