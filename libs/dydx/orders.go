package dydx

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/okex/exchain/libs/dydx/contracts"
)

type Orders struct {
	contracts        *Contracts
	eip712DomainName string
	orders           *contracts.P1Orders
}

func NewOrders(contracts *Contracts) *Orders {
	var orders Orders
	orders.eip712DomainName = DEFAULT_EIP712_DOMAIN_NAME
	orders.contracts = contracts
	orders.orders = contracts.P1Orders
	return &orders
}

func (orders *Orders) Address() common.Address {
	return orders.contracts.Addresses.P1Orders
}

func (orders *Orders) ApproveOrder(order *Order, ops *bind.TransactOpts) (*types.Transaction, error) {
	solOrder := order.ToSolidity()
	return orders.orders.ApproveOrder(combineTxOps(ops, orders.contracts.txOps), *solOrder)
}
