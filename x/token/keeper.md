## StoreKey

### 1. accountStoreKey

KVStoreKey的name是"acc"

|         key          |     value      | number(key)  |     value details     |               Value size               | Clean up | 备注                    |
| :------------------: | :------------: | :----------: | :-------------------: | :------------------------------------: | :------: | ----------------------- |
| prefix(0x01)+address | struct Account | 账户数量 | Account具体结构体如下 | 可能会超过1K，取决于用户下面的币的数量 |    无    | 存用户账号下的token信息 |

value经过Codec.MustMarshalBinaryBare序列化后的[]byte，Account的结构体如下:

```go
type BaseAccount struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         sdk.DecCoins      `json:"coins"`
	PubKey        crypto.PubKey  `json:"public_key"`
	AccountNumber uint64         `json:"account_number"`
	Sequence      uint64         `json:"sequence"`
}
```

### 2. tokenStoreKey 

KVStoreKey的name是"token"

|key                              | value                      | number(key)          | value details | Value size | Clean up |  备注   |
|:-------------------------------:|:-------------------------:|:------:|:------:|:------:|:------:|:------:|
| symbol | struct x/token.Token | 发行币的数量 | Token具体结构体如下 | <1k | 无 | 存的token的信息 |

token的value是经过Codec.MustMarshalBinaryBare序列化后的[]byte，token的结构体如下:

```go
type Token struct {
	Name           string         `json:"name"`							// token的名字
	Symbol         string         `json:"symbol"`						// token的唯一标识
	OriginalSymbol string         `json:"original_symbol"`	// token的原始标识
	TotalSupply    int64          `json:"total_supply"`			// token的总量
	Owner          sdk.AccAddress `json:"owner"`						// token的所有者
	Mintable       bool           `json:"mintable"`					// token是否可以增发
}
```

### 3. freezeStoreKey 

KVStoreKey的name是"freeze"

|   key   |      value       |   number(key)    | value details |               Value size               | Clean up |      备注       |
| :-----: | :--------------: | :--------------: | :-----------: | :------------------------------------: | :------: | :-------------: |
| address | struct sdk.DecCoins | 有冻结币的用户数 |  Coins结构体  | 可能会超过1K，取决于用户冻结的币的数量 |    无    | 存的token的信息 |

token的value是Coins经过Codec.MustMarshalBinaryBare序列化后的[]byte.

## 4. lockStoreKey

KVStoreKey的name是"lock"

|   key   |      value       |    number(key)     | value details |               Value size               | Clean up |      备注       |
| :-----: | :--------------: | :----------------: | :-----------: | :------------------------------------: | :------: | :-------------: |
| address | struct sdk.DecCoins | 有锁定的币的用户数 |  Coins结构体  | 可能会超过1K，取决于用户锁定的币的数量 |    无    | 存的token的信息 |

token的value是Coins经过Codec.MustMarshalBinaryBare序列化后的[]byte.

## 5. tokenPairStoreKey

KVStoreKey的name是"token_pair"

|                 key                  |          value           |   number(key)    |    value details    | Value size |         Clean up         |      备注       |
| :----------------------------------: | :----------------------: | :--------------: | :-----------------: | :--------: | :----------------------: | :-------------: |
| BaseAssetSymbol+"_"+QuoteAssetSymbol | struct x/token.TokenPair | 上交易所的币对数 | TokenPair结构体如下 |    <1k     | 有接口删除上交易所的币对 | 存的token的信息 |

tokenPair的value是经过Codec.MustMarshalBinaryBare序列化后的[]byte，tokenPair的结构体如下:

```go
type TokenPair struct {
	BaseAssetSymbol  string  `json:"base_asset_symbol"`		// 基础货币
	QuoteAssetSymbol string  `json:"quote_asset_symbol"`	// 报价货币
	InitPrice        sdk.Dec `json:"price"`							  // 价格
	MaxPriceDigit    int64   `json:"max_price_digit"`	 	  // 最大交易价格的小数点位数
	MaxQuantityDigit int64   `json:"max_size_digit"`		  // 最大交易数量的小数点位数
	MinQuantity      sdk.Dec `json:"min_trade_size"`		  // 最小交易数量
}
```

## Http api

|       Url       | Method |      读key       |
| :-------------: | :----: | :--------------: |
|     /products     |  GET   | token_pair(遍历) |
|     /tokens      |  GET   |   token(遍历)    |
| /token/{symbol} |  GET   |  token: symbol   |

