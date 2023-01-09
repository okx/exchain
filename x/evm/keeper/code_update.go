package keeper

import (
	"encoding/hex"
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/types"
)

var (
	EventTypeContractUpdateByProposal = "contract-update-by-proposal"
)

// UpdateContractBytecode update contract bytecode
func (k *Keeper) UpdateContractBytecode(ctx sdk.Context, p types.ManagerContractByteCodeProposal) sdk.Error {
	contract := ethcmn.BytesToAddress(p.Contract)
	substituteContract := ethcmn.BytesToAddress(p.SubstituteContract)

	revertContractByteCode := p.Contract.String() == p.SubstituteContract.String()

	preCode := k.EvmStateDb.GetCode(contract)
	contractAcc := k.EvmStateDb.GetAccount(contract)
	if contractAcc == nil {
		return types.ErrNotContracAddress(fmt.Errorf("%s", contract.String()))
	}
	preCodeHash := contractAcc.CodeHash

	var newCode []byte
	if revertContractByteCode {
		newCode = k.getInitContractCode(ctx, p.Contract)
		if len(newCode) == 0 {
			return types.ErrContractCodeNotBeenUpdated(contract.String())
		}
	} else {
		newCode = k.EvmStateDb.GetCode(substituteContract)
	}
	// update code
	k.EvmStateDb.SetCode(contract, newCode)

	// store init code
	k.storeInitContractCode(ctx, p.Contract, preCode)

	// commit evm state db
	k.EvmStateDb.Commit(false)

	return k.AfterUpdateContractByteCode(ctx, contract, substituteContract, preCodeHash, preCode, newCode)
}

func (k *Keeper) AfterUpdateContractByteCode(ctx sdk.Context, contract, substituteContract ethcmn.Address, preCodeHash, preCode, newCode []byte) error {
	contractAfterUpdateCode := k.EvmStateDb.GetAccount(contract)
	if contractAfterUpdateCode == nil {
		return types.ErrNotContracAddress(fmt.Errorf("%s", contractAfterUpdateCode.String()))
	}

	// log
	k.logger.Info("updateContractByteCode", "contract", contract, "preCodeHash", hex.EncodeToString(preCodeHash), "preCodeSize", len(preCode),
		"codeHashAfterUpdateCode", hex.EncodeToString(contractAfterUpdateCode.CodeHash), "codeSizeAfterUpdateCode", len(newCode))
	// emit event
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		EventTypeContractUpdateByProposal,
		sdk.NewAttribute("contract", contract.String()),
		sdk.NewAttribute("preCodeHash", hex.EncodeToString(preCodeHash)),
		sdk.NewAttribute("preCodeSize", fmt.Sprintf("%d", len(preCode))),
		sdk.NewAttribute("SubstituteContract", substituteContract.String()),
		sdk.NewAttribute("codeHashAfterUpdateCode", hex.EncodeToString(contractAfterUpdateCode.CodeHash)),
		sdk.NewAttribute("codeSizeAfterUpdateCode", fmt.Sprintf("%d", len(newCode))),
	))
	// update watcher
	k.EvmStateDb.WithContext(ctx).IteratorCode(func(addr ethcmn.Address, c types.CacheCode) bool {
		ctx.GetWatcher().SaveContractCode(addr, c.Code, uint64(ctx.BlockHeight()))
		ctx.GetWatcher().SaveContractCodeByHash(c.CodeHash, c.Code)
		ctx.GetWatcher().SaveAccount(contractAfterUpdateCode)
		return true
	})
	return nil
}

func (k *Keeper) storeInitContractCode(ctx sdk.Context, addr sdk.AccAddress, code []byte) {
	store := k.paramSpace.CustomKVStore(ctx)
	key := types.GetInitContractCodeKey(addr)
	if !store.Has(key) {
		store.Set(key, code)
	}
}

func (k *Keeper) getInitContractCode(ctx sdk.Context, addr sdk.AccAddress) []byte {
	store := k.paramSpace.CustomKVStore(ctx)
	key := types.GetInitContractCodeKey(addr)
	return store.Get(key)
}
