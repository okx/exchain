package keeper

import (
	"fmt"
	"time"

	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/okex/okexchain/x/common"
	"github.com/okex/okexchain/x/gov/types"
	"github.com/okex/okexchain/x/staking/exported"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/params"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper defines governance keeper
type Keeper struct {
	// The reference to the Param Keeper to get and set Global Params
	paramsKeeper params.Keeper

	// The reference to the Paramstore to get and set gov specific params
	paramSpace params.Subspace

	// The SupplyKeeper to reduce the supply of the network
	supplyKeeper SupplyKeeper

	// The reference to the DelegationSet and ValidatorSet to get information about validators and delegators
	sk StakingKeeper

	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The codec codec for binary encoding/decoding.
	cdc *codec.Codec

	// Reserved codespace
	codespace string

	// Proposal router
	router Router

	// name of the FeeCollector ModuleAccount
	feeCollectorName string

	// The reference to the CoinKeeper to modify balances
	bankKeeper BankKeeper

	// Proposal module parameter router
	proposalHandlerRouter ProposalHandlerRouter
}

// NewKeeper returns a governance keeper. It handles:
// - submitting governance proposals
// - depositing funds into proposals, and activating upon sufficient funds being deposited
// - users voting on proposals, with weight proportional to stake in the system
// - and tallying the result of the vote.
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramsKeeper params.Keeper, paramSpace params.Subspace,
	supplyKeeper SupplyKeeper, sk StakingKeeper, codespace string, rtr Router,
	ck BankKeeper, phr ProposalHandlerRouter, feeCollectorName string,
) Keeper {
	// ensure governance module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// It is vital to seal the governance proposal router here as to not allow
	// further handlers to be registered after the keeper is created since this
	// could create invalid or non-deterministic behavior.
	rtr.Seal()

	keeper := Keeper{
		storeKey:              key,
		paramsKeeper:          paramsKeeper,
		paramSpace:            paramSpace.WithKeyTable(types.ParamKeyTable()),
		supplyKeeper:          supplyKeeper,
		sk:                    sk,
		cdc:                   cdc,
		codespace:             codespace,
		router:                rtr,
		feeCollectorName:      feeCollectorName,
		bankKeeper:            ck,
		proposalHandlerRouter: phr,
	}
	keeper.proposalHandlerRouter = keeper.proposalHandlerRouter.AddRoute(types.RouterKey, keeper)
	return keeper
}

// BankKeeper returns bank keeper in gov keeper
func (keeper Keeper) BankKeeper() BankKeeper {
	return keeper.bankKeeper
}

// ProposalHandlerRouter returns proposal handler router  in gov keeper
func (keeper Keeper) ProposalHandlerRouter() ProposalHandlerRouter {
	return keeper.proposalHandlerRouter
}

// Params

// GetDepositParams returns the current DepositParams from the global param store
func (keeper Keeper) GetDepositParams(ctx sdk.Context) types.DepositParams {
	var depositParams types.DepositParams
	keeper.paramSpace.Get(ctx, types.ParamStoreKeyDepositParams, &depositParams)
	return depositParams
}

// GetVotingParams returns the current VotingParams from the global param store
func (keeper Keeper) GetVotingParams(ctx sdk.Context) types.VotingParams {
	var votingParams types.VotingParams
	keeper.paramSpace.Get(ctx, types.ParamStoreKeyVotingParams, &votingParams)
	return votingParams
}

// GetTallyParams returns the current TallyParams from the global param store
func (keeper Keeper) GetTallyParams(ctx sdk.Context) types.TallyParams {
	var tallyParams types.TallyParams
	keeper.paramSpace.Get(ctx, types.ParamStoreKeyTallyParams, &tallyParams)
	return tallyParams
}

// SetDepositParams sets the current DepositParams to the global param store
func (keeper Keeper) SetDepositParams(ctx sdk.Context, depositParams types.DepositParams) {
	keeper.paramSpace.Set(ctx, types.ParamStoreKeyDepositParams, &depositParams)
}

// SetVotingParams sets the current VotingParams to the global param store
func (keeper Keeper) SetVotingParams(ctx sdk.Context, votingParams types.VotingParams) {
	keeper.paramSpace.Set(ctx, types.ParamStoreKeyVotingParams, &votingParams)
}

// SetTallyParams sets the current TallyParams to the global param store
func (keeper Keeper) SetTallyParams(ctx sdk.Context, tallyParams types.TallyParams) {
	keeper.paramSpace.Set(ctx, types.ParamStoreKeyTallyParams, &tallyParams)
}

// ProposalQueues

// WaitingProposalQueueIterator returns an iterator for all the proposals in the Waiting Queue that expire by endTime
func (keeper Keeper) WaitingProposalQueueIterator(ctx sdk.Context, blockHeight uint64) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return store.Iterator(types.PrefixWaitingProposalQueue,
		sdk.PrefixEndBytes(types.WaitingProposalByBlockHeightKey(blockHeight)))
}

