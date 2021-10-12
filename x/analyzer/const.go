package analyzer

const (
	COMMIT_STATE_DB = "CommitStateDb"
	AVL             = "Iavl"
	EVM             = "EvmCall"
	EVM_Create      = "Create"
	EVM_Call        = "Call"
	READ            = 1
	WRITE           = 2
	EVM_FORMAT      = "read<%dms>, write<%dms>, execute<%dms>, call&create<%dms>"
)

var (
	module        = []string{COMMIT_STATE_DB, EVM}
	STATEDB_WRITE = []string{"AddBalance", "SubBalance", "SetNonce", "SetState", "SetCode", "AddLog", "AddPreimage", "AddRefund", "SubRefund", "AddAddressToAccessList", "AddSlotToAccessList", "PrepareAccessList", "AddressInAccessList", "Suicide", "CreateAccount", "ForEachStorage"}
	STATEDB_READ  = []string{"SlotInAccessList", "GetBalance", "GetNonce", "GetCode", "GetCodeSize", "GetCodeHash", "GetState", "GetCommittedState", "GetRefund", "HasSuicided", "Snapshot", "RevertToSnapshot", "Empty", "Exist"}
	EVM_OPER      = []string{"Create", "Call"}
	dbOper        *DbRecord
)
