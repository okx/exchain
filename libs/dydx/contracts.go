package dydx

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/libs/dydx/contracts"
)

type Contracts struct {
	P1Orders          *contracts.P1Orders
	PerpetualV1       *contracts.PerpetualV1
	PerpetualV1Oracel *contracts.IP1Oracle
	P1MakerOracle     *contracts.P1MakerOracle

	Addresses *ContractsAddressConfig
	txOps     *bind.TransactOpts
	backend   bind.ContractBackend
}

type ContractsAddressConfig struct {
	PerpetualProxy           common.Address
	PerpetualV1              common.Address
	P1FundingOracle          common.Address
	P1InverseFundingOracle   common.Address
	P1ChainlinkOracle        common.Address
	P1MakerOracle            common.Address
	P1MirrorOracle           common.Address
	P1OracleInverter         common.Address
	P1Orders                 common.Address
	P1InverseOrders          common.Address
	P1Deleveraging           common.Address
	P1Liquidation            common.Address
	P1CurrencyConverterProxy common.Address
	P1LiquidatorProxy        common.Address
	P1SoloBridgeProxy        common.Address
	P1WethProxy              common.Address
	ERC20                    common.Address
	WETH                     common.Address
}

func NewContracts(
	config *ContractsAddressConfig,
	defaultTxOps *bind.TransactOpts,
	backend bind.ContractBackend,
) (*Contracts, error) {
	var cons Contracts
	var err error

	cons.PerpetualV1, err = contracts.NewPerpetualV1(config.PerpetualV1, backend)
	if err != nil {
		return nil, err
	}

	cons.P1Orders, err = contracts.NewP1Orders(config.P1Orders, backend)
	if err != nil {
		return nil, err
	}

	cons.P1MakerOracle, err = contracts.NewP1MakerOracle(config.P1MakerOracle, backend)
	if err != nil {
		return nil, err
	}

	cons.txOps = defaultTxOps
	cons.backend = backend
	cons.Addresses = config

	return &cons, nil
}

func (cc *Contracts) GetPerpetualV1OraclePrice() (*big.Int, error) {
	if cc.PerpetualV1Oracel == nil {
		perpetualV1OracleAddress, err := cc.PerpetualV1.GetOracleContract(nil)
		if err != nil {
			return nil, err
		}
		cc.PerpetualV1Oracel, err = contracts.NewIP1Oracle(perpetualV1OracleAddress, cc.backend)
		if err != nil {
			return nil, err
		}
	}

	return cc.PerpetualV1Oracel.GetPrice(&bind.CallOpts{
		From: cc.Addresses.PerpetualV1,
	})
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
