package types

import (
	"github.com/line/lbm-sdk/codec"
	"github.com/line/lbm-sdk/codec/types"
	cryptocodec "github.com/line/lbm-sdk/crypto/codec"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/msgservice"
	"github.com/line/lbm-sdk/x/bank/exported"
)

// RegisterLegacyAminoCodec registers the necessary x/bank interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*exported.SupplyI)(nil), nil)
	cdc.RegisterConcrete(&Supply{}, "lbm-sdk/Supply", nil)
	cdc.RegisterConcrete(&MsgSend{}, "lbm-sdk/MsgSend", nil)
	cdc.RegisterConcrete(&MsgMultiSend{}, "lbm-sdk/MsgMultiSend", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSend{},
		&MsgMultiSend{},
	)

	registry.RegisterInterface(
		"lbm.bank.v1.SupplyI",
		(*exported.SupplyI)(nil),
		&Supply{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/bank module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/staking and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
