package erc20

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/x/erc20/types"
)

// InitGenesis initializes genesis state based on exported genesis
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) []abci.ValidatorUpdate {
	k.SetParams(ctx, data.Params)

	for _, m := range data.ExternalContracts {
		if !types.IsValidIBCDenom(m.Denom) {
			panic(fmt.Sprintf("Invalid denom to map to contract: %s", m.Denom))
		}
		if !common.IsHexAddress(m.Contract) {
			panic(fmt.Sprintf("Invalid contract address: %s", m.Contract))
		}
		if err := k.SetExternalContractForDenom(ctx, m.Denom, common.HexToAddress(m.Contract)); err != nil {
			panic(err)
		}
	}

	for _, m := range data.AutoContracts {
		if !types.IsValidIBCDenom(m.Denom) {
			panic(fmt.Sprintf("Invalid denom to map to contract: %s", m.Denom))
		}
		if !common.IsHexAddress(m.Contract) {
			panic(fmt.Sprintf("Invalid contract address: %s", m.Contract))
		}
		k.SetAutoContractForDenom(ctx, m.Denom, common.HexToAddress(m.Contract))
	}

	return []abci.ValidatorUpdate{}
}

// ExportGenesis exports genesis state of the erc20 module
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return GenesisState{
		Params:            k.GetParams(ctx),
		ExternalContracts: k.GetExternalContracts(ctx),
		AutoContracts:     k.GetAutoContracts(ctx),
	}
}
