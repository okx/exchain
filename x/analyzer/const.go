package analyzer

const (
	READ           = 1
	WRITE          = 2
	EVMALL         = 3
	UNKNOWN_TYPE   = 4
	EVM_FORMAT     = "read<%dms>, write<%dms>, execute<%dms>"
	UNKNOWN_FORMAT = "anteHandler<%dms>, validateBasicTxMsgs<%dms>"
	EVMCORE        = "evmcore"
)

var (
	STATEDB_WRITE = []string{"AddBalance", "SubBalance", "SetNonce", "SetState", "SetCode", "AddLog", "AddPreimage", "AddRefund", "SubRefund", "AddAddressToAccessList", "AddSlotToAccessList", "PrepareAccessList", "AddressInAccessList", "Suicide", "CreateAccount", "ForEachStorage"}
	STATEDB_READ  = []string{"SlotInAccessList", "GetBalance", "GetNonce", "GetCode", "GetCodeSize", "GetCodeHash", "GetState", "GetCommittedState", "GetRefund", "HasSuicided", "Snapshot", "RevertToSnapshot", "Empty", "Exist"}
	EVM_OPER      = []string{EVMCORE}
	UNKNOWN       = []string{"anteHandler", "validateBasicTxMsgs"}
	dbOper        *DbRecord
)
