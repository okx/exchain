package system

const (
	Chain = "okbchain"
	AppName = "OKBChain"
	Server = Chain+"d"
	Client = Chain+"cli"
	ServerHome = "$HOME/."+Server
	ClientHome = "$HOME/."+Client
	ServerLog = Server+".log"
	EnvPrefix = "OKBCHAIN"
	CoinType uint32 = 60
	Currency = "okb"
)
