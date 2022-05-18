package consensus

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestState_unmarshalBlock(t *testing.T) {
	type fields struct{}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{"suppress panic", fields{}, assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &State{}
			tt.wantErr(t, cs.unmarshalBlock(), fmt.Sprintf("unmarshalBlock()"))
		})
	}
}
