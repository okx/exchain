package keeper

import (
	"math/big"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/feesplit/types"
)

var _ evmtypes.EvmHooks = Hooks{}

// Hooks wrapper struct for fees keeper
type Hooks struct {
	k Keeper
}

// Hooks return the wrapper hooks struct for the Keeper
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// PostTxProcessing is a wrapper for calling the EVM PostTxProcessing hook on
// the module keeper
func (h Hooks) PostTxProcessing(ctx sdk.Context, st *evmtypes.StateTransition, receipt *ethtypes.Receipt) error {
	return h.k.PostTxProcessing(ctx, st, receipt)
}

// PostTxProcessing implements EvmHooks.PostTxProcessing. After each successful
// interaction with a registered contract, the contract deployer (or, if set,
// the withdraw address) receives a share from the transaction fees paid by the
// transaction sender.
func (k Keeper) PostTxProcessing(
	ctx sdk.Context,
	st *evmtypes.StateTransition,
	receipt *ethtypes.Receipt,
) error {
	// check if the fees are globally enabled
	params := k.GetParams(ctx)
	if !params.EnableFeeSplit {
		return nil
	}

	contract := st.Recipient
	if contract == nil {
		return nil
	}

	// if the contract is not registered to receive fees, do nothing
	feeSplit, found := k.GetFeeSplit(ctx, *contract)
	if !found {
		return nil
	}

	withdrawer := feeSplit.GetWithdrawerAddr()
	if len(withdrawer) == 0 {
		withdrawer = feeSplit.GetDeployerAddr()
	}

	developerShares := params.DeveloperShares
	// if the contract shares is set by proposal
	shares, found := k.GetContractShare(ctx, *contract)
	if found {
		developerShares = shares
	}

	txFee := new(big.Int).Mul(st.Price, new(big.Int).SetUint64(receipt.GasUsed))
	developerFee := sdk.NewDecFromBigIntWithPrec(txFee, sdk.Precision).Mul(developerShares)
	fees := sdk.Coins{{Denom: sdk.DefaultBondDenom, Amount: developerFee}}

	// distribute the fees to the contract deployer / withdraw address
	k.updateFeeSplitHandler(withdrawer, fees)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeDistributeDevFeeSplit,
				sdk.NewAttribute(sdk.AttributeKeySender, st.Sender.String()),
				sdk.NewAttribute(types.AttributeKeyContract, contract.String()),
				sdk.NewAttribute(types.AttributeKeyWithdrawerAddress, withdrawer.String()),
				sdk.NewAttribute(sdk.AttributeKeyAmount, developerFee.String()),
			),
		},
	)

	return nil
}
