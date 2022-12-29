package v1

import (
	"encoding/json"
	"fmt"
	"github.com/okex/exchain/libs/tendermint/delta/internal"
	"net/url"
)

const (
	Version = 1

	apiPathGetDelta = "api/v1/delta"
)

type APIErrorResult struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Details []interface{} `json:"details"`
}

type APISuccessResult struct {
	Version int    `json:"version"`
	Result  []byte `json:"result"`
}

func MakeGetDeltaRequestPath(base string, height int64) (string, error) {
	key := internal.GenDeltaKey(height)
	return url.JoinPath(base, apiPathGetDelta, key)
}

func ParseResponse(data []byte, height int64) ([]byte, error, int64) {
	var succ APISuccessResult
	if err := json.Unmarshal(data, &succ); err != nil {
		var apiError APIErrorResult
		if err := json.Unmarshal(data, &apiError); err != nil {
			return nil, fmt.Errorf("unknown response format returned by persist-delta server"), 0
		}
		return nil, fmt.Errorf("persist-delta server error: code: %d; message: %s", apiError.Code, apiError.Message), 0
	}

	if succ.Version != Version {
		return nil, fmt.Errorf("unexpect response verion: current is %d but got %d", Version, succ.Version), 0
	}

	// TODO: most recent height(mrh) seems link useless, delete it.
	return succ.Result, nil, height
}
