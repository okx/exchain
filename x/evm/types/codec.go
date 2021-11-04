package types

import (
	"fmt"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/tendermint/go-amino"
)

// ModuleCdc defines the evm module's codec
var ModuleCdc = codec.New()

const (
	MsgEthereumTxName = "ethermint/MsgEthereumTx"
	ChainConfigName   = "ethermint/ChainConfig"
	MsgEthermintName  = "ethermint/MsgEthermint"
	TxDataName        = "ethermint/TxData"

	ManageContractDeploymentWhitelistProposalName = "okexchain/evm/ManageContractDeploymentWhitelistProposal"
	ManageContractBlockedListProposalName         = "okexchain/evm/ManageContractBlockedListProposal"
)

// RegisterCodec registers all the necessary types and interfaces for the
// evm module
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgEthereumTx{}, MsgEthereumTxName, nil)
	cdc.RegisterConcrete(MsgEthermint{}, MsgEthermintName, nil)
	cdc.RegisterConcrete(TxData{}, TxDataName, nil)
	cdc.RegisterConcrete(ChainConfig{}, ChainConfigName, nil)
	cdc.RegisterConcrete(ManageContractDeploymentWhitelistProposal{}, ManageContractDeploymentWhitelistProposalName, nil)
	cdc.RegisterConcrete(ManageContractBlockedListProposal{}, ManageContractBlockedListProposalName, nil)

	cdc.RegisterConcreteUnmarshaller(ChainConfigName, func(c *amino.Codec, bytes []byte) (interface{}, int, error) {
		return UnmarshalChainConfigFromAmino(c, bytes)
	})
	cdc.RegisterConcreteUnmarshaller(MsgEthereumTxName, func(_ *amino.Codec, bytes []byte) (interface{}, int, error) {
		var msg MsgEthereumTx
		err := msg.UnmarshalFromAmino(bytes)
		if err != nil {
			return nil, 0, err
		}
		return &msg, len(bytes), nil
	})
}

func init() {
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}

// UnmarshalChainConfigFromAmino unmarshal a ChainConfig from an amino encoded byte slice
func UnmarshalChainConfigFromAmino(_ *amino.Codec, data []byte) (*ChainConfig, int, error) {
	var dataLen uint64 = 0
	var subData []byte
	var read int
	config := &ChainConfig{}

	for {
		data = data[dataLen:]
		read += int(dataLen)

		if len(data) <= 0 {
			break
		}

		pos, aminoType, err := amino.ParseProtoPosAndTypeMustOneByte(data[0])
		if err != nil {
			return nil, read, err
		}
		data = data[1:]
		read += 1

		if aminoType == amino.Typ3_ByteLength {
			var n int
			dataLen, n, err = amino.DecodeUvarint(data)
			if err != nil {
				return nil, read, err
			}
			data = data[n:]
			read += n
			subData = data[:dataLen]
		}

		switch pos {
		case 1:
			integer, err := sdk.NewIntFromAmino(subData)
			if err != nil {
				return nil, read, err
			}
			config.HomesteadBlock = integer
		case 2:
			integer, err := sdk.NewIntFromAmino(subData)
			if err != nil {
				return nil, read, err
			}
			config.DAOForkBlock = integer
		case 3:
			if data[0] != 0 && data[0] != 1 {
				return nil, read, fmt.Errorf("invalid DAO fork switch")
			}
			config.DAOForkSupport = data[0] == 1
			dataLen = 1
		case 4:
			integer, err := sdk.NewIntFromAmino(subData)
			if err != nil {
				return nil, read, err
			}
			config.EIP150Block = integer
		case 5:
			config.EIP150Hash = string(subData)
		case 6:
			integer, err := sdk.NewIntFromAmino(subData)
			if err != nil {
				return nil, read, err
			}
			config.EIP155Block = integer
		case 7:
			integer, err := sdk.NewIntFromAmino(subData)
			if err != nil {
				return nil, read, err
			}
			config.EIP158Block = integer
		case 8:
			integer, err := sdk.NewIntFromAmino(subData)
			if err != nil {
				return nil, read, err
			}
			config.ByzantiumBlock = integer
		case 9:
			integer, err := sdk.NewIntFromAmino(subData)
			if err != nil {
				return nil, read, err
			}
			config.ConstantinopleBlock = integer
		case 10:
			integer, err := sdk.NewIntFromAmino(subData)
			if err != nil {
				return nil, read, err
			}
			config.PetersburgBlock = integer
		case 11:
			integer, err := sdk.NewIntFromAmino(subData)
			if err != nil {
				return nil, read, err
			}
			config.IstanbulBlock = integer
		case 12:
			integer, err := sdk.NewIntFromAmino(subData)
			if err != nil {
				return nil, read, err
			}
			config.MuirGlacierBlock = integer
		case 13:
			integer, err := sdk.NewIntFromAmino(subData)
			if err != nil {
				return nil, read, err
			}
			config.YoloV2Block = integer
		case 14:
			integer, err := sdk.NewIntFromAmino(subData)
			if err != nil {
				return nil, read, err
			}
			config.EWASMBlock = integer
		default:
			return nil, read, fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return config, read, nil
}
