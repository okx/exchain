package keeper_test

/// This file is a test driver for model-based tests generated from the TLA+ model of token transfer
/// Written by Andrey Kuprianov within the scope of IBC Audit performed by Informal Systems.
/// In case of any questions please don't hesitate to contact andrey@informal.systems.

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/tendermint/crypto"
	"strings"

	// "github.com/tendermint/tendermint/crypto"

	// sdk "github.com/cosmos/cosmos-sdk/types"
	// sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/transfer/types"
	ibctesting "github.com/okex/exchain/libs/ibc-go/testing"
)

type TlaBalance struct {
	Address []string `json:"address"`
	Denom   []string `json:"denom"`
	Amount  int64    `json:"amount"`
}

type TlaFungibleTokenPacketData struct {
	Sender   string   `json:"sender"`
	Receiver string   `json:"receiver"`
	Amount   string   `json:"amount"`
	Denom    []string `json:"denom"`
}

type TlaFungibleTokenPacket struct {
	SourceChannel string                     `json:"sourceChannel"`
	SourcePort    string                     `json:"sourcePort"`
	DestChannel   string                     `json:"destChannel"`
	DestPort      string                     `json:"destPort"`
	Data          TlaFungibleTokenPacketData `json:"data"`
}

type TlaOnRecvPacketTestCase = struct {
	// The required subset of bank balances
	BankBefore []TlaBalance `json:"bankBefore"`
	// The packet to process
	Packet TlaFungibleTokenPacket `json:"packet"`
	// The handler to call
	Handler string `json:"handler"`
	// The expected changes in the bank
	BankAfter []TlaBalance `json:"bankAfter"`
	// Whether OnRecvPacket should fail or not
	Error bool `json:"error"`
}

type FungibleTokenPacket struct {
	SourceChannel string
	SourcePort    string
	DestChannel   string
	DestPort      string
	Data          types.FungibleTokenPacketData
}

type OnRecvPacketTestCase = struct {
	description string
	// The required subset of bank balances
	bankBefore []Balance
	// The packet to process
	packet FungibleTokenPacket
	// The handler to call
	handler string
	// The expected bank state after processing (wrt. bankBefore)
	bankAfter []Balance
	// Whether OnRecvPacket should pass or fail
	pass bool
}

type OwnedCoin struct {
	Address string
	Denom   string
}

type Balance struct {
	Id      string
	Address string
	Denom   string
	Amount  sdk.Int
}

func AddressFromString(address string) string {
	return sdk.AccAddress(crypto.AddressHash([]byte(address))).String()
}

func AddressFromTla(addr []string) string {
	if len(addr) != 3 {
		panic("failed to convert from TLA+ address: wrong number of address components")
	}
	s := ""
	if len(addr[0]) == 0 && len(addr[1]) == 0 {
		// simple address: id
		s = addr[2]
	} else if len(addr[2]) == 0 {
		// escrow address: ics20-1\x00port/channel
		s = fmt.Sprintf("%s\x00%s/%s", types.Version, addr[0], addr[1])
	} else {
		panic("failed to convert from TLA+ address: neither simple nor escrow address")
	}
	return s
}

func DenomFromTla(denom []string) string {
	var i int
	for i = 0; i+1 < len(denom) && len(denom[i]) == 0 && len(denom[i+1]) == 0; i += 2 {
		// skip empty prefixes
	}
	return strings.Join(denom[i:], "/")
}

func BalanceFromTla(balance TlaBalance) Balance {
	return Balance{
		Id:      AddressFromTla(balance.Address),
		Address: AddressFromString(AddressFromTla(balance.Address)),
		Denom:   DenomFromTla(balance.Denom),
		Amount:  sdk.NewInt(balance.Amount),
	}
}

func BalancesFromTla(tla []TlaBalance) []Balance {
	balances := make([]Balance, 0)
	for _, b := range tla {
		balances = append(balances, BalanceFromTla(b))
	}
	return balances
}

func FungibleTokenPacketFromTla(packet TlaFungibleTokenPacket) FungibleTokenPacket {
	return FungibleTokenPacket{
		SourceChannel: packet.SourceChannel,
		SourcePort:    packet.SourcePort,
		DestChannel:   packet.DestChannel,
		DestPort:      packet.DestPort,
		Data: types.NewFungibleTokenPacketData(
			DenomFromTla(packet.Data.Denom),
			packet.Data.Amount,
			AddressFromString(packet.Data.Sender),
			AddressFromString(packet.Data.Receiver)),
	}
}

