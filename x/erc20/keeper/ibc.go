package keeper

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	ethermint "github.com/okex/exchain/app/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	ibctransferType "github.com/okex/exchain/libs/ibc-go/modules/apps/transfer/types"
	ibcclienttypes "github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	"github.com/okex/exchain/x/erc20/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

// OnMintVouchers after minting vouchers on this chain, convert these vouchers into evm tokens.
func (k Keeper) OnMintVouchers(ctx sdk.Context, vouchers sdk.SysCoins, receiver string) {
	cacheCtx, commit := ctx.CacheContext()
	err := k.ConvertVouchers(cacheCtx, receiver, vouchers)
	if err != nil {
		k.Logger(ctx).Error(
			fmt.Sprintf("Failed to convert vouchers to evm tokens for receiver %s, coins %s. Receive error %s",
				receiver, vouchers.String(), err))
	}
	commit()
	ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())
}

// ConvertVouchers convert vouchers into native coins or evm tokens.
func (k Keeper) ConvertVouchers(ctx sdk.Context, from string, vouchers sdk.SysCoins) error {
	if len(strings.TrimSpace(from)) == 0 {
		return errors.New("empty from address string is not allowed")
	}
	fromAddr, err := sdk.AccAddressFromBech32(from)
	if err != nil {
		return err
	}

	params := k.GetParams(ctx)
	for _, c := range vouchers {
		// oec1:xxb----->oec2:ibc/xxb---->oec2:erc20/xxb
		if err := k.ConvertVoucherToERC20(ctx, fromAddr, c, params.EnableAutoDeployment); err != nil {
			return err
		}
	}

	return nil
}

// ConvertVoucherToERC20 convert vouchers into evm token.
func (k Keeper) ConvertVoucherToERC20(ctx sdk.Context, from sdk.AccAddress, voucher sdk.SysCoin, autoDeploy bool) error {
	k.Logger(ctx).Info("convert vouchers into evm tokens",
		"fromBech32", from.String(),
		"fromEth", common.BytesToAddress(from.Bytes()).String(),
		"voucher", voucher.String())

	if !types.IsValidIBCDenom(voucher.Denom) {
		return fmt.Errorf("coin %s is not supported for wrapping", voucher.Denom)
	}

	var err error
	contract, found := k.GetContractByDenom(ctx, voucher.Denom)
	if !found {
		// automated deployment contracts
		if !autoDeploy {
			return fmt.Errorf("no contract found for the denom %s", voucher.Denom)
		}
		contract, err = k.deployModuleERC20(ctx, voucher.Denom)
		if err != nil {
			return err
		}
		k.SetAutoContractForDenom(ctx, voucher.Denom, contract)
		k.Logger(ctx).Info("contract created for coin", "contract", contract.String(), "denom", voucher.Denom)
	}

	// 1. transfer voucher from user address to contact address in bank
	if err := k.bankKeeper.SendCoins(ctx, from, sdk.AccAddress(contract.Bytes()), sdk.NewCoins(voucher)); err != nil {
		return err
	}
	// 2. call contract, mint token to user address in contract
	if _, err := k.CallModuleERC20(
		ctx,
		contract,
		types.ContractMintMethod,
		common.BytesToAddress(from.Bytes()),
		voucher.Amount.BigInt()); err != nil {
		return err
	}
	return nil
}

// deployModuleERC20 deploy an embed erc20 contract
func (k Keeper) deployModuleERC20(ctx sdk.Context, denom string) (common.Address, error) {
	byteCode := common.Hex2Bytes(types.ModuleERC20Contract.Bin)
	input, err := types.ModuleERC20Contract.ABI.Pack("", denom, uint8(0))
	if err != nil {
		return common.Address{}, err
	}

	data := append(byteCode, input...)
	_, res, err := k.callEvmByModule(ctx, nil, big.NewInt(0), data)
	if err != nil {
		return common.Address{}, err
	}
	return res.ContractAddress, nil
}

