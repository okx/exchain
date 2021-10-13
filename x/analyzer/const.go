package analyzer

const (
	READ              = 1
	WRITE             = 2
	EVMALL            = 3
	ANTEHANDLE        = 4
	EVM_FORMAT        = "read<%dms>, write<%dms>, execute<%dms>"
	ANTEHANDLE_FORMAT = "execute<%dms>"
	EVMCORE           = "evmcore"
)

var (
	STATEDB_WRITE = []string{"AddBalance", "SubBalance", "SetNonce", "SetState", "SetCode", "AddLog", "AddPreimage", "AddRefund", "SubRefund", "AddAddressToAccessList", "AddSlotToAccessList", "PrepareAccessList", "AddressInAccessList", "Suicide", "CreateAccount", "ForEachStorage"}
	STATEDB_READ  = []string{"SlotInAccessList", "GetBalance", "GetNonce", "GetCode", "GetCodeSize", "GetCodeHash", "GetState", "GetCommittedState", "GetRefund", "HasSuicided", "Snapshot", "RevertToSnapshot", "Empty", "Exist"}
	EVM_OPER      = []string{EVMCORE}
	ANTE_HANDLE   = []string{"anteHandler"}
	dbOper        *DbRecord
)
