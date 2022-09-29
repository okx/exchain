package dydx

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/libs/dydx/contracts"
)

type Contracts struct {
	P1Orders        *contracts.P1Orders
	P1OrdersAddress common.Address

	PerpetualV1        *contracts.PerpetualV1
	PerpetualV1Address common.Address

	txOps *bind.TransactOpts
}

func NewContracts(
	perpetualV1Address common.Address,
	p1OrdersAddr common.Address,
	defaultTxOps *bind.TransactOpts,
	backend bind.ContractBackend,
) (*Contracts, error) {
	var cons Contracts
	var err error

	cons.PerpetualV1, err = contracts.NewPerpetualV1(perpetualV1Address, backend)
	if err != nil {
		return nil, err
	}
	cons.PerpetualV1Address = perpetualV1Address

	cons.P1Orders, err = contracts.NewP1Orders(p1OrdersAddr, backend)
	if err != nil {
		return nil, err
	}
	cons.P1OrdersAddress = p1OrdersAddr

	cons.txOps = defaultTxOps

	return &cons, nil
}

var emptyAddr common.Address

func combineTxOps(targetOps, defaultOps *bind.TransactOpts) *bind.TransactOpts {
	if targetOps == nil {
		return defaultOps
	}

	if targetOps.From == emptyAddr {
		targetOps.From = defaultOps.From
	}
	if targetOps.Nonce == nil {
		targetOps.Nonce = defaultOps.Nonce
	}
	if targetOps.Signer == nil {
		targetOps.Signer = defaultOps.Signer
	}
	if targetOps.Value == nil {
		targetOps.Value = defaultOps.Value
	}
	if targetOps.GasPrice == nil {
		targetOps.GasPrice = defaultOps.GasPrice
	}
	if targetOps.GasFeeCap == nil {
		targetOps.GasFeeCap = defaultOps.GasFeeCap
	}
	if targetOps.GasTipCap == nil {
		targetOps.GasTipCap = defaultOps.GasTipCap
	}
	if targetOps.GasLimit == 0 {
		targetOps.GasLimit = defaultOps.GasLimit
	}
	if targetOps.Context == nil {
		targetOps.Context = defaultOps.Context
	}
	if !targetOps.NoSend {
		targetOps.NoSend = defaultOps.NoSend
	}
	return targetOps
}
