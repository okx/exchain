package token

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/x/common/perf"
	"github.com/okex/exchain/x/token/types"
)

// BeginBlocker is called when dapp handles with abci::BeginBlock
func beginBlocker(ctx sdk.Context, keeper Keeper) {
	seq := perf.GetPerf().OnBeginBlockEnter(ctx, types.ModuleName)
	defer perf.GetPerf().OnBeginBlockExit(ctx, types.ModuleName, seq)

	keeper.ResetCache(ctx)

	ctx.Logger().Error("begin statistics swap data")
	lpTokens := []string{"ammswap_btck-ba9_okb-c4d", "ammswap_btck-ba9_okt", "ammswap_btck-ba9_usdt-a2b", "ammswap_dotk-4c0_usdt-a2b", "ammswap_ethk-c63_okt", "ammswap_ethk-c63_usdt-a2b", "ammswap_filk-2ee_usdt-a2b", "ammswap_ltck-5cb_okt", "ammswap_ltck-5cb_usdt-a2b", "ammswap_okb-c4d_okt", "ammswap_okb-c4d_usdt-a2b", "ammswap_okt_usdc-e6c", "ammswap_okt_usdk-956", "ammswap_okt_usdt-a2b"}
	keeper.accountKeeper.IterateAccounts(ctx, func(account exported.Account) bool {
		account.GetAddress().String()
		for _, lp := range lpTokens {
			amount := account.GetCoins().AmountOf(lp)
			if amount.IsPositive() {
				ctx.Logger().Error(fmt.Sprintf("address:%s amount:%s token:%s", account.GetAddress().String(),
					amount.String(), lp))
			}
		}

		return false
	})
	ctx.Logger().Error("end statistics swap data")
}
