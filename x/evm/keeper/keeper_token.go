package keeper

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	ethermint "github.com/okex/exchain/app/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	ibctransferType "github.com/okex/exchain/libs/ibc-go/modules/application/transfer/types"
	ibcclienttypes "github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	"github.com/okex/exchain/x/evm/types"
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
	fromAddr, err := sdk.AccAddressFromBech32(from)
	if err != nil {
		return err
	}

	params := k.GetParams(ctx)
	for _, c := range vouchers {
		switch c.Denom {
		case params.IbcDenom:
			// oec1:okt----->oec2:ibc/okt---->oec2:okt
			if err := k.ConvertVoucherToNative(ctx, fromAddr, c); err != nil {
				return err
			}
		default:
			// oec1:xxb----->oec2:ibc/xxb---->oec2:erc20/xxb
			// TODO use autoDeploy boolean in params
			if err := k.ConvertVoucherToERC20(ctx, fromAddr, c, true); err != nil {
				return err
			}
		}
	}
	return nil
}

// ConvertVoucherToNative convert vouchers into native coins.
func (k Keeper) ConvertVoucherToNative(ctx sdk.Context, from sdk.AccAddress, voucher sdk.SysCoin) error {
	// TODO
	return nil
}

// ConvertVoucherToERC20 convert vouchers into evm tokens.
func (k Keeper) ConvertVoucherToERC20(ctx sdk.Context, from sdk.AccAddress, voucher sdk.SysCoin, autoDeploy bool) error {
	err := ibctransferType.ValidateIBCDenom(voucher.Denom)
	if err != nil {
		return ibctransferType.ErrInvalidDenomForTransfer
	}

	contract, found := k.getContractByDenom(ctx, voucher.Denom)
	if !found {
		// automated deployment contracts
		if !autoDeploy {
			return fmt.Errorf("no contract found for the denom %s", voucher.Denom)
		}
		contract, err = k.deployModuleERC20(ctx, voucher.Denom)
		if err != nil {
			return err
		}
		k.setAutoContractForDenom(ctx, voucher.Denom, contract)
		k.Logger(ctx).Info(fmt.Sprintf("contract address %s created for coin denom %s", contract.String(), voucher.Denom))
	}
	// 1. transfer voucher from user address to contact address in bank
	if err := k.bankKeeper.SendCoins(ctx, from, sdk.AccAddress(contract.Bytes()), sdk.NewCoins(voucher)); err != nil {
		return err
	}
	// 2. call contract, mint token to user address in contract
	if _, err := k.callModuleERC20(
		ctx,
		contract,
		"mint_by_oec_module",
		common.BytesToAddress(from.Bytes()),
		voucher.Amount.BigInt()); err != nil {
		return err
	}
	return nil
}

// deployModuleERC20 deploy an embed erc20 contract
func (k Keeper) deployModuleERC20(ctx sdk.Context, denom string) (common.Address, error) {
	// TODO
	//k.callEvmByModule(ctx,nil,big.NewInt(0),data)
	return common.Address{}, nil
}

// callModuleERC20 call a method of ModuleERC20 contract
func (k Keeper) callModuleERC20(ctx sdk.Context, contract common.Address, method string, args ...interface{}) ([]byte, error) {
	// TODO
	return nil, nil
}

// callEvmByModule execute an evm message from native module
func (k Keeper) callEvmByModule(ctx sdk.Context, to *common.Address, value *big.Int, data []byte) (*types.ExecutionResult, *types.ResultData, error) {
	config, found := k.GetChainConfig(ctx)
	if !found {
		return nil, nil, types.ErrChainConfigNotFound
	}

	chainIDEpoch, err := ethermint.ParseChainID(ctx.ChainID())
	if err != nil {
		return nil, nil, err
	}

	acc := k.accountKeeper.GetAccount(ctx, types.EVMModuleBechAddr)

	st := types.StateTransition{
		AccountNonce: acc.GetSequence(),
		Price:        big.NewInt(0),
		GasLimit:     types.DefaultMaxGasLimitPerTx,
		Recipient:    to,
		Amount:       value,
		Payload:      data,
		Csdb:         types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), ctx),
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
	fromAddr, err := sdk.AccAddressFromBech32(from)
	if err != nil {
		return err
	}

	if len(to) == 0 {
		return errors.New("to address cannot be empty")
	}

	//params := k.GetParams(ctx)
	for _, c := range vouchers {
		switch c.Denom {
		case sdk.DefaultBondDenom:
			// oec2:okt----->oec2:ibc/okt---ibc--->oec1:okt
			if err := k.ibcSendEvmDenom(ctx, fromAddr, to, c); err != nil {
				return err
			}
		default:
			if _, found := k.getContractByDenom(ctx, c.Denom); !found {
				return fmt.Errorf("coin %s id not support", c.Denom)
			}
			// oec2:erc20/xxb----->oec2:ibc/xxb---ibc--->oec1:xxb
			if err := k.ibcSendTransfer(ctx, fromAddr, to, c); err != nil {
				return err
			}
		}
	}

	return nil
}

