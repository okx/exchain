package analyzer

const (
	READ           = 1
	WRITE          = 2
	EVMALL         = 3
	UNKNOWN_TYPE   = 4
	EVM_FORMAT     = "【read<%dms>%s], write<%dms>, execute<%dms>"
	EVMCORE        = "evmcore"
)

var (
	STATEDB_WRITE = []string{"AddBalance", "SubBalance", "SetNonce", "SetState", "SetCode", "AddLog", "AddPreimage", "AddRefund", "SubRefund", "AddAddressToAccessList", "AddSlotToAccessList", "PrepareAccessList", "AddressInAccessList", "Suicide", "CreateAccount", "ForEachStorage"}
	STATEDB_READ  = []string{"SlotInAccessList", "GetBalance", "GetNonce", "GetCode", "GetCodeSize", "GetCodeHash", "GetState", "GetCommittedState", "GetRefund", "HasSuicided", "Snapshot", "RevertToSnapshot", "Empty", "Exist"}
	EVM_OPER      = []string{EVMCORE}
	dbOper        *DbRecord
)
