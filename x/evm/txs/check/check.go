package check

import (
	"fmt"

	"github.com/okex/exchain/app/config"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/txs/base"
)

type Tx struct {
	*base.Tx
}

func NewTx(config base.Config) *Tx {
	return &Tx{
		Tx: base.NewTx(config),
	}
}

// AnalyzeStart check Tx do not analyze start
func (t *Tx) AnalyzeStart(tag string) {}

// AnalyzeStop check Tx do not analyze stop
func (t *Tx) AnalyzeStop(tag string) {}

func (tx *Tx) DecorateResult(inResult *base.Result, inErr error) (result *sdk.Result, err error) {
	if inErr != nil {
		return nil, inErr
	}

	if tx.Ctx.EstimateGas() {
		gasEstimated := tx.Ctx.GasMeter().GasConsumed()
		maxGasLimitPerTx := tx.Keeper.GetParams(tx.Ctx).MaxGasLimitPerTx
		if gasEstimated > maxGasLimitPerTx {
			return nil, sdk.ErrOutOfGas(fmt.Sprintf("gas estimated %v greater than system's max gas limit per tx %v", gasEstimated, maxGasLimitPerTx))
		}

		// adjustment estimateGas
		gasBuffer := gasEstimated / 100 * config.GetOecConfig().GetGasLimitBuffer()
		gas := gasEstimated + gasBuffer
		if gas > maxGasLimitPerTx {
			tx.Ctx.GasMeter().ConsumeGas(maxGasLimitPerTx-gasEstimated, "estimate add buffer")
		}
	}

	return inResult.ExecResult.Result, inErr
}