// InsertWaitingProposalQueue inserts a ProposalID into the waiting proposal queue at endTime
func (keeper Keeper) InsertWaitingProposalQueue(ctx sdk.Context, blockHeight, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(proposalID)
	store.Set(types.WaitingProposalQueueKey(proposalID, blockHeight), bz)
}

// RemoveFromWaitingProposalQueue removes a proposalID from the waiting Proposal Queue
func (keeper Keeper) RemoveFromWaitingProposalQueue(ctx sdk.Context, blockHeight, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(types.WaitingProposalQueueKey(proposalID, blockHeight))
}

// Iterators

// IterateProposals iterates over the all the proposals and performs a callback function
func (keeper Keeper) IterateProposals(ctx sdk.Context, cb func(proposal types.Proposal) (stop bool)) {
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ProposalsKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var proposal types.Proposal
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &proposal)

		if cb(proposal) {
			break
		}
	}
}

// IterateActiveProposalsQueue iterates over the proposals in the active proposal queue
// and performs a callback function
func (keeper Keeper) IterateActiveProposalsQueue(
	ctx sdk.Context, endTime time.Time, cb func(proposal types.Proposal,
) (stop bool)) {
	iterator := keeper.ActiveProposalQueueIterator(ctx, endTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		proposalID, _ := types.SplitActiveProposalQueueKey(iterator.Key())
		proposal, found := keeper.GetProposal(ctx, proposalID)
		if !found {
			panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}

		if cb(proposal) {
			break
		}
	}
}

// IterateInactiveProposalsQueue iterates over the proposals in the inactive proposal queue
// and performs a callback function
func (keeper Keeper) IterateInactiveProposalsQueue(
	ctx sdk.Context, endTime time.Time, cb func(proposal types.Proposal,
) (stop bool)) {
	iterator := keeper.InactiveProposalQueueIterator(ctx, endTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		proposalID, _ := types.SplitInactiveProposalQueueKey(iterator.Key())
		proposal, found := keeper.GetProposal(ctx, proposalID)
		if !found {
			panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}

		if cb(proposal) {
			break
		}
	}
}

// IterateWaitingProposalsQueue iterates over the proposals in the waiting proposal queue
// and performs a callback function
func (keeper Keeper) IterateWaitingProposalsQueue(
	ctx sdk.Context, height uint64, cb func(proposal types.Proposal,
) (stop bool)) {
	iterator := keeper.WaitingProposalQueueIterator(ctx, height)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		proposalID, _ := types.SplitWaitingProposalQueueKey(iterator.Key())
		proposal, found := keeper.GetProposal(ctx, proposalID)
		if !found {
			panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}

		if cb(proposal) {
			break
		}
	}
}

// IterateAllWaitingProposals iterates over the all proposals in the waiting proposal queue
// and performs a callback function
func (keeper Keeper) IterateAllWaitingProposals(ctx sdk.Context,
	cb func(proposal types.Proposal, proposalID, height uint64) (stop bool)) {
	store := ctx.KVStore(keeper.storeKey)
	iterator := store.Iterator(types.PrefixWaitingProposalQueue,
		sdk.PrefixEndBytes(types.PrefixWaitingProposalQueue))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		proposalID, height := types.SplitWaitingProposalQueueKey(iterator.Key())
		proposal, found := keeper.GetProposal(ctx, proposalID)
		if !found {
			panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}

		if cb(proposal, proposalID, height) {
			break
		}
	}
}

// IterateAllDeposits iterates over the all the stored deposits and performs a callback function
func (keeper Keeper) IterateAllDeposits(ctx sdk.Context, cb func(deposit types.Deposit) (stop bool)) {
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.DepositsKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var deposit types.Deposit
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &deposit)

		if cb(deposit) {
			break
		}
	}
}

// IterateDeposits iterates over the all the proposals deposits and performs a callback function
func (keeper Keeper) IterateDeposits(ctx sdk.Context, proposalID uint64, cb func(deposit types.Deposit) (stop bool)) {
	iterator := keeper.GetDepositsIterator(ctx, proposalID)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var deposit types.Deposit
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &deposit)

		if cb(deposit) {
			break
		}
	}
}

// IterateVotes iterates over the all the proposals votes and performs a callback function
func (keeper Keeper) IterateVotes(ctx sdk.Context, proposalID uint64, cb func(vote types.Vote) (stop bool)) {
	iterator := keeper.GetVotesIterator(ctx, proposalID)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var vote types.Vote
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &vote)

		if cb(vote) {
			break
		}
	}
}

// GetMinDeposit implement ProposalHandler
// nolint
func (keeper Keeper) GetMinDeposit(ctx sdk.Context, content types.Content) sdk.SysCoins {
	return keeper.GetDepositParams(ctx).MinDeposit
}

// nolint
func (keeper Keeper) GetMaxDepositPeriod(ctx sdk.Context, content types.Content) time.Duration {
	return keeper.GetDepositParams(ctx).MaxDepositPeriod
}

