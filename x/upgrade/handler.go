package upgrade

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/upgrade/types"
)

// EndBlocker does signal maintenance in the end of block
func EndBlocker(ctx sdk.Context, keeper Keeper) {
	logger := ctx.Logger().With("module", types.ModuleName)
	upgradeConfig, ok := keeper.GetAppUpgradeConfig(ctx)
	if ok {
		validator, found := keeper.GetValidatorByConsAddr(ctx, (sdk.ConsAddress)(ctx.BlockHeader().ProposerAddress))
		if !found {
			panic(fmt.Sprintf("validator with consensus-address %s not found",
				(sdk.ConsAddress)(ctx.BlockHeader().ProposerAddress).String()))
		}

		logger.Info(fmt.Sprintf("current version: %v, upgrade version: %v\n",
			ctx.BlockHeader().Version.App, upgradeConfig.ProtocolDef.Version))

		if ctx.BlockHeader().Version.App == upgradeConfig.ProtocolDef.Version {
			keeper.SetSignal(ctx, upgradeConfig.ProtocolDef.Version, validator.ConsAddress().String())

			logger.Info("Validator has downloaded the latest software ",
				"validator", validator.GetOperator().String(), "version", upgradeConfig.ProtocolDef.Version)

		} else {
			ok := keeper.DeleteSignal(ctx, upgradeConfig.ProtocolDef.Version, validator.ConsAddress().String())
			if ok {
				logger.Info("Validator has restarted the old software ",
					"validator", validator.GetOperator().String(), "version", upgradeConfig.ProtocolDef.Version)
			}
		}
		// tally
		if uint64(ctx.BlockHeight())+1 == upgradeConfig.ProtocolDef.Height {
			success := tally(ctx, upgradeConfig.ProtocolDef.Version, keeper, upgradeConfig.ProtocolDef.Threshold)
			if success {
				logger.Info("Software Upgrade is successful.", "version", upgradeConfig.ProtocolDef.Version)
				keeper.SetCurrentVersion(ctx, upgradeConfig.ProtocolDef.Version)
			} else {
				logger.Info("Software Upgrade is failure.", "version", upgradeConfig.ProtocolDef.Version)
				keeper.SetLastFailedVersion(ctx, upgradeConfig.ProtocolDef.Version)
			}

			keeper.AddNewVersionInfo(ctx, types.NewVersionInfo(upgradeConfig, success))
			keeper.ClearSignals(ctx, upgradeConfig.ProtocolDef.Version)
			keeper.ClearUpgradeConfig(ctx)
		}
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(EventTypeUpgradeAppVersion, sdk.NewAttribute(AttributeKeyAppVersion,
			strconv.FormatUint(keeper.GetCurrentVersion(ctx), 10))))
}
