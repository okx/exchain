package keeper

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/wasm/internal/types"
	// authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	// "github.com/okex/okchain/x/wasm/internal/types"
)

// InitGenesis sets supply information for genesis.
//
// CONTRACT: all types of accounts must have been already initialized/created
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	for _, code := range data.Codes {
		newId, err := keeper.Create(ctx, code.CodeInfo.Creator, code.CodesBytes, code.CodeInfo.Source, code.CodeInfo.Builder)
		if err != nil {
			panic(err)
		}
		newInfo := keeper.GetCodeInfo(ctx, newId)
		if !bytes.Equal(code.CodeInfo.CodeHash, newInfo.CodeHash) {
			panic("code hashes not same")
		}
	}

	for _, contract := range data.Contracts {
		keeper.setContractInfo(ctx, contract.ContractAddress, contract.ContractInfo)
		keeper.setContractState(ctx, contract.ContractAddress, contract.ContractState)
	}

}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	var genState types.GenesisState

	maxCodeID := keeper.GetNextCodeID(ctx)
	for i := uint64(1); i < maxCodeID; i++ {
		bytecode, err := keeper.GetByteCode(ctx, i)
		if err != nil {
			panic(err)
		}
		genState.Codes = append(genState.Codes, types.Code{
			CodeInfo:   *keeper.GetCodeInfo(ctx, i),
			CodesBytes: bytecode,
		})
	}

	keeper.ListContractInfo(ctx, func(addr sdk.AccAddress, contract types.ContractInfo) bool {
		contractStateIterator := keeper.GetContractState(ctx, addr)
		var state []types.Model
		for ; contractStateIterator.Valid(); contractStateIterator.Next() {
			m := types.Model{
				Key:   contractStateIterator.Key(),
				Value: contractStateIterator.Value(),
			}
			state = append(state, m)
		}

		genState.Contracts = append(genState.Contracts, types.Contract{
			ContractAddress: addr,
			ContractInfo:    contract,
			ContractState:   state,
		})

		return false
	})

	return genState
}
