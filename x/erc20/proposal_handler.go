package erc20

import (
	ethcmm "github.com/ethereum/go-ethereum/common"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
	"github.com/okx/okbchain/x/common"
	"github.com/okx/okbchain/x/erc20/types"
	govTypes "github.com/okx/okbchain/x/gov/types"
)

// NewProposalHandler handles "gov" type message in "erc20"
func NewProposalHandler(k *Keeper) govTypes.Handler {
	return func(ctx sdk.Context, proposal *govTypes.Proposal) (err sdk.Error) {
		switch content := proposal.Content.(type) {
		case types.TokenMappingProposal:
			return handleTokenMappingProposal(ctx, k, content)
		case types.ProxyContractRedirectProposal:
			return handleProxyContractRedirectProposal(ctx, k, content)
		case types.ContractTemplateProposal:
			return handleContractTemplateProposal(ctx, k, content)
		default:
			return common.ErrUnknownProposalType(types.DefaultCodespace, content.ProposalType())
		}
	}
}

func handleTokenMappingProposal(ctx sdk.Context, k *Keeper, p types.TokenMappingProposal) sdk.Error {
	if p.Denom == sdk.DefaultBondDenom || p.Denom == sdk.DefaultIbcWei {
		return govTypes.ErrInvalidProposalContent("invalid denom, not support okt or wei denom")
	}

	if len(p.Contract) == 0 {
		// delete existing mapping
		k.DeleteContractForDenom(ctx, p.Denom)
	} else {
		// update the mapping
		contract := ethcmm.HexToAddress(p.Contract)
		// contract must already be deployed, to avoid empty contract
		contractAccount, _ := k.GetEthAccount(ctx, contract)
		if contractAccount == nil || !contractAccount.IsContract() {
			return sdkerrors.Wrapf(types.ErrNoContractDeployed, "no contract code found at address %s", p.Contract)
		}

		if err := k.SetContractForDenom(ctx, p.Denom, contract); err != nil {
			return err
		}
	}
	return nil
}

func handleProxyContractRedirectProposal(ctx sdk.Context, k *Keeper, p types.ProxyContractRedirectProposal) sdk.Error {
	address := ethcmm.HexToAddress(p.Addr)

	return k.ProxyContractRedirect(ctx, p.Denom, p.Tp, address)
}

func handleContractTemplateProposal(ctx sdk.Context, k *Keeper, p types.ContractTemplateProposal) sdk.Error {
	return k.SetTemplateContract(ctx, p.ContractType, p.Contract)
}
