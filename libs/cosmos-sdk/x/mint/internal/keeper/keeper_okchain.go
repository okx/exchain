package keeper

import (
	"github.com/pkg/errors"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/mint/internal/types"
)

func (k Keeper) AddYieldFarming(ctx sdk.Context, yieldAmt sdk.Coins) error {
	// todo: verify farmModuleName
	if len(k.farmModuleName) == 0 {
		return nil
	}
	return k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.farmModuleName, yieldAmt)
}

// get the minter custom
func (k Keeper) GetMinterCustom(ctx sdk.Context) (minter types.MinterCustom) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.MinterKey)
	if b != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &minter)
	}
	return
}

// set the minter custom
func (k Keeper) SetMinterCustom(ctx sdk.Context, minter types.MinterCustom) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(minter)
	store.Set(types.MinterKey, b)
}

func (k Keeper) UpdateMinterCustom(ctx sdk.Context, minter *types.MinterCustom, params types.Params) {
	var provisionAmtPerBlock sdk.Dec
	if ctx.BlockHeight() == 0 || minter.NextBlockToUpdate == 0 {
		provisionAmtPerBlock = k.GetOriginalMintedPerBlock()
	} else {
		provisionAmtPerBlock = minter.MintedPerBlock.AmountOf(params.MintDenom).Mul(params.DeflationRate)
	}

	// update new MinterCustom
	minter.MintedPerBlock = sdk.NewDecCoinsFromDec(params.MintDenom, provisionAmtPerBlock)
	minter.NextBlockToUpdate += params.DeflationEpoch * params.BlocksPerYear

	k.SetMinterCustom(ctx, *minter)
}

// GetInflation returns the inflation of the current state of OKC,
// and the calculation of inflation can be found at https://github.com/okex/oec/issues/1628.
func (k Keeper) GetInflation(ctx sdk.Context, minter *types.MinterCustom, params types.Params) sdk.Dec {
	height := uint64(ctx.BlockHeight())
	deflationNum := height / (params.DeflationEpoch * params.BlocksPerYear)
	mpb := uint64(minter.MintedPerBlock.AmountOf(params.MintDenom).TruncateInt64())
	gmpb := uint64(k.originalMintedPerBlock.TruncateInt64())
	genesisSupply := uint64(10000000)

	subNum := 2 * params.DeflationEpoch * params.BlocksPerYear * gmpb
	for i := 0; i < int(deflationNum); i++ {
		subNum /= 2
	}
	addNum := 2 * params.DeflationEpoch * params.BlocksPerYear * gmpb
	num3 := (height % (params.DeflationEpoch * params.BlocksPerYear)) * mpb
	prevSupply := genesisSupply + addNum - subNum + num3

	prevSupplyDec := sdk.NewDec(int64(prevSupply))
	inflation := minter.MintedPerBlock.AmountOf(params.MintDenom).MulInt64(int64(params.BlocksPerYear)).Quo(prevSupplyDec)
	return inflation
}

//______________________________________________________________________

// GetOriginalMintedPerBlock returns the init tokens per block.
func (k Keeper) GetOriginalMintedPerBlock() sdk.Dec {
	return k.originalMintedPerBlock
}

// SetOriginalMintedPerBlock sets the init tokens per block.
func (k Keeper) SetOriginalMintedPerBlock(originalMintedPerBlock sdk.Dec) {
	k.originalMintedPerBlock = originalMintedPerBlock
}

// ValidateMinterCustom validate minter
func ValidateOriginalMintedPerBlock(originalMintedPerBlock sdk.Dec) error {
	if originalMintedPerBlock.IsNegative() {
		panic("init tokens per block must be non-negative")
	}

	return nil
}

// SetTreasures set the treasures to db
func (k Keeper) SetTreasures(ctx sdk.Context, treasures []types.Treasure) {
	store := ctx.KVStore(k.storeKey)
	types.SortTreasures(treasures)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(treasures)
	store.Set(types.TreasuresKey, b)
}

// GetTreasures get the treasures from db
func (k Keeper) GetTreasures(ctx sdk.Context) (treasures []types.Treasure) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.TreasuresKey)
	if b != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &treasures)
	}
	return
}

// AllocateTokenToTreasure allocate token to treasure and return remain
func (k Keeper) AllocateTokenToTreasure(ctx sdk.Context, fees sdk.Coins) (remain sdk.Coins, err error) {
	treasures := k.GetTreasures(ctx)
	remain = sdk.NewCoins()
	remain = remain.Add(fees...)
	for i, _ := range treasures {
		allocated := fees.MulDecTruncate(treasures[i].Proportion)
		if err = k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, treasures[i].Address, allocated); err != nil {
			return
		}
		remain = remain.Sub(allocated)
		if remain.IsAnyNegative() {
			return remain, errors.New("allocate coin is more than mint coin")
		}
		k.Logger(ctx).Debug("allocate treasure", "addr", treasures[i].Address, "proportion", treasures[i].Proportion, "sum coins", fees.String(), "allocated", allocated.String(), "remain", remain.String())
	}
	return
}

func (k Keeper) UpdateTreasures(ctx sdk.Context, treasures []types.Treasure) error {
	src := k.GetTreasures(ctx)
	result := types.InsertAndUpdateTreasures(src, treasures)
	if err := types.ValidateTreasures(result); err != nil {
		return err
	}
	k.SetTreasures(ctx, result)
	return nil
}

func (k Keeper) DeleteTreasures(ctx sdk.Context, treasures []types.Treasure) error {
	src := k.GetTreasures(ctx)
	result, err := types.DeleteTreasures(src, treasures)
	if err != nil {
		return err
	}
	if err := types.ValidateTreasures(result); err != nil {
		return err
	}
	k.SetTreasures(ctx, result)
	return nil
}
