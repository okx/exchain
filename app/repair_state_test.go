package app

import (
	"fmt"
	"regexp"
	"testing"
)

func Test_getRetryHeight(t *testing.T) {

	errReg = regexp.MustCompile("[0-9]+")
	tests := []struct {
		name  string
		err   error
		want  bool
		want1 int64
	}{
		{
			name:  "err",
			err:   fmt.Errorf("failed to load staking Store: wanted to load target 12163103 but only found up to 12163102"),
			want:  true,
			want1: 12163102,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getRetryHeight(tt.err)
			if got != tt.want {
				t.Errorf("getRetryHeight() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getRetryHeight() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
