package legacytx

import (
	"github.com/line/lbm-sdk/codec"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(StdTx{}, "lbm-sdk/StdTx", nil)
}
