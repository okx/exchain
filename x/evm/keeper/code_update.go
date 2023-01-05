package keeper

import (
	"encoding/hex"
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/types"
)

// UpdateContractBytecode update contract bytecode
func (k *Keeper) UpdateContractBytecode(ctx sdk.Context, p types.ManagerContractByteCodeProposal) sdk.Error {
	oldEthAddr := ethcmn.BytesToAddress(p.OldContractAddr)

	newCode := k.EvmStateDb.GetCode(ethcmn.BytesToAddress(p.NewContractAddr))

	oldCode := k.EvmStateDb.GetCode(oldEthAddr)
	oldAcc := k.EvmStateDb.GetAccount(oldEthAddr)
	if oldAcc == nil {
		return fmt.Errorf("%s is null", oldEthAddr.String())
	}
	oldCodeHash := oldAcc.CodeHash

	k.EvmStateDb.SetCode(oldEthAddr, newCode)
	k.EvmStateDb.Commit(false)

	oldAccAfterUpdateCode := k.EvmStateDb.GetAccount(oldEthAddr)
	k.logger.Info("updateContractByteCode", "oldCodeHash", hex.EncodeToString(oldCodeHash), "oldCodeSize", len(oldCode),
		"oldCodeHashAfterUpdateCode", hex.EncodeToString(oldAccAfterUpdateCode.CodeHash), "oldCodeSizeAfterUpdateCode", len(newCode))

	k.EvmStateDb.WithContext(ctx).IteratorCode(func(addr ethcmn.Address, c types.CacheCode) bool {
		ctx.GetWatcher().SaveContractCode(addr, c.Code, uint64(ctx.BlockHeight()))
		ctx.GetWatcher().SaveContractCodeByHash(c.CodeHash, c.Code)
		ctx.GetWatcher().SaveAccount(oldAccAfterUpdateCode)
		return true
	})
	k.Commit(ctx)
	return nil
}
