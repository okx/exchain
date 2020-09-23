package farm

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"
)

const (
	lockedToken = "eth"
	yieldToken  = "xxb"
)

type FarmPool struct {
	PoolName          string
	LockedTokenSymbol string
	// sum of all lockedAmount
	TotalLockedToken       sdk.DecCoin
	YieldingCoins          sdk.DecCoins
	YieldedCoins           sdk.DecCoins
	LastBlockHeightToYield int64
	YieldAmountPerBlock    sdk.Dec
	// sum of all lockedAmount * lockedBlockHeight
	TotalLockedInfo sdk.Dec
}

type LockRecord struct {
	lockAmount       sdk.DecCoin
	startBlockHeight int64
}

type account struct {
	user    string
	balance sdk.DecCoins
}

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

	p.TotalLockedToken = p.TotalLockedToken.Add(amount)
	p.TotalLockedInfo = p.TotalLockedInfo.Add(amount.Amount.MulTruncate(sdk.NewDec(height)))
}

func (p *FarmPool) claim(user *account, height int64) {
	if lockRecords[user.user] == nil {
		return
	}

	if height > p.LastBlockHeightToYield {
		yieldAmt := p.YieldAmountPerBlock.MulTruncate(sdk.NewDec(height - p.LastBlockHeightToYield))
		yieldCoins := sdk.DecCoins{sdk.NewDecCoinFromDec(yieldToken, yieldAmt)}
		p.YieldingCoins = p.YieldingCoins.Sub(yieldCoins)
		p.YieldedCoins = p.YieldedCoins.Add(yieldCoins)
		p.LastBlockHeightToYield = height
	}

	record := lockRecords[user.user]

	numerator := record.lockAmount.Amount.MulTruncate(sdk.NewDec(height - record.startBlockHeight))
	denominator := p.TotalLockedToken.Amount.MulTruncate(sdk.NewDec(height)).Sub(p.TotalLockedInfo)
	yieldCoinsForUser := p.YieldedCoins.MulDecTruncate(numerator).QuoDecTruncate(denominator)

	// claim yield coins
	p.YieldedCoins = p.YieldedCoins.Sub(yieldCoinsForUser)
	user.balance = user.balance.Add(yieldCoinsForUser)
	p.TotalLockedInfo = p.TotalLockedInfo.Sub(record.lockAmount.Amount.MulTruncate(sdk.NewDec(record.startBlockHeight)))

	// update block height
	record.startBlockHeight = height

	p.TotalLockedInfo = p.TotalLockedInfo.Add(record.lockAmount.Amount.MulTruncate(sdk.NewDec(record.startBlockHeight)))
}

func (p *FarmPool) unlock(user *account, height int64) {
	record := lockRecords[user.user]
	if record == nil {
		return
	}
	p.claim(user, height)

	user.balance = user.balance.Add(sdk.DecCoins{record.lockAmount})

	p.TotalLockedToken = p.TotalLockedToken.Sub(record.lockAmount)
	p.TotalLockedInfo = p.TotalLockedInfo.Sub(record.lockAmount.Amount.MulTruncate(sdk.NewDec(record.startBlockHeight)))

	delete(lockRecords, user.user)
}

var lockRecords = make(map[string]*LockRecord)

func TestClaim(t *testing.T) {
	pool := FarmPool{
		PoolName:               "pool-xxb-eth",
		LockedTokenSymbol:      lockedToken,
		TotalLockedToken:       sdk.NewDecCoinFromDec(lockedToken, sdk.ZeroDec()),
		YieldingCoins:          sdk.DecCoins{sdk.NewDecCoinFromDec(yieldToken, sdk.NewDec(100000))},
		YieldAmountPerBlock:    sdk.NewDec(10),
		LastBlockHeightToYield: 0,
		TotalLockedInfo:        sdk.ZeroDec(),
	}

	userA := account{
		"A",
		sdk.DecCoins{sdk.NewDecCoinFromDec(lockedToken, sdk.NewDec(10000))},
	}

	userB := account{
		"B",
		sdk.DecCoins{sdk.NewDecCoinFromDec(lockedToken, sdk.NewDec(10000))},
	}

	pool.lock(&userA, sdk.NewDecCoinFromDec(lockedToken, sdk.NewDec(100)), 10)
	pool.lock(&userB, sdk.NewDecCoinFromDec(lockedToken, sdk.NewDec(150)), 20)

	pool.lock(&userA, sdk.NewDecCoinFromDec(lockedToken, sdk.NewDec(100)), 30)
	fmt.Printf("A balance: %s\n", userA.balance)
	pool.lock(&userB, sdk.NewDecCoinFromDec(lockedToken, sdk.NewDec(200)), 50)
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
