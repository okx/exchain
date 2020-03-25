package app

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/x/slashing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/okex/okchain/app/protocol"
	"github.com/okex/okchain/x/staking"
	abci "github.com/tendermint/tendermint/abci/types"

	tmtypes "github.com/tendermint/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// export the state of okchain for a genesis file
func (app *OKChainApp) ExportAppStateAndValidators(forZeroHeight bool, jailWhiteList []string,
) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	// as if they could withdraw from the start of the next block
	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})

	if forZeroHeight {
		app.prepForZeroHeightGenesis(ctx, jailWhiteList)
	}

	// Get current protocol from engine
	curProtocol := protocol.GetEngine().GetCurrentProtocol()

	//genState := app.mm.ExportGenesis(ctx)
	genesisState := curProtocol.ExportGenesis(ctx)

	//appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	appState, err = codec.MarshalJSONIndent(curProtocol.GetCodec(), genesisState)
	if err != nil {
		return nil, nil, err
	}

	//validators = staking.WriteValidators(ctx, app.stakingKeeper)
	validators = staking.WriteValidators(ctx, curProtocol.GetStakingKeeper())
	return appState, validators, nil
}

// prepare for fresh start at zero height
// NOTE zero height genesis is a temporary feature which will be deprecated
//      in favour of export at a block height
func (app *OKChainApp) prepForZeroHeightGenesis(ctx sdk.Context, jailWhiteList []string) {
	var applyWhiteList bool

	//Check if there is a whitelist
	if len(jailWhiteList) > 0 {
		applyWhiteList = true
	}

	whiteListMap := make(map[string]bool)

	for _, addr := range jailWhiteList {
		_, err := sdk.ValAddressFromBech32(addr)
		if err != nil {
			panic(err)
		}
		whiteListMap[addr] = true
	}

	// Get current protocol from engine
	curProtocol := protocol.GetEngine().GetCurrentProtocol()

	/* Just to be safe, assert the invariants on current state. */
	curProtocol.GetCrisisKeeper().AssertInvariants(ctx)

	/* Handle fee distribution state. */
	// withdraw all validator commission
	curProtocol.GetStakingKeeper().IterateValidators(ctx, func(_ int64, val staking.ValidatorI) (stop bool) {
		_, _ = curProtocol.GetDistrKeeper().WithdrawValidatorCommission(ctx, val.GetOperator())
		return false
	})

	// set context height to zero
	height := ctx.BlockHeight()
	ctx = ctx.WithBlockHeight(0)

	// reset context height
	ctx = ctx.WithBlockHeight(height)

	/* Handle staking state. */
	// Iterate through validators by power descending, reset bond heights, and
	// update bond intra-tx counters.
	store := ctx.KVStore(curProtocol.GetKVStoreKeysMap()[staking.StoreKey])
	iter := sdk.KVStoreReversePrefixIterator(store, staking.ValidatorsKey)
	counter := int16(0)

	var valConsAddrs []sdk.ConsAddress
	for ; iter.Valid(); iter.Next() {
		addr := sdk.ValAddress(iter.Key()[1:])

		validator, found := curProtocol.GetStakingKeeper().GetValidator(ctx, addr)
		if !found {
			panic("didn't find the expected validator")
		}
		validator.UnbondingHeight = 0
		valConsAddrs = append(valConsAddrs, validator.ConsAddress())
		if applyWhiteList && !whiteListMap[addr.String()] {
			validator.Jailed = true
		}

		curProtocol.GetStakingKeeper().SetValidator(ctx, validator)
		counter++
	}

	iter.Close()

	_ = curProtocol.GetStakingKeeper().ApplyAndReturnValidatorSetUpdates(ctx)

	/* Handle slashing state. */
	curProtocol.GetSlashingKeeper().IterateValidatorSigningInfos(
		ctx,
		func(addr sdk.ConsAddress, info slashing.ValidatorSigningInfo) (stop bool) {
			info.StartHeight = 0
			curProtocol.GetSlashingKeeper().SetValidatorSigningInfo(ctx, addr, info)
			return false
		})
}
