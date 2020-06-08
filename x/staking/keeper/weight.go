package keeper

import (
	"fmt"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/staking/types"
)

const (
	// UTC Time: 2000/1/1 00:00:00
	blockTimestampEpoch = int64(946684800)
	secondsPerWeek      = int64(60 * 60 * 24 * 7)
	weeksPerYear        = float64(52)
)

func calculateWeight(nowTime int64, tokens sdk.Dec) (shares types.Shares, sdkErr sdk.Error) {
	nowWeek := (nowTime - blockTimestampEpoch) / secondsPerWeek
	rate := float64(nowWeek) / weeksPerYear
	weight := math.Pow(float64(2), rate)
	weightByDec, sdkErr := sdk.NewDecFromStr(fmt.Sprintf("%.8f", weight))
	if sdkErr == nil {
		shares = tokens.Mul(weightByDec)
	}
	return
}
