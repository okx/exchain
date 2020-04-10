package upgrade

import (
	"fmt"
	"strconv"

	"github.com/okex/okchain/x/upgrade/types"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	flagForQuickTest bool
)

// do signal maintenance
func EndBlocker(ctx sdk.Context, keeper Keeper) {

	ctx = ctx.WithLogger(ctx.Logger().With("handler", "endBlock").With("module", "upgrade"))
	ctx.Logger().Info("enter into upgrade endblocker")
	upgradeConfig, ok := keeper.GetAppUpgradeConfig(ctx)
	if ok {
		validator, found := keeper.GetValidatorByConsAddr(ctx, (sdk.ConsAddress)(ctx.BlockHeader().ProposerAddress))
		if !found {
			panic(fmt.Sprintf("validator with consensus-address %s not found", (sdk.ConsAddress)(ctx.BlockHeader().ProposerAddress).String()))
		}

		ctx.Logger().Info(fmt.Sprintf("current version: %v, upgrade version: %v\n", ctx.BlockHeader().Version.App, upgradeConfig.ProtocolDef.Version))

		///////////////////// original code /////////////////////////////////////////
		//if ctx.BlockHeader().Version.App == upgradeConfig.ProtocolDef.Version {
		///////////////////////////////////////////////////////////////////////////
		// modified for upgrade quick test, plz recover it later!!!!
		if (ctx.BlockHeader().Version.App == upgradeConfig.ProtocolDef.Version) || flagForQuickTest {
			keeper.SetSignal(ctx, upgradeConfig.ProtocolDef.Version, validator.ConsAddress().String())

			ctx.Logger().Info("Validator has downloaded the latest software ",
				"validator", validator.GetOperator().String(), "version", upgradeConfig.ProtocolDef.Version)

		} else {

			ok := keeper.DeleteSignal(ctx, upgradeConfig.ProtocolDef.Version, validator.ConsAddress().String())

			if ok {
				ctx.Logger().Info("Validator has restarted the old software ",
					"validator", validator.GetOperator().String(), "version", upgradeConfig.ProtocolDef.Version)
			}
		}

		// tally
		if uint64(ctx.BlockHeight())+1 == upgradeConfig.ProtocolDef.Height {
			success := tally(ctx, upgradeConfig.ProtocolDef.Version, keeper, upgradeConfig.ProtocolDef.Threshold)

			if success {
				ctx.Logger().Info("Software Upgrade is successful.", "version", upgradeConfig.ProtocolDef.Version)
				keeper.SetCurrentVersion(ctx, upgradeConfig.ProtocolDef.Version)
				// just added 4 upgrade quick test,plz remove it later!!!
				flagForQuickTest = false
				//////////////////////////////////////////////////////////
			} else {
				ctx.Logger().Info("Software Upgrade is failure.", "version", upgradeConfig.ProtocolDef.Version)
				keeper.SetLastFailedVersion(ctx, upgradeConfig.ProtocolDef.Version)
			}

			keeper.AddNewVersionInfo(ctx, types.NewVersionInfo(upgradeConfig, success))
			keeper.ClearSignals(ctx, upgradeConfig.ProtocolDef.Version)
			keeper.ClearUpgradeConfig(ctx)
		}
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeUpgradeAppVersion,
			sdk.NewAttribute(AttributeKeyAppVersion, strconv.FormatUint(keeper.GetCurrentVersion(ctx), 10))))

}

// just 4 test
///////////////////////////////////////////////////////////////////////
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		logger := ctx.Logger().With("module", "upgrade")
		switch msg := msg.(type) {
		case types.MsgUpgradeConfig:
			return handleMsgUpgradeConfig(ctx, keeper, msg, logger)
		default:
			errMsg := fmt.Sprintf("Unrecognized msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgUpgradeConfig(ctx sdk.Context, keeper Keeper, msg types.MsgUpgradeConfig, logger log.Logger) sdk.Result {
	if logger != nil {
		logger.Info("handle msg upgrade config")
	}
	err := keeper.SetAppUpgradeConfig(ctx, msg.ProposalID, msg.Version, msg.Height, msg.Software)
	if err != nil {
		return err.Result()
	} else {
		///////////////////////////////////////////////////////
		// just 4 upgrade quick test,plz remove it later!!!
		flagForQuickTest = true
		///////////////////////////////////////////////////////
		return sdk.Result{}
	}
}

///////////////////////////////////////////////////////////////////////