func (k Keeper) ibcSendEvmDenom(ctx sdk.Context, sender sdk.AccAddress, to string, coin sdk.Coin) error {
	// TODO
	return nil
}

func (k Keeper) ibcSendTransfer(ctx sdk.Context, sender sdk.AccAddress, to string, coin sdk.Coin) error {
	// Coin needs to be a voucher so that we can extract the channel id from the denom
	channelID, err := k.GetSourceChannelID(ctx, coin.Denom)
	if err != nil {
		return nil
	}

	// Transfer coins to receiver through IBC
	// We use current time for timeout timestamp and zero height for timeoutHeight
	// it means it can never fail by timeout
	timeoutTimestamp := uint64(ctx.BlockTime().UnixNano())
	timeoutHeight := ibcclienttypes.ZeroHeight()
	return k.transferKeeper.SendTransfer(
		ctx,
		ibctransferType.PortID,
		channelID,
		coin,
		sender,
		to,
		timeoutHeight,
		timeoutTimestamp,
	)
}

// GetSourceChannelID returns the channel id for an ibc voucher
// The voucher has for format ibc/hash(path)
func (k Keeper) GetSourceChannelID(ctx sdk.Context, ibcVoucherDenom string) (channelID string, err error) {
	path, err := k.transferKeeper.DenomPathFromHash(ctx, ibcVoucherDenom)
	if err != nil {
		return "", err
	}

	// the path has for format port/channelId
	return strings.Split(path, "/")[1], nil
}

// DeleteExternalContractForDenom delete the external contract mapping for native denom,
// returns false if mapping not exists.
func (k Keeper) DeleteExternalContractForDenom(ctx sdk.Context, denom string) bool {
	store := ctx.KVStore(k.storeKey)
	existingContract, found := k.getExternalContractByDenom(ctx, denom)
	if !found {
		return false
	}
	store.Delete(types.ContractToDenomKey(existingContract.Bytes()))
	store.Delete(types.DenomToExternalContractKey(denom))
	return true
}

// SetExternalContractForDenom set the external contract for native denom,
// 1. if any existing for denom, replace the old one.
// 2. if any existing for contract, return error.
func (k Keeper) SetExternalContractForDenom(ctx sdk.Context, denom string, contract common.Address) error {
	// check the contract is not registered already
	_, found := k.getDenomByContract(ctx, contract)
	if found {
		return types.ErrRegisteredContract(contract.String())
	}

	store := ctx.KVStore(k.storeKey)
	existingContract, found := k.getExternalContractByDenom(ctx, denom)
	if found {
		// delete existing mapping
		store.Delete(types.ContractToDenomKey(existingContract.Bytes()))
	}
	store.Set(types.DenomToExternalContractKey(denom), contract.Bytes())
	store.Set(types.ContractToDenomKey(contract.Bytes()), []byte(denom))
	return nil
}

func (k Keeper) setAutoContractForDenom(ctx sdk.Context, denom string, contract common.Address) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DenomToAutoContractKey(denom), contract.Bytes())
	store.Set(types.ContractToDenomKey(contract.Bytes()), []byte(denom))
}

func (k Keeper) getDenomByContract(ctx sdk.Context, contract common.Address) (denom string, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ContractToDenomKey(contract.Bytes()))
	if len(bz) == 0 {
		return "", false
	}
	return string(bz), true
}

// IterateMapping iterates over all the stored mapping and performs a callback function
func (k Keeper) IterateMapping(ctx sdk.Context, cb func(denom, contract string) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixContractToDenom)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		denom := string(iterator.Value())
		conotract := common.BytesToAddress(iterator.Key()).String()

		if cb(denom, conotract) {
			break
		}
	}
}

func (k Keeper) getExternalContractByDenom(ctx sdk.Context, denom string) (contract common.Address, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DenomToExternalContractKey(denom))
	if len(bz) == 0 {
		return common.Address{}, false
	}
	return common.BytesToAddress(bz), true
}

func (k Keeper) getAutoContractByDenom(ctx sdk.Context, denom string) (contract common.Address, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DenomToAutoContractKey(denom))
	if len(bz) == 0 {
		return common.Address{}, false
	}
	return common.BytesToAddress(bz), true
}

func (k Keeper) getContractByDenom(ctx sdk.Context, denom string) (contract common.Address, found bool) {
	contract, found = k.getExternalContractByDenom(ctx, denom)
	if !found {
		contract, found = k.getAutoContractByDenom(ctx, denom)
	}
	return
}
