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
	oldEthAddr := ethcmn.BytesToAddress(p.OldContractAddr)

	newCode := k.EvmStateDb.GetCode(ethcmn.BytesToAddress(p.NewContractAddr))

	oldCode := k.EvmStateDb.GetCode(oldEthAddr)
	oldAcc := k.EvmStateDb.GetAccount(oldEthAddr)
	if oldAcc == nil {
		return fmt.Errorf("unexcepted behavior: oldAcc %s  is null", oldEthAddr.String())
	}
	oldCodeHash := oldAcc.CodeHash

	// update code
	k.EvmStateDb.SetCode(oldEthAddr, newCode)
	// commit evm state db
	k.EvmStateDb.Commit(false)

	oldAccAfterUpdateCode := k.EvmStateDb.GetAccount(oldEthAddr)
	if oldAccAfterUpdateCode == nil {
		return fmt.Errorf("unexcepted behavior: oldAccAfterUpdateCode %s is null", oldEthAddr.String())
	}

	// log
	k.logger.Info("updateContractByteCode", "oldCodeHash", hex.EncodeToString(oldCodeHash), "oldCodeSize", len(oldCode),
		"oldCodeHashAfterUpdateCode", hex.EncodeToString(oldAccAfterUpdateCode.CodeHash), "oldCodeSizeAfterUpdateCode", len(newCode))
	// emit event
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		EventTypeContractUpdateByProposal,
		sdk.NewAttribute("oldContract", p.OldContractAddr.String()),
		sdk.NewAttribute("oldCodeHash", hex.EncodeToString(oldCodeHash)),
		sdk.NewAttribute("oldCodeSize", fmt.Sprintf("%d", len(oldCode))),
		sdk.NewAttribute("newContract", p.NewContractAddr.String()),
		sdk.NewAttribute("oldCodeHashAfterUpdateCode", hex.EncodeToString(oldAccAfterUpdateCode.CodeHash)),
		sdk.NewAttribute("oldCodeSizeAfterUpdateCode", fmt.Sprintf("%d", len(newCode))),
	))
	// update watcher
	k.EvmStateDb.WithContext(ctx).IteratorCode(func(addr ethcmn.Address, c types.CacheCode) bool {
		ctx.GetWatcher().SaveContractCode(addr, c.Code, uint64(ctx.BlockHeight()))
		ctx.GetWatcher().SaveContractCodeByHash(c.CodeHash, c.Code)
		ctx.GetWatcher().SaveAccount(oldAccAfterUpdateCode)
		return true
	})
	return nil
}
