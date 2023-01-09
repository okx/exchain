package keeper

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/x/params"
	"github.com/zjg555543/movevm/api"
	"io/ioutil"

	"github.com/okex/exchain/x/move/types"
	"github.com/zjg555543/movevm"
)

// Keeper of the move store
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	paramSpace params.Subspace
}

// NewKeeper creates a new move Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace) Keeper {

	// ensure move module account is set
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:   key,
		cdc:        cdc,
		paramSpace: paramSpace,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ShortUseByCli)
}

// SetWithdrawAddr sets a new address that will receive the rewards upon withdrawal
func (k Keeper) PublishMove(ctx sdk.Context, delegatorAddr sdk.AccAddress, movePath string) error {

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePublishMove,
			sdk.NewAttribute(types.AttributeKeyPublishMove, movePath),
		),
	)

	version, _ := movevm.Version()
	fmt.Println("finished", version)

	gasMeter := api.NewMockGasMeter(uint64(500_000_000_000))

	store := api.NewLookup(gasMeter)

	pathList := [...]string{
		"/Users/oker/workspace/move/movevm/contracts/readme/build/readme/bytecode_modules/dependencies/MoveNursery/debug.mv",
		"/Users/oker/workspace/move/movevm/contracts/readme/build/readme/bytecode_modules/dependencies/MoveStdlib/fixed_point32.mv",
		"/Users/oker/workspace/move/movevm/contracts/readme/build/readme/bytecode_modules/dependencies/MoveStdlib/hash.mv",
		"/Users/oker/workspace/move/movevm/contracts/readme/build/readme/bytecode_modules/dependencies/MoveStdlib/vector.mv",
		"/Users/oker/workspace/move/movevm/contracts/readme/build/readme/bytecode_modules/dependencies/MoveNursery/errors.mv", "/Users/oker/workspace/move/movevm/contracts/readme/build/readme/bytecode_modules/dependencies/MoveNursery/acl.mv", "/Users/oker/workspace/move/movevm/contracts/readme/build/readme/bytecode_modules/dependencies/MoveStdlib/option.mv",
		"/Users/oker/workspace/move/movevm/contracts/readme/build/readme/bytecode_modules/dependencies/MoveStdlib/ascii.mv",
		"/Users/oker/workspace/move/movevm/contracts/readme/build/readme/bytecode_modules/dependencies/MoveStdlib/bit_vector.mv",
		"/Users/oker/workspace/move/movevm/contracts/readme/build/readme/bytecode_modules/dependencies/MoveStdlib/signer.mv",
		"/Users/oker/workspace/move/movevm/contracts/readme/build/readme/bytecode_modules/dependencies/MoveStdlib/error.mv",
		"/Users/oker/workspace/move/movevm/contracts/readme/build/readme/bytecode_modules/dependencies/MoveNursery/capability.mv",
		"/Users/oker/workspace/move/movevm/contracts/readme/build/readme/bytecode_modules/dependencies/MoveNursery/compare.mv",
		"/Users/oker/workspace/move/movevm/contracts/readme/build/readme/bytecode_modules/dependencies/MoveNursery/guid.mv",
		"/Users/oker/workspace/move/movevm/contracts/readme/build/readme/bytecode_modules/dependencies/MoveStdlib/bcs.mv",
		"/Users/oker/workspace/move/movevm/contracts/readme/build/readme/bytecode_modules/dependencies/MoveNursery/event.mv",
		"/Users/oker/workspace/move/movevm/contracts/readme/build/readme/bytecode_modules/dependencies/MoveNursery/offer.mv",
		"/Users/oker/workspace/move/movevm/contracts/readme/build/readme/bytecode_modules/dependencies/MoveNursery/role.mv",
		"/Users/oker/workspace/move/movevm/contracts/readme/build/readme/bytecode_modules/dependencies/MoveStdlib/string.mv",
		"/Users/oker/workspace/move/movevm/contracts/readme/build/readme/bytecode_modules/dependencies/MoveStdlib/type_name.mv",
		"/Users/oker/workspace/move/movevm/contracts/readme/build/readme/bytecode_modules/dependencies/MoveNursery/vault.mv",
	}
	for _, s := range pathList {
		movevm.Publish(readModule(s), []byte("0x1"), []byte("1234567890"), gasMeter, store, nil, nil, 10000, false)
	}
	logger := k.Logger(ctx)
	logger.Info("Publish move contract ok.......")
	return nil
}

func readModule(path string) []byte {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("read fail", err)
		return nil
	}

	return f
}
