package types

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/vm"
)

// TraceExecutionResult groups all structured logs emitted by the EVM
// while replaying a transaction in debug mode as well as transaction
//// execution status, the amount of gas used and the return value
//type TraceExecutionResult struct {
//	Gas         uint64         `json:"gas"`
//	Failed      bool           `json:"failed"`
//	ReturnValue string         `json:"returnValue"`
//	StructLogs  []StructLogRes `json:"structLogs"`
//	//LogsLen     int             `json:"logsLen""`
//}
//
//// StructLogRes stores a structured log emitted by the EVM while replaying a
//// transaction in debug mode
//type StructLogRes struct {
//	Pc      uint64             `json:"pc"`
//	Op      string             `json:"op"`
//	Gas     uint64             `json:"gas"`
//	GasCost uint64             `json:"gasCost"`
//	Depth   int                `json:"depth"`
//	Error   error              `json:"error,omitempty"`
//	Stack   *[]string          `json:"stack,omitempty"`
//	Memory  *[]string          `json:"memory,omitempty"`
//	Storage *map[string]string `json:"storage,omitempty"`
//}

// FormatLogs formats EVM returned structured logs for json output
func FormatLogs(logs []vm.StructLog) []*StructLogRes {
	formatted := make([]*StructLogRes, len(logs))
	for index, trace := range logs {
		err := ""
		if trace.Err != nil {
			err = trace.Err.Error()
		}
		formatted[index] = &StructLogRes{
			Pc:      trace.Pc,
			Op:      trace.Op.String(),
			Gas:     trace.Gas,
			GasCost: trace.GasCost,
			Depth:   int64(trace.Depth),
			Error:   err,
		}
		if trace.Stack != nil {
			stack := make([]string, len(trace.Stack))
			for i, stackValue := range trace.Stack {
				stack[i] = fmt.Sprintf("%x", math.PaddedBigBytes(stackValue, 32))
			}
			formatted[index].Stack = stack
		}
		//if trace.Memory != nil {
		//	memory := make([]string, 0, (len(trace.Memory)+31)/32)
		//	for i := 0; i+32 <= len(trace.Memory); i += 32 {
		//		memory = append(memory, fmt.Sprintf("%x", trace.Memory[i:i+32]))
		//	}
		//	formatted[index].Memory = memory
		//}
		//if trace.Storage != nil {
		//	storage := make(map[string]string)
		//	for i, storageValue := range trace.Storage {
		//		storage[fmt.Sprintf("%x", i)] = fmt.Sprintf("%x", storageValue)
		//	}
		//	formatted[index].Storage = storage
		//}
	}
	return formatted
}

//
////FormatLog formats EVM returned structured logs for json output
//func FormatLog(log *vm.StructLog) *StructLogRes {
//	formatted := &StructLogRes{
//		Pc:      log.Pc,
//		Op:      log.Op.String(),
//		Gas:     log.Gas,
//		GasCost: log.GasCost,
//		Depth:   log.Depth,
//		Error:   log.Err,
//	}
//	if log.Stack != nil {
//		stack := make([]string, len(log.Stack))
//		for i, stackValue := range log.Stack {
//			stack[i] = fmt.Sprintf("%x", math.PaddedBigBytes(stackValue, 32))
//		}
//		formatted.Stack = &stack
//	}
//	if log.Memory != nil {
//		memory := make([]string, 0, (len(log.Memory)+31)/32)
//		for i := 0; i+32 <= len(log.Memory); i += 32 {
//			memory = append(memory, fmt.Sprintf("%x", log.Memory[i:i+32]))
//		}
//		formatted.Memory = &memory
//	}
//	if log.Storage != nil {
//		storage := make(map[string]string)
//		for i, storageValue := range log.Storage {
//			storage[fmt.Sprintf("%x", i)] = fmt.Sprintf("%x", storageValue)
//		}
//		formatted.Storage = &storage
//	}
//
//	return formatted
//}
