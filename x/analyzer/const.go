package analyzer

const (
	COMMIT_STATE_DB = "CommitStateDb"
	AVL             = "iavl"
	READ            = 1
	WRITE           = 2
	EVM_FORMAT      = "Evm[read<%dms>, write<%dms>, execute<%dms>]"
)

var (
	module        = []string{COMMIT_STATE_DB}
	STATEDB_WRITE = []string{"AddBalance", "SubBalance", "SetNonce", "SetState", "SetCode", "AddLog", "AddPreimage", "AddRefund", "SubRefund", "AddAddressToAccessList", "AddSlotToAccessList", "PrepareAccessList", "AddressInAccessList", "Suicide", "CreateAccount", "ForEachStorage"}
	STATEDB_READ  = []string{"SlotInAccessList", "GetBalance", "GetNonce", "GetCode", "GetCodeSize", "GetCodeHash", "GetState", "GetCommittedState", "GetRefund", "HasSuicided", "Snapshot", "RevertToSnapshot", "Empty", "Exist"}
	dbOper        *DbRecord
)
