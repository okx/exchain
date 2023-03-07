package erc20

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	"github.com/okx/okbchain/x/erc20/types"
)

// InitGenesis initializes genesis state based on exported genesis
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) []abci.ValidatorUpdate {
	k.SetParams(ctx, data.Params)

	for _, m := range data.TokenMappings {
		if !types.IsValidIBCDenom(m.Denom) {
			panic(fmt.Sprintf("Invalid denom to map to contract: %s", m.Denom))
		}
		if !common.IsHexAddress(m.Contract) {
			panic(fmt.Sprintf("Invalid contract address: %s", m.Contract))
		}
		if err := k.SetContractForDenom(ctx, m.Denom, common.HexToAddress(m.Contract)); err != nil {
			panic(err)
		}
	}

	k.InitInternalTemplateContract(ctx)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis exports genesis state of the erc20 module
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return GenesisState{
		Params:        k.GetParams(ctx),
		TokenMappings: k.GetContracts(ctx),
	}
}
