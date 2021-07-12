package farm

import (
	"encoding/csv"
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"os"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	swap "github.com/okex/exchain/x/ammswap/types"
	"github.com/okex/exchain/x/farm/keeper"
	"github.com/okex/exchain/x/farm/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker allocates the native token to the pools in PoolsYieldNativeToken
// according to the value of locked token in pool
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	logger := k.Logger(ctx)
	file, err := os.Create("swap_farm_data.csv")
	if err != nil {
		panic(fmt.Sprintf("create file field:%s", err.Error()))
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	farmHeader := []string{"Farm data", "", ""}
	if err := writer.Write(farmHeader); err != nil {
		panic(fmt.Sprintf("write file field:%s", err.Error()))
	}

	logger.Error("begin statistics farm data")
	farmPools := k.GetFarmPools(ctx)
	farmAccountMap := make(map[string]struct{})
	// accounts in genesis.json
	farmOldAccounts := []string{
		"ex1hv2ut4snuaevufkq074uesfwf0pputu73ckf7h",
		"ex14r4u28s7stngjh2hlszzn6mjwfzgupdpsnuwyu",
		"ex19sx64fl0y7xst8y78rqcsp3n8rz5j9uw39c2d9",
		"ex19luar0c4nr3vfzy9vv59nm6t5sw5chk38kp5vv",
		"ex132klqd60x5kj2pye4273qqxt566csv5j6zzny6",
		"ex14xqakd77k9awmknw9rfmdvf37h33y44qjj8c97",
		"ex1h4lz6j6vu9pd8ne3lgl90usllu484lge2nelal",
		"ex1vvza9zscmvhpdtzwvs639c8j4xjyqnn3kspffl",
		"ex1zc2h0nr76j6l772unncxtwp9pys5e5c2ddhewa",
		"ex1m3mmfalpqy2vym64ak2nu93ky74z5thx7j76fz",
		"ex1354z3nc30rwhm8lp5ahrthzljvj2xg6l9z0xz5",
		"ex15q7wz5f88797hzd5wa8lrs078wtttv5ecxrjnc",
		"ex1t7a9wr0haqfrfwya42vnczfdy0zswearwvfsm8",
		"ex12yxfwpls4udu33mqesn0eymnwx4kaymlazfkfa",
		"ex145v2s97upsqdpakcn4fmmn9hlu8q9wrn7tz8p0",
		"ex1253r38uekl59e4l4lxu3njww63qpmsgecclzkc",
		"ex1gge9ry8d2cakyh8qlhveku4zcruhlmf96pyk5g",
		"ex1lsc6aek2n4769uu06rfnrfnsux0mx2nm43sjau",
		"ex12a2x07ej8zmqe6hd0ahyplxtflfh4aq8tntwkr",
		"ex1yfeuc69zhdg7yfj4yuu6ff8ead8lc7an93n0zq",
		"ex1h68p60ae3yj769u727vqh4wesqvf0rc6vg3p6n",
		"ex1yqtq4smwv585rw7t54w4sxqu35g62epgf78g83",
		"ex1l2nx3s73l02w7kdhnualhyjhrezv8pq5lr66t2",
		"ex1d8w7rdnu34s7dx4cfryg7xhzsct6grlga4amtk",
		"ex1h3szz3f6mwnx6erlf7z23rpdk2mgsdtp3fsfss",
		"ex1nrwfs3c6wpahs2pzdyy6vnwedp67yrsfx8pk65",
		"ex1zjaex3ffkwqp0e75gh4l7a8n0vyv997pn7marq",
		"ex10yvu7kdavphnhue2ux5d6ktpjqmtzpdcfmkvkj",
		"ex17e8m9s3r8kjv4gv985yn67wu6hgr7ljmefhc4p",
		"ex19qykf9mpqgyjquy3wf340k5344ncysxkqlkvkx",
		"ex15qw4ndtwjlj26t9azggfdhzgwlq5hjh9uz70cs",
		"ex10w7ty0a9wwmjf5457l73nmj7c6m8nk8w7hq623",
		"ex12fku56apvxexalpvd5xwet332s6cr7meuyff7d",
		"ex1zs4tqk56qz80f99yatjx2c7rtm73ctvjv3rlus",
		"ex1jx0nw38kzcv2fnr0knud53dpazuf0974anlwjq",
		"ex1f8wv2npwlcy9nmc8nep2zkj5u0tpz0jdxy9epv",
		"ex19zqpr9mh9stuuql0mg8qagl0234gtn5nvuhlvw",
		"ex1mgsng0uk8k63632jdk86ank9vu2yyjltgtnwdx",
		"ex1t3tdptpqr45tmplmkpw4gsffuyqvqk9n3dtmpf",
		"ex1fxns3jy980xrxx976sqauqlzks4lhcxnkl5yjp",
		"ex1ea6vzf8yqyd7uq0h03yd8qztd4vxcvpzdddx6v",
		"ex1nr77ucd3lldteyj4742eke8wmputq5l70ddcce",
		"ex1gc6ssmc3t3p9tjtvrvkrv0hjpszw20ln69xx03",
		"ex13xg6texm63wtj8eq7ukkfx0uqkhzzu5xzw7844",
		"ex18v9x38md6sjgqwvt3pdkpsf49e62rc8qgpwykt",
		"ex1flkw8lu4f6ec5rkrj3pxu8787326d9wlyjmrcj",
		"ex1tmflc79mqvy4w4gdxwhjle74rnseexcxfkkvlg",
		"ex1rau28agfl5y3hvgncsz4a7vssexd87mf7eyc5a",
		"ex1hwfk3z4fxfjpf0uzny4zg695xvn6lvtyvutqrm",
		"ex19we9nuckmx735zx2hwa7rfkwakdssar4vv4huc",
		"ex16gqxz5vdaf2kfccfym79088a8d70cx83d0mcxr",
		"ex12xdx5xtwm2ag603ek0cqfwqz7f4hpt7pq00n03",
		"ex1e3nap5hvypxnxyqquzqkkls7aqswv2jfyjxqxf",
		"ex1c4ygs8llgpp0fgly6avmj8ca0tz2hr3k3a4uem",
		"ex1yldzuqw6g7q4gm9q2nz7h9mdv3j9ulj3z276wy",
		"ex1qv2cdvuzcc9m8mcs6v7dqn4r57xhmvm5adwhyl",
		"ex1aw5usygm5dv349qk8suagyll5esvk7a0hm36h2",
		"ex1yeccscm6222yc7qhhgsdxdts7v98aa67zelegg",
		"ex195py03p7k20dhrn9s38w3s8e40smcn8ujzulyc",
		"ex1azpzljjg700xxlfmcekpwyj5wep0m39520a0mn",
		"ex1eydsgegl3g4clq7qet30uwr2dhk623mqxzz3j8",
		"ex18djhdp2e5g6xnvf7zxm2axnwhmnrqjppy3clca",
		"ex10ukrx3ljpwwf9nl03gckj4z3xf20d886ethwap",
		"ex17m4h6mx6s3q6kv7mqnjhzerp6qezcc5z5rwcnv",
		"ex1wjqgzdduza9mrypuarhr0x4dcmjgl43qgp85kp",
		"ex1wymegpjxhea4lr7295mv8z8yssnnv94ykch5wq",
		"ex1rucgzfytk2g3xxvafgad88ysvkqszjljsmhc94",
		"ex12fek899re55phpdfegugl75xz78axl5rxd4hjm",
		"ex19z503meeeth8ff6eh2vnhqj2fakzcdwhmrz29l",
		"ex1xqm660l7smasm8fz8g67hkz694e6md32ha89ke",
		"ex1m5wrvtzm759ykkq72sexhm0qdtjz5t80zzuwcf",
		"ex1wmzug73us4qnzf6jcuvrqzyy522ntzf66xynqa",
		"ex18mxvgmf5sqh9ra6x8ty8muwl70l4gxjrsy36kf",
		"ex1ae2gd7denpgw96n27vh3jfqmpsju79a5gn4sjq",
		"ex1q3cc42w3j4wjtucwrqa8me3rdlc7370ddtjk77",
		"ex1afq7qfvhuvgccpct8fp4ka43nlulzx3vwsfdyc",
		"ex1enteqf7l929vhkzrp7380ezqgy2ynanpujlz7f",
		"ex1fp43zdajshv2dlu8nc02sced4emdgyqdu3ut6v",
		"ex1sqpwhmjqlm6uk0mgp87z6305y38y5e50ry2s3g",
		"ex10thqayn08lcmqvmpdz8wh0deuytlpxampzewzy",
		"ex158ut4yxcvuhfu5ajrac4l848xlguujpp745znr",
		"ex1wf8qrlj8j7evmkzdjyd93lv6rz9hzlclnl5s9c",
		"ex1g4w3c27wgqnqm4j47gg63g76yg0erc5q5p8q67",
		"ex1pyprtt2gas069dy9zc62qd6z5s8tmr22al9l6e",
		"ex1uw7hayh0e2r7mlp74sek8u2cz0m0nf87ztg57q",
		"ex13s3ss2lamqy9tthwkkmskfxrptf30rxxwt9zkr",
		"ex1eyf4yfwn09f2aduwl6wmdfljvpwkul7kt9yprn",
		"ex100908jq6ywfm8wqkdveq74zjnwcaphg9mg6c6f",
		"ex1c88jplex6rq4n0mtl20608zcvrlz3g8f524a9t",
		"ex1z66u7yqpanw3xdzhd5san4029p4uzls6kv8670",
		"ex1w6ctq82dywsejpslts2rmuqtejqklrkpy73f74",
		"ex1t3dnten99w5a0re5a24ustrjlckvu4fks9ejld",
		"ex19rj3w902reqywf4hrpf0pw7kxvkth332pmr95c",
		"ex1uksy89xh2hjvaqa2rgth9q44cq3k4uklhr2t0e",
		"ex1r2vvy7v20v0gg7e6w229x7g896atd9qcl65rr3",
		"ex156wuy046e2ss4c5q9hhzs5ltmtwa9x7vs55q9k",
		"ex1hjf8z9c797sa2zr6lrgavdnxaxmga0m00pqwcz",
		"ex1flm9xf2tp9ty843q37u4l4ax0kzvnu9xfj2h3p",
		"ex1qqn2tspenxgv4rs4jk3ehk9an76ygjyawu3a85",
		"ex1nm48qjgd3s2p2pcxgggwyqu4x6h5zf6aa5vvp9",
		"ex1m8r0kpfxf6w002nwqxqzkq22v9zxuwjrz53nvn",
		"ex1glgh8a9z86yw52fx5yatt5jl9hfwkt5sx83633",
		"ex1lzquw3ntgyu9m79prjwachjl530clqhreq9axz",
		"ex1xdpxyap5qt6x4l49r37qnezyy9c2nzdafr5cds",
		"ex1a8rzmdyf7y2mmzsps8xnrgrsr86k2ajjc74sls",
		"ex1delh6xsqgm8d7swngr0extxdpu6pkajwwn8sj7",
		"ex1hdzwpflqtusdw78gzhxdq064sj495zdx5awnv5",
		"ex1067rqcywj9n0wc6m0apc5fk5luhq9ch6s98vws",
		"ex1lkmrah920l0nrwdvyr7fxu3zdth4mqztg4fjhm",
		"ex1ck9lzxuk26duczm3zu3ah3jg4n4srx024u8a4k",
		"ex139rpcrl6zwc7eghy480uk9xaehdsupefw75gk2",
		"ex17p3hg85edy9fxlu0gx64aq2kmcv070mm0xpgrv",
		"ex1yrrjjfrrcxgzxf9xuvqfzr5t5367x53aycjmt9",
		"ex13x9670xt4u2yxq2qt0gj2wxdnkh3ysk9rkem96",
		"ex1tqszn8ttcrmlgwq5lq68a793v0ksxfql7wfj8g",
		"ex1gh2myyrmlwyppystm8lmzjsehwguycg7d20r2x",
		"ex1zh0uer4vqgyxwnedvrzv30v5uzerjrwe6rez6y",
		"ex1sccyj8g8j3rwk3rp2x4c5pvu6cun4ud7gw2yvr",
		"ex1shnvexvyql2xfaqf0y9eyxhw34lw9u3n54zsa6",
		"ex1qd8k2ppcdxzdmnpq4vd4dpp32qhqzf48s778n7",
		"ex1925fqpfk0pse7jpsr0zz3f2qm3zg2gxsgqx86n",
		"ex1sveuye373g0f6w4xxx6e3xx3x9z4hqk082w8yl",
		"ex1ec32f6gd4kqcve268m5dzkwnp7z2gjvcrgsal5",
		"ex1dgusmdnmmqx39unk0a4t8n0sujrl76va7n3m8m",
		"ex10a355kst0uel6lutp069hrjx70sxwl8ahu0qat",
		"ex16wp5c9vk790jhhjr77fs9qdtxu4498g379lres",
		"ex13dvxnzwa4ahk4ec5dd0ep7fpyye6w5njk7dl65",
		"ex1tw95jcthxf25gp74hvxg7s9wrfxj3z8309wwd6",
		"ex1x83wdvvatpxujw0sfr5fnt4ed65k44sj7ym8ae",
		"ex1uz62yp08jlc8h2m0n835dt2qp9y4hkuuma6uhd",
		"ex1482smu6jztlqg800pq9qvudf97n8h53jhxln6e",
		"ex1c7kja475e3ycuqy0kqedppamcctx0c44ge79jz",
		"ex1d0050vhpanyexhtq24nz9jw7qcm5j4wd2gnsxp",
		"ex1ry54jmyjlknq7qwv2smzq7ux5477znk2uswqed",
		"ex1y7u6gux9weudr56wsvhzu96njdcqwlcp4dqj9t",
		"ex1pd55sucrh2adkedc6xexp8hf0n4za9zznu53ym",
		"ex16534lp7nwf0gscvu897lmv2j904fs7nqe3vgu6",
		"ex1mk2gkre5ljeq2c8f35r80qrw50v6k7amhjmzjv",
		"ex1kkzchv56z7d9j2pdrqjewqh925ntzqq7sxmcdw",
		"ex1jwxtzfww6hk562p987y3elf5mwy4j9l8gyr6nv",
		"ex1f3x7q3l0zd6s7lglkyxtkag80ckyhqkk8hsucn",
		"ex1ym4992qu9y372snw5x69nuz7epkh882nl0elw0",
		"ex1va6nxsmdt4mexgkhf4pw3k20m9k7x3gtn9w2z5",
		"ex19ues3r8dh2qgp2t4v5epyls03jx4wu0sh8xvvt",
		"ex19sx64fl0y7xst8y78rqcsp3n8rz5j9uw39c2d9",
		"ex1fj42rzf72qxhj4h8ffqwdkfmng9rgquuqw9kr4",
		"ex18smk3xzzxf92k8ltvked4udkvr069cmmq3ug6u",
		"ex1p7hd4x6d8qpn5xw2dppk925szr7388xlx8kg8s",
		"ex1c96hp7mmrk9a5eyqfzkzm6j2aqr3t2lqwjvf9z",
		"ex1npwdsjfm5ak8gaq7shehjz4s0enzd9e2x7asxg",
		"ex1ccmwjrxypn2u9pzv735kywrnh03tggqhhxf6mv",
		"ex1gddmlvulruapl7v3g7y564rn9af69dwqyl6dpr",
		"ex1ua5hkakfscnq8dx39we4d6ch8ksx7z9efs309q",
		"ex1s9rxk4s8ehykelx5zfagvdkgdcmrt8q3d8afzl",
		"ex1rja2yfkrw7z65x22k3km3v56mrd42pe5r934xw",
		"ex186p3893xarpy6dhyp6u7q7tat25hhmc2xzcr6l",
		"ex14uet3ztrl6npnsat2gucukwasszfev9tyjxeyz",
		"ex1p90yfg5f365fkgamyucw0lfd8g5v67z6a8jzjd",
		"ex1ekx0yxcfy69fqprjvtf8r6a5pj2yuuxavx6gh4",
		"ex1tespjplxjutsm8w9vyvdv6qz9jwxw5nmtmmzag",
		"ex1nl0htgq4qehe999cawv8t9pv9es83y37787xj0",
		"ex1pn0cudq5d04yhk2vqwk2j4as6w4yvrs8smrvya",
		"ex1t3fd7weteljk7fzf04cznp6gkx4d4qzxllfkhs",
		"ex10mhdds9w9za8gdgqtyvss5ldg5c7lux8jn9nx8",
		"ex18na9uj5y6tmktt3z9nkn2qy9rfg23y3lerc02v",
		"ex1lfjj2dk348yg8y8sjwx7c0x8ehsnvzn0saqfyl",
		"ex18h6jcsj05p6wep778hjmfqzy9vmd2eva9w2sr6",
		"ex1ug7nqdkcpugrty9g00puqk25m7q57xmk87xxwf",
		"ex12ec65h92tn84ad96ks2cda2yya3w4qjx46fldp",
		"ex12lmjq64crgm4ga687gfhgpw6d3qnhcahk8t0my",
		"ex1ukdc94gkz6659wesqq20n57hluatgtuyq9phmv",
		"ex198t9ly4p7nzye92lxeu757tjqfx0zfvg8jud49",
		"ex1zmafu93ad88whf0t85nvg8hgngks2gp3e2ueau",
		"ex1ynm9d8eyqf9u7rencxd9hntfn4p6v47xqqhy9t",
		"ex10waxmga89z5ent8q0tn4f4ku42ttpx2t6mut4f",
		"ex1tpecvcets3fg8yw60ul39qcwd8uuufy54cu4kw",
		"ex1srpgld39t8hpev7qv3d47ktc3m9t7jt2w7y4wy",
		"ex1zucnq3gul2pd4zqmatmmc0fujpqxhpxm7j5p66",
		"ex1dgjt2mf543hvqw2nrvfgpjlxh27y04skpdkddm",
		"ex10yvu7kdavphnhue2ux5d6ktpjqmtzpdcfmkvkj",
		"ex1x5q9jd2faadkfl3w03n4lmfjd7aclfv6vfdgur",
		"ex1c88jplex6rq4n0mtl20608zcvrlz3g8f524a9t",
		"ex1a28wdeucnm2w0ympc97spsgvedcaph63hqjpv4",
		"ex1s5gvxpwnh5s9xvnljjdyedum4vpfj8507efnt4",
		"ex1cda835v4rm9v24k2fae8zjcsfspc4gqcukqydf",
		"ex1t78taw734pn8c8kk3eqsalc2q2sthxp5srgjv8",
		"ex18lff6f47fkh4yppqevvgrs5rdp5yqyne9wpm64",
		"ex1rja2yfkrw7z65x22k3km3v56mrd42pe5r934xw",
		"ex144lhh2fx3q5c9epdmrcj9h48akxjqt3fx0akyr",
		"ex1mpxynesyjvan3zx9wk3pyn8hrjhmfhy9x4ge5z",
		"ex1rlnu2ztp837mpsmmcvrnv55rz60mxgqpp38kzd",
		"ex1qv2cdvuzcc9m8mcs6v7dqn4r57xhmvm5adwhyl",
		"ex1laa83gu2pqwazfnfmywvyvfhrzuhd9jng78xj9",
		"ex1z7tvssm7sts4mfx94gpdnjfgjqen8vzyj950ey",
		"ex1aw5usygm5dv349qk8suagyll5esvk7a0hm36h2",
		"ex1kk3kz3ke2kafg70vtgml0yad0r6wqa86enxcwa",
		"ex1d0zqntwc58m7c3795gts74g30uydyexs6dpp5g",
		"ex1r6a3p9z9df8fcuscgfe4awrkfn03kgrl662tts",
		"ex158ut4yxcvuhfu5ajrac4l848xlguujpp745znr",
		"ex1ug7nqdkcpugrty9g00puqk25m7q57xmk87xxwf",
		"ex1ddkaw7rxjhdd0r0nqlntx89sqqxv2ectluswlv",
		"ex1eyl246htax76e5kg585c5mufnzkkp58vgrfmsd",
		"ex1243zjspjuz3z7dc3cdl2vnuz3uks9mvl9f8vhn",
		"ex1qd8k2ppcdxzdmnpq4vd4dpp32qhqzf48s778n7",
		"ex1mh5ufppfsz2rrazxd8f9uwlntzmy7jqyjpasvf",
		"ex13cc0heq9lfmw6uygff62qzwzr0fa7alh980ckv",
		"ex1wywd3u3vhhjrmhgks9lafc0v6rw3rz8fyyvfrh",
		"ex16str2zcu7py24l6qhel5n8k3dhmd775ydv9pws",
		"ex1lkmrah920l0nrwdvyr7fxu3zdth4mqztg4fjhm",
		"ex1vj9rjcwft0hx5q295f0xxy6egcvvdahharnkzp",
		"ex12x8gnd60zfav2vqhtjcl8dp972ua38hut4w8nc",
		"ex19zmgp8eyy8ltczwthc2zug5m99a579ytuwq8ry",
		"ex1yth0q2z4gnk4rcl5xk0n7vq2uvffkq97rtexa6",
		"ex1my0f4q36a5y05q2x4evd9mppq95jsm2nphs3d2",
		"ex1zucnq3gul2pd4zqmatmmc0fujpqxhpxm7j5p66",
		"ex18smk3xzzxf92k8ltvked4udkvr069cmmq3ug6u",
		"ex1dgjt2mf543hvqw2nrvfgpjlxh27y04skpdkddm",
		"ex136vegqcvfq8zd8j09wpvucna388tq2tqtt3c4a",
		"ex17e8m9s3r8kjv4gv985yn67wu6hgr7ljmefhc4p",
		"ex1cda835v4rm9v24k2fae8zjcsfspc4gqcukqydf",
		"ex17luwgl2qtg8wngjcelfsah83yfvfrq4wlrjttl",
		"ex1hw0eakq9ucanwxzqdhlp7cygxdhsmvl7nd7j0v",
		"ex1lv5zh5ngqxt42r95ru72xa235u047kctvxf2vc",
		"ex1yldzuqw6g7q4gm9q2nz7h9mdv3j9ulj3z276wy",
		"ex1qv2cdvuzcc9m8mcs6v7dqn4r57xhmvm5adwhyl",
		"ex1semqzalu5xavpskm9fcg3c25g8du3mhzu2amjk",
		"ex1laa83gu2pqwazfnfmywvyvfhrzuhd9jng78xj9",
		"ex1aw5usygm5dv349qk8suagyll5esvk7a0hm36h2",
		"ex1p56slq3u0kyrg97sh9epm6j4hntxrqzu3f6esj",
		"ex1wydar80ydw6v6vxdp9ujl5urhn76xpe0uknxyw",
		"ex1va677rggvqz3wt06j2vmlpkhrwy73zyvujk09m",
		"ex1s5l57gspsx47hghczz87v00hcj928s2wmcu8ek",
		"ex1r6a3p9z9df8fcuscgfe4awrkfn03kgrl662tts",
		"ex1gk6j2km964sjhuxsx007dukhmwfw2jejmv6she",
		"ex1925fqpfk0pse7jpsr0zz3f2qm3zg2gxsgqx86n",
		"ex1wmhdzygxptds3xyzcwxzad3ew37e47cv2pjx95",
		"ex13cc0heq9lfmw6uygff62qzwzr0fa7alh980ckv",
		"ex16udrzupa9ee2hzk9k50ep02xen27rupqh5pr28",
		"ex1ems2sxvvc9ppx596t57m5ehnprlk3xq7e5p97r",
		"ex1cn40ugqh267yy8w6tl3eelvkt4ddmej6qyzrau",
		"ex1vj9rjcwft0hx5q295f0xxy6egcvvdahharnkzp",
		"ex12x8gnd60zfav2vqhtjcl8dp972ua38hut4w8nc",
		"ex1n2ejnmh7s99h5vznh52rkxg4f83wcytrg5g688",
		"ex1flm9xf2tp9ty843q37u4l4ax0kzvnu9xfj2h3p",
		"ex1lnt2gm4dfty0tmymn858tuc8h3hqvc5wq04tm5",
		"ex1cndpna52mk5u3qcyrp6w6m8tqdmxhsrwgre23n",
		"ex1z8urn0ycnagq9a008tjw9mkkntxu7e8fncfscm",
		"ex1py90xptl424ltn50dvqpgtwzvjx94x92nch9dt",
		"ex1ddk9n32ntw0lyn4z2uqqf2ge7qd6k23e3wp9ev",
		"ex1qf5599w7jdvhfj0ty0fzl0m6xm2nrwyu85sfg3",
		"ex1u5l49kvk85u84z2ck7la4a2pty8sc7f0z9uvdl",
		"ex16kpfnwhewdtprc85ed780an364cx9c63cehzmd",
		"ex1xhgznnrl3za89fcqnzd7fh0pqjff3hyuq8y5s7",
		"ex1z58lwae2r8qzh26gc60fet423ntt9pgkgr4z6p",
		"ex1cgztp6qvxhhrl4emsgwtz53r67lcx92ap3js0c",
		"ex1va6nxsmdt4mexgkhf4pw3k20m9k7x3gtn9w2z5",
		"ex1my0f4q36a5y05q2x4evd9mppq95jsm2nphs3d2",
	}

	for _, addr := range farmOldAccounts {
		farmAccountMap[addr] = struct{}{}
	}

	// accounts after upgrade
	for _, pool := range farmPools {
		lockedAccounts := k.GetAccountsLockedTo(ctx, pool.Name)
		for _, account := range lockedAccounts {
			farmAccountMap[account.String()] = struct{}{}
		}
	}

	for userAccount, _ := range farmAccountMap {
		addr, err := sdk.AccAddressFromBech32(userAccount)
		if err != nil {
			panic(fmt.Sprintf("user account error:%s", err.Error()))
		}

		var totalLP sdk.SysCoins
		for _, pool := range farmPools {
			lockInfo, ok := k.GetLockInfo(ctx, addr, pool.Name)
			if !ok {
				continue
			}
			totalLP = append(totalLP, lockInfo.Amount)
		}

		if totalLP == nil {
			// logger.Error(fmt.Sprintf("address:%s has no lp", userAccount))
			continue
		}

		var sumCoins sdk.SysCoins
		for _, lp := range totalLP {
			tokens := strings.SplitN(strings.SplitN(lp.Denom, swap.PoolTokenPrefix, 2)[1], "_", 2)
			coin0, coin1, err := k.SwapKeeper().GetRedeemableAssets(ctx, tokens[0], tokens[1], lp.Amount)
			if err != nil {
				panic(fmt.Sprintf("GetRedeemableAssets error:%s", err.Error()))
			}
			if sumCoins == nil {
				sumCoins = append(sumCoins, coin0, coin1)
			} else {
				sumCoins = sumCoins.Add(coin0, coin1)
			}
		}
		logger.Error(fmt.Sprintf("address:%s farm总量（换算单币之和）:%s farm 总量（lp总量）:%s", userAccount,
			sumCoins.String(), totalLP.String()))
		for _, lp := range totalLP {
			if err := writer.Write([]string{userAccount, lp.Amount.String(), lp.Denom}); err != nil {
				panic(fmt.Sprintf("write file field:%s", err.Error()))
			}
		}

	}

	logger.Error("end statistics farm data")

	ctx.Logger().Error("begin statistics swap data")
	swapHeaders := [][]string{{"", "", ""}, {"", "", ""}, {"", "", ""}, {"Swap data", "", ""}}
	for _, swapHeader := range swapHeaders {
		if err := writer.Write(swapHeader); err != nil {
			panic(fmt.Sprintf("write file field:%s", err.Error()))
		}
	}

	// all accounts
	k.SwapKeeper().AccountKeeper.IterateAccounts(ctx, func(account exported.Account) bool {
		var totalLP sdk.SysCoins
		coins := account.GetCoins()
		for _, coin := range coins {
			if !strings.HasPrefix(coin.Denom, swap.PoolTokenPrefix) {
				continue
			}
			totalLP = append(totalLP, coin)
		}

		if totalLP == nil {
			// logger.Error(fmt.Sprintf("address:%s has no lp", userAccount))
			return false
		}

		var sumCoins sdk.SysCoins
		for _, lp := range totalLP {
			tokens := strings.SplitN(strings.SplitN(lp.Denom, swap.PoolTokenPrefix, 2)[1], "_", 2)
			coin0, coin1, err := k.SwapKeeper().GetRedeemableAssets(ctx, tokens[0], tokens[1], lp.Amount)
			if err != nil {
				panic(fmt.Sprintf("GetRedeemableAssets error:%s", err.Error()))
			}
			if sumCoins == nil {
				sumCoins = append(sumCoins, coin0, coin1)
			} else {
				sumCoins = sumCoins.Add(coin0, coin1)
			}
		}
		logger.Error(fmt.Sprintf("address:%s swap总量（换算单币之和）:%s swap 总量（lp总量）:%s", account.GetAddress().String(),
			sumCoins.String(), totalLP.String()))
		for _, lp := range totalLP {
			if err := writer.Write([]string{account.GetAddress().String(), lp.Amount.String(), lp.Denom}); err != nil {
				panic(fmt.Sprintf("write file field:%s", err.Error()))
			}
		}
		return false
	})
	ctx.Logger().Error("end statistics swap data")
	panic("stop")
	moduleAcc := k.SupplyKeeper().GetModuleAccount(ctx, MintFarmingAccount)
	yieldedNativeTokenAmt := moduleAcc.GetCoins().AmountOf(sdk.DefaultBondDenom)
	logger.Debug(fmt.Sprintf("MintFarmingAccount [%s] balance: %s%s",
		moduleAcc.GetAddress(), yieldedNativeTokenAmt, sdk.DefaultBondDenom))

	if yieldedNativeTokenAmt.LTE(sdk.ZeroDec()) {
		return
	}

	yieldedNativeToken := sdk.NewDecCoinsFromDec(sdk.DefaultBondDenom, yieldedNativeTokenAmt)
	// 0. check the YieldNativeToken parameters
	params := k.GetParams(ctx)
	if !params.YieldNativeToken { // if it is false, only burn the minted native token
		if err := k.SupplyKeeper().BurnCoins(ctx, MintFarmingAccount, yieldedNativeToken); err != nil {
			panic(err)
		}
		return
	}

	// 1. gets all pools in PoolsYieldNativeToken
	lockedPoolValueMap, pools, totalPoolsValue := calculateAllocateInfo(ctx, k)
	if totalPoolsValue.LTE(sdk.ZeroDec()) {
		return
	}

	// 2. allocate native token to pools according to the value
	remainingNativeTokenAmt := yieldedNativeTokenAmt
	for i, pool := range pools {
		var allocatedAmt sdk.Dec
		if i == len(pools)-1 {
			allocatedAmt = remainingNativeTokenAmt
		} else {
			allocatedAmt = lockedPoolValueMap[pool.Name].
				MulTruncate(yieldedNativeTokenAmt).QuoTruncate(totalPoolsValue)
		}
		remainingNativeTokenAmt = remainingNativeTokenAmt.Sub(allocatedAmt)
		logger.Debug(
			fmt.Sprintf("Pool %s allocate %s yielded native token", pool.Name, allocatedAmt.String()),
		)
		allocatedCoins := sdk.NewDecCoinsFromDec(sdk.DefaultBondDenom, allocatedAmt)

		current := k.GetPoolCurrentRewards(ctx, pool.Name)
		current.Rewards = current.Rewards.Add2(allocatedCoins)
		k.SetPoolCurrentRewards(ctx, pool.Name, current)
		logger.Debug(fmt.Sprintf("Pool %s rewards are %s", pool.Name, current.Rewards))

		pool.TotalAccumulatedRewards = pool.TotalAccumulatedRewards.Add2(allocatedCoins)
		k.SetFarmPool(ctx, pool)
	}
	if !remainingNativeTokenAmt.IsZero() {
		panic(fmt.Sprintf("there are some tokens %s not to be allocated", remainingNativeTokenAmt))
	}

	// 3.liquidate native token minted at current block for yield farming
	err = k.SupplyKeeper().SendCoinsFromModuleToModule(ctx, MintFarmingAccount, YieldFarmingAccount, yieldedNativeToken)
	if err != nil {
		panic("should not happen")
	}
}

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {

}

// calculateAllocateInfo gets all pools in PoolsYieldNativeToken
func calculateAllocateInfo(ctx sdk.Context, k keeper.Keeper) (map[string]sdk.Dec, []types.FarmPool, sdk.Dec) {
	lockedPoolValue := make(map[string]sdk.Dec)
	var pools types.FarmPools
	totalPoolsValue := sdk.ZeroDec()

	store := ctx.KVStore(k.StoreKey())
	iterator := sdk.KVStorePrefixIterator(store, types.PoolsYieldNativeTokenPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		poolName := types.SplitPoolsYieldNativeTokenKey(iterator.Key())
		pool, found := k.GetFarmPool(ctx, poolName)
		if !found {
			panic("should not happen")
		}
		poolValue := k.GetPoolLockedValue(ctx, pool)
		if poolValue.LTE(sdk.ZeroDec()) {
			continue
		}
		pools = append(pools, pool)
		lockedPoolValue[poolName] = poolValue
		totalPoolsValue = totalPoolsValue.Add(poolValue)
	}
	return lockedPoolValue, pools, totalPoolsValue
}
