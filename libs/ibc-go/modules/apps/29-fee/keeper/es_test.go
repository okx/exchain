package keeper

import (
	"fmt"
	"testing"

	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/client/utils"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

func TestParse(t *testing.T) {
	str := "0.014170600000000000okt"
	cc, err := sdk.ParseDecCoins(str)
	fmt.Println(err)
	fmt.Println(cc.String())
	timeoutPortion := sdk.NewDecWithPrec(20, 2)
	ret := cc.MulDecTruncate(timeoutPortion)
	vv := utils.CliConvertCoinToCoinAdapters(ret)
	vv2 := sdk.CoinsToCoinAdapters(ret)
	fmt.Println(ret.String())
	fmt.Println(vv.String())
	fmt.Println(vv2.String())
	fmt.Println(vv2.ToCoins().String())
}
