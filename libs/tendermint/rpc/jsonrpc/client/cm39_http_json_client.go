package client

import (
	coretypes "github.com/okex/exchain/libs/tendermint/rpc/core/types"
	"net/http"
	"strings"
)

var (
	defaultCM39ErrorHit = func(e error) bool {
		return strings.Contains(e.Error(), "error unmarshalling result")
	}
	defaultCM39ResultTXBefore  = func() interface{} { return new(coretypes.CM39ResultTx) }
	defaultCM39ResultABCIQuery = func() interface{} { return new(coretypes.CM39ResultABCIQuery) }
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

func (c *Cm39HttpJSONClientAdapter) seal() {
	c.planb = make(map[string]*couple)
	c.planb["tx"] = &couple{
		onHitError: defaultCM39ErrorHit,
		before:     defaultCM39ResultTXBefore,
		after: func(beforeReturnV interface{}, originV interface{}) {
			cm39 := beforeReturnV.(*coretypes.CM39ResultTx)
			cm4 := originV.(*coretypes.ResultTx)
			ConvCM39ToCM4(cm39, cm4)
		},
	}
	c.planb["abci_query"] = &couple{
		onHitError: defaultCM39ErrorHit,
		before:     defaultCM39ResultABCIQuery,
		after: func(beforeReturnV interface{}, originV interface{}) {
			cm39 := beforeReturnV.(*coretypes.CM39ResultABCIQuery)
			cm4 := originV.(*coretypes.ResultABCIQuery)
			ConvTCM39ResultABCIQuery2CM4(cm39, cm4)
		},
	}
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

func (c *Cm39HttpJSONClientAdapter) Call(method string, params map[string]interface{}, result interface{}) (ret interface{}, err error) {
	cp := c.planb[method]

	if cp != nil {
		defer func() {
			if nil == err || !cp.onHitError(err) {
				return
			}
			res := cp.before()
			v, pbErr := c.Client.Call(method, params, res)
			if nil != pbErr {
				return
			}
			cp.after(v, result)
			ret = result
			err = nil
		}()
	}

	ret, err = c.Client.Call(method, params, result)
	if nil != err {
		return nil, err
	}
	return ret, nil
}