// nolint
func (keeper Keeper) GetVotingPeriod(ctx sdk.Context, content types.Content) time.Duration {
	return keeper.GetVotingParams(ctx).VotingPeriod
}

// nolint
func (keeper Keeper) CheckMsgSubmitProposal(ctx sdk.Context, msg types.MsgSubmitProposal) sdk.Error {
	// check initial deposit more than or equal to ratio of MinDeposit
	initDeposit := keeper.GetDepositParams(ctx).MinDeposit.MulDec(sdk.NewDecWithPrec(1, 1))
	err := common.HasSufficientCoins(msg.Proposer, msg.InitialDeposit,
		initDeposit)
	if err != nil {
		return types.ErrInitialDepositNotEnough(initDeposit.String())
	}
	// check proposer has sufficient coins
	err = common.HasSufficientCoins(msg.Proposer, keeper.bankKeeper.GetCoins(ctx, msg.Proposer),
		msg.InitialDeposit)
	if err != nil {
		return common.ErrInsufficientCoins(types.DefaultCodespace, err.Error())
	}
	return nil
}

// nolint
func (keeper Keeper) AfterSubmitProposalHandler(ctx sdk.Context, proposal types.Proposal) {}

// nolint
func (keeper Keeper) VoteHandler(ctx sdk.Context, proposal types.Proposal, vote types.Vote) (string, sdk.Error) {
	return "", nil
}

// nolint
func (keeper Keeper) AfterDepositPeriodPassed(ctx sdk.Context, proposal types.Proposal) {}

// nolint
func (keeper Keeper) RejectedHandler(ctx sdk.Context, content types.Content) {}

// get all current validators except candidate votes
func (keeper Keeper) totalPower(ctx sdk.Context) sdk.Dec {
	totalVoting := sdk.ZeroDec()
	keeper.sk.IterateBondedValidatorsByPower(ctx, func(index int64, validator exported.ValidatorI) (stop bool) {
		totalVoting = totalVoting.Add(validator.GetDelegatorShares())
		return false
	})
	return totalVoting
}

// Logger returns a module-specific logger.
func (keeper Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetGovernanceAccount returns the governance ModuleAccount
func (keeper Keeper) GetGovernanceAccount(ctx sdk.Context) supplyexported.ModuleAccountI {
	return keeper.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}

// Params

// ProposalQueues

// InsertActiveProposalQueue inserts a ProposalID into the active proposal queue at endTime
func (keeper Keeper) InsertActiveProposalQueue(ctx sdk.Context, proposalID uint64, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(proposalID)
	store.Set(types.ActiveProposalQueueKey(proposalID, endTime), bz)
}

// RemoveFromActiveProposalQueue removes a proposalID from the Active Proposal Queue
func (keeper Keeper) RemoveFromActiveProposalQueue(ctx sdk.Context, proposalID uint64, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(types.ActiveProposalQueueKey(proposalID, endTime))
}

// InsertInactiveProposalQueue Inserts a ProposalID into the inactive proposal queue at endTime
func (keeper Keeper) InsertInactiveProposalQueue(ctx sdk.Context, proposalID uint64, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(proposalID)
	store.Set(types.InactiveProposalQueueKey(proposalID, endTime), bz)
}

// RemoveFromInactiveProposalQueue removes a proposalID from the Inactive Proposal Queue
func (keeper Keeper) RemoveFromInactiveProposalQueue(ctx sdk.Context, proposalID uint64, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(types.InactiveProposalQueueKey(proposalID, endTime))
}

// IterateAllVotes iterates over the all the stored votes and performs a callback function
func (keeper Keeper) IterateAllVotes(ctx sdk.Context, cb func(vote types.Vote) (stop bool)) {
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.VotesKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var vote types.Vote
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &vote)

		if cb(vote) {
			break
		}
	}
}

// ActiveProposalQueueIterator returns an sdk.Iterator for all the proposals in the Active Queue that expire by endTime
func (keeper Keeper) ActiveProposalQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return store.Iterator(types.ActiveProposalQueuePrefix, sdk.PrefixEndBytes(types.ActiveProposalByTimeKey(endTime)))
}

// InactiveProposalQueueIterator returns an sdk.Iterator for all the proposals in the Inactive Queue that expire by endTime
func (keeper Keeper) InactiveProposalQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return store.Iterator(types.InactiveProposalQueuePrefix, sdk.PrefixEndBytes(types.InactiveProposalByTimeKey(endTime)))
}

func (keeper Keeper) Cdc() *codec.Codec {
	return keeper.cdc
}

func (keeper Keeper) Router() Router {
	return keeper.router
}

func (keeper Keeper) ProposalHandleRouter() ProposalHandlerRouter {
	return keeper.proposalHandlerRouter
}

func (keeper Keeper) SupplyKeeper() SupplyKeeper {
	return keeper.supplyKeeper
}