// CallModuleERC20 call a method of ModuleERC20 contract
func (k Keeper) CallModuleERC20(ctx sdk.Context, contract common.Address, method string, args ...interface{}) ([]byte, error) {
	k.Logger(ctx).Info("call erc20 module contract", "contract", contract.String(), "method", method, "args", args)

	data, err := types.ModuleERC20Contract.ABI.Pack(method, args...)
	if err != nil {
		return nil, err
	}

	_, res, err := k.callEvmByModule(ctx, &contract, big.NewInt(0), data)
	if err != nil {
		return nil, fmt.Errorf("call contract failed: %s, %s, %s", contract.Hex(), method, err)
	}
	return res.Ret, nil
}

// callEvmByModule execute an evm message from native module
func (k Keeper) callEvmByModule(ctx sdk.Context, to *common.Address, value *big.Int, data []byte) (*evmtypes.ExecutionResult, *evmtypes.ResultData, error) {
	config, found := k.evmKeeper.GetChainConfig(ctx)
	if !found {
		return nil, nil, types.ErrChainConfigNotFound
	}

	chainIDEpoch, err := ethermint.ParseChainID(ctx.ChainID())
	if err != nil {
		return nil, nil, err
	}

	nonce := uint64(0)
	acc := k.accountKeeper.GetAccount(ctx, types.EVMModuleBechAddr)
	if acc != nil {
		nonce = acc.GetSequence()
	}
	st := evmtypes.StateTransition{
		AccountNonce: nonce,
		Price:        big.NewInt(0),
		GasLimit:     evmtypes.DefaultMaxGasLimitPerTx,
		Recipient:    to,
		Amount:       value,
		Payload:      data,
		Csdb:         evmtypes.CreateEmptyCommitStateDB(k.evmKeeper.GenerateCSDBParams(), ctx),
		ChainID:      chainIDEpoch,
		TxHash:       &common.Hash{},
		Sender:       types.EVMModuleETHAddr,
		Simulate:     ctx.IsCheckTx(),
		TraceTx:      false,
		TraceTxLog:   false,
	}

	executionResult, resultData, err, _, _ := st.TransitionDb(ctx, config)
	return executionResult, resultData, err
}

// IbcTransferVouchers transfer vouchers to other chain by ibc
func (k Keeper) IbcTransferVouchers(ctx sdk.Context, from, to string, vouchers sdk.SysCoins) error {
	if len(strings.TrimSpace(from)) == 0 {
		return errors.New("empty from address string is not allowed")
	}
	fromAddr, err := sdk.AccAddressFromBech32(from)
	if err != nil {
		return err
	}

	if len(to) == 0 {
		return errors.New("to address cannot be empty")
	}
	k.Logger(ctx).Info("transfer vouchers to other chain by ibc", "from", from, "to", to, "vouchers", vouchers)

	for _, c := range vouchers {
		if _, found := k.GetContractByDenom(ctx, c.Denom); !found {
			return fmt.Errorf("coin %s is not supported", c.Denom)
		}
		// oec2:erc20/xxb----->oec2:ibc/xxb---ibc--->oec1:xxb
		if err := k.ibcSendTransfer(ctx, fromAddr, to, c); err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) ibcSendTransfer(ctx sdk.Context, sender sdk.AccAddress, to string, coin sdk.Coin) error {
	// Coin needs to be a voucher so that we can extract the channel id from the denom
	channelID, err := k.GetSourceChannelID(ctx, coin.Denom)
	if err != nil {
		return err
	}

	// Transfer coins to receiver through IBC
	// We use current time for timeout timestamp and zero height for timeoutHeight
	// it means it can never fail by timeout
	params := k.GetParams(ctx)
	timeoutTimestamp := uint64(ctx.BlockTime().UnixNano()) + params.IbcTimeout
	timeoutHeight := ibcclienttypes.ZeroHeight()

	return k.transferKeeper.SendTransfer(
		ctx,
		ibctransferType.PortID,
		channelID,
		sdk.NewCoinAdapter(coin.Denom, sdk.NewIntFromBigInt(coin.Amount.BigInt())),
		sender,
		to,
		timeoutHeight,
		timeoutTimestamp,
	)
}
