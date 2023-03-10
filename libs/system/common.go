package system

const (
	Chain = "exchain"
	AppName = "OKExChain"
	Server = Chain+"d"
	Client = Chain+"cli"
	ServerHome = "$HOME/."+Server
	ClientHome = "$HOME/."+Client
	ServerLog = Server+".log"
	EnvPrefix = "OKEXCHAIN"
	CoinType uint32 = 60
)
