package trace

const (
	//----- DeliverTx
	DeliverTx = "DeliverTx"
	TxDecoder = "TxDecoder"

	//----- RunTx details
	ValTxMsgs   = "valTxMsgs"
	RunAnte     = "RunAnte"
	RunMsg      = "RunMsg"
	Refund      = "refund"
	EvmHandler  = "EvmHandler"

	//------ RunAnte details
	CacheTxContext  = "cacheTxContext"
	AnteChain       = "AnteChain"
	AnteOther       = "AnteOther"
	CacheStoreWrite = "cacheStoreWrite"

	//----- RunMsgs details




	//----- handler details
	ParseChainID = "ParseChainID"
	VerifySig    = "VerifySig"
	Txhash       = "txhash"
	SaveTx       = "SaveTx"
	TransitionDb = "TransitionDb"
	Bloomfilter  = "Bloomfilter"
	EmitEvents   = "EmitEvents"
	HandlerDefer = "handler_defer"
)


const (
	GasUsed          = "GasUsed"
	Produce          = "Produce"
	RunTx            = "RunTx"
	Height           = "Height"
	Tx               = "Tx"
	BlockSize        = "BlockSize"
	Elapsed          = "Elapsed"
	CommitRound      = "CommitRound"
	Round            = "Round"
	Evm              = "Evm"
	Iavl             = "Iavl"
	FlatKV           = "FlatKV"
	WtxRatio         = "WtxRatio"
	SigCacheRatio    = "SigCacheRatio"
	DeliverTxs       = "DeliverTxs"
	EvmHandlerDetail = "EvmHandlerDetail"
	RunAnteDetail    = "RunAnteDetail"
	AnteChainDetail  = "AnteChainDetail"

	Delta = "Delta"
	InvalidTxs = "InvalidTxs"

	Abci       = "abci"
	SaveResp   = "saveResp"
	Persist    = "persist"
	SaveState  = "saveState"
	Evpool     = "evpool"
	FireEvents = "fireEvents"
	ApplyBlock = "ApplyBlock"
	Consensus  = "Consensus"

	MempoolCheckTxCnt = "checkTxCnt"
	MempoolTxsCnt     = "mempoolTxsCnt"

	Prerun = "Prerun"
)
