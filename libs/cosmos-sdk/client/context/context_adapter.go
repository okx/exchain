package context

import (
	"encoding/json"
	"github.com/gogo/protobuf/proto"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	interfacetypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"gopkg.in/yaml.v2"
	"os"
)

func (ctx CLIContext) WithProxy(cdc *codec.CodecProxy) CLIContext {
	ctx.CodecProy = cdc
	return ctx
}

func (ctx CLIContext) WithInterfaceRegistry(r interfacetypes.InterfaceRegistry) CLIContext {
	ctx.InterfaceRegistry = r
	return ctx
}

func (ctx CLIContext) PrintProto(toPrint proto.Message) error {
	// always serialize JSON initially because proto json can't be directly YAML encoded
	out, err := ctx.Codec.MarshalJSON(toPrint)
	if err != nil {
		return err
	}
	return ctx.printOutput(out)
}
func (ctx CLIContext) printOutput(out []byte) error {
	if ctx.OutputFormat == "text" {
		// handle text format by decoding and re-encoding JSON as YAML
		var j interface{}

		err := json.Unmarshal(out, &j)
		if err != nil {
			return err
		}

		out, err = yaml.Marshal(j)
		if err != nil {
			return err
		}
	}

	writer := ctx.Output
	if writer == nil {
		writer = os.Stdout
	}

	_, err := writer.Write(out)
	if err != nil {
		return err
	}

	if ctx.OutputFormat != "text" {
		// append new-line for formats besides YAML
		_, err = writer.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}

	return nil
}
