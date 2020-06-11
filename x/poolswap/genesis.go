package poolswap

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/poolswap/types"
)

type GenesisState struct {
	Params               Params          `json:"params"`
	SwapTokenPairRecords []SwapTokenPair `json:"swap_token_pair_records"`
}

func NewGenesisState(swapTokenPairRecords []SwapTokenPair) GenesisState {
	return GenesisState{SwapTokenPairRecords: nil}
}

func ValidateGenesis(data GenesisState) error {
	for _, record := range data.SwapTokenPairRecords {
		if !record.QuotePooledCoin.IsValid() {
			return fmt.Errorf("invalid SwapTokenPairRecord: QuotePooledCoin: %s", record.QuotePooledCoin.String())
		}
		if !record.BasePooledCoin.IsValid() {
			return fmt.Errorf("invalid SwapTokenPairRecord: BasePooledCoin: %s", record.BasePooledCoin)
		}
		if record.PoolTokenName == "" {
			return fmt.Errorf("invalid SwapTokenPairRecord: PoolToken: %s. Error: Missing PoolToken", record.PoolTokenName)
		}
	}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:               types.DefaultParams(),
		SwapTokenPairRecords: []SwapTokenPair{},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)
	for _, record := range data.SwapTokenPairRecords {
		keeper.SetSwapTokenPair(ctx, record.TokenPairName(), record)
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var records []SwapTokenPair
	iterator := k.GetSwapTokenPairsIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {

		quote := string(iterator.Key())
		swapTokenPair, error := k.GetSwapTokenPair(ctx, quote)
		if nil != error {

		}
		records = append(records, swapTokenPair)

	}
	return GenesisState{SwapTokenPairRecords: records}
}
