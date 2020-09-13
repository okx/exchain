## 1. tokenPairStoreKey

KVStoreKey的name是"token_pair"

|       key                |          value           |   number(key) |    value details    | Value size |         Clean up         |      备注       |
| :----------------------: | :----------------------: | :----------: | :-----------------: | :--------: | :----------------------: | :-------------: |
|    0x01+$trading_pair   | struct x/dex.TokenPair    |   币对数      | TokenPair结构体如下 |    <1k     | 有接口删除上交易所的币对 | 存的交易币对的详细信息 |
| "tokenPairNumberKey"     |        uint64            |      1        |  dex运营方在okexchain发行的币对数量  |  <1k |                | 存的okexchain交易币对的数量 |
|     0x02+$owner_addr     |        WithdrawInfo      |  在途赎回数   |  dex运营方在okexchain发行的币对数量  |  <1k |                | 存的okexchain交易币对的数量 |
|     0x03+$trading_pair   |        product lock      |   币对数      |  dex运营方在okexchain发行的币对数量  |  <1k |                | 存的okexchain交易币对的数量 |

tokenPair的value是经过Codec.MustMarshalBinaryBare序列化后的[]byte，tokenPair的结构体如下:

```go
type TokenPair struct {
	BaseAssetSymbol  string  `json:"base_asset_symbol"`		  // 基础货币
	QuoteAssetSymbol string  `json:"quote_asset_symbol"`	  // 报价货币
	InitPrice        sdk.Dec `json:"price"`					  // 价格
	MaxPriceDigit    int64   `json:"max_price_digit"`	 	  // 最大交易价格的小数点位数
	MaxQuantityDigit int64   `json:"max_size_digit"`		  // 最大交易数量的小数点位数
	MinQuantity      sdk.Dec `json:"min_trade_size"`		  // 最小交易数量
	Id               uint64  `json:"token_pair_id"`
	Delisting        bool    `json:"delisting"` 		      // 该TokenPair是否处于提案下线中 delisting
	Owner            sdk.AccAddress `json:"owner" v2:"owner"`  // token的所有者
	Deposits		 sdk.DecCoins `json:"deposits"`            // 优先撮合成交金
}
```



## Http api

|       Url        | Method |       读key       |
| :--------------: | :----: | :--------------: |
| /dex/products    |  GET   | token_pair       |
| /dex/deposits    |  GET   | token_pair       |
| /dex/match_order |  GET   | token_pair       |


