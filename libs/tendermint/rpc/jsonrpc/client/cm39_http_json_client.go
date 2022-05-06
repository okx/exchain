package client

import (
	"bytes"
	"encoding/json"
	"github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	coretypes "github.com/okex/exchain/libs/tendermint/rpc/core/types"
	"github.com/okex/exchain/libs/tendermint/rpc/jsonrpc/types"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	defaultCM39ErrorHit = func(e error) bool {
		return strings.Contains(e.Error(), "error unmarshalling result")
	}
	defaultCM39ResultTXBefore  = func() interface{} { return new(CM39ResultTx) }
	defaultCM39ResultABCIQuery = func() interface{} { return new(CM39ResultABCIQuery) }
)

type couple struct {
	onHitError func(e error) bool
	before     func() interface{}
	after      func(beforeReturnV interface{}, originV interface{})
}

type Cm39HttpJSONClientAdapter struct {
	planb map[string]*couple
	*Client
}

func NewCm39HttpJSONClient(remote string, client *http.Client) (*Cm39HttpJSONClientAdapter, error) {
	p, err := NewWithHTTPClient(remote, client)
	if nil != err {
		return nil, err
	}
	ret := &Cm39HttpJSONClientAdapter{Client: p}
	ret.seal()
	return ret, nil
}

func (c *Cm39HttpJSONClientAdapter) seal() {
	c.planb = make(map[string]*couple)
	c.planb["tx"] = &couple{
		onHitError: defaultCM39ErrorHit,
		before:     defaultCM39ResultTXBefore,
		after: func(beforeReturnV interface{}, originV interface{}) {
			cm39 := beforeReturnV.(*CM39ResultTx)
			cm4 := originV.(*coretypes.ResultTx)
			ConvTCM392CM4(cm39, cm4)
		},
	}
	c.planb["abci_query"] = &couple{
		onHitError: defaultCM39ErrorHit,
		before:     defaultCM39ResultABCIQuery,
		after: func(beforeReturnV interface{}, originV interface{}) {
			cm39 := beforeReturnV.(*CM39ResultABCIQuery)
			cm4 := originV.(*coretypes.ResultABCIQuery)
			ConvTCM39ResultABCIQuery2CM4(cm39, cm4)
		},
	}
	c.planb["broadcast_tx_commit"] = &couple{
		onHitError: defaultCM39ErrorHit,
		before: func() interface{} {
			return new(CM39ResultBroadcastTxCommit)
		},
		after: func(beforeReturnV interface{}, originV interface{}) {
			cm39 := beforeReturnV.(*CM39ResultBroadcastTxCommit)
			cm4 := originV.(*coretypes.ResultBroadcastTxCommit)
			ConvTCM39BroadcastCommitTx2CM4(cm39, cm4)
		},
	}
}

func (c *Cm39HttpJSONClientAdapter) Call(method string, params map[string]interface{}, result interface{}) (ret interface{}, err error) {
	return c.call(method, params, result, c.planb[method])
}

func (c *Cm39HttpJSONClientAdapter) call(method string, params map[string]interface{}, result interface{}, coup *couple) (interface{}, error) {
	id := c.nextRequestID()

	request, err := types.MapToRequest(c.cdc, id, method, params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode params")
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request")
	}

	requestBuf := bytes.NewBuffer(requestBytes)
	httpRequest, err := http.NewRequest(http.MethodPost, c.address, requestBuf)
	if err != nil {
		return nil, errors.Wrap(err, "Request failed")
	}
	httpRequest.Header.Set("Content-Type", "text/json")
	if c.username != "" || c.password != "" {
		httpRequest.SetBasicAuth(c.username, c.password)
	}
	httpResponse, err := c.client.Do(httpRequest)
	if err != nil {
		return nil, errors.Wrap(err, "Post failed")
	}
	defer httpResponse.Body.Close() // nolint: errcheck

	responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	ret, err := unmarshalResponseBytes(c.cdc, responseBytes, id, result)
	if nil != err && coup != nil && coup.onHitError(err) {
		bef := coup.before()
		if ret2, err2 := unmarshalResponseBytes(c.cdc, responseBytes, id, bef); nil != err2 {
			err = errors.Wrap(err, "second error:"+err2.Error())
		} else {
			coup.after(ret2, result)
			err = nil
		}
	}
	return ret, err
}
