package farm

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/types"
	"testing"
)

const (
	symbolLocked = "eth"
	yieldToken   = "xxb"
)

type LockRecord struct {
	lockAmount       sdk.DecCoin
	startBlockHeight int64
}

type account struct {
	user    string
	balance sdk.DecCoins
}

type FarmPool types.FarmPool

func (p *FarmPool) lock(user *account, amount sdk.DecCoin, height int64) {
	var record *LockRecord
	if lockRecords[user.user] != nil {
		p.claim(user, height)
		record := lockRecords[user.user]
		record.lockAmount = record.lockAmount.Add(amount)
		record.startBlockHeight = height
	} else {
		record = &LockRecord{
			lockAmount:       amount,
			startBlockHeight: height,
		}
		lockRecords[user.user] = record
	}

	user.balance = user.balance.Sub(sdk.DecCoins{amount})

	p.TotalValueLocked = p.TotalValueLocked.Add(amount)
	p.TotalLockedWeight = p.TotalLockedWeight.Add(amount.Amount.MulTruncate(sdk.NewDec(height)))
}

func (p *FarmPool) claim(user *account, height int64) {
	if lockRecords[user.user] == nil {
		return
	}

	if height > p.LastClaimedBlockHeight {
		var yieldedCoins sdk.DecCoins
		for _, YieldedToken := range p.YieldedTokenInfos {
			if height > YieldedToken.StartBlockHeightToYield {
				yieldAmt := YieldedToken.AmountYieldedPerBlock.MulTruncate(sdk.NewDec(height - p.LastClaimedBlockHeight))
				amountYielded := sdk.NewDecCoinFromDec(yieldToken, yieldAmt)
				YieldedToken.RemainingAmount = YieldedToken.RemainingAmount.Sub(amountYielded)
				yieldedCoins = yieldedCoins.Add(sdk.DecCoins{amountYielded})
			}
		}

		p.AmountYielded = p.AmountYielded.Add(yieldedCoins)
		p.LastClaimedBlockHeight = height
	}

	record := lockRecords[user.user]

	numerator := record.lockAmount.Amount.MulTruncate(sdk.NewDec(height - record.startBlockHeight))
	denominator := p.TotalValueLocked.Amount.MulTruncate(sdk.NewDec(height)).Sub(p.TotalLockedWeight)
	yieldCoinsForUser := p.AmountYielded.MulDecTruncate(numerator).QuoDecTruncate(denominator)

	// claimRewards yield coins
	p.AmountYielded = p.AmountYielded.Sub(yieldCoinsForUser)
	user.balance = user.balance.Add(yieldCoinsForUser)
	p.TotalLockedWeight = p.TotalLockedWeight.Sub(record.lockAmount.Amount.MulTruncate(sdk.NewDec(record.startBlockHeight)))

	// update block height
	record.startBlockHeight = height

	p.TotalLockedWeight = p.TotalLockedWeight.Add(record.lockAmount.Amount.MulTruncate(sdk.NewDec(record.startBlockHeight)))
}

func (p *FarmPool) unlock(user *account, height int64) {
	record := lockRecords[user.user]
	if record == nil {
		return
	}
	p.claim(user, height)

	user.balance = user.balance.Add(sdk.DecCoins{record.lockAmount})

	p.TotalValueLocked = p.TotalValueLocked.Sub(record.lockAmount)
	p.TotalLockedWeight = p.TotalLockedWeight.Sub(record.lockAmount.Amount.MulTruncate(sdk.NewDec(record.startBlockHeight)))

	delete(lockRecords, user.user)
}

var lockRecords = make(map[string]*LockRecord)

func TestClaim(t *testing.T) {
	pool := FarmPool{
		Name:                   "pool-xxb-eth",
		SymbolLocked:           symbolLocked,
		TotalValueLocked:       sdk.NewDecCoinFromDec(symbolLocked, sdk.ZeroDec()),
		LastClaimedBlockHeight: 0,
		TotalLockedWeight:      sdk.ZeroDec(),
	}

	YieldedToken := types.YieldedTokenInfo{
		RemainingAmount:         sdk.NewDecCoinFromDec(yieldToken, sdk.NewDec(100000)),
		StartBlockHeightToYield: 0,
		AmountYieldedPerBlock:   sdk.NewDec(10),
	}

	pool.YieldedTokenInfos = append(pool.YieldedTokenInfos, YieldedToken)

	userA := account{
		"A",
		sdk.DecCoins{sdk.NewDecCoinFromDec(symbolLocked, sdk.NewDec(10000))},
	}

	userB := account{
		"B",
		sdk.DecCoins{sdk.NewDecCoinFromDec(symbolLocked, sdk.NewDec(10000))},
	}

	pool.lock(&userA, sdk.NewDecCoinFromDec(symbolLocked, sdk.NewDec(100)), 10)
	pool.lock(&userB, sdk.NewDecCoinFromDec(symbolLocked, sdk.NewDec(150)), 20)

	pool.lock(&userA, sdk.NewDecCoinFromDec(symbolLocked, sdk.NewDec(100)), 30)
	fmt.Printf("A balance: %s\n", userA.balance)
	pool.lock(&userB, sdk.NewDecCoinFromDec(symbolLocked, sdk.NewDec(200)), 50)
	fmt.Printf("B balance: %s\n", userB.balance)

	for user, record := range lockRecords {
		fmt.Printf("%s: locked: %s, height: %d\n", user, record.lockAmount, record.startBlockHeight)
	}

	pool.unlock(&userA, 100)
	fmt.Printf("A balance: %s\n", userA.balance)
	pool.unlock(&userB, 100)
	fmt.Printf("B balance: %s\n", userB.balance)

	for user, record := range lockRecords {
		fmt.Printf("%s: locked: %s, height: %d\n", user, record.lockAmount, record.startBlockHeight)
	}
}
