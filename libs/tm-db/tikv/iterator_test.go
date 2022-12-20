package tikv

import (
	"testing"

	"github.com/tikv/client-go/v2/rawkv"
)

func TestIterator_Valid(t *testing.T) {
	type fields struct {
		client    *rawkv.Client
		curKey    []byte
		curValue  []byte
		start     []byte
		end       []byte
		isReverse bool
		err       error
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"err is nil", fields{client: new(rawkv.Client)}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Iterator{
				client:    tt.fields.client,
				curKey:    tt.fields.curKey,
				curValue:  tt.fields.curValue,
				start:     tt.fields.start,
				end:       tt.fields.end,
				isReverse: tt.fields.isReverse,
				err:       tt.fields.err,
			}
			if got := i.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
