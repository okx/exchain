package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type CaseRewriteOperator func(c *Instance)

type Instance struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Workspace   string    `json:"workspace"`
	Network     Network   `json:"network"`
	Nodes       []Node    `json:"nodes"`
	Accounts    []Account `json:"accounts"`
}

func NewInstance() *Instance {
	return &Instance{
		Name:        "",
		Description: "",
		Network:     NewNetwork(),
		Nodes:       make([]Node, 0),
		Accounts:    make([]Account, 0),
	}
}

// Rewrite case to file with rewrite operators
func (c *Instance) Rewrite(file *os.File, operators ...CaseRewriteOperator) error {
	for _, operator := range operators {
		operator(c)
	}
	document, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}
	_, err = file.Write(document)
	return err
}

func (c *Instance) IsSeedMode(name string) bool {
	for _, v := range c.Network.Seeds {
		if v == name {
			return true
		}
	}
	return false
}

func (c *Instance) IsValidatorMode(name string) bool {
	for _, v := range c.Network.Validators {
		if v == name {
			return true
		}
	}
	return false
}

func (c *Instance) IsRpcMode(name string) bool {
	for _, v := range c.Network.Rpcs {
		if v == name {
			return true
		}
	}
	return false
}

func (c *Instance) Unique() bool {
	m := make(map[string]interface{}, 0)
	for _, node := range c.Nodes {
		if _, ok := m[node.Name]; ok {
			return false
		}
		m[node.Name] = new(interface{})
	}
	return true
}

type Network struct {
	P2PBase    int      `json:"p2p"`        // p2p network base port 26656
	RpcBase    int      `json:"rpc"`        // consensus system port 26657
	RestBase   int      `json:"rest"`       // rest port 8545
	ChainID    string   `json:"chain_id"`   // chain id
	IP         string   `json:"ip"`         // ip of this machine
	Seeds      []string `json:"seeds"`      // seeds node
	Validators []string `json:"validators"` // validator node
	Rpcs       []string `json:"rpcs"`       // rpc node
	Whitelist  []string `json:"whitelist"`  // whitelist node key list
}

func NewNetwork() Network {
	return Network{
		P2PBase:    26656,
		RpcBase:    26657,
		RestBase:   8545,
		ChainID:    "exchain-67",
		IP:         "127.0.0.1",
		Seeds:      []string{},
		Validators: []string{},
		Rpcs:       []string{},
		Whitelist:  []string{},
	}
}

type Node struct {
	Name       string   `json:"name"`
	Branch     string   `json:"branch"`
	Executable string   `json:"executable"`
	ServerHome string   `json:"server_home"`
	Wtx        bool     `json:"wtx"`
	White      bool     `json:"white"`
	Flags      []string `json:"flags"`
}

// modify the p2p seeds modify
func (node Node) WithP2PSeeds(seeds, ip, port string) Node {
	index := -1
	for i, flag := range node.Flags {
		if flag == "--p2p.seeds" {
			index = i
		}
	}
	if index == -1 {
		return node
	}
	position := index + 1
	if position == len(node.Flags) {
		return node
	}
	node.Flags[position] = fmt.Sprintf("%s@%s:%s", seeds, ip, port)
	return node
}

type Account struct {
	Mnemonic string `json:"mnemonic"`
	Balance  string `json:"balance"`
}
