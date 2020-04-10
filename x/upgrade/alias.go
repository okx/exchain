package upgrade

import (
	"github.com/okex/okchain/x/upgrade/keeper"
	"github.com/okex/okchain/x/upgrade/types"
)

const (
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	RouterKey         = types.RouterKey
	QuerierRoute      = types.QuerierRoute
	DefaultParamspace = types.DefaultParamspace
)

var (
	// functions aliases
	NewQuerier     = keeper.NewQuerier
	NewKeeper      = keeper.NewKeeper
	NewVersionInfo = types.NewVersionInfo
	NewAppUpgradeProposalHandler = keeper.NewAppUpgradeProposalHandler

	// variable aliases
	ModuleCdc                    = types.ModuleCdc
	EventTypeUpgradeAppVersion   = types.EventTypeUpgradeAppVersion
	EventTypeUpgradeFailure      = types.EventTypeUpgradeFailure
	AttributeKeyAppVersion       = types.AttributeKeyAppVersion
)

type (
	Keeper = keeper.Keeper
	VersionInfo = types.VersionInfo
)
