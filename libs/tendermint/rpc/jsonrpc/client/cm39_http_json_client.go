package client

import (
	"bytes"
	"encoding/json"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	coretypes "github.com/okex/exchain/libs/tendermint/rpc/core/types"
	"github.com/okex/exchain/libs/tendermint/rpc/jsonrpc/types"
	"io/ioutil"
	"net/http"
)

type decoder interface {
	decode(cdc *codec.Codec, id types.JSONRPCIntID, data []byte, result interface{}) (interface{}, error)
}

type Cm39HttpJSONClientAdapter struct {
	decoders map[string]decoder
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

var (
	_ decoder = (*defaultFuncDecoder)(nil)
)

type defaultFuncDecoder struct {
	secondChance      func() interface{}
	convSecondToFirst func(s interface{}, f interface{})
}

func newDefaultFuncDecoder(secondChance func() interface{}, convSecondToFirst func(s interface{}, f interface{})) *defaultFuncDecoder {
	return &defaultFuncDecoder{secondChance: secondChance, convSecondToFirst: convSecondToFirst}
}

func (d defaultFuncDecoder) decode(cdc *codec.Codec, id types.JSONRPCIntID, data []byte, result interface{}) (ret interface{}, err error) {
	ret, err = unmarshalResponseBytes(cdc, data, id, result)
	if nil == err {
		return
	}
	another := d.secondChance()
	ret2, err2 := unmarshalResponseBytes(cdc, data, id, another)
	if nil != err2 {
		err = errors.Wrap(err, "second error:"+err2.Error())
		return
	}
	d.convSecondToFirst(ret2, result)
	ret = result
	err = nil
	return
}

func (c *Cm39HttpJSONClientAdapter) seal() {
	c.decoders = make(map[string]decoder)
	c.decoders["tx"] = newDefaultFuncDecoder(func() interface{} {
		return new(CM39ResultTx)
	}, func(s interface{}, f interface{}) {
		cm39 := s.(*CM39ResultTx)
		cm4 := f.(*coretypes.ResultTx)
		ConvTCM392CM4(cm39, cm4)
	})
	c.decoders["abci_query"] = newDefaultFuncDecoder(func() interface{} {
		return new(CM39ResultABCIQuery)
	}, func(s interface{}, f interface{}) {
		cm39 := s.(*CM39ResultABCIQuery)
		cm4 := f.(*coretypes.ResultABCIQuery)
		ConvTCM39ResultABCIQuery2CM4(cm39, cm4)
	})
	c.decoders["broadcast_tx_commit"] = newDefaultFuncDecoder(func() interface{} {
		return new(CM39ResultBroadcastTxCommit)
	}, func(s interface{}, f interface{}) {
		cm39 := s.(*CM39ResultBroadcastTxCommit)
		cm4 := f.(*coretypes.ResultBroadcastTxCommit)
		ConvTCM39BroadcastCommitTx2CM4(cm39, cm4)
	})
}

func (c *Cm39HttpJSONClientAdapter) Call(method string, params map[string]interface{}, result interface{}) (ret interface{}, err error) {
	return c.call(method, params, result, c.decoders[method])
}

func (c *Cm39HttpJSONClientAdapter) call(method string, params map[string]interface{}, result interface{}, dec decoder) (interface{}, error) {
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

	if dec != nil {
		return dec.decode(c.cdc, id, responseBytes, result)
	}
	return unmarshalResponseBytes(c.cdc, responseBytes, id, result)
}
