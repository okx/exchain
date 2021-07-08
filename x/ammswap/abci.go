package ammswap

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/exchain/x/ammswap/types"
	"strings"
)

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block
func BeginBlocker(ctx sdk.Context, keeper Keeper) {
	ctx.Logger().Error("begin statistics swap data")
	userAccounts := []string{
		"ex1er8hq8es45ga8m580h8dp4m54vk6j9vctsrlv5",
		"ex1rzwzzuduw094zcmf2zdjhy3ejhfrwytg0uctwr",
		"ex1rvjdngnc907jfjld09ygkr54h20wudyucrgftc",
		"ex1tsk07nte44c5uscj03h6amqh270vsnl9zj6xrk",
		"ex165fu2a25pxuykhpmffjpxntd5vejmvttylmyge",
		"ex19apkzla8kjccxum4rjsaar25paeu665yn78dqy",
		"ex1rhx4p9dar5khyzptvg93gl66v60f9jtepgxuhc",
		"ex1eytduadennwled3rx0209rttavgjjm2g67yqga",
		"ex1a49lsz3htyd3rvjfvrtxg7fup70ejj2gw7qasc",
		"ex19crw89z4cclsuxzrevg4rhalnw8jjz32pn0en5",
		"ex1y9hn4uqkhwhftj5fqkkmnp6vva4x0tl2ksat06",
		"ex139fzfukfy566e37r45lzsqmwhtmxhlyuu0g49f",
		"ex1vhnu683ztf6ggk07wawaqere6lvrjy4gt3k6r0",
		"ex1vz4aj9gtftltz3g0w4vmg5pd2dmw0w2cg2atjj",
		"ex1rdlrzf8fwuwck8qzu0t9pu70ghdrwcd3xmtca4",
		"ex19jmudnj5wnekswv0nlkwf0dkmw23lhanaed58m",
		"ex1hk90szz0hw90z0smfw6zldey66u06z29cwc7f8",
		"ex1zxc5nhuu2lvffvlpd4c345l7tfr7tqma99g992",
		"ex18x4yzahz0rs2c3dstcrfpttk09lr2hpwhqgsat",
		"ex1lvuc7ueq9x8mtyqwm4z3e8wvfkrj62gpscurlr",
		"ex1nfushs6ea3x87ne2yty2v7jla7q8qfg7vjr079",
		"ex1xj34a5klkyre0vwwplt806hrj5txy5xn9lan6h",
		"ex1p9ej7pdjz62yq0g95wsvcxhqwzhehu2evxerfg",
		"ex1e3nap5hvypxnxyqquzqkkls7aqswv2jfyjxqxf",
		"ex1yfrklx6cav55vznzdfw2vhy0e3x5j7jspdz86y",
		"ex1wyyg5qmxy446dpnj27kr80mz2k5awkc9u4rqha",
		"ex1ejw6zvmajlgyw3xdap69jzxjpte6lnyfxmf4kc",
		"ex13ps5v0yeh5lhdgsquznan3y3xn44p35q08tz3w",
		"ex19te4zcqenvtmute6pk454gpcevyxxnyux80fhg",
		"ex1tm76y22wyhgewzxcnwvev5pw80w3achmdsr9hg",
		"ex1zecfdg2hv86qldkhxx70kukp47h3rschqn00r8",
		"ex16359hn56q320c53g8vkxr0m2n9rj44gwmj4my4",
		"ex1m5q97uvxkd86vxraygklztth48zxr7de0pghpk",
		"ex1rnp3p2z7emanhzww8quyflskxdpzymwfgaggyv",
		"ex1lgjwfxl53cpa0ta3pyqpaxtnx6txzdxt0jxqlc",
		"ex19sx64fl0y7xst8y78rqcsp3n8rz5j9uw39c2d9",
		"ex12s98wz369m5y7mlzggq5mt6fw7sqachc2xeueg",
		"ex14ldyvl47gh0gxuk4089gqhgqeyfgd9dll7cgqv",
		"ex1a8w4lkh594rt7q8qfg6s94p62x2whkekc986nq",
		"ex123u9qnmdfdgh82j9kmg4r3442l46xg2xf2srtc",
		"ex146h5sujcxq6plmak4zpa0lpv7468uuxrd5u90x",
		"ex12r6ywfnn6zrfh67fftr6t35lyqkcp5jkyneqr3",
		"ex1cceq8xhygsyr2whrx4u8juwrgqux48ll7hrexm",
		"ex1t3ff2jqjpal7zqnr2h58jtvjxy2y7p4zv33rtp",
		"ex1c2p4lynq9726jmrj5w7jhp6sct4fmyht7affmk",
		"ex1u8j4dk3623hpg9agz0nfz3u3yh3rdhjpj2ku7w",
		"ex1t99tppnfg3dlwahtxc06wt6zl6qg5qx5z0w9ce",
		"ex142xl92p50s34p7wp4putuu45cfper3j0m8c9uz",
		"ex1uw5vjd2q6fewwrs0xtme9kf2ddf4wqhe5esghv",
		"ex1qjysw79e5mak98c3fvq6l8q578sme9zngr902s",
		"ex17jrkcrchfnpmjeyt0rww504246j7x2vfaqrnyc",
		"ex10zzq9jxfwaqm7rre8uagsy0agh0zgsqqvfkqur",
		"ex18g7clpenamtqqtkej4s4mpxw39qp60pux5tdc8",
		"ex108lhfpf33ha2uhv6ryhvcpl7r8wuhy7nnrnxdc",
		"ex1q3www762cvh377kv259d2l5p9yn8jpl39zeu4n",
		"ex100n77xmx56ff6f4jl4mvdyyclu2jzrhmf7asr6",
		"ex18va54g7mtxmy80g9vgqpm7z8kz0eh67ux4xdfw",
		"ex1eyf4yfwn09f2aduwl6wmdfljvpwkul7kt9yprn",
		"ex1phy2t74w34jx9ljh2vtdydtmqhnxm0vwh6wkrs",
		"ex15a9xjkyyhtqlph05g2pg785xeyteequarm9h0m",
		"ex1ec32f6gd4kqcve268m5dzkwnp7z2gjvcrgsal5",
		"ex1qgekve9pvxqt0uxgztyge8q2vkjngs3u3ttsk6",
		"ex1fa6ps295udzdm4gm0jt3q6h0yumv5sdjfz4w7j",
		"ex15l3n89697cnkx7y7pq6vdqtfl7p74grktph3kr",
		"ex1amrpaw2uc0mhrpvjp87z2zvu62ylut0yxlhkwt",
		"ex12yxfwpls4udu33mqesn0eymnwx4kaymlazfkfa",
		"ex1tcrq3vpau4d0pepk8vlhfq3k6pzvwy5avesqz9",
		"ex18djhdp2e5g6xnvf7zxm2axnwhmnrqjppy3clca",
		"ex19z503meeeth8ff6eh2vnhqj2fakzcdwhmrz29l",
		"ex1lrvescvdgkf5u5jy04n82nq5v0fwxsn3cwk6wn",
		"ex1yth0q2z4gnk4rcl5xk0n7vq2uvffkq97rtexa6",
		"ex1wmhdzygxptds3xyzcwxzad3ew37e47cv2pjx95",
		"ex1kkzchv56z7d9j2pdrqjewqh925ntzqq7sxmcdw",
		"ex1mgsng0uk8k63632jdk86ank9vu2yyjltgtnwdx",
		"ex1zh0uer4vqgyxwnedvrzv30v5uzerjrwe6rez6y",
		"ex1nrwfs3c6wpahs2pzdyy6vnwedp67yrsfx8pk65",
		"ex19rj3w902reqywf4hrpf0pw7kxvkth332pmr95c",
	}
	for _, userAccount := range userAccounts {
		addr, err := sdk.AccAddressFromBech32(userAccount)
		if err != nil {
			panic(fmt.Sprintf("user account error:%s", err.Error()))
		}

		var totalLP sdk.SysCoins
		account := keeper.AccountKeeper.GetAccount(ctx, addr)
		coins := account.GetCoins()
		for _, coin := range coins {
			if !strings.HasPrefix(coin.Denom, types.PoolTokenPrefix) {
				continue
			}
			totalLP = append(totalLP, coin)
		}

		if totalLP == nil {
			keeper.Logger(ctx).Error(fmt.Sprintf("address:%s has no lp", userAccount))
			continue
		}

		var sumCoins sdk.SysCoins
		for _, lp := range totalLP {
			tokens := strings.SplitN(strings.SplitN(lp.Denom, types.PoolTokenPrefix, 2)[1], "_", 2)
			coin0, coin1, err := keeper.GetRedeemableAssets(ctx, tokens[0], tokens[1], lp.Amount)
			if err != nil {
				panic(fmt.Sprintf("GetRedeemableAssets error:%s", err.Error()))
			}
			if sumCoins == nil {
				sumCoins = append(sumCoins, coin0, coin1)
			} else {
				sumCoins = sumCoins.Add(coin0, coin1)
			}
		}
		keeper.Logger(ctx).Error(fmt.Sprintf("address:%s sum amount:%s total lp:%s", account.String(),
			sumCoins.String(), totalLP.String()))
	}

	ctx.Logger().Error("end statistics swap data")
}

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k Keeper) {
}
