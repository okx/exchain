package dydx

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/okex/exchain/libs/dydx/contracts"
)

func BnToBytes32(value *big.Int) string {
	str := value.Text(16)
	needPad := 64 - len(str)
	var sb strings.Builder
	if needPad > 0 {
		sb.Grow(2 + 64)
	} else {
		sb.Grow(2 + len(str))
	}
	sb.WriteString("0x")
	for needPad > 0 {
		sb.WriteByte('0')
		needPad--
	}
	sb.WriteString(str)
	return sb.String()
}

func StripHexPrefix(value string) string {
	if strings.HasPrefix(value, "0x") {
		return value[2:]
	}
	return value
}

func CombineHexString(args ...string) string {
	var sb strings.Builder
	size := 2
	for _, arg := range args {
		size += len(arg)
	}
	sb.Grow(size)
	sb.WriteString("0x")
	for _, arg := range args {
		sb.WriteString(StripHexPrefix(arg))
	}
	return sb.String()
}

func Bytes32ToBalance(bz *[32]byte) contracts.P1TypesBalance {
	var balance contracts.P1TypesBalance
	balance.Position = new(big.Int).SetBytes(bz[17:32])
	balance.Margin = new(big.Int).SetBytes(bz[1:16])
	balance.PositionIsPositive = bz[16]&0x01 == 0x01
	balance.MarginIsPositive = bz[0]&0x01 == 0x01
	return balance
}

func Bytes32ToIndex(bz *[32]byte) contracts.P1TypesIndex {
	var index contracts.P1TypesIndex
	index.Value = new(big.Int).SetBytes(bz[16:32])
	index.IsPositive = bz[15]&0x01 == 0x01
	index.Timestamp = uint32(new(big.Int).SetBytes(bz[11:15]).Uint64())
	return index
}

type P1TypesBalanceStringer contracts.P1TypesBalance

func (b P1TypesBalanceStringer) String() string {
	margin := b.Margin.String()
	position := b.Position.String()
	if !b.PositionIsPositive {
		position = "-" + position
	}
	if !b.MarginIsPositive {
		margin = "-" + margin
	}
	return fmt.Sprintf("margin: %s, position: %s", margin, position)
}
