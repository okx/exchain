package keeper

import (
	"fmt"
	"math"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/staking/types"
)

const (
	// UTC Time: 2000/1/1 00:00:00
	blockTimestampEpoch = int64(946684800)
	secondsPerWeek      = int64(60 * 60 * 24 * 7)
	weeksPerYear        = float64(52)
	constTimeStamp      = int64(1685577600) // 2023-06-01 00:00:00 GMT+0
	constNowWeek        = (constTimeStamp - blockTimestampEpoch) / secondsPerWeek
)

func calculateWeight(nowTime int64, tokens sdk.Dec, height int64, fixedValue bool) (shares types.Shares, sdkErr error) {
	var nowWeek int64
	if fixedValue {
		nowWeek = constNowWeek
	} else {
		nowWeek = (nowTime - blockTimestampEpoch) / secondsPerWeek
	}

	rate := float64(nowWeek) / weeksPerYear
	weight := math.Pow(float64(2), rate)

	precision := fmt.Sprintf("%d", sdk.Precision)

	weightByDec, sdkErr := sdk.NewDecFromStr(fmt.Sprintf("%."+precision+"f", weight))
	if sdkErr == nil {
		shares = tokens.Mul(weightByDec)
	}
	return
}

func SimulateWeight(nowTime int64, tokens sdk.Dec, height int64) (votes types.Shares, sdkErr error) {
	return calculateWeight(nowTime, tokens, height, false)
}
