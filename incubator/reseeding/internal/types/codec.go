package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSeed{}, "cosmos-sdk/MsgSeed", nil)
}

// ModuleCdc generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	codec.RegisterCrypto(ModuleCdc)
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
