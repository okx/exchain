package iavl

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_fastNodeChangesWithVersion_checkAdditions(t *testing.T) {
	type fields struct {
		mtx      sync.RWMutex
		versions []int64
		fncMap   map[int64]*fastNodeChanges
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"one version has key", fields{
			versions: []int64{1},
			fncMap: map[int64]*fastNodeChanges{
				1: {
					additions: map[string]*FastNode{"key1": NewFastNode([]byte("key1"), []byte("value1"), 1)},
					removals:  map[string]interface{}{"key2": true},
				}}},
			args{key: "key1"},
			true},
		{"one version has no key 1", fields{
			versions: []int64{1},
			fncMap: map[int64]*fastNodeChanges{
				1: {
					additions: map[string]*FastNode{"key1": NewFastNode([]byte("key1"), []byte("value1"), 1)},
					removals:  map[string]interface{}{"key2": true},
				}}},
			args{key: "key2"},
			false},
		{"one version has no key 2", fields{
			versions: []int64{1},
			fncMap: map[int64]*fastNodeChanges{
				1: {
					additions: map[string]*FastNode{"key1": NewFastNode([]byte("key1"), []byte("value1"), 1)},
					removals:  map[string]interface{}{"key2": true},
				}}},
			args{key: "key3"},
			false},
		{"two versions latest version has", fields{
			versions: []int64{1, 2},
			fncMap: map[int64]*fastNodeChanges{
				1: {
					//additions: map[string]*FastNode{"key1": NewFastNode([]byte("key1"), []byte("value1"), 1)},
					removals: map[string]interface{}{"key1": true},
				},
				2: {
					additions: map[string]*FastNode{"key1": NewFastNode([]byte("key1"), []byte("value1"), 1)},
					removals:  map[string]interface{}{"key2": true},
				}}},
			args{key: "key1"},
			true},
		{"two versions latest version remove ", fields{
			versions: []int64{1, 2},
			fncMap: map[int64]*fastNodeChanges{
				1: {
					additions: map[string]*FastNode{"key1": NewFastNode([]byte("key1"), []byte("value1"), 1)},
					//removals: map[string]interface{}{"key1": true},
				},
				2: {
					//additions: map[string]*FastNode{"key1": NewFastNode([]byte("key1"), []byte("value1"), 1)},
					removals: map[string]interface{}{"key1": true},
				}}},
			args{key: "key1"},
			false},
		{"three versions latest version remove ", fields{
			versions: []int64{1, 2, 3},
			fncMap: map[int64]*fastNodeChanges{
				1: {
					//additions: map[string]*FastNode{"key3": NewFastNode([]byte("key1"), []byte("value1"), 1)},
					removals: map[string]interface{}{"key1": true},
				},
				2: {
					additions: map[string]*FastNode{"key1": NewFastNode([]byte("key1"), []byte("value1"), 1)},
					removals:  map[string]interface{}{"key2": true},
				},
				3: {
					additions: map[string]*FastNode{"key0": NewFastNode([]byte("key0"), []byte("value0"), 1)},
					removals:  map[string]interface{}{"key100": true},
				},
			}},
			args{key: "key1"},
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fncv := &fastNodeChangesWithVersion{
				mtx:      tt.fields.mtx,
				versions: tt.fields.versions,
				fncMap:   tt.fields.fncMap,
			}
			assert.Equalf(t, tt.want, fncv.checkAdditions(tt.args.key), "checkAdditions(%v)", tt.args.key)
		})
	}
}

func Test_fastNodeChangesWithVersion_checkRemovals(t *testing.T) {
	type fields struct {
		mtx      sync.RWMutex
		versions []int64
		fncMap   map[int64]*fastNodeChanges
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"one version removals", fields{
			versions: []int64{1},
			fncMap: map[int64]*fastNodeChanges{
				1: {
					additions: map[string]*FastNode{"key1": NewFastNode([]byte("key1"), []byte("value1"), 1)},
					removals:  map[string]interface{}{"key2": true},
				}}},
			args{key: "key2"},
			true},
		{"one version no removals but additions", fields{
			versions: []int64{1},
			fncMap: map[int64]*fastNodeChanges{
				1: {
					additions: map[string]*FastNode{"key1": NewFastNode([]byte("key1"), []byte("value1"), 1)},
					removals:  map[string]interface{}{"key2": true},
				}}},
			args{key: "key1"},
			true},
		{"one version no related", fields{
			versions: []int64{1},
			fncMap: map[int64]*fastNodeChanges{
				1: {
					additions: map[string]*FastNode{"key1": NewFastNode([]byte("key1"), []byte("value1"), 1)},
					removals:  map[string]interface{}{"key2": true},
				}}},
			args{key: "key3"},
			false},
		{"two versions latest version removed", fields{
			versions: []int64{1, 2},
			fncMap: map[int64]*fastNodeChanges{
				1: {
					//additions: map[string]*FastNode{"key1": NewFastNode([]byte("key1"), []byte("value1"), 1)},
					removals: map[string]interface{}{"key1": true},
				},
				2: {
					additions: map[string]*FastNode{"key1": NewFastNode([]byte("key1"), []byte("value1"), 1)},
					removals:  map[string]interface{}{"key2": true},
				}}},
			args{key: "key2"},
			true},
		{"two versions first version removed ", fields{
			versions: []int64{1, 2},
			fncMap: map[int64]*fastNodeChanges{
				1: {
					//additions: map[string]*FastNode{"key1": NewFastNode([]byte("key1"), []byte("value1"), 1)},
					removals: map[string]interface{}{"key1": true},
				},
				2: {
					additions: map[string]*FastNode{"key1": NewFastNode([]byte("key1"), []byte("value1"), 1)},
					//removals: map[string]interface{}{"key1": true},
				}}},
			args{key: "key1"},
			false},
		{"three versions mid version remove ", fields{
			versions: []int64{1, 2, 3},
			fncMap: map[int64]*fastNodeChanges{
				1: {
					//additions: map[string]*FastNode{"key3": NewFastNode([]byte("key1"), []byte("value1"), 1)},
					removals: map[string]interface{}{"key1": true},
				},
				2: {
					additions: map[string]*FastNode{"key1": NewFastNode([]byte("key1"), []byte("value1"), 1)},
					removals:  map[string]interface{}{"key2": true},
				},
				3: {
					additions: map[string]*FastNode{"key0": NewFastNode([]byte("key0"), []byte("value0"), 1)},
					removals:  map[string]interface{}{"key100": true},
				},
			}},
			args{key: "key2"},
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fncv := &fastNodeChangesWithVersion{
				mtx:      tt.fields.mtx,
				versions: tt.fields.versions,
				fncMap:   tt.fields.fncMap,
			}
			assert.Equalf(t, tt.want, fncv.checkRemovals(tt.args.key), "checkRemovals(%v)", tt.args.key)
		})
	}
}
