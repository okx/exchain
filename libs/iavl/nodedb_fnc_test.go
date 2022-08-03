package iavl

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_fastNodeChangesWithVersion_expand(t *testing.T) {
	type fields struct {
		mtx      sync.RWMutex
		versions []int64
		fncMap   map[int64]*fastNodeChanges
	}
	type args struct {
		changes *fastNodeChanges
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *fastNodeChanges
	}{
		{"add3 del12", fields{
			versions: []int64{1},
			fncMap: map[int64]*fastNodeChanges{
				1: {
					additions: map[string]*FastNode{"key1": NewFastNode([]byte("key1"), []byte("value1"), 1)},
					removals:  map[string]interface{}{"key2": true},
				}}},
			args{
				changes: &fastNodeChanges{
					additions: map[string]*FastNode{"key3": NewFastNode([]byte("key3"), []byte("value3"), 2)},
					removals:  map[string]interface{}{"key1": true},
				}},
			&fastNodeChanges{
				additions: map[string]*FastNode{"key3": NewFastNode([]byte("key3"), []byte("value3"), 2)},
				removals:  map[string]interface{}{"key1": true, "key2": true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fncv := &fastNodeChangesWithVersion{
				mtx:      tt.fields.mtx,
				versions: tt.fields.versions,
				fncMap:   tt.fields.fncMap,
			}
			assert.Equalf(t, tt.want, fncv.expand(tt.args.changes), "expand(%v)", tt.args.changes)
		})
	}
}
