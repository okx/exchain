package trace

const (
	//----- DeliverTx
	DeliverTx = "DeliverTx"
	TxDecoder = "TxDecoder"

	//----- RunTx details
	ValTxMsgs  = "valTxMsgs"
	RunAnte    = "RunAnte"
	RunMsg     = "RunMsg"
	Refund     = "refund"
	EvmHandler = "EvmHandler"

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
	GasUsed     = "GasUsed"
	SimGasUsed  = "SimGasUsed"
	Produce     = "Produce"
	RunTx       = "RunTx"
	LastRun     = "lastRun"
	Height      = "Height"
	Tx          = "Tx"
	SimTx       = "SimTx"
	BlockSize   = "BlockSize"
	Elapsed     = "Elapsed"
	CommitRound = "CommitRound"
	Round       = "Round"
	BlockParts  = "BlockParts"
	Evm         = "Evm"
	Iavl        = "Iavl"
	FlatKV      = "FlatKV"
	//RecvBlock        = "RecvBlock"
	First2LastPart = "First2LastPart"

	SigCacheRatio    = "SigCacheRatio"
	DeliverTxs       = "DeliverTxs"
	EvmHandlerDetail = "EvmHandlerDetail"
	RunAnteDetail    = "RunAnteDetail"
	AnteChainDetail  = "AnteChainDetail"

	Delta      = "Delta"
	InvalidTxs = "InvalidTxs"

	Abci = "abci"
	//SaveResp        = "saveResp"
	Persist        = "persist"
	PersistDetails = "persistDetails"
	PreChange      = "preChange"
	FlushCache     = "flushCache"
	CommitStores   = "commitStores"
	FlushMeta      = "flushMeta"

	//MempoolUpdate   = "mpUpdate"
	//SaveState       = "saveState"
	ApplyBlock    = "ApplyBlock"
	Consensus     = "Consensus"
	LastBlockTime = "LastBlockTime"
	BTInterval    = "BTInterval"
	RecommendedGP = "RecommendedGP"
	IsCongested   = "IsCongested"
	UpdateState   = "UpdateState"
	Waiting       = "Waiting"

	MempoolCheckTxCnt  = "CheckTx"
	MempoolTxsCnt      = "MempoolTxs"
	MempoolCheckTxTime = "CheckTxTime"

	CompressBlock   = "Compress"
	UncompressBlock = "Uncompress"
	Prerun          = "Prerun"
	IavlRuntime     = "IavlRuntime"

	BlockPartsP2P = "BlockPartsP2P"

	Workload = "Workload"
	ACOffset = "ACOffset"
)

const (
	READ         = 1
	WRITE        = 2
	EVMALL       = 3
	UNKNOWN_TYPE = 4
	EVM_FORMAT   = "read<%dms>, write<%dms>, execute<%dms>"
	EVMCORE      = "evmcore"
)

var (
	STATEDB_WRITE = []string{"AddBalance", "SubBalance", "SetNonce", "SetState", "SetCode", "AddLog",
		"AddPreimage", "AddRefund", "SubRefund", "AddAddressToAccessList", "AddSlotToAccessList",
		"PrepareAccessList", "AddressInAccessList", "Suicide", "CreateAccount", "ForEachStorage"}

	STATEDB_READ = []string{"SlotInAccessList", "GetBalance", "GetNonce", "GetCode", "GetCodeSize",
		"GetCodeHash", "GetState", "GetCommittedState", "GetRefund",
		"HasSuicided", "Snapshot", "RevertToSnapshot", "Empty", "Exist"}

	EVM_OPER = []string{EVMCORE}
	dbOper   *DbRecord
)
