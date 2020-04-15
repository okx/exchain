package types

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/crypto"
)

var (
	regOriginalSymbol = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]{0,5}$")
	reWholeName       = `[a-zA-Z0-9[:space:]]{1,30}`
	reWhole           = regexp.MustCompile(fmt.Sprintf(`^%s$`, reWholeName))
)

func WholeNameCheck(wholeName string) (newName string, isValid bool) {
	wholeName = strings.TrimSpace(wholeName)
	strs := strings.Fields(wholeName)
	wholeName = strings.Join(strs, " ")
	if !reWhole.MatchString(wholeName) {
		return wholeName, false
	}
	return wholeName, true
}

func wholeNameValid(wholeName string) bool {
	return reWhole.MatchString(wholeName)
}

type BaseAccount struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         sdk.Coins      `json:"coins"`
	PubKey        crypto.PubKey  `json:"public_key"`
	AccountNumber uint64         `json:"account_number"`
	Sequence      uint64         `json:"sequence"`
}

type DecAccount struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         sdk.DecCoins   `json:"coins"`
	PubKey        crypto.PubKey  `json:"public_key"`
	AccountNumber uint64         `json:"account_number"`
	Sequence      uint64         `json:"sequence"`
}

// String implements fmt.Stringer
func (acc DecAccount) String() string {
	var pubkey string

	if acc.PubKey != nil {
		pubkey = sdk.MustBech32ifyAccPub(acc.PubKey)
	}

	return fmt.Sprintf(`Account:
 Address:       %s
 Pubkey:        %s
 Coins:         %v
 AccountNumber: %d
 Sequence:      %d`,
		acc.Address, pubkey, acc.Coins, acc.AccountNumber, acc.Sequence,
	)
}

func ValidOriginalSymbol(name string) bool {
	return regOriginalSymbol.MatchString(name)
}

// Convert a formatted json string into a TransferUnit array
// e.g.) [{"to": "addr", "amount": "1BNB,2BTC"}, ...]
func StrToTransfers(str string) (transfers []TransferUnit, err error) {
	var transfer []Transfer
	err = json.Unmarshal([]byte(str), &transfer)
	if err != nil {
		return transfers, err
	}

	for _, trans := range transfer {
		var t TransferUnit
		to, err := sdk.AccAddressFromBech32(trans.To)
		if err != nil {
			return transfers, fmt.Errorf("invalid address：%s", trans.To)
		}
		t.To = to
		t.Coins, err = sdk.ParseDecCoins(trans.Amount)
		if err != nil {
			return transfers, err
		}
		transfers = append(transfers, t)
	}
	return transfers, nil
}

func BaseAccountToDecAccount(account auth.BaseAccount) DecAccount {
	var decCoins sdk.DecCoins
	for _, coin := range account.Coins {
		dec := coin.Amount
		decCoin := sdk.NewDecCoinFromDec(coin.Denom, dec)
		decCoins = append(decCoins, decCoin)
	}
	decAccount := DecAccount{
		Address:       account.Address,
		PubKey:        account.PubKey,
		Coins:         decCoins,
		AccountNumber: account.AccountNumber,
		Sequence:      account.Sequence,
	}
	return decAccount
}

func (acc *DecAccount) ToBaseAccount() *auth.BaseAccount {
	decAccount := auth.BaseAccount{
		Address:       acc.Address,
		PubKey:        acc.PubKey,
		Coins:         acc.Coins,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
	}
	return &decAccount
}

func DecAccountArrToBaseAccountArr(decAccounts []DecAccount) (baseAccountArr []auth.Account) {
	for _, decAccount := range decAccounts {
		baseAccountArr = append(baseAccountArr, decAccount.ToBaseAccount())
	}
	return baseAccountArr
}

func MergeCoinInfo(availableCoins, lockCoins sdk.DecCoins) (coinsInfo CoinsInfo) {
	m := make(map[string]CoinInfo)

	for _, availableCoin := range availableCoins {
		coinInfo, ok := m[availableCoin.Denom]
		if !ok {
			coinInfo.Symbol = availableCoin.Denom

			coinInfo.Available = availableCoin.Amount.String()
			coinInfo.Locked = "0"
			m[availableCoin.Denom] = coinInfo
		}
	}

	for _, lockCoin := range lockCoins {
		coinInfo, ok := m[lockCoin.Denom]
		if ok {
			coinInfo.Locked = lockCoin.Amount.String()
			m[lockCoin.Denom] = coinInfo
		} else {
			coinInfo.Symbol = lockCoin.Denom
			coinInfo.Available = "0"
			coinInfo.Locked = lockCoin.Amount.String()

			m[lockCoin.Denom] = coinInfo
		}
	}

	for _, coinInfo := range m {
		coinsInfo = append(coinsInfo, coinInfo)
	}
	sort.Sort(coinsInfo)
	return coinsInfo
}
