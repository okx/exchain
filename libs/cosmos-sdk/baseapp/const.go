package baseapp

const (
	//----- DeliverTx
	DeliverTx = "DeliverTx"
	TxDecoder = "TxDecoder"
	RunTx     = "RunTx"

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
