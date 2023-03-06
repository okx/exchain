package keeper

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	ethermint "github.com/okx/okbchain/app/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/tendermint/libs/log"
	"github.com/okx/okbchain/x/feesplit/types"
	"github.com/okx/okbchain/x/params"
)

// Keeper of this module maintains collections of fee splits for contracts
// registered to receive transaction fees.
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	paramSpace types.Subspace

	evmKeeper             types.EvmKeeper
	govKeeper             types.GovKeeper
	supplyKeeper          types.SupplyKeeper
	accountKeeper         types.AccountKeeper
	updateFeeSplitHandler sdk.UpdateFeeSplitHandler
}

// NewKeeper creates new instances of the fees Keeper
func NewKeeper(
	storeKey sdk.StoreKey,
	cdc *codec.Codec,
	ps params.Subspace,
	ek types.EvmKeeper,
	sk types.SupplyKeeper,
	ak types.AccountKeeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		paramSpace:    ps,
		evmKeeper:     ek,
		supplyKeeper:  sk,
		accountKeeper: ak,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetEthAccount returns an eth account.
func (k Keeper) GetEthAccount(ctx sdk.Context, addr common.Address) (*ethermint.EthAccount, bool) {
	cosmosAddr := sdk.AccAddress(addr.Bytes())
	acct := k.accountKeeper.GetAccount(ctx, cosmosAddr)
	if acct == nil {
		return nil, false
	}

	ethAcct, _ := acct.(*ethermint.EthAccount)
	return ethAcct, true
}

// SetGovKeeper sets keeper of gov
func (k *Keeper) SetGovKeeper(gk types.GovKeeper) {
	k.govKeeper = gk
}
