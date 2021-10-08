package analyzer

const (
	COMMIT_STATE_DB = "CommitStateDb"
	AVL             = "iavl"
	BLOCK_FORMAT    = "Block: Height<%d>, allCost<%dms>, evmCost<%dms>, Tx[ %s ]"
	TX_FORMAT       = "{TxID: %d, allCost<%dms>, db:[read<%dms>, write<%dms>], evmCost<%dms>} "
	TX_DEBUG_FORMAT = "[oper: %s , callTimes<%d>, cost:<%dms>]"
	READ            = 1
	WRITE           = 2
)

var (
	module        = []string{COMMIT_STATE_DB}
	STATEDB_WRITE = []string{"AddBalance", "SubBalance", "SetNonce", "SetState", "SetCode", "AddLog", "AddPreimage", "AddRefund", "SubRefund", "AddAddressToAccessList", "AddSlotToAccessList", "PrepareAccessList", "AddressInAccessList", "Suicide", "CreateAccount", "ForEachStorage"}
	STATEDB_READ  = []string{"SlotInAccessList", "GetBalance", "GetNonce", "GetCode", "GetCodeSize", "GetCodeHash", "GetState", "GetCommittedState", "GetRefund", "HasSuicided", "Snapshot", "RevertToSnapshot", "Empty", "Exist"}
	dbOper        *DbRecord
)
