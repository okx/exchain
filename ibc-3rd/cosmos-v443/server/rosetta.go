package server

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/okex/exchain/ibc-3rd/cosmos-v443/codec"
	codectypes "github.com/okex/exchain/ibc-3rd/cosmos-v443/codec/types"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/server/rosetta"
)

// RosettaCommand builds the rosetta root command given
// a protocol buffers serializer/deserializer
func RosettaCommand(ir codectypes.InterfaceRegistry, cdc codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rosetta",
		Short: "spin up a rosetta server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Println("WARNING: The Rosetta server is still a beta feature. Please do not use it in production.")

			conf, err := rosetta.FromFlags(cmd.Flags())
			if err != nil {
				return err
			}

			protoCodec, ok := cdc.(*codec.ProtoCodec)
			if !ok {
				return fmt.Errorf("exoected *codec.ProtoMarshaler, got: %T", cdc)
			}
			conf.WithCodec(ir, protoCodec)

			rosettaSrv, err := rosetta.ServerFromConfig(conf)
			if err != nil {
				return err
			}
			return rosettaSrv.Start()
		},
	}
	rosetta.SetFlags(cmd.Flags())

	return cmd
}
