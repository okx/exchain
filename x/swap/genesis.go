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
		if record.Quote == "" {
			return fmt.Errorf("invalid SwapTokenPairRecord: Quote: %s. Error: Missing Quote", record.Quote)
		}
		if record.QuoteAmount.IsPositive() {
			return fmt.Errorf("invalid SwapTokenPairRecord: QuoteAmount: %s. Error: Missing QuoteAmount", record.QuoteAmount)
		}
		if record.BaseAmount.IsPositive() {
			return fmt.Errorf("invalid SwapTokenPairRecord: BaseAmount: %s. Error: Missing BaseAmount", record.BaseAmount)
		}
		if record.PoolToken == "" {
			return fmt.Errorf("invalid SwapTokenPairRecord: PoolToken: %s. Error: Missing PoolToken", record.PoolToken)
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
		keeper.SetSwapTokenPair(ctx, record.Quote, record)
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
