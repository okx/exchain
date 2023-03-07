package evm

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/x/common"
	"github.com/okx/okbchain/x/evm/types"
	"github.com/okx/okbchain/x/evm/watcher"
	govTypes "github.com/okx/okbchain/x/gov/types"
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
			return handleManageSysContractAddressProposal(ctx, k, content)
		case types.ManageContractByteCodeProposal:
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

func handleManageContractBytecodeProposal(ctx sdk.Context, k *Keeper, p types.ManageContractByteCodeProposal) error {
	csdb := types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), ctx)
	return csdb.UpdateContractBytecode(ctx, p)
}
