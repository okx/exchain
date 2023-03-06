package feesplit

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
	"github.com/okx/okbchain/x/feesplit/keeper"
	"github.com/okx/okbchain/x/feesplit/types"
)

// NewHandler defines the fees module handler instance
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx.SetEventManager(sdk.NewEventManager())

		params := k.GetParams(ctx)
		if !params.EnableFeeSplit {
			return nil, types.ErrFeeSplitDisabled
		}

		switch msg := msg.(type) {
		case types.MsgRegisterFeeSplit:
			return handleMsgRegisterFeeSplit(ctx, msg, k, params)
		case types.MsgUpdateFeeSplit:
			return handleMsgUpdateFeeSplit(ctx, msg, k)
		case types.MsgCancelFeeSplit:
			return handleMsgCancelFeeSplit(ctx, msg, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

// handleMsgRegisterFeeSplit registers a contract to receive transaction fees
func handleMsgRegisterFeeSplit(
	ctx sdk.Context,
	msg types.MsgRegisterFeeSplit,
	k keeper.Keeper,
	params types.Params,
) (*sdk.Result, error) {
	contract := common.HexToAddress(msg.ContractAddress)
	if k.IsFeeSplitRegistered(ctx, contract) {
		return nil, sdkerrors.Wrapf(
			types.ErrFeeSplitAlreadyRegistered,
			"contract is already registered %s", contract,
		)
	}

	deployer := sdk.MustAccAddressFromBech32(msg.DeployerAddress)
	deployerAccount, isExist := k.GetEthAccount(ctx, common.BytesToAddress(deployer))
	if !isExist {
		return nil, sdkerrors.Wrapf(
			types.ErrFeeAccountNotFound,
			"deployer account not found %s", msg.DeployerAddress,
		)
	}

	if deployerAccount != nil && deployerAccount.IsContract() {
		return nil, sdkerrors.Wrapf(
			types.ErrFeeSplitDeployerIsNotEOA,
			"deployer cannot be a contract %s", msg.DeployerAddress,
		)
	}

	// contract must already be deployed, to avoid spam registrations
	contractAccount, _ := k.GetEthAccount(ctx, contract)
	if contractAccount == nil || !contractAccount.IsContract() {
		return nil, sdkerrors.Wrapf(
			types.ErrFeeSplitNoContractDeployed,
			"no contract code found at address %s", msg.ContractAddress,
		)
	}

	if msg.WithdrawerAddress == "" {
		msg.WithdrawerAddress = msg.DeployerAddress
	}
	withdrawer := sdk.MustAccAddressFromBech32(msg.WithdrawerAddress)

	derivedContract := common.BytesToAddress(deployer)

	// the contract can be directly deployed by an EOA or created through one
	// or more factory contracts. If it was deployed by an EOA account, then
	// msg.Nonces contains the EOA nonce for the deployment transaction.
	// If it was deployed by one or more factories, msg.Nonces contains the EOA
	// nonce for the origin factory contract, then the nonce of the factory
	// for the creation of the next factory/contract.
	for _, nonce := range msg.Nonces {
		ctx.GasMeter().ConsumeGas(
			params.AddrDerivationCostCreate,
			"fee split registration: address derivation CREATE opcode",
		)

		derivedContract = crypto.CreateAddress(derivedContract, nonce)
	}

	if contract != derivedContract {
		return nil, sdkerrors.Wrapf(
			types.ErrDerivedNotMatched,
			"not contract deployer or wrong nonce: expected %s instead of %s",
			derivedContract, msg.ContractAddress,
		)
	}

	// prevent storing the same address for deployer and withdrawer
	feeSplit := types.NewFeeSplit(contract, deployer, withdrawer)
	k.SetFeeSplit(ctx, feeSplit)
	k.SetDeployerMap(ctx, deployer, contract)
	k.SetWithdrawerMap(ctx, withdrawer, contract)

	k.Logger(ctx).Debug(
		"registering contract for transaction fees",
		"contract", msg.ContractAddress, "deployer", msg.DeployerAddress,
		"withdraw", msg.WithdrawerAddress,
	)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeRegisterFeeSplit,
				sdk.NewAttribute(sdk.AttributeKeySender, msg.DeployerAddress),
				sdk.NewAttribute(types.AttributeKeyContract, msg.ContractAddress),
				sdk.NewAttribute(types.AttributeKeyWithdrawerAddress, msg.WithdrawerAddress),
			),
		},
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgUpdateFeeSplit updates the withdraw address of a given FeeSplit. If the given
// withdraw address is empty or the same as the deployer address, the withdraw
// address is removed.
func handleMsgUpdateFeeSplit(
	ctx sdk.Context,
	msg types.MsgUpdateFeeSplit,
	k keeper.Keeper,
) (*sdk.Result, error) {
	contract := common.HexToAddress(msg.ContractAddress)
	feeSplit, found := k.GetFeeSplit(ctx, contract)
	if !found {
		return nil, sdkerrors.Wrapf(
			types.ErrFeeSplitContractNotRegistered,
			"contract %s is not registered", msg.ContractAddress,
		)
	}

	// error if the msg deployer address is not the same as the fee's deployer
	if !sdk.MustAccAddressFromBech32(msg.DeployerAddress).Equals(feeSplit.DeployerAddress) {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrUnauthorized,
			"%s is not the contract deployer", msg.DeployerAddress,
		)
	}

	var withdrawer sdk.AccAddress
	withdrawer = sdk.MustAccAddressFromBech32(msg.WithdrawerAddress)

	// fee split with the given withdraw address is already registered
	if withdrawer.Equals(feeSplit.WithdrawerAddress) {
		return nil, sdkerrors.Wrapf(
			types.ErrFeeSplitAlreadyRegistered,
			"fee split with withdraw address %s", msg.WithdrawerAddress,
		)
	}

	k.DeleteWithdrawerMap(ctx, feeSplit.WithdrawerAddress, contract)
	k.SetWithdrawerMap(ctx, withdrawer, contract)
	// update fee split
	feeSplit.WithdrawerAddress = withdrawer
	k.SetFeeSplit(ctx, feeSplit)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeUpdateFeeSplit,
				sdk.NewAttribute(types.AttributeKeyContract, msg.ContractAddress),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.DeployerAddress),
				sdk.NewAttribute(types.AttributeKeyWithdrawerAddress, msg.WithdrawerAddress),
			),
		},
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgCancelFeeSplit deletes the FeeSplit for a given contract
func handleMsgCancelFeeSplit(
	ctx sdk.Context,
	msg types.MsgCancelFeeSplit,
	k keeper.Keeper,
) (*sdk.Result, error) {
	contract := common.HexToAddress(msg.ContractAddress)
	fee, found := k.GetFeeSplit(ctx, contract)
	if !found {
		return nil, sdkerrors.Wrapf(
			types.ErrFeeSplitContractNotRegistered,
			"contract %s is not registered", msg.ContractAddress,
		)
	}

	if !sdk.MustAccAddressFromBech32(msg.DeployerAddress).Equals(fee.DeployerAddress) {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrUnauthorized,
			"%s is not the contract deployer", msg.DeployerAddress,
		)
	}

	k.DeleteFeeSplit(ctx, fee)
	k.DeleteDeployerMap(ctx, fee.DeployerAddress, contract)
	k.DeleteWithdrawerMap(ctx, fee.WithdrawerAddress, contract)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeCancelFeeSplit,
				sdk.NewAttribute(sdk.AttributeKeySender, msg.DeployerAddress),
				sdk.NewAttribute(types.AttributeKeyContract, msg.ContractAddress),
			),
		},
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
