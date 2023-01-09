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
	store      movevm.KVStore
	gasMeter   movevm.GasMeter
}

// NewKeeper creates a new move Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace) Keeper {
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	gasMeter := api.NewMockGasMeter(uint64(500_000_000_000))

	return Keeper{
		storeKey:   key,
		cdc:        cdc,
		paramSpace: paramSpace,
		gasMeter:   gasMeter,
		store:      api.NewLookup(gasMeter),
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ShortUseByCli)
}

func (k Keeper) PublishMove(ctx sdk.Context, delegatorAddr sdk.AccAddress, movePath string) error {

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePublishMove,
			sdk.NewAttribute(types.AttributeKeyPublishMove, movePath),
		),
	)

	version, _ := movevm.Version()
	fmt.Println("finished", version)

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
		movevm.Publish(readModule(s), []byte("0x1"), []byte("1234567890"), k.gasMeter, k.store, nil, nil, 10000, false)
	}

	testByte := []byte("1234567890")
	moduleBytes2 := readModule("/Users/oker/workspace/move/movevm/contracts/readme/build/readme/bytecode_modules/Test.mv")
	movevm.Publish(moduleBytes2, []byte("0x2"), testByte, k.gasMeter, k.store, nil, nil, 10000, false)

	moduleBytes3 := readModule("/Users/oker/workspace/move/movevm/contracts/readme/build/readme/bytecode_modules/Caller.mv")
	movevm.Publish(moduleBytes3, []byte("0x3"), testByte, k.gasMeter, k.store, nil, nil, 10000, false)
	logger := k.Logger(ctx)

	logger.Info("Publish move contract ok.......")
	return nil
}

func (k Keeper) RunMove(ctx sdk.Context, delegatorAddr sdk.AccAddress, movePath string) error {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRunMove,
			sdk.NewAttribute(types.AttributeKeyPublishMove, movePath),
		),
	)

	scriptBytes := readModule("/Users/oker/workspace/move/movevm/contracts/readme/build/readme/bytecode_scripts/test_script.mv")
	movevm.Run(scriptBytes, []byte("0xF"), []byte("1234567890"), k.gasMeter, k.store, nil, nil, 10000, false)

	logger := k.Logger(ctx)

	logger.Info("Run move contract ok.......")
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
