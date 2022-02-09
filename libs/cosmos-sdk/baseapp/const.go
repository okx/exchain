package baseapp

const (
	//----- DeliverTx
	DeliverTx = "DeliverTx"
	TxDecoder = "TxDecoder"
	RunTx     = "RunTx"

	//----- run_tx
	InitCtx     = "initCtx"
	ValTxMsgs   = "valTxMsgs"
	AnteHandler = "anteHandler"
	RunMsgs     = "runMsgs"
	Refund      = "refund"
	ConsumeGas  = "ConsumeGas"
	Recover     = "recover"

	//----- handler
	EvmHandler   = "evm_handler"
	ParseChainID = "ParseChainID"
	VerifySig    = "VerifySig"
	Txhash       = "txhash"
	SaveTx       = "SaveTx"
	TransitionDb = "TransitionDb"
	Bloomfilter  = "Bloomfilter"
	EmitEvents   = "EmitEvents"
	HandlerDefer = "handler_defer"
)
