package ibc

import (
	"github.com/okx/okbchain/libs/ibc-go/modules/core/common"
	"github.com/okx/okbchain/libs/ibc-go/modules/core/keeper"
	"github.com/okx/okbchain/libs/ibc-go/modules/core/types"
)

type (
	Keeper   = keeper.FacadedKeeper
	V2Keeper = keeper.Keeper
)

const ()

var (
	NewKeeper           = keeper.NewKeeper
	NewV4Keeper         = keeper.NewV4Keeper
	NewFacadedKeeper    = keeper.NewFacadedKeeper
	ModuleCdc           = types.ModuleCdc
	DefaultGenesisState = types.DefaultGenesisState
)

const (
	IBCV4 common.SelectVersion = 4.0
	IBCV2 common.SelectVersion = 2.0
)
