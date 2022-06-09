package types

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

var contractVerifierPool = sync.Pool{
	New: func() interface{} {
		return new(ContractVerifier)
	},
}

func initContractVerifier(
	cv *ContractVerifier,
	params Params) {
	cv.SetParams(params)
}

var vmConfigPool = sync.Pool{
	New: func() interface{} {
		return new(vm.Config)
	},
}

func initVmConfig(
	vmConfig *vm.Config,
	extraEips []int, // Additional EIPS that are to be enabled
	debug bool,
	tracer vm.Tracer,
	contractVerifier *ContractVerifier) {

	vmConfig.ExtraEips = extraEips
	vmConfig.Debug = debug
	vmConfig.Tracer = tracer
	vmConfig.ContractVerifier = contractVerifier

	vmConfig.NoRecursion = false
	vmConfig.NoBaseFee = false
	vmConfig.EnablePreimageRecording = false
	vmConfig.JumpTable = vm.JumpTable{}
}

var vmBlockCtxPool = sync.Pool{
	New: func() interface{} {
		return new(vm.BlockContext)
	},
}

func initBlockCtx(
	vmBlockContext *vm.BlockContext,
	canTransfer vm.CanTransferFunc, // Transfer transfers ether from one account to the other
	transfer vm.TransferFunc, // GetHash returns the hash corresponding to n
	getHash vm.GetHashFunc,
	coinbase common.Address, // Provides information for COINBASE
	gasLimit uint64, // Provides information for GASLIMIT
	blockNumber *big.Int, // Provides information for NUMBER
	time *big.Int, // Provides information for TIME
	difficulty *big.Int, // Provides information for DIFFICULTY
) {
	vmBlockContext.CanTransfer = canTransfer
	vmBlockContext.Transfer = transfer
	vmBlockContext.GetHash = getHash
	vmBlockContext.Coinbase = coinbase
	vmBlockContext.GasLimit = gasLimit
	vmBlockContext.BlockNumber = blockNumber
	vmBlockContext.Time = time
	vmBlockContext.Difficulty = difficulty
	vmBlockContext.BaseFee = nil
}

var vmTxContextPool = sync.Pool{
	New: func() interface{} {
		return new(vm.TxContext)
	},
}

func initTxContext(
	vmTxContext *vm.TxContext,
	origin common.Address, // Provides information for ORIGIN
	gasPrice *big.Int, // Provides information for GASPRICE
) {
	vmTxContext.Origin = origin
	vmTxContext.GasPrice = gasPrice
}
