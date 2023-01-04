package evm

import (
	"encoding/hex"
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/evm/watcher"
	govTypes "github.com/okex/exchain/x/gov/types"
	"reflect"
)

// NewManageContractDeploymentWhitelistProposalHandler handles "gov" type message in "evm"
func NewManageContractDeploymentWhitelistProposalHandler(k *Keeper) govTypes.Handler {
	return func(ctx sdk.Context, proposal *govTypes.Proposal) (err sdk.Error) {
		if watcher.IsWatcherEnabled() {
			ctx.SetWatcher(watcher.NewTxWatcher())
		}

		defer func() {
			if err == nil {
				ctx.GetWatcher().Finalize()
			}
		}()
		switch content := proposal.Content.(type) {
		case types.ManageContractDeploymentWhitelistProposal:
			return handleManageContractDeploymentWhitelistProposal(ctx, k, content)
		case types.ManageContractBlockedListProposal:
			return handleManageContractBlockedlListProposal(ctx, k, content)
		case types.ManageContractMethodBlockedListProposal:
			return handleManageContractMethodBlockedlListProposal(ctx, k, content)
		case types.ManageSysContractAddressProposal:
			if tmtypes.HigherThanVenus3(ctx.BlockHeight()) {
				return handleManageSysContractAddressProposal(ctx, k, content)
			}
			return common.ErrUnknownProposalType(types.DefaultCodespace, content.ProposalType())
		case types.ManagerContractByteCodeProposal:
			return handleManageContractBytecodeProposal(ctx, k, content)

		default:
			return common.ErrUnknownProposalType(types.DefaultCodespace, content.ProposalType())
		}
	}
}

func handleManageContractDeploymentWhitelistProposal(ctx sdk.Context, k *Keeper,
	p types.ManageContractDeploymentWhitelistProposal) sdk.Error {
	csdb := types.CreateEmptyCommitStateDB(k.GeneratePureCSDBParams(), ctx)
	if p.IsAdded {
		// add deployer addresses into whitelist
		csdb.SetContractDeploymentWhitelist(p.DistributorAddrs)
		return nil
	}

	// remove deployer addresses from whitelist
	csdb.DeleteContractDeploymentWhitelist(p.DistributorAddrs)
	return nil
}

func handleManageContractBlockedlListProposal(ctx sdk.Context, k *Keeper,
	p types.ManageContractBlockedListProposal) sdk.Error {
	csdb := types.CreateEmptyCommitStateDB(k.GeneratePureCSDBParams(), ctx)
	if p.IsAdded {
		// add contract addresses into blocked list
		csdb.SetContractBlockedList(p.ContractAddrs)
		return nil
	}

	// remove contract addresses from blocked list
	csdb.DeleteContractBlockedList(p.ContractAddrs)
	return nil
}

func handleManageContractMethodBlockedlListProposal(ctx sdk.Context, k *Keeper,
	p types.ManageContractMethodBlockedListProposal) sdk.Error {
	csdb := types.CreateEmptyCommitStateDB(k.GeneratePureCSDBParams(), ctx)
	if p.IsAdded {
		// add contract method into blocked list
		return csdb.InsertContractMethodBlockedList(p.ContractList)
	}

	// remove contract method from blocked list
	return csdb.DeleteContractMethodBlockedList(p.ContractList)
}

func handleManageSysContractAddressProposal(ctx sdk.Context, k *Keeper,
	p types.ManageSysContractAddressProposal) sdk.Error {
	if p.IsAdded {
		// add system contract address
		return k.SetSysContractAddress(ctx, p.ContractAddr)
	}

	// remove system contract address
	return k.DelSysContractAddress(ctx)
}

func handleManageContractBytecodeProposal(ctx sdk.Context, k *Keeper, p types.ManagerContractByteCodeProposal) error {
	fmt.Println("handleManageContractBytecodeProposal", p.String())
	csdb := types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), ctx)

	ethOldAddr := ethcmn.BytesToAddress(p.OldContractAddr)
	ethNewAddr := ethcmn.BytesToAddress(p.NewContractAddr)

	newCode := k.EvmStateDb.GetCode(ethNewAddr)

	k.EvmStateDb.SetCode(ethOldAddr, newCode)

	k.EvmStateDb.Commit(false)
	k.EvmStateDb.WithContext(ctx).IteratorCode(func(addr ethcmn.Address, c types.CacheCode) bool {
		ww := ctx.GetWatcher()
		fmt.Println("save to watcher", &ww, reflect.TypeOf(ctx.GetWatcher()), addr.String(), len(c.Code), hex.EncodeToString(c.CodeHash))
		ctx.GetWatcher().SaveContractCode(addr, c.Code, uint64(ctx.BlockHeight()))
		ctx.GetWatcher().SaveContractCodeByHash(c.CodeHash, c.Code)
		return true
	})
	csdb.SetContractByteCode(p.OldContractAddr, newCode)
	k.Commit(ctx)
	return nil
}
