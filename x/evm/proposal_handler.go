package evm

import (
	ethcmm "github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/evm/types"
	govTypes "github.com/okex/exchain/x/gov/types"
)

// NewManageContractDeploymentWhitelistProposalHandler handles "gov" type message in "evm"
func NewManageContractDeploymentWhitelistProposalHandler(k *Keeper) govTypes.Handler {
	return func(ctx sdk.Context, proposal *govTypes.Proposal) (err sdk.Error) {
		switch content := proposal.Content.(type) {
		case types.ManageContractDeploymentWhitelistProposal:
			return handleManageContractDeploymentWhitelistProposal(ctx, k, content)
		case types.ManageContractBlockedListProposal:
			return handleManageContractBlockedlListProposal(ctx, k, content)
		case types.ManageContractMethodBlockedListProposal:
			return handleManageContractMethodBlockedlListProposal(ctx, k, content)
		case types.TokenMappingProposal:
			return handleTokenMappingProposal(ctx, k, content)
		default:
			return common.ErrUnknownProposalType(types.DefaultCodespace, content.ProposalType())
		}
	}
}

func handleManageContractDeploymentWhitelistProposal(ctx sdk.Context, k *Keeper, p types.ManageContractDeploymentWhitelistProposal) sdk.Error {
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

func handleManageContractBlockedlListProposal(ctx sdk.Context, k *Keeper, p types.ManageContractBlockedListProposal) sdk.Error {
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

func handleManageContractMethodBlockedlListProposal(ctx sdk.Context, k *Keeper, p types.ManageContractMethodBlockedListProposal) sdk.Error {
	csdb := types.CreateEmptyCommitStateDB(k.GeneratePureCSDBParams(), ctx)
	if p.IsAdded {
		// add contract method into blocked list
		return csdb.InsertContractMethodBlockedList(p.ContractList)
	}

	// remove contract method from blocked list
	return csdb.DeleteContractMethodBlockedList(p.ContractList)
}

func handleTokenMappingProposal(ctx sdk.Context, k *Keeper, p types.TokenMappingProposal) sdk.Error {
	if len(p.Contract) == 0 {
		// delete existing mapping
		k.DeleteExternalContractForDenom(ctx, p.Denom)
	} else {
		// update the mapping
		contract := ethcmm.HexToAddress(p.Contract)
		if err := k.SetExternalContractForDenom(ctx, p.Denom, contract); err != nil {
			return err
		}
	}
	return nil
}
