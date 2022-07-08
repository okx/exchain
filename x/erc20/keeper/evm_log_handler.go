package keeper

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/erc20/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

var (
	_ evmtypes.EvmLogHandler = SendToIbcEventHandler{}
)

const (
	SendToIbcEventName         = "__OKCSendToIbc"
	SendNative20ToIbcEventName = "__OKCSendNative20ToIbc"
	SendToWasmEventName        = "__OKCSendToWasm"
)

// SendToIbcEvent represent the signature of
// `event __OKCSendToIbc(string recipient, uint256 amount)`
var SendToIbcEvent abi.Event

// SendNative20ToIbcEvent represent the signature of
// `event __OKCSendNative20ToIbc(string recipient, uint256 amount, string portID, string channelID)`
var SendNative20ToIbcEvent abi.Event

func init() {
	addressType, _ := abi.NewType("address", "", nil)
	uint256Type, _ := abi.NewType("uint256", "", nil)
	stringType, _ := abi.NewType("string", "", nil)

	SendToIbcEvent = abi.NewEvent(
		SendToIbcEventName,
		SendToIbcEventName,
		false,
		abi.Arguments{abi.Argument{
			Name:    "sender",
			Type:    addressType,
			Indexed: false,
		}, abi.Argument{
			Name:    "recipient",
			Type:    stringType,
			Indexed: false,
		}, abi.Argument{
			Name:    "amount",
			Type:    uint256Type,
			Indexed: false,
		}},
	)

	SendNative20ToIbcEvent = abi.NewEvent(
		SendNative20ToIbcEventName,
		SendNative20ToIbcEventName,
		false,
		abi.Arguments{abi.Argument{
			Name:    "sender",
			Type:    addressType,
			Indexed: false,
		}, abi.Argument{
			Name:    "recipient",
			Type:    stringType,
			Indexed: false,
		}, abi.Argument{
			Name:    "amount",
			Type:    uint256Type,
			Indexed: false,
		}, abi.Argument{
			Name:    "portID",
			Type:    stringType,
			Indexed: false,
		}, abi.Argument{
			Name:    "channelID",
			Type:    stringType,
			Indexed: false,
		}},
	)
}

type SendToIbcEventHandler struct {
	Keeper
}

func NewSendToIbcEventHandler(k Keeper) *SendToIbcEventHandler {
	return &SendToIbcEventHandler{k}
}

// EventID Return the id of the log signature it handles
func (h SendToIbcEventHandler) EventID() common.Hash {
	return SendToIbcEvent.ID
}

// Handle Process the log
func (h SendToIbcEventHandler) Handle(ctx sdk.Context, contract common.Address, data []byte) error {
	h.Logger(ctx).Info("trigger evm event", "event", SendToIbcEvent.Name, "contract", contract)
	// first confirm that the contract address and denom are registered,
	// to avoid unpacking any contract '__OKCSendToIbc' event, which consumes performance
	denom, found := h.Keeper.GetDenomByContract(ctx, contract)
	if !found {
		return fmt.Errorf("contract %s is not connected to native token", contract)
	}
	if !types.IsValidIBCDenom(denom) {
		return fmt.Errorf("the native token associated with the contract %s is not an ibc voucher", contract)
	}

	unpacked, err := SendToIbcEvent.Inputs.Unpack(data)
	if err != nil {
		// log and ignore
		h.Keeper.Logger(ctx).Error("log signature matches but failed to decode", "error", err)
		return nil
	}

	contractAddr := sdk.AccAddress(contract.Bytes())
	sender := sdk.AccAddress(unpacked[0].(common.Address).Bytes())
	recipient := unpacked[1].(string)
	amount := sdk.NewIntFromBigInt(unpacked[2].(*big.Int))
	amountDec := sdk.NewDecFromIntWithPrec(amount, sdk.Precision)
	vouchers := sdk.NewCoins(sdk.NewCoin(denom, amountDec))

	// 1. transfer IBC coin to user so that he will be the refunded address if transfer fails
	if err = h.bankKeeper.SendCoins(ctx, contractAddr, sender, vouchers); err != nil {
		return err
	}
	// 2. Initiate IBC transfer from sender account
	if err = h.Keeper.IbcTransferVouchers(ctx, sender.String(), recipient, vouchers); err != nil {
		return err
	}
	return nil
}

type SendNative20ToIbcEventHandler struct {
	Keeper
}

func NewSendNative20ToIbcEventHandler(k Keeper) *SendNative20ToIbcEventHandler {
	return &SendNative20ToIbcEventHandler{k}
}

// EventID Return the id of the log signature it handles
func (h SendNative20ToIbcEventHandler) EventID() common.Hash {
	return SendNative20ToIbcEvent.ID
}

// Handle Process the log
func (h SendNative20ToIbcEventHandler) Handle(ctx sdk.Context, contract common.Address, data []byte) error {
	h.Logger(ctx).Info("trigger evm event", "event", SendNative20ToIbcEvent.Name, "contract", contract)
	// first confirm that the contract address and denom are registered,
	// to avoid unpacking any contract '__OKCSendNative20ToIbc' event, which consumes performance
	denom, found := h.Keeper.GetDenomByContract(ctx, contract)
	if !found {
		return fmt.Errorf("contract %s is not connected to native token", contract)
	}
	if err := sdk.ValidateDenom(denom); err != nil {
		return fmt.Errorf("the native token associated with the contract %s is not an valid token", contract)
	}

	unpacked, err := SendNative20ToIbcEvent.Inputs.Unpack(data)
	if err != nil {
		// log and ignore
		h.Keeper.Logger(ctx).Error("log signature matches but failed to decode", "error", err)
		return nil
	}

	//contractAddr := sdk.AccAddress(contract.Bytes())
	sender := sdk.AccAddress(unpacked[0].(common.Address).Bytes())
	recipient := unpacked[1].(string)
	amount := sdk.NewIntFromBigInt(unpacked[2].(*big.Int))
	portID := unpacked[3].(string)
	channelID := unpacked[4].(string)

	amountDec := sdk.NewDecFromIntWithPrec(amount, sdk.Precision)
	native20s := sdk.NewCoins(sdk.NewCoin(denom, amountDec))

	// 1. mint new tokens to user so that he will be the refunded address if transfer fails
	if err = h.supplyKeeper.MintCoins(ctx, types.ModuleName, native20s); err != nil {
		return err
	}
	if err = h.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, native20s); err != nil {
		return err
	}

	// 2. Initiate IBC transfer from sender account
	if err = h.Keeper.IbcTransferNative20(ctx, sender.String(), recipient, native20s, portID, channelID); err != nil {
		return err
	}
	return nil
}
