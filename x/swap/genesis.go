package swap

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	SwapTokenPairRecords []SwapTokenPair `json:"whois_records"`
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
		SwapTokenPairRecords: []SwapTokenPair{},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
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
