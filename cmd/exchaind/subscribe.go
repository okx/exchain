package main

import (
	"fmt"
	"github.com/okex/exchain/app/logevents"
	"github.com/spf13/cobra"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
)

func subscribeCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscribe",
		Short: "subscribe oec logs from kafka",
	}

	cmd.AddCommand(
		subscriber(),
	)

	return cmd

}

func subscriber() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs [urls] [topic]",
		Short: "logs urls topic",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s, %s\n", args[0], args[1])
			subscriber := logevents.NewSubscriber()
			subscriber.Init(args[0], args[1])
			subscriber.Run()
		},
	}
	return cmd
}


