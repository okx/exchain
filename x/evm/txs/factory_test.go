package txs

import (
	"reflect"
	"testing"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/txs/base"
	"github.com/okex/exchain/x/evm/txs/check"
	"github.com/okex/exchain/x/evm/txs/deliver"
	"github.com/okex/exchain/x/evm/txs/tracetxlog"
)

func Test_factory_CreateTx(t *testing.T) {
	type fields struct {
		Config base.Config
	}
	traceTxConfig := base.Config{Keeper: &base.Keeper{}, Ctx: *(&sdk.Context{}).SetRunTxMode(sdk.RunTxModeTrace).SetNeedTraceTxLog(true)}
	checkTxConfig := base.Config{Keeper: &base.Keeper{}, Ctx: *(&sdk.Context{}).SetRunTxMode(sdk.RunTxModeCheck)}
	deliverTxConfig := base.Config{Keeper: &base.Keeper{}, Ctx: *(&sdk.Context{}).SetRunTxMode(sdk.RunTxModeDeliver)}
	tests := []struct {
		name    string
		fields  fields
		want    Tx
		wantErr bool
	}{
		{"1. factory keeper is nil", fields{Config: base.Config{Keeper: nil, Ctx: sdk.Context{}}}, nil, true},
		{"2. create trace tx log", fields{Config: traceTxConfig}, tracetxlog.NewTx(traceTxConfig), false},
		{"3. create check tx log", fields{Config: checkTxConfig}, check.NewTx(checkTxConfig), false},
		{"4. create deliver(default) tx log", fields{Config: deliverTxConfig}, deliver.NewTx(deliverTxConfig), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := &factory{
				Config: tt.fields.Config,
			}
			got, err := factory.CreateTx()
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateTx() got = %v, want %v", got, tt.want)
			}
		})
	}
}