func OnRecvPacketTestCaseFromTla(tc TlaOnRecvPacketTestCase) OnRecvPacketTestCase {
	return OnRecvPacketTestCase{
		description: "auto-generated",
		bankBefore:  BalancesFromTla(tc.BankBefore),
		packet:      FungibleTokenPacketFromTla(tc.Packet),
		handler:     tc.Handler,
		bankAfter:   BalancesFromTla(tc.BankAfter), // TODO different semantics
		pass:        !tc.Error,
	}
}

var addressMap = make(map[string]string)

type Bank struct {
	balances map[OwnedCoin]sdk.Int
}

// Make an empty bank
func MakeBank() Bank {
	return Bank{balances: make(map[OwnedCoin]sdk.Int)}
}

// Subtract other bank from this bank
func (bank *Bank) Sub(other *Bank) Bank {
	diff := MakeBank()
	for coin, amount := range bank.balances {
		otherAmount, exists := other.balances[coin]
		if exists {
			diff.balances[coin] = amount.Sub(otherAmount)
		} else {
			diff.balances[coin] = amount
		}
	}
	for coin, amount := range other.balances {
		if _, exists := bank.balances[coin]; !exists {
			diff.balances[coin] = amount.Neg()
		}
	}
	return diff
}

// Set specific bank balance
func (bank *Bank) SetBalance(address string, denom string, amount sdk.Int) {
	bank.balances[OwnedCoin{address, denom}] = amount
}

// Set several balances at once
func (bank *Bank) SetBalances(balances []Balance) {
	for _, balance := range balances {
		bank.balances[OwnedCoin{balance.Address, balance.Denom}] = balance.Amount
		addressMap[balance.Address] = balance.Id
	}
}

func NullCoin() OwnedCoin {
	return OwnedCoin{
		Address: AddressFromString(""),
		Denom:   "",
	}
}

// Set several balances at once
func BankFromBalances(balances []Balance) Bank {
	bank := MakeBank()
	for _, balance := range balances {
		coin := OwnedCoin{balance.Address, balance.Denom}
		if coin != NullCoin() { // ignore null coin
			bank.balances[coin] = balance.Amount
			addressMap[balance.Address] = balance.Id
		}
	}
	return bank
}

// String representation of all bank balances
func (bank *Bank) String() string {
	str := ""
	for coin, amount := range bank.balances {
		str += coin.Address
		if addressMap[coin.Address] != "" {
			str += "(" + addressMap[coin.Address] + ")"
		}
		str += " : " + coin.Denom + " = " + amount.String() + "\n"
	}
	return str
}

// String representation of non-zero bank balances
func (bank *Bank) NonZeroString() string {
	str := ""
	for coin, amount := range bank.balances {
		if !amount.IsZero() {
			str += coin.Address + " : " + coin.Denom + " = " + amount.String() + "\n"
		}
	}
	return str
}

// Construct a bank out of the chain bank
func BankOfChain(chain ibctesting.TestChainI) Bank {
	bank := MakeBank()
	// todo how to Iterate all balance
	//	chain.GetSimApp().BankKeeper.IterateAllBalances(chain.GetContext(), func(address sdk.AccAddress, coin sdk.Coin) (stop bool) {
	//		fullDenom := coin.Denom
	//		if strings.HasPrefix(coin.Denom, "ibc/") {
	//			fullDenom, _ = chain.GetSimApp().TransferKeeper.DenomPathFromHash(chain.GetContext(), coin.Denom)
	//		}
	//		bank.SetBalance(address.String(), fullDenom, coin.Amount)
	//		return false
	//	})
	return bank
}

// Check that the state of the bank is the bankBefore + expectedBankChange
func (suite *KeeperTestSuite) CheckBankBalances(chain ibctesting.TestChainI, bankBefore *Bank, expectedBankChange *Bank) error {
	bankAfter := BankOfChain(chain)
	bankChange := bankAfter.Sub(bankBefore)
	diff := bankChange.Sub(expectedBankChange)
	NonZeroString := diff.NonZeroString()
	if len(NonZeroString) != 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "Unexpected changes in the bank: \n"+NonZeroString)
	}
	return nil
}
